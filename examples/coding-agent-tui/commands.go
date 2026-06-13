package main

import (
	"fmt"
	"strings"
)

type slashCommand struct {
	name    string
	summary string
}

var slashCommands = []slashCommand{
	{name: "/help", summary: "Show available commands."},
	{name: "/local", summary: "Run future prompts in the local workspace."},
	{name: "/cloud", summary: "Run future prompts in Cursor cloud."},
	{name: "/model", summary: "Open a picker with available Cursor models."},
	{name: "/reset", summary: "Start a fresh agent and clear context."},
	{name: "/exit", summary: "Exit the TUI."},
	{name: "/quit", summary: "Exit the TUI."},
}

func getSlashCommand(input string) string {
	parts := strings.Fields(strings.TrimSpace(input))
	if len(parts) == 0 {
		return ""
	}
	cmd := parts[0]
	for _, c := range slashCommands {
		if c.name == cmd {
			return cmd
		}
	}
	return ""
}

func formatSlashHelp() string {
	parts := make([]string, 0, len(slashCommands))
	for _, c := range slashCommands {
		parts = append(parts, fmt.Sprintf("%s - %s", c.name, c.summary))
	}
	return strings.Join(parts, "  ")
}
