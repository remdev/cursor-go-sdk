package cursor

import (
	"context"
	"os"
	"sync"

	"github.com/remdev/cursor-go-sdk/internal/bridge"
)

var (
	defaultMu      sync.Mutex
	defaultClient  *Client
	defaultBridge  *bridge.Bridge
)

// DefaultClient returns the process-wide default client, launching a bridge if needed.
func DefaultClient(ctx context.Context) (*Client, error) {
	defaultMu.Lock()
	defer defaultMu.Unlock()
	if defaultClient != nil {
		return defaultClient, nil
	}
	if url := os.Getenv("CURSOR_SDK_BRIDGE_URL"); url != "" {
		token := os.Getenv("CURSOR_SDK_BRIDGE_TOKEN")
		if token == "" {
			token = os.Getenv("CURSOR_SDK_BRIDGE_AUTH_TOKEN")
		}
		if token != "" {
			c, err := Connect(url, token, WithAllowAPIKeyEnvFallback(false))
			if err != nil {
				return nil, err
			}
			defaultClient = c
			return c, nil
		}
	}
	c, err := LaunchBridge(ctx, WithAllowAPIKeyEnvFallback(true))
	if err != nil {
		return nil, err
	}
	defaultBridge = c.ownedBridge
	defaultClient = c
	return c, nil
}

// CloseDefaultClient closes the process-wide default client and bridge.
func CloseDefaultClient() error {
	defaultMu.Lock()
	c := defaultClient
	b := defaultBridge
	defaultClient = nil
	defaultBridge = nil
	defaultMu.Unlock()
	var err error
	if c != nil {
		err = c.Close()
	}
	if b != nil {
		if closeErr := b.Close(); closeErr != nil && err == nil {
			err = closeErr
		}
	}
	return err
}

// EnsureBridgeInstalled reports whether bridge/ has npm dependencies installed.
func EnsureBridgeInstalled(ctx context.Context) error {
	return bridge.EnsureInstalled(ctx)
}

func defaultOrLaunchClient(ctx context.Context, opts ...ClientOption) (*Client, error) {
	// If caller passed explicit connect/launch options, create a dedicated client.
	if len(opts) > 0 {
		cfg := clientConfig{}
		for _, opt := range opts {
			opt(&cfg)
		}
		if cfg.baseURL != "" && cfg.authToken != "" {
			return Connect(cfg.baseURL, cfg.authToken, opts...)
		}
		if cfg.endpoint != nil || cfg.ownedBridge != nil {
			return newClient(cfg)
		}
	}
	return DefaultClient(ctx)
}
