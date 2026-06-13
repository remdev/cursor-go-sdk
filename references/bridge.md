# Bridge runtime

The Go SDK talks to **cursor-sdk-bridge** (Connect RPC on loopback). The bridge glue lives in `bridge/`; runtime deps come from npm.

## Install (once per machine / clone)

```bash
cd bridge && npm install
```

Dependencies are declared in `bridge/package.json`:

- `@cursor/sdk@1.0.18` — agent runtime (same as TypeScript SDK)
- `@cursor/sdk-<platform>@1.0.18` — optional native helpers (rg, sandbox)

Requires **Node.js >= 18** and **npm**.

Helper:

```bash
go run ./cmd/setup-bridge
export CURSOR_SDK_BRIDGE_ROOT="$(go list -f '{{.Dir}}' -m github.com/remdev/cursor-go-sdk)/bridge"
```

## Environment

| Variable | Purpose |
|----------|---------|
| `CURSOR_SDK_BRIDGE_ROOT` | Directory with `bin/`, `dist/`, `node_modules/` |
| `CURSOR_SDK_BRIDGE_BIN` | Path to launcher script (skips root discovery) |
| `CURSOR_SDK_NODE_BIN` | Node binary for launcher |
| `CURSOR_SDK_USE_REMOTE_BRIDGE=1` | Skip bundled bridge lookup; use PATH only |

## Discovery order (Go)

1. `CURSOR_SDK_BRIDGE_BIN`
2. `CURSOR_SDK_BRIDGE_ROOT`
3. `bridge/` next to module source (GOMODCACHE or checkout)
4. `bridge/` in cwd or ancestors
5. `cursor-sdk-bridge` on PATH

## What bridge calls

Bridge is a thin Connect server over `@cursor/sdk` in-process:

- **Local**: agent loop on your machine → tools, MCP, `api2.cursor.sh`
- **Cloud**: `api.cursor.com` `/v1/agents`, …

Not the same as `cursor-agent` CLI.

## Dev artifacts

Local experiments (venv, downloaded npm tarballs, `fetch-bridge` for `-tags embedbridge`) go in `.artifacts/` (gitignored). Embed helper: `.artifacts/tools/fetch-bridge`.
