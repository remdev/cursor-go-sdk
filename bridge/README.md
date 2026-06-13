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

From a clone of this repository (development):

```bash
cd bridge && npm ci && npm link
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
| `vendor/anysphere-proto/` | Vendored `@anysphere/proto` wire definitions |
| `bin/cursor-sdk-bridge` | Shell launcher → `dist/bin/cursor-sdk-bridge.js` |
| `dist/` | Connect RPC server |

## Go SDK discovery

1. `CURSOR_SDK_BRIDGE_BIN` — absolute path to the launcher
2. `CURSOR_SDK_BRIDGE_ROOT` — directory containing `bin/cursor-sdk-bridge`
3. `cursor-sdk-bridge` on `PATH` (from global install or `npm link`)

See [references/bridge.md](../references/bridge.md).
