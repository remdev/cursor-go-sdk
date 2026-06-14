//go:build e2e

package e2e_test

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestExampleBasic(t *testing.T) {
	requireE2E(t)
	root := moduleRoot(t)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	cmd := exec.CommandContext(ctx, "go", "run", "./examples/basic")
	cmd.Dir = root
	cmd.Env = append(os.Environ(),
		apiKeyEnvVar+"="+e2eAPIKey(t),
		"CURSOR_MODEL="+e2eModel(t),
	)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		t.Fatalf("go run ./examples/basic failed: %v\nstderr:\n%s\nstdout:\n%s", err, stderr.String(), stdout.String())
	}
	combined := stdout.String() + stderr.String()
	if !strings.Contains(combined, "status:") || !containsPong(combined) {
		t.Fatalf("unexpected output:\n%s", combined)
	}
}

func TestExampleCodingAgentCLI(t *testing.T) {
	requireE2E(t)
	root := moduleRoot(t)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	cmd := exec.CommandContext(ctx, "go", "run", "./examples/coding-agent-cli", "--", "Reply with exactly: pong")
	cmd.Dir = root
	cmd.Env = append(os.Environ(),
		apiKeyEnvVar+"="+e2eAPIKey(t),
		"CURSOR_MODEL="+e2eModel(t),
	)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		t.Fatalf("go run ./examples/coding-agent-cli failed: %v\nstderr:\n%s\nstdout:\n%s", err, stderr.String(), stdout.String())
	}
	combined := stdout.String() + stderr.String()
	if !containsPong(combined) {
		t.Fatalf("unexpected output:\nstdout:\n%s\nstderr:\n%s", stdout.String(), stderr.String())
	}
}

func TestExampleQuickstartBuild(t *testing.T) {
	requireE2EEnabled(t)
	root := moduleRoot(t)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "go", "build", "-o", filepath.Join(t.TempDir(), "quickstart"), "./examples/quickstart")
	cmd.Dir = root
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("go build ./examples/quickstart: %v\n%s", err, out)
	}
}
