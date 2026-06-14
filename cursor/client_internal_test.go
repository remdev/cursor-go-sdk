package cursor

import "testing"

func TestBuildCreateAgentRequestLocalOmitsIdempotencyKey(t *testing.T) {
	opts := AgentOptions{
		Model:  "composer-2.5",
		APIKey: "key",
		Local:  LocalOptions("/tmp"),
	}
	wire := opts.ToWire()
	req := map[string]any{"options": stripEmptyMessage(wire)}
	if opts.Cloud != nil {
		req["idempotencyKey"] = newID()
	}
	if _, ok := req["idempotencyKey"]; ok {
		t.Fatalf("local create should not set idempotencyKey: %v", req)
	}
	if _, ok := wire["cloud"]; ok {
		t.Fatalf("local wire should omit cloud: %v", wire)
	}
}

func TestBuildCreateAgentRequestCloudSetsIdempotencyKey(t *testing.T) {
	opts := AgentOptions{
		Model:  "composer-2.5",
		APIKey: "key",
		Cloud:  &CloudAgentOptions{},
	}
	req := map[string]any{"options": stripEmptyMessage(opts.ToWire())}
	if opts.Cloud != nil {
		if req["idempotencyKey"] == nil {
			req["idempotencyKey"] = newID()
		}
	}
	if req["idempotencyKey"] == "" {
		t.Fatalf("cloud create should set idempotencyKey: %v", req)
	}
}
