package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

type Agent struct {
	client  *openai.Client
	model   string
	history []openai.ChatCompletionMessage
}

func NewAgent(cli *openai.Client, model string) *Agent {
	return &Agent{
		client: cli,
		model:  model,
		history: []openai.ChatCompletionMessage{
			{Role: "system", Content: `Assistants name is Care Research Economist Agent.
You are a helpful AI assistant specialized in economic research and analysis. You think and act like a PhD economist.
You are able to perform economic research and analysis on a wide range of topics. Some areas of your expertise are:
- Economic theory
- Econometrics theory and applications
- Microeconomics
- Macroeconomics
- International economics theory and applications
- Development economics theory and applications
- Policy analysis and evaluation
- Data analysis and visualization
- Machine learning and AI applications in economics
- Econometric modeling and forecasting
- Economic policy design and implementation
- Economic research and analysis

You have one tool available:

1. "read_file" – return the contents of a local file.

When you want to call a function, respond with the OpenAI tool calling format; set name to the function name and provide JSON arguments.
After you receive the function response (role:"tool"), continue reasoning and then answer the user in natural language.
`},
		},
	}
}

func (a *Agent) Ask(ctx context.Context, user string) error {
	if strings.TrimSpace(user) != "" {
		// append user message to running history
		a.history = append(a.history, openai.ChatCompletionMessage{
			Role: "user", Content: user,
		})
	}

	// Stream the answer, providing our tool/function definitions on every call
	stream, err := a.client.CreateChatCompletionStream(ctx, openai.ChatCompletionRequest{
		Model:     a.model,
		Stream:    true,
		Messages:  a.history,
		Functions: FunctionDefinitions,
	})
	if err != nil {
		return err
	}
	defer stream.Close()

	// We will accumulate the assistant's full message (content + any tool calls)
	var fullMsg openai.ChatCompletionMessage
	fullMsg.Role = openai.ChatMessageRoleAssistant

	// helper to track toolCall accumulation
	toolCallIndex := map[string]int{}

	for {
		chunk, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		delta := chunk.Choices[0].Delta

		// stream textual content to the terminal as we go
		if delta.Content != "" {
			fmt.Print(delta.Content)
			fullMsg.Content += delta.Content
		}

		// merge tool call deltas
		for _, tc := range delta.ToolCalls {
			idx, ok := toolCallIndex[tc.ID]
			if !ok {
				// first time seeing this call
				fullMsg.ToolCalls = append(fullMsg.ToolCalls, tc)
				toolCallIndex[tc.ID] = len(fullMsg.ToolCalls) - 1
			} else {
				existing := &fullMsg.ToolCalls[idx]
				// merge streamed arguments (they may arrive in pieces)
				var sb strings.Builder
				sb.WriteString(existing.Function.Arguments)
				sb.WriteString(tc.Function.Arguments)
				existing.Function.Arguments = sb.String()
			}
		}
	}

	fmt.Println() // newline after streaming output

	// If the assistant triggered a tool, handle it and loop once more
	if len(fullMsg.ToolCalls) > 0 {
		for _, tc := range fullMsg.ToolCalls {
			tool, found := toolMap[tc.Function.Name]
			if !found {
				fmt.Printf("⚠️ unknown tool: %s\n", tc.Function.Name)
				continue
			}

			out, err := tool.Call(json.RawMessage(tc.Function.Arguments))
			if err != nil {
				fmt.Printf("🛠  error: %v\n", err)
				continue
			}

			// Send the tool response back to the model so it can keep reasoning
			a.history = append(a.history, openai.ChatCompletionMessage{
				Role:       openai.ChatMessageRoleTool,
				ToolCallID: tc.ID,
				Content:    out,
			})
		}

		// Re-enter the loop with no new user input so the model can use the tool output
		return a.Ask(ctx, "")
	}

	// No tool calls → this is the final assistant answer. Save it to history.
	a.history = append(a.history, fullMsg)
	return nil
}
