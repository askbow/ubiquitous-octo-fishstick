```
>go run main.go
Chat with Agent (use 'ctrl-c' to quit)
You: What's going on in main.go? Be concise! No more that 144 characters per function, or struct
tool: read_file({"path":"main.go"})
Agent: **main()**: Creates Anthropic client, sets up stdin scanner, initializes agent with read_file tool, runs agent loop.

**Agent**: Holds client, user input func, and tools. Manages conversation loop between user and Claude with tool execution.

**NewAgent()**: Constructor for Agent struct.

**runInference()**: Converts tools to API format, calls Claude API with conversation history, returns response.

**executeTool()**: Finds tool by name, executes it with input JSON, returns result block (success or error).

**Run()**: Main loop: gets user input, calls Claude, handles text/tool_use responses, appends tool results, repeats.

**ToolDefinition**: Tool metadata (name, desc, schema) + execution function.

**GenerateSchema[T]()**: Reflects Go struct into JSON schema for tool input validation.

**ReadFileDefinition**: Tool config for reading files.

**ReadFileInput**: Schema struct for file path parameter.

**ReadFileTool()**: Unmarshals input, reads file from path, returns contents as string.# ubiquitous-octo-fishstick
```