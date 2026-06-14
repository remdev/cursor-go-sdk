# Bridge runtime

`@cursor-go-sdk/cursor-sdk-bridge` is the **Node adapter** that lets Go call [`@cursor/sdk`](https://www.npmjs.com/package/@cursor/sdk) over Connect JSON on loopback. It is a **prerequisite** for the Go SDK, distributed as an npm package with a `cursor-sdk-bridge` binary.

```
Go (cursor/)  ──Connect JSON (proto/sdk.v1)──►  cursor-sdk-bridge  ──@cursor/sdk──►  agents
```

Protobuf schema is owned in [`bridge/proto/`](../bridge/proto/). TypeScript stubs: `npm run generate`.

## Install (prerequisite)

```bash
npm install -g @cursor-go-sdk/cursor-sdk-bridge@0.0.2
```

Requires **Node.js >= 18** and **npm**.

Development from this repo:

```bash
go run ./cmd/setup --local
```

Verify:

```bash
cursor-sdk-bridge --help
```

```go
if err := cursor.EnsureBridgeInstalled(ctx); err != nil { /* bridge not on PATH */ }
```

## What the adapter does

The Connect server in `bridge/src/` (compiled to `dist/`) translates RPC requests into `@cursor/sdk` calls:

- **Local agents:** tools, MCP, store, `api2.cursor.sh` — all inside the npm SDK
- **Cloud agents:** `api.cursor.com` via the npm SDK

This is **not** the `cursor-agent` CLI.

Custom tools: Go starts a loopback `ToolCallbackServer`; bridge forwards tool execution to Go via `--tool-callback-*`.

## Environment

| Variable | Purpose |
|----------|---------|
| `CURSOR_SDK_BRIDGE_BIN` | Override launcher path |
| `CURSOR_SDK_BRIDGE_ROOT` | Bridge package root (uses `dist/bin/cursor-sdk-bridge.js` when built) |
| `CURSOR_SDK_NODE_BIN` | Node binary for launcher |
| `CURSOR_SDK_BRIDGE_URL` | Connect to an existing bridge (skip local launch) |
| `CURSOR_SDK_BRIDGE_TOKEN` | Token for an existing bridge |
| `CURSOR_SDK_USE_REMOTE_BRIDGE=1` | Skip `PATH` lookup; use URL or explicit bin |

## Go discovery order

1. `CURSOR_SDK_BRIDGE_BIN`
2. `CURSOR_SDK_BRIDGE_ROOT`
3. `cursor-sdk-bridge` on `PATH`

## Maintenance

- **Bump runtime:** edit `bridge/package.json` → publish new npm version → test Go SDK.
- **Bump adapter:** edit `bridge/src/` when Connect proto or `@cursor/sdk` API changes; run `npm run generate && npm run build`.
- **Trim inventory:** [bridge-trim.md](bridge-trim.md).

Adapter source: [../bridge/README.md](../bridge/README.md).
