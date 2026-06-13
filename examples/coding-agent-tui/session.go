package main

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/remdev/cursor-go-sdk/cursor"
	"github.com/remdev/cursor-go-sdk/examples/agentutil"
)

var (
	ErrNoActiveRun      = errors.New("no active run to cancel")
	ErrRunNotCancelable = errors.New("this run cannot be cancelled")
)

const agentInstructions = `You are a lightweight coding agent running from a terminal.
Work in the configured workspace.
Help the user inspect, edit, and validate code with small focused changes.
Before changing files, understand the surrounding code and preserve unrelated user work.
Keep progress updates concise and summarize the result clearly.`

type executionMode string

const (
	modeLocal executionMode = "local"
	modeCloud executionMode = "cloud"
)

type eventKind int

const (
	eventAssistantDelta eventKind = iota
	eventThinking
	eventTool
	eventStatus
	eventTask
	eventResult
	eventError
)

type modelChoice struct {
	label string
	value cursor.ModelSelection
}

type agentEvent struct {
	kind          eventKind
	text          string
	toolName      string
	toolStatus    string
	toolParams    string
	status        string
	statusMessage string
	taskStatus    string
	taskText      string
	resultStatus  string
	durationMS    int
	err           error
}

type session struct {
	apiKey    string
	cwd       string
	force     bool
	mode      executionMode
	model     cursor.ModelSelection
	cloudRepo *cloudRepository

	mu         sync.Mutex
	agent      *cursor.Agent
	currentRun *cursor.Run
	agentKey   string
}

func newSession(apiKey, cwd, modelID string, force bool) (*session, error) {
	s := &session{
		apiKey: apiKey,
		cwd:    cwd,
		force:  force,
		mode:   modeLocal,
		model:  cursor.ModelFromString(modelID),
	}
	if err := s.replaceAgent(context.Background()); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *session) executionTarget() string {
	if s.mode == modeLocal {
		return s.cwd
	}
	if s.cloudRepo != nil {
		return formatCloudRepository(*s.cloudRepo)
	}
	return "(cloud)"
}

func (s *session) modelLabel() string {
	return formatModelLabel(s.model)
}

func (s *session) listModels(ctx context.Context) ([]modelChoice, error) {
	models, err := cursor.Models.List(ctx, cursor.CursorRequestOptions{APIKey: s.apiKey})
	if err != nil {
		return nil, err
	}
	choices := make([]modelChoice, 0)
	seen := map[string]struct{}{}
	for _, m := range models {
		for _, c := range modelToChoices(m) {
			key := modelSelectionKey(c.value)
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}
			choices = append(choices, c)
		}
	}
	if len(choices) == 0 {
		return []modelChoice{{label: s.model.ID, value: s.model}}, nil
	}
	return choices, nil
}

func (s *session) setModel(model cursor.ModelSelection) {
	s.model = model
}

func (s *session) reset(ctx context.Context) error {
	return s.replaceAgent(ctx)
}

func (s *session) setExecutionMode(ctx context.Context, mode executionMode) error {
	s.mu.Lock()
	if s.currentRun != nil {
		s.mu.Unlock()
		return fmt.Errorf("wait for the current run to finish before switching execution mode")
	}
	if s.mode == mode {
		s.mu.Unlock()
		return nil
	}
	prev := s.mode
	s.mode = mode
	s.mu.Unlock()

	if err := s.replaceAgent(ctx); err != nil {
		s.mu.Lock()
		s.mode = prev
		s.mu.Unlock()
		return err
	}
	return nil
}

func (s *session) close(ctx context.Context) {
	s.mu.Lock()
	agent := s.agent
	s.agent = nil
	s.mu.Unlock()
	if agent != nil {
		_ = agent.Close(ctx)
	}
}

func (s *session) cancelCurrentRun(ctx context.Context) error {
	s.mu.Lock()
	run := s.currentRun
	s.mu.Unlock()
	if run == nil {
		return ErrNoActiveRun
	}
	if !run.Supports(cursor.RunOpCancel) {
		return ErrRunNotCancelable
	}
	return run.Cancel(ctx)
}

