// Package bridge launches cursor-sdk-bridge, the npm-installed adapter over @cursor/sdk.
package bridge

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/remdev/cursor-go-sdk/internal/connect"
)

const readyLinePrefix = "cursor-sdk-bridge ready "

// Endpoint describes a running bridge instance.
type Endpoint struct {
	URL           string
	AuthToken     string
	SchemaVersion int
	ServerVersion string
	PID           int
	WorkspaceRef  string
	StateRoot     string
}

// Bridge wraps a managed bridge subprocess.
type Bridge struct {
	Endpoint Endpoint
	cmd      *exec.Cmd
	stderr   io.ReadCloser
	mu       sync.Mutex
	closed   bool
}

// LaunchConfig configures bridge startup.
type LaunchConfig struct {
	Command     []string
	Workspace   string
	StateRoot   string
	Host        string
	Port        int
	Timeout     time.Duration
	ExtraArgs   []string
	LocalStore  map[string]any
	ToolCallbackURL    string
	ToolCallbackToken  string
}

// Launch starts cursor-sdk-bridge and waits for discovery on stderr.
func Launch(ctx context.Context, cfg LaunchConfig) (*Bridge, error) {
	if cfg.Timeout <= 0 {
		cfg.Timeout = 30 * time.Second
	}
	argv := cfg.Command
	if len(argv) == 0 {
		path, err := ResolvePath()
		if err != nil {
			return nil, err
		}
		argv = []string{path}
	}
	if cfg.Workspace != "" {
		argv = append(argv, "--workspace", cfg.Workspace)
	}
	if cfg.StateRoot != "" {
		argv = append(argv, "--state-root", cfg.StateRoot)
	}
	if cfg.Host != "" {
		argv = append(argv, "--host", cfg.Host)
	}
	if cfg.Port > 0 {
		argv = append(argv, "--port", fmt.Sprintf("%d", cfg.Port))
	}
	if cfg.ToolCallbackURL != "" {
		argv = append(argv, "--tool-callback-url", cfg.ToolCallbackURL, "--tool-callback-auth-token", cfg.ToolCallbackToken)
	}
	argv = append(argv, cfg.ExtraArgs...)

	cmd := exec.CommandContext(ctx, argv[0], argv[1:]...)
	cmd.Env = bridgeEnv()
	cmd.Stdout = io.Discard
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		stderr.Close()
		return nil, fmt.Errorf("start bridge: %w", err)
	}

	discovery, err := readDiscovery(stderr, cmd, cfg.Timeout)
	if err != nil {
		terminate(cmd)
		stderr.Close()
		return nil, err
	}
	endpoint, err := endpointFromDiscovery(discovery)
	if err != nil {
		terminate(cmd)
		stderr.Close()
		return nil, err
	}
	return &Bridge{
		Endpoint: endpoint,
		cmd:      cmd,
		stderr:   stderr,
	}, nil
}

// Close shuts down the bridge gracefully.
func (b *Bridge) Close() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.closed {
		return nil
	}
	b.closed = true
	connect.PostShutdown(b.Endpoint.URL, b.Endpoint.AuthToken, 0)
	if b.cmd != nil && b.cmd.Process != nil {
		terminate(b.cmd)
	}
	if b.stderr != nil {
		b.stderr.Close()
	}
	return nil
}

func bridgeEnv() []string {
	env := os.Environ()
	env = setDefaultEnv(env, "CURSOR_SDK_CLIENT_LANGUAGE", "go")
	return env
}

func setDefaultEnv(env []string, key, value string) []string {
	prefix := key + "="
	for _, e := range env {
		if strings.HasPrefix(e, prefix) {
			return env
		}
	}
	return append(env, prefix+value)
}

func readDiscovery(stderr io.Reader, cmd *exec.Cmd, timeout time.Duration) (map[string]any, error) {
	deadline := time.Now().Add(timeout)
	scanner := bufio.NewScanner(stderr)
	var lines []string
	for time.Now().Before(deadline) {
		if cmd.ProcessState != nil && cmd.ProcessState.Exited() {
			break
		}
		if scanner.Scan() {
			line := scanner.Text() + "\n"
			lines = append(lines, line)
			if d, ok := parseDiscoveryLine(line); ok {
				return d, nil
			}
			continue
		}
		if err := scanner.Err(); err != nil {
			return nil, err
		}
		if cmd.Process != nil {
			if state, err := cmd.Process.Wait(); err == nil && state.Exited() {
				break
			}
		}
		time.Sleep(50 * time.Millisecond)
	}
	return nil, fmt.Errorf("timed out waiting for bridge discovery: %s", strings.Join(lines, ""))
}

