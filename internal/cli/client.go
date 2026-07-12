package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client is an HTTP client for the talon-core control plane.
type Client struct {
	baseURL    string
	httpClient *http.Client
}

// NewClient builds a client against coreBase (e.g. http://localhost:8000).
func NewClient(coreBase string, timeout time.Duration) (*Client, error) {
	base := strings.TrimRight(strings.TrimSpace(coreBase), "/")
	if base == "" {
		return nil, fmt.Errorf("core URL is empty")
	}
	if _, err := url.ParseRequestURI(base); err != nil {
		return nil, fmt.Errorf("invalid core URL %q: %w", coreBase, err)
	}
	if timeout <= 0 {
		timeout = 30 * time.Second
	}
	return &Client{
		baseURL: base,
		httpClient: &http.Client{
			Timeout: timeout,
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				DialContext: (&net.Dialer{
					Timeout:   5 * time.Second,
					KeepAlive: 30 * time.Second,
				}).DialContext,
				MaxIdleConns:          10,
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   5 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
			},
		},
	}, nil
}

// BaseURL returns the configured core URL.
func (c *Client) BaseURL() string { return c.baseURL }

// StartRequest is POST /input/start.
type StartRequest struct {
	IP          string `json:"ip"`
	CVEID       string `json:"cve_id,omitempty"`
	ServiceName string `json:"service_name,omitempty"`
	Description string `json:"description,omitempty"`
	LHOST       string `json:"lhost,omitempty"`
	LPORT       int    `json:"lport,omitempty"`
}

// StartResponse is the body returned by POST /input/start.
type StartResponse struct {
	RunID   string `json:"run_id"`
	Message string `json:"message"`
}

// PendingInterrupt mirrors core.PendingInterrupt over the wire.
// JSON keys match the exported Go field names currently emitted by
// encoding/json on the control server (no json tags on the struct).
type PendingInterrupt struct {
	ToolName string         `json:"ToolName"`
	Args     map[string]any `json:"Args"`
}

// StatusResponse is GET /output/status/{run_id}.
type StatusResponse struct {
	Status    string            `json:"status"`
	Output    string            `json:"output"`
	Interrupt *PendingInterrupt `json:"interrupt"`
}

// ResumeRequest is POST /output/resume/{run_id}.
type ResumeRequest struct {
	Decision   string         `json:"decision"`
	EditedArgs map[string]any `json:"edited_args,omitempty"`
}

// ResumeResponse is the body returned by POST /output/resume/{run_id}.
type ResumeResponse struct {
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
}

// ToolCallRecord is one logged tool invocation (GET /monitor/tools).
type ToolCallRecord struct {
	Index    int            `json:"Index"`
	ToolName string         `json:"ToolName"`
	Args     map[string]any `json:"Args"`
	Output   string         `json:"Output"`
}

// ToolsResponse is GET /monitor/tools?run_id=...
type ToolsResponse struct {
	ToolLog []ToolCallRecord `json:"tool_log"`
}

// TracesResponse is GET /monitor/traces/{run_id}.
type TracesResponse struct {
	History []string `json:"history"`
}

// APIError is a non-2xx response from talon-core.
type APIError struct {
	StatusCode int
	Detail     string
	Body       string
}

func (e *APIError) Error() string {
	if e.Detail != "" {
		return fmt.Sprintf("talon-core HTTP %d: %s", e.StatusCode, e.Detail)
	}
	if e.Body != "" {
		return fmt.Sprintf("talon-core HTTP %d: %s", e.StatusCode, e.Body)
	}
	return fmt.Sprintf("talon-core HTTP %d", e.StatusCode)
}

// ProbeCore checks that something is listening on the control plane.
// Prefers GET /health; falls back to any HTTP response on GET /.
func (c *Client) ProbeCore(ctx context.Context) error {
	for _, path := range []string{"/health", "/"} {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+path, nil)
		if err != nil {
			return err
		}
		resp, err := c.httpClient.Do(req)
		if err != nil {
			return fmt.Errorf("unreachable: %w", err)
		}
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
		if path == "/health" && resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return nil
		}
		if path == "/" {
			// Any HTTP response means the process is up.
			return nil
		}
	}
	return nil
}

