# AGENTS.md — cursor-go-sdk

Go SDK for Cursor agents. Parity target: TypeScript `@cursor/sdk` and Python `cursor-sdk`, idiomatic Go API.

**Unofficial community project** — see [DISCLAIMER.md](DISCLAIMER.md). Not affiliated with Cursor or Anysphere.

## Before coding

1. Read [references/README.md](references/README.md) — index of local docs and upstream links.
2. Bridge must be installed: `cd bridge && npm install` (see [references/bridge.md](references/bridge.md)).
3. Set `CURSOR_API_KEY` for live runs.
4. Run `go test ./...` after changes.

## Layout

| Path | Purpose |
|------|---------|
| `cursor/` | Public Go API (`CreateAgent`, `Prompt`, `Run`, …) |
| `internal/connect/` | Connect JSON RPC client |
| `internal/bridge/` | Bridge subprocess launcher + path resolution |
| `bridge/` | Node adapter: Connect server + `package.json` npm deps (`@cursor/sdk`) |
| `examples/` | Ports of [cursor/cookbook](https://github.com/cursor/cookbook) SDK examples |
| `references/` | Curated docs for agents (TS SDK, cookbook mapping, bridge) |
| `.artifacts/` | Local dev junk (venv, npm tarballs) — gitignored |

## Bridge (`bridge/`) — `@cursor/sdk` adapter

Go cannot load the npm SDK directly. **`bridge/`** is a Node Connect server that calls `@cursor/sdk` in-process; Go talks to it over loopback.

```bash
cd bridge && npm install   # installs @cursor/sdk + Connect deps
export CURSOR_SDK_BRIDGE_ROOT="$(pwd)/bridge"
```

- **npm runtime:** `bridge/package.json` → `@cursor/sdk`, platform packages, Connect/protobuf.
- **Adapter glue:** `bridge/dist/` — maintained in **this repo** for Go (see [references/bridge.md](references/bridge.md)).
- **Go launcher:** `internal/bridge/` — subprocess, discovery, callbacks.

Update path: bump `@cursor/sdk` in `package.json` → `npm install` → adjust `bridge/dist/` if API changed → `go test ./...`.

Or: `go run ./cmd/setup-bridge`

Requires **Node.js >= 18**. Override: `CURSOR_SDK_BRIDGE_BIN`, `CURSOR_SDK_BRIDGE_ROOT`, `CURSOR_SDK_NODE_BIN`.

## API quick map (TypeScript → Go)

| TypeScript | Go |
|------------|-----|
| `Agent.create(opts)` | `cursor.CreateAgent(ctx, opts)` |
| `Agent.prompt(msg, opts)` | `cursor.Prompt(ctx, msg, opts)` |
| `Agent.resume(id, opts)` | `cursor.ResumeAgent(ctx, id, opts)` |
| `agent.send(msg, opts)` | `agent.Send(ctx, msg, opts)` |
| `run.stream()` / `run.messages()` | `run.Stream(ctx)` / `run.Messages(ctx)` |
| `run.wait()` | `run.Wait(ctx)` |
| `run.iterText()` | `run.IterText(ctx)` |
| `await using agent` | `defer agent.Close(ctx)` |
| `Cursor.models.list()` | `cursor.Models.List(ctx, opts)` |

Full table: [references/go-api-map.md](references/go-api-map.md)

## Cookbook examples (upstream → Go)

| Cookbook | Go example |
|----------|------------|
| [sdk/quickstart](https://github.com/cursor/cookbook/tree/main/sdk/quickstart) | `examples/quickstart` |
| [sdk/coding-agent-cli](https://github.com/cursor/cookbook/tree/main/sdk/coding-agent-cli) (plain mode) | `examples/coding-agent-cli` |
| [sdk/coding-agent-cli](https://github.com/cursor/cookbook/tree/main/sdk/coding-agent-cli) (TUI) | `examples/coding-agent-tui` |
| One-shot prompt | `examples/basic` |

Other cookbook projects (app-builder, agent-kanban, dag-task-runner) are web/UI-heavy; port only when asked.

## Conventions

- Go 1.26+, module `github.com/remdev/cursor-go-sdk`
- Wire encoding: protobuf JSON field names via `ToWire()` / `wire.go`
- Errors: startup → typed `*cursor.AgentError`; run failure → `result.Status == "error"`
- Do not commit `bridge/node_modules/` or `.artifacts/` contents
- Prefer extending existing types over new abstractions

## Verification

```bash
go test ./...
go build -o /dev/null ./examples/...
# Live (needs API key + bridge):
go run ./examples/quickstart
```

## Upstream references

- [TypeScript SDK docs](https://cursor.com/docs/sdk/typescript)
- [Python SDK docs](https://cursor.com/docs/sdk/python)
- [Cookbook repo](https://github.com/cursor/cookbook)
- [npm @cursor/sdk](https://www.npmjs.com/package/@cursor/sdk)
