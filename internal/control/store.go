// Package control is the HTTP layer in front of the agent orchestrator,
// exposing start/status/resume/monitor routes for a long-running pentest
// validation session.
package control

import (
	"fmt"
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
	// JudgeVerdict is set when a completed run included a judge assessment.
	// Only meaningful when Status is "completed".
	JudgeVerdict bool
	// JudgeSet is true when JudgeVerdict was populated (false means "no
	// verdict yet / judge skipped", not "judge said false").
	JudgeSet bool
}

// Store is a thread-safe (RWMutex-protected) in-memory session table.
type Store struct {
	mu       sync.RWMutex
	sessions map[string]*Session
}

func NewStore() *Store {
	return &Store{sessions: make(map[string]*Session)}
}

// Create starts a new session in the "initializing" state.
func (s *Store) Create(runID string, input core.RunInput) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[runID] = &Session{
		Status:   "initializing",
		RunInput: input,
		History:  []string{historyLine("created target=%s", input.TargetIP)},
	}
}

// Get returns a copy of the session's current fields, or ok=false if unknown.
func (s *Store) Get(runID string) (Session, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	sess, ok := s.sessions[runID]
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
	return out, true
}

// SetStatus updates just the status field.
func (s *Store) SetStatus(runID, status string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if sess, ok := s.sessions[runID]; ok {
		if sess.Status != status {
			sess.Status = status
			sess.History = append(sess.History, historyLine("status=%s", status))
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
	sess.JudgeVerdict = result.JudgeVerdict
	sess.JudgeSet = true
	sess.History = append(sess.History, historyLine("completed judge=%v tools=%d", result.JudgeVerdict, len(sess.ToolLog)))
}

// SetError records a run's failure.
func (s *Store) SetError(runID string, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if sess, ok := s.sessions[runID]; ok {
		sess.Status = "error"
		sess.Output = err.Error()
		sess.PendingInterrupt = nil
		sess.History = append(sess.History, historyLine("error: %s", err.Error()))
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
	return out, true
}

// AppendHistory records an operator-visible event for a run.
func (s *Store) AppendHistory(runID, format string, args ...any) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if sess, ok := s.sessions[runID]; ok {
		sess.History = append(sess.History, historyLine(format, args...))
	}
}

// ToolLog returns the accumulated tool-call log for one run, served by
// GET /monitor/tools?run_id=... (per-run rather than a single global log,
// since ToolCallRecord lives on RunResult per-run).
func (s *Store) ToolLog(runID string) ([]core.ToolCallRecord, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	sess, ok := s.sessions[runID]
	if !ok {
		return nil, false
	}
	return append([]core.ToolCallRecord(nil), sess.ToolLog...), true
}

func historyLine(format string, args ...any) string {
	return time.Now().UTC().Format(time.RFC3339) + " " + fmt.Sprintf(format, args...)
}