// Start begins a new validation run.
func (c *Client) Start(ctx context.Context, body StartRequest) (*StartResponse, error) {
	var out StartResponse
	if err := c.doJSON(ctx, http.MethodPost, "/input/start", body, &out); err != nil {
		return nil, err
	}
	if out.RunID == "" {
		return nil, fmt.Errorf("start succeeded but run_id was empty")
	}
	return &out, nil
}

// Status fetches run status.
func (c *Client) Status(ctx context.Context, runID string) (*StatusResponse, error) {
	var out StatusResponse
	if err := c.doJSON(ctx, http.MethodGet, "/output/status/"+url.PathEscape(runID), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Resume sends a HITL decision for a paused run.
func (c *Client) Resume(ctx context.Context, runID string, body ResumeRequest) (*ResumeResponse, error) {
	var out ResumeResponse
	if err := c.doJSON(ctx, http.MethodPost, "/output/resume/"+url.PathEscape(runID), body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Tools fetches the tool log for a run.
func (c *Client) Tools(ctx context.Context, runID string) (*ToolsResponse, error) {
	path := "/monitor/tools?run_id=" + url.QueryEscape(runID)
	var out ToolsResponse
	if err := c.doJSON(ctx, http.MethodGet, path, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// Traces fetches message history for a run.
func (c *Client) Traces(ctx context.Context, runID string) (*TracesResponse, error) {
	var out TracesResponse
	if err := c.doJSON(ctx, http.MethodGet, "/monitor/traces/"+url.PathEscape(runID), nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) doJSON(ctx context.Context, method, path string, body any, out any) error {
	var rdr io.Reader
	if body != nil {
		raw, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("encode request: %w", err)
		}
		rdr = bytes.NewReader(raw)
	}
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, rdr)
	if err != nil {
		return err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "talon-cli/"+Version)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(io.LimitReader(resp.Body, 16<<20)) // 16 MiB cap
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var detail struct {
			Detail string `json:"detail"`
		}
		_ = json.Unmarshal(raw, &detail)
		return &APIError{
			StatusCode: resp.StatusCode,
			Detail:     detail.Detail,
			Body:       strings.TrimSpace(string(raw)),
		}
	}

	if out == nil {
		return nil
	}
	if err := json.Unmarshal(raw, out); err != nil {
		// Retry after escaping raw control chars inside JSON strings
		// (nmap/nuclei blobs sometimes arrive unescaped from the control plane).
		if err2 := json.Unmarshal(sanitizeJSONStrings(raw), out); err2 != nil {
			return fmt.Errorf("decode response (HTTP %d): %w", resp.StatusCode, err)
		}
	}
	return nil
}

// sanitizeJSONStrings walks JSON text and escapes raw U+0000–U+001F bytes
// that appear inside string literals only. Structural whitespace is left alone.
func sanitizeJSONStrings(raw []byte) []byte {
	out := make([]byte, 0, len(raw)+32)
	inString := false
	escaped := false
	for i := 0; i < len(raw); i++ {
		b := raw[i]
		if !inString {
			out = append(out, b)
			if b == '"' {
				inString = true
				escaped = false
			}
			continue
		}
		if escaped {
			out = append(out, b)
			escaped = false
			continue
		}
		if b == '\\' {
			out = append(out, b)
			escaped = true
			continue
		}
		if b == '"' {
			out = append(out, b)
			inString = false
			continue
		}
		if b < 0x20 {
			switch b {
			case '\n':
				out = append(out, '\\', 'n')
			case '\r':
				out = append(out, '\\', 'r')
			case '\t':
				out = append(out, '\\', 't')
			default:
				out = append(out, fmt.Sprintf("\\u%04x", b)...)
			}
			continue
		}
		out = append(out, b)
	}
	return out
}
