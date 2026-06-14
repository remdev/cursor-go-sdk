package cursor

import (
	"context"
	"encoding/base64"
	"io"
	"os"
	"sync"
	"time"

	"github.com/remdev/cursor-go-sdk/internal/bridge"
	conn "github.com/remdev/cursor-go-sdk/internal/connect"
)

const (
	agentService         = "sdk.v1.SdkAgentService"
	bridgeControlService = "sdk.v1.SdkBridgeControlService"
	cursorService        = "sdk.v1.SdkCursorService"
)

// Client talks to the Cursor SDK bridge.
type Client struct {
	transport                *conn.Transport
	ownedBridge              *bridge.Bridge
	allowAPIKeyEnvFallback   bool
	toolCallback             *ToolCallbackServer
	mu                       sync.Mutex

	Agents       *AgentsResource
	Models       *ModelsResource
	Repositories *RepositoriesResource
}

// ClientOption configures a Client.
type ClientOption func(*clientConfig)

type clientConfig struct {
	baseURL                string
	authToken              string
	endpoint               *bridge.Endpoint
	transport              *conn.Transport
	ownedBridge            *bridge.Bridge
	timeout                time.Duration
	unaryTimeout           time.Duration
	streamTimeout          time.Duration
	maxRetries             int
	allowAPIKeyEnvFallback *bool
}

// Connect returns a client attached to an existing bridge.
func Connect(baseURL, authToken string, opts ...ClientOption) (*Client, error) {
	cfg := clientConfig{baseURL: baseURL, authToken: authToken}
	for _, opt := range opts {
		opt(&cfg)
	}
	return newClient(cfg)
}

// LaunchBridge starts a local bridge and returns a client.
func LaunchBridge(ctx context.Context, opts ...ClientOption) (*Client, error) {
	cfg := clientConfig{}
	for _, opt := range opts {
		opt(&cfg)
	}
	launchCfg := bridge.LaunchConfig{
		Timeout: 30 * time.Second,
	}
	if cfg.baseURL != "" {
		launchCfg.Workspace = cfg.baseURL
	}
	if launchCfg.Workspace == "" {
		if wd, err := os.Getwd(); err == nil {
			launchCfg.Workspace = wd
		}
	}

	toolSrv := NewToolCallbackServer()
	launchCfg.ToolCallbackURL = toolSrv.URL
	launchCfg.ToolCallbackToken = toolSrv.Token

	b, err := bridge.Launch(ctx, launchCfg)
	if err != nil {
		toolSrv.Close()
		return nil, wrapError(err)
	}
	cfg.ownedBridge = b
	cfg.endpoint = &b.Endpoint
	allow := true
	cfg.allowAPIKeyEnvFallback = &allow
	c, err := newClient(cfg)
	if err != nil {
		b.Close()
		toolSrv.Close()
		return nil, err
	}
	c.toolCallback = toolSrv
	return c, nil
}

// WithWorkspace sets the bridge workspace directory.
func WithWorkspace(path string) ClientOption {
	return func(c *clientConfig) { c.baseURL = path }
}

// WithAllowAPIKeyEnvFallback allows reading CURSOR_API_KEY from the environment.
func WithAllowAPIKeyEnvFallback(allow bool) ClientOption {
	return func(c *clientConfig) { c.allowAPIKeyEnvFallback = &allow }
}

// WithUnaryTimeout sets unary RPC timeout.
func WithUnaryTimeout(d time.Duration) ClientOption {
	return func(c *clientConfig) { c.unaryTimeout = d }
}

// WithStreamTimeout sets streaming RPC timeout.
func WithStreamTimeout(d time.Duration) ClientOption {
	return func(c *clientConfig) { c.streamTimeout = d }
}

// WithMaxRetries sets retry count for idempotent unary RPCs.
func WithMaxRetries(n int) ClientOption {
	return func(c *clientConfig) { c.maxRetries = n }
}

func newClient(cfg clientConfig) (*Client, error) {
	var transport *conn.Transport
	if cfg.transport != nil {
		transport = cfg.transport
	} else {
		baseURL := cfg.baseURL
		token := cfg.authToken
		if cfg.endpoint != nil {
			baseURL = cfg.endpoint.URL
			token = cfg.endpoint.AuthToken
		}
		if baseURL == "" || token == "" {
			return nil, configurationError(
				"Client requires a bridge endpoint or base URL with auth token",
				"missing_bridge_endpoint",
			)
		}
		tOpts := []conn.TransportOption{}
		if cfg.unaryTimeout > 0 {
			tOpts = append(tOpts, conn.WithUnaryTimeout(cfg.unaryTimeout))
		}
		if cfg.streamTimeout > 0 {
			tOpts = append(tOpts, conn.WithStreamTimeout(cfg.streamTimeout))
		}
		if cfg.maxRetries > 0 {
			tOpts = append(tOpts, conn.WithMaxRetries(cfg.maxRetries))
		}
		transport = conn.NewTransport(baseURL, token, tOpts...)
	}
	allowFallback := false
	if cfg.allowAPIKeyEnvFallback != nil {
		allowFallback = *cfg.allowAPIKeyEnvFallback
	}
	c := &Client{
		transport:              transport,
		ownedBridge:            cfg.ownedBridge,
		allowAPIKeyEnvFallback: allowFallback,
	}
	c.installResources()
	return c, nil
}

