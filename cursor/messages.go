package cursor

import "encoding/json"

// TextBlock is assistant/user text content.
type TextBlock struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// ToolUseBlock is a tool invocation in assistant content.
type ToolUseBlock struct {
	Type  string `json:"type"`
	ID    string `json:"id"`
	Name  string `json:"name"`
	Input any    `json:"input,omitempty"`
}

// SDKMessage is a typed stream message during a run.
type SDKMessage struct {
	Type               string
	AgentID            string
	RunID              string
	Subtype            string
	Model              *ModelSelection
	Tools              []string
	UserContent        []TextBlock
	AssistantContent   []any
	ThinkingText       string
	ThinkingDurationMS *int
	ToolCallID         string
	ToolName           string
	ToolStatus         string
	ToolArgs           any
	ToolResult         any
	Status             string
	StatusMessage      string
	TaskStatus         string
	TaskText           string
	RequestID          string
	Raw                map[string]any
}

// SDKMessageFromJSON parses a wire SDK message.
func SDKMessageFromJSON(m map[string]any) SDKMessage {
	payload := m
	if inner, ok := m["message"].(map[string]any); ok {
		if _, hasType := inner["type"]; hasType {
			payload = inner
		}
	}
	msgType := stringField(payload, "type")
	if msgType == "" {
		msgType = stringField(m, "type")
	}
	agentID := stringField(payload, "agent_id", "agentId")
	runID := stringField(payload, "run_id", "runId")
	out := SDKMessage{
		Type: msgType, AgentID: agentID, RunID: runID, Raw: payload,
	}
	switch msgType {
	case "system":
		out.Subtype = stringField(payload, "subtype")
		if model, ok := payload["model"].(map[string]any); ok {
			ms := parseModelSelection(model)
			out.Model = &ms
		}
		if tools, ok := payload["tools"].([]any); ok {
			for _, t := range tools {
				out.Tools = append(out.Tools, fmtString(t))
			}
		}
	case "user":
		if content, ok := payload["message"].(map[string]any); ok {
			out.UserContent = parseTextBlocks(content["content"])
		}
	case "assistant":
		if content, ok := payload["message"].(map[string]any); ok {
			out.AssistantContent = parseContentBlocks(content["content"])
		}
	case "thinking":
		out.ThinkingText = stringField(payload, "text")
		if v := intField(payload, "thinkingDurationMs"); v > 0 {
			out.ThinkingDurationMS = &v
		}
	case "tool_call":
		out.ToolCallID = stringField(payload, "call_id", "callId")
		out.ToolName = stringField(payload, "name")
		out.ToolStatus = stringField(payload, "status")
		out.ToolArgs = payload["args"]
		out.ToolResult = payload["result"]
	case "status":
		out.Status = stringField(payload, "status")
		out.StatusMessage = stringField(payload, "message")
	case "task":
		out.TaskStatus = stringField(payload, "status")
		out.TaskText = stringField(payload, "text")
	case "request":
		out.RequestID = stringField(payload, "request_id", "requestId")
	}
	return out
}

func parseTextBlocks(v any) []TextBlock {
	raw, ok := v.([]any)
	if !ok {
		return nil
	}
	out := make([]TextBlock, 0, len(raw))
	for _, item := range raw {
		if m, ok := item.(map[string]any); ok && stringField(m, "type") == "text" {
			out = append(out, TextBlock{Type: "text", Text: stringField(m, "text")})
		}
	}
	return out
}

func parseContentBlocks(v any) []any {
	raw, ok := v.([]any)
	if !ok {
		return nil
	}
	out := make([]any, 0, len(raw))
	for _, item := range raw {
		m, ok := item.(map[string]any)
		if !ok {
			continue
		}
		switch stringField(m, "type") {
		case "text":
			out = append(out, TextBlock{Type: "text", Text: stringField(m, "text")})
		case "tool_use":
			out = append(out, ToolUseBlock{
				Type: "tool_use", ID: stringField(m, "id"),
				Name: stringField(m, "name"), Input: m["input"],
			})
		default:
			out = append(out, m)
		}
	}
	return out
}

// InteractionUpdate is a delta event during streaming.
type InteractionUpdate struct {
	Type               string
	Text               string
	ThinkingDurationMS int
	CallID             string
	ToolCall           map[string]any
	ModelCallID        string
	Tokens             int
	StepID             int
	StepDurationMS     int
	Usage              map[string]any
	UserMessage        map[string]any
	Summary            string
	Event              map[string]any
	Raw                map[string]any
}

