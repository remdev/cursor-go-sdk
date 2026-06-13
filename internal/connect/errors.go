package connect

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// Error is the base error type for bridge failures.
type Error struct {
	Message       string
	Code          string
	Status        int
	Details       []map[string]any
	IsRetryable   bool
	ProtoErrorCode string
	RequestID     string
	Headers       map[string]string
	RetryAfter    string
}

func (e *Error) Error() string {
	if e.Code != "" {
		return e.Code + ": " + e.Message
	}
	return e.Message
}

// APIError is a backward-compatible alias.
type APIError = Error

func networkError(cause error) error {
	msg := fmt.Sprintf("bridge request failed: %T: %v", cause, cause)
	if isTimeout(cause) {
		return &Error{Message: msg, IsRetryable: true, Code: "timeout"}
	}
	return &Error{Message: msg, IsRetryable: true, Code: "network"}
}

func isTimeout(err error) bool {
	if err == nil {
		return false
	}
	s := err.Error()
	return strings.Contains(s, "timeout") || strings.Contains(s, "deadline exceeded")
}

func isRetryableError(err error) bool {
	var e *Error
	if AsError(err, &e) {
		return e.IsRetryable
	}
	return false
}

func isRetryableStatus(status int) bool {
	return status == 408 || status == 429 || status >= 500
}

func sleepBeforeRetry(attempt int) {
	delay := time.Duration(250*(1<<attempt)) * time.Millisecond
	if delay > 2*time.Second {
		delay = 2 * time.Second
	}
	time.Sleep(delay)
}

func parseHTTPError(status int, body []byte, headers http.Header) error {
	if len(body) == 0 {
		return &Error{
			Message: fmt.Sprintf("bridge request failed with HTTP %d", status),
			Status:  status,
			Headers: headerMap(headers),
		}
	}
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		text := string(body)
		if len(text) > 200 {
			text = text[:200] + "..."
		}
		return &Error{
			Message: fmt.Sprintf("bridge request failed with HTTP %d: %s", status, text),
			Status:  status,
			Headers: headerMap(headers),
		}
	}
	errObj, _ := payload["error"].(map[string]any)
	if errObj == nil {
		errObj = payload
	}
	return connectError(errObj, status, headers)
}

func connectError(err map[string]any, status int, headers http.Header) error {
	message := stringValue(err["message"], "bridge request failed")
	code := stringValue(err["code"], "")
	details := detailList(err["details"])
	parsed := errorFromConnectPayload(message, code, details, status)
	if headers != nil {
		if parsed.Headers == nil {
			parsed.Headers = headerMap(headers)
		}
		if ra := headers.Get("Retry-After"); ra != "" && parsed.RetryAfter == "" {
			parsed.RetryAfter = ra
		}
	}
	return parsed
}

func headerMap(h http.Header) map[string]string {
	out := make(map[string]string, len(h))
	for k, v := range h {
		if len(v) > 0 {
			out[k] = v[0]
		}
	}
	return out
}

func detailList(v any) []map[string]any {
	raw, ok := v.([]any)
	if !ok {
		return nil
	}
	out := make([]map[string]any, 0, len(raw))
	for _, item := range raw {
		if m, ok := item.(map[string]any); ok {
			out = append(out, m)
		}
	}
	return out
}

func stringValue(v any, fallback string) string {
	if v == nil {
		return fallback
	}
	if s, ok := v.(string); ok {
		return s
	}
	return fmt.Sprint(v)
}

// AsError reports whether err is a *Error.
func AsError(err error, target **Error) bool {
	if err == nil {
		return false
	}
	if e, ok := err.(*Error); ok {
		*target = e
		return true
	}
	return false
}

