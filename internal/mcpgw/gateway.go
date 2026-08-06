// Package mcpgw implements an MCP tool gateway adapted from Cloudflare
// OS's MCP integration pattern. It provides tool classification
// (read-only vs write), trust tiers (byo vs vetted), and automatic
// routing through the approval system for write actions.
package mcpgw

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/anubhavg-icpl/talon/internal/approval"
	"github.com/google/uuid"
)

// TrustTier classifies how much trust a tool endpoint has earned.
type TrustTier string

const (
	TierBYO    TrustTier = "byo"    // user-provided, unvetted
	TierVetted TrustTier = "vetted" // admin-approved, known-safe endpoint
)

// ToolClassification determines how a tool's calls are handled.
type ToolClassification string

const (
	ClassObservation ToolClassification = "observation" // read-only, auto-allowed
	ClassAction      ToolClassification = "action"      // has side-effects, needs approval
)

// ToolDescriptor describes a registered MCP tool.
type ToolDescriptor struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Endpoint    string                 `json:"endpoint"`
	Tier        TrustTier              `json:"tier"`
	Class       ToolClassification     `json:"class"`
	Hints       map[string]bool        `json:"hints,omitempty"`       // readOnlyHint, destructiveHint, etc.
	Schema      map[string]interface{} `json:"schema,omitempty"`      // input schema
	Vetted      bool                   `json:"vetted"`                // true if admin has vetted
}

// CallRequest represents a tool invocation through the gateway.
type CallRequest struct {
	ToolName string                 `json:"tool_name"`
	Args     map[string]interface{} `json:"args"`
	RunID    string                 `json:"run_id"`
	Caller   string                 `json:"caller"` // "agent" or "user"
}

// CallResult is the outcome of a tool call.
type CallResult struct {
	Success     bool                   `json:"success"`
	Output      interface{}            `json:"output,omitempty"`
	Error       string                 `json:"error,omitempty"`
	ApprovalID  string                 `json:"approval_id,omitempty"`
	AutoApproved bool                  `json:"auto_approved,omitempty"`
	Simulated   bool                   `json:"simulated,omitempty"` // true when waiting for approval
	Timestamp   time.Time              `json:"timestamp"`
}

// Gateway routes MCP tool calls through classification and approval.
type Gateway struct {
	mu      sync.RWMutex
	tools   map[string]*ToolDescriptor
	store   *approval.ActionStore
	client  *http.Client
	mcpBase string // base URL of MCP server
}

// New creates a Gateway wired to the given approval store.
func New(store *approval.ActionStore, mcpBaseURL string) *Gateway {
	return &Gateway{
		tools:   make(map[string]*ToolDescriptor),
		store:   store,
		client:  &http.Client{Timeout: 120 * time.Second},
		mcpBase: strings.TrimRight(mcpBaseURL, "/"),
	}
}

// Register adds or updates a tool descriptor.
func (g *Gateway) Register(td ToolDescriptor) {
	g.mu.Lock()
	defer g.mu.Unlock()
	// Auto-classify if not set
	if td.Class == "" {
		td.Class = classifyTool(td)
	}
	if td.Tier == "" {
		td.Tier = TierBYO
	}
	g.tools[td.Name] = &td
}

// Unregister removes a tool.
func (g *Gateway) Unregister(name string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	delete(g.tools, name)
}

// List returns all registered tools.
func (g *Gateway) List() []ToolDescriptor {
	g.mu.RLock()
	defer g.mu.RUnlock()
	result := make([]ToolDescriptor, 0, len(g.tools))
	for _, td := range g.tools {
		result = append(result, *td)
	}
	return result
}

// Get returns a tool descriptor by name.
func (g *Gateway) Get(name string) (*ToolDescriptor, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	td, ok := g.tools[name]
	if !ok {
		return nil, fmt.Errorf("mcpgw: tool %q not registered", name)
	}
	return td, nil
}

// Call executes a tool through the gateway, applying classification and
// approval rules. For actions that need approval, it returns a simulated
// result immediately and queues the action for async approval (adapted
// from CF-OS Gatekeeper pattern).
func (g *Gateway) Call(ctx context.Context, req CallRequest) (*CallResult, error) {
	td, err := g.Get(req.ToolName)
	if err != nil {
		return nil, err
	}

	// Observations are always allowed
	if td.Class == ClassObservation {
		return g.executeDirect(ctx, req, td)
	}

	// Actions need approval
	risk := approval.RiskMedium
	if td.Hints["destructiveHint"] || approval.IsDangerous(req.ToolName) {
		risk = approval.RiskHigh
	}

	// Check auto-approve
	if g.store != nil && g.store.AutoApprove(risk) {
		result, err := g.executeDirect(ctx, req, td)
		if err != nil {
			return &CallResult{
				Success:   false,
				Error:     err.Error(),
				Timestamp: time.Now(),
			}, nil
		}
		result.AutoApproved = true
		return result, nil
	}

	// Queue for approval - create action entry
	argsJSON, _ := json.Marshal(req.Args)
	action := &approval.Action{
		ID:       uuid.NewString(),
		RunID:    req.RunID,
		ToolName: req.ToolName,
		Args:     argsJSON,
		State:    approval.StatePending,
		RiskLevel: risk,
		Summary:  fmt.Sprintf("MCP tool %s called by %s", req.ToolName, req.Caller),
	}

	if g.store != nil {
		if err := g.store.Create(action); err != nil {
			return nil, fmt.Errorf("mcpgw: create approval: %w", err)
		}
	}

	// Return simulated result (CF-OS pattern: agent proceeds with simulated data)
	return &CallResult{
		Success:    true,
		Simulated:  true,
		ApprovalID: action.ID,
		Output:     map[string]string{"status": "queued for approval"},
		Timestamp:  time.Now(),
	}, nil
}

