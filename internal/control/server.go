package control

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/anubhavg-icpl/talon/internal/config"
	"github.com/anubhavg-icpl/talon/internal/core"
	"github.com/anubhavg-icpl/talon/internal/llm"
	"github.com/anubhavg-icpl/talon/internal/mcpclient"
	"github.com/coder/websocket"
	"github.com/google/uuid"
)

// defaultRunTimeout is the wall-clock budget for one start/resume workflow
// segment. Override with TALON_RUN_TIMEOUT (Go duration, e.g. "20m").
const defaultRunTimeout = 20 * time.Minute

func runTimeout() time.Duration {
	if v := os.Getenv("TALON_RUN_TIMEOUT"); v != "" {
		if d, err := time.ParseDuration(v); err == nil && d > 0 {
			return d
		}
	}
	return defaultRunTimeout
}

// Server wires the Store to an core.Orchestrator over HTTP routes:
// /input/start, /runs, /monitor/traces/{run_id}, /monitor/tools,
// /monitor/stream/{run_id} (SSE), /monitor/ws/{run_id} (WebSocket),
// /output/status/{run_id}, /output/resume/{run_id}.
type Server struct {
	orch  *core.Orchestrator
	store *Store
	// analyzer, when non-nil, powers POST /analyze/{run_id} — typically the
	// same LLM the orchestrator uses.
	analyzer llm.ChatModel
	// auth, when non-nil, gates every route except /health* and /auth/login.
	auth *Auth
	// db, when non-nil, lets /health/services probe postgres directly.
	db *DB
	// cache, when non-nil, accelerates health probes, analysis and sessions.
	cache Cache
	// settings resolves dashboard-managed config (nil-safe env fallback).
	settings *Settings
	// tools, when non-nil, backs GET /mcp/servers.
	tools *mcpclient.Multi
	// platform: targets, scope, schedules, notify, credentials, evidence, budget.
	platform *Platform
}

