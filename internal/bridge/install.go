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
	b.WriteString("Install prerequisites:\n\n")
	b.WriteString("\tgo run github.com/remdev/cursor-go-sdk/cmd/setup@latest\n\n")
	b.WriteString("Or manually:\n\n")
	b.WriteString("\tnpm install -g ")
	b.WriteString(bridgePackage)
	b.WriteString("@")
	b.WriteString(BridgeNpmVersion)
	b.WriteString("\n\n")
	b.WriteString("From a clone of cursor-go-sdk:\n\n")
	b.WriteString("\tgo run ./cmd/setup --local\n\n")
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
