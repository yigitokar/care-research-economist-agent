package tooling

import (
    "encoding/json"
    "os"
)

func init() {
    Register(editFileTool)
}

var editFileTool = Tool{
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
    Call: func(input json.RawMessage) (string, error) {
        var req struct {
            Path    string `json:"path"`
            Content string `json:"content"`
        }
        if err := json.Unmarshal(input, &req); err != nil {
            return "", err
        }
        if len(req.Content) > 32_000 {
            req.Content = req.Content[:32_000]
        }
        if err := os.WriteFile(req.Path, []byte(req.Content), 0o644); err != nil {
            return "", err
        }
        return "ok", nil
    },
}
