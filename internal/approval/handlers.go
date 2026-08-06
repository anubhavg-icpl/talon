// Handlers for the approval HITL store. These are self-contained HTTP handlers
// (package approval) that internal/control/server.go can import and register on
// its mux, e.g.:
//
//	ah := approval.NewHandler(store)
//	ah.Register(mux)
//
// or wired individually:
//
//	mux.HandleFunc("POST /approvals", ah.HandleCreate)
//	mux.HandleFunc("GET /approvals/{id}", ah.HandleGet)
//
// The handlers own their JSON helpers so they have no dependency on the control
// package, keeping the import graph one-directional (control -> approval).
package approval

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Handler exposes the ActionStore over HTTP. It is safe for concurrent use.
type Handler struct {
	store *ActionStore
	// NewID, when non-nil, generates IDs for created actions (defaults to a
	// random UUID). Useful for deterministic tests.
	NewID func() string
}

// NewHandler wraps an ActionStore with HTTP handlers.
func NewHandler(store *ActionStore) *Handler {
	return &Handler{
		store: store,
		NewID: uuid.NewString,
	}
}

// Register mounts all approval routes on the given mux using Go 1.22+ method
// routing. Routes:
//
//	POST   /approvals              create a new pending action
//	GET    /approvals              list actions (?run_id=&pending=true)
//	GET    /approvals/pending      list pending actions (?run_id=)
//	GET    /approvals/{id}         fetch one action
//	POST   /approvals/{id}/claim   claim an action (claim-before-dispatch)
//	POST   /approvals/{id}/approve approve an action with a result body
//	POST   /approvals/{id}/reject  reject an action with a reason
//	POST   /approvals/{id}/fail    mark a dispatched action failed
func (h *Handler) Register(mux *http.ServeMux) {
	if mux == nil {
		return
	}
	mux.HandleFunc("POST /approvals", h.HandleCreate)
	mux.HandleFunc("GET /approvals", h.HandleList)
	mux.HandleFunc("GET /approvals/pending", h.HandleListPending)
	mux.HandleFunc("GET /approvals/{id}", h.HandleGet)
	mux.HandleFunc("POST /approvals/{id}/claim", h.HandleClaim)
	mux.HandleFunc("POST /approvals/{id}/approve", h.HandleApprove)
	mux.HandleFunc("POST /approvals/{id}/reject", h.HandleReject)
	mux.HandleFunc("POST /approvals/{id}/fail", h.HandleFail)
}

// --- request/response shapes -------------------------------------------------

// createRequest is the body for POST /approvals.
type createRequest struct {
	RunID     string          `json:"run_id"`
	ToolName  string          `json:"tool_name"`
	Args      json.RawMessage `json:"args,omitempty"`
	RiskLevel RiskLevel       `json:"risk_level,omitempty"`
	Summary   string          `json:"summary,omitempty"`
	// AutoApprove, when true, allows the handler to auto-resolve low/medium-risk
	// actions if the store is configured for it (TALON_AUTO_APPROVE_RISK).
	AutoApprove bool `json:"auto_approve,omitempty"`
}

// approveRequest is the body for POST /approvals/{id}/approve.
type approveRequest struct {
	Result json.RawMessage `json:"result,omitempty"`
}

// rejectRequest is the body for POST /approvals/{id}/reject.
type rejectRequest struct {
	Reason string `json:"reason"`
}

// failRequest is the body for POST /approvals/{id}/fail.
type failRequest struct {
	Detail string `json:"detail"`
}

// --- handlers ----------------------------------------------------------------