func errorFromConnectPayload(message, code string, details []map[string]any, status int) *Error {
	sdk := findSDKErrorDetails(details)
	if sdk != nil {
		if m := stringValue(sdk["message"], ""); m != "" {
			message = m
		}
		proto := stringValue(sdk["sdkErrorCode"], "")
		reqID := stringValue(sdk["requestId"], "")
		retryAfter := stringValue(sdk["retryAfter"], "")
		isRetryable := retryAfter != ""
		if isIntegrationError(proto, code) {
			return &Error{
				Message: message, Code: code, Status: status, Details: details,
				IsRetryable: isRetryable, ProtoErrorCode: proto, RequestID: reqID,
				RetryAfter: retryAfter,
			}
		}
		return typedError(proto, code, message, code, status, details, isRetryable, proto, reqID, retryAfter)
	}
	return typedError("", code, message, code, status, details, false, "", "", "")
}

func findSDKErrorDetails(details []map[string]any) map[string]any {
	var fallback map[string]any
	for _, detail := range details {
		for _, candidate := range detailCandidates(detail) {
			if _, ok := candidate["sdkErrorCode"]; ok {
				return candidate
			}
			if fallback == nil {
				if _, ok := candidate["requestId"]; ok {
					fallback = candidate
				}
			}
		}
	}
	return fallback
}

func detailCandidates(detail map[string]any) []map[string]any {
	out := []map[string]any{detail}
	if debug, ok := detail["debug"].(map[string]any); ok {
		out = append(out, debug)
	}
	if value, ok := detail["value"].(map[string]any); ok {
		out = append(out, value)
	}
	return out
}

func isIntegrationError(proto, connectCode string) bool {
	n := normalizeErrorCode(firstNonEmpty(proto, connectCode))
	switch n {
	case "integration_not_connected", "repository_access", "pr_resolution_failed":
		return true
	}
	return false
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}

func normalizeErrorCode(code string) string {
	return strings.TrimPrefix(strings.ToLower(code), "sdk_error_code_")
}

func typedError(proto, connectCode, message, code string, status int, details []map[string]any, retryable bool, protoCode, reqID, retryAfter string) *Error {
	n := normalizeErrorCode(firstNonEmpty(proto, connectCode))
	e := &Error{
		Message: message, Code: code, Status: status, Details: details,
		IsRetryable: retryable, ProtoErrorCode: protoCode, RequestID: reqID, RetryAfter: retryAfter,
	}
	switch n {
	case "unauthorized", "unauthenticated", "api_key_not_found":
		e.Code = "authentication"
	case "role_forbidden", "permission_denied", "forbidden":
		e.Code = "permission_denied"
	case "rate_limit_exceeded", "usage_limit_exceeded", "resource_exhausted", "rate_limited":
		e.Code = "rate_limit"
	case "invalid_argument", "validation_error", "bad_request":
		e.Code = "bad_request"
	case "invalid_model", "invalid_branch_name", "not_found":
		e.Code = "configuration"
	case "agent_busy":
		e.Code = "agent_busy"
	case "agent_not_found":
		e.Code = "agent_not_found"
	case "run_not_found":
		e.Code = "not_found"
	case "timeout":
		e.Code = "timeout"
		e.IsRetryable = true
	case "upstream_error", "internal_error", "internal", "internal_server_error":
		e.Code = "internal"
		e.IsRetryable = true
	case "unavailable", "deadline_exceeded", "data_loss":
		e.Code = "network"
		e.IsRetryable = true
	}
	return e
}

// Typed error predicates for callers.

func IsAuthentication(err error) bool {
	var e *Error
	return AsError(err, &e) && (e.Code == "authentication" || e.Code == "permission_denied")
}

func IsRateLimit(err error) bool {
	var e *Error
	return AsError(err, &e) && e.Code == "rate_limit"
}

func IsNotFound(err error) bool {
	var e *Error
	return AsError(err, &e) && (e.Code == "not_found" || e.Code == "agent_not_found")
}

func IsConfiguration(err error) bool {
	var e *Error
	return AsError(err, &e) && (e.Code == "configuration" || e.Code == "bad_request")
}
