// Package strike is a msgpack-RPC client for Metasploit's msfrpcd, plus the
// 12 MCP tools built on top of it (see tools.go).
package strike

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"net/http"

	"github.com/anubhavg-icpl/talon/internal/config"
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
	// msfrpcd serves SSL with a self-signed cert that carries no SAN/CN, so
	// Go's default verifier rejects it ("certificate is not valid for any
	// names"). msfrpcd is a local, network-isolated daemon whose real auth
	// boundary is the RPC password, so skip TLS verification here -- the same
	// way every other msfrpcd client (msfrpc, metasploit's own) handles it.
	if cfg.SSL {
		c.http.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	loginResp, err := c.doCall(ctx, "auth.login", []any{"auth.login", msfDefaultUsername, cfg.Password})
	if err != nil {
		return nil, fmt.Errorf("msfrpc: login failed: %w", err)
	}
	// msfrpcd encodes map keys/values as msgpack bin, so vmihailenco often
	// yields []byte rather than string for "result" / "token".
	if result := asString(loginResp["result"]); result != "success" {
		return nil, fmt.Errorf("msfrpc: login failed: unexpected result %v", loginResp["result"])
	}
	loginToken := asString(loginResp["token"])
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

	// msfrpcd returns maps with non-string keys for some methods
	// (session.list / job.list use integer session/job IDs). vmihailenco's
	// default interface decoder assumes string keys and fails with
	// "msgpack: invalid code=N decoding string/bytes length", which made
	// pollForSession miss every real shell session.
	out, err := decodeMSFMap(raw)
	if err != nil {
		return nil, fmt.Errorf("msfrpc: decode %s response (HTTP %d): %w", method, resp.StatusCode, err)
	}

	if errFlag, _ := out["error"].(bool); errFlag {
		msg := firstNonEmptyString(out["error_message"], out["error_string"], "unknown RPC error")
		return out, fmt.Errorf("msfrpc: %s: %s", method, msg)
	}
	return out, nil
}

// decodeMSFMap decodes an msfrpcd msgpack response into map[string]any,
// accepting integer / bin map keys and normalizing bin string values.
func decodeMSFMap(raw []byte) (map[string]any, error) {
	dec := msgpack.NewDecoder(bytes.NewReader(raw))
	dec.SetMapDecoder(func(d *msgpack.Decoder) (any, error) {
		n, err := d.DecodeMapLen()
		if err != nil {
			return nil, err
		}
		if n < 0 {
			n = 0
		}
		m := make(map[string]any, n)
		for i := 0; i < n; i++ {
			var key any
			if err := d.Decode(&key); err != nil {
				return nil, err
			}
			var val any
			if err := d.Decode(&val); err != nil {
				return nil, err
			}
			m[msgpackKeyString(key)] = normalizeMsgpack(val)
		}
		return m, nil
	})

	var decoded any
	if err := dec.Decode(&decoded); err != nil {
		return nil, err
	}
	return asStringKeyedMap(decoded)
}

// asStringKeyedMap normalizes msgpack maps so callers always see
// map[string]any. Integer keys (session/job IDs) and []byte keys/values
// (msfrpcd bin encoding) become strings.
func asStringKeyedMap(v any) (map[string]any, error) {
	if v == nil {
		return map[string]any{}, nil
	}
	norm := normalizeMsgpack(v)
	switch t := norm.(type) {
	case map[string]any:
		return t, nil
	default:
		return nil, fmt.Errorf("expected map, got %T", v)
	}
}

func normalizeMsgpack(v any) any {
	switch t := v.(type) {
	case map[string]any:
		out := make(map[string]any, len(t))
		for k, val := range t {
			out[k] = normalizeMsgpack(val)
		}
		return out
	case map[any]any:
		out := make(map[string]any, len(t))
		for k, val := range t {
			out[msgpackKeyString(k)] = normalizeMsgpack(val)
		}
		return out
	case []any:
		for i := range t {
			t[i] = normalizeMsgpack(t[i])
		}
		return t
	case []byte:
		// msfrpcd commonly encodes strings as msgpack bin.
		return string(t)
	default:
		return v
	}
}

func msgpackKeyString(k any) string {
	switch t := k.(type) {
	case string:
		return t
	case []byte:
		return string(t)
	default:
		return fmt.Sprint(t)
	}
}

// CoreVersion mirrors core.version, useful to verify connectivity.
func (c *Client) CoreVersion(ctx context.Context) (map[string]any, error) {
	return c.Call(ctx, "core.version")
}

// asString normalizes msgpack-decoded strings: msfrpcd often sends bin
// (decodes as []byte) while older servers use str (decodes as string).
func asString(v any) string {
	switch t := v.(type) {
	case string:
		return t
	case []byte:
		return string(t)
	default:
		return ""
	}
}

func firstNonEmptyString(vals ...any) string {
	for _, v := range vals {
		if s := asString(v); s != "" {
			return s
		}
	}
	return ""
}
