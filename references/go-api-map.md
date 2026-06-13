# Go API map

TypeScript / Python → Go. See also root `README.md`.

## Client

| TS / Python | Go |
|-------------|-----|
| `Client.launchBridge(...)` | `cursor.LaunchBridge(ctx, opts...)` |
| `Client.connect(url, token)` | `cursor.Connect(url, token, opts...)` |
| default client | `cursor.DefaultClient(ctx)` |

## Agent

| TS / Python | Go |
|-------------|-----|
| `Agent.create(opts)` | `cursor.CreateAgent(ctx, opts)` |
| `Agent.prompt(msg, opts)` | `cursor.Prompt(ctx, msg, opts)` |
| `Agent.resume(id, opts)` | `cursor.ResumeAgent(ctx, id, opts)` |
| `agent.send(msg, opts?)` | `agent.Send(ctx, msg, sendOpts)` |
| `agent.close()` / dispose | `agent.Close(ctx)` |
| `agent.reload()` | `agent.Reload(ctx)` |

## Run

| TS / Python | Go |
|-------------|-----|
| `run.stream()` / `run.messages()` | `run.Stream(ctx)` / `run.Messages(ctx)` |
| `run.events()` | `run.Events(ctx)` |
| `run.wait()` | `run.Wait(ctx)` |
| `run.text()` | `run.Text(ctx)` |
| `run.iterText()` | `run.IterText(ctx)` |
| `run.cancel()` | `run.Cancel(ctx)` |
| `run.conversation()` | `run.Conversation(ctx)` |
| `run.requestId` | `run.RequestID` |
| `run.onDidChangeStatus(fn)` | `run.OnDidChangeStatus(fn)` |
| `send({ onDelta, onStep })` | `SendOptions{OnDelta, OnStep}` |

## Static

| TS / Python | Go |
|-------------|-----|
| `Agent.get(id)` | `cursor.GetAgent(ctx, id, opts)` |
| `Agent.list(opts)` | `cursor.ListAgents(ctx, opts)` |
| `Agent.getRun(id, opts)` | `cursor.GetRun(ctx, id, opts)` |
| `Agent.listRuns(agentId, opts)` | `cursor.ListRuns(ctx, agentId, opts)` |
| `Agent.archive / unarchive / delete` | `cursor.ArchiveAgent`, … |

## Cursor namespace

| TS / Python | Go |
|-------------|-----|
| `Cursor.models.list()` | `cursor.Models.List(ctx, opts)` |
| `Cursor.repositories.list()` | `cursor.Repositories.List(ctx, opts)` |
| `Cursor.me()` | `cursor.Cursor.Me(ctx, opts)` |
| `Cursor.configure(opts)` | `cursor.Configure(opts)` |

## Options naming

Go uses exported struct fields (`APIKey`, `Local`, `Cloud`); wire JSON uses camelCase via `ToWire()`.

## Bridge

| Python | Go |
|--------|-----|
| `Client.launch_bridge(...)` | `cursor.LaunchBridge(ctx, ...)` |
| `resolve_bridge_path()` | `bridge` package resolves via env + `bridge/` dir |
| Bridge + `@cursor/sdk` | `cd bridge && npm install` |
