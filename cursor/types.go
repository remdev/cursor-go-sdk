package cursor

import (
	"encoding/base64"
	"mime"
	"os"
	"path/filepath"
	"strings"
)

// RunStatus is the lifecycle state of a run.
type RunStatus string

const (
	RunStatusRunning   RunStatus = "running"
	RunStatusFinished  RunStatus = "finished"
	RunStatusError     RunStatus = "error"
	RunStatusCancelled RunStatus = "cancelled"
	RunStatusExpired   RunStatus = "expired"
)

// RunResultStatus is the terminal status returned by Wait.
type RunResultStatus = RunStatus

// AgentMode selects agent or plan mode.
type AgentMode string

const (
	AgentModeAgent AgentMode = "agent"
	AgentModePlan  AgentMode = "plan"
)

// SettingSource controls which ambient settings are loaded for local agents.
type SettingSource string

const (
	SettingSourceProject SettingSource = "project"
	SettingSourceUser    SettingSource = "user"
	SettingSourceTeam    SettingSource = "team"
	SettingSourceMDM     SettingSource = "mdm"
	SettingSourcePlugins SettingSource = "plugins"
	SettingSourceAll     SettingSource = "all"
)

// ModelParameterValue is a per-model option such as reasoning effort.
type ModelParameterValue struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}

func (p ModelParameterValue) ToWire() map[string]any {
	return map[string]any{"id": p.ID, "value": p.Value}
}

// ModelSelection identifies a model and optional parameters.
type ModelSelection struct {
	ID     string                `json:"id"`
	Params []ModelParameterValue `json:"params,omitempty"`
}

func ModelFromString(id string) ModelSelection {
	return ModelSelection{ID: id}
}

func (m ModelSelection) ToWire() map[string]any {
	params := make([]map[string]any, 0, len(m.Params))
	for _, p := range m.Params {
		params = append(params, p.ToWire())
	}
	return dropEmpty(map[string]any{
		"id":     m.ID,
		"params": params,
	})
}

func parseModelSelection(v any) ModelSelection {
	switch t := v.(type) {
	case string:
		return ModelFromString(t)
	case ModelSelection:
		return t
	case map[string]any:
		ms := ModelSelection{ID: stringField(t, "id")}
		if raw, ok := t["params"].([]any); ok {
			for _, item := range raw {
				if m, ok := item.(map[string]any); ok {
					ms.Params = append(ms.Params, ModelParameterValue{
						ID: stringField(m, "id"), Value: stringField(m, "value"),
					})
				}
			}
		}
		return ms
	default:
		return ModelSelection{}
	}
}

// SandboxOptions configures local sandbox behavior.
type SandboxOptions struct {
	Enabled *bool `json:"enabled,omitempty"`
}

func (s SandboxOptions) ToWire() map[string]any {
	return dropEmpty(map[string]any{"enabled": s.Enabled})
}

// LocalAgentStoreConfig configures a custom local store.
type LocalAgentStoreConfig struct {
	Type    string `json:"type"`
	RootDir string `json:"rootDir,omitempty"`
}

func (c LocalAgentStoreConfig) ToWire() map[string]any {
	return dropEmpty(map[string]any{"type": c.Type, "rootDir": c.RootDir})
}

// CustomTool defines a host-executed tool for local agents.
type CustomTool struct {
	Description  string
	InputSchema  map[string]any
	Execute      func(args map[string]any, ctx CustomToolContext) (any, error)
}

// CustomToolContext carries per-invocation metadata.
type CustomToolContext struct {
	ToolCallID string
}

// LocalAgentOptions configures a local agent runtime.
type LocalAgentOptions struct {
	CWD            []string
	SettingSources []SettingSource
	SandboxOptions *SandboxOptions
	Store          *LocalAgentStoreConfig
	AutoReview     *bool
	CustomTools    map[string]CustomTool
}

