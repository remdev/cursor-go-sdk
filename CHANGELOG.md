# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed

- Minimum Go version for the root module lowered to 1.23 (`iter.Seq2` streaming API); TUI example requires Go 1.24.2 (`charmbracelet/bubbles`).

## [0.1.0] - 2026-06-15

### Added

- `cmd/setup` and `cursor.Setup()` — one-command install for bridge prerequisites (`go run github.com/remdev/cursor-go-sdk/cmd/setup@latest`)
- Opt-in e2e test suite (`e2e/`, `scripts/run-e2e.sh`) with `-tags=e2e` and `CURSOR_E2E=1`
- Bridge version gate (`>= 0.0.2`) and improved resolution (`CURSOR_SDK_BRIDGE_BIN`, `CURSOR_SDK_BRIDGE_ROOT`, PATH)
- `cursor.EnsureBridgeInstalled()` helper
- Unit tests for wire encoding and bridge lifecycle

### Changed

- README install path: prefer `go run .../cmd/setup@latest` over manual `npm link`
- `@cursor-go-sdk/cursor-sdk-bridge` npm `bin` points at compiled JS entrypoint (fixes global/Homebrew install)
- Bridge subprocess uses `exec.Command` with startup-only context (not tied to caller RPC context)

### Fixed

- `dropEmpty` no longer emits `cloud: {}` for local-only agents (prevented spurious idempotency keys)
- Idempotency keys only set for cloud agent create/send
- Global npm install of bridge 0.0.1: `MODULE_NOT_FOUND` when shell wrapper resolved wrong `dist/` path
- `examples/basic`: explicit API key check, `CURSOR_MODEL` env support

## [0.0.1] - 2026-06-13

### Added

- Initial public Go module `github.com/remdev/cursor-go-sdk/cursor`
- Connect RPC client over `@cursor-go-sdk/cursor-sdk-bridge`
- API parity with TypeScript `@cursor/sdk` and Python `cursor-sdk` (see `references/go-api-map.md`)
- Examples: `quickstart`, `coding-agent-cli`, `coding-agent-tui`, `basic`
- CI workflow (Go tests, bridge build, example builds)
- DISCLAIMER, CONTRIBUTING, SECURITY policy

[Unreleased]: https://github.com/remdev/cursor-go-sdk/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/remdev/cursor-go-sdk/releases/tag/v0.1.0
[0.0.1]: https://github.com/remdev/cursor-go-sdk/commit/51668d3
