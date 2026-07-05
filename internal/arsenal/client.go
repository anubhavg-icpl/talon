// Package hexstrike proxies the HexStrike AI HTTP API server as MCP tools,
// porting pentest_core/hexstrike_mcp.py. HexStrike itself (a prebuilt image,
// ghcr.io/suryaaramesh/suryaarc-hexstrike) is untouched -- this package only
// replaces the Python MCP client that sits in front of it.
package arsenal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	baseURL string
	http    *http.Client
}

func NewClient(baseURL string, timeout time.Duration) *Client {
	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		http:    &http.Client{Timeout: timeout},
	}
}

func (c *Client) url(endpoint string) string {
	return c.baseURL + "/" + strings.TrimLeft(endpoint, "/")
}

// Get mirrors HexStrikeClient.safe_get: GET with query params, never returns
// a Go error for HTTP/network failures -- it folds them into the result map
// as {"error": ..., "success": false}, matching the Python client's contract
// so callers don't need special-case error handling per tool.
func (c *Client) Get(endpoint string, params map[string]any) map[string]any {
	q := url.Values{}
	for k, v := range params {
		if v == nil {
			continue
		}
		q.Set(k, fmt.Sprintf("%v", v))
	}
	full := c.url(endpoint)
	if len(q) > 0 {
		full += "?" + q.Encode()
	}
	resp, err := c.http.Get(full)
	if err != nil {
		return map[string]any{"error": fmt.Sprintf("Request failed: %v", err), "success": false}
	}
	defer resp.Body.Close()
	return decodeOrError(resp)
}

// Post mirrors HexStrikeClient.safe_post.
func (c *Client) Post(endpoint string, data map[string]any) map[string]any {
	body, err := json.Marshal(data)
	if err != nil {
		return map[string]any{"error": fmt.Sprintf("Request failed: %v", err), "success": false}
	}
	resp, err := c.http.Post(c.url(endpoint), "application/json", bytes.NewReader(body))
	if err != nil {
		return map[string]any{"error": fmt.Sprintf("Request failed: %v", err), "success": false}
	}
	defer resp.Body.Close()
	return decodeOrError(resp)
}

func decodeOrError(resp *http.Response) map[string]any {
	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return map[string]any{
			"error":   fmt.Sprintf("HTTP %d: %s", resp.StatusCode, string(raw)),
			"success": false,
		}
	}
	var out map[string]any
	if err := json.Unmarshal(raw, &out); err != nil {
		return map[string]any{"error": fmt.Sprintf("Unexpected error: %v", err), "success": false}
	}
	return out
}

func (c *Client) ExecuteCommand(command string, useCache bool) map[string]any {
	return c.Post("api/command", map[string]any{"command": command, "use_cache": useCache})
}

func (c *Client) CheckHealth() map[string]any {
	return c.Get("health", nil)
}
