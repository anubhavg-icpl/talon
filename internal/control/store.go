// Package control is the HTTP layer in front of the agent orchestrator,
// exposing start/status/resume/monitor routes for a long-running pentest
// validation session.
package control

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/anubhavg-icpl/talon/internal/core"
)

// Session is one run's state.
type Session struct {
	Status           string
	Output           string
	PendingInterrupt *core.PendingInterrupt
	RunInput         core.RunInput
	History          []string
	ToolLog          []core.ToolCallRecord
	// StartedAt is when the run was created (UTC), used for listing/sorting.
	StartedAt time.Time
	// EndedAt is when the run reached a terminal state (completed/error), UTC.
	// Zero while the run is still initializing/running/awaiting_approval. Lets
	// the UI freeze the elapsed timer instead of counting forever.
	EndedAt time.Time
	// JudgeVerdict is set when a completed run included a judge assessment.
	// Only meaningful when Status is "completed".
	JudgeVerdict bool
	// JudgeSet is true when JudgeVerdict was populated (false means "no
	// verdict yet / judge skipped", not "judge said false").
	JudgeSet bool
	// Findings are structured security findings (CyberStrike-inspired).
	Findings []core.Finding
	// Report is the multi-section structured validation report.
	Report *core.StructuredReport
	// KillChain is derived attack-path analysis.
	KillChain *core.KillChainAnalysis
	// Methodology is stage coverage.
	Methodology *core.MethodologyState
	// Notes are operator annotations (HITL comments, engagement context).
	Notes []core.OperatorNote
}

// Store is a thread-safe (RWMutex-protected) in-memory session table.
// Persistence: Postgres when EnablePostgres succeeded, else JSON file when
// EnablePersistence was called — every mutation is flushed so runs survive
// restarts.
type Store struct {
	mu       sync.RWMutex
	sessions map[string]*Session
	// persistPath, when non-empty, is the JSON file runs are flushed to.
	persistPath string
	// pg/pgCtx, when set, flush to Postgres instead of JSON.
	pg    *DB
	pgCtx context.Context
}

func NewStore() *Store {
	return &Store{sessions: make(map[string]*Session)}
}

// Create starts a new session in the "initializing" state.
func (s *Store) Create(runID string, input core.RunInput) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[runID] = &Session{
		Status:    "initializing",
		RunInput:  input,
		History:   []string{historyLine("created target=%s", input.TargetIP)},
		StartedAt: time.Now().UTC(),
	}
	s.saveLocked(runID)
}

// Get returns a copy of the session's current fields, or ok=false if unknown.
// Memory (active runs) first, then Postgres (historical runs) when enabled.
func (s *Store) Get(runID string) (Session, bool) {
	s.mu.RLock()
	pg := s.pg
	pgCtx := s.pgCtx
	sess, ok := s.sessions[runID]
	s.mu.RUnlock()
	if !ok && pg != nil {
		if row, found := pg.getRun(pgCtx, runID); found {
			sess = row
			ok = true
		}
	}
	if !ok {
		return Session{}, false
	}
	out := *sess
	if sess.History != nil {
		out.History = append([]string(nil), sess.History...)
	}
	if sess.ToolLog != nil {
		out.ToolLog = append([]core.ToolCallRecord(nil), sess.ToolLog...)
	}
	if sess.PendingInterrupt != nil {
		pi := *sess.PendingInterrupt
		out.PendingInterrupt = &pi
	}
	if sess.Findings != nil {
		out.Findings = append([]core.Finding(nil), sess.Findings...)
	}
	if sess.Report != nil {
		r := *sess.Report
		if sess.Report.Findings != nil {
			r.Findings = append([]core.Finding(nil), sess.Report.Findings...)
		}
		out.Report = &r
	}
	if sess.KillChain != nil {
		kc := *sess.KillChain
		out.KillChain = &kc
	}
	if sess.Methodology != nil {
		m := *sess.Methodology
		out.Methodology = &m
	}
	if sess.Notes != nil {
		out.Notes = append([]core.OperatorNote(nil), sess.Notes...)
	}
	return out, true
}

// GlobalFinding is a finding with run context for the global findings view.
type GlobalFinding struct {
	RunID  string       `json:"run_id"`
	Target string       `json:"target"`
	Finding core.Finding `json:"finding"`
}

