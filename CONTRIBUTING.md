# Contributing

Thanks for your interest in `cursor-go-sdk`. See [DISCLAIMER.md](DISCLAIMER.md) for how this project relates to Cursor and upstream SDKs.

## Prerequisites

- Go 1.23+ (see `go.mod`; TUI example in `examples/coding-agent-tui` requires Go 1.24.2+)
- Node.js >= 18 and npm (for the bridge runtime)
- A [Cursor API key](https://cursor.com/docs) for manual live testing

## Setup

```bash
git clone https://github.com/remdev/cursor-go-sdk.git
cd cursor-go-sdk

go run ./cmd/setup --local

# After editing bridge/proto/*.proto:
# cd bridge && npm run generate && npm run build

export CURSOR_API_KEY="cursor_..."   # optional, for live examples only
```

## Development

```bash
# Unit tests (root module)
go test ./...

# Local e2e (real API key; not run in CI)
export CURSOR_E2E=1
export CURSOR_API_KEY="cursor_..."
./scripts/run-e2e.sh
# or: go test -tags=e2e -count=1 -timeout=5m -v ./e2e/...

# TUI example module
cd examples/coding-agent-tui && go test ./... && cd ../..

# Vet
go vet ./...
go vet ./examples/...

# Run an example
go run ./examples/quickstart
```

CI runs the same checks without `CURSOR_API_KEY`. Bridge tests require `cursor-sdk-bridge` on `PATH` (CI runs `npm ci && npm run build` in `bridge/`). E2e tests live in [`e2e/`](e2e/) and require `CURSOR_E2E=1`, `-tags=e2e`, and a valid API key.

## Pull requests

1. Fork the repository and create a feature branch.
2. Keep changes focused; match existing style in `cursor/` and examples.
3. Run `go test ./...` and `go vet ./...` before opening a PR.
4. Update [README.md](README.md) or [references/](references/) if behavior or setup changes.
5. Do not commit `bridge/node_modules/`, `bridge/dist/`, secrets, or `.artifacts/` contents.

## API parity

When adding or changing public API:

- Check [references/go-api-map.md](references/go-api-map.md) and the official TypeScript SDK docs.
- Prefer extending existing types (`ToWire()`, typed errors) over new abstractions.

## Questions

Open a [GitHub Discussion](https://github.com/remdev/cursor-go-sdk/discussions) for design questions, or an [issue](https://github.com/remdev/cursor-go-sdk/issues) for bugs and feature requests.

Do not use this repository for Cursor account, billing, or product support — contact Cursor directly.
