package tooling

import (
    "encoding/json"
    "os/exec"
    "strings"
)

func init() {
    Register(bashTool)
}

var bashTool = Tool{
    Name:        "bash",
    Description: "Execute a bash command and return its output",
    Parameters: map[string]interface{}{
        "type":     "object",
        "required": []string{"command"},
        "properties": map[string]interface{}{
            "command": map[string]interface{}{
                "type":        "string",
                "description": "The bash command to execute",
            },
        },
    },
    Call: func(input json.RawMessage) (string, error) {
        var req struct {
            Command string `json:"command"`
        }
        if err := json.Unmarshal(input, &req); err != nil {
            return "", err
        }
        
        // Execute the command
        cmd := exec.Command("bash", "-c", req.Command)
        output, err := cmd.CombinedOutput()
        
        result := strings.TrimSpace(string(output))
        if err != nil {
            return result + "\nError: " + err.Error(), err
        }
        
        if result == "" {
            return "(command completed successfully with no output)", nil
        }
        
        return result, nil
    },
}