func NewServer(orch *core.Orchestrator, store *Store, opts ...ServerOption) *Server {
	s := &Server{orch: orch, store: store}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// ServerOption customizes a Server (keeps NewServer call sites stable).
type ServerOption func(*Server)

// WithAnalyzer enables POST /analyze/{run_id} with the given model.
func WithAnalyzer(m llm.ChatModel) ServerOption {
	return func(s *Server) { s.analyzer = m }
}

// WithAuth enables session auth on all non-open routes.
func WithAuth(a *Auth) ServerOption {
	return func(s *Server) { s.auth = a }
}

// WithDB lets /health/services probe postgres.
func WithDB(db *DB) ServerOption {
	return func(s *Server) { s.db = db }
}

// WithCache enables Redis acceleration (health probes, analysis, sessions).
func WithCache(c Cache) ServerOption {
	return func(s *Server) { s.cache = c }
}

// WithSettings enables GET/PUT /config.
func WithSettings(st *Settings) ServerOption {
	return func(s *Server) { s.settings = st }
}

// WithTools enables GET /mcp/servers.
func WithTools(m *mcpclient.Multi) ServerOption {
	return func(s *Server) { s.tools = m }
}

// WithPlatform wires targets/scope/schedules/notify/credentials.
func WithPlatform(p *Platform) ServerOption {
	return func(s *Server) {
		s.platform = p
		if p != nil {
			p.SetStartFn(s.startRunFromPlatform)
			p.StartScheduler()
		}
	}
}

// Mux builds the http.ServeMux, using Go 1.22+ method+pattern routing so no
// router dependency is needed.
func (s *Server) Mux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", s.handleHealth)
	mux.HandleFunc("GET /health/services", s.handleServiceHealth)
	mux.HandleFunc("POST /auth/login", s.handleLogin)
	mux.HandleFunc("POST /auth/logout", s.handleLogout)
	mux.HandleFunc("GET /auth/me", s.handleMe)
	mux.HandleFunc("POST /input/start", s.handleStart)
	mux.HandleFunc("GET /runs", s.handleListRuns)
	mux.HandleFunc("GET /runs/summary", s.handleRunsSummary)
	mux.HandleFunc("GET /monitor/traces/{run_id}", s.handleTraces)
	mux.HandleFunc("GET /monitor/tools", s.handleToolLog)
	mux.HandleFunc("GET /monitor/stream/{run_id}", s.handleStream)
	mux.HandleFunc("GET /monitor/ws/{run_id}", s.handleStreamWS)
	mux.HandleFunc("GET /output/status/{run_id}", s.handleStatus)
	mux.HandleFunc("POST /output/resume/{run_id}", s.handleResume)
	mux.HandleFunc("POST /analyze/{run_id}", s.handleAnalyze)
	mux.HandleFunc("GET /config", s.handleGetConfig)
	mux.HandleFunc("PUT /config", s.handlePutConfig)
	mux.HandleFunc("GET /mcp/servers", s.handleMCPServers)

	// /shell/* — SSO reverse-proxy to the ttyd web terminal inside
	// arsenal_engine. ttyd runs no-auth bound to loopback; this route is gated
	// by the same session-auth middleware as the rest of the console (the
	// talon_session cookie — no second login), and Go proxies the WebSocket
	// upgrade natively. Only reachable through this authed proxy.
	shellTarget := os.Getenv("TTYD_URL")
	if shellTarget == "" {
		shellTarget = "http://127.0.0.1:7681"
	}
	if u, err := url.Parse(shellTarget); err == nil {
		proxy := httputil.NewSingleHostReverseProxy(u)
		mux.Handle("/shell/", proxy)
		mux.HandleFunc("GET /shell", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/shell/", http.StatusTemporaryRedirect)
		})
	}

	// Structured findings + report (CyberStrike-inspired) + skills catalog.
	mux.HandleFunc("GET /runs/{run_id}/findings", s.handleFindings)
	mux.HandleFunc("POST /runs/{run_id}/findings/{finding_id}/triage", s.handleTriageFinding)
	mux.HandleFunc("GET /runs/{run_id}/report", s.handleReport)
	mux.HandleFunc("GET /runs/{run_id}/killchain", s.handleKillChain)
	mux.HandleFunc("GET /runs/{run_id}/methodology", s.handleMethodology)
	mux.HandleFunc("GET /findings", s.handleGlobalFindings)
	mux.HandleFunc("GET /skills", s.handleSkills)
	mux.HandleFunc("GET /skills/{id}", s.handleSkillByID)
	mux.HandleFunc("GET /agents", s.handleAgents)
	// Wave 6: playbooks, compare, export, notes, intel, timeline, batch
	mux.HandleFunc("GET /playbooks", s.handlePlaybooks)
	mux.HandleFunc("GET /intel", s.handleIntel)
	mux.HandleFunc("GET /runs/compare", s.handleCompareRuns)
	mux.HandleFunc("GET /runs/{run_id}/export", s.handleExport)
	mux.HandleFunc("GET /runs/{run_id}/timeline", s.handleTimeline)
	mux.HandleFunc("GET /runs/{run_id}/notes", s.handleGetNotes)
	mux.HandleFunc("POST /runs/{run_id}/notes", s.handleAddNote)
	mux.HandleFunc("POST /input/batch", s.handleBatchStart)
	// Platform: scope, targets, schedules, notify, credentials, evidence, budget, retest, HTML report, OpenAPI
	mux.HandleFunc("GET /scope", s.handleGetScope)
	mux.HandleFunc("PUT /scope", s.handlePutScope)
	mux.HandleFunc("GET /targets", s.handleListTargets)
	mux.HandleFunc("POST /targets", s.handleUpsertTarget)
	mux.HandleFunc("DELETE /targets/{id}", s.handleDeleteTarget)
	mux.HandleFunc("GET /schedules", s.handleListSchedules)
	mux.HandleFunc("POST /schedules", s.handleUpsertSchedule)
	mux.HandleFunc("DELETE /schedules/{id}", s.handleDeleteSchedule)
	mux.HandleFunc("GET /notify", s.handleGetNotify)
	mux.HandleFunc("PUT /notify", s.handlePutNotify)
	mux.HandleFunc("GET /credentials", s.handleListCredentials)
	mux.HandleFunc("POST /credentials", s.handleAddCredential)
	mux.HandleFunc("DELETE /credentials/{id}", s.handleDeleteCredential)
	mux.HandleFunc("GET /evidence", s.handleListEvidence)
	mux.HandleFunc("POST /evidence", s.handleAddEvidence)
	mux.HandleFunc("GET /budget", s.handleBudget)
	mux.HandleFunc("POST /runs/{run_id}/retest", s.handleRetest)
	mux.HandleFunc("GET /runs/{run_id}/report.html", s.handleReportHTML)
	mux.HandleFunc("GET /openapi.yaml", s.handleOpenAPI)
	return mux
}

// Handler wraps the mux with permissive CORS so the web dashboard (and other
// local tooling) can call the API cross-origin. Matches the existing no-auth
// local-operator posture.
func (s *Server) Handler() http.Handler {
	mux := s.Mux()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		// Session gate: everything except open routes requires a valid
		// cookie or bearer token when auth is enabled.
		if s.auth != nil && !openPath(r.URL.Path) {
			if s.auth.usernameForRequest(r) == "" {
				writeError(w, http.StatusUnauthorized, "unauthorized")
				return
			}
		}
		mux.ServeHTTP(w, r)
	})
}

// handleHealth is GET /health — liveness for operators and the talon CLI.
func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"service": "talon-core",
	})
}

// targetRequest is the POST /input/start request body.
type targetRequest struct {
	IP          string `json:"ip"`
	CVEID       string `json:"cve_id"`
	ServiceName string `json:"service_name"`
	Description string `json:"description"`
	LHOST       string `json:"lhost"`
	LPORT       int    `json:"lport"`
	// AgentMode: full|recon|exploit|web|network|post (CyberStrike-style specialist).
	AgentMode string `json:"agent_mode"`
	// PlaybookID optionally fills defaults from a builtin playbook.
	PlaybookID string `json:"playbook_id"`
}

