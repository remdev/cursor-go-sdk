package cursor

// McpServerFromMap normalizes a flat MCP server config map to wire format.
// Accepts both nested ({http: ...}, {stdio: ...}) and flat ({type, url, command}) shapes.
func McpServerFromMap(m map[string]any) map[string]any {
	if m == nil {
		return map[string]any{}
	}
	if _, ok := m["http"]; ok {
		return m
	}
	if _, ok := m["stdio"]; ok {
		return m
	}
	transport := stringField(m, "type")
	if transport == "http" || transport == "sse" || (transport == "" && m["url"] != nil) {
		t := transport
		if t == "" {
			t = "http"
		}
		auth := m["auth"]
		return map[string]any{
			"http": dropEmpty(map[string]any{
				"type":    enumValue(t, "HTTP_MCP_TRANSPORT_TYPE_"),
				"url":     stringField(m, "url"),
				"headers": m["headers"],
				"auth":    normalizeMcpAuth(auth),
			}, "headers"),
		}
	}
	if transport == "stdio" || m["command"] != nil {
		return map[string]any{
			"stdio": dropEmpty(map[string]any{
				"command": stringField(m, "command"),
				"args":    m["args"],
				"env":     m["env"],
				"cwd":     stringField(m, "cwd"),
			}, "env"),
		}
	}
	return m
}

func normalizeMcpAuth(v any) map[string]any {
	if v == nil {
		return map[string]any{}
	}
	if a, ok := v.(McpAuth); ok {
		return a.ToWire()
	}
	if m, ok := v.(map[string]any); ok {
		out := dropEmpty(map[string]any{
			"clientId":     firstPresent(m, "clientId", "CLIENT_ID"),
			"clientSecret": firstPresent(m, "clientSecret", "CLIENT_SECRET"),
			"scopes":       firstPresent(m, "scopes", "SCOPES"),
		})
		for k, val := range m {
			if k == "clientId" || k == "CLIENT_ID" || k == "clientSecret" || k == "CLIENT_SECRET" || k == "scopes" || k == "SCOPES" {
				continue
			}
			if _, exists := out[k]; !exists && val != nil {
				out[k] = val
			}
		}
		return out
	}
	return map[string]any{}
}

func firstPresent(m map[string]any, keys ...string) any {
	for _, k := range keys {
		if v, ok := m[k]; ok {
			return v
		}
	}
	return nil
}

// FlatMcpServer implements McpServer for map-based configs.
type FlatMcpServer map[string]any

func (f FlatMcpServer) ToWire() map[string]any {
	return McpServerFromMap(f)
}
