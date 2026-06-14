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
	v, ok := parseSemverParts(version)
	if !ok {
		return false
	}
	m, ok := parseSemverParts(min)
	if !ok {
		return false
	}
	if v.major != m.major {
		return v.major > m.major
	}
	if v.minor != m.minor {
		return v.minor > m.minor
	}
	if v.patch != m.patch {
		return v.patch > m.patch
	}
	if v.prerelease == "" {
		return true
	}
	if m.prerelease == "" {
		return false
	}
	return comparePrerelease(v.prerelease, m.prerelease) >= 0
}

type semverParts struct {
	major, minor, patch int
	prerelease          string
}

func parseSemverParts(version string) (semverParts, bool) {
	version = strings.TrimPrefix(strings.TrimSpace(version), "v")
	if i := strings.Index(version, "+"); i >= 0 {
		version = version[:i]
	}
	prerelease := ""
	if i := strings.Index(version, "-"); i >= 0 {
		prerelease = version[i+1:]
		version = version[:i]
	}
	parts := strings.Split(version, ".")
	if len(parts) < 3 {
		return semverParts{}, false
	}
	var err error
	var p semverParts
	if p.major, err = strconv.Atoi(parts[0]); err != nil {
		return semverParts{}, false
	}
	if p.minor, err = strconv.Atoi(parts[1]); err != nil {
		return semverParts{}, false
	}
	if p.patch, err = strconv.Atoi(parts[2]); err != nil {
		return semverParts{}, false
	}
	p.prerelease = prerelease
	return p, true
}

func comparePrerelease(a, b string) int {
	if a == b {
		return 0
	}
	aParts := strings.Split(a, ".")
	bParts := strings.Split(b, ".")
	n := len(aParts)
	if len(bParts) > n {
		n = len(bParts)
	}
	for i := 0; i < n; i++ {
		var ap, bp string
		if i < len(aParts) {
			ap = aParts[i]
		}
		if i < len(bParts) {
			bp = bParts[i]
		}
		if ap == bp {
			continue
		}
		if ap == "" {
			return -1
		}
		if bp == "" {
			return 1
		}
		if cmp := comparePrereleaseIdentifier(ap, bp); cmp != 0 {
			return cmp
		}
	}
	return 0
}

func comparePrereleaseIdentifier(a, b string) int {
	aNum, aErr := strconv.Atoi(a)
	bNum, bErr := strconv.Atoi(b)
	if aErr == nil && bErr == nil {
		return aNum - bNum
	}
	if aErr == nil {
		return -1
	}
	if bErr == nil {
		return 1
	}
	return strings.Compare(a, b)
}

func parseSemver(version string) (major, minor, patch int, ok bool) {
	p, ok := parseSemverParts(version)
	if !ok {
		return 0, 0, 0, false
	}
	return p.major, p.minor, p.patch, true
}

// SemverAtLeastForTest exposes semverAtLeast for tests.
func SemverAtLeastForTest(version, min string) bool {
	return semverAtLeast(version, min)
}

// BridgePackageVersionForTest exposes bridgePackageVersion for tests.
func BridgePackageVersionForTest(entry string) (string, error) {
	return bridgePackageVersion(entry)
}