// AllFindings returns recent findings across runs (newest runs first).
func (s *Store) AllFindings(limit int, severity string) []GlobalFinding {
	if limit <= 0 {
		limit = 100
	}
	runs, _, _ := s.PaginatedList(200, 0)
	var out []GlobalFinding
	for _, rs := range runs {
		sess, ok := s.Get(rs.RunID)
		if !ok {
			continue
		}
		for _, f := range sess.Findings {
			if severity != "" && strings.ToLower(f.Severity) != severity {
				continue
			}
			out = append(out, GlobalFinding{RunID: rs.RunID, Target: rs.Target, Finding: f})
			if len(out) >= limit {
				return out
			}
		}
	}
	return out
}

// AddNote appends an operator note to a run.
func (s *Store) AddNote(runID, author, body string) (core.OperatorNote, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	sess, ok := s.sessions[runID]
	if !ok && s.pg != nil {
		if row, found := s.pg.getRun(s.pgCtx, runID); found {
			s.sessions[runID] = row
			sess = row
			ok = true
		}
	}
	if !ok || sess == nil {
		return core.OperatorNote{}, fmt.Errorf("run not found")
	}
	body = strings.TrimSpace(body)
	if body == "" {
		return core.OperatorNote{}, fmt.Errorf("note body required")
	}
	n := core.OperatorNote{
		ID:        fmt.Sprintf("NOTE-%03d", len(sess.Notes)+1),
		Author:    strings.TrimSpace(author),
		Body:      body,
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
	sess.Notes = append(sess.Notes, n)
	sess.History = append(sess.History, historyLine("note %s by %s", n.ID, n.Author))
	s.saveLocked(runID)
	return n, nil
}

// Notes returns a copy of operator notes for a run.
func (s *Store) Notes(runID string) ([]core.OperatorNote, bool) {
	sess, ok := s.Get(runID)
	if !ok {
		return nil, false
	}
	return append([]core.OperatorNote(nil), sess.Notes...), true
}

// Snapshot builds a core.RunSnapshot for compare/export.
func (s *Store) Snapshot(runID string) (core.RunSnapshot, bool) {
	sess, ok := s.Get(runID)
	if !ok {
		return core.RunSnapshot{}, false
	}
	var judge *bool
	if sess.JudgeSet {
		v := sess.JudgeVerdict
		judge = &v
	}
	return core.BuildSnapshot(sess.RunInput, sess.Findings, sess.ToolLog, judge, sess.KillChain, sess.Methodology), true
}

// ExportBundle builds a portable export package.
func (s *Store) ExportBundle(runID string) (core.ExportBundle, bool) {
	sess, ok := s.Get(runID)
	if !ok {
		return core.ExportBundle{}, false
	}
	snap, _ := s.Snapshot(runID)
	md := ""
	if sess.Report != nil {
		md = sess.Report.Markdown
	}
	return core.ExportBundle{
		Version:  "1.0",
		RunID:    runID,
		Snapshot: snap,
		ReportMD: md,
		History:  append([]string(nil), sess.History...),
		ToolLog:  append([]core.ToolCallRecord(nil), sess.ToolLog...),
		Notes:    append([]core.OperatorNote(nil), sess.Notes...),
	}, true
}

// IntelEvent is one row in the global intel feed.
type IntelEvent struct {
	At     string `json:"at"`
	RunID  string `json:"run_id"`
	Target string `json:"target"`
	Kind   string `json:"kind"` // finding | history | tool
	Label  string `json:"label"`
	Detail string `json:"detail,omitempty"`
	Severity string `json:"severity,omitempty"`
}

// IntelFeed returns recent cross-run intel (newest first).
func (s *Store) IntelFeed(limit int) []IntelEvent {
	if limit <= 0 {
		limit = 50
	}
	runs, _, _ := s.PaginatedList(100, 0)
	var out []IntelEvent
	for _, rs := range runs {
		sess, ok := s.Get(rs.RunID)
		if !ok {
			continue
		}
		// Recent findings
		for i := len(sess.Findings) - 1; i >= 0; i-- {
			f := sess.Findings[i]
			at := ""
			if !f.CreatedAt.IsZero() {
				at = f.CreatedAt.UTC().Format(time.RFC3339)
			} else {
				at = sess.StartedAt.UTC().Format(time.RFC3339)
			}
			out = append(out, IntelEvent{
				At: at, RunID: rs.RunID, Target: rs.Target,
				Kind: "finding", Label: f.Title, Detail: f.Description, Severity: f.Severity,
			})
		}
		// Last few history lines
		start := len(sess.History) - 3
		if start < 0 {
			start = 0
		}
		for _, h := range sess.History[start:] {
			out = append(out, IntelEvent{
				At: sess.StartedAt.UTC().Format(time.RFC3339),
				RunID: rs.RunID, Target: rs.Target,
				Kind: "history", Label: h,
			})
		}
		if len(out) >= limit*2 {
			break
		}
	}
	// Sort by At desc (string RFC3339 sorts)
	sort.Slice(out, func(i, j int) bool { return out[i].At > out[j].At })
	if len(out) > limit {
		out = out[:limit]
	}
	return out
}

// TriageFinding updates a finding's status on a completed run.
func (s *Store) TriageFinding(runID, findingID, status, duplicateOf string) (core.Finding, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	sess, ok := s.sessions[runID]
	if !ok {
		// try pg hydrate into memory for mutation
		if s.pg != nil {
			if row, found := s.pg.getRun(s.pgCtx, runID); found {
				s.sessions[runID] = row
				sess = row
				ok = true
			}
		}
	}
	if !ok || sess == nil {
		return core.Finding{}, fmt.Errorf("run not found")
	}
	status = strings.ToLower(strings.TrimSpace(status))
	switch status {
	case core.FindingStatusApproved, core.FindingStatusDup, core.FindingStatusOpen,
		"fixed", "ignored", core.FindingStatusNew:
	default:
		return core.Finding{}, fmt.Errorf("invalid status %q", status)
	}
	for i := range sess.Findings {
		if sess.Findings[i].ID == findingID {
			sess.Findings[i].Status = status
			if status == core.FindingStatusDup && duplicateOf != "" &&
				!strings.Contains(sess.Findings[i].Description, "duplicate_of=") {
				sess.Findings[i].Description += " [duplicate_of=" + duplicateOf + "]"
			}
			f := sess.Findings[i]
			s.saveLocked(runID)
			return f, nil
		}
	}
	return core.Finding{}, fmt.Errorf("finding %s not found", findingID)
}

// SetStatus updates just the status field.
func (s *Store) SetStatus(runID, status string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if sess, ok := s.sessions[runID]; ok {
		if sess.Status != status {
			sess.Status = status
			sess.History = append(sess.History, historyLine("status=%s", status))
			s.saveLocked(runID)
		}
	}
}

// SetResult records a run's outcome: history, tool log, and either a pending
// interrupt (status "awaiting_approval") or a final output (status
// "completed").
func (s *Store) SetResult(runID string, result core.RunResult) {
	s.mu.Lock()
	defer s.mu.Unlock()
	sess, ok := s.sessions[runID]
	if !ok {
		return
	}
	defer s.saveLocked(runID)
	// Replace (do not append): the orchestrator tracker already holds the
	// full run log across HITL resume cycles. Appending duplicated entries
	// on every interrupt and made tool counts look inflated.
	if result.ToolLog != nil {
		sess.ToolLog = append([]core.ToolCallRecord(nil), result.ToolLog...)
	}
	if result.Interrupted {
		sess.Status = "awaiting_approval"
		sess.PendingInterrupt = result.Interrupt
		tool := ""
		if result.Interrupt != nil {
			tool = result.Interrupt.ToolName
		}
		sess.History = append(sess.History, historyLine("awaiting_approval tool=%s tools=%d", tool, len(sess.ToolLog)))
		return
	}
	sess.Status = "completed"
	sess.Output = result.FinalMessage
	sess.PendingInterrupt = nil
	sess.EndedAt = time.Now().UTC()
	sess.JudgeVerdict = result.JudgeVerdict
	// Prefer explicit JudgeSet. Also accept legacy callers that only set
	// JudgeVerdict=true, and report-embedded verdicts from finalizeResult.
	sess.JudgeSet = result.JudgeSet || result.JudgeVerdict
	if result.Report != nil && result.Report.JudgeVerdict != nil {
		sess.JudgeSet = true
		sess.JudgeVerdict = *result.Report.JudgeVerdict
	}
	if result.Findings != nil {
		sess.Findings = append([]core.Finding(nil), result.Findings...)
	}
	if result.Report != nil {
		r := *result.Report
		if result.Report.Findings != nil {
			r.Findings = append([]core.Finding(nil), result.Report.Findings...)
		}
		sess.Report = &r
	}
	if result.KillChain != nil {
		kc := *result.KillChain
		sess.KillChain = &kc
	}
	if result.Methodology != nil {
		m := *result.Methodology
		sess.Methodology = &m
	}
	// Prefer structured report markdown in Output when agent text is thin.
	if sess.Report != nil && sess.Report.Markdown != "" &&
		(sess.Output == "" || len(sess.Output) < 80) {
		sess.Output = sess.Report.Markdown
	}
	nFind := len(sess.Findings)
	sess.History = append(sess.History, historyLine("completed judge=%v tools=%d findings=%d", result.JudgeVerdict, len(sess.ToolLog), nFind))
}

// SetError records a run's failure.
func (s *Store) SetError(runID string, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if sess, ok := s.sessions[runID]; ok {
		sess.Status = "error"
		sess.Output = err.Error()
		sess.PendingInterrupt = nil
		sess.EndedAt = time.Now().UTC()
		sess.History = append(sess.History, historyLine("error: %s", err.Error()))
		s.saveLocked(runID)
	}
}

// SetToolLog updates the live tool log while a run is still "running" so
// operators can poll progress without waiting for HITL/completion.
func (s *Store) SetToolLog(runID string, toolLog []core.ToolCallRecord) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if sess, ok := s.sessions[runID]; ok {
		prev := len(sess.ToolLog)
		sess.ToolLog = append([]core.ToolCallRecord(nil), toolLog...)
		// Record only newly completed tools so history stays readable.
		for i := prev; i < len(sess.ToolLog); i++ {
			rec := sess.ToolLog[i]
			sess.History = append(sess.History, historyLine("tool[%d]=%s", rec.Index, rec.ToolName))
		}
		s.saveLocked(runID)
	}
}