// batchStartRequest is POST /input/batch — launch multiple hosts with shared context.
type batchStartRequest struct {
	IPs         []string `json:"ips"`
	CVEID       string   `json:"cve_id"`
	ServiceName string   `json:"service_name"`
	Description string   `json:"description"`
	LHOST       string   `json:"lhost"`
	LPORT       int      `json:"lport"`
	AgentMode   string   `json:"agent_mode"`
	PlaybookID  string   `json:"playbook_id"`
}

// resumeRequest is the POST /output/resume/{run_id} request body.
type resumeRequest struct {
	Decision   string         `json:"decision"`
	EditedArgs map[string]any `json:"edited_args"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, detail string) {
	writeJSON(w, status, map[string]string{"detail": detail})
}

// handleStart is POST /input/start.
func (s *Server) handleStart(w http.ResponseWriter, r *http.Request) {
	var req targetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	applyStartDefaults(&req)
	if strings.TrimSpace(req.IP) == "" {
		writeError(w, http.StatusBadRequest, "ip required")
		return
	}

	// Scope / ROE gate
	active, _, _ := s.store.PaginatedList(500, 0)
	activeN := 0
	for _, rs := range active {
		if rs.Status == "running" || rs.Status == "awaiting_approval" || rs.Status == "initializing" {
			activeN++
		}
	}
	if err := s.ensurePlatform().CheckStart(req.IP, req.Description, activeN); err != nil {
		writeError(w, http.StatusForbidden, err.Error())
		return
	}

	runID := uuid.NewString()
	input := core.RunInput{
		SessionID:   runID,
		TargetIP:    req.IP,
		CVEID:       req.CVEID,
		ServiceName: req.ServiceName,
		Description: req.Description,
		AgentMode:   core.NormalizeAgentMode(req.AgentMode),
		Context:     config.Context{LHOST: req.LHOST, LPORT: req.LPORT},
	}
	s.store.Create(runID, input)
	s.ensurePlatform().IncBudget("start", 1)
	s.ensurePlatform().TouchTarget(req.IP, runID, "running")

	// Runs in the background; on a HITL interrupt it stores the pending
	// interrupt and returns rather than blocking -- POST
	// /output/resume/{run_id} is what drives the workflow forward from
	// there (see handleResume).
	go s.runWorkflow(runID, input)

	writeJSON(w, http.StatusOK, map[string]string{
		"run_id":     runID,
		"message":    "Agent execution started",
		"agent_mode": input.AgentMode,
	})
}

func applyStartDefaults(req *targetRequest) {
	if req.PlaybookID != "" {
		if pb, ok := core.GetPlaybook(req.PlaybookID); ok {
			if req.AgentMode == "" {
				req.AgentMode = pb.AgentMode
			}
			if req.Description == "" {
				req.Description = pb.Prompt
			}
			if req.PlaybookID == "cve-lab" && req.CVEID == "" {
				req.CVEID = "CVE-2011-2523"
			}
			if req.PlaybookID == "cve-lab" && req.ServiceName == "" {
				req.ServiceName = "vsftpd 2.3.4"
			}
		}
	}
	if req.LHOST == "" {
		if e := os.Getenv("LHOST"); e != "" {
			req.LHOST = e
		} else {
			req.LHOST = "127.0.0.1"
		}
	}
	if req.LPORT == 0 {
		if e := os.Getenv("LPORT"); e != "" {
			if n, err := strconv.Atoi(e); err == nil && n > 0 {
				req.LPORT = n
			}
		}
		if req.LPORT == 0 {
			req.LPORT = 4444
		}
	}
}

// handleBatchStart is POST /input/batch — start one run per IP.
func (s *Server) handleBatchStart(w http.ResponseWriter, r *http.Request) {
	var req batchStartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if len(req.IPs) == 0 {
		writeError(w, http.StatusBadRequest, "ips required")
		return
	}
	if len(req.IPs) > 50 {
		writeError(w, http.StatusBadRequest, "max 50 ips per batch")
		return
	}
	type started struct {
		RunID string `json:"run_id"`
		IP    string `json:"ip"`
	}
	out := make([]started, 0, len(req.IPs))
	for _, ip := range req.IPs {
		ip = strings.TrimSpace(ip)
		if ip == "" {
			continue
		}
		tr := targetRequest{
			IP: ip, CVEID: req.CVEID, ServiceName: req.ServiceName,
			Description: req.Description, LHOST: req.LHOST, LPORT: req.LPORT,
			AgentMode: req.AgentMode, PlaybookID: req.PlaybookID,
		}
		applyStartDefaults(&tr)
		runID := uuid.NewString()
		input := core.RunInput{
			SessionID: runID, TargetIP: tr.IP, CVEID: tr.CVEID, ServiceName: tr.ServiceName,
			Description: tr.Description, AgentMode: core.NormalizeAgentMode(tr.AgentMode),
			Context: config.Context{LHOST: tr.LHOST, LPORT: tr.LPORT},
		}
		s.store.Create(runID, input)
		go s.runWorkflow(runID, input)
		out = append(out, started{RunID: runID, IP: ip})
	}
	writeJSON(w, http.StatusOK, map[string]any{"started": out, "count": len(out)})
}

// handlePlaybooks is GET /playbooks.
func (s *Server) handlePlaybooks(w http.ResponseWriter, r *http.Request) {
	pbs := core.ListPlaybooks()
	writeJSON(w, http.StatusOK, map[string]any{"playbooks": pbs, "count": len(pbs)})
}

// handleIntel is GET /intel — cross-run feed.
func (s *Server) handleIntel(w http.ResponseWriter, r *http.Request) {
	limit := 50
	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 200 {
			limit = n
		}
	}
	writeJSON(w, http.StatusOK, map[string]any{"events": s.store.IntelFeed(limit)})
}

// handleCompareRuns is GET /runs/compare?a=&b=.
func (s *Server) handleCompareRuns(w http.ResponseWriter, r *http.Request) {
	aID := r.URL.Query().Get("a")
	bID := r.URL.Query().Get("b")
	if aID == "" || bID == "" {
		writeError(w, http.StatusBadRequest, "query params a and b (run ids) required")
		return
	}
	sa, oka := s.store.Snapshot(aID)
	sb, okb := s.store.Snapshot(bID)
	if !oka || !okb {
		writeError(w, http.StatusNotFound, "one or both runs not found")
		return
	}
	writeJSON(w, http.StatusOK, core.CompareRuns(sa, sb))
}

// handleExport is GET /runs/{id}/export — portable JSON bundle.
func (s *Server) handleExport(w http.ResponseWriter, r *http.Request) {
	runID := r.PathValue("run_id")
	bundle, ok := s.store.ExportBundle(runID)
	if !ok {
		writeError(w, http.StatusNotFound, "Run not found")
		return
	}
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="talon-export-%s.json"`, runID[:8]))
	writeJSON(w, http.StatusOK, bundle)
}

