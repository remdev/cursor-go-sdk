// Coding-agent CLI (plain mode) ports cursor/cookbook sdk/coding-agent-cli runPlainPrompt:
// https://github.com/cursor/cookbook/tree/main/sdk/coding-agent-cli
//
// Usage:
//
//	go run ./examples/coding-agent-cli -- "Explain the auth flow"
//	printf "Review changes\n" | go run ./examples/coding-agent-cli
//
// Requires CURSOR_API_KEY and bridge setup (cd bridge && npm install).
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/remdev/cursor-go-sdk/cursor"
	"github.com/remdev/cursor-go-sdk/examples/agentutil"
)

func main() {
	cwd := flag.String("cwd", "", "workspace directory (default: cwd)")
	model := flag.String("model", agentutil.EnvOr("CURSOR_MODEL", "composer-2"), "model id")
	force := flag.Bool("force", false, "reserved for stuck-run recovery (not wired yet)")
	flag.Parse()

	if *force {
		fmt.Fprintln(os.Stderr, "warning: --force is not implemented in this minimal port")
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

	prompt := strings.TrimSpace(strings.Join(flag.Args(), " "))
	if prompt == "" && agentutil.IsTerminal(os.Stdin) {
		printHelp()
		fmt.Fprintln(os.Stderr, "\nFor interactive TUI: go run ./examples/coding-agent-tui")
		os.Exit(1)
	}
	if prompt == "" {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		prompt = strings.TrimSpace(string(data))
	}
	if prompt == "" {
		fmt.Fprintln(os.Stderr, "no prompt provided")
		os.Exit(1)
	}

	ctx := context.Background()
	agent, err := cursor.CreateAgent(ctx, cursor.AgentOptions{
		Model:  *model,
		APIKey: apiKey,
		Local:  &cursor.LocalAgentOptions{CWD: []string{workdir}},
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer agent.Close(ctx)

	run, err := agent.Send(ctx, prompt, cursor.SendOptions{})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	assistantNL := true
	annotate := func(line string) {
		if !assistantNL {
			fmt.Fprintln(os.Stderr)
		}
		fmt.Fprintln(os.Stderr, line)
		assistantNL = true
	}

	for msg, err := range run.Messages(ctx) {
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		switch msg.Type {
		case "assistant":
			text := cursor.AssistantText(msg)
			if text != "" {
				fmt.Print(text)
				assistantNL = strings.HasSuffix(text, "\n")
			}
		case "thinking":
			if t := strings.TrimSpace(msg.ThinkingText); t != "" {
				annotate("[thinking] " + agentutil.Compact(t))
			}
		case "tool":
			annotate(fmt.Sprintf("[tool] %s %s", msg.ToolStatus, msg.ToolName))
		case "status":
			if msg.Status != "" && msg.Status != "FINISHED" {
				line := "[status] " + msg.Status
				if msg.StatusMessage != "" {
					line += " " + msg.StatusMessage
				}
				annotate(line)
			}
		case "task":
			if msg.TaskText != "" || msg.TaskStatus != "" {
				annotate("[task] " + agentutil.Compact(msg.TaskStatus+" "+msg.TaskText))
			}
		}
	}

	result, err := run.Wait(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	annotate(fmt.Sprintf("[done] status=%s durationMs=%d", result.Status, result.DurationMS))
	if result.Status == cursor.RunStatusError {
		os.Exit(2)
	}
}

func printHelp() {
	fmt.Fprintf(os.Stderr, `Lightweight coding agent CLI (Go port, plain mode)

Usage:
  coding-agent-cli [flags] "your task"
  coding-agent-cli [flags]   # reads prompt from stdin

Flags:
  -cwd string     Workspace directory (default: cwd)
  -model string   Model id (default: CURSOR_MODEL or composer-2)

Examples:
  go run ./examples/coding-agent-cli -- "Explain the auth flow"
  printf "Review changes" | go run ./examples/coding-agent-cli
`)
}
