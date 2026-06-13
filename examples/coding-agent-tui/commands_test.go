package main

import (
	"strings"
	"testing"
)

func TestGetSlashCommand(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{name: "help", input: "/help", want: "/help"},
		{name: "help with spaces", input: "  /help  ", want: "/help"},
		{name: "prompt not command", input: "explain auth", want: ""},
		{name: "unknown", input: "/unknown", want: ""},
		{name: "model", input: "/model", want: "/model"},
		{name: "quit alias", input: "/quit", want: "/quit"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getSlashCommand(tt.input); got != tt.want {
				t.Fatalf("got %q want %q", got, tt.want)
			}
		})
	}
}

func TestFormatSlashHelp(t *testing.T) {
	help := formatSlashHelp()
	if help == "" {
		t.Fatal("expected non-empty help")
	}
	if !containsAll(help, "/help", "/model", "/exit") {
		t.Fatalf("missing commands in help: %q", help)
	}
}

func containsAll(s string, parts ...string) bool {
	for _, p := range parts {
		if !strings.Contains(s, p) {
			return false
		}
	}
	return true
}