func (c *Client) installResources() {
	c.Agents = &AgentsResource{client: c}
	c.Models = &ModelsResource{client: c}
	c.Repositories = &RepositoriesResource{client: c}
}

// Close releases bridge and transport resources.
func (c *Client) Close() error {
	if c.toolCallback != nil {
		c.toolCallback.Close()
		c.toolCallback = nil
	}
	if c.transport != nil {
		c.transport.Close()
	}
	if c.ownedBridge != nil {
		return c.ownedBridge.Close()
	}
	return nil
}

// WithOptions returns a client sharing the transport with overridden timeouts.
func (c *Client) WithOptions(opts ...ClientOption) (*Client, error) {
	cfg := clientConfig{
		transport:              c.transport,
		allowAPIKeyEnvFallback: &c.allowAPIKeyEnvFallback,
	}
	for _, opt := range opts {
		opt(&cfg)
	}
	clone, err := newClient(cfg)
	if err != nil {
		return nil, err
	}
	clone.toolCallback = c.toolCallback
	return clone, nil
}

// Ping checks bridge connectivity.
func (c *Client) Ping(ctx context.Context) (string, error) {
	resp, err := c.controlUnary(ctx, "Ping", map[string]any{})
	if err != nil {
		return "", err
	}
	return stringField(resp, "message"), nil
}

// GetVersion returns bridge version metadata.
func (c *Client) GetVersion(ctx context.Context) (map[string]any, error) {
	return c.controlUnary(ctx, "GetVersion", map[string]any{})
}

// Shutdown requests graceful bridge shutdown.
func (c *Client) Shutdown(ctx context.Context, graceSeconds int) error {
	_, err := c.controlUnary(ctx, "Shutdown", map[string]any{"graceSeconds": graceSeconds})
	return err
}

func (c *Client) createAgent(ctx context.Context, opts AgentOptions, idempotencyKey string) (*Agent, error) {
	applyConfigureDefaults(&opts)
	wire := opts.ToWire()
	key := apiKeyOrEnv(opts.APIKey, c.allowAPIKeyEnvFallback)
	if key == "" {
		return nil, configurationError(
			"Agent.create requires api_key (set CURSOR_API_KEY or pass APIKey)",
			"missing_api_key",
		)
	}
	wire["apiKey"] = key
	if wire["model"] == nil || len(wire["model"].(map[string]any)) == 0 {
		return nil, configurationError(
			`Agent.create requires model (e.g. model="composer-2.5")`,
			"missing_model",
		)
	}
	if c.toolCallback != nil {
		if local, ok := wire["local"].(map[string]any); ok {
			if tools, ok := local["customTools"].(map[string]any); ok && len(tools) > 0 {
				agentID := stringField(wire, "agentId")
				if agentID == "" {
					agentID = "pending"
				}
				c.toolCallback.RegisterAgent(agentID, opts.Local.CustomTools)
			}
		}
	}
	req := map[string]any{"options": stripEmptyMessage(wire)}
	if opts.Cloud != nil {
		if idempotencyKey == "" {
			idempotencyKey = newID()
		}
		req["idempotencyKey"] = idempotencyKey
	} else if idempotencyKey != "" {
		return nil, configurationError(
			"Agent.create idempotency_key is only supported for cloud agents",
			"local_idempotency_key_unsupported",
		)
	}
	resp, err := c.agentUnary(ctx, "CreateAgent", req)
	if err != nil {
		return nil, err
	}
	agentID := stringField(resp, "agentId")
	var model *ModelSelection
	if m, ok := resp["model"].(map[string]any); ok {
		ms := parseModelSelection(m)
		model = &ms
	}
	return &Agent{
		client: c, AgentID: agentID, Model: model, apiKey: key,
	}, nil
}

func (c *Client) resumeAgent(ctx context.Context, agentID string, opts AgentOptions) (*Agent, error) {
	wire := opts.ToWire()
	resp, err := c.agentUnary(ctx, "ResumeAgent", map[string]any{
		"agentId": agentID, "options": stripEmptyMessage(wire),
	})
	if err != nil {
		return nil, err
	}
	var model *ModelSelection
	if m, ok := resp["model"].(map[string]any); ok {
		ms := parseModelSelection(m)
		model = &ms
	}
	return &Agent{
		client: c,
		AgentID: stringField(resp, "agentId"),
		Model:   model,
		apiKey:  stringField(wire, "apiKey"),
	}, nil
}