func parseDiscoveryLine(line string) (map[string]any, bool) {
	if !strings.HasPrefix(line, readyLinePrefix) {
		return nil, false
	}
	payload := strings.TrimSpace(line[len(readyLinePrefix):])
	var loaded map[string]any
	if err := json.Unmarshal([]byte(payload), &loaded); err != nil {
		return nil, false
	}
	return loaded, true
}

func endpointFromDiscovery(d map[string]any) (Endpoint, error) {
	if intValue(d["schemaVersion"]) != 1 {
		return Endpoint{}, fmt.Errorf("unsupported bridge discovery schema: %v", d["schemaVersion"])
	}
	if strValue(d["transport"]) != "tcp" || strValue(d["protocol"]) != "connect" {
		return Endpoint{}, fmt.Errorf("unsupported bridge transport discovery payload")
	}
	url := strValue(d["url"])
	if url == "" {
		host := strValue(d["host"])
		port := intValue(d["port"])
		if host == "" || port == 0 {
			return Endpoint{}, fmt.Errorf("bridge discovery payload is missing a URL")
		}
		if strings.Contains(host, ":") && !strings.HasPrefix(host, "[") {
			host = "[" + host + "]"
		}
		url = fmt.Sprintf("http://%s:%d", host, port)
	}
	token := strValue(d["authToken"])
	if token == "" {
		if tf := strValue(d["authTokenFile"]); tf != "" {
			b, err := os.ReadFile(tf)
			if err != nil {
				return Endpoint{}, err
			}
			token = strings.TrimSpace(string(b))
		}
	}
	if token == "" {
		return Endpoint{}, fmt.Errorf("bridge discovery payload is missing an auth token")
	}
	ep := Endpoint{
		URL:           url,
		AuthToken:     token,
		SchemaVersion: 1,
		ServerVersion: strValue(d["serverVersion"]),
		WorkspaceRef:  strValue(d["workspaceRef"]),
		StateRoot:     strValue(d["stateRoot"]),
	}
	if pid := intValue(d["pid"]); pid > 0 {
		ep.PID = pid
	}
	return ep, nil
}

func terminate(cmd *exec.Cmd) {
	if cmd == nil || cmd.Process == nil {
		return
	}
	cmd.Process.Signal(syscall.SIGTERM)
	done := make(chan struct{})
	go func() {
		cmd.Wait()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		cmd.Process.Kill()
		<-done
	}
}

func strValue(v any) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return fmt.Sprint(v)
}

func intValue(v any) int {
	switch n := v.(type) {
	case float64:
		return int(n)
	case int:
		return n
	case json.Number:
		i, _ := n.Int64()
		return int(i)
	default:
		return 0
	}
}

// ResolvePath locates the cursor-sdk-bridge launcher installed via npm or env override.
func ResolvePath() (string, error) {
	if override := strings.TrimSpace(os.Getenv("CURSOR_SDK_BRIDGE_BIN")); override != "" {
		if st, err := os.Stat(override); err != nil || st.IsDir() {
			return "", fmt.Errorf("CURSOR_SDK_BRIDGE_BIN=%q does not point to a file", override)
		}
		return override, nil
	}
	if root := strings.TrimSpace(os.Getenv("CURSOR_SDK_BRIDGE_ROOT")); root != "" {
		if path, err := launcherInRoot(root); err == nil {
			return path, nil
		}
	}
	if !optedOutOfRemoteBridgeLookup() {
		if path, err := exec.LookPath("cursor-sdk-bridge"); err == nil {
			return path, nil
		}
	}
	return "", bridgeNotFoundError()
}

func optedOutOfRemoteBridgeLookup() bool {
	v := strings.ToLower(strings.TrimSpace(os.Getenv("CURSOR_SDK_USE_REMOTE_BRIDGE")))
	return v == "1" || v == "true" || v == "yes" || v == "on"
}

// ParseDiscoveryLine parses a bridge ready line (exported for tests).
func ParseDiscoveryLine(line string) (map[string]any, error) {
	if d, ok := parseDiscoveryLine(line); ok {
		return d, nil
	}
	return nil, fmt.Errorf("not a discovery line")
}
