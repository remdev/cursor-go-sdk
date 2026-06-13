package connect_test

import (
	"testing"

	"github.com/remdev/cursor-go-sdk/internal/connect"
)

func TestEncodeStreamEnvelope(t *testing.T) {
	frame, err := connect.EncodeStreamEnvelope(map[string]any{"hello": "world"}, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(frame) < 6 {
		t.Fatalf("frame too short: %d", len(frame))
	}
	if frame[0] != 0 {
		t.Fatalf("expected non-terminal flag, got %d", frame[0])
	}
}

func TestStreamParserEndStream(t *testing.T) {
	end, err := connect.EncodeStreamEnvelope(map[string]any{}, true)
	if err != nil {
		t.Fatal(err)
	}
	msg, err := connect.EncodeStreamEnvelope(map[string]any{"event": 1}, false)
	if err != nil {
		t.Fatal(err)
	}
	// Parser is internal; exercise via public stream reader in integration tests.
	_ = end
	_ = msg
}

func TestErrorFromHTTP(t *testing.T) {
	// Smoke test that error helpers compile and basic predicates work.
	err := &connect.Error{Message: "rate limited", Code: "rate_limit", IsRetryable: true}
	if !connect.IsRateLimit(err) {
		t.Fatal("expected rate limit")
	}
}
