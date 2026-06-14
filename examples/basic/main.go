// Example basic demonstrates a one-shot local prompt (smoke test).
//
// Requires CURSOR_API_KEY and cursor-sdk-bridge on PATH (npm install -g @cursor-go-sdk/cursor-sdk-bridge).
// Optional: CURSOR_MODEL (default composer-2).
//
//	go run ./examples/basic
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

	wd, _ := os.Getwd()

	result, err := cursor.Prompt(ctx, "Reply with exactly: pong", cursor.AgentOptions{
		Model:  model,
		APIKey: apiKey,
		Local:  &cursor.LocalAgentOptions{CWD: []string{wd}},
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("status:", result.Status)
	fmt.Println("result:", result.Result)
}
