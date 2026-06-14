package cursor

import (
	"encoding/json"
	"os"
	"strings"
)

func dropEmpty(m map[string]any, preserveEmptyStringMaps ...string) map[string]any {
	preserve := map[string]struct{}{}
	for _, k := range preserveEmptyStringMaps {
		preserve[k] = struct{}{}
	}
	out := make(map[string]any, len(m))
	for k, v := range m {
		if v == nil {
			continue
		}
		switch t := v.(type) {
		case string:
			if t == "" {
				continue
			}
		case []any:
			if len(t) == 0 {
				continue
			}
		case map[string]any:
			if len(t) == 0 {
				if _, ok := preserve[k]; ok {
					out[k] = t
				}
				continue
			}
			nested := dropEmpty(t, preserveEmptyStringMaps...)
			if len(nested) == 0 {
				continue
			}
			out[k] = nested
			continue
		}
		out[k] = v
	}
	return out
}

func enumValue(value, prefix string) string {
	if value == "" {
		return ""
	}
	if strings.HasPrefix(value, prefix) {
		return value
	}
	return prefix + strings.ToUpper(value)
}

func stringField(m map[string]any, keys ...string) string {
	for _, k := range keys {
		if v, ok := m[k]; ok && v != nil {
			return fmtString(v)
		}
	}
	return ""
}

func intField(m map[string]any, key string) int {
	v, ok := m[key]
	if !ok || v == nil {
		return 0
	}
	switch n := v.(type) {
	case float64:
		return int(n)
	case int:
		return n
	case json.Number:
		i, _ := n.Int64()
		return int(i)
	default:
		return 0
	}
}

func boolField(m map[string]any, key string) bool {
	v, ok := m[key]
	if !ok || v == nil {
		return false
	}
	b, _ := v.(bool)
	return b
}

func fmtString(v any) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func normalizeRunStatus(value string) RunStatus {
	n := strings.ToLower(strings.TrimPrefix(value, "RUN_LIFECYCLE_STATUS_"))
	switch n {
	case "running", "finished", "error", "cancelled", "expired":
		return RunStatus(n)
	case "creating":
		return RunStatusRunning
	case "unspecified":
		return RunStatusError
	default:
		return RunStatusError
	}
}

func normalizeAgentStatus(value string) string {
	n := strings.ToLower(strings.TrimPrefix(value, "AGENT_INFO_STATUS_"))
	switch n {
	case "running", "finished", "error":
		return n
	default:
		return ""
	}
}

func apiKeyOrEnv(explicit string, allowEnv bool) string {
	if explicit != "" {
		return explicit
	}
	if allowEnv {
		return os.Getenv("CURSOR_API_KEY")
	}
	return ""
}

func deepMerge(base, overrides map[string]any) map[string]any {
	merged := make(map[string]any, len(base)+len(overrides))
	for k, v := range base {
		merged[k] = v
	}
	for k, v := range overrides {
		if k == "model" {
			merged[k] = v
			continue
		}
		if existing, ok := merged[k].(map[string]any); ok {
			if overrideMap, ok := v.(map[string]any); ok {
				merged[k] = deepMerge(existing, overrideMap)
				continue
			}
		}
		merged[k] = v
	}
	return merged
}

func stripEmptyMessage(m map[string]any) map[string]any {
	return dropEmpty(m, "envVars", "headers")
}

func hasCloudOptions(opts map[string]any) bool {
	if opts == nil {
		return false
	}
	cloud, ok := opts["cloud"].(map[string]any)
	return ok && len(cloud) > 0
}