// HandleCreate is POST /approvals — record a new pending action. If
// AutoApprove is requested in the body and the store's policy permits the
// risk level, the action is created and immediately approved.
func (h *Handler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	var req createRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONErr(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if strings.TrimSpace(req.RunID) == "" {
		writeJSONErr(w, http.StatusBadRequest, "run_id required")
		return
	}
	if strings.TrimSpace(req.ToolName) == "" {
		writeJSONErr(w, http.StatusBadRequest, "tool_name required")
		return
	}
	if req.RiskLevel == "" {
		// Default risk from the danger classification: dangerous tools start
		// high, everything else low.
		if IsDangerous(req.ToolName) {
			req.RiskLevel = RiskHigh
		} else {
			req.RiskLevel = RiskLow
		}
	}

	action := &Action{
		ID:        h.newID(),
		RunID:     req.RunID,
		ToolName:  req.ToolName,
		Args:      req.Args,
		RiskLevel: req.RiskLevel,
		Summary:   req.Summary,
		State:     StatePending,
		CreatedAt: time.Now().UTC(),
	}
	if err := h.store.Create(action); err != nil {
		writeJSONErr(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Optional auto-approval path for low/medium-risk actions.
	autoApproved := false
	if req.AutoApprove && h.store.AutoApprove(req.RiskLevel) {
		if err := h.store.Approve(action.ID, json.RawMessage(`{"auto":true}`)); err == nil {
			autoApproved = true
			if updated, err := h.store.Get(action.ID); err == nil {
				action = updated
			}
		}
	}
	writeJSON(w, http.StatusCreated, map[string]any{
		"action":       action,
		"auto_approved": autoApproved,
	})
}

// HandleGet is GET /approvals/{id}.
func (h *Handler) HandleGet(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	a, err := h.store.Get(id)
	if err != nil {
		status := http.StatusInternalServerError
		if err == ErrNotFound {
			status = http.StatusNotFound
		}
		writeJSONErr(w, status, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"action": a})
}

// HandleList is GET /approvals — list actions, optionally filtered by ?run_id=
// and limited to pending via ?pending=true.
func (h *Handler) HandleList(w http.ResponseWriter, r *http.Request) {
	runID := r.URL.Query().Get("run_id")
	pending := r.URL.Query().Has("pending")
	if pending {
		h.listPending(w, runID)
		return
	}
	actions, err := h.store.ListAll(runID)
	if err != nil {
		writeJSONErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"actions": actions,
		"count":   len(actions),
	})
}

// HandleListPending is GET /approvals/pending — list only pending actions.
func (h *Handler) HandleListPending(w http.ResponseWriter, r *http.Request) {
	runID := r.URL.Query().Get("run_id")
	h.listPending(w, runID)
}

func (h *Handler) listPending(w http.ResponseWriter, runID string) {
	actions, err := h.store.ListPending(runID)
	if err != nil {
		writeJSONErr(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"actions": actions,
		"count":   len(actions),
	})
}

// HandleClaim is POST /approvals/{id}/claim — claim an action before dispatch.
func (h *Handler) HandleClaim(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	a, err := h.store.Claim(id)
	if err != nil {
		status := http.StatusInternalServerError
		switch err {
		case ErrNotClaimable:
			status = http.StatusConflict
		case ErrNotFound:
			status = http.StatusNotFound
		}
		writeJSONErr(w, status, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"action": a})
}

// HandleApprove is POST /approvals/{id}/approve.
func (h *Handler) HandleApprove(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req approveRequest
	// Body is optional; an empty result is allowed.
	_ = json.NewDecoder(r.Body).Decode(&req)

	if err := h.store.Approve(id, req.Result); err != nil {
		status := http.StatusInternalServerError
		switch err {
		case ErrAlreadyResolved:
			status = http.StatusConflict
		case ErrNotFound:
			status = http.StatusNotFound
		}
		writeJSONErr(w, status, err.Error())
		return
	}
	a, _ := h.store.Get(id)
	writeJSON(w, http.StatusOK, map[string]any{"action": a})
}

// HandleReject is POST /approvals/{id}/reject.
func (h *Handler) HandleReject(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req rejectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONErr(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if strings.TrimSpace(req.Reason) == "" {
		writeJSONErr(w, http.StatusBadRequest, "reason required")
		return
	}
	if err := h.store.Reject(id, req.Reason); err != nil {
		status := http.StatusInternalServerError
		switch err {
		case ErrAlreadyResolved:
			status = http.StatusConflict
		case ErrNotFound:
			status = http.StatusNotFound
		}
		writeJSONErr(w, status, err.Error())
		return
	}
	a, _ := h.store.Get(id)
	writeJSON(w, http.StatusOK, map[string]any{"action": a})
}

// HandleFail is POST /approvals/{id}/fail — record a dispatch/execution error.
func (h *Handler) HandleFail(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req failRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONErr(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if err := h.store.Fail(id, req.Detail); err != nil {
		status := http.StatusInternalServerError
		switch err {
		case ErrAlreadyResolved:
			status = http.StatusConflict
		case ErrNotFound:
			status = http.StatusNotFound
		}
		writeJSONErr(w, status, err.Error())
		return
	}
	a, _ := h.store.Get(id)
	writeJSON(w, http.StatusOK, map[string]any{"action": a})
}

// --- helpers -----------------------------------------------------------------

func (h *Handler) newID() string {
	if h.NewID != nil {
		return h.NewID()
	}
	return uuid.NewString()
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeJSONErr(w http.ResponseWriter, status int, detail string) {
	writeJSON(w, status, map[string]string{"detail": detail})
}
