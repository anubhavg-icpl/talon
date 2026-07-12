package control

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/anubhavg-icpl/talon/internal/config"
	"github.com/anubhavg-icpl/talon/internal/core"
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

// Server wires the Store to an core.Orchestrator over 5 routes:
// /input/start, /monitor/traces/{run_id}, /monitor/tools,
// /output/status/{run_id}, /output/resume/{run_id}.
type Server struct {
	orch  *core.Orchestrator
	store *Store
}

func NewServer(orch *core.Orchestrator, store *Store) *Server {
	return &Server{orch: orch, store: store}
}

// Mux builds the http.ServeMux, using Go 1.22+ method+pattern routing so no
// router dependency is needed.
func (s *Server) Mux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", s.handleHealth)
	mux.HandleFunc("POST /input/start", s.handleStart)
	mux.HandleFunc("GET /monitor/traces/{run_id}", s.handleTraces)
	mux.HandleFunc("GET /monitor/tools", s.handleToolLog)
	mux.HandleFunc("GET /output/status/{run_id}", s.handleStatus)
	mux.HandleFunc("POST /output/resume/{run_id}", s.handleResume)
	return mux
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
	writeJSON(w, http.StatusOK, map[string]any{
		"status":    sess.Status,
		"output":    sess.Output,
		"interrupt": sess.PendingInterrupt,
	})
}

// handleResume is POST /output/resume/{run_id}. The decision is decoded
// straight into a core.Decision and fed to orchestrator.Resume in a fresh
// goroutine, driving the workflow forward from its paused point.
func (s *Server) handleResume(w http.ResponseWriter, r *http.Request) {
	runID := r.PathValue("run_id")
	sess, ok := s.store.Get(runID)
	if !ok || sess.PendingInterrupt == nil {
		writeError(w, http.StatusBadRequest, "No pending interrupt")
		return
	}

	var req resumeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if !isValidDecision(req.Decision) {
		writeError(w, http.StatusBadRequest, "decision must be one of approve, reject, edit")
		return
	}

	decision := core.Decision{Type: req.Decision, EditedArgs: req.EditedArgs}
	s.store.ClearInterrupt(runID)
	s.store.SetStatus(runID, "running")

	go s.resumeWorkflow(runID, sess.RunInput, decision)

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

func isValidDecision(d string) bool {
	switch strings.ToLower(d) {
	case "approve", "reject", "edit":
		return true
	default:
		return false
	}
}
