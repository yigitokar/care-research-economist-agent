package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/common-nighthawk/go-figure"
	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"
)

func main() {
	// Load .env so GROQ_API_KEY becomes visible
	_ = godotenv.Load()

	// Display ASCII art banner
	myFigure := figure.NewFigure("Economist Agent", "speed", true)
	myFigure.Print()
	fmt.Println()

	key := os.Getenv("GROQ_API_KEY")
	if key == "" {
		log.Fatalln("GROQ_API_KEY not set—add it to .env")
	}

	// Build a client that talks to Groq (OpenAI‑compatible)
	cfg := openai.DefaultConfig(key)
	cfg.BaseURL = "https://api.groq.com/openai/v1"
	cli := openai.NewClientWithConfig(cfg)

	// Spin up our agent using the model you chose
	agent := NewAgent(cli, "moonshotai/kimi-k2-instruct")

	// Simple REPL
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("\n› ")
		line, err := reader.ReadString('\n')
		if err != nil { // Ctrl‑D = EOF → quit
			fmt.Println("\nBye!")
			break
		}

		cmd := strings.TrimSpace(line)
		if cmd == "/quit" || cmd == "/exit" {
			fmt.Println("Bye!")
			break
		}

		if err := agent.Ask(context.Background(), line); err != nil {
			fmt.Println("error:", err)
		}
	}

}
