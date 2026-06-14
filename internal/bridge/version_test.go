package bridge_test

import (
	"testing"

	"github.com/remdev/cursor-go-sdk/internal/bridge"
)

func TestSemverAtLeast(t *testing.T) {
	cases := []struct {
		version string
		min     string
		want    bool
	}{
		{"0.0.2", "0.0.2", true},
		{"0.0.3", "0.0.2", true},
		{"0.0.1", "0.0.2", false},
		{"1.0.0", "0.0.2", true},
	}
	for _, tc := range cases {
		if got := bridge.SemverAtLeastForTest(tc.version, tc.min); got != tc.want {
			t.Fatalf("%s >= %s: got %v want %v", tc.version, tc.min, got, tc.want)
		}
	}
}
