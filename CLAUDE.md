# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

**Build and Run:**
```bash
go build              # Build the application
go run .              # Run the application directly
go mod tidy           # Clean up module dependencies
go mod download       # Download all dependencies
```

**Development:**
```bash
go test ./...         # Run all tests (when tests are added)
go fmt ./...          # Format all Go files
go vet ./...          # Run Go vet for static analysis
```

## Architecture Overview

This is a CLI-based economic research agent that provides an interactive REPL interface for economic analysis and research assistance.

**Core Components:**
- `main.go`: Entry point implementing a REPL loop with ASCII banner display
- `agent.go`: AI agent implementation using Groq's API (OpenAI-compatible) with streaming chat completions
- `tools.go`: Defines AI tools (currently `read_file` for file access)

**Key Design Patterns:**
- The agent maintains conversation history across interactions
- Uses OpenAI-style function/tool calling for extending capabilities
- Implements streaming responses for real-time output
- Environment-based configuration for API credentials

**API Integration:**
- Uses Groq's API endpoint: `https://api.groq.com/openai/v1`
- Model: `moonshotai/kimi-k2-instruct`
- Requires `GROQ_API_KEY` in `.env` file

**Tool System:**
The agent supports extensible tools through the OpenAI function calling interface. New tools can be added in `tools.go` by:
1. Defining the tool schema in `getTools()`
2. Implementing the handler in `executeTool()`