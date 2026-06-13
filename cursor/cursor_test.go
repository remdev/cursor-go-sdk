package cursor_test

import (
	"testing"

	"github.com/remdev/cursor-go-sdk/cursor"
)

func TestModelSelectionWire(t *testing.T) {
	m := cursor.ModelSelection{
		ID: "composer-2.5",
		Params: []cursor.ModelParameterValue{{ID: "thinking", Value: "high"}},
	}
	w := m.ToWire()
	if w["id"] != "composer-2.5" {
		t.Fatalf("unexpected id: %v", w["id"])
	}
}

func TestUserMessageWire(t *testing.T) {
	msg := cursor.UserMessageFromText("hello")
	w := msg.ToWire()
	if w["text"] != "hello" {
		t.Fatalf("unexpected text: %v", w["text"])
	}
}

func TestParseRunStreamEvent(t *testing.T) {
	ev := cursor.ParseRunStreamEvent(map[string]any{
		"sdkMessage": map[string]any{
			"type":    "assistant",
			"agentId": "agent-1",
			"runId":   "run-1",
			"message": map[string]any{
				"role": "assistant",
				"content": []any{
					map[string]any{"type": "text", "text": "hi"},
				},
			},
		},
	})
	if ev.Kind != "sdk_message" {
		t.Fatalf("kind=%q", ev.Kind)
	}
	if cursor.AssistantText(ev.SDKMessage) != "hi" {
		t.Fatalf("text=%q", cursor.AssistantText(ev.SDKMessage))
	}
}

func TestAgentOptionsWire(t *testing.T) {
	autoReview := true
	opts := cursor.AgentOptions{
		Model:  "composer-2.5",
		APIKey: "test-key",
		Local: &cursor.LocalAgentOptions{
			CWD:        []string{"/tmp"},
			AutoReview: &autoReview,
		},
	}
	w := opts.ToWire()
	if w["apiKey"] != "test-key" {
		t.Fatalf("apiKey=%v", w["apiKey"])
	}
	local, ok := w["local"].(map[string]any)
	if !ok {
		t.Fatal("missing local")
	}
	cwd, ok := local["cwd"].([]string)
	if !ok || len(cwd) != 1 || cwd[0] != "/tmp" {
		t.Fatalf("cwd=%v", local["cwd"])
	}
}
