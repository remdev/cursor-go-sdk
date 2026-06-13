package bridge_test

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/remdev/cursor-go-sdk/internal/bridge"
)

func TestCurrentPlatformMatchesRuntime(t *testing.T) {
	t.Parallel()
	_, err := bridge.EnsureInstalledForTestPlatform()
	if err != nil {
		t.Fatalf("platform mapping: %v", err)
	}
}

func TestNpmPlatformPackageDarwinARM64(t *testing.T) {
	t.Parallel()
	if runtime.GOOS != "darwin" || runtime.GOARCH != "arm64" {
		t.Skip("platform-specific assertion")
	}
	pkg, err := bridge.NpmPlatformPackageForTest()
	if err != nil {
		t.Fatal(err)
	}
	if pkg != "@cursor/sdk-darwin-arm64" {
		t.Fatalf("got %q", pkg)
	}
}

func TestResolvePathFromBridgeRoot(t *testing.T) {
	repoBridge, err := filepath.Abs(filepath.Join("..", "..", "bridge"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(repoBridge, "package.json")); err != nil {
		t.Skip("bridge/ not present")
	}
	t.Setenv("CURSOR_SDK_BRIDGE_BIN", "")
	t.Setenv("CURSOR_SDK_BRIDGE_ROOT", repoBridge)
	t.Setenv("PATH", t.TempDir())

	resolved, err := bridge.ResolvePath()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(resolved, "cursor-sdk-bridge") {
		t.Fatalf("unexpected path: %q", resolved)
	}
}

func TestResolvePathMissingBridge(t *testing.T) {
	t.Setenv("PATH", t.TempDir())
	t.Setenv("CURSOR_SDK_BRIDGE_BIN", "")
	t.Setenv("CURSOR_SDK_BRIDGE_ROOT", "")
	_, err := bridge.ResolvePath()
	if err == nil {
		t.Fatal("expected error when bridge is missing")
	}
}
