package bridge_test

import (
	"testing"

	"github.com/remdev/cursor-go-sdk/internal/bridge"
)

func TestParseDiscoveryLine(t *testing.T) {
	line := `cursor-sdk-bridge ready {"schemaVersion":1,"transport":"tcp","protocol":"connect","url":"http://127.0.0.1:1234","authToken":"secret"}`
	d, err := bridge.ParseDiscoveryLine(line)
	if err != nil {
		t.Fatal(err)
	}
	if d["schemaVersion"].(float64) != 1 {
		t.Fatalf("unexpected schema: %v", d["schemaVersion"])
	}
	if d["authToken"] != "secret" {
		t.Fatalf("unexpected token: %v", d["authToken"])
	}
}
