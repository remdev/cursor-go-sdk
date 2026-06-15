package cursor

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateAgentReregistersCustomToolsFromPending(t *testing.T) {
	toolSrv := NewToolCallbackServer()
	defer toolSrv.Close()

	const realAgentID = "agent-real-123"
	bridge := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/sdk.v1.SdkAgentService/CreateAgent" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"agentId": realAgentID})
	}))
	defer bridge.Close()

	c, err := Connect(bridge.URL, "test-token", WithAllowAPIKeyEnvFallback(true))
	if err != nil {
		t.Fatal(err)
	}
	c.toolCallback = toolSrv

	agent, err := c.createAgent(context.Background(), AgentOptions{
		Model:  "composer-2.5",
		APIKey: "test-key",
		Local: &LocalAgentOptions{
			CWD:         []string{"."},
			CustomTools: echoTool(),
		},
	}, "")
	if err != nil {
		t.Fatal(err)
	}
	if agent.AgentID != realAgentID {
		t.Fatalf("agentId=%q", agent.AgentID)
	}

	resp := callToolCallback(t, toolSrv, realAgentID, "echo", map[string]any{"text": "pong"})
	if resp["result"] != "pong" {
		t.Fatalf("result=%v", resp["result"])
	}

	pendingResp := callToolCallback(t, toolSrv, "pending", "echo", map[string]any{"text": "lost"})
	if pendingResp["error"] == nil {
		t.Fatal("expected pending registration to be removed")
	}
}

func TestCreateAgentKeepsExplicitAgentIDRegistration(t *testing.T) {
	toolSrv := NewToolCallbackServer()
	defer toolSrv.Close()

	const explicitAgentID = "agent-explicit"
	bridge := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"agentId": explicitAgentID})
	}))
	defer bridge.Close()

	c, err := Connect(bridge.URL, "test-token", WithAllowAPIKeyEnvFallback(true))
	if err != nil {
		t.Fatal(err)
	}
	c.toolCallback = toolSrv

	_, err = c.createAgent(context.Background(), AgentOptions{
		Model:   "composer-2.5",
		APIKey:  "test-key",
		AgentID: explicitAgentID,
		Local: &LocalAgentOptions{
			CWD:         []string{"."},
			CustomTools: echoTool(),
		},
	}, "")
	if err != nil {
		t.Fatal(err)
	}

	resp := callToolCallback(t, toolSrv, explicitAgentID, "echo", map[string]any{"text": "ok"})
	if resp["result"] != "ok" {
		t.Fatalf("result=%v", resp["result"])
	}
}

func TestCreateAgentCreateFailureUnregistersPending(t *testing.T) {
	toolSrv := NewToolCallbackServer()
	defer toolSrv.Close()

	bridge := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `{"code":"internal","message":"boom"}`, http.StatusInternalServerError)
	}))
	defer bridge.Close()

	c, err := Connect(bridge.URL, "test-token", WithAllowAPIKeyEnvFallback(true))
	if err != nil {
		t.Fatal(err)
	}
	c.toolCallback = toolSrv

	_, err = c.createAgent(context.Background(), AgentOptions{
		Model:  "composer-2.5",
		APIKey: "test-key",
		Local: &LocalAgentOptions{
			CWD:         []string{"."},
			CustomTools: echoTool(),
		},
	}, "")
	if err == nil {
		t.Fatal("expected create failure")
	}

	toolSrv.mu.Lock()
	_, pendingExists := toolSrv.agents["pending"]
	toolSrv.mu.Unlock()
	if pendingExists {
		t.Fatal("pending registration should be cleaned up on create failure")
	}
}

func TestResumeAgentRegistersCustomTools(t *testing.T) {
	toolSrv := NewToolCallbackServer()
	defer toolSrv.Close()

	const agentID = "agent-resumed"
	bridge := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/sdk.v1.SdkAgentService/ResumeAgent" {
			http.NotFound(w, r)
			return
		}
		body, _ := io.ReadAll(r.Body)
		var req map[string]any
		_ = json.Unmarshal(body, &req)
		if req["agentId"] != agentID {
			t.Errorf("agentId=%v", req["agentId"])
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{"agentId": agentID})
	}))
	defer bridge.Close()

	c, err := Connect(bridge.URL, "test-token", WithAllowAPIKeyEnvFallback(true))
	if err != nil {
		t.Fatal(err)
	}
	c.toolCallback = toolSrv

	_, err = c.resumeAgent(context.Background(), agentID, AgentOptions{
		Model:  "composer-2.5",
		APIKey: "test-key",
		Local: &LocalAgentOptions{
			CWD:         []string{"."},
			CustomTools: echoTool(),
		},
	})
	if err != nil {
		t.Fatal(err)
	}

	resp := callToolCallback(t, toolSrv, agentID, "echo", map[string]any{"text": "resume"})
	if resp["result"] != "resume" {
		t.Fatalf("result=%v", resp["result"])
	}
}
