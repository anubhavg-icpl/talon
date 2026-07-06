// Package strike is a msgpack-RPC client for Metasploit's msfrpcd, plus the
// 12 MCP tools built on top of it (see tools.go).
package strike

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/anubhavg-icpl/ talon/internal/config"
	"github.com/google/uuid"
	"github.com/vmihailenco/msgpack/v5"
)

// msfDefaultUsername is the RPC username msfrpcd expects when none is
// configured -- config.MSFConfig has no Username field since it's never
// overridden in practice.
const msfDefaultUsername = "msf"

// Client is a msgpack-RPC client for msfrpcd.
type Client struct {
	baseURL        string
	http           *http.Client
	token          string
	PayloadSaveDir string
}

// NewClient dials server, authenticates with auth.login, then mints a fresh
// UUID token and registers it via auth.token_add (sent under the login
// token, per the "every call except auth.login is [method, token,
// ...params]" rule), switching to the new token for every subsequent call.
// This is how msfrpcd expects a long-lived client to behave.
func NewClient(ctx context.Context, cfg config.MSFConfig) (*Client, error) {
	if cfg.Password == "" {
		return nil, fmt.Errorf("msfrpc: MSF_PASSWORD is required")
	}
	scheme := "http"
	if cfg.SSL {
		scheme = "https"
	}
	c := &Client{
		baseURL:        fmt.Sprintf("%s://%s:%s/api/", scheme, cfg.Server, cfg.Port),
		http:           &http.Client{},
		PayloadSaveDir: cfg.PayloadSaveDir,
	}

	loginResp, err := c.doCall(ctx, "auth.login", []any{"auth.login", msfDefaultUsername, cfg.Password})
	if err != nil {
		return nil, fmt.Errorf("msfrpc: login failed: %w", err)
	}
	if result, _ := loginResp["result"].(string); result != "success" {
		return nil, fmt.Errorf("msfrpc: login failed: unexpected result %v", loginResp["result"])
	}
	loginToken, _ := loginResp["token"].(string)
	if loginToken == "" {
		return nil, fmt.Errorf("msfrpc: login succeeded but no token was returned")
	}

	newToken := uuid.NewString()
	c.token = loginToken
	if _, err := c.Call(ctx, "auth.token_add", newToken); err != nil {
		return nil, fmt.Errorf("msfrpc: auth.token_add failed: %w", err)
	}
	c.token = newToken

	return c, nil
}

// Call issues an authenticated RPC: [method, token, params...].
func (c *Client) Call(ctx context.Context, method string, params ...any) (map[string]any, error) {
	args := make([]any, 0, len(params)+2)
	args = append(args, method, c.token)
	args = append(args, params...)
	return c.doCall(ctx, method, args)
}

func (c *Client) doCall(ctx context.Context, method string, args []any) (map[string]any, error) {
	body, err := msgpack.Marshal(args)
	if err != nil {
		return nil, fmt.Errorf("msfrpc: encode %s request: %w", method, err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("msfrpc: build %s request: %w", method, err)
	}
	req.Header.Set("Content-Type", "binary/message-pack")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("msfrpc: %s request failed: %w", method, err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("msfrpc: read %s response: %w", method, err)
	}

	var out map[string]any
	if err := msgpack.Unmarshal(raw, &out); err != nil {
		return nil, fmt.Errorf("msfrpc: decode %s response (HTTP %d): %w", method, resp.StatusCode, err)
	}

	if errFlag, _ := out["error"].(bool); errFlag {
		msg := firstNonEmptyString(out["error_message"], out["error_string"], "unknown RPC error")
		return out, fmt.Errorf("msfrpc: %s: %s", method, msg)
	}
	return out, nil
}

// CoreVersion mirrors core.version, useful to verify connectivity.
func (c *Client) CoreVersion(ctx context.Context) (map[string]any, error) {
	return c.Call(ctx, "core.version")
}

func firstNonEmptyString(vals ...any) string {
	for _, v := range vals {
		if s, ok := v.(string); ok && s != "" {
			return s
		}
	}
	return ""
}
