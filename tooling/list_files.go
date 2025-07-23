package tooling

import (
    "encoding/json"
    "os"
    "strings"
)

func init() {
    Register(listFilesTool)
}

var listFilesTool = Tool{
    Name:        "list_files",
    Description: "Return a newline separated list of files and directories in the given local path.",
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
        return result, nil
    },
}
