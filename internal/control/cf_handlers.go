// Package control — cf_handlers.go: HTTP handlers for CF-derived features.
// These expose VFS, Approval, Gatekeeper, Blueprint, Audit, MCP Gateway,
// and Sharing systems via REST routes.
package control

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/anubhavg-icpl/talon/internal/approval"
	"github.com/anubhavg-icpl/talon/internal/audit"
	"github.com/anubhavg-icpl/talon/internal/blueprint"
	"github.com/anubhavg-icpl/talon/internal/gatekeeper"
	"github.com/anubhavg-icpl/talon/internal/mcpgw"
	"github.com/anubhavg-icpl/talon/internal/sharing"
	"github.com/anubhavg-icpl/talon/internal/vfs"
	"github.com/google/uuid"
)

// CFIntegration holds all the CF-derived subsystem instances.
type CFIntegration struct {
	VFS         *vfs.VFS
	Approvals   *approval.ActionStore
	Gatekeepers *gatekeeper.GatekeeperRegistry
	Blueprints  *blueprint.BlueprintStore
	Audit       *audit.AuditStore
	MCPGateway  *mcpgw.Gateway
	Sharing     *sharing.Store
}

// InitCFIntegration initializes all CF-derived subsystems.
func InitCFIntegration(dataDir string) (*CFIntegration, error) {
	if dataDir == "" {
		dataDir = filepath.Join(os.Getenv("HOME"), ".talon", "data")
	}
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		return nil, fmt.Errorf("cf: mkdir data dir: %w", err)
	}

	v, err := vfs.Open(filepath.Join(dataDir, "vfs.db"))
	if err != nil {
		return nil, fmt.Errorf("cf: open vfs: %w", err)
	}

	apStore, err := approval.Open(filepath.Join(dataDir, "approvals.db"))
	if err != nil {
		return nil, fmt.Errorf("cf: open approvals: %w", err)
	}

	gkDB := openSQLite(filepath.Join(dataDir, "gatekeepers.db"))
	gkRegistry, err := gatekeeper.NewGatekeeperRegistry(gkDB)
	if err != nil {
		return nil, fmt.Errorf("cf: gatekeeper registry: %w", err)
	}

	bpDB := openSQLite(filepath.Join(dataDir, "blueprints.db"))
	bpStore := blueprint.NewBlueprintStore(bpDB)
	if err := bpStore.Migrate(); err != nil {
		return nil, fmt.Errorf("cf: blueprint migrate: %w", err)
	}
	if err := bpStore.SeedDefaults(); err != nil {
		fmt.Printf("cf: warning: seed blueprints: %v\n", err)
	}

	auDB := openSQLite(filepath.Join(dataDir, "audit.db"))
	auStore := audit.NewAuditStore(auDB)
	if err := auStore.Migrate(); err != nil {
		return nil, fmt.Errorf("cf: audit migrate: %w", err)
	}

	mcpBase := os.Getenv("TALON_MCP_BASE_URL")
	if mcpBase == "" {
		mcpBase = "http://localhost:8888"
	}
	gw := mcpgw.New(apStore, mcpBase)

	shareSecret := os.Getenv("TALON_SHARE_SECRET")
	if shareSecret == "" {
		shareSecret = uuid.NewString()
	}
	sh := sharing.New(shareSecret)

	return &CFIntegration{
		VFS:         v,
		Approvals:   apStore,
		Gatekeepers: gkRegistry,
		Blueprints:  bpStore,
		Audit:       auStore,
		MCPGateway:  gw,
		Sharing:     sh,
	}, nil
}

func openSQLite(path string) *sql.DB {
	db, err := sql.Open("sqlite3", path+"?_journal=WAL&_busy_timeout=5000")
	if err != nil {
		panic(fmt.Sprintf("failed to open sqlite %s: %v", path, err))
	}
	return db
}

// WithCFIntegration is a ServerOption that attaches CF subsystems.
func WithCFIntegration(cf *CFIntegration) ServerOption {
	return func(s *Server) {
		s.cf = cf
	}
}

// SetCFIntegration attaches CF subsystems after server creation.
func (s *Server) SetCFIntegration(cf *CFIntegration) {
	s.cf = cf
}

