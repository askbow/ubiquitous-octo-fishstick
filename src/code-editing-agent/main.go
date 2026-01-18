package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"text/scanner"

	"github.com/anthropics/anthropic-sdk-go"
)

func main() {
    client := anthropic.NewClient(os.Getenv("ANTHROPIC_API_KEY"))

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

type Agent struct {
	client *anthropic.Client
	getUserMessage func() (string, bool)
}

func NewAgent(client *anthropic.Client, getUserMessage func() (string, bool)) *Agent {
	return &Agent{
		client: client,
		getUserMessage: getUserMessage,
	}
}
