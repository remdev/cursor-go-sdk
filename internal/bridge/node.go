package bridge

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func ensureNodeRuntime() error {
	nodeBin := strings.TrimSpace(os.Getenv("CURSOR_SDK_NODE_BIN"))
	if nodeBin == "" {
		var err error
		nodeBin, err = exec.LookPath("node")
		if err != nil {
			return fmt.Errorf("node not found on PATH: install Node.js >= 18 or set CURSOR_SDK_NODE_BIN")
		}
	}
	cmd := exec.Command(nodeBin, "--version")
	out, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("node runtime check failed: %w", err)
	}
	major, err := parseNodeMajorVersion(strings.TrimSpace(string(out)))
	if err != nil {
		return err
	}
	if major < 18 {
		return fmt.Errorf("node %s is too old; need >= 18", strings.TrimSpace(string(out)))
	}
	return nil
}

func parseNodeMajorVersion(version string) (int, error) {
	version = strings.TrimPrefix(version, "v")
	parts := strings.SplitN(version, ".", 2)
	if len(parts) == 0 || parts[0] == "" {
		return 0, fmt.Errorf("invalid node version %q", version)
	}
	var major int
	if _, err := fmt.Sscanf(parts[0], "%d", &major); err != nil {
		return 0, fmt.Errorf("invalid node version %q", version)
	}
	return major, nil
}
