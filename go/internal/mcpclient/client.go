// Package mcpclient wraps mark3labs/mcp-go's stdio client so the agent
// orchestrator can spawn the hexstrike-mcp and metasploit-mcp binaries and
// call their tools, mirroring langchain_mcp_adapters.MultiServerMCPClient
// used in pentest_core/final.py and pentest_core/integration.py.
package mcpclient

import (
	"context"
	"fmt"
	"sync"

	"github.com/anubhavg-icpl/pentester2/go/internal/llm"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

// ServerSpec describes one stdio MCP server to launch.
type ServerSpec struct {
	Name    string
	Command string
	Args    []string
}

// Multi manages a set of MCP servers and exposes their tools under a single
// flat namespace, resolving each ToolSpec.Name back to its owning server on
// Call, same as MultiServerMCPClient.get_tools() + tool.ainvoke().
type Multi struct {
	mu      sync.RWMutex
	clients map[string]*client.Client // server name -> connected client
	owner   map[string]string         // tool name -> server name
	tools   []llm.ToolSpec
}

func NewMulti(ctx context.Context, specs []ServerSpec) (*Multi, error) {
	m := &Multi{
		clients: make(map[string]*client.Client),
		owner:   make(map[string]string),
	}
	for _, spec := range specs {
		c, err := client.NewStdioMCPClient(spec.Command, nil, spec.Args...)
		if err != nil {
			return nil, fmt.Errorf("mcpclient: start %s: %w", spec.Name, err)
		}
		if _, err := c.Initialize(ctx, mcp.InitializeRequest{}); err != nil {
			return nil, fmt.Errorf("mcpclient: initialize %s: %w", spec.Name, err)
		}
		listing, err := c.ListTools(ctx, mcp.ListToolsRequest{})
		if err != nil {
			return nil, fmt.Errorf("mcpclient: list tools %s: %w", spec.Name, err)
		}
		m.clients[spec.Name] = c
		for _, t := range listing.Tools {
			schema := map[string]any{}
			if t.InputSchema.Properties != nil {
				schema["properties"] = t.InputSchema.Properties
				schema["required"] = t.InputSchema.Required
				schema["type"] = "object"
			}
			m.tools = append(m.tools, llm.ToolSpec{
				Name:        t.Name,
				Description: t.Description,
				InputSchema: schema,
			})
			m.owner[t.Name] = spec.Name
		}
	}
	return m, nil
}

// Tools returns every tool discovered across all connected servers.
func (m *Multi) Tools() []llm.ToolSpec {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return append([]llm.ToolSpec(nil), m.tools...)
}

// Subset returns only the named tools, preserving discovery order -- mirrors
// final.py's tools_named(*names) helper used to scope each subagent.
func (m *Multi) Subset(names ...string) []llm.ToolSpec {
	want := make(map[string]bool, len(names))
	for _, n := range names {
		want[n] = true
	}
	var out []llm.ToolSpec
	for _, t := range m.Tools() {
		if want[t.Name] {
			out = append(out, t)
		}
	}
	return out
}

// Call invokes a tool by name on whichever server owns it.
func (m *Multi) Call(ctx context.Context, name string, args map[string]any) (string, error) {
	m.mu.RLock()
	serverName, ok := m.owner[name]
	m.mu.RUnlock()
	if !ok {
		return "", fmt.Errorf("mcpclient: unknown tool %q", name)
	}
	c := m.clients[serverName]

	req := mcp.CallToolRequest{}
	req.Params.Name = name
	req.Params.Arguments = args

	res, err := c.CallTool(ctx, req)
	if err != nil {
		return "", fmt.Errorf("mcpclient: call %s: %w", name, err)
	}
	var out string
	for _, content := range res.Content {
		if tc, ok := mcp.AsTextContent(content); ok {
			out += tc.Text
		}
	}
	if res.IsError {
		return out, fmt.Errorf("mcpclient: tool %s reported error", name)
	}
	return out, nil
}

// Close shuts down every underlying MCP server process.
func (m *Multi) Close() error {
	var firstErr error
	for _, c := range m.clients {
		if err := c.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}