func (o LocalAgentOptions) ToWire() map[string]any {
	cwd := o.CWD
	if len(cwd) == 0 {
		cwd = []string{}
	}
	sources := make([]string, 0, len(o.SettingSources))
	for _, s := range o.SettingSources {
		sources = append(sources, enumValue(string(s), "SETTING_SOURCE_"))
	}
	payload := map[string]any{
		"cwd":            cwd,
		"settingSources": sources,
		"autoReview":     o.AutoReview,
	}
	if o.SandboxOptions != nil {
		payload["sandboxOptions"] = o.SandboxOptions.ToWire()
	}
	if o.Store != nil {
		payload["store"] = o.Store.ToWire()
	}
	if len(o.CustomTools) > 0 {
		tools := make(map[string]any, len(o.CustomTools))
		for name, tool := range o.CustomTools {
			tools[name] = dropEmpty(map[string]any{
				"description":  tool.Description,
				"inputSchema": tool.InputSchema,
			})
		}
		payload["customTools"] = tools
	}
	return dropEmpty(payload)
}

// CloudEnvironment selects cloud vs self-hosted pool.
type CloudEnvironment struct {
	Type string `json:"type"`
	Name string `json:"name,omitempty"`
}

func (e CloudEnvironment) ToWire() map[string]any {
	return dropEmpty(map[string]any{
		"type": enumValue(e.Type, "CLOUD_ENVIRONMENT_TYPE_"),
		"name": e.Name,
	})
}

// CloudRepository identifies a repo for cloud agents.
type CloudRepository struct {
	URL         string `json:"url"`
	StartingRef string `json:"startingRef,omitempty"`
	PRURL       string `json:"prUrl,omitempty"`
}

func (r CloudRepository) ToWire() map[string]any {
	return dropEmpty(map[string]any{
		"url": r.URL, "startingRef": r.StartingRef, "prUrl": r.PRURL,
	})
}

// CloudAgentOptions configures a cloud agent runtime.
type CloudAgentOptions struct {
	Env                  *CloudEnvironment
	Repos                []CloudRepository
	WorkOnCurrentBranch  *bool
	AutoCreatePR         *bool
	SkipReviewerRequest  *bool
	EnvVars              map[string]string
}

func (o CloudAgentOptions) ToWire() map[string]any {
	repos := make([]map[string]any, 0, len(o.Repos))
	for _, r := range o.Repos {
		repos = append(repos, r.ToWire())
	}
	envVars := map[string]string{}
	for k, v := range o.EnvVars {
		envVars[k] = v
	}
	payload := map[string]any{
		"repos":                repos,
		"envVars":              envVars,
		"workOnCurrentBranch":  o.WorkOnCurrentBranch,
		"autoCreatePr":         o.AutoCreatePR,
		"skipReviewerRequest":  o.SkipReviewerRequest,
	}
	if o.Env != nil {
		payload["env"] = o.Env.ToWire()
	}
	return dropEmpty(payload, "envVars")
}

// McpAuth holds OAuth2 client credentials for HTTP MCP servers.
type McpAuth struct {
	ClientID     string   `json:"clientId"`
	ClientSecret string   `json:"clientSecret,omitempty"`
	Scopes       []string `json:"scopes,omitempty"`
}

func (a McpAuth) ToWire() map[string]any {
	return dropEmpty(map[string]any{
		"clientId": a.ClientID, "clientSecret": a.ClientSecret, "scopes": a.Scopes,
	})
}

// HttpMcpServerConfig configures an HTTP or SSE MCP server.
type HttpMcpServerConfig struct {
	URL     string
	Type    string // "http" or "sse"
	Headers map[string]string
	Auth    *McpAuth
}

func (c HttpMcpServerConfig) ToWire() map[string]any {
	t := c.Type
	if t == "" {
		t = "http"
	}
	return map[string]any{
		"http": dropEmpty(map[string]any{
			"type":    enumValue(t, "HTTP_MCP_TRANSPORT_TYPE_"),
			"url":     c.URL,
			"headers": c.Headers,
			"auth":    authWire(c.Auth),
		}, "headers"),
	}
}

// StdioMcpServerConfig configures a stdio MCP server.
type StdioMcpServerConfig struct {
	Command string
	Args    []string
	Env     map[string]string
	CWD     string
}

func (c StdioMcpServerConfig) ToWire() map[string]any {
	return map[string]any{
		"stdio": dropEmpty(map[string]any{
			"command": c.Command,
			"args":    c.Args,
			"env":     c.Env,
			"cwd":     c.CWD,
		}, "env"),
	}
}

