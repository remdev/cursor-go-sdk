# References index

Curated material for agents working on **cursor-go-sdk**. Read these before implementing features or porting cookbook examples.

## Local

| File | Contents |
|------|----------|
| [bridge.md](bridge.md) | Bridge setup, npm deps, env vars |
| [go-api-map.md](go-api-map.md) | TypeScript / Python → Go API table |
| [cookbook.md](cookbook.md) | Upstream cookbook examples + Go ports |
| [typescript-sdk.md](typescript-sdk.md) | TS SDK concepts, runtime model, streaming |

## Upstream (always current)

- [TypeScript SDK](https://cursor.com/docs/sdk/typescript)
- [Python SDK](https://cursor.com/docs/sdk/python)
- [Cookbook repo](https://github.com/cursor/cookbook)
- [@cursor/sdk on npm](https://www.npmjs.com/package/@cursor/sdk) — version pinned in `bridge/package.json`
- [Cloud Agents REST API](https://cursor.com/docs/cloud-agent/api/endpoints) — alternative to SDK for cloud-only

## Cookbook source files (TypeScript)

Use these as the spec when porting to Go:

| Example | Upstream path |
|---------|---------------|
| Quickstart | [sdk/quickstart/src/index.ts](https://github.com/cursor/cookbook/blob/main/sdk/quickstart/src/index.ts) |
| Coding agent CLI | [sdk/coding-agent-cli/src/agent.ts](https://github.com/cursor/cookbook/blob/main/sdk/coding-agent-cli/src/agent.ts) |
| App builder | [sdk/app-builder](https://github.com/cursor/cookbook/tree/main/sdk/app-builder) |
| Agent kanban | [sdk/agent-kanban](https://github.com/cursor/cookbook/tree/main/sdk/agent-kanban) |
| DAG task runner | [sdk/dag-task-runner](https://github.com/cursor/cookbook/tree/main/sdk/dag-task-runner) |
