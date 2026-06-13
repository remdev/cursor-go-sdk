package cursor

import (
	"context"
	"crypto/rand"
	"fmt"
)

// Agent is a durable handle for multi-turn agent conversations.
type Agent struct {
	client  *Client
	AgentID string
	Model   *ModelSelection
	apiKey  string
}

// CreateAgent creates a new agent with the given options.
func CreateAgent(ctx context.Context, opts AgentOptions, clientOpts ...ClientOption) (*Agent, error) {
	c, err := defaultOrLaunchClient(ctx, clientOpts...)
	if err != nil {
		return nil, err
	}
	return c.createAgent(ctx, opts, opts.IdempotencyKey)
}

// ResumeAgent resumes an existing agent by ID.
func ResumeAgent(ctx context.Context, agentID string, opts AgentOptions, clientOpts ...ClientOption) (*Agent, error) {
	c, err := defaultOrLaunchClient(ctx, clientOpts...)
	if err != nil {
		return nil, err
	}
	return c.resumeAgent(ctx, agentID, opts)
}

// Prompt is a one-shot helper: create agent, send, wait, and close.
func Prompt(ctx context.Context, message any, opts AgentOptions, clientOpts ...ClientOption) (RunResult, error) {
	agent, err := CreateAgent(ctx, opts, clientOpts...)
	if err != nil {
		return RunResult{}, err
	}
	defer agent.Close(context.Background())
	run, err := agent.Send(ctx, message, SendOptions{})
	if err != nil {
		return RunResult{}, err
	}
	return run.Wait(ctx)
}

// Close closes the agent on the bridge.
func (a *Agent) Close(ctx context.Context) error {
	if a.client.toolCallback != nil {
		a.client.toolCallback.UnregisterAgent(a.AgentID)
	}
	_, err := a.client.agentUnary(ctx, "CloseAgent", map[string]any{"agentId": a.AgentID})
	return err
}

// Send starts a new run for this agent.
func (a *Agent) Send(ctx context.Context, message any, opts SendOptions) (*Run, error) {
	if !a.client.allowAPIKeyEnvFallback && stringsHasPrefix(a.AgentID, "bc-") && a.apiKey == "" {
		return nil, configurationError(
			"Cloud agent send through a caller-supplied bridge requires an explicit api_key",
			"missing_api_key",
		)
	}
	wireOpts := opts.ToWire()
	if a.apiKey != "" && wireOpts["apiKey"] == nil {
		wireOpts = copyMap(wireOpts)
		// apiKey goes in request options via optionsWithAPIKey pattern
	}
	req := map[string]any{
		"agentId": a.AgentID,
		"message": parseUserMessage(message).ToWire(),
		"options": optionsWithAPIKey(a.apiKey, wireOpts),
	}
	if opts.IdempotencyKey != "" {
		req["idempotencyKey"] = opts.IdempotencyKey
	}
	stream, err := a.client.agentStream(ctx, "Send", req)
	if err != nil {
		return nil, err
	}
	run := newRunFromStream(a.client, a.AgentID, stream, &opts)
	run.apiKey = a.apiKey
	if opts.Local != nil && len(opts.Local.CustomTools) > 0 && a.client.toolCallback != nil {
		a.client.toolCallback.RegisterAgent(a.AgentID, opts.Local.CustomTools)
	}
	if opts.Model != nil {
		ms := parseModelSelection(opts.Model)
		a.Model = &ms
	}
	return run, nil
}

// Reload reloads agent configuration from disk/settings.
func (a *Agent) Reload(ctx context.Context) error {
	_, err := a.client.agentUnary(ctx, "ReloadAgent", map[string]any{"agentId": a.AgentID})
	return err
}

// ListMessages returns persisted agent messages.
func (a *Agent) ListMessages(ctx context.Context) ([]AgentMessage, error) {
	resp, err := a.client.agentUnary(ctx, "ListAgentMessages", map[string]any{
		"agentId": a.AgentID,
		"options": optionsWithAPIKey(a.apiKey, nil),
	})
	if err != nil {
		return nil, err
	}
	out := make([]AgentMessage, 0)
	for _, item := range sliceMaps(resp["messages"]) {
		out = append(out, AgentMessage{
			Type: stringField(item, "type"), UUID: stringField(item, "uuid"),
			AgentID: stringField(item, "agentId"), Message: item["message"],
		})
	}
	return out, nil
}

