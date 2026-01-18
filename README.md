# A minimal coding agent
Largely based on https://ampcode.com/how-to-build-an-agent tutorial with some minor updates and personal touches

Requires `ANTHROPIC_API_KEY` to be set.


### Read and analyse a file:
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
### Listing and reading files:
```
Chat with Agent (use 'ctrl-c' to quit)
You: Describe what you see in this directory. Examine every file and provide a summary description, no longer than 144 characters per file or directory.
Agent: I'll explore the directory structure and examine all files to provide you with a comprehensive summary.
tool: list_files({})
Agent: Now let me examine each file:
tool: read_file({"path":"go.mod"})
tool: read_file({"path":"go.sum"})
tool: read_file({"path":"main.go"})
tool: read_file({"path":"test-text-file.txt"})
Agent: ## Directory Summary

**go.mod** - Go module file defining "agent" project with Anthropic SDK v1.19.0 & jsonschema v0.13.0 dependencies for Go 1.24.5

**go.sum** - Checksums for Go dependencies including Anthropic SDK, jsonschema, JSON parsing libs, and testing utilities

**main.go** - AI agent implementation using Anthropic's Claude API with tools for file ops; includes REPL loop for chat interaction

**test-text-file.txt** - Contains a riddle: "what animal is the most disagreeable because it always says neigh?" (Answer: a horse)
```
## Creating new code
```
> go run main.go
Chat with Agent (use 'ctrl-c' to quit)
You: create fizzbuzz.py that I can run wth python3 and that implements a short fizzbuzz solution and executes it
Agent: I'll create a fizzbuzz.py file with a short implementation that executes automatically.
tool: edit_file({"path":"fizzbuzz.py","old_str":"placeholder","new_str":"#!/usr/bin/env python3\n\nfor i in range(1, 101):\n    print(\"Fizz\"*(i%3==0) + \"Buzz\"*(i%5==0) or i)\n"})
Agent: Let me create the file properly by providing empty content first:
tool: edit_file({"path":"fizzbuzz.py","old_str":"","new_str":"#!/usr/bin/env python3\n\nfor i in range(1, 101):\n    print(\"Fizz\"*(i%3==0) + \"Buzz\"*(i%5==0) or i)\n"})
Agent: Perfect! I've created `fizzbuzz.py` with a short, concise FizzBuzz implementation.

The script:
- Iterates through numbers 1 to 100
- Uses a clever one-liner that concatenates "Fizz" (if divisible by 3) and "Buzz" (if divisible by 5)
- Uses the `or` operator to print the number itself if the string is empty
- Automatically executes when run

You can now run it with:
```bash
python3 fizzbuzz.py
```

This will print:
- "Fizz" for multiples of 3
- "Buzz" for multiples of 5
- "FizzBuzz" for multiples of both 3 and 5
- The number itself otherwise
You:
PS C:\Users\Askbow\Documents\GitHub\ubiquitous-octo-fishstick\src\code-editing-agent> python .\fizzbuzz.py
1
2
Fizz
4
Buzz
Fizz
7
8
Fizz
Buzz
11
Fizz
13
14
FizzBuzz
16
17
Fizz
19
Buzz
Fizz
22
23
Fizz
Buzz
26
Fizz
28
29
FizzBuzz
31
32
Fizz
34
Buzz
Fizz
37
38
Fizz
Buzz
41
Fizz
43
44
FizzBuzz
46
47
Fizz
49
Buzz
Fizz
52
53
Fizz
Buzz
56
Fizz
58
59
FizzBuzz
61
62
Fizz
64
Buzz
Fizz
67
68
Fizz
Buzz
71
Fizz
73
74
FizzBuzz
76
77
Fizz
79
Buzz
Fizz
82
83
Fizz
Buzz
86
Fizz
88
89
FizzBuzz
91
92
Fizz
94
Buzz
Fizz
97
98
Fizz
Buzz
```
### Updating code
```
PS C:\Users\Askbow\Documents\GitHub\ubiquitous-octo-fishstick\src\code-editing-agent> go run main.go
Chat with Agent (use 'ctrl-c' to quit)
You: Update fizzbuzz.py to only print results until 15
Agent: I'll help you update the fizzbuzz.py file. Let me first read the current content to see what we're working with.
tool: read_file({"path":"fizzbuzz.py"})
Agent: Now I'll update the file to only print results until 15 by changing the range from `range(1, 101)` to `range(1, 16)`:
tool: edit_file({"path":"fizzbuzz.py","old_str":"for i in range(1, 101):\n    print(\"Fizz\"*(i%3==0) + \"Buzz\"*(i%5==0) or i)","new_str":"for i in range(1, 16):\n    print(\"Fizz\"*(i%3==0) + \"Buzz\"*(i%5==0) or i)"})
Agent: Done! I've updated fizzbuzz.py to only print results from 1 to 15 (inclusive) instead of 1 to 100. The range has been changed from `range(1, 101)` to `range(1, 16)`.
You:
```