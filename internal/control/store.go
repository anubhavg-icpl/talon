// Package control is the HTTP layer in front of the agent orchestrator,
// exposing start/status/resume/monitor routes for a long-running pentest
// validation session.
package control

import (
	"sync"

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
	s.sessions[runID] = &Session{Status: "initializing", RunInput: input}
}

// Get returns a copy of the session's current fields, or ok=false if unknown.
func (s *Store) Get(runID string) (Session, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	sess, ok := s.sessions[runID]
	if !ok {
		return Session{}, false
	}
	return *sess, true
}

// SetStatus updates just the status field.
func (s *Store) SetStatus(runID, status string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if sess, ok := s.sessions[runID]; ok {
		sess.Status = status
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
	sess.ToolLog = append(sess.ToolLog, result.ToolLog...)
	if result.Interrupted {
		sess.Status = "awaiting_approval"
		sess.PendingInterrupt = result.Interrupt
		return
	}
	sess.Status = "completed"
	sess.Output = result.FinalMessage
	sess.PendingInterrupt = nil
}

// SetError records a run's failure.
func (s *Store) SetError(runID string, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if sess, ok := s.sessions[runID]; ok {
		sess.Status = "error"
		sess.Output = err.Error()
	}
}

// ClearInterrupt drops the pending interrupt after a resume decision has
// been accepted.
func (s *Store) ClearInterrupt(runID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if sess, ok := s.sessions[runID]; ok {
		sess.PendingInterrupt = nil
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
