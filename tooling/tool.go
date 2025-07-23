package tooling

import (
    "encoding/json"

    openai "github.com/sashabaranov/go-openai"
)

// Tool represents a single function/tool the LLM can call.
// It mirrors the struct used earlier in main but lives in its own package
// so that all tools can register themselves cleanly.
//
// Call should execute the tool and return a string result that will be fed
// back to the model.
// Definition returns the corresponding OpenAI function definition.
// (We generate it ad-hoc instead of keeping a separate value in every file.)
//
// Note: we do not embed the JSON schema generator; each concrete tool
// provides its own FunctionDefinition value instead.
type Tool struct {
    Name        string
    Description string
    Parameters  map[string]interface{}
    Call        func(input json.RawMessage) (string, error)
}

var (
    registry []Tool
    toolMap  = map[string]Tool{}
)

// Register makes a tool available to the agent.
func Register(t Tool) {
    registry = append(registry, t)
    toolMap[t.Name] = t
}

// Definition converts the Tool meta-data to an openai.FunctionDefinition.
func (t Tool) Definition() openai.FunctionDefinition {
    return openai.FunctionDefinition{
        Name:        t.Name,
        Description: t.Description,
        Parameters:  t.Parameters,
    }
}

// Definitions returns all registered tools as OpenAI function definitions.
func Definitions() []openai.FunctionDefinition {
    defs := make([]openai.FunctionDefinition, 0, len(registry))
    for _, t := range registry {
        defs = append(defs, t.Definition())
    }
    return defs
}

// Map gives access to the underlying name→Tool mapping.
func Map() map[string]Tool {
    return toolMap
}
