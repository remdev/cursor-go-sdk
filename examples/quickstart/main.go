// Quickstart ports cursor/cookbook sdk/quickstart:
// https://github.com/cursor/cookbook/blob/main/sdk/quickstart/src/index.ts
//
// Requires CURSOR_API_KEY and bridge setup (cd bridge && npm install).
//
//	go run ./examples/quickstart
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/remdev/cursor-go-sdk/cursor"
)

func main() {
	ctx := context.Background()
	apiKey := os.Getenv("CURSOR_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "set CURSOR_API_KEY")
		os.Exit(1)
	}

	model := os.Getenv("CURSOR_MODEL")
	if model == "" {
		model = "composer-2"
	}

	wd, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	agent, err := cursor.CreateAgent(ctx, cursor.AgentOptions{
		Name:   "SDK quickstart",
		Model:  model,
		APIKey: apiKey,
		Local:  &cursor.LocalAgentOptions{CWD: []string{wd}},
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer agent.Close(ctx)

	run, err := agent.Send(ctx, "Explain this project in one paragraph.", cursor.SendOptions{})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for msg, err := range run.Messages(ctx) {
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		if msg.Type != "assistant" {
			continue
		}
		fmt.Print(cursor.AssistantText(msg))
	}

	result, err := run.Wait(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if result.Status != cursor.RunStatusFinished {
		fmt.Fprintf(os.Stderr, "\nrun ended with status %s\n", result.Status)
		os.Exit(2)
	}
	fmt.Println()
}
