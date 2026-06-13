package cursor

import (
	"context"
	"io"
	"iter"
	"sync"

	conn "github.com/remdev/cursor-go-sdk/internal/connect"
)

var terminalRunStatuses = map[RunStatus]struct{}{
	RunStatusFinished: {}, RunStatusError: {}, RunStatusCancelled: {}, RunStatusExpired: {},
}

// Run is a handle for one prompt submission and its stream.
type Run struct {
	client          *Client
	ID              string
	RequestID       string
	AgentID         string
	Status          RunStatus
	Result          string
	Model           *ModelSelection
	DurationMS      int
	Git             *RunGitInfo
	CreatedAt       string
	apiKey          string
	sendOptions     *SendOptions

	mu              sync.Mutex
	buffer          []RunStreamEvent
	stream          *conn.StreamReader
	terminalResult  *RunResult
	statusListeners []func(RunStatus)
}

func newRunFromSnapshot(c *Client, snap RunResult) *Run {
	return &Run{
		client: c, ID: snap.ID, RequestID: snap.RequestID, AgentID: snap.AgentID, Status: snap.Status,
		Result: snap.Result, Model: snap.Model, DurationMS: snap.DurationMS,
		Git: snap.Git, CreatedAt: snap.CreatedAt,
		terminalResult: &snap,
	}
}

func newRunFromStream(c *Client, agentID string, stream *conn.StreamReader, opts *SendOptions) *Run {
	r := &Run{
		client: c, AgentID: agentID, Status: RunStatusRunning,
		stream: stream, sendOptions: opts,
	}
	r.primeStream(context.Background())
	return r
}

func (r *Run) primeStream(ctx context.Context) {
	ev, ok := r.nextStreamEvent(ctx)
	if !ok {
		return
	}
	r.handleEvent(ev)
	r.buffer = append(r.buffer, ev)
}

// Supports reports whether an operation is available on this run handle.
func (r *Run) Supports(op RunOperation) bool {
	return r.SupportsString(string(op))
}

// SupportsString reports whether a named operation is available.
func (r *Run) SupportsString(op string) bool {
	switch op {
	case string(RunOpStream), string(RunOpWait), string(RunOpCancel), string(RunOpConversation):
		return true
	default:
		return false
	}
}

// UnsupportedReason explains why an operation is unavailable.
func (r *Run) UnsupportedReason(op RunOperation) string {
	return r.UnsupportedReasonString(string(op))
}

// UnsupportedReasonString explains why a named operation is unavailable.
func (r *Run) UnsupportedReasonString(op string) string {
	if r.SupportsString(op) {
		return ""
	}
	return "Run operation " + op + " is not supported"
}

// OnDidChangeStatus registers a listener for status transitions (TypeScript alias).
func (r *Run) OnDidChangeStatus(fn func(RunStatus)) func() {
	return r.OnStatusChange(fn)
}

// OnStatusChange registers a listener for status transitions.
func (r *Run) OnStatusChange(fn func(RunStatus)) func() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.statusListeners = append(r.statusListeners, fn)
	idx := len(r.statusListeners) - 1
	return func() {
		r.mu.Lock()
		defer r.mu.Unlock()
		if idx < len(r.statusListeners) {
			r.statusListeners = append(r.statusListeners[:idx], r.statusListeners[idx+1:]...)
		}
	}
}

// Events returns low-level stream envelopes.
func (r *Run) Events(ctx context.Context) iter.Seq2[RunStreamEvent, error] {
	return func(yield func(RunStreamEvent, error) bool) {
		for _, ev := range r.drainBuffer() {
			if !yield(ev, nil) {
				return
			}
		}
		for {
			ev, ok := r.nextStreamEvent(ctx)
			if !ok {
				return
			}
			r.handleEvent(ev)
			if !yield(ev, nil) {
				return
			}
		}
	}
}

// Messages returns typed SDK messages from the stream.
func (r *Run) Messages(ctx context.Context) iter.Seq2[SDKMessage, error] {
	return func(yield func(SDKMessage, error) bool) {
		for ev, err := range r.Events(ctx) {
			if err != nil {
				yield(SDKMessage{}, err)
				return
			}
			if ev.SDKMessage.Type != "" {
				if !yield(ev.SDKMessage, nil) {
					return
				}
			}
		}
	}
}

