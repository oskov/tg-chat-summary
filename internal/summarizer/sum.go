package summarizer

import (
	"fmt"

	"github.com/oskov/tg-chat-summary/internal/ollama"
)

type Summarizer struct {
	client *ollama.OllamaClient
}

// NewSummarizer creates a new Summarizer instance with the given OllamaClient.
func NewSummarizer(client *ollama.OllamaClient) *Summarizer {
	return &Summarizer{
		client: client,
	}
}

const template = `Based on a given chat history, write a story that summarizes the chat.
Chat history:
%s`

// SendPrompt sends a prompt to the Ollama server and returns the response or an error.
func (s *Summarizer) SendPrompt(prompt string) (string, error) {
	payload := fmt.Sprintf(template, prompt)
	// Use the OllamaClient to send the request
	return s.client.Generate("deepseek-r1:14b", payload)
}
