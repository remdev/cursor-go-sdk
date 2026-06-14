# Go client for Cursor agents

[![CI](https://github.com/remdev/cursor-go-sdk/actions/workflows/ci.yml/badge.svg)](https://github.com/remdev/cursor-go-sdk/actions/workflows/ci.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/remdev/cursor-go-sdk/cursor.svg)](https://pkg.go.dev/github.com/remdev/cursor-go-sdk/cursor)
[![npm bridge](https://img.shields.io/npm/v/@cursor-go-sdk/cursor-sdk-bridge)](https://www.npmjs.com/package/@cursor-go-sdk/cursor-sdk-bridge)
[![License: MIT](https://img.shields.io/github/license/remdev/cursor-go-sdk)](LICENSE)
[![Go Version](https://img.shields.io/github/go-mod/go-version/remdev/cursor-go-sdk)](go.mod)

Go client for [Cursor agents](https://cursor.com/docs/sdk/typescript). API parity with TypeScript `@cursor/sdk` and Python `cursor-sdk`.

Local agents run through **`cursor-sdk-bridge`** — a Node adapter over [`@cursor/sdk`](https://www.npmjs.com/package/@cursor/sdk). Install the bridge once; the Go SDK launches it automatically.

**Requirements:** Go 1.23+ (TUI example: Go 1.24.2+), Node.js >= 18, npm, [`@cursor-go-sdk/cursor-sdk-bridge`](bridge/) `>= 0.0.2` on `PATH`.

## How is this different?

| | cursor-go-sdk | REST-only Go SDKs | Official `@cursor/sdk` |
|--|---------------|-------------------|------------------------|
| API | Agent SDK parity | Cloud Agents REST | Agent SDK (official) |
| Local agents | Yes (via bridge) | No | Yes |
| Cloud agents | Yes | Yes | Yes |
| Language | Go | Go | TypeScript |

Not a REST wrapper for the Cloud Agents API — this mirrors `Agent.create`, `agent.send`, `run.stream`, and related SDK surface. See [docs/which-go-client.md](docs/which-go-client.md) for a fuller comparison.

## Install

```bash
go run github.com/remdev/cursor-go-sdk/cmd/setup@latest
go get github.com/remdev/cursor-go-sdk/cursor
```

This installs `@cursor-go-sdk/cursor-sdk-bridge` via npm (requires Node.js >= 18).

Manual alternative:

```bash
npm install -g @cursor-go-sdk/cursor-sdk-bridge
go get github.com/remdev/cursor-go-sdk/cursor
```

Check readiness:

```go
if err := cursor.EnsureBridgeInstalled(ctx); err != nil {
    // cursor-sdk-bridge not on PATH
}
```

Development from a clone:

```bash
go run ./cmd/setup --local
```

## Authentication

```bash
export CURSOR_API_KEY="cursor_..."
```

## Quick start

```go
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/remdev/cursor-go-sdk/cursor"
)

func main() {
	ctx := context.Background()

	agent, err := cursor.CreateAgent(ctx, cursor.AgentOptions{
		Model:  "composer-2.5",
		APIKey: os.Getenv("CURSOR_API_KEY"),
		Local:  &cursor.LocalAgentOptions{CWD: []string{"."}},
	})
	if err != nil {
		panic(err)
	}
	defer agent.Close(ctx)

	run, err := agent.Send(ctx, "Summarize this repository", cursor.SendOptions{})
	if err != nil {
		panic(err)
	}

	for msg, err := range run.Messages(ctx) {
		if err != nil {
			panic(err)
		}
		if msg.Type == "assistant" {
			fmt.Print(cursor.AssistantText(msg))
		}
	}

	result, err := run.Wait(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println("\nstatus:", result.Status)
}
```

One-shot helper:

```go
result, err := cursor.Prompt(ctx, "Explain main.go", cursor.AgentOptions{
	Model:  "composer-2.5",
	APIKey: os.Getenv("CURSOR_API_KEY"),
	Local:  &cursor.LocalAgentOptions{CWD: []string{"."}},
})
```

## Examples

Ports of [cursor/cookbook](https://github.com/cursor/cookbook) SDK samples:

| Example | Description |
|---------|-------------|
| [`examples/quickstart`](examples/quickstart) | Create agent, send, stream, wait |
| [`examples/coding-agent-cli`](examples/coding-agent-cli) | Non-interactive CLI (flags, stdin) |
| [`examples/coding-agent-tui`](examples/coding-agent-tui) | Interactive terminal UI |
| [`examples/basic`](examples/basic) | Minimal smoke test |

```bash
export CURSOR_API_KEY=...
go run ./examples/quickstart
go run ./examples/coding-agent-cli -- "Explain the auth flow"
go run ./examples/coding-agent-tui
```

### Local e2e tests

Opt-in integration tests against the real API (not run in CI):

```bash
export CURSOR_E2E=1
export CURSOR_API_KEY=cursor_...
./scripts/run-e2e.sh
```

Optional: `CURSOR_E2E_MODEL`, `CURSOR_E2E_WORKSPACE`, `CURSOR_E2E_TIMEOUT` (default `5m` in the script, `3m` per test).

## Configuration

| Variable | Purpose |
|----------|---------|
| `CURSOR_API_KEY` | API key |
| `CURSOR_E2E` | Set to `1` to enable local e2e tests in `e2e/` |
| `CURSOR_E2E_MODEL` | Model for e2e (default `auto`, falls back to `CURSOR_MODEL`) |
| `CURSOR_SDK_BRIDGE_BIN` | Override bridge launcher binary |
| `CURSOR_SDK_BRIDGE_ROOT` | Bridge package root (prefers `dist/bin/cursor-sdk-bridge.js`, else `bin/cursor-sdk-bridge`) |
| `CURSOR_SDK_NODE_BIN` | Override Node.js binary |
| `CURSOR_SDK_BRIDGE_URL` | Connect to an existing bridge |
| `CURSOR_SDK_BRIDGE_TOKEN` | Token for an existing bridge |
| `CURSOR_SDK_USE_REMOTE_BRIDGE` | Skip local bridge discovery on PATH |

Connect to a bridge that is already running:

```go
client, err := cursor.Connect(os.Getenv("CURSOR_SDK_BRIDGE_URL"), os.Getenv("CURSOR_SDK_BRIDGE_TOKEN"))
```

## Error handling

Startup errors are typed (`AuthenticationError`, `RateLimitError`, `AgentBusyError`, …) as `*cursor.AgentError`. A started run may finish with `result.Status == cursor.RunStatusError`.

```go
result, err := run.Wait(ctx)
if err != nil {
    var ae *cursor.AgentError
    if errors.As(err, &ae) {
        // ae.IsRetryable, ae.RetryAfter, ae.RequestID
    }
}
```

## Documentation

- [Changelog](CHANGELOG.md)
- [Which Go client for Cursor?](docs/which-go-client.md)
- [Disclaimer](DISCLAIMER.md)
- [Contributing](CONTRIBUTING.md)
- [Security policy](SECURITY.md)
- [API mapping (TS/Python → Go)](references/go-api-map.md)
- [Bridge npm package](references/bridge.md)
- [AGENTS.md](AGENTS.md) — contributor and agent notes

## License

[MIT](LICENSE)
