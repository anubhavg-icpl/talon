package core

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/anubhavg-icpl/talon/internal/llm"
)

// evidenceToolSpecs returns the agent-facing evidence tools.
// These let the model inspect stored tool outputs without re-running tools.
func evidenceToolSpecs() []llm.ToolSpec {
	return []llm.ToolSpec{
		{
			Name:        "evidence_list",
			Description: "List all stored evidence from previous tool calls. Each entry has an id, tool name, summary, and size. Use evidence_view to see full content.",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"limit": map[string]any{
						"type":        "integer",
						"description": "Max results (default 20, max 50)",
					},
				},
			},
		},
		{
			Name:        "evidence_view",
			Description: "View full raw content of a stored evidence record by id. Use when you need the complete output that was truncated in the preview.",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"id": map[string]any{
						"type":        "string",
						"description": "Evidence id from evidence_list (e.g. e001)",
					},
					"max_chars": map[string]any{
						"type":        "integer",
						"description": "Optional truncate limit (default 8000, max 20000)",
					},
				},
				"required": []any{"id"},
			},
		},
		{
			Name:        "evidence_search",
			Description: "Search stored evidence by substring. Returns matching evidence records with previews. Useful for finding specific responses, flags, or error messages.",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"q": map[string]any{
						"type":        "string",
						"description": "Search query (case-insensitive substring)",
					},
					"limit": map[string]any{
						"type":        "integer",
						"description": "Max results (default 10, max 20)",
					},
				},
				"required": []any{"q"},
			},
		},
	}
}

// handleEvidenceList handles the evidence_list tool call.
func handleEvidenceList(store *EvidenceStore, args map[string]any, tr *tracker) (string, bool) {
	if store == nil {
		return "evidence store not configured", true
	}
	limit := 20
	if v, ok := args["limit"].(float64); ok {
		limit = int(v)
	}
	if limit <= 0 {
		limit = 20
	}
	if limit > 50 {
		limit = 50
	}

	all := store.List()
	if limit < len(all) {
		all = all[:limit]
	}

	type item struct {
		ID       string `json:"id"`
		Tool     string `json:"tool"`
		Summary  string `json:"summary"`
		Size     int    `json:"size"`
		Status   int    `json:"status"`
		DuplicateOf string `json:"duplicate_of,omitempty"`
	}
	items := make([]item, 0, len(all))
	for _, e := range all {
		items = append(items, item{
			ID: e.ID, Tool: e.Tool, Summary: e.Summary, Size: e.Size(),
			Status: e.Status, DuplicateOf: e.DuplicateOf,
		})
	}
	payload := map[string]any{
		"total":   store.Count(),
		"showing": len(items),
		"items":   items,
	}
	raw, _ := json.MarshalIndent(payload, "", "  ")
	if tr != nil {
		tr.record("evidence_list", args, fmt.Sprintf("total=%d showing=%d", store.Count(), len(items)))
	}
	return string(raw), false
}

// handleEvidenceView handles the evidence_view tool call.
func handleEvidenceView(store *EvidenceStore, args map[string]any, tr *tracker) (string, bool) {
	if store == nil {
		return "evidence store not configured", true
	}
	id, _ := args["id"].(string)
	id = strings.TrimSpace(id)
	if id == "" {
		return "evidence_view: id is required", true
	}
	maxChars := 8000
	if v, ok := args["max_chars"].(float64); ok {
		maxChars = int(v)
	}
	if maxChars <= 0 {
		maxChars = 8000
	}
	if maxChars > 20000 {
		maxChars = 20000
	}

	rec, ok := store.Get(id)
	if !ok {
		return fmt.Sprintf("evidence %q not found", id), true
	}

	// If duplicate, resolve to original
	content := rec.Content
	if rec.DuplicateOf != "" {
		if orig, ok := store.Get(rec.DuplicateOf); ok {
			content = orig.Content
		}
	}

	truncated := false
	if len(content) > maxChars {
		content = content[:maxChars] + "\n…[truncated]"
		truncated = true
	}
	payload := map[string]any{
		"id":        rec.ID,
		"tool":      rec.Tool,
		"key_args":  rec.KeyArgs,
		"status":    rec.Status,
		"size":      rec.Size(),
		"truncated": truncated,
		"content":   content,
	}
	if rec.DuplicateOf != "" {
		payload["duplicate_of"] = rec.DuplicateOf
	}
	raw, _ := json.MarshalIndent(payload, "", "  ")
	if tr != nil {
		tr.record("evidence_view", args, fmt.Sprintf("id=%s bytes=%d", rec.ID, len(content)))
	}
	return string(raw), false
}

// handleEvidenceSearch handles the evidence_search tool call.
func handleEvidenceSearch(store *EvidenceStore, args map[string]any, tr *tracker) (string, bool) {
	if store == nil {
		return "evidence store not configured", true
	}
	q, _ := args["q"].(string)
	q = strings.TrimSpace(q)
	if q == "" {
		return "evidence_search: q is required", true
	}
	limit := 10
	if v, ok := args["limit"].(float64); ok {
		limit = int(v)
	}
	if limit <= 0 {
		limit = 10
	}
	if limit > 20 {
		limit = 20
	}

	results := store.Search(q, limit)
	type item struct {
		ID      string `json:"id"`
		Tool    string `json:"tool"`
		Summary string `json:"summary"`
		Size    int    `json:"size"`
	}
	items := make([]item, 0, len(results))
	for _, e := range results {
		items = append(items, item{
			ID: e.ID, Tool: e.Tool, Summary: e.Summary, Size: e.Size(),
		})
	}
	payload := map[string]any{
		"query":   q,
		"total":   len(results),
		"showing": len(items),
		"items":   items,
	}
	raw, _ := json.MarshalIndent(payload, "", "  ")
	if tr != nil {
		tr.record("evidence_search", args, fmt.Sprintf("q=%q hits=%d", q, len(results)))
	}
	return string(raw), false
}

// evidenceHybridExec wraps a base exec function to intercept evidence_* tool
// calls and route them to the evidence store. Other calls pass through.
func evidenceHybridExec(base toolExecFunc, store *EvidenceStore, tr *tracker) toolExecFunc {
	if store == nil {
		return base
	}
	return func(ctx context.Context, call llm.ToolCall) (string, bool) {
		switch call.Name {
		case "evidence_list":
			return handleEvidenceList(store, call.Args, tr)
		case "evidence_view":
			return handleEvidenceView(store, call.Args, tr)
		case "evidence_search":
			return handleEvidenceSearch(store, call.Args, tr)
		default:
			return base(ctx, call)
		}
	}
}

// recordToolEvidence captures a tool result into the evidence store.
// Called after every tool execution in the subagent loop.
func recordToolEvidence(store *EvidenceStore, toolName string, args map[string]any, output string, isErr bool) {
	if store == nil {
		return
	}
	keyArgs := fmt.Sprintf("%v", args)
	if len(keyArgs) > 200 {
		keyArgs = keyArgs[:200]
	}
	status := ParseStatusCode(output)
	store.Record(toolName, keyArgs, output, status)
}