// RegisterCFRoutes adds all CF-derived feature routes to the mux.
func (s *Server) RegisterCFRoutes(mux *http.ServeMux) {
	if s.cf == nil {
		return
	}

	// VFS
	mux.HandleFunc("GET /vfs", s.handleVFSList)
	mux.HandleFunc("POST /vfs/file", s.handleVFSWrite)
	mux.HandleFunc("GET /vfs/file", s.handleVFSRead)
	mux.HandleFunc("DELETE /vfs/file", s.handleVFSDelete)
	mux.HandleFunc("POST /vfs/mkdir", s.handleVFSMkdir)
	mux.HandleFunc("GET /vfs/stat", s.handleVFSStat)

	// Approvals
	mux.HandleFunc("POST /approvals", s.handleApprovalCreate)
	mux.HandleFunc("GET /approvals", s.handleApprovalList)
	mux.HandleFunc("GET /approvals/pending", s.handleApprovalListPending)
	mux.HandleFunc("GET /approvals/{id}", s.handleApprovalGet)
	mux.HandleFunc("POST /approvals/{id}/approve", s.handleApprovalApprove)
	mux.HandleFunc("POST /approvals/{id}/reject", s.handleApprovalReject)
	mux.HandleFunc("GET /approvals/check/{tool}", s.handleApprovalIsDangerous)

	// Gatekeepers
	mux.HandleFunc("GET /gatekeepers", s.handleGatekeeperList)
	mux.HandleFunc("POST /gatekeepers", s.handleGatekeeperRegister)
	mux.HandleFunc("DELETE /gatekeepers/{name}", s.handleGatekeeperRemove)
	mux.HandleFunc("POST /gatekeepers/{name}/access", s.handleGatekeeperRequestAccess)
	mux.HandleFunc("GET /gatekeepers/{name}/actions", s.handleGatekeeperActions)
	mux.HandleFunc("POST /gatekeepers/{name}/sessions/{sid}/revoke", s.handleGatekeeperRevoke)

	// Blueprints
	mux.HandleFunc("GET /blueprints", s.handleBlueprintList)
	mux.HandleFunc("POST /blueprints", s.handleBlueprintCreate)
	mux.HandleFunc("GET /blueprints/{id}", s.handleBlueprintGet)
	mux.HandleFunc("PUT /blueprints/{id}", s.handleBlueprintUpdate)
	mux.HandleFunc("DELETE /blueprints/{id}", s.handleBlueprintDelete)

	// Audit
	mux.HandleFunc("GET /audit/{run_id}", s.handleAuditList)
	mux.HandleFunc("POST /audit", s.handleAuditLog)
	mux.HandleFunc("GET /audit/{run_id}/export", s.handleAuditExport)
	mux.HandleFunc("GET /audit/{run_id}/stats", s.handleAuditStats)

	// MCP Gateway
	mux.HandleFunc("GET /mcp/tools", s.handleMCPList)
	mux.HandleFunc("POST /mcp/tools", s.handleMCPRegister)
	mux.HandleFunc("POST /mcp/call", s.handleMCPCall)
	mux.HandleFunc("POST /mcp/approve/{id}", s.handleMCPExecuteApproved)
	mux.HandleFunc("POST /mcp/reject/{id}", s.handleMCPReject)
	mux.HandleFunc("POST /mcp/vet/{tool}", s.handleMCPVet)
	mux.HandleFunc("GET /mcp/stats", s.handleMCPStats)

	// Sharing
	mux.HandleFunc("GET /engagements", s.handleSharingListEngagements)
	mux.HandleFunc("POST /engagements", s.handleSharingCreateEngagement)
	mux.HandleFunc("GET /engagements/{id}", s.handleSharingGetEngagement)
	mux.HandleFunc("POST /engagements/{id}/shares", s.handleSharingCreateLink)
	mux.HandleFunc("GET /engagements/{id}/shares", s.handleSharingListLinks)
	mux.HandleFunc("POST /engagements/{id}/shares/{linkID}/revoke", s.handleSharingRevokeLink)
	mux.HandleFunc("POST /share/accept", s.handleSharingAccept)
	mux.HandleFunc("GET /engagements/{id}/collaborators", s.handleSharingListCollaborators)
	mux.HandleFunc("DELETE /engagements/{id}/collaborators/{userID}", s.handleSharingRemoveCollab)
	mux.HandleFunc("POST /engagements/{id}/runs/{runID}", s.handleSharingAddRun)
}

// ==================== VFS ====================

func (s *Server) handleVFSList(w http.ResponseWriter, r *http.Request) {
	dir := r.URL.Query().Get("dir")
	if dir == "" {
		dir = "/"
	}
	entries, err := s.cf.VFS.List(dir)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, entries)
}

