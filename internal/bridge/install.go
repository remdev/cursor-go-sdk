package bridge

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const bridgePackage = "@cursor-go-sdk/cursor-sdk-bridge"

// EnsureInstalled verifies that cursor-sdk-bridge is available.
func EnsureInstalled(context.Context) error {
	_, err := ResolvePath()
	if err != nil {
		return err
	}
	return ensureNodeRuntime()
}

func bridgeNotFoundError() error {
	var b strings.Builder
	b.WriteString("cursor-sdk-bridge not found.\n\n")
	b.WriteString("Install the bridge prerequisite:\n\n")
	b.WriteString("\tnpm install -g ")
	b.WriteString(bridgePackage)
	b.WriteString("\n\n")
	b.WriteString("From a clone of cursor-go-sdk:\n\n")
	b.WriteString("\tcd bridge && npm ci && npm link\n\n")
	b.WriteString("Or set CURSOR_SDK_BRIDGE_BIN to the launcher path, or CURSOR_SDK_BRIDGE_ROOT to an installed bridge directory.\n")
	b.WriteString("Requires Node.js >= 18.\n")
	return fmt.Errorf("%s", strings.TrimRight(b.String(), "\n"))
}

func launcherInRoot(root string) (string, error) {
	path := filepath.Join(root, "bin", launcherName())
	st, err := os.Stat(path)
	if err != nil {
		return "", err
	}
	if st.IsDir() {
		return "", fmt.Errorf("%q is a directory", path)
	}
	return path, nil
}

func launcherName() string {
	if runtime.GOOS == "windows" {
		return "cursor-sdk-bridge.cmd"
	}
	return "cursor-sdk-bridge"
}
