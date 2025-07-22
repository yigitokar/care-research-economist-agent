package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	openai "github.com/sashabaranov/go-openai"
)

func main() {
	// 1. read .env so GROQ_API_KEY becomes an env‑var
	_ = godotenv.Load() // silently does nothing if .env isn’t there

	// 2. grab the key (fail fast if it’s missing)
	key := os.Getenv("GROQ_API_KEY")
	if key == "" {
		log.Fatalln("GROQ_API_KEY not set—put it in .env like GROQ_API_KEY=sk-...")
	}

	// 3. build a client that talks to Groq, not OpenAI
	cfg := openai.DefaultConfig(key)
	cfg.BaseURL = "https://api.groq.com/openai/v1" //  ← the crucial swap
	cli := openai.NewClientWithConfig(cfg)

	// // 4. send one “Hello” request to the 8‑b Llama‑3
	// resp, err := cli.CreateChatCompletion(
	// 	context.Background(),
	// 	openai.ChatCompletionRequest{
	// 		Model: "moonshotai/kimi-k2-instruct",
	// 		Messages: []openai.ChatCompletionMessage{
	// 			{Role: openai.ChatMessageRoleUser, Content: "Hello from Go!"},
	// 		},
	// 	},
	// )
	// if err != nil {
	// 	log.Fatalln("API call failed:", err)
	// }

	// // 5. print Groq’s reply
	// fmt.Println(resp.Choices[0].Message.Content)
	// 4. send the request *as a stream* so we see tokens live
	stream, err := cli.CreateChatCompletionStream(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:  "moonshotai/kimi-k2-instruct",
			Stream: true, // ← key flag
			Messages: []openai.ChatCompletionMessage{
				{Role: openai.ChatMessageRoleUser, Content: "Tell me a joke about Go."},
			},
		},
	)
	if err != nil {
		log.Fatalln("API call failed:", err)
	}
	defer stream.Close()

	// 5. read each chunk and print immediately
	for {
		chunk, err := stream.Recv()
		if err != nil { // io.EOF when stream ends
			break
		}
		fmt.Print(chunk.Choices[0].Delta.Content)
	}
	fmt.Println() // newline after the stream finishes
}