// ListArtifacts lists files produced by the agent.
func (a *Agent) ListArtifacts(ctx context.Context) ([]SdkArtifact, error) {
	resp, err := a.client.agentUnary(ctx, "ListArtifacts", map[string]any{
		"agentId": a.AgentID,
		"options": optionsWithAPIKey(a.apiKey, nil),
	})
	if err != nil {
		return nil, err
	}
	out := make([]SdkArtifact, 0)
	for _, item := range sliceMaps(resp["artifacts"]) {
		out = append(out, SdkArtifact{
			Path: stringField(item, "path"), SizeBytes: intField(item, "sizeBytes"),
			UpdatedAt: stringField(item, "updatedAt"),
		})
	}
	return out, nil
}

// DownloadArtifact downloads an artifact by path.
func (a *Agent) DownloadArtifact(ctx context.Context, path string) ([]byte, error) {
	return a.client.downloadArtifact(ctx, a.AgentID, path, a.apiKey)
}

// Archive archives the agent.
func (a *Agent) Archive(ctx context.Context, opts AgentOperationOptions) error {
	_, err := a.client.agentUnary(ctx, "ArchiveAgent", map[string]any{
		"agentId": a.AgentID, "options": optionsWithAPIKey(a.apiKey, opts.ToWire()),
	})
	return err
}

// Unarchive restores an archived agent.
func (a *Agent) Unarchive(ctx context.Context, opts AgentOperationOptions) error {
	_, err := a.client.agentUnary(ctx, "UnarchiveAgent", map[string]any{
		"agentId": a.AgentID, "options": optionsWithAPIKey(a.apiKey, opts.ToWire()),
	})
	return err
}

// Delete permanently deletes the agent.
func (a *Agent) Delete(ctx context.Context, opts AgentOperationOptions) error {
	_, err := a.client.agentUnary(ctx, "DeleteAgent", map[string]any{
		"agentId": a.AgentID, "options": optionsWithAPIKey(a.apiKey, opts.ToWire()),
	})
	return err
}

// GetAgent fetches agent metadata by ID.
func GetAgent(ctx context.Context, agentID string, opts GetAgentOptions, clientOpts ...ClientOption) (SDKAgentInfo, error) {
	c, err := defaultOrLaunchClient(ctx, clientOpts...)
	if err != nil {
		return SDKAgentInfo{}, err
	}
	return c.getAgent(ctx, agentID, opts.CWD, apiKeyOrEnv(opts.APIKey, c.allowAPIKeyEnvFallback))
}

// ListAgents lists agents with optional filters.
func ListAgents(ctx context.Context, opts ListAgentsOptions, clientOpts ...ClientOption) (ListResult[SDKAgentInfo], error) {
	c, err := defaultOrLaunchClient(ctx, clientOpts...)
	if err != nil {
		return ListResult[SDKAgentInfo]{}, err
	}
	return c.listAgents(ctx, opts.ToWire())
}

// ListAgentsOptions configures agent listing.
type ListAgentsOptions struct {
	Runtime         string
	CWD             string
	IncludeArchived *bool
	PRURL           string
	Limit           int
	APIKey          string
	Cursor          string
}

func (o ListAgentsOptions) ToWire() map[string]any {
	m := map[string]any{}
	if o.Runtime != "" {
		m["runtime"] = o.Runtime
	}
	if o.CWD != "" {
		m["cwd"] = o.CWD
	}
	if o.IncludeArchived != nil {
		m["includeArchived"] = *o.IncludeArchived
	}
	if o.PRURL != "" {
		m["prUrl"] = o.PRURL
	}
	if o.Limit > 0 {
		m["limit"] = o.Limit
	}
	if o.APIKey != "" {
		m["apiKey"] = o.APIKey
	}
	if o.Cursor != "" {
		m["cursor"] = o.Cursor
	}
	return m
}

func optionsWithAPIKey(apiKey string, opts map[string]any) map[string]any {
	merged := copyMap(opts)
	if apiKey != "" && merged["apiKey"] == nil {
		merged["apiKey"] = apiKey
	}
	return merged
}

func randomUUID() string {
	var b [16]byte
	_, _ = rand.Read(b[:])
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}
