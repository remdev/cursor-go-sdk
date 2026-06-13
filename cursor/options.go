package cursor

import "sync"

var (
	configureMu sync.RWMutex
	configure   ConfigureOptions
)

// ConfigureOptions sets process-wide SDK defaults for later Agent calls.
// Per-call options override these values.
type ConfigureOptions struct {
	Local *ConfigureLocalOptions
}

// ConfigureLocalOptions configures local agent defaults.
type ConfigureLocalOptions struct {
	Store *LocalAgentStoreConfig
}

// Configure sets process-wide SDK defaults matching Cursor.configure() in TypeScript.
func Configure(opts ConfigureOptions) {
	configureMu.Lock()
	defer configureMu.Unlock()
	configure = opts
}

func applyConfigureDefaults(opts *AgentOptions) {
	if opts == nil || opts.Local == nil {
		return
	}
	configureMu.RLock()
	cfg := configure
	configureMu.RUnlock()
	if cfg.Local == nil || cfg.Local.Store == nil {
		return
	}
	if opts.Local.Store == nil {
		store := *cfg.Local.Store
		opts.Local.Store = &store
	}
}

// LocalOptions builds LocalAgentOptions from one or more workspace paths.
func LocalOptions(cwd ...string) *LocalAgentOptions {
	return &LocalAgentOptions{CWD: cwd}
}

// GetAgentOptions configures Agent.get lookups.
type GetAgentOptions struct {
	CWD    string
	APIKey string
}

func (o GetAgentOptions) ToWire() map[string]any {
	return dropEmpty(map[string]any{
		"cwd": o.CWD, "apiKey": o.APIKey,
	})
}

// GetRunOptions configures Agent.getRun lookups.
type GetRunOptions struct {
	Runtime string
	CWD     string
	AgentID string
	APIKey  string
}

func (o GetRunOptions) ToWire() map[string]any {
	return dropEmpty(map[string]any{
		"runtime": o.Runtime,
		"cwd":     o.CWD,
		"agentId": o.AgentID,
		"apiKey":  o.APIKey,
	})
}

// AgentOperationOptions configures archive/unarchive/delete operations.
type AgentOperationOptions struct {
	CWD    string
	APIKey string
}

func (o AgentOperationOptions) ToWire() map[string]any {
	return dropEmpty(map[string]any{
		"cwd": o.CWD, "apiKey": o.APIKey,
	})
}

// RunOperation names a run capability.
type RunOperation string

const (
	RunOpStream       RunOperation = "stream"
	RunOpWait         RunOperation = "wait"
	RunOpCancel       RunOperation = "cancel"
	RunOpConversation RunOperation = "conversation"
)

// SDKUserMessage is the structured send message type (alias for UserMessage).
type SDKUserMessage = UserMessage

// SDKArtifact is an alias matching the TypeScript export name.
type SDKArtifact = SdkArtifact

// CursorRequestOptions carries optional apiKey for Cursor namespace reads.
type CursorRequestOptions struct {
	APIKey string
}

// DeltaCallback receives raw interaction updates during streaming.
type DeltaCallback func(InteractionUpdate)

// StepCallback receives completed conversation steps during streaming.
type StepCallback func(ConversationStep)