// ExecuteApproved runs a previously-queued action after it's been approved.
func (g *Gateway) ExecuteApproved(ctx context.Context, actionID string) (*CallResult, error) {
	if g.store == nil {
		return nil, fmt.Errorf("mcpgw: no approval store configured")
	}

	action, err := g.store.Claim(actionID)
	if err != nil {
		return nil, fmt.Errorf("mcpgw: claim action: %w", err)
	}

	td, err := g.Get(action.ToolName)
	if err != nil {
		_ = g.store.Fail(actionID, fmt.Sprintf("tool not found: %s", action.ToolName))
		return nil, err
	}

	var args map[string]interface{}
	if err := json.Unmarshal(action.Args, &args); err != nil {
		_ = g.store.Fail(actionID, "invalid args JSON")
		return nil, err
	}

	req := CallRequest{
		ToolName: action.ToolName,
		Args:     args,
		RunID:    action.RunID,
		Caller:   "approved",
	}

	result, err := g.executeDirect(ctx, req, td)
	if err != nil {
		_ = g.store.Fail(actionID, err.Error())
		return nil, err
	}

	resultJSON, _ := json.Marshal(result)
	_ = g.store.Approve(actionID, resultJSON)
	return result, nil
}

// RejectAction marks a queued action as rejected.
func (g *Gateway) RejectAction(actionID, reason string) error {
	if g.store == nil {
		return fmt.Errorf("mcpgw: no approval store")
	}
	return g.store.Reject(actionID, reason)
}

// classifyTool determines if a tool is an observation (read-only) or
// action (has side-effects), based on its hints. This is the single
// trust-boundary function, adapted from CF-OS's classifyTool().
func classifyTool(td ToolDescriptor) ToolClassification {
	// readOnlyHint must be strictly true (matching CF-OS strict === true check)
	if td.Hints["readOnlyHint"] == true {
		return ClassObservation
	}
	// Destructive tools are always actions
	if td.Hints["destructiveHint"] == true {
		return ClassAction
	}
	// Unannotated tools default to action (safe default)
	return ClassAction
}

// executeDirect calls the MCP endpoint without approval checks.
func (g *Gateway) executeDirect(ctx context.Context, req CallRequest, td *ToolDescriptor) (*CallResult, error) {
	payload, err := json.Marshal(map[string]interface{}{
		"tool":  req.ToolName,
		"args":  req.Args,
	})
	if err != nil {
		return nil, err
	}

	url := g.mcpBase + "/tools/call"
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(string(payload)))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := g.client.Do(httpReq)
	if err != nil {
		return &CallResult{
			Success:   false,
			Error:     err.Error(),
			Timestamp: time.Now(),
		}, nil
	}
	defer resp.Body.Close()

	var output interface{}
	if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
		output = nil
	}

	success := resp.StatusCode >= 200 && resp.StatusCode < 300
	result := &CallResult{
		Success:   success,
		Output:    output,
		Timestamp: time.Now(),
	}
	if !success {
		result.Error = fmt.Sprintf("HTTP %d", resp.StatusCode)
	}
	return result, nil
}

// Vet marks a tool as vetted by an admin (promotes from byo to vetted tier).
func (g *Gateway) Vet(toolName string) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	td, ok := g.tools[toolName]
	if !ok {
		return fmt.Errorf("mcpgw: tool %q not found", toolName)
	}
	td.Tier = TierVetted
	td.Vetted = true
	return nil
}

// Stats returns summary statistics about registered tools.
func (g *Gateway) Stats() map[string]int {
	g.mu.RLock()
	defer g.mu.RUnlock()
	stats := map[string]int{
		"total":       len(g.tools),
		"observations": 0,
		"actions":     0,
		"vetted":      0,
		"byo":         0,
	}
	for _, td := range g.tools {
		if td.Class == ClassObservation {
			stats["observations"]++
		} else {
			stats["actions"]++
		}
		if td.Tier == TierVetted {
			stats["vetted"]++
		} else {
			stats["byo"]++
		}
	}
	return stats
}