func (s *session) sendPrompt(ctx context.Context, prompt string, onEvent func(agentEvent)) error {
	if err := s.ensureAgentFresh(ctx); err != nil {
		return err
	}

	s.mu.Lock()
	agent := s.agent
	s.mu.Unlock()

	emitMessage := func(msg cursor.SDKMessage) {
		switch msg.Type {
		case "assistant":
			text := cursor.AssistantText(msg)
			if text != "" {
				onEvent(agentEvent{kind: eventAssistantDelta, text: text})
			}
		case "thinking":
			if t := strings.TrimSpace(msg.ThinkingText); t != "" {
				onEvent(agentEvent{kind: eventThinking, text: agentutil.Compact(t)})
			}
		case "tool":
			onEvent(agentEvent{
				kind:       eventTool,
				toolName:   msg.ToolName,
				toolStatus: msg.ToolStatus,
				toolParams: summarizeTool(msg.ToolName, msg.ToolArgs),
			})
		case "status":
			if msg.Status != "" && msg.Status != "FINISHED" {
				onEvent(agentEvent{
					kind:          eventStatus,
					status:        msg.Status,
					statusMessage: msg.StatusMessage,
				})
			}
		case "task":
			onEvent(agentEvent{
				kind:       eventTask,
				taskStatus: msg.TaskStatus,
				taskText:   msg.TaskText,
			})
		}
	}

	sendOpts := cursor.SendOptions{
		OnDelta: func(u cursor.InteractionUpdate) {
			switch u.Type {
			case "text-delta":
				if u.Text != "" {
					onEvent(agentEvent{kind: eventAssistantDelta, text: u.Text})
				}
			case "thinking-delta":
				if t := strings.TrimSpace(u.Text); t != "" {
					onEvent(agentEvent{kind: eventThinking, text: agentutil.Compact(t)})
				}
			}
		},
	}
	if s.mode == modeLocal {
		sendOpts.Model = &s.model
		if s.force {
			force := true
			sendOpts.Local = &cursor.LocalSendOptions{Force: &force}
		}
	}

	run, err := agent.Send(ctx, buildPrompt(prompt), sendOpts)
	if err != nil {
		return err
	}

	s.mu.Lock()
	s.currentRun = run
	s.mu.Unlock()

	defer func() {
		s.mu.Lock()
		if s.currentRun == run {
			s.currentRun = nil
		}
		s.mu.Unlock()
	}()

	for msg, err := range run.Messages(ctx) {
		if err != nil {
			onEvent(agentEvent{kind: eventError, err: err})
			return err
		}
		emitMessage(msg)
	}

	result, err := run.Wait(ctx)
	if err != nil {
		onEvent(agentEvent{kind: eventError, err: err})
		return err
	}
	onEvent(agentEvent{
		kind:         eventResult,
		resultStatus: string(result.Status),
		durationMS:   result.DurationMS,
	})
	return nil
}

func (s *session) replaceAgent(ctx context.Context) error {
	s.mu.Lock()
	prev := s.agent
	s.mu.Unlock()

	agent, cloudRepo, err := s.createAgent(ctx)
	if err != nil {
		return err
	}

	s.mu.Lock()
	s.agent = agent
	s.cloudRepo = cloudRepo
	s.agentKey = s.currentAgentKey()
	s.mu.Unlock()

	if prev != nil {
		_ = prev.Close(ctx)
	}
	return nil
}

func (s *session) ensureAgentFresh(ctx context.Context) error {
	s.mu.Lock()
	key := s.currentAgentKey()
	stale := s.agentKey != key
	s.mu.Unlock()
	if stale {
		return s.replaceAgent(ctx)
	}
	return nil
}

func (s *session) createAgent(ctx context.Context) (*cursor.Agent, *cloudRepository, error) {
	opts := cursor.AgentOptions{
		APIKey: s.apiKey,
		Name:   "Lightweight coding agent",
		Model:  s.model,
	}

	var cloudRepo *cloudRepository
	if s.mode == modeCloud {
		repo, err := detectCloudRepository(s.cwd)
		if err != nil {
			return nil, nil, err
		}
		cloudRepo = &repo
		opts.Cloud = &cursor.CloudAgentOptions{
			Repos: []cursor.CloudRepository{{
				URL:         repo.URL,
				StartingRef: repo.StartingRef,
			}},
		}
	} else {
		opts.Local = &cursor.LocalAgentOptions{CWD: []string{s.cwd}}
	}

	agent, err := cursor.CreateAgent(ctx, opts)
	if err != nil {
		return nil, nil, err
	}
	return agent, cloudRepo, nil
}

