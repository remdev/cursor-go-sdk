// Package connect implements the Connect JSON protocol used by cursor-sdk-bridge.
package connect

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	ProtocolVersion = "1"
	EndStreamFlag   = 0x02

	DefaultUnaryTimeout  = 60 * time.Second
	DefaultStreamTimeout = 600 * time.Second
)

var retrySafeMethods = map[string]struct{}{
	"Ping": {}, "GetVersion": {}, "GetAgent": {}, "ListAgents": {},
	"GetRun": {}, "ListRuns": {}, "WaitLiveRun": {}, "GetRunConversation": {},
	"Me": {}, "ListModels": {}, "ListRepositories": {},
	"ListAgentMessages": {}, "ListArtifacts": {},
}

// Transport speaks Connect JSON over HTTP to the SDK bridge.
type Transport struct {
	BaseURL       string
	AuthToken     string
	UnaryTimeout  time.Duration
	StreamTimeout time.Duration
	MaxRetries    int
	HTTPClient    *http.Client
	ownsClient    bool
}

// NewTransport creates a transport with SDK defaults.
func NewTransport(baseURL, authToken string, opts ...TransportOption) *Transport {
	t := &Transport{
		BaseURL:       strings.TrimRight(baseURL, "/") + "/",
		AuthToken:     authToken,
		UnaryTimeout:  DefaultUnaryTimeout,
		StreamTimeout: DefaultStreamTimeout,
		MaxRetries:    0,
	}
	for _, opt := range opts {
		opt(t)
	}
	if t.HTTPClient == nil {
		t.HTTPClient = &http.Client{Timeout: t.UnaryTimeout}
		t.ownsClient = true
	}
	return t
}

// TransportOption configures a Transport.
type TransportOption func(*Transport)

// WithHTTPClient sets a custom HTTP client. The transport will not close it.
func WithHTTPClient(c *http.Client) TransportOption {
	return func(t *Transport) {
		t.HTTPClient = c
		t.ownsClient = false
	}
}

// WithUnaryTimeout sets the default timeout for unary RPCs.
func WithUnaryTimeout(d time.Duration) TransportOption {
	return func(t *Transport) { t.UnaryTimeout = d }
}

// WithStreamTimeout sets the default timeout for streaming RPCs.
func WithStreamTimeout(d time.Duration) TransportOption {
	return func(t *Transport) { t.StreamTimeout = d }
}

// WithMaxRetries sets retry count for idempotent unary RPCs.
func WithMaxRetries(n int) TransportOption {
	return func(t *Transport) {
		if n < 0 {
			n = 0
		}
		t.MaxRetries = n
	}
}

// Close closes the owned HTTP client.
func (t *Transport) Close() {
	if t.ownsClient && t.HTTPClient != nil {
		t.HTTPClient.CloseIdleConnections()
	}
}

// WithOptions returns a transport view with overridden defaults sharing the HTTP client.
func (t *Transport) WithOptions(opts ...TransportOption) *Transport {
	clone := *t
	clone.ownsClient = false
	for _, opt := range opts {
		opt(&clone)
	}
	return &clone
}

// Unary performs a Connect unary JSON RPC.
func (t *Transport) Unary(ctx context.Context, service, method string, message map[string]any) (map[string]any, error) {
	body, err := json.Marshal(message)
	if err != nil {
		return nil, err
	}
	url := t.BaseURL + service + "/" + method
	headers := t.headers("application/json", "application/json")

	maxRetries := t.MaxRetries
	if _, safe := retrySafeMethods[method]; !safe {
		maxRetries = 0
	}

	var lastErr error
	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			sleepBeforeRetry(attempt)
		}
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
		if err != nil {
			return nil, err
		}
		for k, v := range headers {
			req.Header.Set(k, v)
		}
		resp, err := t.doUnary(ctx, req)
		if err != nil {
			lastErr = err
			if attempt < maxRetries && isRetryableError(err) {
				continue
			}
			return nil, err
		}
		defer resp.Body.Close()
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode >= 400 {
			parsed := parseHTTPError(resp.StatusCode, respBody, resp.Header)
			lastErr = parsed
			if attempt < maxRetries && isRetryableStatus(resp.StatusCode) {
				continue
			}
			return nil, parsed
		}
		if len(respBody) == 0 {
			return map[string]any{}, nil
		}
		var out map[string]any
		if err := json.Unmarshal(respBody, &out); err != nil {
			return nil, fmt.Errorf("bridge returned invalid JSON: %w", err)
		}
		return out, nil
	}
	return nil, lastErr
}

func (t *Transport) doUnary(ctx context.Context, req *http.Request) (*http.Response, error) {
	client := t.HTTPClient
	if client == nil {
		client = http.DefaultClient
	}
	ctx, cancel := context.WithTimeout(ctx, t.UnaryTimeout)
	defer cancel()
	req = req.WithContext(ctx)
	resp, err := client.Do(req)
	if err != nil {
		return nil, networkError(err)
	}
	return resp, nil
}

// Stream performs a Connect server-streaming JSON RPC.
func (t *Transport) Stream(ctx context.Context, service, method string, message map[string]any) (*StreamReader, error) {
	payload, err := EncodeStreamEnvelope(message, false)
	if err != nil {
		return nil, err
	}
	url := t.BaseURL + service + "/" + method
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(payload))
	if err != nil {
		return nil, err
	}
	for k, v := range t.headers("application/connect+json", "application/connect+json") {
		req.Header.Set(k, v)
	}
	client := t.HTTPClient
	if client == nil {
		client = http.DefaultClient
	}
	ctx, cancel := context.WithTimeout(ctx, t.StreamTimeout)
	req = req.WithContext(ctx)
	resp, err := client.Do(req)
	if err != nil {
		cancel()
		return nil, networkError(err)
	}
	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		cancel()
		return nil, parseHTTPError(resp.StatusCode, body, resp.Header)
	}
	return &StreamReader{
		resp:   resp,
		parser: newStreamParser(resp.StatusCode),
		cancel: cancel,
	}, nil
}

func (t *Transport) headers(accept, contentType string) map[string]string {
	return map[string]string{
		"Accept":                   accept,
		"Authorization":            "Bearer " + t.AuthToken,
		"Connect-Protocol-Version": ProtocolVersion,
		"Content-Type":             contentType,
	}
}

// PostShutdown sends a best-effort Shutdown RPC to the bridge.
func PostShutdown(baseURL, authToken string, graceSeconds int) {
	url := strings.TrimRight(baseURL, "/") + "/sdk.v1.SdkBridgeControlService/Shutdown"
	body, _ := json.Marshal(map[string]any{"graceSeconds": graceSeconds})
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+authToken)
	req.Header.Set("Connect-Protocol-Version", ProtocolVersion)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: DefaultUnaryTimeout}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	io.Copy(io.Discard, resp.Body)
	resp.Body.Close()
}
