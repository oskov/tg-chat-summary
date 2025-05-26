package tgclient

import (
	"math"
	"path/filepath"
	"time"

	"github.com/zelenin/go-tdlib/client"
)

type Client struct {
	tdlibClient *client.Client
}

// NewClient initializes a new TDLib client and returns it, or an error if initialization fails.
func NewClient(apiId int, apiHash string) (*Client, error) {
	tdlibParameters := &client.SetTdlibParametersRequest{
		UseTestDc:           false,
		DatabaseDirectory:   filepath.Join(".tdlib", "database"),
		FilesDirectory:      filepath.Join(".tdlib", "files"),
		UseFileDatabase:     true,
		UseChatInfoDatabase: true,
		UseMessageDatabase:  true,
		UseSecretChats:      false,
		ApiId:               int32(apiId),
		ApiHash:             apiHash,
		SystemLanguageCode:  "en",
		DeviceModel:         "Server",
		SystemVersion:       "1.0.0",
		ApplicationVersion:  "1.0.0",
	}

	// Initialize the authorizer
	authorizer := client.ClientAuthorizer(tdlibParameters)
	go client.CliInteractor(authorizer)

	// Set log verbosity level
	_, err := client.SetLogVerbosityLevel(&client.SetLogVerbosityLevelRequest{
		NewVerbosityLevel: 1,
	})
	if err != nil {
		return nil, err
	}

	// Create the TDLib client
	tdlibClient, err := client.NewClient(authorizer)
	if err != nil {
		return nil, err
	}

	return &Client{tdlibClient: tdlibClient}, nil
}

type VersionInfo struct {
	Version string
	Commit  string
}

func (c *Client) GetVersionInfo() (VersionInfo, error) {
	versionOption, err := client.GetOption(&client.GetOptionRequest{
		Name: "version",
	})
	if err != nil {
		return VersionInfo{}, err
	}

	commitOption, err := client.GetOption(&client.GetOptionRequest{
		Name: "commit_hash",
	})
	if err != nil {
		return VersionInfo{}, err
	}

	versionInfo := VersionInfo{
		Version: versionOption.(*client.OptionValueString).Value,
		Commit:  commitOption.(*client.OptionValueString).Value,
	}

	return versionInfo, nil
}

type MeInfo struct {
	FirstName string
	LastName  string
}

func (c *Client) GetMe() (MeInfo, error) {
	me, err := c.tdlibClient.GetMe()
	if err != nil {
		return MeInfo{}, err
	}

	return MeInfo{
		FirstName: me.FirstName,
		LastName:  me.LastName,
	}, nil
}

type ChatType string

const (
	ChatTypePrivate    ChatType = client.TypeChatTypePrivate
	ChatTypeGroup      ChatType = client.TypeChatTypeBasicGroup
	ChatTypeSuperGroup ChatType = client.TypeChatTypeSupergroup
	ChatTypeSecret     ChatType = client.TypeChatTypeSecret
)

type Chat struct {
	Id    int64
	Title string
	Type  ChatType
}

func (c *Client) GetChats() ([]Chat, error) {
	chats, err := c.tdlibClient.GetChats(&client.GetChatsRequest{
		Limit: math.MaxInt32,
	})
	if err != nil {
		return nil, err
	}

	var result []Chat
	for _, chatId := range chats.ChatIds {
		chat, err := c.tdlibClient.GetChat(&client.GetChatRequest{
			ChatId: chatId,
		})
		if err != nil {
			continue
		}
		chatType := ChatType(chat.Type.ChatTypeType())
		if chatType == ChatTypeSecret {
			continue
		}
		result = append(result, Chat{
			Id:    chat.Id,
			Title: chat.Title,
			Type:  chatType,
		})
	}
	return result, nil
}

type Message struct {
	Sender string
	Text   string
	Date   time.Time
}

type MessageSenderType string

const (
	MessageSenderTypeUser MessageSenderType = client.TypeMessageSenderUser
	MessageSenderTypeChat MessageSenderType = client.TypeMessageSenderChat
)

func (c *Client) ReadHistory(chatId int64, limit int) ([]Message, error) {
	var messages []Message
	var fromMessageId int64 = 0
	const maxLimit int32 = 100
	remaining := limit

	for remaining > 0 {
		batchLimit := int32(math.Min(float64(remaining), float64(maxLimit)))

		history, err := c.tdlibClient.GetChatHistory(&client.GetChatHistoryRequest{
			ChatId:        chatId,
			FromMessageId: fromMessageId,
			Offset:        0,
			Limit:         batchLimit,
			OnlyLocal:     false,
		})
		if err != nil {
			return nil, err
		}

		if len(history.Messages) == 0 {
			break
		}

		for _, message := range history.Messages {
			sender := "unknown"
			if message.SenderId.MessageSenderType() == string(MessageSenderTypeUser) {
				userId := message.SenderId.(*client.MessageSenderUser).UserId
				if userId != 0 {
					user, err := c.tdlibClient.GetUser(&client.GetUserRequest{
						UserId: userId,
					})
					if err == nil {
						sender = user.FirstName
						if user.LastName != "" {
							sender += " " + user.LastName
						}
					}
				}
			} else {
				sender = "SystemChat"
			}
			if message.Content.MessageContentType() == client.TypeMessageText {
				text := message.Content.(*client.MessageText).Text.Text
				messages = append(messages, Message{
					Sender: sender,
					Text:   text,
					Date:   time.Unix(int64(message.Date), 0),
				})
			}
		}

		fromMessageId = history.Messages[len(history.Messages)-1].Id
		remaining -= len(history.Messages)
	}

	return messages, nil
}

func (c *Client) Close() error {
	_, err := c.tdlibClient.Close()
	return err
}
