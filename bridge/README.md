# `@cursor-go-sdk/cursor-sdk-bridge`

The Go SDK cannot load `@cursor/sdk` (Node). **`@cursor-go-sdk/cursor-sdk-bridge`** is the prerequisite npm package:

```
cursor/ (Go)  ──Connect JSON──►  cursor-sdk-bridge  ──in-process──►  @cursor/sdk (npm)
```

Install the bridge **before** using the Go SDK. The Go client launches the **`cursor-sdk-bridge`** binary from `PATH`.

## Install

```bash
npm install -g @cursor-go-sdk/cursor-sdk-bridge
```

Or via Go (from any project):

```bash
go run github.com/remdev/cursor-go-sdk/cmd/setup@latest
```

From a clone of this repository (development):

```bash
go run ./cmd/setup --local
```

Requires **Node.js >= 18**.

Verify:

```bash
cursor-sdk-bridge --help
```

## Publish (maintainers)

See [PUBLISHING.md](PUBLISHING.md).

## Layout

| Path | Role |
|------|------|
| `package.json` | npm package `@cursor-go-sdk/cursor-sdk-bridge` |
| `proto/` | Owned protobuf wire schema (`sdk.v1`) |
| `gen/ts/` | Generated Connect + ES modules (`npm run generate`) |
| `src/` | TypeScript source (Connect RPC handlers → `@cursor/sdk`) |
| `dist/` | Compiled output (`npm run build`, gitignored) |
| `bin/cursor-sdk-bridge` | Dev/local shell launcher → `dist/bin/cursor-sdk-bridge.js` |
| npm `bin` | Points at `dist/bin/cursor-sdk-bridge.js` (works with `npm install -g`) |

## Wire protocol

Protobuf schema: [`proto/`](proto/) (`sdk.v1`). Regenerate TypeScript: `npm run generate` → `gen/ts/`.

Connect services: `SdkAgentService`, `SdkBridgeControlService`, `SdkCursorService`. Go client: `internal/connect/`.

## Go SDK discovery

1. `CURSOR_SDK_BRIDGE_BIN` — absolute path to the launcher
2. `CURSOR_SDK_BRIDGE_ROOT` — directory containing `bin/cursor-sdk-bridge`
3. `cursor-sdk-bridge` on `PATH` (from global install or `npm link`)

See [references/bridge.md](../references/bridge.md).
