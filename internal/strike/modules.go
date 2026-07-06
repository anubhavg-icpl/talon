package strike

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// ListExploits mirrors module.exploits: {"modules": [...]}.
func (c *Client) ListExploits(ctx context.Context) ([]string, error) {
	return c.listModuleNames(ctx, "module.exploits")
}

// ListPayloads mirrors module.payloads.
func (c *Client) ListPayloads(ctx context.Context) ([]string, error) {
	return c.listModuleNames(ctx, "module.payloads")
}

// ListAuxiliary mirrors module.auxiliary.
func (c *Client) ListAuxiliary(ctx context.Context) ([]string, error) {
	return c.listModuleNames(ctx, "module.auxiliary")
}

// ListPost mirrors module.post.
func (c *Client) ListPost(ctx context.Context) ([]string, error) {
	return c.listModuleNames(ctx, "module.post")
}

func (c *Client) listModuleNames(ctx context.Context, method string) ([]string, error) {
	resp, err := c.Call(ctx, method)
	if err != nil {
		return nil, err
	}
	raw, _ := resp["modules"].([]any)
	names := make([]string, 0, len(raw))
	for _, v := range raw {
		if s, ok := v.(string); ok {
			names = append(names, s)
		}
	}
	return names, nil
}

// moduleTypePrefixes are the path prefixes _get_module_object strips off a
// module name before use -- module.execute wants the base name (e.g.
// "windows/smb/ms17_010_eternalblue", not "exploit/windows/...").
var moduleTypePrefixes = map[string]bool{
	"exploit": true, "payload": true, "post": true,
	"auxiliary": true, "encoder": true, "nop": true,
}

func normalizeModuleName(name string) string {
	if !strings.Contains(name, "/") {
		return name
	}
	parts := strings.SplitN(name, "/", 2)
	if moduleTypePrefixes[parts[0]] {
		return parts[1]
	}
	return name
}

// Execute runs a module as a background job via module.execute -- the only
// execution path implemented; see tools.go for why a console-based
// fallback was left out. A client-side module.info/module.options
// validation round trip is skipped too: module.execute validates
// server-side and fails the same way.
func (c *Client) Execute(ctx context.Context, modtype, modname string, options map[string]any) (map[string]any, error) {
	return c.Call(ctx, "module.execute", modtype, normalizeModuleName(modname), coerceOptionTypes(options))
}

// ParseOptionsGracefully tolerates a few common ways module options arrive:
// a map (already correct), a "key=value,key2=value2" string, a JSON object
// string (common when a caller must pass options through a string-typed
// field), or nil/empty.
func ParseOptionsGracefully(v any) (map[string]any, error) {
	switch opts := v.(type) {
	case nil:
		return map[string]any{}, nil
	case map[string]any:
		return opts, nil
	case string:
		s := strings.TrimSpace(opts)
		if s == "" {
			return map[string]any{}, nil
		}
		if strings.HasPrefix(s, "{") {
			var m map[string]any
			if err := json.Unmarshal([]byte(s), &m); err == nil {
				return m, nil
			}
		}
		parsed := map[string]any{}
		for _, pair := range strings.Split(s, ",") {
			pair = strings.TrimSpace(pair)
			if pair == "" {
				continue
			}
			kv := strings.SplitN(pair, "=", 2)
			if len(kv) != 2 {
				return nil, fmt.Errorf("invalid option format: %q (missing '=')", pair)
			}
			key := strings.TrimSpace(kv[0])
			if key == "" {
				return nil, fmt.Errorf("invalid option format: %q (empty key)", pair)
			}
			value := strings.TrimSpace(kv[1])
			if len(value) >= 2 {
				if (value[0] == '"' && value[len(value)-1] == '"') || (value[0] == '\'' && value[len(value)-1] == '\'') {
					value = value[1 : len(value)-1]
				}
			}
			parsed[key] = coerceScalar(value)
		}
		return parsed, nil
	default:
		return nil, fmt.Errorf("options must be a map or 'key=value,key2=value2' string, got %T", v)
	}
}

// coerceOptionTypes mirrors _set_module_options' basic type guessing on
// every value before sending it to module.execute.
func coerceOptionTypes(options map[string]any) map[string]any {
	out := make(map[string]any, len(options))
	for k, v := range options {
		if s, ok := v.(string); ok {
			out[k] = coerceScalar(s)
		} else {
			out[k] = v
		}
	}
	return out
}

func coerceScalar(value string) any {
	lower := strings.ToLower(value)
	if lower == "true" || lower == "false" {
		return lower == "true"
	}
	if isAllDigits(value) {
		if n, err := strconv.Atoi(value); err == nil {
			return n
		}
	}
	return value
}

func isAllDigits(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}
