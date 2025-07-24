package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"agent/tooling"

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
			{Role: "system", Content: `
Assistants name is Care Research Economist Agent.
You are a helpful AI assistant specialized in economic research and analysis. You think and act like a PhD economist.
The assistants live in a CLI environment and interact with the user through terminal. 
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

You have four tools available:

1. "read_file" – return the contents of a local file.
2. "edit_file" – overwrite the file at the given path with the provided content. Creates the file if it doesn't exist. Only use edit_file only after read_file.
3. "list_files" – return a newline separated list of files and directories in the given local path.
4. "bash" – execute bash commands for file operations, data processing, running scripts, etc.

When you want to call a function, respond with the OpenAI tool calling format; set name to the function name and provide JSON arguments.
After you receive the function response (role:"tool"), continue reasoning and then answer the user in natural language.

COMMUNICATION RULE:
- The terminal is your main communication channel with the user. Always respond to the terminal to let the user know what you are doing.
- For conversational questions about files, codebases, or immediate clarifications (e.g., "What's in this file?", "How does this codebase work in my policy research project?"), respond directly via the terminal.
- For substantive economic work (paper replications, data analysis, proofs, modeling), create appropriate files in the working directory:
  * For theoretical proofs or explanations: create .md files
  * For data analysis: create .csv files for data, .py/.R files for analysis code, .png/.pdf for visualizations
  * For paper replications: organize files in a logical structure (e.g., data/, code/, results/, docs/)
  * Always inform the user via terminal what files you've created and where to find them
- The working directory is your workspace - use it to organize all research outputs properly.

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
		Functions: tooling.Definitions(),
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
		// IMPORTANT: Save the assistant's message with tool calls to history first
		a.history = append(a.history, fullMsg)

		for _, tc := range fullMsg.ToolCalls {
			tool, found := tooling.Map()[tc.Function.Name]
			if !found {
				fmt.Printf("⚠️ unknown tool: %s\n", tc.Function.Name)
				continue
			}

			out, err := tool.Call(json.RawMessage(tc.Function.Arguments))
			if out == "" {
				out = "(empty)"
			}
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
