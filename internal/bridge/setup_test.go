package bridge_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/remdev/cursor-go-sdk/internal/bridge"
)

func TestFindBridgeDirFromRepoRoot(t *testing.T) {
	repoRoot, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(repoRoot, "bridge", "package.json")); err != nil {
		t.Skip("bridge/ not present")
	}
	dir, ok := bridge.FindBridgeDirForTest(repoRoot)
	if !ok {
		t.Fatal("expected to find bridge/")
	}
	if filepath.Base(dir) != "bridge" {
		t.Fatalf("unexpected dir: %q", dir)
	}
}

func TestFindBridgeDirMissing(t *testing.T) {
	dir, ok := bridge.FindBridgeDirForTest(t.TempDir())
	if ok {
		t.Fatalf("unexpected bridge dir: %q", dir)
	}
}

func TestLocalBridgeDirFromEnv(t *testing.T) {
	root := t.TempDir()
	bridgeDir := filepath.Join(root, "bridge")
	if err := os.MkdirAll(bridgeDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(
		filepath.Join(bridgeDir, "package.json"),
		[]byte(`{"name":"@cursor-go-sdk/cursor-sdk-bridge","version":"0.0.2"}`),
		0o644,
	); err != nil {
		t.Fatal(err)
	}

	t.Setenv("CURSOR_SDK_BRIDGE_ROOT", bridgeDir)
	t.Chdir(t.TempDir())

	got, err := bridge.LocalBridgeDirForTest(bridge.SetupOptions{Local: true})
	if err != nil {
		t.Fatal(err)
	}
	if got != bridgeDir {
		t.Fatalf("got %q want %q", got, bridgeDir)
	}
}

func TestLocalBridgeDirPrefersBridgeDirOption(t *testing.T) {
	t.Setenv("CURSOR_SDK_BRIDGE_ROOT", "/env/should/not/win")
	got, err := bridge.LocalBridgeDirForTest(bridge.SetupOptions{
		Local:     true,
		BridgeDir: "/explicit/bridge",
	})
	if err != nil {
		t.Fatal(err)
	}
	if got != "/explicit/bridge" {
		t.Fatalf("got %q", got)
	}
}

func TestValidateSetupVersionRejectsBelowMinimum(t *testing.T) {
	_, err := bridge.ValidateSetupVersionForTest("0.0.1")
	if err == nil {
		t.Fatal("expected error for version below minimum")
	}
	if !strings.Contains(err.Error(), "below minimum") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestValidateSetupVersionDefaultsToPinned(t *testing.T) {
	got, err := bridge.ValidateSetupVersionForTest("")
	if err != nil {
		t.Fatal(err)
	}
	if got != bridge.BridgeNpmVersion {
		t.Fatalf("got %q want %q", got, bridge.BridgeNpmVersion)
	}
}
