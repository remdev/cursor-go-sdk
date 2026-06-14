//go:build e2e

package e2e_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/remdev/cursor-go-sdk/cursor"
)

const (
	e2eEnvVar    = "CURSOR_E2E"
	apiKeyEnvVar = "CURSOR_API_KEY"
	modelEnvVar  = "CURSOR_E2E_MODEL"
)

func TestMain(m *testing.M) {
	code := m.Run()
	if os.Getenv(e2eEnvVar) == "1" {
		_ = cursor.CloseDefaultClient()
	}
	os.Exit(code)
}

func requireE2E(t *testing.T) {
	t.Helper()
	requireE2EEnabled(t)
	requireAPIKey(t)
}

func requireE2EEnabled(t *testing.T) {
	t.Helper()
	if os.Getenv(e2eEnvVar) != "1" {
		t.Skip("local e2e disabled: set CURSOR_E2E=1")
	}
}

func requireAPIKey(t *testing.T) {
	t.Helper()
	if strings.TrimSpace(os.Getenv(apiKeyEnvVar)) == "" {
		t.Skip("CURSOR_API_KEY is not set")
	}
}

func e2eContext(t *testing.T) context.Context {
	t.Helper()
	timeout := 3 * time.Minute
	if raw := strings.TrimSpace(os.Getenv("CURSOR_E2E_TIMEOUT")); raw != "" {
		if d, err := time.ParseDuration(raw); err == nil && d > 0 {
			timeout = d
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	t.Cleanup(cancel)
	return ctx
}

func e2eClient(t *testing.T) *cursor.Client {
	t.Helper()
	requireBridge(t, context.Background())
	client, err := cursor.DefaultClient(context.Background())
	if err != nil {
		t.Fatalf("DefaultClient: %v", err)
	}
	return client
}

func e2eAPIKey(t *testing.T) string {
	t.Helper()
	requireE2E(t)
	return strings.TrimSpace(os.Getenv(apiKeyEnvVar))
}

func e2eModel(t *testing.T) string {
	t.Helper()
	if model := strings.TrimSpace(os.Getenv(modelEnvVar)); model != "" {
		return model
	}
	if model := strings.TrimSpace(os.Getenv("CURSOR_MODEL")); model != "" {
		return model
	}
	return "auto"
}

func e2eWorkspace(t *testing.T) string {
	t.Helper()
	if dir := strings.TrimSpace(os.Getenv("CURSOR_E2E_WORKSPACE")); dir != "" {
		return dir
	}
	return moduleRoot(t)
}

func localAgentOptions(t *testing.T) cursor.AgentOptions {
	t.Helper()
	return cursor.AgentOptions{
		Model:  e2eModel(t),
		APIKey: e2eAPIKey(t),
		Local:  &cursor.LocalAgentOptions{CWD: []string{e2eWorkspace(t)}},
	}
}

func moduleRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), ".."))
}

func requireBridge(t *testing.T, ctx context.Context) {
	t.Helper()
	if err := cursor.EnsureBridgeInstalled(ctx); err != nil {
		t.Fatalf("cursor-sdk-bridge not available: %v\n\nRun: go run ./cmd/setup --local", err)
	}
}

func containsPong(text string) bool {
	return strings.Contains(strings.ToLower(text), "pong")
}

func runLocalPromptText(t *testing.T, ctx context.Context, prompt string, opts cursor.AgentOptions) (cursor.RunResult, string) {
	t.Helper()
	agent, err := cursor.CreateAgent(ctx, opts)
	if err != nil {
		t.Fatalf("CreateAgent: %v", err)
	}
	defer agent.Close(ctx)

	run, err := agent.Send(ctx, prompt, cursor.SendOptions{})
	if err != nil {
		t.Fatalf("Send: %v", err)
	}

	var assistant strings.Builder
	for msg, err := range run.Messages(ctx) {
		if err != nil {
			t.Fatalf("Messages: %v", err)
		}
		if msg.Type == "assistant" {
			assistant.WriteString(cursor.AssistantText(msg))
		}
	}

	result, err := run.Wait(ctx)
	if err != nil {
		t.Fatalf("Wait: %v", err)
	}

	text := strings.TrimSpace(result.Result)
	if text == "" {
		text = strings.TrimSpace(assistant.String())
	}
	return result, text
}

func versionField(version map[string]any, keys ...string) string {
	for _, key := range keys {
		if v, ok := version[key]; ok && v != nil {
			if s := strings.TrimSpace(fmt.Sprint(v)); s != "" {
				return s
			}
		}
	}
	return ""
}
