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
	if req.LHOST == "" {
		// Prefer operator env (compose sets LHOST=127.0.0.1 for local lab).
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

	runID := uuid.NewString()
	input := core.RunInput{
		SessionID:   runID,
		TargetIP:    req.IP,
		CVEID:       req.CVEID,
		ServiceName: req.ServiceName,
		Description: req.Description,
		Context:     config.Context{LHOST: req.LHOST, LPORT: req.LPORT},
	}
	s.store.Create(runID, input)

	// Runs in the background; on a HITL interrupt it stores the pending
	// interrupt and returns rather than blocking -- POST
	// /output/resume/{run_id} is what drives the workflow forward from
	// there (see handleResume).
	go s.runWorkflow(runID, input)

	writeJSON(w, http.StatusOK, map[string]string{
		"run_id":  runID,
		"message": "Agent execution started",
	})
}

func (s *Server) runWorkflow(runID string, input core.RunInput) {
	ctx, cancel := context.WithTimeout(context.Background(), runTimeout())
	defer cancel()
	ctx = core.WithProgress(ctx, func(toolLog []core.ToolCallRecord) {
		s.store.SetToolLog(runID, toolLog)
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
			s.store.SetResult(runID, result)
			return
		}
		s.store.SetError(runID, err)
		return
	}
	log.Printf("talon-core: run %s finished interrupted=%v tools=%d", runID, result.Interrupted, len(result.ToolLog))
	s.store.SetResult(runID, result)
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
		"status":    sess.Status,
		"output":    sess.Output,
		"interrupt": sess.PendingInterrupt,
	}
	if sess.JudgeSet {
		body["judge_verdict"] = sess.JudgeVerdict
	}
	writeJSON(w, http.StatusOK, body)
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

	log.Printf("talon-core: resume %s decision=%s", runID, decision.Type)
	result, err := s.orch.Resume(ctx, input, decision)
	if err != nil {
		log.Printf("talon-core: resume %s: %v", runID, err)
		if ctx.Err() != nil && len(result.ToolLog) > 0 {
			result.FinalMessage = strings.TrimSpace(result.FinalMessage + "\n[run stopped: wall-clock timeout]")
			result.Interrupted = false
			s.store.SetResult(runID, result)
			return
		}
		s.store.SetError(runID, err)
		return
	}
	log.Printf("talon-core: resume %s finished interrupted=%v tools=%d", runID, result.Interrupted, len(result.ToolLog))
	s.store.SetResult(runID, result)
}

// handleMCPServers is GET /mcp/servers — connected MCP servers + their tools.
func (s *Server) handleMCPServers(w http.ResponseWriter, r *http.Request) {
	if s.tools == nil {
		writeJSON(w, http.StatusOK, map[string]any{"servers": []any{}})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"servers": s.tools.Servers()})
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

		if status != lastStatus {
			body := map[string]any{"status": status}
			if found && terminalStatus(status) {
				body["output"] = sess.Output
				if sess.JudgeSet {
					body["judge_verdict"] = sess.JudgeVerdict
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

		if status != lastStatus {
			body := map[string]any{"status": status}
			if found && terminalStatus(status) {
				body["output"] = sess.Output
				if sess.JudgeSet {
					body["judge_verdict"] = sess.JudgeVerdict
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