func (s *Server) handleVFSWrite(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Path    string `json:"path"`
		Content string `json:"content"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := s.cf.VFS.WriteFile(req.Path, []byte(req.Content)); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "written", "path": req.Path})
}

func (s *Server) handleVFSRead(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Query().Get("path")
	if p == "" {
		writeError(w, http.StatusBadRequest, "path required")
		return
	}
	data, err := s.cf.VFS.ReadFile(p)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(data)
}

func (s *Server) handleVFSDelete(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Query().Get("path")
	if p == "" {
		writeError(w, http.StatusBadRequest, "path required")
		return
	}
	if err := s.cf.VFS.Delete(p); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

func (s *Server) handleVFSMkdir(w http.ResponseWriter, r *http.Request) {
	var req struct{ Path string `json:"path"` }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := s.cf.VFS.Mkdir(req.Path); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "created"})
}

func (s *Server) handleVFSStat(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Query().Get("path")
	if p == "" {
		writeError(w, http.StatusBadRequest, "path required")
		return
	}
	entry, err := s.cf.VFS.Stat(p)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, entry)
}

// ==================== Approvals ====================

func (s *Server) handleApprovalCreate(w http.ResponseWriter, r *http.Request) {
	var action approval.Action
	if err := json.NewDecoder(r.Body).Decode(&action); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if action.ID == "" {
		action.ID = uuid.NewString()
	}
	if err := s.cf.Approvals.Create(&action); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, action)
}

func (s *Server) handleApprovalList(w http.ResponseWriter, r *http.Request) {
	runID := r.URL.Query().Get("run_id")
	if runID == "" {
		writeError(w, http.StatusBadRequest, "run_id required")
		return
	}
	actions, err := s.cf.Approvals.ListAll(runID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, actions)
}

func (s *Server) handleApprovalListPending(w http.ResponseWriter, r *http.Request) {
	runID := r.URL.Query().Get("run_id")
	if runID == "" {
		writeError(w, http.StatusBadRequest, "run_id required")
		return
	}
	actions, err := s.cf.Approvals.ListPending(runID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, actions)
}

func (s *Server) handleApprovalGet(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	action, err := s.cf.Approvals.Get(id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, action)
}

func (s *Server) handleApprovalApprove(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req struct{ Result json.RawMessage `json:"result"` }
	json.NewDecoder(r.Body).Decode(&req)
	if err := s.cf.Approvals.Approve(id, req.Result); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "approved"})
}

func (s *Server) handleApprovalReject(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req struct{ Reason string `json:"reason"` }
	json.NewDecoder(r.Body).Decode(&req)
	if err := s.cf.Approvals.Reject(id, req.Reason); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "rejected"})
}

func (s *Server) handleApprovalIsDangerous(w http.ResponseWriter, r *http.Request) {
	tool := r.PathValue("tool")
	writeJSON(w, http.StatusOK, map[string]any{
		"tool":      tool,
		"dangerous": approval.IsDangerous(tool),
	})
}

// ==================== Gatekeepers ====================

func (s *Server) handleGatekeeperList(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.cf.Gatekeepers.List())
}

func (s *Server) handleGatekeeperRegister(w http.ResponseWriter, r *http.Request) {
	var config gatekeeper.GatekeeperConfig
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	s.cf.Gatekeepers.Register(config)
	writeJSON(w, http.StatusCreated, map[string]string{"status": "registered", "name": config.Name})
}

func (s *Server) handleGatekeeperRemove(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	if err := s.cf.Gatekeepers.Remove(name); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "removed"})
}

func (s *Server) handleGatekeeperRequestAccess(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	gk, err := s.cf.Gatekeepers.Get(name)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	var cap gatekeeper.Capability
	if err := json.NewDecoder(r.Body).Decode(&cap); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	session, err := gk.RequestAccess(cap)
	if err != nil {
		writeError(w, http.StatusForbidden, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, session)
}

func (s *Server) handleGatekeeperActions(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	gk, err := s.cf.Gatekeepers.Get(name)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	sessionID := r.URL.Query().Get("session")
	actions, err := gk.GetActions(sessionID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, actions)
}

func (s *Server) handleGatekeeperRevoke(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	sessionID := r.PathValue("sid")
	gk, err := s.cf.Gatekeepers.Get(name)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	if err := gk.Revoke(sessionID); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "revoked"})
}

// ==================== Blueprints ====================

func (s *Server) handleBlueprintList(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	blueprints, err := s.cf.Blueprints.List(category)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, blueprints)
}

func (s *Server) handleBlueprintCreate(w http.ResponseWriter, r *http.Request) {
	var bp blueprint.Blueprint
	if err := json.NewDecoder(r.Body).Decode(&bp); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := s.cf.Blueprints.Create(bp); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, bp)
}

func (s *Server) handleBlueprintGet(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	bp, err := s.cf.Blueprints.Get(id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, bp)
}

func (s *Server) handleBlueprintUpdate(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var bp blueprint.Blueprint
	if err := json.NewDecoder(r.Body).Decode(&bp); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	bp.ID = id
	if err := s.cf.Blueprints.Update(bp); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, bp)
}

func (s *Server) handleBlueprintDelete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := s.cf.Blueprints.Delete(id); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "deleted"})
}

// ==================== Audit ====================

func (s *Server) handleAuditList(w http.ResponseWriter, r *http.Request) {
	runID := r.PathValue("run_id")
	entries, err := s.cf.Audit.List(runID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, entries)
}

func (s *Server) handleAuditLog(w http.ResponseWriter, r *http.Request) {
	var entry audit.AuditEntry
	if err := json.NewDecoder(r.Body).Decode(&entry); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := s.cf.Audit.Log(entry); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, entry)
}

func (s *Server) handleAuditExport(w http.ResponseWriter, r *http.Request) {
	runID := r.PathValue("run_id")
	data, err := s.cf.Audit.Export(runID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=audit-%s.json", runID))
	w.Write(data)
}

func (s *Server) handleAuditStats(w http.ResponseWriter, r *http.Request) {
	runID := r.PathValue("run_id")
	stats, err := s.cf.Audit.Stats(runID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, stats)
}

// ==================== MCP Gateway ====================

func (s *Server) handleMCPList(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.cf.MCPGateway.List())
}

func (s *Server) handleMCPRegister(w http.ResponseWriter, r *http.Request) {
	var td mcpgw.ToolDescriptor
	if err := json.NewDecoder(r.Body).Decode(&td); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	s.cf.MCPGateway.Register(td)
	writeJSON(w, http.StatusCreated, td)
}

func (s *Server) handleMCPCall(w http.ResponseWriter, r *http.Request) {
	var req mcpgw.CallRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	result, err := s.cf.MCPGateway.Call(r.Context(), req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (s *Server) handleMCPExecuteApproved(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	result, err := s.cf.MCPGateway.ExecuteApproved(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (s *Server) handleMCPReject(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req struct{ Reason string `json:"reason"` }
	json.NewDecoder(r.Body).Decode(&req)
	if err := s.cf.MCPGateway.RejectAction(id, req.Reason); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "rejected"})
}

func (s *Server) handleMCPVet(w http.ResponseWriter, r *http.Request) {
	tool := r.PathValue("tool")
	if err := s.cf.MCPGateway.Vet(tool); err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "vetted", "tool": tool})
}

func (s *Server) handleMCPStats(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.cf.MCPGateway.Stats())
}

// ==================== Sharing ====================

func (s *Server) handleSharingListEngagements(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		userID = "anonymous"
	}
	engagements, err := s.cf.Sharing.ListEngagements(userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, engagements)
}

func (s *Server) handleSharingCreateEngagement(w http.ResponseWriter, r *http.Request) {
	var eng sharing.Engagement
	if err := json.NewDecoder(r.Body).Decode(&eng); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	eng.OwnerID = r.Header.Get("X-User-ID")
	if eng.OwnerID == "" {
		eng.OwnerID = "anonymous"
	}
	if err := s.cf.Sharing.CreateEngagement(&eng); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, eng)
}

func (s *Server) handleSharingGetEngagement(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	eng, err := s.cf.Sharing.GetEngagement(id)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, eng)
}

func (s *Server) handleSharingCreateLink(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req struct {
		Role  string `json:"role"`
		Label string `json:"label"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	createdBy := r.Header.Get("X-User-ID")
	if createdBy == "" {
		createdBy = "anonymous"
	}
	link, err := s.cf.Sharing.CreateShareLink(id, createdBy, sharing.Role(req.Role), req.Label, nil)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, link)
}

