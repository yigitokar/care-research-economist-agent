package main

import (
	"encoding/json"
	"os"

	openai "github.com/sashabaranov/go-openai"
)

type Tool struct {
	Name string
	Call func(input json.RawMessage) (string, error)
}

var toolMap = map[string]Tool{}

// FunctionDefinitions is sent with every chat completion request so the model
// understands which tools are available and their JSON schema.
var FunctionDefinitions = []openai.FunctionDefinition{
	{
		Name:        "read_file",
		Description: "Return the contents of a local file as text.",
		Parameters: map[string]interface{}{
			"type":     "object",
			"required": []string{"path"},
			"properties": map[string]interface{}{
				"path": map[string]interface{}{"type": "string"},
			},
		},
	},
}

// helper to register tools once in init()
func addTool(t Tool) { toolMap[t.Name] = t }

// ---- concrete implementation of read_file ----
func init() {
	addTool(Tool{
		Name: "read_file",
		Call: func(input json.RawMessage) (string, error) {
			var req struct {
				Path string `json:"path"`
			}
			if err := json.Unmarshal(input, &req); err != nil {
				return "", err
			}
			data, err := os.ReadFile(req.Path)
			if err != nil {
				return "", err
			}
			if len(data) > 16_000 {
				data = data[:16_000]
			}
			return string(data), nil // → will go back to the LLM
		},
	})
}
