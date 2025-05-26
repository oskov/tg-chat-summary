package tui

import (
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	bubbletea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
	"github.com/oskov/tg-chat-summary/internal/summarizer"
	"github.com/oskov/tg-chat-summary/internal/tgclient"
)

type UI struct {
	client *tgclient.Client
	summr  *summarizer.Summarizer
}

func NewUI(client *tgclient.Client, summr *summarizer.Summarizer) *UI {
	return &UI{
		client: client,
		summr:  summr,
	}
}

type bubbleteaModel interface {
	bubbletea.Model
	Help() []string
}

type mainModel struct {
	err     error
	loading bool

	chatList    *chatListModel
	chatSummary *chatSummaryModel
	spinner     spinner.Model

	currentModel bubbleteaModel
}

func (m *mainModel) Init() bubbletea.Cmd {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	m.spinner = s

	return m.chatList.Init()
}
func (m *mainModel) Update(msg bubbletea.Msg) (bubbletea.Model, bubbletea.Cmd) {
	switch msg := msg.(type) {
	case loadingMsg:
		m.loading = msg.loading
		if msg.loading {
			return m, m.spinner.Tick
		}
	case fetchChatsMsg:
		m.currentModel = m.chatList
		m.currentModel.Update(msg)
		return m, nil
	case fetchChatHistorySummaryMsg:
		m.currentModel = m.chatSummary
		m.currentModel.Update(msg)
		return m, nil
	case errMsg:
		m.err = msg.err
		m.currentModel = nil
		return m, nil
	case spinner.TickMsg:
		var cmd bubbletea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	case bubbletea.KeyMsg:
		switch msg.String() {
		case "esc": // Return to the main menu
			m.currentModel = m.chatList
			return m, nil
		case "q": // Quit the program
			return m, bubbletea.Quit
		case "ctrl+c": // Quit the program on Ctrl+C
			return m, bubbletea.Quit
		}
	}
	if m.currentModel != nil {
		_, cmd := m.currentModel.Update(msg)
		return m, cmd
	}
	if m.loading {
		return m, m.spinner.Tick
	}
	return m, nil
}

func (m *mainModel) View() string {
	mainView := "Welcome to the Telegram Chat Summary Tool!"
	if m.err != nil {
		mainView = fmt.Sprintf("Error: %v\nPress 'q' to quit.", m.err)
	} else if m.loading {
		mainView = fmt.Sprintf("Loading... %s\n", m.spinner.View())
	} else if m.currentModel != nil {
		mainView = m.currentModel.View()
	}

	mainView += "\n\n" + m.Help()
	return mainView
}

func (m *mainModel) Help() string {
	var help []string
	if m.currentModel != nil {
		help = m.currentModel.Help()
	}
	help = append(help, "Press ESC to return to the main menu.")
	help = append(help, "Press 'q' to quit.")

	for i := range help {
		help[i] = fmt.Sprintf("  :: %s", help[i])
	}

	return strings.Join(help, "\n")
}

type chatListModel struct {
	client *tgclient.Client
	summr  *summarizer.Summarizer

	chats  []tgclient.Chat
	cursor int // Index of the currently selected chat
}

func (m *chatListModel) Init() bubbletea.Cmd {
	return fetchChatsCmd(m.client)
}

func (m *chatListModel) Update(msg bubbletea.Msg) (bubbletea.Model, bubbletea.Cmd) {
	switch msg := msg.(type) {
	case fetchChatsMsg:
		m.chats = msg.chats
		return m, nil
	case bubbletea.KeyMsg:
		switch msg.String() {
		case "up": // Move cursor up
			if m.cursor > 0 {
				m.cursor--
			}
			return m, nil
		case "down": // Move cursor down
			if m.cursor < len(m.chats)-1 {
				m.cursor++
			}
			return m, nil
		case "enter":
			if m.cursor >= 0 && m.cursor < len(m.chats) {
				return m, fetchChatHistorySummaryCmd(m.client, m.summr, m.chats[m.cursor].Id)
			}
		}
	}
	return m, nil
}

func (m *chatListModel) View() string {
	if len(m.chats) == 0 {
		return "Loading chats...\nPress 'q' to quit."
	}

	const maxVisible = 10 // Maximum number of chats to display at once
	start := m.cursor - maxVisible/2
	if start < 0 {
		start = 0
	}
	end := start + maxVisible
	if end > len(m.chats) {
		end = len(m.chats)
		start = end - maxVisible
		if start < 0 {
			start = 0
		}
	}

	view := "Chats:\n"
	if start > 0 {
		view += "...\n" // Indicate there are more chats above
	}
	for i := start; i < end; i++ {
		cursor := " " // No cursor by default
		if i == m.cursor {
			cursor = ">" // Highlight the selected chat
		}
		view += fmt.Sprintf("%s %s (%s)\n", cursor, m.chats[i].Title, m.chats[i].Type)
	}
	if end < len(m.chats) {
		view += "...\n" // Indicate there are more chats below
	}
	return view
}

func (m *chatListModel) Help() []string {
	return []string{
		"Use ↑/↓ to navigate through the chat list.",
		"Press Enter to select a chat and view it's summary.",
	}
}

type chatSummaryModel struct {
	text string
}

func (m *chatSummaryModel) Init() bubbletea.Cmd {
	m.text = "No summary available."
	return nil
}

func (m *chatSummaryModel) Update(msg bubbletea.Msg) (bubbletea.Model, bubbletea.Cmd) {
	switch msg := msg.(type) {
	case fetchChatHistorySummaryMsg:
		m.text = msg.summary
		return m, nil
	}
	return m, nil
}

func (m *chatSummaryModel) View() string {
	sumText := wordwrap.String(m.text, 80)
	return fmt.Sprintf("Summary:\n%s", sumText)
}

func (m *chatSummaryModel) Help() []string {
	return []string{}
}

// Run starts the Bubble Tea program
func (ui *UI) Run() {
	// Initialize the main model
	mainModel := &mainModel{
		chatList:    &chatListModel{client: ui.client, summr: ui.summr},
		chatSummary: &chatSummaryModel{},
	}
	p := bubbletea.NewProgram(mainModel)
	if _, err := p.Run(); err != nil {
		log.Fatalf("Error running UI: %v", err)
	}
}