// handleTimeline is GET /runs/{id}/timeline.
func (s *Server) handleTimeline(w http.ResponseWriter, r *http.Request) {
	runID := r.PathValue("run_id")
	sess, ok := s.store.Get(runID)
	if !ok {
		writeError(w, http.StatusNotFound, "Run not found")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"run_id":   runID,
		"timeline": core.BuildTimeline(sess.ToolLog, sess.Findings),
	})
}

// handleGetNotes is GET /runs/{id}/notes.
func (s *Server) handleGetNotes(w http.ResponseWriter, r *http.Request) {
	runID := r.PathValue("run_id")
	notes, ok := s.store.Notes(runID)
	if !ok {
		writeError(w, http.StatusNotFound, "Run not found")
		return
	}
	if notes == nil {
		notes = []core.OperatorNote{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"run_id": runID, "notes": notes})
}

// handleAddNote is POST /runs/{id}/notes.
func (s *Server) handleAddNote(w http.ResponseWriter, r *http.Request) {
	runID := r.PathValue("run_id")
	var body struct {
		Body   string `json:"body"`
		Author string `json:"author"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	if body.Author == "" {
		body.Author = "operator"
	}
	n, err := s.store.AddNote(runID, body.Author, body.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, n)
}

func (s *Server) runWorkflow(runID string, input core.RunInput) {
	ctx, cancel := context.WithTimeout(context.Background(), runTimeout())
	defer cancel()
	ctx = core.WithProgress(ctx, func(toolLog []core.ToolCallRecord) {
		s.store.SetToolLog(runID, toolLog)
	})
	ctx = core.WithFindingsProgress(ctx, func(findings []core.Finding) {
		s.store.SetFindings(runID, findings)
	})

	s.store.SetStatus(runID, "running")
	log.Printf("talon-core: run %s starting (timeout=%s target=%s)", runID, runTimeout(), input.TargetIP)
	result, err := s.orch.Run(ctx, input)
	if err != nil {
		log.Printf("talon-core: run %s: %v", runID, err)
		// If we timed out but still have tool output, surface it as completed
		// with a note rather than a bare error (better for operators).
		if (ctx.Err() != nil) && len(result.ToolLog) > 0 {
			result.FinalMessage = strings.TrimSpace(result.FinalMessage + "\n[run stopped: wall-clock timeout]")
			result.Interrupted = false
			s.store.SetResult(runID, ensureFindings(input, result))
			return
		}
		s.store.SetError(runID, err)
		return
	}
	log.Printf("talon-core: run %s finished interrupted=%v tools=%d findings=%d",
		runID, result.Interrupted, len(result.ToolLog), len(result.Findings))
	final := ensureFindings(input, result)
	s.store.SetResult(runID, final)
	s.afterRunHooks(runID, input, final)
}

func (s *Server) afterRunHooks(runID string, input core.RunInput, result core.RunResult) {
	plat := s.ensurePlatform()
	if result.Interrupted {
		plat.Fire("hitl", map[string]any{
			"run_id": runID, "target": input.TargetIP,
			"tool": func() string {
				if result.Interrupt != nil {
					return result.Interrupt.ToolName
				}
				return ""
			}(),
		})
		// Lab auto-approve private nmap (scope policy)
		if result.Interrupt != nil && result.Interrupt.ToolName == "nmap_scan" && plat.AutoApproveNmap(input.TargetIP) {
			log.Printf("talon-core: auto-approving nmap for private target %s run %s", input.TargetIP, runID)
			go func() {
				time.Sleep(400 * time.Millisecond)
				if _, ok := s.store.ClaimInterrupt(runID); ok {
					s.resumeWorkflow(runID, input, core.Decision{Type: "approve"})
				}
			}()
		}
		return
	}
	plat.IncBudget("complete", 1)
	plat.IncBudget("tool", int64(len(result.ToolLog)))
	crit := 0
	for _, f := range result.Findings {
		if f.Severity == core.SeverityCritical {
			crit++
		}
	}
	if crit > 0 {
		plat.IncBudget("critical", int64(crit))
		plat.Fire("critical", map[string]any{"run_id": runID, "target": input.TargetIP, "critical": crit})
	}
	status := "completed"
	if !result.Interrupted && result.FinalMessage != "" {
		// ok
	}
	plat.TouchTarget(input.TargetIP, runID, status)
	plat.Fire("complete", map[string]any{
		"run_id": runID, "target": input.TargetIP,
		"judge": result.JudgeVerdict, "findings": len(result.Findings),
	})
}

// ensureFindings attaches structured findings/report if the orchestrator
// returned a completed result without them (timeout soft-fail paths).
func ensureFindings(input core.RunInput, result core.RunResult) core.RunResult {
	if result.Interrupted {
		return result
	}
	if result.Report != nil && result.Findings != nil && result.KillChain != nil && result.Methodology != nil {
		return result
	}
	findings := result.Findings
	if findings == nil {
		findings = core.ExtractFindings(input, result.ToolLog, result.FinalMessage, result.JudgeVerdict, result.JudgeSet)
	}
	rep := core.BuildReport(input, result.ToolLog, result.FinalMessage, findings, result.JudgeVerdict, result.JudgeSet)
	kc := core.AnalyzeKillChain(findings)
	meth := core.ComputeMethodology(result.ToolLog, input.AgentMode)
	if rep.Markdown != "" && result.KillChain == nil {
		rep.Markdown = rep.Markdown + "\n" + kc.Summary + "\n" + core.FormatMethodologyMarkdown(meth)
	}
	result.Findings = findings
	result.Report = &rep
	if result.KillChain == nil {
		result.KillChain = &kc
	}
	if result.Methodology == nil {
		result.Methodology = &meth
	}
	return result
}

// handleTraces is GET /monitor/traces/{run_id}.
func (s *Server) handleTraces(w http.ResponseWriter, r *http.Request) {
	runID := r.PathValue("run_id")
	sess, ok := s.store.Get(runID)
	if !ok {
		writeError(w, http.StatusNotFound, "Run not found")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"history": sess.History})
}

// handleToolLog is GET /monitor/tools?run_id=... -- ToolCallRecord is
// returned per-run on RunResult, so this accumulates each run's tool log
// into its Session and serves it keyed by run_id rather than as a single
// global log.
func (s *Server) handleToolLog(w http.ResponseWriter, r *http.Request) {
	runID := r.URL.Query().Get("run_id")
	if runID == "" {
		writeError(w, http.StatusBadRequest, "run_id query parameter required")
		return
	}
	log, ok := s.store.ToolLog(runID)
	if !ok {
		writeError(w, http.StatusNotFound, "Run not found")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"tool_log": log})
}

// handleStatus is GET /output/status/{run_id}.
func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	runID := r.PathValue("run_id")
	sess, ok := s.store.Get(runID)
	if !ok {
		writeJSON(w, http.StatusOK, map[string]string{"status": "not_found"})
		return
	}
	body := map[string]any{
		"status":     sess.Status,
		"output":     sess.Output,
		"interrupt":  sess.PendingInterrupt,
		"agent_mode": sess.RunInput.AgentMode,
	}
	if sess.JudgeSet {
		body["judge_verdict"] = sess.JudgeVerdict
	}
	if len(sess.Findings) > 0 {
		body["findings_summary"] = core.SummarizeFindings(sess.Findings)
		body["findings_count"] = len(sess.Findings)
	}
	if sess.Report != nil {
		body["has_report"] = true
	}
	if sess.Methodology != nil {
		body["methodology_percent"] = sess.Methodology.Percent
	}
	writeJSON(w, http.StatusOK, body)
}

// handleFindings is GET /runs/{run_id}/findings — structured findings list.
func (s *Server) handleFindings(w http.ResponseWriter, r *http.Request) {
	runID := r.PathValue("run_id")
	sess, ok := s.store.Get(runID)
	if !ok {
		writeError(w, http.StatusNotFound, "Run not found")
		return
	}
	findings := sess.Findings
	if findings == nil {
		findings = []core.Finding{}
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"run_id":   runID,
		"findings": findings,
		"summary":  core.SummarizeFindings(findings),
	})
}

// handleReport is GET /runs/{run_id}/report — structured multi-section report.
func (s *Server) handleReport(w http.ResponseWriter, r *http.Request) {
	runID := r.PathValue("run_id")
	sess, ok := s.store.Get(runID)
	if !ok {
		writeError(w, http.StatusNotFound, "Run not found")
		return
	}
	// Lazy-build report for older sessions that only have tool log + output.
	if sess.Report == nil && sess.Status == "completed" {
		findings := sess.Findings
		if findings == nil {
			findings = core.ExtractFindings(sess.RunInput, sess.ToolLog, sess.Output, sess.JudgeVerdict, sess.JudgeSet)
		}
		rep := core.BuildReport(sess.RunInput, sess.ToolLog, sess.Output, findings, sess.JudgeVerdict, sess.JudgeSet)
		writeJSON(w, http.StatusOK, rep)
		return
	}
	if sess.Report == nil {
		writeJSON(w, http.StatusOK, map[string]any{
			"markdown": "",
			"findings": []core.Finding{},
			"summary":  core.FindingsSummary{},
			"target":   sess.RunInput.TargetIP,
			"message":  "report not available yet — run still in progress",
		})
		return
	}
	writeJSON(w, http.StatusOK, sess.Report)
}

// handleSkills is GET /skills — paginated catalog (CyberStrike pack + builtins).
// Query: brief=1, stage=, category=, q=, limit=, offset=
func (s *Server) handleSkills(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	limit, _ := strconv.Atoi(q.Get("limit"))
	offset, _ := strconv.Atoi(q.Get("offset"))
	// Default brief for large catalogs when limit not set to avoid huge payloads.
	brief := q.Get("brief") == "1" || q.Get("brief") == "true"
	if q.Get("brief") == "" && q.Get("full") != "1" {
		brief = true // default: metadata only (bodies via GET /skills/{id})
	}
	if limit == 0 {
		limit = 100
	}
	result := core.QuerySkills(core.SkillQuery{
		Stage:    q.Get("stage"),
		Category: q.Get("category"),
		Q:        q.Get("q"),
		Brief:    brief,
		Limit:    limit,
		Offset:   offset,
	})
	writeJSON(w, http.StatusOK, result)
}

// handleSkillByID is GET /skills/{id} — full skill body for UI detail pane.
func (s *Server) handleSkillByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	// IDs may contain path-like segments; PathValue gets the full {id} for Go 1.22+
	// when registered as {id} — for IDs with slashes we'd need {...id}. Our ids use dashes.
	sk, ok := core.GetSkill(id)
	if !ok {
		// try URL-decoded
		if u, err := url.PathUnescape(id); err == nil {
			sk, ok = core.GetSkill(u)
		}
	}
	if !ok {
		writeError(w, http.StatusNotFound, "skill not found")
		return
	}
	writeJSON(w, http.StatusOK, sk)
}

// handleAgents is GET /agents — specialist agent catalog.
func (s *Server) handleAgents(w http.ResponseWriter, r *http.Request) {
	agents := core.ListAgents()
	writeJSON(w, http.StatusOK, map[string]any{"agents": agents, "count": len(agents)})
}

// handleKillChain is GET /runs/{run_id}/killchain.
func (s *Server) handleKillChain(w http.ResponseWriter, r *http.Request) {
	runID := r.PathValue("run_id")
	sess, ok := s.store.Get(runID)
	if !ok {
		writeError(w, http.StatusNotFound, "Run not found")
		return
	}
	if sess.KillChain != nil {
		writeJSON(w, http.StatusOK, sess.KillChain)
		return
	}
	findings := sess.Findings
	if findings == nil {
		findings = core.ExtractFindings(sess.RunInput, sess.ToolLog, sess.Output, sess.JudgeVerdict, sess.JudgeSet)
	}
	writeJSON(w, http.StatusOK, core.AnalyzeKillChain(findings))
}

// handleMethodology is GET /runs/{run_id}/methodology.
func (s *Server) handleMethodology(w http.ResponseWriter, r *http.Request) {
	runID := r.PathValue("run_id")
	sess, ok := s.store.Get(runID)
	if !ok {
		writeError(w, http.StatusNotFound, "Run not found")
		return
	}
	if sess.Methodology != nil {
		writeJSON(w, http.StatusOK, sess.Methodology)
		return
	}
	writeJSON(w, http.StatusOK, core.ComputeMethodology(sess.ToolLog, sess.RunInput.AgentMode))
}

// handleTriageFinding is POST /runs/{run_id}/findings/{finding_id}/triage.
func (s *Server) handleTriageFinding(w http.ResponseWriter, r *http.Request) {
	runID := r.PathValue("run_id")
	findingID := r.PathValue("finding_id")
	var body struct {
		Status      string `json:"status"`
		DuplicateOf string `json:"duplicate_of"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid body")
		return
	}
	f, err := s.store.TriageFinding(runID, findingID, body.Status, body.DuplicateOf)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, f)
}

