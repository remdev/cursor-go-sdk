package cursor_test

import (
	"testing"

	"github.com/remdev/cursor-go-sdk/cursor"
)

func TestMcpServerFromMapFlatHTTP(t *testing.T) {
	wire := cursor.McpServerFromMap(map[string]any{
		"type": "http",
		"url":  "https://example.com/mcp",
		"headers": map[string]string{"Authorization": "Bearer x"},
	})
	http, ok := wire["http"].(map[string]any)
	if !ok {
		t.Fatal("missing http block")
	}
	if http["url"] != "https://example.com/mcp" {
		t.Fatalf("url=%v", http["url"])
	}
}

func TestMcpServerFromMapFlatStdio(t *testing.T) {
	wire := cursor.McpServerFromMap(map[string]any{
		"command": "npx",
		"args":    []string{"-y", "server"},
	})
	stdio, ok := wire["stdio"].(map[string]any)
	if !ok {
		t.Fatal("missing stdio block")
	}
	if stdio["command"] != "npx" {
		t.Fatalf("command=%v", stdio["command"])
	}
}

func TestGetRunOptionsWire(t *testing.T) {
	w := cursor.GetRunOptions{
		Runtime: "cloud",
		AgentID: "bc-123",
		APIKey:  "key",
	}.ToWire()
	if w["runtime"] != "cloud" || w["agentId"] != "bc-123" {
		t.Fatalf("wire=%v", w)
	}
}

func TestLocalOptionsHelper(t *testing.T) {
	local := cursor.LocalOptions("/repo")
	if len(local.CWD) != 1 || local.CWD[0] != "/repo" {
		t.Fatalf("cwd=%v", local.CWD)
	}
}

func TestSendOptionsDeltaStepFlags(t *testing.T) {
	calledDelta := false
	calledStep := false
	opts := cursor.SendOptions{
		OnDelta: func(cursor.InteractionUpdate) { calledDelta = true },
		OnStep:  func(cursor.ConversationStep) { calledStep = true },
	}
	w := opts.ToWire()
	if w["enableDeltas"] != true || w["enableSteps"] != true {
		t.Fatalf("wire=%v", w)
	}
	_ = calledDelta
	_ = calledStep
}

func TestAgentDefinitionMcpServers(t *testing.T) {
	def := cursor.AgentDefinition{
		Description: "reviewer",
		Prompt:      "review code",
		McpServers:  []any{"docs", cursor.HttpMcpServerConfig{URL: "https://x/mcp"}},
	}
	w := def.ToWire()
	servers, ok := w["mcpServers"].([]map[string]any)
	if !ok || len(servers) != 2 {
		t.Fatalf("mcpServers=%v", w["mcpServers"])
	}
}

func TestListResultNextPageInfo(t *testing.T) {
	page := cursor.ListResult[string]{NextCursor: "abc"}
	info := page.NextPageInfo()
	if info["cursor"] != "abc" {
		t.Fatalf("info=%v", info)
	}
}

func TestConfigure(t *testing.T) {
	cursor.Configure(cursor.ConfigureOptions{
		Local: &cursor.ConfigureLocalOptions{
			Store: &cursor.LocalAgentStoreConfig{Type: "jsonl", RootDir: "/tmp/store"},
		},
	})
	cursor.Cursor.Configure(cursor.ConfigureOptions{})
}
