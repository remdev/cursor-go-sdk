package cursor

import (
	"context"
)

// Cursor provides account-level SDK operations.
var Cursor = cursorNamespace{}

type cursorNamespace struct{}

// Me returns information about the authenticated API key.
func (cursorNamespace) Me(ctx context.Context, opts CursorRequestOptions, clientOpts ...ClientOption) (SDKUser, error) {
	c, err := defaultOrLaunchClient(ctx, clientOpts...)
	if err != nil {
		return SDKUser{}, err
	}
	return c.me(ctx, opts.APIKey)
}

// Configure sets process-wide SDK defaults.
func (cursorNamespace) Configure(opts ConfigureOptions) {
	Configure(opts)
}

// Models exposes model listing.
var Models = modelsNamespace{}

type modelsNamespace struct{}

func (modelsNamespace) List(ctx context.Context, opts CursorRequestOptions, clientOpts ...ClientOption) ([]SDKModel, error) {
	c, err := defaultOrLaunchClient(ctx, clientOpts...)
	if err != nil {
		return nil, err
	}
	return c.listModels(ctx, opts.APIKey)
}

// Repositories exposes connected repository listing.
var Repositories = repositoriesNamespace{}

type repositoriesNamespace struct{}

func (repositoriesNamespace) List(ctx context.Context, opts CursorRequestOptions, clientOpts ...ClientOption) ([]SDKRepository, error) {
	c, err := defaultOrLaunchClient(ctx, clientOpts...)
	if err != nil {
		return nil, err
	}
	return c.listRepositories(ctx, opts.APIKey)
}

// AgentsResource groups agent operations on a client.
type AgentsResource struct {
	client *Client
}

func (r *AgentsResource) Create(ctx context.Context, opts AgentOptions) (*Agent, error) {
	return r.client.createAgent(ctx, opts, opts.IdempotencyKey)
}

func (r *AgentsResource) Resume(ctx context.Context, agentID string, opts AgentOptions) (*Agent, error) {
	return r.client.resumeAgent(ctx, agentID, opts)
}

func (r *AgentsResource) List(ctx context.Context, opts ListAgentsOptions) (ListResult[SDKAgentInfo], error) {
	return r.client.listAgents(ctx, opts.ToWire())
}

func (r *AgentsResource) Get(ctx context.Context, agentID string, opts GetAgentOptions) (SDKAgentInfo, error) {
	return r.client.getAgent(ctx, agentID, opts.CWD, apiKeyOrEnv(opts.APIKey, r.client.allowAPIKeyEnvFallback))
}

func (r *AgentsResource) ListRuns(ctx context.Context, agentID string, opts ListRunsOptions) (ListResult[*Run], error) {
	return r.client.listRuns(ctx, agentID, opts.ToWire())
}

func (r *AgentsResource) GetRun(ctx context.Context, runID string, opts GetRunOptions) (*Run, error) {
	return r.client.getRun(ctx, runID, opts.ToWire())
}

func (r *AgentsResource) CancelRun(ctx context.Context, runID, agentID string) error {
	return r.client.cancelRun(ctx, runID, agentID)
}

func (r *AgentsResource) Archive(ctx context.Context, agentID string, opts AgentOperationOptions) error {
	_, err := r.client.agentUnary(ctx, "ArchiveAgent", map[string]any{
		"agentId": agentID, "options": opts.ToWire(),
	})
	return err
}

func (r *AgentsResource) Unarchive(ctx context.Context, agentID string, opts AgentOperationOptions) error {
	_, err := r.client.agentUnary(ctx, "UnarchiveAgent", map[string]any{
		"agentId": agentID, "options": opts.ToWire(),
	})
	return err
}

func (r *AgentsResource) Delete(ctx context.Context, agentID string, opts AgentOperationOptions) error {
	_, err := r.client.agentUnary(ctx, "DeleteAgent", map[string]any{
		"agentId": agentID, "options": opts.ToWire(),
	})
	return err
}

// ListRunsOptions configures run listing.
type ListRunsOptions struct {
	Runtime string
	CWD     string
	Limit   int
	APIKey  string
	Cursor  string
}

func (o ListRunsOptions) ToWire() map[string]any {
	m := map[string]any{}
	if o.Runtime != "" {
		m["runtime"] = o.Runtime
	}
	if o.CWD != "" {
		m["cwd"] = o.CWD
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

// ModelsResource lists models through a client.
type ModelsResource struct {
	client *Client
}

func (r *ModelsResource) List(ctx context.Context, opts CursorRequestOptions) ([]SDKModel, error) {
	return r.client.listModels(ctx, opts.APIKey)
}

// RepositoriesResource lists repositories through a client.
type RepositoriesResource struct {
	client *Client
}

func (r *RepositoriesResource) List(ctx context.Context, opts CursorRequestOptions) ([]SDKRepository, error) {
	return r.client.listRepositories(ctx, opts.APIKey)
}

// CursorClient is an alias for Client matching the Python SDK name.
type CursorClient = Client