func (s *session) currentAgentKey() string {
	if s.mode == modeCloud {
		return fmt.Sprintf("cloud:%s", modelSelectionKey(s.model))
	}
	return fmt.Sprintf("local:%s", modelSelectionKey(s.model))
}

func buildPrompt(prompt string) string {
	return agentInstructions + "\n\nUser task:\n" + prompt
}

func formatModelLabel(m cursor.ModelSelection) string {
	if len(m.Params) == 0 {
		return m.ID
	}
	parts := make([]string, 0, len(m.Params))
	for _, p := range m.Params {
		if p.Value != "" {
			parts = append(parts, p.Value)
		}
	}
	if len(parts) == 0 {
		return m.ID
	}
	return m.ID + " (" + strings.Join(parts, ", ") + ")"
}

func modelSelectionKey(m cursor.ModelSelection) string {
	if len(m.Params) == 0 {
		return m.ID
	}
	parts := make([]string, 0, len(m.Params))
	for _, p := range m.Params {
		parts = append(parts, p.ID+"="+p.Value)
	}
	return m.ID + "?" + strings.Join(parts, "&")
}

func modelToChoices(m cursor.SDKModel) []modelChoice {
	base := m.DisplayName
	if base == "" {
		base = m.ID
	}
	if len(m.Variants) == 0 {
		return []modelChoice{{label: base, value: cursor.ModelFromString(m.ID)}}
	}
	out := make([]modelChoice, 0, len(m.Variants))
	for _, v := range m.Variants {
		out = append(out, modelChoice{
			label: buildVariantLabel(base, v.DisplayName),
			value: cursor.ModelSelection{ID: m.ID, Params: append([]cursor.ModelParameterValue(nil), v.Params...)},
		})
	}
	return out
}

func buildVariantLabel(base, variantDisplayName string) string {
	variantLabel := strings.TrimSpace(variantDisplayName)
	if variantLabel == "" || strings.EqualFold(base, variantLabel) {
		return base
	}
	return base + " - " + variantLabel
}

func summarizeTool(name string, args any) string {
	record, ok := args.(map[string]any)
	if !ok || len(record) == 0 {
		return ""
	}
	keys := toolSummaryKeys(name)
	parts := make([]string, 0, 2)
	for _, group := range keys {
		for _, key := range group {
			if v, ok := record[key]; ok {
				if s := formatArgValue(v); s != "" {
					parts = append(parts, key+"="+s)
					break
				}
			}
		}
	}
	return strings.Join(parts, " ")
}

func toolSummaryKeys(toolName string) [][]string {
	name := strings.ToLower(toolName)
	switch {
	case strings.Contains(name, "read"):
		return [][]string{{"path", "filePath", "target_file"}, {"offset"}, {"limit"}}
	case strings.Contains(name, "glob"):
		return [][]string{{"pattern", "glob"}, {"path", "cwd"}}
	case strings.Contains(name, "grep"), strings.Contains(name, "search"):
		return [][]string{{"pattern", "query"}, {"path"}, {"glob"}}
	case strings.Contains(name, "shell"), strings.Contains(name, "terminal"), strings.Contains(name, "command"):
		return [][]string{{"command", "cmd"}, {"cwd"}}
	case strings.Contains(name, "edit"), strings.Contains(name, "write"), strings.Contains(name, "patch"):
		return [][]string{{"path", "target_file", "file"}, {"instruction"}}
	default:
		return [][]string{{"path", "file"}, {"pattern", "query", "command"}}
	}
}

func formatArgValue(v any) string {
	switch t := v.(type) {
	case string:
		return shorten(strings.Join(strings.Fields(t), " "))
	case float64:
		return fmt.Sprintf("%v", t)
	case bool:
		return fmt.Sprintf("%v", t)
	case []any:
		items := make([]string, 0, 3)
		for i, item := range t {
			if i >= 3 {
				break
			}
			if s := formatArgValue(item); s != "" {
				items = append(items, s)
			}
		}
		if len(items) > 0 {
			return "[" + strings.Join(items, ",") + "]"
		}
	}
	return ""
}

func shorten(s string) string {
	const max = 80
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}