func (s *Server) handleSharingListLinks(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	links, err := s.cf.Sharing.ListShareLinks(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, links)
}

func (s *Server) handleSharingRevokeLink(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	linkID := r.PathValue("linkID")
	revokedBy := r.Header.Get("X-User-ID")
	if revokedBy == "" {
		revokedBy = "anonymous"
	}
	if err := s.cf.Sharing.RevokeShareLink(linkID, revokedBy); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "revoked", "engagement": id})
}

func (s *Server) handleSharingAccept(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Token    string `json:"token"`
		Username string `json:"username"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		userID = "anonymous"
	}
	collab, err := s.cf.Sharing.AcceptShareLink(req.Token, userID, req.Username)
	if err != nil {
		writeError(w, http.StatusForbidden, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, collab)
}

func (s *Server) handleSharingListCollaborators(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	collabs, err := s.cf.Sharing.ListCollaborators(id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, collabs)
}

func (s *Server) handleSharingRemoveCollab(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	userID := r.PathValue("userID")
	removedBy := r.Header.Get("X-User-ID")
	if removedBy == "" {
		removedBy = "anonymous"
	}
	if err := s.cf.Sharing.RemoveCollaborator(id, userID, removedBy); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "removed"})
}

func (s *Server) handleSharingAddRun(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	runID := r.PathValue("runID")
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		userID = "anonymous"
	}
	if err := s.cf.Sharing.AddRunToEngagement(id, runID, userID); err != nil {
		writeError(w, http.StatusForbidden, err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "added"})
}
