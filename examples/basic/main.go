// Example basic demonstrates a one-shot local prompt (smoke test).
//
// Requires CURSOR_API_KEY and bridge setup (cd bridge && npm install).
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
	wd, _ := os.Getwd()

	result, err := cursor.Prompt(ctx, "Reply with exactly: pong", cursor.AgentOptions{
		Model:  "composer-2.5",
		APIKey: os.Getenv("CURSOR_API_KEY"),
		Local:  &cursor.LocalAgentOptions{CWD: []string{wd}},
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("status:", result.Status)
	fmt.Println("result:", result.Result)
}
