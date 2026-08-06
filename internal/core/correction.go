package core

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"
)

// CorrectionLayer is a lightweight tool-call lifecycle observer.
// Ported from pentest agent agent/correction_layer.py + anti_loop.py.
//
// It tracks:
//   - Duplicate tool-call detection (same tool + args hash within a window)
//   - Consecutive failure tracking per tool (degraded health)
//   - Stall detection (only re-reading evidence, no new tool families)
//   - Completion gate: validates claimed findings against stored evidence
type CorrectionLayer struct {
	mu sync.Mutex

	// callHistory stores recent tool+argsHash for duplicate detection.
	callHistory []callEntry
	historyCap  int

	// toolHealth tracks consecutive failures per tool.
	toolHealth map[string]*toolHealthEntry

	// toolFamilySeen records which tool families have been exercised.
	toolFamilySeen map[string]bool

	// evidenceSeen tracks evidence IDs the model has viewed (redundancy check).
	evidenceSeen map[string]int

	// rejectionCount tracks how many times the completion gate rejected an answer.
	rejectionCount int
}

type callEntry struct {
	tool      string
	argsHash  string
	timestamp time.Time
}

type toolHealthEntry struct {
	consecutiveFailures int
	status              string // "healthy" | "degraded"
	lastFailureAt       time.Time
}

const (
	degradedThreshold     = 3
	duplicateWindow       = 8
	maxRejectionRetries   = 2
	argsHashTruncateBytes = 16
)

// NewCorrectionLayer creates a fresh correction layer for one run.
func NewCorrectionLayer() *CorrectionLayer {
	return &CorrectionLayer{
		historyCap:     100,
		toolHealth:     make(map[string]*toolHealthEntry),
		toolFamilySeen: make(map[string]bool),
		evidenceSeen:   make(map[string]int),
	}
}

// BeforeToolCall returns a soft pre-tool hint. It never blocks execution.
// The hint is appended to the model's context for self-correction.
func (c *CorrectionLayer) BeforeToolCall(tool string, args map[string]any) string {
	c.mu.Lock()
	defer c.mu.Unlock()

	var hints []string

	// Check for duplicate calls within the window
	argsHash := hashArgs(args)
	repeated := c.countRecentCalls(tool, argsHash, duplicateWindow)
	if repeated >= 2 {
		hints = append(hints, fmt.Sprintf(
			"Exact tool call %s has already appeared %d times recently; "+
				"recent repetition may be low-value unless there is a new reason.",
			tool, repeated))
	}

	// Check for evidence_view redundancy
	if tool == "evidence_view" {
		if id, ok := args["id"].(string); ok {
			if c.evidenceSeen[id] > 0 {
				hints = append(hints, fmt.Sprintf(
					"evidence_view %s was already called %d times; "+
						"the content is already in context.",
					id, c.evidenceSeen[id]))
			}
		}
	}

	// Check degraded tool health
	if h, ok := c.toolHealth[tool]; ok && h.status == "degraded" {
		hints = append(hints, fmt.Sprintf(
			"%s is currently degraded after %d consecutive failure(s); "+
				"account for this tool-health signal when choosing actions.",
			tool, h.consecutiveFailures))
	}

	return strings.Join(hints, " ")
}

// AfterToolCall records post-tool health/progress facts and returns model-visible hints.
func (c *CorrectionLayer) AfterToolCall(tool string, args map[string]any, output string, isErr bool, duration time.Duration) string {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Record call in history
	argsHash := hashArgs(args)
	c.callHistory = append(c.callHistory, callEntry{
		tool: tool, argsHash: argsHash, timestamp: time.Now(),
	})
	if len(c.callHistory) > c.historyCap {
		c.callHistory = c.callHistory[len(c.callHistory)-c.historyCap:]
	}

	// Record evidence views
	if tool == "evidence_view" {
		if id, ok := args["id"].(string); ok {
			c.evidenceSeen[id]++
		}
	}

	// Update tool health
	h, ok := c.toolHealth[tool]
	if !ok {
		h = &toolHealthEntry{status: "healthy"}
		c.toolHealth[tool] = h
	}

	// Determine success/failure from output markers
	failed := isErr || containsFailureMarker(output)
	if failed {
		h.consecutiveFailures++
		h.lastFailureAt = time.Now()
		if h.consecutiveFailures >= degradedThreshold {
			h.status = "degraded"
		}
	} else {
		h.consecutiveFailures = 0
		h.status = "healthy"
	}

	// Track tool family
	c.toolFamilySeen[toolFamily(tool)] = true

	// Build hints
	var hints []string
	if h.status == "degraded" {
		hints = append(hints, fmt.Sprintf(
			"%s failed again (%d consecutive failures). Consider an alternative approach.",
			tool, h.consecutiveFailures))
	}

	// Stall detection: model is only reading evidence, not trying new tools
	if c.isStalled() {
		hints = append(hints,
			"Stall detected: only evidence re-reading observed recently. "+
				"Try a new tool or attack vector to make progress.")
	}

	return strings.Join(hints, " ")
}

