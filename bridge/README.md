# Bridge — `@cursor/sdk` adapter for Go

The Go SDK cannot load `@cursor/sdk` (Node). **`bridge/`** is this project's Node adapter:

```
cursor/ (Go)  ──Connect JSON──►  bridge/dist  ──in-process──►  @cursor/sdk (npm)
```

- **Runtime dependency:** [`@cursor/sdk`](https://www.npmjs.com/package/@cursor/sdk) and platform packages — installed via `npm install` in this directory.
- **Adapter code:** `bridge/dist/` + `bin/cursor-sdk-bridge` — Connect RPC server that wraps the npm SDK. Maintained **in this repository** for Go; not fetched from Python or any other SDK package.

Trimmed for Go-only use — see [references/bridge-trim.md](../references/bridge-trim.md).

## Layout

| Path | Role |
|------|------|
| `package.json` | Pins `@cursor/sdk`, Connect, protobuf |
| `bin/cursor-sdk-bridge` | Shell launcher (`CURSOR_SDK_NODE_BIN`) → `dist/bin/cursor-sdk-bridge.js` |
| `dist/bin/cursor-sdk-bridge.js` | CLI: callbacks, discovery line on stderr |
| `dist/server.js` | HTTP + Connect router |
| `dist/sdk-service.js` | RPC handlers → `Agent`, `Cursor` from `@cursor/sdk` |
| `dist/sdk-converters.js` | Protobuf ↔ SDK types |
| `dist/bridge-custom-tools.js` | Custom tools → Go `ToolCallbackServer` |
| `dist/bridge-local-agent-store.js` | `jsonl` / default sqlite `local.store` |
| `dist/tool-callback-config.js` | Tool callback URL/token from CLI |
| `dist/auth.js`, `registry.js`, `sdk-error-interceptor.js`, `constants.js` | Auth, registry, errors |

## Install

```bash
cd bridge && npm install
```

Go discovers `bridge/` via `CURSOR_SDK_BRIDGE_ROOT`, module path, or cwd. See [references/bridge.md](../references/bridge.md).

## Updating

1. Bump `@cursor/sdk` (and `@cursor/sdk-*` platform pins) in `package.json`.
2. `npm install`
3. If the npm SDK API or proto shapes changed, update `bridge/dist/` glue accordingly.
4. `go test ./...`

There is **no** separate npm package for this adapter today — it lives here. Do not sync from `cursor-sdk` (Python); npm is the only upstream runtime.

## Wire protocol

Same Connect services as Cursor's official SDK bridge (`SdkAgentService`, `SdkBridgeControlService`, `SdkCursorService`). Go client in `internal/connect/` speaks that protocol.
