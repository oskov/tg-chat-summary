package main

import (
	"log"
	"os"
	"strconv"

	"github.com/oskov/tg-chat-summary/internal/ollama"
	"github.com/oskov/tg-chat-summary/internal/summarizer"
	"github.com/oskov/tg-chat-summary/internal/tgclient"
	"github.com/oskov/tg-chat-summary/internal/tui"
)

func main() {
	// Read environment variables
	apiId, err := strconv.Atoi(os.Getenv("TG_API_ID"))
	if err != nil {
		log.Fatalf("Invalid TG_API_ID: %v", err)
	}
	apiHash := os.Getenv("TG_API_HASH")
	ollamaHost := os.Getenv("OLLAMA_HOST")

	if apiHash == "" || ollamaHost == "" {
		log.Fatalf("Environment variables TG_API_ID, TG_API_HASH, and OLLAMA_HOST must be set")
	}

	// Initialize the TDLib client
	client, err := tgclient.NewClient(apiId, apiHash)
	if err != nil {
		log.Fatalf("Failed to create TDLib client: %v", err)
	}
	defer client.Close()

	ollamaClient := ollama.NewOllamaClient(ollamaHost)
	summarizer := summarizer.NewSummarizer(ollamaClient)
	gui := tui.NewUI(client, summarizer)
	gui.Run()

	if err := client.Close(); err != nil {
		log.Printf("Error closing client: %v", err)
	} else {
		log.Println("Client closed successfully.")
	}
}
