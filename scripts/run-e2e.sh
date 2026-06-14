#!/usr/bin/env bash
set -euo pipefail

root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$root"

if [[ "${CURSOR_E2E:-}" != "1" ]]; then
  echo "Set CURSOR_E2E=1 to run local end-to-end tests." >&2
  exit 1
fi

if [[ -z "${CURSOR_API_KEY:-}" ]]; then
  echo "Set CURSOR_API_KEY to a valid Cursor API key." >&2
  exit 1
fi

timeout="${CURSOR_E2E_TIMEOUT:-5m}"
exec go test -tags=e2e -parallel=1 -count=1 -timeout="$timeout" -v ./e2e/...
