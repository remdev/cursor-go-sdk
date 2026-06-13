package bridge

import (
	"archive/tar"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

var (
	resolveMu    sync.Mutex
	resolvedRoot string
	moduleBridge string
)

func init() {
	_, file, _, ok := runtime.Caller(0)
	if ok {
		moduleBridge = filepath.Clean(filepath.Join(filepath.Dir(file), "..", "..", "bridge"))
	}
}

// EnsureInstalled verifies that the bridge runtime is present (npm deps installed).
func EnsureInstalled(context.Context) error {
	_, err := ResolvePath()
	return err
}

// BridgeRoot returns the bridge directory containing bin/, dist/, and node_modules/.
func BridgeRoot() (string, error) {
	root, err := resolveBridgeRoot()
	if err != nil {
		return "", err
	}
	return root, nil
}

func resolveBridgeRoot() (string, error) {
	resolveMu.Lock()
	defer resolveMu.Unlock()
	if resolvedRoot != "" && validBridgeRoot(resolvedRoot) {
		return resolvedRoot, nil
	}

	candidates, err := bridgeRootCandidates()
	if err != nil {
		return "", err
	}
	for _, root := range candidates {
		if validBridgeRoot(root) {
			resolvedRoot = root
			return root, nil
		}
	}
	return "", bridgeNotReadyError(candidates)
}

func bridgeRootCandidates() ([]string, error) {
	var out []string
	seen := map[string]struct{}{}
	add := func(path string) {
		path = filepath.Clean(path)
		if path == "" || path == "." {
			return
		}
		if _, ok := seen[path]; ok {
			return
		}
		seen[path] = struct{}{}
		out = append(out, path)
	}

	if root := strings.TrimSpace(os.Getenv("CURSOR_SDK_BRIDGE_ROOT")); root != "" {
		add(root)
	}
	if moduleBridge != "" {
		add(moduleBridge)
	}
	if cwd, err := os.Getwd(); err == nil {
		for dir := cwd; ; dir = filepath.Dir(dir) {
			add(filepath.Join(dir, "bridge"))
			parent := filepath.Dir(dir)
			if parent == dir {
				break
			}
		}
	}
	if cache := embeddedBridgeCacheRoot(); cache != "" {
		add(cache)
	}
	return out, nil
}

func bridgeNotReadyError(candidates []string) error {
	var b strings.Builder
	b.WriteString("cursor-sdk-bridge is not installed.\n\n")
	b.WriteString("Install npm dependencies:\n\n")
	b.WriteString("\tcd bridge && npm install\n\n")
	b.WriteString("When using the module from GOPATH/GOMODCACHE:\n\n")
	b.WriteString("\tcd \"$(go list -f '{{.Dir}}' -m github.com/remdev/cursor-go-sdk)/bridge\" && npm install\n")
	b.WriteString("\texport CURSOR_SDK_BRIDGE_ROOT=\"$(go list -f '{{.Dir}}' -m github.com/remdev/cursor-go-sdk)/bridge\"\n\n")
	b.WriteString("Requires Node.js >= 18 (@cursor/sdk). Then set CURSOR_SDK_BRIDGE_ROOT if needed.\n")
	if len(candidates) > 0 {
		b.WriteString("\nChecked:\n")
		for _, c := range candidates {
			b.WriteString("\t- ")
			b.WriteString(c)
			b.WriteByte('\n')
		}
	}
	return fmt.Errorf("%s", strings.TrimRight(b.String(), "\n"))
}

func validBridgeRoot(root string) bool {
	sdkModule := filepath.Join(root, "node_modules", "@cursor", "sdk", "package.json")
	if _, err := os.Stat(sdkModule); err != nil {
		return false
	}
	launcher, err := launcherInRoot(root)
	if err != nil || launcher == "" {
		return false
	}
	manifestPath := filepath.Join(root, "manifest.json")
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return true
	}
	var manifest struct {
		SDKVersion string `json:"sdkVersion"`
	}
	if err := json.Unmarshal(data, &manifest); err != nil {
		return false
	}
	return manifest.SDKVersion == "" || manifest.SDKVersion == BundledSDKVersion
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

func embeddedBridgeCacheRoot() string {
	if len(embeddedBridgeTarball()) == 0 {
		return ""
	}
	if override := strings.TrimSpace(os.Getenv("CURSOR_SDK_BRIDGE_CACHE")); override != "" {
		return override
	}
	base, err := defaultCacheBase()
	if err != nil {
		return ""
	}
	return filepath.Join(base, "cursor-go-sdk", "bridge", "embedded-"+BundledSDKVersion)
}

func extractEmbeddedBridgeIfNeeded(root string) error {
	if len(embeddedBridgeTarball()) == 0 {
		return nil
	}
	if validBridgeRoot(root) {
		return nil
	}
	if err := os.MkdirAll(root, 0o755); err != nil {
		return err
	}
	return extractTarGz(embeddedBridgeTarball(), root)
}

func extractTarGz(data []byte, destDir string) error {
	gz, err := gzip.NewReader(newBytesReader(data))
	if err != nil {
		return err
	}
	defer gz.Close()
	tr := tar.NewReader(gz)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		target := filepath.Join(destDir, filepath.FromSlash(hdr.Name))
		if !strings.HasPrefix(filepath.Clean(target), filepath.Clean(destDir)+string(os.PathSeparator)) {
			return fmt.Errorf("invalid tar path %q", hdr.Name)
		}
		switch hdr.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, os.FileMode(hdr.Mode)); err != nil {
				return err
			}
		case tar.TypeReg, tar.TypeRegA:
			if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
				return err
			}
			out, err := os.OpenFile(target, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.FileMode(hdr.Mode))
			if err != nil {
				return err
			}
			if _, err := io.Copy(out, tr); err != nil {
				out.Close()
				return err
			}
			if err := out.Close(); err != nil {
				return err
			}
		default:
			continue
		}
	}
}

type bytesReader struct {
	b   []byte
	off int
}

func newBytesReader(b []byte) *bytesReader {
	return &bytesReader{b: b}
}

func (r *bytesReader) Read(p []byte) (int, error) {
	if r.off >= len(r.b) {
		return 0, io.EOF
	}
	n := copy(p, r.b[r.off:])
	r.off += n
	return n, nil
}

func defaultCacheBase() (string, error) {
	if xdg := strings.TrimSpace(os.Getenv("XDG_CACHE_HOME")); xdg != "" {
		return xdg, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	switch runtime.GOOS {
	case "darwin":
		return filepath.Join(home, "Library", "Caches"), nil
	default:
		return filepath.Join(home, ".cache"), nil
	}
}
