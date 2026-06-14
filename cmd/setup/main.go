// Command setup installs cursor-go-sdk runtime prerequisites (Node bridge via npm).
//
// Usage:
//
//	go run github.com/remdev/cursor-go-sdk/cmd/setup@latest
//	go run ./cmd/setup --local
package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/remdev/cursor-go-sdk/cursor"
	"github.com/remdev/cursor-go-sdk/internal/bridge"
)

func main() {
	local := flag.Bool("local", false, "install from ./bridge in a repository clone (npm ci && npm run build && npm link)")
	bridgeDir := flag.String("bridge-dir", "", "path to bridge/ (implies --local when set)")
	version := flag.String(
		"version",
		"",
		fmt.Sprintf("npm version for global install (default: pinned release; must be >= %s)", bridge.BridgeNpmVersion),
	)
	flag.Parse()

	ctx := context.Background()
	opts := cursor.SetupOptions{
		Local:     *local || *bridgeDir != "",
		BridgeDir: *bridgeDir,
		Version:   *version,
	}
	if err := cursor.Setup(ctx, opts); err != nil {
		fmt.Fprintf(os.Stderr, "setup failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("cursor-go-sdk prerequisites ready (cursor-sdk-bridge on PATH)")
}
