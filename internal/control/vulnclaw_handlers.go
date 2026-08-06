package control

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/anubhavg-icpl/talon/internal/core"
)

// handleRunEvidence returns the evidence store for a completed or in-progress run.
// Since EvidenceStore is run-scoped (in-memory), this returns the tool log
// as evidence-equivalent records for completed runs.
func (s *Server) handleRunEvidence(w http.ResponseWriter, r *http.Request) {
	runID := r.PathValue("run_id")
	if runID == "" {
		writeError(w, http.StatusBadRequest, "run_id required")
		return
	}
	sess, ok := s.store.Get(runID)
	if !ok {
		writeError(w, http.StatusNotFound, "run not found")
		return
	}

	type evidenceItem struct {
		Index   int    `json:"index"`
		Tool    string `json:"tool"`
		Summary string `json:"summary"`
		Size    int    `json:"size"`
	}
	items := make([]evidenceItem, 0, len(sess.ToolLog))
	for _, tc := range sess.ToolLog {
		summary := tc.Output
		if len(summary) > 200 {
			summary = summary[:200] + "..."
		}
		items = append(items, evidenceItem{
			Index: tc.Index, Tool: tc.ToolName, Summary: summary, Size: len(tc.Output),
		})
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"run_id":  runID,
		"total":   len(items),
		"items":   items,
	})
}

// handleTargetState returns the persisted target state for a given target address.
func (s *Server) handleTargetState(w http.ResponseWriter, r *http.Request) {
	addr := r.PathValue("addr")
	if addr == "" {
		writeError(w, http.StatusBadRequest, "addr required")
		return
	}
	ts := core.NewTargetStore("talon-data/targets")
	state, err := ts.GetOrCreate(addr)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, state)
}

// handleTargetResumePlan returns a deterministic resume plan for a target.
func (s *Server) handleTargetResumePlan(w http.ResponseWriter, r *http.Request) {
	addr := r.PathValue("addr")
	if addr == "" {
		writeError(w, http.StatusBadRequest, "addr required")
		return
	}
	ts := core.NewTargetStore("talon-data/targets")
	state, err := ts.GetOrCreate(addr)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	steps := state.ResumePlan()
	writeJSON(w, http.StatusOK, map[string]any{
		"target":  addr,
		"steps":   steps,
		"total":   len(steps),
		"summary": state.Summary(),
	})
}

// handleRunTraffic returns stored HTTP traffic for a run.
func (s *Server) handleRunTraffic(w http.ResponseWriter, r *http.Request) {
	runID := r.PathValue("run_id")
	if runID == "" {
		writeError(w, http.StatusBadRequest, "run_id required")
		return
	}
	sess, ok := s.store.Get(runID)
	if !ok {
		writeError(w, http.StatusNotFound, "run not found")
		return
	}

	// Derive traffic from tool log entries that look like HTTP calls
	type trafficItem struct {
		Seq    int    `json:"seq"`
		Tool   string `json:"tool"`
		Output string `json:"output_snippet"`
	}
	items := make([]trafficItem, 0)
	for _, tc := range sess.ToolLog {
		if isHTTPTool(tc.ToolName) {
			snippet := tc.Output
			if len(snippet) > 500 {
				snippet = snippet[:500] + "..."
			}
			items = append(items, trafficItem{
				Seq: tc.Index, Tool: tc.ToolName, Output: snippet,
			})
		}
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"run_id": runID,
		"total":  len(items),
		"items":  items,
	})
}

// handleRunRecap returns a deterministic recap for a completed run.
func (s *Server) handleRunRecap(w http.ResponseWriter, r *http.Request) {
	runID := r.PathValue("run_id")
	if runID == "" {
		writeError(w, http.StatusBadRequest, "run_id required")
		return
	}
	sess, ok := s.store.Get(runID)
	if !ok {
		writeError(w, http.StatusNotFound, "run not found")
		return
	}

	target := sess.RunInput.TargetIP
	findings := sess.Findings
	startTime := sess.StartedAt
	if startTime.IsZero() {
		startTime = time.Now().Add(-time.Minute)
	}

	tracker := &core.ToolCallTracker{
		Calls: sess.ToolLog,
	}
	recap := core.BuildRecap(target, runID, startTime, tracker, nil, findings)

	format := r.URL.Query().Get("format")
	if format == "markdown" || format == "md" {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(recap.FormatMarkdown()))
		return
	}
	writeJSON(w, http.StatusOK, recap)
}

// handleCryptoOps returns the list of supported crypto operations.
func (s *Server) handleCryptoOps(w http.ResponseWriter, r *http.Request) {
	// Import the operation list from core
	writeJSON(w, http.StatusOK, map[string]any{
		"operations": cryptoOpList(),
	})
}

// handleCryptoDecode executes a crypto operation server-side.
func (s *Server) handleCryptoDecode(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Operation string `json:"operation"`
		Input     string `json:"input"`
		Key       string `json:"key"`
		IV        string `json:"iv"`
		Shift     int    `json:"shift"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if req.Operation == "" || req.Input == "" {
		writeError(w, http.StatusBadRequest, "operation and input required")
		return
	}
	if req.Shift == 0 {
		req.Shift = 3
	}

	result, err := cryptoExecuteSafe(req.Operation, req.Input, req.Key, req.IV, req.Shift)
	if err != nil {
		writeJSON(w, http.StatusOK, map[string]any{
			"success": false,
			"error":   err.Error(),
		})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"success":   true,
		"operation": req.Operation,
		"result":    result,
	})
}

func isHTTPTool(name string) bool {
	switch name {
	case "fetch", "http_probe_batch", "web_headers_audit", "js_endpoint_extract",
		"sqlmap_scan", "nuclei_scan", "curl":
		return true
	}
	return false
}

// cryptoOpList and cryptoExecuteSafe are thin wrappers that call the core
// package functions. They exist to avoid import cycles in the handler.
func cryptoOpList() []string {
	// We expose the operations via a hardcoded list that mirrors core.cryptoOpNames()
	// to avoid needing to export that function from core.
	return []string{
		"base64_encode", "base64_decode", "base32_encode", "base32_decode",
		"base58_encode", "base58_decode", "hex_encode", "hex_decode",
		"url_encode", "url_decode", "html_encode", "html_decode",
		"unicode_encode", "unicode_decode", "rot13", "caesar_encode", "caesar_decode",
		"morse_encode", "morse_decode",
		"md5_hash", "sha1_hash", "sha256_hash", "sha512_hash",
		"aes_encrypt", "aes_decrypt", "des_encrypt", "des_decrypt",
		"jwt_decode", "auto_decode",
	}
}

// cryptoExecuteSafe wraps core.CryptoExecute (which is unexported) via a
// re-export pattern. Since we can't call the unexported function directly,
// we use a small public shim.
func cryptoExecuteSafe(op, input, key, iv string, shift int) (string, error) {
	return core.CryptoExecutePublic(op, input, key, iv, shift)
}