// handleGlobalFindings is GET /findings — aggregate across all runs.
func (s *Server) handleGlobalFindings(w http.ResponseWriter, r *http.Request) {
	sev := strings.ToLower(r.URL.Query().Get("severity"))
	limit := 100
	if v := r.URL.Query().Get("limit"); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 && n <= 500 {
			limit = n
		}
	}
	items := s.store.AllFindings(limit, sev)
	writeJSON(w, http.StatusOK, map[string]any{"findings": items, "count": len(items)})
}

// handleResume is POST /output/resume/{run_id}. The decision is normalized
// into a core.Decision and fed to orchestrator.Resume in a fresh goroutine.
func (s *Server) handleResume(w http.ResponseWriter, r *http.Request) {
	runID := r.PathValue("run_id")

	var req resumeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	decision, err := core.NormalizeDecision(core.Decision{
		Type:       req.Decision,
		EditedArgs: req.EditedArgs,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, "decision must be one of approve, reject, edit (edit requires edited_args)")
		return
	}

	// Atomically claim the interrupt so concurrent resume requests cannot
	// both launch resumeWorkflow (double-execution / lost session race).
	sess, ok := s.store.ClaimInterrupt(runID)
	if !ok {
		writeError(w, http.StatusBadRequest, "No pending interrupt")
		return
	}

	// Ensure RunInput carries SessionID for orchestrator session lookup.
	input := sess.RunInput
	if input.SessionID == "" {
		input.SessionID = runID
	}
	s.store.AppendHistory(runID, "decision=%s", decision.Type)

	go s.resumeWorkflow(runID, input, decision)

	writeJSON(w, http.StatusOK, map[string]string{"message": "Decision received, resuming orchestrator..."})
}

