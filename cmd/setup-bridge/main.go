package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	bridgeDir := flag.String("dir", "", "bridge directory (default: ./bridge or module bridge/)")
	flag.Parse()

	dir := *bridgeDir
	if dir == "" {
		dir = detectBridgeDir()
	}
	if dir == "" {
		fmt.Fprintln(os.Stderr, "setup-bridge: could not find bridge/ directory; pass -dir")
		os.Exit(1)
	}

	npm, err := exec.LookPath("npm")
	if err != nil {
		fmt.Fprintln(os.Stderr, "setup-bridge: npm not found on PATH")
		os.Exit(1)
	}

	cmd := exec.Command(npm, "install", "--omit=dev", "--no-fund", "--no-audit")
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		os.Exit(1)
	}

	fmt.Printf("\nBridge ready at: %s\n", dir)
	fmt.Printf("export CURSOR_SDK_BRIDGE_ROOT=%q\n", dir)
}

func detectBridgeDir() string {
	if cwd, err := os.Getwd(); err == nil {
		candidate := filepath.Join(cwd, "bridge")
		if st, err := os.Stat(filepath.Join(candidate, "package.json")); err == nil && !st.IsDir() {
			return candidate
		}
	}
	if modDir, err := exec.Command("go", "list", "-f", "{{.Dir}}", "-m", "github.com/remdev/cursor-go-sdk").Output(); err == nil {
		candidate := filepath.Join(string(trimLine(modDir)), "bridge")
		if st, err := os.Stat(filepath.Join(candidate, "package.json")); err == nil && !st.IsDir() {
			return candidate
		}
	}
	return ""
}

func trimLine(b []byte) string {
	for len(b) > 0 && (b[len(b)-1] == '\n' || b[len(b)-1] == '\r') {
		b = b[:len(b)-1]
	}
	return string(b)
}
