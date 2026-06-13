# Contributing

Thanks for your interest in `cursor-go-sdk`. This is a community project; see [DISCLAIMER.md](DISCLAIMER.md) for the relationship to Cursor and Anysphere.

## Prerequisites

- Go 1.26+ (see `go.mod`)
- Node.js >= 18 and npm (for the bridge runtime)
- A [Cursor API key](https://cursor.com/docs) for manual live testing

## Setup

```bash
git clone https://github.com/remdev/cursor-go-sdk.git
cd cursor-go-sdk

cd bridge && npm install && cd ..

export CURSOR_API_KEY="cursor_..."   # optional, for live examples only
```

## Development

```bash
# Unit tests (root module)
go test ./...

# TUI example module
cd examples/coding-agent-tui && go test ./... && cd ../..

# Vet
go vet ./...
go vet ./examples/...

# Run an example
go run ./examples/quickstart
```

CI runs the same checks without `CURSOR_API_KEY`. Tests that need a installed bridge skip automatically when `bridge/node_modules` is missing locally before `npm install`.

## Pull requests

1. Fork the repository and create a feature branch.
2. Keep changes focused; match existing style in `cursor/` and examples.
3. Run `go test ./...` and `go vet ./...` before opening a PR.
4. Update [README.md](README.md) or [references/](references/) if behavior or setup changes.
5. Do not commit `bridge/node_modules/`, secrets, or `.artifacts/` contents.

## API parity

When adding or changing public API:

- Check [references/go-api-map.md](references/go-api-map.md) and the official TypeScript SDK docs.
- Prefer extending existing types (`ToWire()`, typed errors) over new abstractions.

## Questions

Open a [GitHub Discussion](https://github.com/remdev/cursor-go-sdk/discussions) for design questions, or an [issue](https://github.com/remdev/cursor-go-sdk/issues) for bugs and feature requests.

Do not use this repository for Cursor account, billing, or product support — contact Cursor directly.
