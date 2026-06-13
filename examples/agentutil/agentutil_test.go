package agentutil

import "testing"

func TestCompact(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"  hello   world  ", "hello world"},
		{"", ""},
		{"single", "single"},
	}

	for _, tt := range tests {
		if got := Compact(tt.in); got != tt.want {
			t.Fatalf("Compact(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestEnvOr(t *testing.T) {
	t.Setenv("AGENTUTIL_TEST_KEY", "  value  ")
	if got := EnvOr("AGENTUTIL_TEST_KEY", "fallback"); got != "value" {
		t.Fatalf("got %q", got)
	}
	if got := EnvOr("AGENTUTIL_MISSING_KEY", "fallback"); got != "fallback" {
		t.Fatalf("got %q", got)
	}
}
