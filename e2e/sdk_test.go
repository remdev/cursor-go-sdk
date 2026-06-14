//go:build e2e

package e2e_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/remdev/cursor-go-sdk/cursor"
)

func TestBridgeGetVersion(t *testing.T) {
	requireE2EEnabled(t)
	client := e2eClient(t)
	ctx := e2eContext(t)

	version, err := client.GetVersion(ctx)
	if err != nil {
		t.Fatalf("GetVersion: %v", err)
	}
	if versionField(version, "bridgeVersion", "bridge_version") == "" {
		t.Fatalf("GetVersion returned no bridge version: %v", version)
	}
	if versionField(version, "protocolVersion", "protocol_version") == "" {
		t.Fatalf("GetVersion returned no protocol version: %v", version)
	}
}

func TestBridgeSurvivesAfterGetVersion(t *testing.T) {
	requireE2E(t)
	client := e2eClient(t)

	shortCtx, cancel := context.WithTimeout(context.Background(), time.Second)
	if _, err := client.GetVersion(shortCtx); err != nil {
		t.Fatalf("GetVersion: %v", err)
	}
	cancel()

	if _, err := client.GetVersion(e2eContext(t)); err != nil {
		t.Fatalf("bridge died after caller context was canceled: %v", err)
	}
}

func TestCursorMe(t *testing.T) {
	requireE2E(t)
	e2eClient(t)
	ctx := e2eContext(t)

	user, err := cursor.Cursor.Me(ctx, cursor.CursorRequestOptions{APIKey: e2eAPIKey(t)})
	if err != nil {
		t.Fatalf("Cursor.Me: %v", err)
	}
	if strings.TrimSpace(user.UserEmail) == "" && user.UserID == nil {
		t.Fatalf("Cursor.Me returned empty user: %+v", user)
	}
}

func TestPromptPong(t *testing.T) {
	requireE2E(t)
	e2eClient(t)
	ctx := e2eContext(t)

	result, text := runLocalPromptText(t, ctx, "Reply with exactly: pong", localAgentOptions(t))
	if result.Status != cursor.RunStatusFinished {
		t.Fatalf("status=%q want finished", result.Status)
	}
	if !containsPong(text) {
		t.Fatalf("response=%q want pong", text)
	}
}

func TestCreateAgentSendWait(t *testing.T) {
	requireE2E(t)
	e2eClient(t)
	ctx := e2eContext(t)

	opts := localAgentOptions(t)
	opts.Name = "e2e send/wait"

	result, text := runLocalPromptText(t, ctx, "Reply with one short sentence confirming the SDK works.", opts)
	if result.Status != cursor.RunStatusFinished {
		t.Fatalf("status=%q want finished", result.Status)
	}
	if strings.TrimSpace(text) == "" {
		t.Fatal("empty assistant response")
	}
}
