package arsenal

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Register wires every generated proxy tool plus the hand-ported composite
// tools (tools_manual.go) onto srv, all backed by client.
func Register(srv *server.MCPServer, client *Client) {
	for _, spec := range generatedTools {
		srv.AddTool(buildMCPTool(spec), buildHandler(spec, client))
	}
	registerManualTools(srv, client)
}

func buildMCPTool(spec toolSpec) mcp.Tool {
	opts := []mcp.ToolOption{mcp.WithDescription(spec.Description)}
	for _, p := range spec.Params {
		var popts []mcp.PropertyOption
		if p.Required {
			popts = append(popts, mcp.Required())
		}
		switch p.Kind {
		case "integer":
			if n, ok := asInt(p.Default); ok {
				popts = append(popts, mcp.DefaultNumber(n))
			}
			opts = append(opts, mcp.WithNumber(p.Name, popts...))
		case "number":
			if f, ok := p.Default.(float64); ok {
				popts = append(popts, mcp.DefaultNumber(f))
			}
			opts = append(opts, mcp.WithNumber(p.Name, popts...))
		case "boolean":
			if b, ok := p.Default.(bool); ok {
				popts = append(popts, mcp.DefaultBool(b))
			}
			opts = append(opts, mcp.WithBoolean(p.Name, popts...))
		case "object":
			opts = append(opts, mcp.WithObject(p.Name, popts...))
		case "array":
			opts = append(opts, mcp.WithArray(p.Name, popts...))
		default: // "string"
			if s, ok := p.Default.(string); ok {
				popts = append(popts, mcp.DefaultString(s))
			}
			opts = append(opts, mcp.WithString(p.Name, popts...))
		}
	}
	return mcp.NewTool(spec.Name, opts...)
}

func asInt(v any) (int64, bool) {
	switch n := v.(type) {
	case int:
		return int64(n), true
	case int64:
		return n, true
	case float64:
		return int64(n), true
	default:
		return 0, false
	}
}

// buildHandler implements the uniform proxy behavior every generated tool
// shares: gather declared params (falling back to their defaults), map
// them onto the Arsenal Engine's request shape, and forward.
func buildHandler(spec toolSpec, client *Client) server.ToolHandlerFunc {
	return func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := req.GetArguments()

		values := make(map[string]any, len(spec.Params))
		for _, p := range spec.Params {
			if v, ok := args[p.Name]; ok {
				values[p.Name] = v
			} else {
				values[p.Name] = p.Default
			}
		}

		body := map[string]any{}
		for k, v := range spec.StaticFields {
			body[k] = v
		}
		for _, bf := range spec.BodyFields {
			body[bf.JSONKey] = values[bf.ParamName]
		}

		endpoint := spec.Endpoint
		if len(spec.EndpointVars) > 0 {
			vals := make([]any, len(spec.EndpointVars))
			for i, name := range spec.EndpointVars {
				vals[i] = values[name]
			}
			endpoint = fmt.Sprintf(endpoint, vals...)
		}

		var result map[string]any
		if spec.HTTPMethod == "GET" {
			result = client.Get(endpoint, body)
		} else {
			result = client.Post(endpoint, body)
		}

		out, err := json.Marshal(result)
		if err != nil {
			return mcp.NewToolResultErrorFromErr("marshal result", err), nil
		}
		return mcp.NewToolResultText(string(out)), nil
	}
}
