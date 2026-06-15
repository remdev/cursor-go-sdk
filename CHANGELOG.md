# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.0.2] - 2026-06-15

First tagged Go module release. Community-maintained client for the Cursor Agent SDK (`@cursor/sdk` parity).

### Added

- Go module `github.com/remdev/cursor-go-sdk/cursor` with API parity to TypeScript `@cursor/sdk` and Python `cursor-sdk` (see `references/go-api-map.md`)
- Connect RPC client over `@cursor-go-sdk/cursor-sdk-bridge`
- `cmd/setup` and `cursor.Setup()` — one-command bridge install (`go run github.com/remdev/cursor-go-sdk/cmd/setup@latest`)
- `cursor.EnsureBridgeInstalled()` helper
- Bridge version gate (`>= 0.0.2`) and resolution via `CURSOR_SDK_BRIDGE_BIN`, `CURSOR_SDK_BRIDGE_ROOT`, or `PATH`
- Opt-in e2e test suite (`e2e/`, `scripts/run-e2e.sh`) with `-tags=e2e` and `CURSOR_E2E=1`
- Examples ported from [cursor/cookbook](https://github.com/cursor/cookbook): `quickstart`, `coding-agent-cli`, `coding-agent-tui`, `basic`
- npm package `@cursor-go-sdk/cursor-sdk-bridge` — Connect adapter over `@cursor/sdk` (published **0.0.2**)
- Owned protobuf schema in `bridge/proto/` with buf codegen
- CI workflow (Go tests, bridge build, example builds)
- `CHANGELOG.md`, `docs/which-go-client.md`, `DISCLAIMER.md`, `CONTRIBUTING.md`, `SECURITY.md`

### Changed

- Bridge migrated to TypeScript with trimmed RPC surface for Go client needs
- README install path: `go run .../cmd/setup@latest` (replaces manual `npm link` for most users)
- npm bridge `bin` points at `dist/bin/cursor-sdk-bridge.js` (works with global/Homebrew install)
- Bridge subprocess uses `exec.Command` with startup-only context (not tied to caller RPC `context`)
- Minimum Go version: **1.23** for the root module; **1.24.2** for the TUI example (`charmbracelet/bubbles`)
- Disclaimer wording consolidated into `DISCLAIMER.md` (removed repetitive notices across docs)

### Fixed

- `dropEmpty` no longer emits `cloud: {}` for local-only agents (prevented spurious idempotency keys)
- Idempotency keys only set for cloud agent create/send
- Global npm install of bridge 0.0.1: `MODULE_NOT_FOUND` when shell wrapper resolved wrong `dist/` path
- Bridge subprocess killed on caller context cancel (`connection refused` after RPC/tests)
- Setup rejects bridge versions below minimum (`0.0.2`)
- `examples/basic`: explicit API key check, `CURSOR_MODEL` env support

[Unreleased]: https://github.com/remdev/cursor-go-sdk/compare/v0.0.2...HEAD
[0.0.2]: https://github.com/remdev/cursor-go-sdk/releases/tag/v0.0.2