// SetFindings updates mid-run structured findings (from report_finding tool).
func (s *Store) SetFindings(runID string, findings []core.Finding) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if sess, ok := s.sessions[runID]; ok {
		prev := len(sess.Findings)
		sess.Findings = append([]core.Finding(nil), findings...)
		if len(findings) > prev {
			sess.History = append(sess.History, historyLine("findings=%d", len(findings)))
		}
		s.saveLocked(runID)
	}
}

// ClaimInterrupt atomically takes ownership of a pending interrupt for
// resume. Returns a snapshot of the session and true if there was a
// pending interrupt; false if another resume already claimed it or none
// exists. Prevents double-resume races after ClearInterrupt.
func (s *Store) ClaimInterrupt(runID string) (Session, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	sess, ok := s.sessions[runID]
	if !ok || sess.PendingInterrupt == nil {
		return Session{}, false
	}
	out := *sess
	if sess.PendingInterrupt != nil {
		pi := *sess.PendingInterrupt
		out.PendingInterrupt = &pi
	}
	if sess.History != nil {
		out.History = append([]string(nil), sess.History...)
	}
	if sess.ToolLog != nil {
		out.ToolLog = append([]core.ToolCallRecord(nil), sess.ToolLog...)
	}
	sess.PendingInterrupt = nil
	sess.Status = "running"
	sess.History = append(sess.History, historyLine("resume_claimed"))
	s.saveLocked(runID)
	return out, true
}