func authWire(a *McpAuth) map[string]any {
	if a == nil {
		return map[string]any{}
	}
	return a.ToWire()
}

// McpServer is implemented by MCP server config types.
type McpServer interface {
	ToWire() map[string]any
}

// AgentDefinition describes a sub-agent available to the main agent.
type AgentDefinition struct {
	Description string
	Prompt      string
	Model       any // string, ModelSelection, or "inherit"
	McpServers  []any // string names or inline configs
}

func (d AgentDefinition) ToWire() map[string]any {
	inherit := d.Model == nil || d.Model == "inherit"
	var modelWire map[string]any
	if !inherit {
		modelWire = parseModelSelection(d.Model).ToWire()
	}
	mcpServers := make([]map[string]any, 0, len(d.McpServers))
	for _, srv := range d.McpServers {
		switch t := srv.(type) {
		case string:
			mcpServers = append(mcpServers, map[string]any{"name": t})
		case McpServer:
			mcpServers = append(mcpServers, map[string]any{"inlineConfig": t.ToWire()})
		case map[string]any:
			if _, ok := t["name"]; ok {
				mcpServers = append(mcpServers, t)
			} else {
				mcpServers = append(mcpServers, map[string]any{"inlineConfig": McpServerFromMap(t)})
			}
		}
	}
	return dropEmpty(map[string]any{
		"description":  d.Description,
		"prompt":       d.Prompt,
		"model":        modelWire,
		"inheritModel": inherit,
		"mcpServers":   mcpServers,
	})
}

// AgentOptions configures agent creation or resume.
type AgentOptions struct {
	Model          any // string or ModelSelection
	APIKey         string
	Name           string
	Local          *LocalAgentOptions
	Cloud          *CloudAgentOptions
	McpServers     map[string]McpServer
	Agents         map[string]AgentDefinition
	AgentID        string
	IdempotencyKey string
	Mode           AgentMode
}

func (o AgentOptions) ToWire() map[string]any {
	mcp := make(map[string]any, len(o.McpServers))
	for name, srv := range o.McpServers {
		mcp[name] = srv.ToWire()
	}
	agents := make(map[string]any, len(o.Agents))
	for name, def := range o.Agents {
		agents[name] = def.ToWire()
	}
	var localWire, cloudWire map[string]any
	if o.Local != nil {
		localWire = o.Local.ToWire()
	}
	if o.Cloud != nil {
		cloudWire = o.Cloud.ToWire()
	}
	var modelWire map[string]any
	if o.Model != nil {
		modelWire = parseModelSelection(o.Model).ToWire()
	}
	return dropEmpty(map[string]any{
		"model":          modelWire,
		"apiKey":         o.APIKey,
		"name":           o.Name,
		"local":          localWire,
		"cloud":          cloudWire,
		"mcpServers":     mcp,
		"agents":         agents,
		"agentId":        o.AgentID,
		"mode":           enumValue(string(o.Mode), "AGENT_MODE_OPTION_"),
	})
}

// LocalSendOptions configures local-only send overrides.
type LocalSendOptions struct {
	Force       *bool
	CustomTools map[string]CustomTool
}

func (o LocalSendOptions) ToWire() map[string]any {
	payload := map[string]any{"force": o.Force}
	if len(o.CustomTools) > 0 {
		tools := make(map[string]any, len(o.CustomTools))
		for name, tool := range o.CustomTools {
			tools[name] = dropEmpty(map[string]any{
				"description": tool.Description, "inputSchema": tool.InputSchema,
			})
		}
		payload["customTools"] = tools
	}
	return dropEmpty(payload)
}

// SendOptions configures a single prompt submission.
type SendOptions struct {
	Model          any
	McpServers     map[string]McpServer
	Local          *LocalSendOptions
	IdempotencyKey string
	Mode           AgentMode
	OnDelta        DeltaCallback
	OnStep         StepCallback
}

