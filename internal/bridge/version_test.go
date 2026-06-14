package bridge_test

import (
	"os"
	"path/filepath"
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
		{"0.0.2+meta", "0.0.2", true},
		{"0.0.2-beta", "0.0.2", true},
		{"0.0.1-rc1", "0.0.2", false},
		{"v0.0.3", "0.0.2", true},
	}
	for _, tc := range cases {
		if got := bridge.SemverAtLeastForTest(tc.version, tc.min); got != tc.want {
			t.Fatalf("%s >= %s: got %v want %v", tc.version, tc.min, got, tc.want)
		}
	}
}

func TestBridgePackageVersionWalkUp(t *testing.T) {
	root := t.TempDir()
	pkg := filepath.Join(root, "lib", "node_modules", "@cursor-go-sdk", "cursor-sdk-bridge")
	jsDir := filepath.Join(pkg, "dist", "bin")
	if err := os.MkdirAll(jsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	js := filepath.Join(jsDir, "cursor-sdk-bridge.js")
	if err := os.WriteFile(js, []byte("#!/usr/bin/env node\n"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(
		filepath.Join(pkg, "package.json"),
		[]byte(`{"name":"@cursor-go-sdk/cursor-sdk-bridge","version":"0.0.2+build"}`),
		0o644,
	); err != nil {
		t.Fatal(err)
	}

	version, err := bridge.BridgePackageVersionForTest(js)
	if err != nil {
		t.Fatal(err)
	}
	if version != "0.0.2+build" {
		t.Fatalf("got version %q", version)
	}
	if !bridge.SemverAtLeastForTest(version, "0.0.2") {
		t.Fatal("expected 0.0.2+build to satisfy >= 0.0.2")
	}
}