// AppendHistory records an operator-visible event for a run.
func (s *Store) AppendHistory(runID, format string, args ...any) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if sess, ok := s.sessions[runID]; ok {
		sess.History = append(sess.History, historyLine(format, args...))
		s.saveLocked(runID)
	}
}

// ToolLog returns the accumulated tool-call log for one run, served by
// GET /monitor/tools?run_id=... (per-run rather than a single global log,
// since ToolCallRecord lives on RunResult per-run).
func (s *Store) ToolLog(runID string) ([]core.ToolCallRecord, bool) {
	// Use Get so Postgres-hydrated completed runs (not only hot in-memory
	// sessions) still expose their tool log to GET /monitor/tools and Assist.
	sess, ok := s.Get(runID)
	if !ok {
		return nil, false
	}
	if sess.ToolLog == nil {
		return []core.ToolCallRecord{}, true
	}
	return append([]core.ToolCallRecord(nil), sess.ToolLog...), true
}

func historyLine(format string, args ...any) string {
	return time.Now().UTC().Format(time.RFC3339) + " " + fmt.Sprintf(format, args...)
}

// RunSummary is the per-run listing record served by GET /runs.
type RunSummary struct {
	RunID         string    `json:"run_id"`
	Target        string    `json:"target"`
	CVEID         string    `json:"cve_id,omitempty"`
	ServiceName   string    `json:"service_name,omitempty"`
	Status        string    `json:"status"`
	JudgeVerdict  *bool     `json:"judge_verdict,omitempty"`
	ToolCalls     int       `json:"tool_calls"`
	FindingsCount int       `json:"findings_count"`
	AgentMode     string    `json:"agent_mode,omitempty"`
	StartedAt     time.Time `json:"started_at"`
	// EndedAt is the terminal-state timestamp; omitted (zero) while active.
	EndedAt time.Time `json:"ended_at,omitzero"`
}