func (s *Server) resumeWorkflow(runID string, input core.RunInput, decision core.Decision) {
	ctx, cancel := context.WithTimeout(context.Background(), runTimeout())
	defer cancel()
	ctx = core.WithProgress(ctx, func(toolLog []core.ToolCallRecord) {
		s.store.SetToolLog(runID, toolLog)
	})
	ctx = core.WithFindingsProgress(ctx, func(findings []core.Finding) {
		s.store.SetFindings(runID, findings)
	})

	log.Printf("talon-core: resume %s decision=%s", runID, decision.Type)
	result, err := s.orch.Resume(ctx, input, decision)
	if err != nil {
		log.Printf("talon-core: resume %s: %v", runID, err)
		if ctx.Err() != nil && len(result.ToolLog) > 0 {
			result.FinalMessage = strings.TrimSpace(result.FinalMessage + "\n[run stopped: wall-clock timeout]")
			result.Interrupted = false
			s.store.SetResult(runID, ensureFindings(input, result))
			return
		}
		s.store.SetError(runID, err)
		return
	}
	log.Printf("talon-core: resume %s finished interrupted=%v tools=%d findings=%d",
		runID, result.Interrupted, len(result.ToolLog), len(result.Findings))
	s.store.SetResult(runID, ensureFindings(input, result))
}

