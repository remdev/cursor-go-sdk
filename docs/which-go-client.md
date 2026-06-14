# Which Go client for Cursor?

This repo is one of several community Go options. They target different APIs and use cases.

| | **cursor-go-sdk** (this repo) | REST-only Go clients | **Official `@cursor/sdk`** |
|--|-------------------------------|----------------------|----------------------------|
| **Language** | Go | Go | TypeScript |
| **API surface** | Agent SDK parity (`Agent.create`, `agent.send`, `run.stream`, …) | [Cloud Agents REST API](https://cursor.com/docs/cloud-agent/api/endpoints) | Agent SDK (official) |
| **Local agents** | Yes — via Node bridge over `@cursor/sdk` | No | Yes — native |
| **Cloud agents** | Yes — same Go types | Yes — HTTP | Yes |
| **Node.js required** | Yes (bridge runtime) | No | Yes |
| **Maintained by** | Community (unofficial) | Community (unofficial) | Cursor / Anysphere |

## Community REST clients (cloud-only)

Useful when you only need the Cloud Agents HTTP API and do not need local agents or the full Agent SDK harness:

- [PippaOS/cursor-go](https://github.com/PippaOS/cursor-go) — OpenAPI-generated client + wrapper
- [donvito/cursor-sdk-go](https://github.com/donvito/cursor-sdk-go) — Cloud Agents API v1

## When to use cursor-go-sdk

- Your service or CLI is written in Go and you want **the same Agent SDK interface** as TypeScript and Python.
- You need **local agents** (codebase indexing, MCP, skills) from Go.
- You are porting [cursor/cookbook](https://github.com/cursor/cookbook) SDK examples to Go.

## When to use the official TypeScript or Python SDK instead

- You can run Node or Python in your stack without a bridge.
- You want the officially supported SDK path with the fastest upstream updates.

## Disclaimer

All community projects listed here are **unofficial** and not affiliated with Cursor or Anysphere. See [DISCLAIMER.md](../DISCLAIMER.md).
