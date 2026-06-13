# Bridge runtime

`bridge/` is the **Node adapter** that lets Go call [`@cursor/sdk`](https://www.npmjs.com/package/@cursor/sdk) over Connect JSON on loopback. It is part of **cursor-go-sdk**, not a separate product.

```
Go (cursor/)  ──Connect JSON──►  bridge/dist  ──require()──►  @cursor/sdk in node_modules
```

## Install

```bash
cd bridge && npm install
```

Requires **Node.js >= 18** and **npm**.

Pinned in `bridge/package.json`:

- `@cursor/sdk@…` — agent runtime (npm; this is the upstream SDK)
- `@cursor/sdk-<platform>@…` — native helpers (rg, sandbox)
- `@connectrpc/connect`, `@connectrpc/connect-node`, `@bufbuild/protobuf` — Connect stack (1.x, compatible with `@cursor/sdk`)

Helper:

```bash
go run ./cmd/setup-bridge
export CURSOR_SDK_BRIDGE_ROOT="$(go list -f '{{.Dir}}' -m github.com/remdev/cursor-go-sdk)/bridge"
```

Verify:

```go
if err := cursor.EnsureBridgeInstalled(ctx); err != nil { /* npm install missing */ }
```

## What the adapter does

The Connect server in `bridge/dist/` translates RPC requests into `@cursor/sdk` calls:

- **Local agents:** tools, MCP, store, `api2.cursor.sh` — all inside the npm SDK
- **Cloud agents:** `api.cursor.com` via the npm SDK

This is **not** the `cursor-agent` CLI.

Custom tools: Go starts a loopback `ToolCallbackServer`; bridge forwards tool execution to Go via `--tool-callback-*`.

## Environment

| Variable | Purpose |
|----------|---------|
| `CURSOR_SDK_BRIDGE_ROOT` | Directory with `bin/`, `dist/`, `node_modules/` |
| `CURSOR_SDK_BRIDGE_BIN` | Override launcher script |
| `CURSOR_SDK_NODE_BIN` | Node binary for launcher |
| `CURSOR_SDK_USE_REMOTE_BRIDGE=1` | Skip bundled lookup; use PATH only |

## Go discovery order

1. `CURSOR_SDK_BRIDGE_BIN`
2. `CURSOR_SDK_BRIDGE_ROOT`
3. `bridge/` next to module source
4. `bridge/` in cwd or ancestors
5. `cursor-sdk-bridge` on PATH

## Maintenance

- **Bump runtime:** edit `bridge/package.json` → `npm install` → test.
- **Bump adapter:** edit `bridge/dist/` when Connect proto or `@cursor/sdk` API changes.
- **Trim inventory:** [bridge-trim.md](bridge-trim.md) — what was removed from the generic bridge for Go-only use.

Adapter source: [../bridge/README.md](../bridge/README.md).

## Optional embed

`go build -tags embedbridge` can ship a bridge tarball from `internal/bridge/prebuilt/` (see `.artifacts/tools/fetch-bridge`). Not required for normal use.