// handleMCPServers is GET /mcp/servers — connected MCP servers + their tools,
// plus the in-process "talon-core" virtual tools (skills, findings, A2A delegates).
func (s *Server) handleMCPServers(w http.ResponseWriter, r *http.Request) {
	var servers []mcpclient.ServerInfo
	if s.tools != nil {
		servers = s.tools.Servers()
	}
	// Virtual in-process agent tools (not separate MCP processes).
	coreTools := []string{
		"skill_search", "skill_get",
		"report_finding", "triage_finding",
		"delegate_recon", "delegate_exploit", "delegate_post_exploit",
		"delegate_codegen", "delegate_report",
	}
	servers = append(servers, mcpclient.ServerInfo{
		Name:  "talon-core (in-process)",
		Tools: coreTools,
	})
	// Skill pack availability for operators.
	stats := core.SkillStats()
	writeJSON(w, http.StatusOK, map[string]any{
		"servers":     servers,
		"skill_stats": stats,
		"agent_to_agent": map[string]any{
			"model": "orchestrator → subagents via delegate_* tools",
			"notes": []string{
				"Subagents do not call each other directly.",
				"Orchestrator sequences recon → exploit → post_exploit → codegen → report.",
				"Subagents share context only through orchestrator instructions + return text.",
				"CyberStrike skills: skill_search / skill_get on every subagent.",
				"MCP tool servers: hexstrike (arsenal) + metasploit (strike) over stdio.",
			},
		},
	})
}

