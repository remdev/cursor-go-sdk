package bridge_test

import (
	"context"
	"testing"
	"time"

	"github.com/remdev/cursor-go-sdk/internal/bridge"
	conn "github.com/remdev/cursor-go-sdk/internal/connect"
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

func TestLaunchFromGlobalBinOnPATH(t *testing.T) {
	if _, err := bridge.ResolvePath(); err != nil {
		t.Skip(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	b, err := bridge.Launch(ctx, bridge.LaunchConfig{})
	if err != nil {
		t.Fatal(err)
	}
	defer b.Close()
	if b.Endpoint.URL == "" {
		t.Fatal("empty bridge URL")
	}
}

func TestLaunchSurvivesCallerContextCancel(t *testing.T) {
	if _, err := bridge.ResolvePath(); err != nil {
		t.Skip(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	b, err := bridge.Launch(ctx, bridge.LaunchConfig{})
	if err != nil {
		t.Fatal(err)
	}
	defer b.Close()
	cancel()

	transport := conn.NewTransport(b.Endpoint.URL, b.Endpoint.AuthToken)
	if _, err := transport.Unary(context.Background(), "sdk.v1.SdkBridgeControlService", "Ping", map[string]any{}); err != nil {
		t.Fatalf("bridge died after caller context was canceled: %v", err)
	}
}
