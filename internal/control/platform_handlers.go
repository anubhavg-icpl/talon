package control

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/anubhavg-icpl/talon/internal/config"
	"github.com/anubhavg-icpl/talon/internal/core"
	"github.com/google/uuid"
)

func (s *Server) ensurePlatform() *Platform {
	if s.platform == nil {
		// lazy empty platform if main forgot to wire
		s.platform = NewPlatform("talon-data")
	}
	return s.platform
}

func (s *Server) handleGetScope(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.ensurePlatform().GetScope())
}

func (s *Server) handlePutScope(w http.ResponseWriter, r *http.Request) {
	var sc ScopePolicy
	if err := json.NewDecoder(r.Body).Decode(&sc); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	writeJSON(w, http.StatusOK, s.ensurePlatform().PutScope(sc))
}

func (s *Server) handleListTargets(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"targets": s.ensurePlatform().ListTargets()})
}

func (s *Server) handleUpsertTarget(w http.ResponseWriter, r *http.Request) {
	var t Target
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if t.Address == "" && t.URL == "" {
		writeError(w, http.StatusBadRequest, "address or url required")
		return
	}
	writeJSON(w, http.StatusOK, s.ensurePlatform().UpsertTarget(t))
}

func (s *Server) handleDeleteTarget(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if !s.ensurePlatform().DeleteTarget(id) {
		writeError(w, http.StatusNotFound, "target not found")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"deleted": id})
}

func (s *Server) handleListSchedules(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"schedules": s.ensurePlatform().ListSchedules()})
}

func (s *Server) handleUpsertSchedule(w http.ResponseWriter, r *http.Request) {
	var sch Schedule
	if err := json.NewDecoder(r.Body).Decode(&sch); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if sch.Name == "" || sch.TargetAddr == "" {
		writeError(w, http.StatusBadRequest, "name and target required")
		return
	}
	writeJSON(w, http.StatusOK, s.ensurePlatform().UpsertSchedule(sch))
}

func (s *Server) handleDeleteSchedule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if !s.ensurePlatform().DeleteSchedule(id) {
		writeError(w, http.StatusNotFound, "schedule not found")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"deleted": id})
}

func (s *Server) handleGetNotify(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.ensurePlatform().GetNotify())
}

func (s *Server) handlePutNotify(w http.ResponseWriter, r *http.Request) {
	var n NotifyConfig
	if err := json.NewDecoder(r.Body).Decode(&n); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	writeJSON(w, http.StatusOK, s.ensurePlatform().PutNotify(n))
}

func (s *Server) handleListCredentials(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"credentials": s.ensurePlatform().ListCredentials()})
}

