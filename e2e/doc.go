// Package e2e holds opt-in local integration tests against a real Cursor API key.
//
// Run manually on a developer machine:
//
//	export CURSOR_E2E=1
//	export CURSOR_API_KEY=cursor_...
//	go test -tags=e2e -count=1 -timeout=5m -v ./e2e/...
//
// Or: ./scripts/run-e2e.sh
//
// Test sources use //go:build e2e and are excluded from default go test ./...
package e2e
