# Bridge trim inventory

What the **Go-focused** `@cursor/sdk` adapter keeps vs what was removed from the generic Cursor bridge glue.

## Stack

```
Go client ──Connect JSON──► bridge/dist (this repo) ──in-process──► @cursor/sdk (npm)
```

## RPC surface used by Go SDK

### SdkBridgeControlService

| Method | Used by Go |
|--------|------------|
| Ping | yes |
| GetVersion | yes |
| Shutdown | yes |
| SetToolCallback | no (Go passes `--tool-callback-*` at launch) |

### SdkCursorService

| Method | Used by Go |
|--------|------------|
| Me | yes |
| ListModels | yes |
| ListRepositories | yes |

### SdkAgentService

| Method | Used by Go |
|--------|------------|
| CreateAgent, ResumeAgent, ReloadAgent, CloseAgent | yes |
| Send | yes |
| WaitLiveRun, GetRun, ListRuns, CancelRun, ObserveRun, GetRunConversation | yes |
| GetAgent, ListAgents | yes |
| ArchiveAgent, UnarchiveAgent, DeleteAgent | yes |
| ListAgentMessages, ListArtifacts, DownloadArtifact | yes |

All agent RPCs remain in `sdk-service.js` so the Go module can expose full API parity with `@cursor/sdk`.

## Removed (Go adapter only)

| Item | Reason |
|------|--------|
| `dist/**/*.d.ts`, `*.d.ts.map` | Not needed at Node runtime |
| `dist/index.js` | Library barrel; launcher does not import it |
| `dist/process-error-survivors.js` | E2E test harness only; Go never enables it |
| `dist/store-callback-config.js` | Host-owned store via loopback RPC — Go uses jsonl/sqlite on the wire instead |
| Host store callback in `bridge-local-agent-store.js` | `local.store.type: "custom"` not supported; Go uses jsonl or default sqlite |
| Store callback CLI/env in launcher | Not used by Go launcher |

## Kept (required for Go)

| Item | Reason |
|------|--------|
| `bridge-custom-tools.js` | Go `ToolCallbackServer` for custom tools |
| `tool-callback-config.js` | `--tool-callback-url` from `LaunchBridge` |
| `bridge-local-agent-store.js` (slim) | Optional jsonl `local.store` |
| Full `sdk-service.js` + `sdk-converters.js` | Complete RPC mapping to `@cursor/sdk` |

## Not trimmed

- **Per-RPC handlers** — partial cuts would fork the adapter and break Go API parity on every `@cursor/sdk` bump.
- **`SetToolCallback` RPC** — optional; attach-to-running-bridge scenarios.
- **Cloud + local converter paths** — both used by Go SDK.

## Flags Go actually passes

From `internal/bridge/bridge.go` / `cursor.LaunchBridge`:

- `--workspace`, optional `--state-root`, `--host`, `--port`
- `--tool-callback-url`, `--tool-callback-auth-token`

Not passed: `--local-store`, `--max-concurrent-agents`, `--max-message-bytes` (available for advanced use).
