package cursor

import (
	"context"

	"github.com/remdev/cursor-go-sdk/internal/bridge"
)

// SetupOptions configures prerequisite installation for the Go SDK.
type SetupOptions struct {
	Local     bool
	BridgeDir string
	Version   string
}

// Setup installs runtime prerequisites (cursor-sdk-bridge via npm).
// If the bridge is already on PATH, Setup is a no-op.
func Setup(ctx context.Context, opts SetupOptions) error {
	return bridge.Setup(ctx, bridge.SetupOptions{
		Local:     opts.Local,
		BridgeDir: opts.BridgeDir,
		Version:   opts.Version,
	})
}
