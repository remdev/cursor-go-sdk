package cursor

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
	"sync/atomic"
)

const toolCallbackService = "sdk.v1.SdkCustomToolCallbackService"

// ToolCallbackServer hosts custom tool execute handlers for the bridge.
type ToolCallbackServer struct {
	URL    string
	Token  string
	server *http.Server
	mu     sync.Mutex
	agents map[string]map[string]CustomTool
	nextID atomic.Uint64
}

// NewToolCallbackServer starts a loopback Connect/JSON server for custom tools.
func NewToolCallbackServer() *ToolCallbackServer {
	mux := http.NewServeMux()
	s := &ToolCallbackServer{
		Token:  randomUUID(),
		agents: make(map[string]map[string]CustomTool),
	}
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return s
	}
	s.URL = fmt.Sprintf("http://%s", ln.Addr().String())
	mux.HandleFunc("/", s.handle)
	s.server = &http.Server{Handler: mux}
	go s.server.Serve(ln)
	return s
}

// RegisterAgent registers custom tools for an agent ID.
func (s *ToolCallbackServer) RegisterAgent(agentID string, tools map[string]CustomTool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	copyTools := make(map[string]CustomTool, len(tools))
	for k, v := range tools {
		copyTools[k] = v
	}
	s.agents[agentID] = copyTools
}

// UnregisterAgent removes custom tools for an agent ID.
func (s *ToolCallbackServer) UnregisterAgent(agentID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.agents, agentID)
}

// Close stops the callback server.
func (s *ToolCallbackServer) Close() {
	if s.server != nil {
		s.server.Close()
	}
}

func (s *ToolCallbackServer) handle(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Authorization") != "Bearer "+s.Token {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var req map[string]any
	if err := json.Unmarshal(body, &req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Execute RPC based on path: /sdk.v1.SdkCustomToolCallbackService/Execute
	agentID := stringField(req, "agentId")
	toolName := stringField(req, "toolName")
	args, _ := req["args"].(map[string]any)
	toolCallID := stringField(req, "toolCallId")

	s.mu.Lock()
	tools := s.agents[agentID]
	var tool CustomTool
	if tools != nil {
		tool = tools[toolName]
	}
	s.mu.Unlock()

	resp := map[string]any{}
	if tool.Execute == nil {
		resp["error"] = map[string]any{"message": "tool not found"}
	} else {
		result, err := tool.Execute(args, CustomToolContext{ToolCallID: toolCallID})
		if err != nil {
			resp["error"] = map[string]any{"message": err.Error()}
		} else {
			resp["result"] = result
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// unused but documents service name for bridge integration
var _ = toolCallbackService
