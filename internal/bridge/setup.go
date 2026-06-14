package bridge

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// SetupOptions configures bridge prerequisite installation.
type SetupOptions struct {
	// Local installs from a repository clone (npm ci && npm run build && npm link).
	Local bool
	// BridgeDir is the bridge package root; defaults to auto-discovery when Local is true.
	BridgeDir string
	// Version is the npm version for global install; defaults to BridgeNpmVersion.
	Version string
}

// Setup installs cursor-sdk-bridge when it is not already on PATH at the required version.
func Setup(ctx context.Context, opts SetupOptions) error {
	if err := ensureNodeRuntime(); err != nil {
		return err
	}
	if _, err := ResolvePath(); err == nil {
		return nil
	}

	if opts.Local || opts.BridgeDir != "" {
		return setupLocal(ctx, opts)
	}
	return setupGlobal(ctx, opts.Version)
}

func setupGlobal(ctx context.Context, version string) error {
	resolved, err := validateSetupVersion(version)
	if err != nil {
		return err
	}
	pkg := bridgePackage + "@" + resolved
	npm, err := npmBin()
	if err != nil {
		return err
	}
	cmd := exec.CommandContext(ctx, npm, "install", "-g", pkg, "--force")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("npm install -g %s: %w", pkg, err)
	}
	if _, err := ResolvePath(); err != nil {
		return setupGlobalResolveError(pkg, err)
	}
	return nil
}

func validateSetupVersion(version string) (string, error) {
	if version == "" {
		return BridgeNpmVersion, nil
	}
	if !semverAtLeast(version, BridgeNpmVersion) {
		return "", fmt.Errorf(
			"requested bridge version %s is below minimum %s",
			version,
			BridgeNpmVersion,
		)
	}
	return version, nil
}

func setupGlobalResolveError(pkg string, err error) error {
	if strings.Contains(err.Error(), "need >=") {
		return fmt.Errorf("installed %s but bridge version check failed: %w", pkg, err)
	}
	return fmt.Errorf("installed %s but cursor-sdk-bridge is still not on PATH: %w", pkg, err)
}

// ValidateSetupVersionForTest exposes validateSetupVersion for tests.
func ValidateSetupVersionForTest(version string) (string, error) {
	return validateSetupVersion(version)
}

func setupLocal(ctx context.Context, opts SetupOptions) error {
	bridgeDir, err := localBridgeDir(opts)
	if err != nil {
		return err
	}
	bridgeDir, err = filepath.Abs(bridgeDir)
	if err != nil {
		return err
	}
	if _, err := os.Stat(filepath.Join(bridgeDir, "package.json")); err != nil {
		return fmt.Errorf("bridge package not found at %s: %w", bridgeDir, err)
	}

	npm, err := npmBin()
	if err != nil {
		return err
	}
	steps := [][]string{
		{npm, "ci"},
		{npm, "run", "build"},
		{npm, "link"},
	}
	for _, step := range steps {
		cmd := exec.CommandContext(ctx, step[0], step[1:]...)
		cmd.Dir = bridgeDir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("%s in %s: %w", strings.Join(step, " "), bridgeDir, err)
		}
	}
	if _, err := ResolvePath(); err != nil {
		return fmt.Errorf("linked bridge from %s but cursor-sdk-bridge is still not on PATH: %w", bridgeDir, err)
	}
	return nil
}

func findBridgeDir(start string) (string, bool) {
	dir := start
	if dir == "" {
		var err error
		dir, err = os.Getwd()
		if err != nil {
			return "", false
		}
	}
	for {
		bridgeDir := filepath.Join(dir, "bridge")
		if isBridgePackage(filepath.Join(bridgeDir, "package.json")) {
			return bridgeDir, true
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", false
		}
		dir = parent
	}
}

func isBridgePackage(path string) bool {
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	var meta struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal(data, &meta); err != nil {
		return false
	}
	return meta.Name == bridgePackage
}

func npmBin() (string, error) {
	name := "npm"
	if runtime.GOOS == "windows" {
		name = "npm.cmd"
	}
	path, err := exec.LookPath(name)
	if err != nil {
		return "", fmt.Errorf("npm not found on PATH: install Node.js >= 18 (includes npm) or add npm to PATH")
	}
	return path, nil
}

func localBridgeDir(opts SetupOptions) (string, error) {
	if opts.BridgeDir != "" {
		return opts.BridgeDir, nil
	}
	if root := strings.TrimSpace(os.Getenv("CURSOR_SDK_BRIDGE_ROOT")); root != "" {
		return root, nil
	}
	bridgeDir, ok := findBridgeDir("")
	if !ok {
		return "", fmt.Errorf(
			"bridge/ not found; run from a cursor-go-sdk clone, set CURSOR_SDK_BRIDGE_ROOT, or pass --bridge-dir",
		)
	}
	return bridgeDir, nil
}

// FindBridgeDirForTest exposes findBridgeDir for tests.
func FindBridgeDirForTest(start string) (string, bool) {
	return findBridgeDir(start)
}

// LocalBridgeDirForTest exposes localBridgeDir for tests.
func LocalBridgeDirForTest(opts SetupOptions) (string, error) {
	return localBridgeDir(opts)
}
