package tooling

import (
    "encoding/json"
    "os"
)

func init() {
    Register(readFileTool)
}

var readFileTool = Tool{
    Name:        "read_file",
    Description: "Return the contents of a local file as text.",
    Parameters: map[string]interface{}{
        "type":     "object",
        "required": []string{"path"},
        "properties": map[string]interface{}{
            "path": map[string]interface{}{"type": "string"},
        },
    },
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
        return string(data), nil
    },
}