// Stream is an alias for Messages.
func (r *Run) Stream(ctx context.Context) iter.Seq2[SDKMessage, error] {
	return r.Messages(ctx)
}

// IterText yields assistant text chunks as they arrive.
func (r *Run) IterText(ctx context.Context) iter.Seq2[string, error] {
	return func(yield func(string, error) bool) {
		for msg, err := range r.Messages(ctx) {
			if err != nil {
				yield("", err)
				return
			}
			if msg.Type != "assistant" {
				continue
			}
			text := AssistantText(msg)
			if text != "" && !yield(text, nil) {
				return
			}
		}
	}
}

// Text blocks until the run completes and returns the final assistant text.
func (r *Run) Text(ctx context.Context) (string, error) {
	if r.Result != "" {
		return r.Result, nil
	}
	res, err := r.Wait(ctx)
	if err != nil {
		return "", err
	}
	return res.Result, nil
}

// Wait blocks until the run reaches a terminal state.
func (r *Run) Wait(ctx context.Context) (RunResult, error) {
	r.mu.Lock()
	if r.terminalResult != nil {
		res := *r.terminalResult
		r.mu.Unlock()
		return res, nil
	}
	r.mu.Unlock()

	for ev, err := range r.Events(ctx) {
		_ = ev
		if err != nil && err != io.EOF {
			return RunResult{}, err
		}
	}
	r.mu.Lock()
	if r.terminalResult != nil {
		res := *r.terminalResult
		r.mu.Unlock()
		return res, nil
	}
	r.mu.Unlock()

	res, err := r.client.waitLiveRun(ctx, r.ID)
	if err != nil {
		return RunResult{}, err
	}
	r.applyResult(res)
	return res, nil
}

// Cancel requests cancellation of an in-flight run.
func (r *Run) Cancel(ctx context.Context) error {
	if _, terminal := terminalRunStatuses[r.Status]; terminal {
		return newUnsupportedRunOperation("cancel",
			"Run is already in terminal status "+string(r.Status)+"; cancel is not applicable.")
	}
	return r.client.cancelRun(ctx, r.ID, r.AgentID)
}

// Observe replays or continues observing a run from an offset.
func (r *Run) Observe(ctx context.Context, afterOffset string) iter.Seq2[RunStreamEvent, error] {
	return func(yield func(RunStreamEvent, error) bool) {
		stream, err := r.client.observeRun(ctx, r.ID, afterOffset)
		if err != nil {
			yield(RunStreamEvent{}, err)
			return
		}
		defer stream.Close()
		for {
			raw, err := stream.Next()
			if err == io.EOF {
				return
			}
			if err != nil {
				yield(RunStreamEvent{}, wrapError(err))
				return
			}
			ev := ParseRunStreamEvent(raw)
			if !yield(ev, nil) {
				return
			}
		}
	}
}

// ConversationJSON returns the raw conversation JSON from the bridge.
func (r *Run) ConversationJSON(ctx context.Context) (string, error) {
	resp, err := r.client.agentUnary(ctx, "GetRunConversation", map[string]any{"runId": r.ID})
	if err != nil {
		return "", err
	}
	return stringField(resp, "conversationJson"), nil
}

// Conversation returns parsed conversation turns.
func (r *Run) Conversation(ctx context.Context) ([]ConversationTurn, error) {
	raw, err := r.ConversationJSON(ctx)
	if err != nil {
		return nil, err
	}
	return ConversationFromJSON(raw)
}

func (r *Run) drainBuffer() []RunStreamEvent {
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(r.buffer) == 0 {
		return nil
	}
	out := r.buffer
	r.buffer = nil
	return out
}

func (r *Run) nextStreamEvent(ctx context.Context) (RunStreamEvent, bool) {
	r.mu.Lock()
	stream := r.stream
	r.mu.Unlock()
	if stream == nil {
		return RunStreamEvent{}, false
	}
	raw, err := stream.Next()
	if err == io.EOF {
		r.mu.Lock()
		r.stream = nil
		r.mu.Unlock()
		return RunStreamEvent{}, false
	}
	if err != nil {
		return RunStreamEvent{}, false
	}
	return ParseRunStreamEvent(raw), true
}

