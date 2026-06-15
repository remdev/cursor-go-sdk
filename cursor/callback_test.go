package cursor

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

func echoTool() map[string]CustomTool {
	return map[string]CustomTool{
		"echo": {
			Description: "Echo back the input",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"text": map[string]any{"type": "string"},
				},
				"required": []any{"text"},
			},
			Execute: func(args map[string]any, _ CustomToolContext) (any, error) {
				return args["text"], nil
			},
		},
	}
}

func callToolCallback(t *testing.T, srv *ToolCallbackServer, agentID, toolName string, args map[string]any) map[string]any {
	t.Helper()
	body, err := json.Marshal(map[string]any{
		"agentId":  agentID,
		"toolName": toolName,
		"args":     args,
	})
	if err != nil {
		t.Fatal(err)
	}
	req, err := http.NewRequest(http.MethodPost, srv.URL, bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+srv.Token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status=%d", resp.StatusCode)
	}
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	var out map[string]any
	if err := json.Unmarshal(raw, &out); err != nil {
		t.Fatal(err)
	}
	return out
}

func TestToolCallbackServerExecutesRegisteredTool(t *testing.T) {
	srv := NewToolCallbackServer()
	defer srv.Close()

	srv.RegisterAgent("agent-1", echoTool())

	resp := callToolCallback(t, srv, "agent-1", "echo", map[string]any{"text": "pong"})
	if resp["error"] != nil {
		t.Fatalf("unexpected error: %v", resp["error"])
	}
	if resp["result"] != "pong" {
		t.Fatalf("result=%v", resp["result"])
	}
}

func TestToolCallbackServerReturnsNotFoundForUnknownAgent(t *testing.T) {
	srv := NewToolCallbackServer()
	defer srv.Close()

	srv.RegisterAgent("pending", echoTool())

	resp := callToolCallback(t, srv, "real-agent-id", "echo", map[string]any{"text": "pong"})
	errObj, ok := resp["error"].(map[string]any)
	if !ok {
		t.Fatalf("expected error object, got %v", resp)
	}
	if errObj["message"] != "tool not found" {
		t.Fatalf("message=%v", errObj["message"])
	}
}

func TestMigratePendingCustomToolsRegistration(t *testing.T) {
	srv := NewToolCallbackServer()
	defer srv.Close()

	tools := echoTool()
	srv.RegisterAgent("pending", tools)
	srv.RegisterAgent("real-agent-id", tools)
	srv.UnregisterAgent("pending")

	pendingResp := callToolCallback(t, srv, "pending", "echo", map[string]any{"text": "lost"})
	if pendingResp["error"] == nil {
		t.Fatal("expected pending lookup to fail after migration")
	}

	resp := callToolCallback(t, srv, "real-agent-id", "echo", map[string]any{"text": "pong"})
	if resp["result"] != "pong" {
		t.Fatalf("result=%v", resp["result"])
	}
}