func (o SendOptions) ToWire() map[string]any {
	mcp := make(map[string]any, len(o.McpServers))
	for name, srv := range o.McpServers {
		mcp[name] = srv.ToWire()
	}
	var modelWire map[string]any
	if o.Model != nil {
		modelWire = parseModelSelection(o.Model).ToWire()
	}
	var localWire map[string]any
	if o.Local != nil {
		localWire = o.Local.ToWire()
	}
	payload := map[string]any{
		"model": modelWire, "mcpServers": mcp, "local": localWire,
	}
	if o.OnDelta != nil {
		payload["enableDeltas"] = true
	}
	if o.OnStep != nil {
		payload["enableSteps"] = true
	}
	if o.Mode != "" {
		payload["mode"] = enumValue(string(o.Mode), "AGENT_MODE_OPTION_")
	}
	return dropEmpty(payload)
}

// UserMessage is a prompt with optional images.
type UserMessage struct {
	Text   string
	Images []SDKImage
}

func UserMessageFromText(text string) UserMessage {
	return UserMessage{Text: text}
}

func (m UserMessage) ToWire() map[string]any {
	images := make([]map[string]any, 0, len(m.Images))
	for _, img := range m.Images {
		images = append(images, img.ToWire())
	}
	return dropEmpty(map[string]any{"text": m.Text, "images": images})
}

// SDKImageDimension describes image size metadata.
type SDKImageDimension struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

// SDKImage is an image attachment for user messages.
type SDKImage struct {
	URL       string
	Data      string
	MimeType  string
	Dimension *SDKImageDimension
}

func ImageFromURL(url string, dim *SDKImageDimension) SDKImage {
	return SDKImage{URL: url, Dimension: dim}
}

func ImageFromFile(path string, mimeType string, dim *SDKImageDimension) (SDKImage, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return SDKImage{}, err
	}
	if mimeType == "" {
		mimeType = mime.TypeByExtension(filepath.Ext(path))
	}
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}
	return SDKImage{
		Data: base64.StdEncoding.EncodeToString(data), MimeType: mimeType, Dimension: dim,
	}, nil
}

func (i SDKImage) ToWire() map[string]any {
	out := map[string]any{}
	if i.URL != "" {
		out["url"] = map[string]any{"url": i.URL}
	} else if i.Data != "" {
		out["data"] = map[string]any{"data": i.Data, "mimeType": i.MimeType}
	}
	if i.Dimension != nil {
		out["dimension"] = map[string]any{"width": i.Dimension.Width, "height": i.Dimension.Height}
	}
	return dropEmpty(out)
}

// RunGitBranchInfo describes a branch touched by a run.
type RunGitBranchInfo struct {
	RepoURL string `json:"repoUrl"`
	Branch  string `json:"branch"`
	PRURL   string `json:"prUrl"`
}

// RunGitInfo aggregates git metadata from a run.
type RunGitInfo struct {
	Branches []RunGitBranchInfo `json:"branches"`
}

// RunResult is the terminal outcome of a run.
type RunResult struct {
	ID         string          `json:"runId"`
	RequestID  string          `json:"requestId,omitempty"`
	AgentID    string          `json:"agentId"`
	Status     RunStatus       `json:"status"`
	Result     string          `json:"result"`
	Model      *ModelSelection `json:"model,omitempty"`
	DurationMS int             `json:"durationMs"`
	Git        *RunGitInfo     `json:"git,omitempty"`
	CreatedAt  string          `json:"createdAt,omitempty"`
}

func parseRunResult(m map[string]any) RunResult {
	r := RunResult{
		ID:         stringField(m, "runId"),
		RequestID:  stringField(m, "requestId"),
		AgentID:    stringField(m, "agentId"),
		Status:     normalizeRunStatus(stringField(m, "status")),
		Result:     stringField(m, "result"),
		DurationMS: intField(m, "durationMs"),
		CreatedAt:  stringField(m, "createdAt"),
	}
	if model, ok := m["model"].(map[string]any); ok {
		ms := parseModelSelection(model)
		r.Model = &ms
	}
	if git, ok := m["git"].(map[string]any); ok {
		r.Git = parseRunGitInfo(git)
	}
	return r
}

