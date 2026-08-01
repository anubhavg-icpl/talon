// Package mcpclient wraps mark3labs/mcp-go's stdio client so the agent
// orchestrator can spawn the talon-arsenal and talon-strike binaries and
// call their tools through a single multiplexed client.
package mcpclient

import (
	"context"
	"fmt"
	"log"
	"sort"
	"sync"
	"time"

	"github.com/anubhavg-icpl/talon/internal/llm"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

// Connect retry: a required server (e.g. metasploit → msfrpcd) may not be
// reachable the instant core boots. Retry initialize/list so core waits for it
// to come up instead of crashing and relying on a docker restart loop.
// ~60s of headroom, which comfortably covers msfrpcd startup.
const (
	connectAttempts = 20
	connectDelay    = 3 * time.Second
)

// ServerSpec describes one stdio MCP server to launch.
type ServerSpec struct {
	Name    string
	Command string
	Args    []string
	// Optional marks a best-effort server: if it fails to start, initialize,
	// or list tools, the failure is logged and the server is skipped rather
	// than aborting the whole orchestrator. Core tools (hexstrike/metasploit)
	// stay strict; add-on servers (e.g. lightpanda) are optional.
	Optional bool
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
		var (
			c       *client.Client
			tools   []mcp.Tool
			lastErr error
		)
		for attempt := 1; attempt <= connectAttempts; attempt++ {
			cc, err := client.NewStdioMCPClient(spec.Command, nil, spec.Args...)
			if err != nil {
				// Spawn failure (binary missing / can't exec) is a hard error —
				// retrying won't fix it, so stop and let the caller skip/fail.
				lastErr = fmt.Errorf("start: %w", err)
				break
			}
			if _, err = cc.Initialize(ctx, mcp.InitializeRequest{}); err != nil {
				lastErr = fmt.Errorf("initialize: %w", err)
				_ = cc.Close()
			} else if listing, lerr := cc.ListTools(ctx, mcp.ListToolsRequest{}); lerr != nil {
				lastErr = fmt.Errorf("list tools: %w", lerr)
				_ = cc.Close()
			} else {
				c, tools, lastErr = cc, listing.Tools, nil
				break
			}
			if attempt < connectAttempts {
				select {
				case <-ctx.Done():
					return nil, ctx.Err()
				case <-time.After(connectDelay):
				}
			}
		}
		if c == nil {
			if spec.Optional {
				log.Printf("mcpclient: optional server %s unavailable: %v", spec.Name, lastErr)
				continue
			}
			return nil, fmt.Errorf("mcpclient: %s: %w", spec.Name, lastErr)
		}
		m.clients[spec.Name] = c
		for _, t := range tools {
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

// Subset returns only the named tools, preserving discovery order -- used
// to scope each subagent to its own tool set.
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

// ServerInfo describes one connected MCP server for the dashboard.
type ServerInfo struct {
	Name  string   `json:"name"`
	Tools []string `json:"tools"`
}

// Servers returns the connected MCP servers with their tool names.
func (m *Multi) Servers() []ServerInfo {
	m.mu.RLock()
	defer m.mu.RUnlock()
	byServer := map[string][]string{}
	for _, t := range m.tools {
		byServer[m.owner[t.Name]] = append(byServer[m.owner[t.Name]], t.Name)
	}
	out := make([]ServerInfo, 0, len(m.clients))
	for name := range m.clients {
		out = append(out, ServerInfo{Name: name, Tools: byServer[name]})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Name < out[j].Name })
	return out
}

// Call invokes a tool by name on whichever server owns it.
func (m *Multi) Call(ctx context.Context, name string, args map[string]any) (string, error) {
	m.mu.RLock()
	serverName, ok := m.owner[name]
	c := m.clients[serverName]
	m.mu.RUnlock()
	if !ok || c == nil {
		return "", fmt.Errorf("mcpclient: unknown tool %q", name)
	}

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