// handleListRuns is GET /runs — paginated summaries, newest first.
// Query params: limit (default 100, max 500), offset (default 0).
// Response: {runs, total, limit, offset} so tables can page server-side.
func (s *Server) handleListRuns(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	limit, _ := strconv.Atoi(q.Get("limit"))
	offset, _ := strconv.Atoi(q.Get("offset"))
	if limit <= 0 {
		limit = 100
	}
	if limit > 500 {
		limit = 500
	}
	if offset < 0 {
		offset = 0
	}
	runs, total, err := s.store.PaginatedList(limit, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list runs")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"runs":   runs,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

// handleRunsSummary is GET /runs/summary — aggregate counts for the overview.
func (s *Server) handleRunsSummary(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.store.Summary())
}

// terminalStatus reports whether a run status is final (stream can close).
func terminalStatus(status string) bool {
	switch status {
	case "completed", "error", "not_found":
		return true
	}
	return false
}

// handleStream is GET /monitor/stream/{run_id} — Server-Sent Events stream of
// a run's live state. Emits:
//
//	event: tool    data: ToolCallRecord JSON (one per newly completed tool)
//	event: status  data: {status, output?, judge_verdict?} JSON
//
// It polls the store once per second (the store is the single source of
// truth for both this and the polling endpoints) and closes after sending
// the terminal status. Clients should fall back to polling
// /output/status + /monitor/tools if the stream drops.
func (s *Server) handleStream(w http.ResponseWriter, r *http.Request) {
	runID := r.PathValue("run_id")
	flusher, ok := w.(http.Flusher)
	if !ok {
		writeError(w, http.StatusInternalServerError, "streaming unsupported")
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.WriteHeader(http.StatusOK)

	sentTools := 0
	lastStatus := ""
	lastFindings := -1
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		sess, found := s.store.Get(runID)
		status := "not_found"
		if found {
			status = sess.Status
		}

		if found && len(sess.ToolLog) > sentTools {
			for _, rec := range sess.ToolLog[sentTools:] {
				data, err := json.Marshal(rec)
				if err == nil {
					fmt.Fprintf(w, "event: tool\ndata: %s\n\n", data)
				}
			}
			sentTools = len(sess.ToolLog)
		}

		// Live findings stream (mid-run report_finding). Skip initial empty snapshot.
		if found {
			nFind := len(sess.Findings)
			if nFind != lastFindings {
				if nFind > 0 || lastFindings > 0 {
					body := map[string]any{
						"findings_count":   nFind,
						"findings_summary": core.SummarizeFindings(sess.Findings),
						"findings":         sess.Findings,
					}
					if data, err := json.Marshal(body); err == nil {
						fmt.Fprintf(w, "event: findings\ndata: %s\n\n", data)
					}
				}
				lastFindings = nFind
			}
		}

		if status != lastStatus {
			body := map[string]any{"status": status}
			if found {
				body["agent_mode"] = sess.RunInput.AgentMode
				body["findings_count"] = len(sess.Findings)
				if len(sess.Findings) > 0 {
					body["findings_summary"] = core.SummarizeFindings(sess.Findings)
				}
			}
			if found && terminalStatus(status) {
				body["output"] = sess.Output
				if sess.JudgeSet {
					body["judge_verdict"] = sess.JudgeVerdict
				}
				if sess.Report != nil {
					body["has_report"] = true
				}
			}
			if found && status == "awaiting_approval" {
				body["interrupt"] = sess.PendingInterrupt
			}
			data, err := json.Marshal(body)
			if err == nil {
				fmt.Fprintf(w, "event: status\ndata: %s\n\n", data)
			}
			lastStatus = status
		}
		flusher.Flush()

		if terminalStatus(status) {
			return
		}

		select {
		case <-r.Context().Done():
			return
		case <-ticker.C:
		}
	}
}

// wsMessage is the JSON envelope sent over the WebSocket stream:
// {"type":"tool","data":ToolCallRecord} or {"type":"status","data":{...}}.
type wsMessage struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}

// handleStreamWS is GET /monitor/ws/{run_id} — WebSocket variant of the SSE
// stream (same payloads, same 1s store poll, closes on terminal status).
// The dashboard prefers this and falls back to SSE/polling if unavailable.
func (s *Server) handleStreamWS(w http.ResponseWriter, r *http.Request) {
	runID := r.PathValue("run_id")
	// OriginPatterns "*": the dashboard (and CLI tooling) connect cross-origin;
	// matches the permissive CORS posture.
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{OriginPatterns: []string{"*"}})
	if err != nil {
		log.Printf("talon-core: ws accept %s: %v", runID, err)
		return
	}
	defer conn.Close(websocket.StatusNormalClosure, "done")

	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Minute)
	defer cancel()

	send := func(msg wsMessage) error {
		data, err := json.Marshal(msg)
		if err != nil {
			return nil
		}
		return conn.Write(ctx, websocket.MessageText, data)
	}

	sentTools := 0
	lastStatus := ""
	lastFindings := -1
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		sess, found := s.store.Get(runID)
		status := "not_found"
		if found {
			status = sess.Status
		}

		if found && len(sess.ToolLog) > sentTools {
			for _, rec := range sess.ToolLog[sentTools:] {
				if err := send(wsMessage{Type: "tool", Data: rec}); err != nil {
					return
				}
			}
			sentTools = len(sess.ToolLog)
		}

		if found {
			nFind := len(sess.Findings)
			if nFind != lastFindings {
				if nFind > 0 || lastFindings > 0 {
					body := map[string]any{
						"findings_count":   nFind,
						"findings_summary": core.SummarizeFindings(sess.Findings),
						"findings":         sess.Findings,
					}
					if err := send(wsMessage{Type: "findings", Data: body}); err != nil {
						return
					}
				}
				lastFindings = nFind
			}
		}

		if status != lastStatus {
			body := map[string]any{"status": status}
			if found {
				body["agent_mode"] = sess.RunInput.AgentMode
				body["findings_count"] = len(sess.Findings)
				if len(sess.Findings) > 0 {
					body["findings_summary"] = core.SummarizeFindings(sess.Findings)
				}
			}
			if found && terminalStatus(status) {
				body["output"] = sess.Output
				if sess.JudgeSet {
					body["judge_verdict"] = sess.JudgeVerdict
				}
				if sess.Report != nil {
					body["has_report"] = true
				}
			}
			if found && status == "awaiting_approval" {
				body["interrupt"] = sess.PendingInterrupt
			}
			if err := send(wsMessage{Type: "status", Data: body}); err != nil {
				return
			}
			lastStatus = status
		}

		if terminalStatus(status) {
			return
		}

		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}


