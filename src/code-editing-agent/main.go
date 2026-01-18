package main

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
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

	agent := NewAgent(&client, getUserMessage)
	err := agent.Run(context.TODO())

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		// no os.Exit(1)
	}
}

type Agent struct {
	client         *anthropic.Client
	getUserMessage func() (string, bool)
}

func NewAgent(client *anthropic.Client, getUserMessage func() (string, bool)) *Agent {
	return &Agent{
		client:         client,
		getUserMessage: getUserMessage,
	}
}

func (a *Agent) runInference(ctx context.Context, conversation []anthropic.MessageParam) (*anthropic.Message, error) {
	message, err := a.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.ModelClaudeSonnet4_5, //https://github.com/anthropics/anthropic-sdk-go/blob/main/message.go#L2081
		MaxTokens: int64(1024),
		Messages:  conversation,
	})
	return message, err
}

func (a *Agent) Run(ctx context.Context) error {
	// main interaction loop
	conversation := []anthropic.MessageParam{}
	fmt.Println("Chat with Agent (use 'ctrl-c' to quit)")

	for {
		fmt.Print("\u001b[94mYou\u001b[0m: ")
		userInputRaw, ok := a.getUserMessage()
		if !ok {
			break
		}

		userMessage := anthropic.NewUserMessage(anthropic.NewTextBlock(userInputRaw))
		conversation = append(conversation, userMessage)

		message, err := a.runInference(ctx, conversation)
		if err != nil {
			return err
		}
		conversation = append(conversation, message.ToParam())
		for _, content := range message.Content {
			switch content.Type {
			case "text":
				fmt.Printf("\u001b[93mAgent\u001b[0m: %s\n", content.Text)
			}
		}
	}
	return nil
}
