//go:build legacytools
// +build legacytools

package main

import (
	"encoding/json"
	"os"
	"strings"

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
	{
		Name:        "edit_file",
		Description: "Overwrite the file at the given path with the provided content. Creates the file if it doesn't exist. Only use edit_file only after read_file.",
		Parameters: map[string]interface{}{
			"type":     "object",
			"required": []string{"path", "content"},
			"properties": map[string]interface{}{
				"path":    map[string]interface{}{"type": "string"},
				"content": map[string]interface{}{"type": "string"},
			},
		},
	},
	{
		Name:        "list_files",
		Description: "Return a newline separated list of files and directories in the given local path.",
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
	addTool(Tool{
		Name: "edit_file",
		Call: func(input json.RawMessage) (string, error) {
			var req struct {
				Path    string `json:"path"`
				Content string `json:"content"`
			}
			if err := json.Unmarshal(input, &req); err != nil {
				return "", err
			}
			// Limit content size to prevent abuse
			if len(req.Content) > 32_000 {
				req.Content = req.Content[:32_000]
			}
			if err := os.WriteFile(req.Path, []byte(req.Content), 0o644); err != nil {
				return "", err
			}
			return "ok", nil // simple acknowledgement
		},
	})
	addTool(Tool{
		Name: "list_files",
		Call: func(input json.RawMessage) (string, error) {
			var req struct {
				Path string `json:"path"`
			}
			if err := json.Unmarshal(input, &req); err != nil {
				return "", err
			}
			entries, err := os.ReadDir(req.Path)
			if err != nil {
				return "", err
			}
			names := make([]string, 0, len(entries))
			for _, e := range entries {
				names = append(names, e.Name())
			}
			result := strings.Join(names, "\n")
			if len(result) > 16_000 {
				result = result[:16_000]
			}
			return result, nil // → will go back to the LLM
		},
	})
}