// ParseInteractionUpdate parses a wire interaction update.
func ParseInteractionUpdate(m map[string]any) InteractionUpdate {
	payload := m
	if inner, ok := m["update"].(map[string]any); ok {
		payload = inner
	}
	t := stringField(payload, "type")
	if t == "" {
		t = stringField(m, "type")
	}
	u := InteractionUpdate{Type: t, Raw: payload}
	switch t {
	case "text-delta":
		u.Text = stringField(payload, "text")
	case "thinking-delta":
		u.Text = stringField(payload, "text")
	case "thinking-completed":
		u.ThinkingDurationMS = intField(payload, "thinkingDurationMs")
	case "tool-call-started", "tool-call-completed", "partial-tool-call":
		u.CallID = stringField(payload, "callId")
		if tc, ok := payload["toolCall"].(map[string]any); ok {
			u.ToolCall = tc
		}
		u.ModelCallID = stringField(payload, "modelCallId")
	case "token-delta":
		u.Tokens = intField(payload, "tokens")
	case "step-started":
		u.StepID = intField(payload, "stepId")
	case "step-completed":
		u.StepID = intField(payload, "stepId")
		u.StepDurationMS = intField(payload, "stepDurationMs")
	case "turn-ended":
		if usage, ok := payload["usage"].(map[string]any); ok {
			u.Usage = usage
		}
	case "user-message-appended":
		if um, ok := payload["userMessage"].(map[string]any); ok {
			u.UserMessage = um
		}
	case "summary":
		u.Summary = stringField(payload, "summary")
	case "shell-output-delta":
		if ev, ok := payload["event"].(map[string]any); ok {
			u.Event = ev
		}
	}
	return u
}

// ConversationStep is one step in a run conversation.
type ConversationStep struct {
	Type    string
	Text    string
	Message map[string]any
	Raw     map[string]any
}

// ParseConversationStep parses a conversation step envelope.
func ParseConversationStep(m map[string]any) ConversationStep {
	payload := m
	if step, ok := m["step"].(map[string]any); ok {
		payload = step
	}
	t := stringField(payload, "type")
	step := ConversationStep{Type: t, Raw: payload}
	if msg, ok := payload["message"].(map[string]any); ok {
		switch t {
		case "assistantMessage":
			step.Text = stringField(msg, "text")
		case "toolCall":
			step.Message = msg
		case "thinkingMessage":
			step.Text = stringField(msg, "text")
		}
	}
	return step
}

// ConversationTurn is a turn in run conversation history.
type ConversationTurn struct {
	Type string
	Turn any
}

// ParseConversationTurn parses a conversation turn from JSON bytes.
func ParseConversationTurn(m map[string]any) ConversationTurn {
	t := stringField(m, "type")
	turnRaw, _ := m["turn"].(map[string]any)
	return ConversationTurn{Type: t, Turn: turnRaw}
}

// ParseRunStreamEvent parses a stream envelope.
func ParseRunStreamEvent(m map[string]any) RunStreamEvent {
	ev := RunStreamEvent{Offset: stringField(m, "offset")}
	if sm, ok := m["sdkMessage"].(map[string]any); ok {
		ev.Kind = "sdk_message"
		msg := SDKMessageFromJSON(sm)
		ev.SDKMessage = msg
		return ev
	}
	if iu, ok := m["interactionUpdate"].(map[string]any); ok {
		ev.Kind = "interaction_update"
		ev.InteractionUpdate = ParseInteractionUpdate(iu)
		return ev
	}
	if step, ok := m["step"].(map[string]any); ok {
		ev.Kind = "step"
		ev.Step = ParseConversationStep(step)
		return ev
	}
	if result, ok := m["result"].(map[string]any); ok {
		ev.Kind = "result"
		if inner, ok := result["result"].(map[string]any); ok {
			ev.Result = inner
			ev.ResultIsFull = true
		} else {
			ev.Result = result
		}
		return ev
	}
	if done, ok := m["done"].(map[string]any); ok {
		ev.Kind = "done"
		ev.Done = done
		return ev
	}
	ev.Kind = "unknown"
	return ev
}

// AssistantText extracts text from an assistant SDKMessage.
func AssistantText(msg SDKMessage) string {
	var b []byte
	for _, block := range msg.AssistantContent {
		switch t := block.(type) {
		case TextBlock:
			b = append(b, t.Text...)
		case map[string]any:
			if stringField(t, "type") == "text" {
				b = append(b, stringField(t, "text")...)
			}
		}
	}
	return string(b)
}

// ConversationFromJSON parses conversation JSON returned by GetRunConversation.
func ConversationFromJSON(raw string) ([]ConversationTurn, error) {
	if raw == "" {
		return nil, nil
	}
	var loaded []map[string]any
	if err := json.Unmarshal([]byte(raw), &loaded); err != nil {
		return nil, err
	}
	out := make([]ConversationTurn, 0, len(loaded))
	for _, item := range loaded {
		out = append(out, ParseConversationTurn(item))
	}
	return out, nil
}