func (c *Client) getAgent(ctx context.Context, agentID, cwd, apiKey string) (SDKAgentInfo, error) {
	resp, err := c.agentUnary(ctx, "GetAgent", map[string]any{
		"agentId": agentID,
		"options": map[string]any{"cwd": cwd, "apiKey": apiKey},
	})
	if err != nil {
		return SDKAgentInfo{}, err
	}
	if agent, ok := resp["agent"].(map[string]any); ok {
		return parseSDKAgentInfo(agent), nil
	}
	return SDKAgentInfo{}, nil
}

func (c *Client) listAgents(ctx context.Context, options map[string]any) (ListResult[SDKAgentInfo], error) {
	resp, err := c.agentUnary(ctx, "ListAgents", map[string]any{"options": options})
	if err != nil {
		return ListResult[SDKAgentInfo]{}, err
	}
	items := make([]SDKAgentInfo, 0)
	for _, item := range sliceMaps(resp["items"]) {
		items = append(items, parseSDKAgentInfo(item))
	}
	base := copyMap(options)
	return ListResult[SDKAgentInfo]{
		Items:      items,
		NextCursor: stringField(resp, "nextCursor"),
		fetchNext: func(cursor string) (ListResult[SDKAgentInfo], error) {
			nextOpts := copyMap(base)
			nextOpts["cursor"] = cursor
			return c.listAgents(ctx, nextOpts)
		},
	}, nil
}

func (c *Client) getRun(ctx context.Context, runID string, options map[string]any) (*Run, error) {
	resp, err := c.agentUnary(ctx, "GetRun", map[string]any{
		"runId": runID, "options": options,
	})
	if err != nil {
		return nil, err
	}
	snap := parseRunResult(mapFrom(resp["run"]))
	return newRunFromSnapshot(c, snap), nil
}

func (c *Client) listRuns(ctx context.Context, agentID string, options map[string]any) (ListResult[*Run], error) {
	resp, err := c.agentUnary(ctx, "ListRuns", map[string]any{
		"agentId": agentID, "options": options,
	})
	if err != nil {
		return ListResult[*Run]{}, err
	}
	items := make([]*Run, 0)
	for _, item := range sliceMaps(resp["items"]) {
		items = append(items, newRunFromSnapshot(c, parseRunResult(item)))
	}
	base := copyMap(options)
	return ListResult[*Run]{
		Items:      items,
		NextCursor: stringField(resp, "nextCursor"),
		fetchNext: func(cursor string) (ListResult[*Run], error) {
			nextOpts := copyMap(base)
			nextOpts["cursor"] = cursor
			return c.listRuns(ctx, agentID, nextOpts)
		},
	}, nil
}

func (c *Client) cancelRun(ctx context.Context, runID, agentID string) error {
	_, err := c.agentUnary(ctx, "CancelRun", map[string]any{
		"runId": runID, "agentId": agentID,
	})
	return err
}

func (c *Client) waitLiveRun(ctx context.Context, runID string) (RunResult, error) {
	resp, err := c.agentUnary(ctx, "WaitLiveRun", map[string]any{"runId": runID})
	if err != nil {
		return RunResult{}, err
	}
	return parseRunResult(mapFrom(resp["result"])), nil
}

func (c *Client) observeRun(ctx context.Context, runID, afterOffset string) (*conn.StreamReader, error) {
	return c.agentStream(ctx, "ObserveRun", map[string]any{
		"runId": runID, "afterOffset": afterOffset,
	})
}

func (c *Client) me(ctx context.Context, apiKey string) (SDKUser, error) {
	resp, err := c.cursorUnary(ctx, "Me", map[string]any{
		"options": map[string]any{
			"apiKey": apiKeyOrEnv(apiKey, c.allowAPIKeyEnvFallback),
		},
	})
	if err != nil {
		return SDKUser{}, err
	}
	return parseSDKUser(mapFrom(resp["user"])), nil
}

func (c *Client) listModels(ctx context.Context, apiKey string) ([]SDKModel, error) {
	resp, err := c.cursorUnary(ctx, "ListModels", map[string]any{
		"options": map[string]any{
			"apiKey": apiKeyOrEnv(apiKey, c.allowAPIKeyEnvFallback),
		},
	})
	if err != nil {
		return nil, err
	}
	out := make([]SDKModel, 0)
	for _, item := range sliceMaps(resp["items"]) {
		out = append(out, parseSDKModel(item))
	}
	return out, nil
}

