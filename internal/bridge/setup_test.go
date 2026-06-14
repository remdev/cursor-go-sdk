package bridge_test

import (
	"os"
	"path/filepath"
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
