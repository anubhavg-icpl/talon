package control

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/anubhavg-icpl/talon/internal/llm"
)

// analyzeTimeout bounds the one-shot LLM analysis call.
const analyzeTimeout = 120 * time.Second

// analysisSystemPrompt steers the one-shot analyst model.
const analysisSystemPrompt = `You are a senior penetration-test analyst embedded in an ops console.
Given one agent run's metadata, final report, and tool-call log, produce a tight operator briefing:

1. EXECUTIVE SUMMARY — 2-3 sentences: target, outcome, confidence.
2. KILL CHAIN — what actually worked, step by step (cite tool names).
3. FINDINGS — vulnerabilities confirmed, with evidence (command output excerpts).
4. IMPACT — what an attacker gains (be concrete: shells, creds, data).
5. REMEDIATION — prioritized, actionable fixes.
6. FALSE-POSITIVE CHECK — anything unproven or contradictory, say so plainly.

Use plain text with short uppercase section headers. Be terse, technical, and honest — PoC || GTFO.`

// handleAnalyze is POST /analyze/{run_id} — one-shot LLM analysis of a run's
// report + tool log, for the dashboard's AI ANALYSIS panel. Requires an
// analyzer model (the server's configured LLM); 503 when unconfigured.
func (s *Server) handleAnalyze(w http.ResponseWriter, r *http.Request) {
	if s.analyzer == nil {
		writeError(w, http.StatusServiceUnavailable, "analysis unavailable: no LLM configured on talon-core")
		return
	}
	if s.settings != nil && !s.settings.FeatureEnabled(r.Context(), "FEATURE_AI_ANALYSIS", true) {
		writeError(w, http.StatusForbidden, "AI analysis is disabled in settings")
		return
	}
	runID := r.PathValue("run_id")

	// Cached briefing? The LLM call costs tens of seconds — serve Redis hits
	// instantly (keyed per run; long TTL since terminal runs don't change).
	if s.cache != nil {
		if cached, ok := s.cache.Get(r.Context(), cacheKeyAnalyzePrefix+runID); ok {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("X-Cache", "hit")
			_, _ = w.Write([]byte(cached))
			return
		}
	}

	sess, ok := s.store.Get(runID)
	if !ok {
		writeError(w, http.StatusNotFound, "Run not found")
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), analyzeTimeout)
	defer cancel()

	resp, err := s.analyzer.Converse(ctx, analysisSystemPrompt, []llm.Message{
		llm.UserMessage(buildAnalysisPrompt(sess)),
	}, nil)
	if err != nil {
		writeError(w, http.StatusBadGateway, fmt.Sprintf("analysis failed: %v", err))
		return
	}
	if s.cache != nil {
		if raw, err := json.Marshal(map[string]string{"run_id": runID, "analysis": resp.Text}); err == nil {
			s.cache.Set(r.Context(), cacheKeyAnalyzePrefix+runID, string(raw), analyzeTTL)
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("X-Cache", "miss")
			_, _ = w.Write(raw)
			return
		}
	}
	writeJSON(w, http.StatusOK, map[string]string{"run_id": runID, "analysis": resp.Text})
}

// buildAnalysisPrompt renders a session into a compact analyst input,
// truncating tool outputs so the prompt stays within model budgets.
func buildAnalysisPrompt(sess Session) string {
	var b strings.Builder
	fmt.Fprintf(&b, "TARGET: %s\n", sess.RunInput.TargetIP)
	if sess.RunInput.CVEID != "" {
		fmt.Fprintf(&b, "CVE: %s\n", sess.RunInput.CVEID)
	}
	if sess.RunInput.ServiceName != "" {
		fmt.Fprintf(&b, "SERVICE: %s\n", sess.RunInput.ServiceName)
	}
	if sess.RunInput.Description != "" {
		fmt.Fprintf(&b, "DESCRIPTION: %s\n", sess.RunInput.Description)
	}
	fmt.Fprintf(&b, "STATUS: %s\n", sess.Status)
	if sess.JudgeSet {
		fmt.Fprintf(&b, "JUDGE VERDICT (compromise proven): %v\n", sess.JudgeVerdict)
	}
	if sess.Output != "" {
		fmt.Fprintf(&b, "\nFINAL REPORT:\n%s\n", truncate(sess.Output, 4000))
	}
	b.WriteString("\nTOOL LOG:\n")
	for _, rec := range sess.ToolLog {
		args, _ := json.Marshal(rec.Args)
		fmt.Fprintf(&b, "\n[%d] %s %s\n%s\n", rec.Index, rec.ToolName, truncate(string(args), 500), truncate(rec.Output, 1200))
	}
	return b.String()
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "\n…[truncated]"
}