func parseRunGitInfo(m map[string]any) *RunGitInfo {
	raw, ok := m["branches"].([]any)
	if !ok {
		return &RunGitInfo{}
	}
	branches := make([]RunGitBranchInfo, 0, len(raw))
	for _, item := range raw {
		if b, ok := item.(map[string]any); ok {
			branches = append(branches, RunGitBranchInfo{
				RepoURL: stringField(b, "repoUrl"),
				Branch:  stringField(b, "branch"),
				PRURL:   stringField(b, "prUrl"),
			})
		}
	}
	return &RunGitInfo{Branches: branches}
}

// RunSnapshot is a point-in-time view of a run.
type RunSnapshot = RunResult

// SDKAgentInfo describes a known agent.
type SDKAgentInfo struct {
	AgentID      string            `json:"agentId"`
	Name         string            `json:"name"`
	Summary      string            `json:"summary"`
	LastModified string            `json:"lastModified,omitempty"`
	Status       string            `json:"status,omitempty"`
	CreatedAt    string            `json:"createdAt,omitempty"`
	Archived     bool              `json:"archived"`
	Runtime      string            `json:"runtime,omitempty"`
	CWD          string            `json:"cwd,omitempty"`
	Env          *CloudEnvironment `json:"env,omitempty"`
	Repos        []string          `json:"repos,omitempty"`
}

func parseSDKAgentInfo(m map[string]any) SDKAgentInfo {
	info := SDKAgentInfo{
		AgentID:      stringField(m, "agentId"),
		Name:         stringField(m, "name"),
		Summary:      stringField(m, "summary"),
		LastModified: stringField(m, "lastModified"),
		Status:       normalizeAgentStatus(stringField(m, "status")),
		CreatedAt:    stringField(m, "createdAt"),
		Archived:     boolField(m, "archived"),
	}
	if local, ok := m["local"].(map[string]any); ok {
		info.Runtime = "local"
		info.CWD = stringField(local, "cwd")
	} else if cloud, ok := m["cloud"].(map[string]any); ok {
		info.Runtime = "cloud"
		if env, ok := cloud["env"].(map[string]any); ok {
			e := CloudEnvironment{
				Type: cloudEnvType(stringField(env, "type")),
				Name: stringField(env, "name"),
			}
			info.Env = &e
		}
		if repos, ok := cloud["repos"].([]any); ok {
			for _, r := range repos {
				info.Repos = append(info.Repos, fmtString(r))
			}
		}
	}
	return info
}

func cloudEnvType(raw string) string {
	return strings.ToLower(strings.TrimPrefix(raw, "CLOUD_ENVIRONMENT_TYPE_"))
}

// ModelParameterDefinition describes a configurable model parameter.
type ModelParameterDefinition struct {
	ID          string
	DisplayName string
	Values      []ModelParameterDefinitionValue
}

// ModelParameterDefinitionValue is one allowed value for a model parameter.
type ModelParameterDefinitionValue struct {
	Value       string
	DisplayName string
}

// ModelVariant is a preset parameter combination for a model.
type ModelVariant struct {
	Params      []ModelParameterValue
	DisplayName string
	Description string
	IsDefault   bool
}

// SDKModel describes an available model.
type SDKModel struct {
	ID          string                     `json:"id"`
	DisplayName string                     `json:"displayName"`
	Description string                     `json:"description"`
	Parameters  []ModelParameterDefinition `json:"parameters,omitempty"`
	Variants    []ModelVariant             `json:"variants,omitempty"`
}

func parseSDKModel(m map[string]any) SDKModel {
	model := SDKModel{
		ID: stringField(m, "id"), DisplayName: stringField(m, "displayName"),
		Description: stringField(m, "description"),
	}
	if raw, ok := m["parameters"].([]any); ok {
		for _, item := range raw {
			if p, ok := item.(map[string]any); ok {
				def := ModelParameterDefinition{
					ID: stringField(p, "id"), DisplayName: stringField(p, "displayName"),
				}
				if vals, ok := p["values"].([]any); ok {
					for _, v := range vals {
						if vm, ok := v.(map[string]any); ok {
							def.Values = append(def.Values, ModelParameterDefinitionValue{
								Value: stringField(vm, "value"), DisplayName: stringField(vm, "displayName"),
							})
						}
					}
				}
				model.Parameters = append(model.Parameters, def)
			}
		}
	}
	if raw, ok := m["variants"].([]any); ok {
		for _, item := range raw {
			if v, ok := item.(map[string]any); ok {
				variant := ModelVariant{
					DisplayName: stringField(v, "displayName"),
					Description: stringField(v, "description"),
					IsDefault:   boolField(v, "isDefault"),
				}
				if params, ok := v["params"].([]any); ok {
					for _, p := range params {
						if pm, ok := p.(map[string]any); ok {
							variant.Params = append(variant.Params, ModelParameterValue{
								ID: stringField(pm, "id"), Value: stringField(pm, "value"),
							})
						}
					}
				}
				model.Variants = append(model.Variants, variant)
			}
		}
	}
	return model
}

