// Interactive coding-agent TUI ports cursor/cookbook sdk/coding-agent-cli TUI mode:
// https://github.com/cursor/cookbook/tree/main/sdk/coding-agent-cli
//
// Usage:
//
//	export CURSOR_API_KEY=...
//	go run ./examples/coding-agent-tui
//
// Requires cursor-sdk-bridge on PATH (npm install -g @cursor-go-sdk/cursor-sdk-bridge) and a TTY.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/remdev/cursor-go-sdk/examples/agentutil"
)

func main() {
	cwd := flag.String("cwd", "", "workspace directory (default: cwd)")
	model := flag.String("model", agentutil.EnvOr("CURSOR_MODEL", "composer-2"), "model id")
	force := flag.Bool("force", false, "pass local force on send (local mode)")
	flag.Parse()

	if !agentutil.IsTerminal(os.Stdin) || !agentutil.IsTerminal(os.Stdout) {
		fmt.Fprintln(os.Stderr, "coding-agent-tui requires an interactive terminal")
		os.Exit(1)
	}

	apiKey := os.Getenv("CURSOR_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "set CURSOR_API_KEY")
		os.Exit(1)
	}

	workdir := *cwd
	if workdir == "" {
		var err error
		workdir, err = os.Getwd()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	sess, err := newSession(apiKey, workdir, *model, *force)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer sess.close(context.Background())

	if _, err := tea.NewProgram(newAppModel(sess), tea.WithAltScreen()).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
