# TypeScript SDK — concepts for Go ports

Source: [cursor.com/docs/sdk/typescript](https://cursor.com/docs/sdk/typescript)

## Runtime

| Mode | TS | Go equivalent |
|------|-----|---------------|
| Local | `local: { cwd }` | `Local: &cursor.LocalAgentOptions{CWD: []string{...}}` |
| Cloud | `cloud: { repos, autoCreatePR }` | `Cloud: &cursor.CloudAgentOptions{...}` |

"Local" = agent loop on disk, **not** local LLM. Inference still uses Cursor hosted models.

## Agent lifecycle

```typescript
await using agent = await Agent.create({ ... });
const run = await agent.send("...");
for await (const event of run.stream()) { ... }
await run.wait();
```

Go:

```go
agent, err := cursor.CreateAgent(ctx, opts)
defer agent.Close(ctx)
run, err := agent.Send(ctx, "...", cursor.SendOptions{})
for msg, err := range run.Messages(ctx) { ... }
result, err := run.Wait(ctx)
```

## One-shot

`Agent.prompt(msg, opts)` — create, send, wait, dispose.

Go: `cursor.Prompt(ctx, msg, opts)`.

## Streaming messages

TS `run.stream()` yields `SDKMessage` with `type`: `assistant`, `tool`, `thinking`, `status`, etc.

Go: `run.Messages(ctx)` → `cursor.SDKMessage`. Use `cursor.AssistantText(msg)` for assistant text blocks.

## Errors

- Thrown before run starts → `CursorAgentError` / Go typed `*cursor.AgentError`
- Run started but failed → `result.status === "error"` / `result.Status == cursor.RunStatusError`

## Models

List: `Cursor.models.list({ apiKey })` → Go `cursor.Models.List(ctx, cursor.CursorRequestOptions{APIKey: ...})`.

Default model in cookbook quickstart: `composer-2` or env `CURSOR_MODEL`.

## What Go SDK does not port (Node-only)

- `SqliteLocalAgentStore`, in-process `@cursor/sdk` (TS runs without bridge)
- Web/TUI cookbook apps as-is (app-builder, kanban TUI)
