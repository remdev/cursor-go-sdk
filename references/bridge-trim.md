# Bridge trim inventory

What the **Go-focused** `@cursor/sdk` adapter keeps vs what was removed from the generic Cursor bridge glue.

## Stack

```
Go client ──Connect JSON (bridge/proto)──► bridge/src (dist/) ──@cursor/sdk──► npm SDK
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

All agent RPCs remain in `sdk-service.ts` so the Go module can expose full API parity with `@cursor/sdk`.

## Removed (Go adapter only)

| Item | Reason |
|------|--------|
| `dist/index.js` | Library barrel; launcher does not import it |
| `process-error-survivors` | E2E test harness only; Go never enables it |
| `store-callback-config` | Host-owned store via loopback RPC — Go uses jsonl/sqlite on the wire instead |
| Host store callback in `bridge-local-agent-store` | `local.store.type: "custom"` not supported; Go uses jsonl or default sqlite |
| Store callback CLI/env in launcher | Not used by Go launcher |
| `@anysphere/proto` vendor + `bridge/vendor/` | Replaced by owned `bridge/proto/` + buf codegen |

## Kept (required for Go)

| Item | Reason |
|------|--------|
| `bridge-custom-tools.ts` | Go `ToolCallbackServer` for custom tools |
| `tool-callback-config.ts` | `--tool-callback-url` from `LaunchBridge` |
| `bridge-local-agent-store.ts` (slim) | Optional jsonl `local.store` |
| Full `sdk-service.ts` + `sdk-converters.ts` | Complete RPC mapping to `@cursor/sdk` |

## Not trimmed

- **Per-RPC handlers** — partial cuts would fork the adapter and break Go API parity on every `@cursor/sdk` bump.
- **`SetToolCallback` RPC** — optional; attach-to-running-bridge scenarios.
- **Cloud + local converter paths** — both used by Go SDK.

## Flags Go actually passes

From `internal/bridge/bridge.go` / `cursor.LaunchBridge`:

- `--workspace`, optional `--state-root`, `--host`, `--port`
- `--tool-callback-url`, `--tool-callback-auth-token`

Not passed: `--local-store`, `--max-concurrent-agents`, `--max-message-bytes` (available for advanced use).
