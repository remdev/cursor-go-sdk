package cursor

import (
	"errors"
	"fmt"

	conn "github.com/remdev/cursor-go-sdk/internal/connect"
)

// AgentError is the base error type for Cursor SDK failures.
type AgentError struct {
	Message        string
	Code           string
	Status         int
	Details        []map[string]any
	IsRetryable    bool
	ProtoErrorCode string
	RequestID      string
	Headers        map[string]string
	RetryAfter     string
}

func (e *AgentError) Error() string {
	if e.Code != "" {
		return e.Code + ": " + e.Message
	}
	return e.Message
}

// SDKError is a backward-compatible alias.
type SDKError = AgentError

// CursorAgentError is a backward-compatible alias matching TypeScript.
type CursorAgentError = AgentError

// Typed errors matching the TypeScript SDK hierarchy.

type AuthenticationError struct{ *AgentError }
type PermissionDeniedError struct{ *AgentError }
type RateLimitError struct{ *AgentError }
type ConfigurationError struct{ *AgentError }
type AgentBusyError struct{ *AgentError }
type BadRequestError struct{ *AgentError }
type NetworkError struct{ *AgentError }
type APITimeoutError struct{ *AgentError }
type InternalServerError struct{ *AgentError }
type NotFoundError struct{ *AgentError }
type AgentNotFoundError struct{ *AgentError }
type UnknownAgentError struct{ *AgentError }

type IntegrationNotConnectedError struct {
	*AgentError
	Provider string
	HelpURL  string
}

type UnsupportedRunOperationError struct {
	*AgentError
	Operation string
}

func wrapError(err error) error {
	if err == nil {
		return nil
	}
	var ce *conn.Error
	if !errors.As(err, &ce) {
		return err
	}
	base := &AgentError{
		Message:        ce.Message,
		Code:           ce.Code,
		Status:         ce.Status,
		Details:        ce.Details,
		IsRetryable:    ce.IsRetryable,
		ProtoErrorCode: ce.ProtoErrorCode,
		RequestID:      ce.RequestID,
		Headers:        ce.Headers,
		RetryAfter:     ce.RetryAfter,
	}
	return typedFromAgentError(base)
}

func typedFromAgentError(base *AgentError) error {
	switch base.Code {
	case "authentication", "permission_denied":
		if base.Code == "permission_denied" {
			return &PermissionDeniedError{AgentError: base}
		}
		return &AuthenticationError{AgentError: base}
	case "rate_limit":
		return &RateLimitError{AgentError: base}
	case "agent_busy":
		return &AgentBusyError{AgentError: base}
	case "agent_not_found":
		return &AgentNotFoundError{AgentError: base}
	case "not_found":
		return &NotFoundError{AgentError: base}
	case "bad_request":
		return &BadRequestError{AgentError: base}
	case "configuration":
		return &ConfigurationError{AgentError: base}
	case "timeout":
		return &APITimeoutError{AgentError: base}
	case "internal":
		return &InternalServerError{AgentError: base}
	case "network":
		return &NetworkError{AgentError: base}
	case "unsupported_run_operation":
		return &UnsupportedRunOperationError{AgentError: base}
	}
	if provider, helpURL := integrationDetails(base.Details); provider != "" {
		return &IntegrationNotConnectedError{AgentError: base, Provider: provider, HelpURL: helpURL}
	}
	if base.Code == "" {
		return base
	}
	return &UnknownAgentError{AgentError: base}
}

func integrationDetails(details []map[string]any) (provider, helpURL string) {
	for _, detail := range details {
		if p := stringField(detail, "provider"); p != "" {
			provider = p
		}
		if u := stringField(detail, "helpUrl"); u != "" {
			helpURL = u
		}
	}
	return provider, helpURL
}

func configurationError(message, code string) error {
	return &ConfigurationError{AgentError: &AgentError{Message: message, Code: code}}
}

func newUnsupportedRunOperation(op string, message string) error {
	if message == "" {
		message = fmt.Sprintf("Run operation %q is not supported", op)
	}
	return &UnsupportedRunOperationError{
		AgentError: &AgentError{Message: message, Code: "unsupported_run_operation"},
		Operation:  op,
	}
}

// AsAgentError reports whether err is an *AgentError (including typed wrappers).
func AsAgentError(err error) (*AgentError, bool) {
	var e *AgentError
	if errors.As(err, &e) {
		return e, true
	}
	return nil, false
}

// IsAuthentication reports authentication/permission errors.
func IsAuthentication(err error) bool {
	var e *AgentError
	return errors.As(err, &e) && (e.Code == "authentication" || e.Code == "permission_denied")
}

// IsRateLimit reports rate limit errors.
func IsRateLimit(err error) bool {
	var e *AgentError
	return errors.As(err, &e) && e.Code == "rate_limit"
}

// IsConfiguration reports configuration errors.
func IsConfiguration(err error) bool {
	var e *AgentError
	return errors.As(err, &e) && (e.Code == "configuration" || e.Code == "bad_request")
}

// IsNotFound reports not-found errors.
func IsNotFound(err error) bool {
	var e *AgentError
	return errors.As(err, &e) && (e.Code == "not_found" || e.Code == "agent_not_found")
}

// IsAgentBusy reports agent busy errors.
func IsAgentBusy(err error) bool {
	var e *AgentError
	return errors.As(err, &e) && e.Code == "agent_busy"
}

// IsIntegrationNotConnected reports integration connection errors.
func IsIntegrationNotConnected(err error) bool {
	var e *IntegrationNotConnectedError
	return errors.As(err, &e)
}

// IsUnsupportedRunOperation reports unsupported run operation errors.
func IsUnsupportedRunOperation(err error) bool {
	var e *UnsupportedRunOperationError
	return errors.As(err, &e)
}
