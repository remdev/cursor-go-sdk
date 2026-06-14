package cursor

import "testing"

func TestDropEmptyNilMap(t *testing.T) {
	var cloudWire map[string]any
	w := dropEmpty(map[string]any{"cloud": cloudWire})
	if _, ok := w["cloud"]; ok {
		t.Fatalf("nil cloud map should be omitted, got %v", w)
	}
}

func TestDropEmptyEmptyMap(t *testing.T) {
	w := dropEmpty(map[string]any{"cloud": map[string]any{}})
	if _, ok := w["cloud"]; ok {
		t.Fatalf("empty cloud map should be omitted, got %v", w)
	}
}

func TestDropEmptyPreservesNamedEmptyMaps(t *testing.T) {
	w := dropEmpty(map[string]any{
		"cloud": map[string]any{
			"envVars": map[string]any{},
		},
	}, "envVars")
	cloud, ok := w["cloud"].(map[string]any)
	if !ok {
		t.Fatalf("cloud=%v", w["cloud"])
	}
	if _, ok := cloud["envVars"]; !ok {
		t.Fatalf("envVars should be preserved, cloud=%v", cloud)
	}
}

func TestHasCloudOptions(t *testing.T) {
	if hasCloudOptions(map[string]any{"cloud": map[string]any{}}) {
		t.Fatal("empty cloud block is not cloud runtime")
	}
	if !hasCloudOptions(map[string]any{"cloud": map[string]any{"repos": []any{}}}) {
		t.Fatal("expected cloud runtime")
	}
}