// countRecentCalls counts how many times this exact tool+args appeared in the last N calls.
func (c *CorrectionLayer) countRecentCalls(tool, argsHash string, window int) int {
	count := 0
	start := len(c.callHistory) - window
	if start < 0 {
		start = 0
	}
	for i := start; i < len(c.callHistory); i++ {
		if c.callHistory[i].tool == tool && c.callHistory[i].argsHash == argsHash {
			count++
		}
	}
	return count
}

// isStalled detects when the model is only re-reading evidence without new tool calls.
func (c *CorrectionLayer) isStalled() bool {
	if len(c.callHistory) < 4 {
		return false
	}
	// Check last 4 calls
	recent := c.callHistory[len(c.callHistory)-4:]
	for _, e := range recent {
		fam := toolFamily(e.tool)
		if fam != "evidence" && fam != "skill" {
			return false // Found an active tool call
		}
	}
	// Also check that we've seen at least two families before (active + evidence)
	return len(c.toolFamilySeen) > 1
}

// ValidateCompletion checks if a claimed finding/flag/verified exploit is
// backed by evidence in the store. Returns (ok, feedback).
// After maxRejectionRetries rejections, it accepts with an "unverified" flag.
func (c *CorrectionLayer) ValidateCompletion(claim string, evidence *EvidenceStore) (bool, string) {
	c.mu.Lock()
	c.rejectionCount++
	rejections := c.rejectionCount
	c.mu.Unlock()

	if evidence == nil {
		return true, "" // No store available, can't validate
	}

	// Check if claimed artifacts appear in stored evidence
	if c.claimBackedByEvidence(claim, evidence) {
		return true, ""
	}

	// After max retries, accept with unverified flag
	if rejections > maxRejectionRetries {
		return true, "ACCEPTED WITH UNVERIFIED FLAG: claim could not be verified against stored evidence after " +
			fmt.Sprintf("%d attempts.", rejections)
	}

	// Reject with feedback
	return false, fmt.Sprintf(
		"REJECTION %d/%d: The claimed result could not be verified against stored evidence. "+
			"Ensure the claimed artifact (flag, response, output) actually appears in a tool result. "+
			"Use evidence_search to verify, then re-state your conclusion with evidence_id references.",
		rejections, maxRejectionRetries)
}

// claimBackedByEvidence checks if key tokens from the claim appear in evidence.
func (c *CorrectionLayer) claimBackedByEvidence(claim string, evidence *EvidenceStore) bool {
	// Extract flag-like tokens
	flags := ExtractFlags(claim)
	if len(flags) > 0 {
		for _, f := range flags {
			results := evidence.Search(f, 5)
			if len(results) > 0 {
				return true
			}
		}
	}

	// Extract URL-like tokens
	urls := extractURLs(claim)
	for _, u := range urls {
		results := evidence.Search(u, 3)
		if len(results) > 0 {
			return true
		}
	}

	// Check for common evidence references
	if strings.Contains(claim, "evidence_id:") || strings.Contains(claim, "evidence:") {
		return true
	}

	// Check for common completion keywords with evidence backing
	claimLower := strings.ToLower(claim)
	if strings.Contains(claimLower, "flag{") || strings.Contains(claimLower, "ctf{") {
		flags := ExtractFlags(claim)
		for _, f := range flags {
			results := evidence.Search(f, 3)
			if len(results) > 0 {
				return true
			}
		}
		return false
	}

	// For non-flag claims, accept if there are at least some evidence records
	if evidence.Count() > 0 {
		return true
	}

	return false
}

// --- Helpers ---

var failureMarkers = []string{
	"[!]", "failed locally", "traceback", "exception", "timed out", "timeout",
	"connection refused", "cancellederror", "constraint_violation",
	"role_tool_violation", "error:", "failed:", "fatal:",
}

func containsFailureMarker(output string) bool {
	lower := strings.ToLower(output)
	for _, m := range failureMarkers {
		if strings.Contains(lower, m) {
			return true
		}
	}
	return false
}

func toolFamily(toolName string) string {
	switch {
	case strings.HasPrefix(toolName, "evidence_"):
		return "evidence"
	case strings.HasPrefix(toolName, "skill_"):
		return "skill"
	case strings.HasPrefix(toolName, "report_") || strings.HasPrefix(toolName, "triage_"):
		return "finding"
	default:
		return "active"
	}
}

func hashArgs(args map[string]any) string {
	// Simple deterministic hash of args
	keys := make([]string, 0, len(args))
	for k := range args {
		keys = append(keys, k)
	}
	// Sort keys for determinism
	for i := 0; i < len(keys); i++ {
		for j := i + 1; j < len(keys); j++ {
			if keys[i] > keys[j] {
				keys[i], keys[j] = keys[j], keys[i]
			}
		}
	}
	var sb strings.Builder
	for _, k := range keys {
		sb.WriteString(k)
		sb.WriteString("=")
		sb.WriteString(fmt.Sprintf("%v", args[k]))
		sb.WriteString(";")
	}
	h := sha256.Sum256([]byte(sb.String()))
	return hex.EncodeToString(h[:argsHashTruncateBytes])
}

var urlRe = regexp.MustCompile(`https?://[^\s'"<>)]+`)

func extractURLs(text string) []string {
	matches := urlRe.FindAllString(text, -1)
	seen := make(map[string]bool)
	var out []string
	for _, m := range matches {
		if !seen[m] {
			seen[m] = true
			out = append(out, m)
		}
	}
	return out
}