func (r *Run) handleEvent(ev RunStreamEvent) {
	r.applyEventState(ev)
}

func (r *Run) applyEventState(ev RunStreamEvent) {
	if ev.SDKMessage.Type != "" {
		if ev.SDKMessage.AgentID != "" {
			r.AgentID = ev.SDKMessage.AgentID
		}
		if ev.SDKMessage.RunID != "" {
			r.ID = ev.SDKMessage.RunID
		}
	}
	if ev.Kind == "interaction_update" && r.sendOptions != nil && r.sendOptions.OnDelta != nil {
		r.sendOptions.OnDelta(ev.InteractionUpdate)
	}
	if ev.Kind == "step" && r.sendOptions != nil && r.sendOptions.OnStep != nil {
		r.sendOptions.OnStep(ev.Step)
	}
	if ev.Kind == "result" && ev.Result != nil {
		res := parseRunResult(ev.Result)
		r.applyResult(res)
	}
}

func (r *Run) applyResult(res RunResult) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if res.ID != "" {
		r.ID = res.ID
	}
	if res.RequestID != "" {
		r.RequestID = res.RequestID
	}
	if res.AgentID != "" {
		r.AgentID = res.AgentID
	}
	r.Result = res.Result
	r.Model = res.Model
	r.DurationMS = res.DurationMS
	r.Git = res.Git
	r.CreatedAt = res.CreatedAt
	old := r.Status
	r.Status = res.Status
	r.terminalResult = &res
	if old != res.Status {
		for _, fn := range r.statusListeners {
			fn(res.Status)
		}
	}
}

// GetRun fetches a run by ID.
func GetRun(ctx context.Context, runID string, opts GetRunOptions, clientOpts ...ClientOption) (*Run, error) {
	c, err := defaultOrLaunchClient(ctx, clientOpts...)
	if err != nil {
		return nil, err
	}
	return c.getRun(ctx, runID, opts.ToWire())
}

// ListRuns lists runs for an agent.
func ListRuns(ctx context.Context, agentID string, opts ListRunsOptions, clientOpts ...ClientOption) (ListResult[*Run], error) {
	c, err := defaultOrLaunchClient(ctx, clientOpts...)
	if err != nil {
		return ListResult[*Run]{}, err
	}
	return c.listRuns(ctx, agentID, opts.ToWire())
}

// ArchiveAgent archives a cloud agent by ID.
func ArchiveAgent(ctx context.Context, agentID string, opts AgentOperationOptions, clientOpts ...ClientOption) error {
	c, err := defaultOrLaunchClient(ctx, clientOpts...)
	if err != nil {
		return err
	}
	_, err = c.agentUnary(ctx, "ArchiveAgent", map[string]any{
		"agentId": agentID, "options": opts.ToWire(),
	})
	return err
}

// UnarchiveAgent restores an archived cloud agent.
func UnarchiveAgent(ctx context.Context, agentID string, opts AgentOperationOptions, clientOpts ...ClientOption) error {
	c, err := defaultOrLaunchClient(ctx, clientOpts...)
	if err != nil {
		return err
	}
	_, err = c.agentUnary(ctx, "UnarchiveAgent", map[string]any{
		"agentId": agentID, "options": opts.ToWire(),
	})
	return err
}

// DeleteAgent permanently deletes a cloud agent.
func DeleteAgent(ctx context.Context, agentID string, opts AgentOperationOptions, clientOpts ...ClientOption) error {
	c, err := defaultOrLaunchClient(ctx, clientOpts...)
	if err != nil {
		return err
	}
	_, err = c.agentUnary(ctx, "DeleteAgent", map[string]any{
		"agentId": agentID, "options": opts.ToWire(),
	})
	return err
}

// CancelRun cancels a run by ID.
func CancelRun(ctx context.Context, runID, agentID string, clientOpts ...ClientOption) error {
	c, err := defaultOrLaunchClient(ctx, clientOpts...)
	if err != nil {
		return err
	}
	return c.cancelRun(ctx, runID, agentID)
}
