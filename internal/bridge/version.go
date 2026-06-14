package bridge

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const bridgeJSEntry = "dist/bin/cursor-sdk-bridge.js"

func ensureBridgeVersion(entry string) error {
	version, err := bridgePackageVersion(entry)
	if err != nil {
		return err
	}
	if semverAtLeast(version, BridgeNpmVersion) {
		return nil
	}
	return fmt.Errorf(
		"cursor-sdk-bridge %s is installed; need >= %s\n\nRun:\n\n\tgo run github.com/remdev/cursor-go-sdk/cmd/setup@latest\n",
		version,
		BridgeNpmVersion,
	)
}

func bridgePackageVersion(entry string) (string, error) {
	root, err := bridgePackageRoot(entry)
	if err != nil {
		return "", err
	}
	data, err := os.ReadFile(filepath.Join(root, "package.json"))
	if err != nil {
		return "", fmt.Errorf("read bridge package.json: %w", err)
	}
	var meta struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	}
	if err := json.Unmarshal(data, &meta); err != nil {
		return "", fmt.Errorf("parse bridge package.json: %w", err)
	}
	if meta.Name != bridgePackage {
		return "", fmt.Errorf("unexpected bridge package name %q", meta.Name)
	}
	if meta.Version == "" {
		return "", fmt.Errorf("bridge package.json is missing version")
	}
	return meta.Version, nil
}

func bridgePackageRoot(entry string) (string, error) {
	abs, err := filepath.Abs(entry)
	if err != nil {
		return "", err
	}
	jsSuffix := filepath.Join("dist", "bin", "cursor-sdk-bridge.js")
	if strings.HasSuffix(abs, jsSuffix) {
		return filepath.Dir(filepath.Dir(filepath.Dir(abs))), nil
	}
	binSuffix := filepath.Join("bin", launcherName())
	if strings.HasSuffix(abs, binSuffix) {
		return filepath.Dir(filepath.Dir(abs)), nil
	}
	dir := filepath.Dir(abs)
	for {
		if isBridgePackage(filepath.Join(dir, "package.json")) {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", fmt.Errorf("cannot locate bridge package root from %q", entry)
}

func semverAtLeast(version, min string) bool {
	vMajor, vMinor, vPatch, ok := parseSemver(version)
	if !ok {
		return false
	}
	mMajor, mMinor, mPatch, ok := parseSemver(min)
	if !ok {
		return false
	}
	if vMajor != mMajor {
		return vMajor > mMajor
	}
	if vMinor != mMinor {
		return vMinor > mMinor
	}
	return vPatch >= mPatch
}

func parseSemver(version string) (major, minor, patch int, ok bool) {
	version = strings.TrimPrefix(strings.TrimSpace(version), "v")
	if i := strings.IndexAny(version, "+-"); i >= 0 {
		version = version[:i]
	}
	parts := strings.Split(version, ".")
	if len(parts) < 3 {
		return 0, 0, 0, false
	}
	var err error
	if major, err = strconv.Atoi(parts[0]); err != nil {
		return 0, 0, 0, false
	}
	if minor, err = strconv.Atoi(parts[1]); err != nil {
		return 0, 0, 0, false
	}
	if patch, err = strconv.Atoi(parts[2]); err != nil {
		return 0, 0, 0, false
	}
	return major, minor, patch, true
}

// SemverAtLeastForTest exposes semverAtLeast for tests.
func SemverAtLeastForTest(version, min string) bool {
	return semverAtLeast(version, min)
}

// BridgePackageVersionForTest exposes bridgePackageVersion for tests.
func BridgePackageVersionForTest(entry string) (string, error) {
	return bridgePackageVersion(entry)
}