// List returns one summary per run, newest first.
func (s *Store) List() []RunSummary {
	runs, _, _ := s.PaginatedList(0, 0)
	return runs
}

// PaginatedList returns one page of summaries plus the total count.
// With Postgres enabled the page is produced in SQL (scales to any registry
// size); otherwise it slices the in-memory table. limit <= 0 means "all"
// (memory mode only; PG mode clamps to a sane page size).
func (s *Store) PaginatedList(limit, offset int) ([]RunSummary, int, error) {
	s.mu.RLock()
	pg := s.pg
	pgCtx := s.pgCtx
	s.mu.RUnlock()

	if pg != nil {
		if limit <= 0 {
			limit = 500
		}
		return pg.listRunsPaginated(pgCtx, limit, offset)
	}

	s.mu.RLock()
	defer s.mu.RUnlock()
	all := make([]RunSummary, 0, len(s.sessions))
	for id, sess := range s.sessions {
		sum := RunSummary{
			RunID:         id,
			Target:        sess.RunInput.TargetIP,
			CVEID:         sess.RunInput.CVEID,
			ServiceName:   sess.RunInput.ServiceName,
			Status:        sess.Status,
			ToolCalls:     len(sess.ToolLog),
			FindingsCount: len(sess.Findings),
			AgentMode:     sess.RunInput.AgentMode,
			StartedAt:     sess.StartedAt,
			EndedAt:       sess.EndedAt,
		}
		if sess.JudgeSet {
			v := sess.JudgeVerdict
			sum.JudgeVerdict = &v
		}
		all = append(all, sum)
	}
	sort.Slice(all, func(i, j int) bool {
		if all[i].StartedAt.Equal(all[j].StartedAt) {
			return all[i].RunID > all[j].RunID // stable tie-break (newest id preference)
		}
		return all[i].StartedAt.After(all[j].StartedAt)
	})
	total := len(all)
	if limit > 0 {
		if offset > len(all) {
			offset = len(all)
		}
		end := offset + limit
		if end > len(all) {
			end = len(all)
		}
		all = all[offset:end]
	}
	return all, total, nil
}

// Summary aggregates registry counts for the overview dashboard.
func (s *Store) Summary() RunsSummary {
	s.mu.RLock()
	pg := s.pg
	pgCtx := s.pgCtx
	s.mu.RUnlock()

	if pg != nil {
		sum, err := pg.runsSummary(pgCtx)
		if err == nil {
			return sum
		}
		log.Printf("control: pg summary: %v", err)
	}

	s.mu.RLock()
	defer s.mu.RUnlock()
	var sum RunsSummary
	for _, sess := range s.sessions {
		sum.Total++
		switch sess.Status {
		case "initializing", "running", "awaiting_approval":
			sum.Active++
		case "completed":
			sum.Completed++
		case "error":
			sum.Errored++
		}
		if sess.Status == "awaiting_approval" {
			sum.AwaitingApproval++
		}
		if sess.JudgeSet && sess.JudgeVerdict {
			sum.Compromised++
		}
	}
	return sum
}