func (s *Server) handleAddCredential(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Name     string `json:"name"`
		Kind     string `json:"kind"`
		Username string `json:"username"`
		Secret   string `json:"secret"`
		Scope    string `json:"scope"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if body.Name == "" || body.Secret == "" {
		writeError(w, http.StatusBadRequest, "name and secret required")
		return
	}
	if body.Kind == "" {
		body.Kind = "password"
	}
	c, err := s.ensurePlatform().AddCredential(body.Name, body.Kind, body.Username, body.Secret, body.Scope)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, c)
}

func (s *Server) handleDeleteCredential(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if !s.ensurePlatform().DeleteCredential(id) {
		writeError(w, http.StatusNotFound, "not found")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"deleted": id})
}

func (s *Server) handleListEvidence(w http.ResponseWriter, r *http.Request) {
	runID := r.URL.Query().Get("run_id")
	writeJSON(w, http.StatusOK, map[string]any{"evidence": s.ensurePlatform().ListEvidence(runID)})
}

func (s *Server) handleAddEvidence(w http.ResponseWriter, r *http.Request) {
	var e EvidenceItem
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if e.RunID == "" || e.Title == "" {
		writeError(w, http.StatusBadRequest, "run_id and title required")
		return
	}
	if e.Kind == "" {
		e.Kind = "note"
	}
	writeJSON(w, http.StatusOK, s.ensurePlatform().AddEvidence(e))
}

func (s *Server) handleBudget(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.ensurePlatform().GetBudget())
}

// handleRetest starts a new run seeded from a previous run's finding.
func (s *Server) handleRetest(w http.ResponseWriter, r *http.Request) {
	runID := r.PathValue("run_id")
	sess, ok := s.store.Get(runID)
	if !ok {
		writeError(w, http.StatusNotFound, "run not found")
		return
	}
	var body struct {
		FindingID string `json:"finding_id"`
	}
	_ = json.NewDecoder(r.Body).Decode(&body)

	desc := "RETEST AUTHORIZED — verify prior finding is remediated. "
	if body.FindingID != "" {
		for _, f := range sess.Findings {
			if f.ID == body.FindingID {
				desc += fmt.Sprintf("Prior finding [%s] %s: %s. Endpoint=%s. Confirm FIXED or still VULNERABLE with 3-gate evidence.",
					f.Severity, f.Title, f.Description, f.Endpoint)
				break
			}
		}
	} else {
		desc += "Re-validate entire prior engagement scope against same target."
	}

	active, _, _ := s.store.PaginatedList(500, 0)
	activeN := 0
	for _, rs := range active {
		if rs.Status == "running" || rs.Status == "awaiting_approval" || rs.Status == "initializing" {
			activeN++
		}
	}
	if err := s.ensurePlatform().CheckStart(sess.RunInput.TargetIP, desc, activeN); err != nil {
		writeError(w, http.StatusForbidden, err.Error())
		return
	}

	newID := uuid.NewString()
	input := sess.RunInput
	input.SessionID = newID
	input.Description = desc
	input.AgentMode = core.NormalizeAgentMode(input.AgentMode)
	if input.AgentMode == "" {
		input.AgentMode = core.AgentModeFull
	}
	s.store.Create(newID, input)
	s.ensurePlatform().IncBudget("start", 1)
	go s.runWorkflow(newID, input)
	writeJSON(w, http.StatusOK, map[string]string{
		"run_id":  newID,
		"message": "Retest started",
		"source":  runID,
	})
}

// handleReportHTML returns a print-friendly HTML client report (browser → PDF).
func (s *Server) handleReportHTML(w http.ResponseWriter, r *http.Request) {
	runID := r.PathValue("run_id")
	sess, ok := s.store.Get(runID)
	if !ok {
		writeError(w, http.StatusNotFound, "run not found")
		return
	}
	md := ""
	if sess.Report != nil {
		md = sess.Report.Markdown
	}
	if md == "" {
		md = sess.Output
	}
	// Minimal HTML with print CSS — open in browser and Print → PDF
	html := fmt.Sprintf(`<!DOCTYPE html>
<html><head><meta charset="utf-8"><title>Talon Report %s</title>
<style>
body{font-family:ui-monospace,Menlo,monospace;background:#0a0608;color:#f2e8ea;padding:2rem;max-width:900px;margin:0 auto;line-height:1.5}
h1,h2,h3{color:#ff2b2b} pre{white-space:pre-wrap;background:#140c10;padding:1rem;border:1px solid #3a1a22}
.meta{color:#a8989c;font-size:12px;margin-bottom:2rem}
@media print{body{background:#fff;color:#111} h1,h2,h3{color:#b00} pre{border-color:#ccc}}
</style></head><body>
<div class="meta">Talon AI · Authorized engagement report · Run %s · Target %s</div>
<pre>%s</pre>
<script>/* optional: window.print() */</script>
</body></html>`, runID[:8], runID, htmlEscape(sess.RunInput.TargetIP), htmlEscape(md))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(html))
}

func htmlEscape(s string) string {
	r := strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;", `"`, "&quot;")
	return r.Replace(s)
}

// handleOpenAPI serves a minimal OpenAPI 3 document for the control plane.
func (s *Server) handleOpenAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/yaml")
	_, _ = w.Write([]byte(openAPIYAML))
}

// startRunFromPlatform is used by scheduler — reuses store + workflow.
func (s *Server) startRunFromPlatform(input core.RunInput) string {
	if input.SessionID == "" {
		input.SessionID = uuid.NewString()
	}
	if input.Context.LHOST == "" {
		input.Context = config.Context{LHOST: "127.0.0.1", LPORT: 4444}
	}
	// scope check
	active, _, _ := s.store.PaginatedList(500, 0)
	n := 0
	for _, rs := range active {
		if rs.Status == "running" || rs.Status == "awaiting_approval" {
			n++
		}
	}
	if err := s.ensurePlatform().CheckStart(input.TargetIP, input.Description, n); err != nil {
		return ""
	}
	s.store.Create(input.SessionID, input)
	s.ensurePlatform().IncBudget("start", 1)
	go s.runWorkflow(input.SessionID, input)
	return input.SessionID
}

const openAPIYAML = `openapi: 3.0.3
info:
  title: Talon Control Plane
  version: "1.0"
  description: AI pentest orchestration API
servers:
  - url: http://localhost:8000
paths:
  /health:
    get:
      summary: Liveness
  /input/start:
    post:
      summary: Start a validation run
  /input/batch:
    post:
      summary: Start runs for multiple IPs
  /runs:
    get:
      summary: List runs
  /runs/{run_id}/findings:
    get:
      summary: Structured findings
  /runs/{run_id}/report:
    get:
      summary: Structured report
  /runs/{run_id}/export:
    get:
      summary: JSON export bundle
  /runs/{run_id}/report.html:
    get:
      summary: Printable HTML report (Save as PDF)
  /runs/{run_id}/retest:
    post:
      summary: Launch retest from prior run
  /skills:
    get:
      summary: CyberStrike skill catalog
  /playbooks:
    get:
      summary: Engagement playbooks
  /targets:
    get:
      summary: Target inventory
    post:
      summary: Upsert target
  /scope:
    get:
      summary: Scope / ROE policy
    put:
      summary: Update scope policy
  /schedules:
    get:
      summary: List schedules
    post:
      summary: Upsert schedule
  /notify:
    get:
      summary: Notification config
    put:
      summary: Update webhooks
  /credentials:
    get:
      summary: List credentials (no secrets)
    post:
      summary: Add encrypted credential
  /evidence:
    get:
      summary: Evidence vault
    post:
      summary: Add evidence
  /budget:
    get:
      summary: Token/tool budget counters
  /intel:
    get:
      summary: Cross-run intel feed
  /openapi.yaml:
    get:
      summary: This document
`
