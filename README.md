# Code Editing Agent

A CLI-based economic research agent that provides an interactive REPL interface for economic analysis and research assistance.

## Features

- **Interactive REPL Interface**: Chat with an AI economist specialized in economic research and analysis
- **File Operations**: Read, edit, and manage files for research projects
- **Bash Integration**: Execute shell commands for data processing and automation
- **Streaming Responses**: Real-time output as the agent processes your requests
- **Persistent Conversation**: Maintains context throughout your research session

## Prerequisites

- Go 1.19 or higher
- Groq API key

## Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/code-editing-agent.git
cd code-editing-agent
```

2. Create a `.env` file in the root directory:
```bash
GROQ_API_KEY=your_groq_api_key_here
```

3. Install dependencies:
```bash
go mod download
```

## Usage

Run the agent:
```bash
go run .
```

Or build and run:
```bash
go build
./code-editing-agent
```

## Available Tools

The agent has access to four tools:

1. **read_file**: Read contents of local files
2. **edit_file**: Create or overwrite files with new content
3. **list_files**: List files and directories in a given path
4. **bash**: Execute bash commands for advanced file operations and data processing

## Example Usage

```
🤖 Care Research Economist Agent
> help me analyze the relationship between GDP and unemployment

I'll help you analyze the relationship between GDP and unemployment. Let me create a structured analysis for you...

> create a literature review on scaling laws in economics

I'll create a comprehensive literature review on scaling laws in economics...
```

## Architecture

- `main.go`: Entry point with REPL loop and ASCII banner
- `agent.go`: Core AI agent implementation using Groq's API
- `tooling/`: Package containing all available tools
  - `tool.go`: Tool interface and registry
  - `read_file.go`: File reading functionality
  - `edit_file.go`: File writing functionality
  - `list_files.go`: Directory listing functionality
  - `bash.go`: Bash command execution

## Configuration

The agent uses the following configuration:
- **API**: Groq's OpenAI-compatible endpoint
- **Model**: `moonshotai/kimi-k2-instruct`
- **Environment**: Requires `GROQ_API_KEY` in `.env` file

## Contributing

Feel free to open issues or submit pull requests for improvements.

## License

[Add your license here]