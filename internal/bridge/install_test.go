package bridge_test

import (
	"os"
	"path/filepath"
	"runtime"
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

func TestModuleBridgeDir(t *testing.T) {
	dir := bridge.ModuleBridgeDirForTest()
	if dir == "" {
		t.Fatal("empty module bridge dir")
	}
	if _, err := os.Stat(filepath.Join(dir, "package.json")); err != nil {
		t.Fatalf("bridge/package.json: %v", err)
	}
}

func TestValidBridgeRootRequiresNodeModules(t *testing.T) {
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "manifest.json"), []byte(`{"sdkVersion":"1.0.18"}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if bridge.ValidBridgeRootForTest(root) {
		t.Fatal("expected invalid without node_modules")
	}
}

func TestResolvePathFromRepoBridge(t *testing.T) {
	root := bridge.ModuleBridgeDirForTest()
	if !bridge.ValidBridgeRootForTest(root) {
		t.Skip("run `cd bridge && npm install` to enable integration test")
	}
	path, err := bridge.ResolvePath()
	if err != nil {
		t.Fatal(err)
	}
	if path == "" {
		t.Fatal("empty bridge path")
	}
}
