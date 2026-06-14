package bridge

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ResolveLaunchArgv returns argv to start the bridge.
func ResolveLaunchArgv() ([]string, error) {
	entry, err := ResolvePath()
	if err != nil {
		return nil, err
	}
	return launchArgvForEntrypoint(entry)
}

func launchArgvForEntrypoint(entry string) ([]string, error) {
	if strings.HasSuffix(entry, ".js") {
		node, err := nodeBin()
		if err != nil {
			return nil, err
		}
		return []string{node, entry}, nil
	}
	return []string{entry}, nil
}

func resolveBridgeEntry(launcher string) (string, error) {
	resolved := launcher
	if link, err := filepath.EvalSymlinks(launcher); err == nil {
		resolved = link
	}
	abs, err := filepath.Abs(resolved)
	if err != nil {
		return "", err
	}
	st, err := os.Stat(abs)
	if err != nil {
		return "", fmt.Errorf("cursor-sdk-bridge entry %q: %w", abs, err)
	}
	if st.IsDir() {
		return "", fmt.Errorf("cursor-sdk-bridge entry %q is a directory", abs)
	}
	return abs, nil
}

func nodeBin() (string, error) {
	if override := strings.TrimSpace(os.Getenv("CURSOR_SDK_NODE_BIN")); override != "" {
		return override, nil
	}
	path, err := exec.LookPath("node")
	if err != nil {
		return "", fmt.Errorf("node not found on PATH: install Node.js >= 18 or set CURSOR_SDK_NODE_BIN")
	}
	return path, nil
}

// ResolveLaunchArgvForTest exposes ResolveLaunchArgv for tests.
func ResolveLaunchArgvForTest() ([]string, error) {
	return ResolveLaunchArgv()
}