// SDKRepository is a connected repository.
type SDKRepository struct {
	URL string `json:"url"`
}

// SDKUser is the authenticated API key owner.
type SDKUser struct {
	APIKeyName    string `json:"apiKeyName"`
	CreatedAt     string `json:"createdAt"`
	UserID        *int   `json:"userId,omitempty"`
	UserEmail     string `json:"userEmail"`
	UserFirstName string `json:"userFirstName"`
	UserLastName  string `json:"userLastName"`
}

func parseSDKUser(m map[string]any) SDKUser {
	u := SDKUser{
		APIKeyName: stringField(m, "apiKeyName"), CreatedAt: stringField(m, "createdAt"),
		UserEmail: stringField(m, "userEmail"), UserFirstName: stringField(m, "userFirstName"),
		UserLastName: stringField(m, "userLastName"),
	}
	if _, ok := m["userId"]; ok {
		id := intField(m, "userId")
		u.UserID = &id
	}
	return u
}

// SdkArtifact is a file produced by an agent run.
type SdkArtifact struct {
	Path      string `json:"path"`
	SizeBytes int    `json:"sizeBytes"`
	UpdatedAt string `json:"updatedAt"`
}

// AgentMessage is a persisted agent message.
type AgentMessage struct {
	Type    string `json:"type"`
	UUID    string `json:"uuid"`
	AgentID string `json:"agentId"`
	Message any    `json:"message"`
}

// ListResult is a paginated list response.
type ListResult[T any] struct {
	Items      []T
	NextCursor string
	fetchNext  func(cursor string) (ListResult[T], error)
}

func (l ListResult[T]) HasNextPage() bool { return l.NextCursor != "" }

// NextPageInfo returns pagination metadata for the next page.
func (l ListResult[T]) NextPageInfo() map[string]string {
	if l.NextCursor == "" {
		return map[string]string{}
	}
	return map[string]string{"cursor": l.NextCursor}
}

func (l ListResult[T]) NextPage() (ListResult[T], error) {
	if l.NextCursor == "" {
		return ListResult[T]{}, nil
	}
	if l.fetchNext == nil {
		return ListResult[T]{}, configurationError(
			"This ListResult was not created by a pageable SDK method.",
			"pagination_unavailable",
		)
	}
	return l.fetchNext(l.NextCursor)
}

// AllPages iterates all items across pages.
func (l ListResult[T]) AllPages() ([]T, error) {
	var all []T
	page := l
	for {
		all = append(all, page.Items...)
		if !page.HasNextPage() {
			return all, nil
		}
		next, err := page.NextPage()
		if err != nil {
			return nil, err
		}
		page = next
	}
}

// Content blocks and SDK messages are defined in messages.go.

// RunStreamEvent is a low-level stream envelope.
type RunStreamEvent struct {
	Kind              string
	Offset            string
	SDKMessage        SDKMessage
	InteractionUpdate InteractionUpdate
	Step              ConversationStep
	Result            map[string]any
	Done              map[string]any
	ResultIsFull      bool
}

// Conversation types are in messages.go.

func parseUserMessage(v any) UserMessage {
	switch t := v.(type) {
	case string:
		return UserMessageFromText(t)
	case UserMessage:
		return t
	case map[string]any:
		msg := UserMessage{Text: stringField(t, "text")}
		if imgs, ok := t["images"].([]any); ok {
			for _, img := range imgs {
				if m, ok := img.(map[string]any); ok {
					msg.Images = append(msg.Images, SDKImage{URL: stringField(m, "url")})
				}
			}
		}
		return msg
	default:
		return UserMessage{}
	}
}
