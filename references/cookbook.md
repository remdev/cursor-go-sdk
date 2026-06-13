# Cookbook → Go examples

Upstream: [github.com/cursor/cookbook](https://github.com/cursor/cookbook)

## Ported

### Quickstart — `examples/quickstart`

TS: [sdk/quickstart/src/index.ts](https://github.com/cursor/cookbook/blob/main/sdk/quickstart/src/index.ts)

- `Agent.create` with `local.cwd`, optional `name`
- `agent.send(prompt)`
- Stream assistant text from `run.stream()`
- `run.wait()`

Run:

```bash
export CURSOR_API_KEY=...
go run ./examples/quickstart
```

### Coding agent CLI (plain mode) — `examples/coding-agent-cli`

TS: [sdk/coding-agent-cli](https://github.com/cursor/cookbook/tree/main/sdk/coding-agent-cli) — `runPlainPrompt` path only (no TUI).

Flags: `--cwd`, `--model`, `--force`, positional prompt or stdin.

Run:

```bash
go run ./examples/coding-agent-cli -- "Explain the auth flow"
printf "Review changes" | go run ./examples/coding-agent-cli
```

### Coding agent CLI (TUI) — `examples/coding-agent-tui`

TS: [sdk/coding-agent-cli](https://github.com/cursor/cookbook/tree/main/sdk/coding-agent-cli) — interactive TUI path (`App.tsx`).

- Scrollable transcript, streaming assistant output
- Slash commands: `/help`, `/model`, `/local`, `/cloud`, `/reset`, `/exit`
- Ctrl+C cancels an in-flight run or exits
- Multi-turn session via `agent.Send`

Run:

```bash
export CURSOR_API_KEY=...
go run ./examples/coding-agent-tui
```

Uses [Bubble Tea](https://github.com/charmbracelet/bubbletea) (separate `go.mod` under the example).

### Basic one-shot — `examples/basic`

Minimal `cursor.Prompt` smoke test.

## Not ported (by design)

| Cookbook | Reason |
|----------|--------|
| app-builder | Next.js web UI |
| agent-kanban | React kanban + cloud UI |
| dag-task-runner | Canvas + skill; heavy TS/Node |

Port web/canvas examples only on explicit request; use cloud REST API for headless cloud orchestration if SDK surface is insufficient.

## Adding a new example

1. Read upstream TS source in cookbook.
2. Map API via [go-api-map.md](go-api-map.md).
3. Add `examples/<name>/main.go` + short comment with upstream URL.
4. List in this file and `AGENTS.md`.
