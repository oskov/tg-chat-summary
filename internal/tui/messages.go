package tui

import (
	"fmt"
	"strings"

	bubbletea "github.com/charmbracelet/bubbletea"
	"github.com/oskov/tg-chat-summary/internal/summarizer"
	"github.com/oskov/tg-chat-summary/internal/tgclient"
)

type errMsg struct {
	err error
}

// Command to fetch chats
type fetchChatsMsg struct {
	chats []tgclient.Chat
}

func fetchChatsCmd(client *tgclient.Client) bubbletea.Cmd {
	return func() bubbletea.Msg {
		chats, err := client.GetChats()
		if err != nil {
			return errMsg{err: err}
		}
		return fetchChatsMsg{chats: chats}
	}
}

type fetchChatHistorySummaryMsg struct {
	summary string
}

func fetchChatHistorySummaryCmd(client *tgclient.Client, summr *summarizer.Summarizer, chatId int64) bubbletea.Cmd {
	fetchCmd := func() bubbletea.Msg {
		msgs, err := client.ReadHistory(chatId, 100) // Read the last 100 messages
		if err != nil {
			return errMsg{err: err}
		}

		var messages []string

		for i := len(msgs) - 1; i >= 0; i-- {
			dateStr := msgs[i].Date.Format("2006-01-02 15:04:05")
			messages = append(messages, fmt.Sprintf("(%s) %s: `%s`", dateStr, msgs[i].Sender, msgs[i].Text))
		}

		resp, err := summr.SendPrompt(strings.Join(messages, "\n"))
		if err != nil {
			return errMsg{err: err}
		}

		return fetchChatHistorySummaryMsg{summary: resp}
	}
	return bubbletea.Sequence(loadingCmd(true), fetchCmd, loadingCmd(false))
}

type loadingMsg struct {
	loading bool
}

func loadingCmd(loading bool) bubbletea.Cmd {
	return func() bubbletea.Msg {
		return loadingMsg{loading: loading}
	}
}