func (c *Client) listRepositories(ctx context.Context, apiKey string) ([]SDKRepository, error) {
	resp, err := c.cursorUnary(ctx, "ListRepositories", map[string]any{
		"options": map[string]any{
			"apiKey": apiKeyOrEnv(apiKey, c.allowAPIKeyEnvFallback),
		},
	})
	if err != nil {
		return nil, err
	}
	out := make([]SDKRepository, 0)
	for _, item := range sliceMaps(resp["items"]) {
		out = append(out, SDKRepository{URL: stringField(item, "url")})
	}
	return out, nil
}

func (c *Client) downloadArtifact(ctx context.Context, agentID, path, apiKey string) ([]byte, error) {
	stream, err := c.agentStream(ctx, "DownloadArtifact", map[string]any{
		"agentId": agentID, "path": path,
		"options": map[string]any{"apiKey": apiKey},
	})
	if err != nil {
		return nil, err
	}
	defer stream.Close()
	var parts [][]byte
	for {
		chunk, err := stream.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, wrapError(err)
		}
		data := stringField(chunk, "data")
		if data != "" {
			b, err := base64.StdEncoding.DecodeString(data)
			if err != nil {
				return nil, err
			}
			parts = append(parts, b)
		}
	}
	return concatBytes(parts), nil
}

func (c *Client) agentUnary(ctx context.Context, method string, msg map[string]any) (map[string]any, error) {
	if err := c.requireCloudAPIKey(method, msg); err != nil {
		return nil, err
	}
	resp, err := c.transport.Unary(ctx, agentService, method, msg)
	return resp, wrapError(err)
}

func (c *Client) agentStream(ctx context.Context, method string, msg map[string]any) (*conn.StreamReader, error) {
	if err := c.requireCloudAPIKey(method, msg); err != nil {
		return nil, err
	}
	stream, err := c.transport.Stream(ctx, agentService, method, msg)
	return stream, wrapError(err)
}

func (c *Client) controlUnary(ctx context.Context, method string, msg map[string]any) (map[string]any, error) {
	resp, err := c.transport.Unary(ctx, bridgeControlService, method, msg)
	return resp, wrapError(err)
}

func (c *Client) cursorUnary(ctx context.Context, method string, msg map[string]any) (map[string]any, error) {
	resp, err := c.transport.Unary(ctx, cursorService, method, msg)
	return resp, wrapError(err)
}

func (c *Client) requireCloudAPIKey(method string, msg map[string]any) error {
	if c.allowAPIKeyEnvFallback {
		return nil
	}
	if !isCloudAgentRPC(method, msg) {
		return nil
	}
	opts, _ := msg["options"].(map[string]any)
	if opts != nil && stringField(opts, "apiKey") != "" {
		return nil
	}
	return configurationError(
		"Cloud agent SDK calls through a caller-supplied bridge require an explicit api_key",
		"missing_api_key",
	)
}

func isCloudAgentRPC(method string, msg map[string]any) bool {
	opts, _ := msg["options"].(map[string]any)
	agentID := stringField(msg, "agentId")
	switch method {
	case "CreateAgent":
		return opts != nil && hasCloudOptions(opts)
	case "ResumeAgent":
		return stringsHasPrefix(agentID, "bc-") || (opts != nil && hasCloudOptions(opts))
	case "GetAgent", "ArchiveAgent", "UnarchiveAgent", "DeleteAgent",
		"ListAgentMessages", "ListArtifacts", "DownloadArtifact":
		return stringsHasPrefix(agentID, "bc-")
	case "ListAgents", "ListRuns", "GetRun":
		if opts == nil {
			return false
		}
		rt := stringField(opts, "runtime")
		return rt == "cloud" || rt == "RUNTIME_CLOUD"
	case "Send":
		return stringsHasPrefix(agentID, "bc-")
	case "CancelRun", "WaitLiveRun", "ObserveRun", "GetRunConversation":
		return true
	}
	return false
}

func stringsHasPrefix(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}

func sliceMaps(v any) []map[string]any {
	raw, ok := v.([]any)
	if !ok {
		return nil
	}
	out := make([]map[string]any, 0, len(raw))
	for _, item := range raw {
		if m, ok := item.(map[string]any); ok {
			out = append(out, m)
		}
	}
	return out
}

func mapFrom(v any) map[string]any {
	m, _ := v.(map[string]any)
	if m == nil {
		return map[string]any{}
	}
	return m
}

func copyMap(m map[string]any) map[string]any {
	out := make(map[string]any, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}

func concatBytes(parts [][]byte) []byte {
	n := 0
	for _, p := range parts {
		n += len(p)
	}
	out := make([]byte, 0, n)
	for _, p := range parts {
		out = append(out, p...)
	}
	return out
}

func newID() string {
	return randomUUID()
}
