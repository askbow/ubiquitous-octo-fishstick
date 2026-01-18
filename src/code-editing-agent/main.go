package main

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"encoding/json"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/invopop/jsonschema"
)

func main() {
	client := anthropic.NewClient(
		// setting variables is scoped to User or System on windows
		// it is astonishing there is no `export` like in Linux
		option.WithAPIKey(os.Getenv("ANTHROPIC_API_KEY")), // defaults to os.LookupEnv("ANTHROPIC_API_KEY")
	)
	scanner := bufio.NewScanner(os.Stdin)

	getUserMessage := func() (string, bool) {
		if !scanner.Scan() {
			return "", false
		}
		return scanner.Text(), true
	}

	tools := []ToolDefinition{
		ReadFileDefinition,
	}

	agent := NewAgent(&client, getUserMessage, tools)
	err := agent.Run(context.TODO())

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		// no os.Exit(1)
	}
}

// The Agent main loop stuff:

type Agent struct {
	client         *anthropic.Client
	getUserMessage func() (string, bool)
	tools          []ToolDefinition
}

func NewAgent(
	client *anthropic.Client,
	getUserMessage func() (string, bool),
	tools []ToolDefinition,
) *Agent {
	return &Agent{
		client:         client,
		getUserMessage: getUserMessage,
		tools:          tools,
	}
}

func (a *Agent) runInference(ctx context.Context, conversation []anthropic.MessageParam) (*anthropic.Message, error) {
	anthropicTools := []anthropic.ToolUnionParam{}
	for _, tool := range a.tools {
		anthropicTools = append(anthropicTools, anthropic.ToolUnionParam{
			OfTool: &anthropic.ToolParam{
				Name:        tool.Name,
				Description: anthropic.String(tool.Description),
				InputSchema: tool.InputSchema,
			},
		})
	}

	message, err := a.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.ModelClaudeSonnet4_5, //https://github.com/anthropics/anthropic-sdk-go/blob/main/message.go#L2081
		MaxTokens: int64(1024),
		Messages:  conversation,
		Tools:     anthropicTools,
	})
	return message, err
}

func (a *Agent) executeTool(id, name string, input json.RawMessage) anthropic.ContentBlockParamUnion {
	var toolDef ToolDefinition
	var found bool //is this a common pattern in go? do they have hashmaps?
	for _, tool := range a.tools {
		if tool.Name == name {
			toolDef = tool
			found = true
			break
		}
	}
	if !found {
		return anthropic.NewToolResultBlock(id, "tool not found", true)
	}
	fmt.Printf("\u001b[92mtool\u001b[0m: %s(%s)\n", name, input)
	response, err := toolDef.Function(input)
	if err != nil {
		return anthropic.NewToolResultBlock(id, err.Error(), true)
	}
	return anthropic.NewToolResultBlock(id, response, false)
}

func (a *Agent) Run(ctx context.Context) error {
	// main interaction loop
	conversation := []anthropic.MessageParam{}
	fmt.Println("Chat with Agent (use 'ctrl-c' to quit)")
	readUserInput := true
	for {
		if readUserInput {
			fmt.Print("\u001b[94mYou\u001b[0m: ")
			userInputRaw, ok := a.getUserMessage()
			if !ok {
				break
			}

			userMessage := anthropic.NewUserMessage(anthropic.NewTextBlock(userInputRaw))
			conversation = append(conversation, userMessage)
		}
		message, err := a.runInference(ctx, conversation)
		if err != nil {
			return err
		}
		conversation = append(conversation, message.ToParam())

		toolResults := []anthropic.ContentBlockParamUnion{}

		for _, content := range message.Content {
			switch content.Type {
			case "text":
				fmt.Printf("\u001b[93mAgent\u001b[0m: %s\n", content.Text)

			case "tool_use":
				result := a.executeTool(content.ID, content.Name, content.Input)
				toolResults = append(toolResults, result)
			}
		}
		if len(toolResults) == 0 {
			readUserInput = true
			continue
		}
		readUserInput = false
		conversation = append(conversation, anthropic.NewUserMessage(toolResults...))

	}
	return nil
}

// Tools:

type ToolDefinition struct {
	Name        string                                      `json:"name"`
	Description string                                      `json:"description"`
	InputSchema anthropic.ToolInputSchemaParam              `json:"input_schema"`
	Function    func(input json.RawMessage) (string, error) `json:"-"`
}

func GenerateSchema[T any]() anthropic.ToolInputSchemaParam {
	reflector := jsonschema.Reflector{
		AllowAdditionalProperties: false,
		DoNotReference:            true,
	}
	var v T
	schema := reflector.Reflect(v)
	return anthropic.ToolInputSchemaParam{
		Properties: schema.Properties,
	}
}

// ReadFileTool
var ReadFileDefinition = ToolDefinition{
	Name: "read_file",
	Description: "Read the contents from the given relative file path. " +
		"Use this when you want to see what's inside a file." +
		" Do not use this with directory names.",
	InputSchema: ReadFileInputSchema,
	Function:    ReadFileTool,
}

type ReadFileInput struct {
	Path string `json:"path" jsonschema_description:"The relative path of a file in the working directory."`
}

var ReadFileInputSchema = GenerateSchema[ReadFileInput]()

func ReadFileTool(input json.RawMessage) (string, error) {
	// Read File tool implementation
	readFileInput := ReadFileInput{}
	err := json.Unmarshal(input, &readFileInput)
	if err != nil {
		panic(err)
	}
	content, err := os.ReadFile(readFileInput.Path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
