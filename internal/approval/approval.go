// Package approval implements a human-in-the-loop (HITL) action store for the
// Talon pentest platform.
//
// It is adapted from Cloudflare OS's action-store pattern: every potentially
// dangerous tool invocation is recorded as an Action that must be claimed
// (claim-before-dispatch) and resolved (approved/rejected) before it takes
// effect. The store is backed by SQLite via database/sql and uses only
// parameterized queries.
//
// Lifecycle:
//
//	pending  -> applying (Claim, atomic)  -> applied   (Approve)
//	                                  \----> rejected  (Reject)
//	applying -> unknown  (MarkUnknown, crash recovery — non-retryable)
//	pending  -> failed   (execution error after dispatch)
//
// Claim-before-dispatch gives at-most-once semantics: only one dispatcher
// wins the atomic "UPDATE ... WHERE state='pending'" and flips the row to
// 'applying'. If the process crashes after claiming but before resolving, the
// row stays 'applying'; on restart MarkUnknown marks it 'unknown' so it is
// never silently retried.
package approval

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

// ActionState is the lifecycle state of an approval-gated action.
type ActionState string

const (
	// StatePending: created, awaiting a dispatcher to claim it.
	StatePending ActionState = "pending"
	// StateApplying: claimed by a dispatcher, executing or about to.
	StateApplying ActionState = "applying"
	// StateApplied: approved and the result recorded (terminal).
	StateApplied ActionState = "applied"
	// StateRejected: a human rejected the action (terminal).
	StateRejected ActionState = "rejected"
	// StateFailed: dispatch/execution failed (terminal).
	StateFailed ActionState = "failed"
	// StateUnknown: claimed but never resolved (crash recovery, terminal &
	// non-retryable).
	StateUnknown ActionState = "unknown"
)

// IsTerminal reports whether the state is a final, non-retryable state.
func (s ActionState) IsTerminal() bool {
	switch s {
	case StateApplied, StateRejected, StateFailed, StateUnknown:
		return true
	}
	return false
}

// RiskLevel classifies how dangerous an action is. Ordering matters for
// auto-approval: low < medium < high < critical.
type RiskLevel string

const (
	RiskLow      RiskLevel = "low"
	RiskMedium   RiskLevel = "medium"
	RiskHigh     RiskLevel = "high"
	RiskCritical RiskLevel = "critical"
)

// rank returns a comparable ordering for RiskLevel. Unknown levels rank 0
// (treated as the least risky by AutoApprove only when configured, but callers
// should normalize inputs first).
func (r RiskLevel) rank() int {
	switch r {
	case RiskLow:
		return 1
	case RiskMedium:
		return 2
	case RiskHigh:
		return 3
	case RiskCritical:
		return 4
	}
	return 0
}

// Action is one approval-gated tool invocation.
type Action struct {
	// ID is the unique action identifier (caller-provided, e.g. a UUID).
	ID string `json:"id"`
	// RunID is the engagement/run this action belongs to.
	RunID string `json:"run_id"`
	// ToolName is the MCP tool or CLI invoked (e.g. "run_exploit", "sqlmap").
	ToolName string `json:"tool_name"`
	// Args is the raw JSON arguments for the tool call.
	Args json.RawMessage `json:"args"`
	// State is the current lifecycle state.
	State ActionState `json:"state"`
	// RiskLevel is the assessed danger of this action.
	RiskLevel RiskLevel `json:"risk_level"`
	// Summary is a short human-readable description shown in the approval UI.
	Summary string `json:"summary"`
	// Result holds the outcome JSON: the tool result on Approve, a reason on
	// Reject, or the error on Fail.
	Result json.RawMessage `json:"result,omitempty"`
	// CreatedAt is when the action was recorded (UTC, RFC3339).
	CreatedAt time.Time `json:"created_at"`
	// ResolvedAt is set when the action reaches a terminal state (UTC).
	ResolvedAt *time.Time `json:"resolved_at,omitempty"`
	// ClaimedAt is set when a dispatcher claims the action (UTC).
	ClaimedAt *time.Time `json:"claimed_at,omitempty"`
}

// ErrNotFound is returned when an action ID does not exist.
var ErrNotFound = errors.New("approval: action not found")

// ErrNotClaimable is returned by Claim when the action is missing or no longer
// pending (already claimed/resolved).
var ErrNotClaimable = errors.New("approval: action not pending or already claimed")

// ErrAlreadyResolved is returned by Approve/Reject when the action is already
// in a terminal state.
var ErrAlreadyResolved = errors.New("approval: action already resolved")

// ActionStore persists approval-gated actions in SQLite via database/sql.
// It is safe for concurrent use: writes are serialized by a mutex and the
// claim transition is an atomic conditional UPDATE.
type ActionStore struct {
	db *sql.DB
	mu sync.Mutex
	// autoApproveMax, when non-empty, is the highest RiskLevel that AutoApprove
	// will permit without human review (e.g. RiskMedium auto-approves low+medium).
	autoApproveMax RiskLevel
}

// New wraps an already-open *sql.DB and runs the schema migration. Use this
// when you want to share a connection or register the SQLite driver yourself.
func New(db *sql.DB) (*ActionStore, error) {
	s := &ActionStore{db: db}
	if err := s.Migrate(); err != nil {
		return nil, err
	}
	return s, nil
}

// Open creates a new SQLite-backed store at dbPath, registers the schema, and
// returns it. It uses the "sqlite3" driver (github.com/mattn/go-sqlite3), which
// the calling binary must import — typically via a blank import in main:
//
//	import _ "github.com/mattn/go-sqlite3"
//
// The DSN enables WAL journaling and a 5s busy timeout for safe concurrent
// access, matching the rest of the platform's SQLite usage.
func Open(dbPath string) (*ActionStore, error) {
	db, err := sql.Open("sqlite3", dbPath+"?_journal=WAL&_busy_timeout=5000")
	if err != nil {
		return nil, fmt.Errorf("approval: open %s: %w", dbPath, err)
	}
	s, err := New(db)
	if err != nil {
		db.Close()
		return nil, err
	}
	return s, nil
}

// Close closes the underlying database handle.
func (s *ActionStore) Close() error {
	if s.db == nil {
		return nil
	}
	return s.db.Close()
}

// SetAutoApprove configures the maximum RiskLevel eligible for automatic
// approval. An empty level disables auto-approval entirely (the default).
// Example: SetAutoApprove(RiskMedium) auto-approves "low" and "medium".
func (s *ActionStore) SetAutoApprove(max RiskLevel) {
	s.mu.Lock()
	s.autoApproveMax = max
	s.mu.Unlock()
}

// ConfigureAutoApproveFromEnv reads TALON_AUTO_APPROVE_RISK
// (one of low|medium|high|critical) and applies it. Empty/unset disables
// auto-approval. Invalid values are ignored with no effect.
func (s *ActionStore) ConfigureAutoApproveFromEnv() {
	v := strings.ToLower(strings.TrimSpace(os.Getenv("TALON_AUTO_APPROVE_RISK")))
	switch RiskLevel(v) {
	case RiskLow, RiskMedium, RiskHigh, RiskCritical:
		s.SetAutoApprove(RiskLevel(v))
	}
}

// AutoApprove reports whether an action with the given risk level may be
// auto-approved without human review. It returns true only when auto-approval
// is configured (via SetAutoApprove) and riskLevel is at or below the
// configured threshold. Unknown risk levels are never auto-approved.
func (s *ActionStore) AutoApprove(riskLevel RiskLevel) bool {
	s.mu.Lock()
	max := s.autoApproveMax
	s.mu.Unlock()
	if max == "" {
		return false
	}
	r := riskLevel.rank()
	if r == 0 {
		return false // unknown risk — require human review
	}
	return r <= max.rank()
}

// Migrate creates the approval_actions table if it does not already exist.
// It is idempotent and called automatically by New/Open.
func (s *ActionStore) Migrate() error {
	const ddl = `
CREATE TABLE IF NOT EXISTS approval_actions (
    id          TEXT PRIMARY KEY,
    run_id      TEXT NOT NULL,
    tool_name   TEXT NOT NULL,
    args        TEXT NOT NULL,            -- JSON arguments
    state       TEXT NOT NULL,
    risk_level  TEXT NOT NULL,
    summary     TEXT NOT NULL DEFAULT '',
    result      TEXT NOT NULL DEFAULT '', -- JSON result/reason/error
    created_at  TEXT NOT NULL,            -- RFC3339
    resolved_at TEXT,                     -- RFC3339, nullable
    claimed_at  TEXT                      -- RFC3339, nullable
);
CREATE INDEX IF NOT EXISTS idx_approval_actions_run   ON approval_actions(run_id);
CREATE INDEX IF NOT EXISTS idx_approval_actions_state ON approval_actions(state);
`
	_, err := s.db.Exec(ddl)
	if err != nil {
		return fmt.Errorf("approval: migrate: %w", err)
	}
	return nil
}

// Create records a new action in the pending state. The Action.ID, RunID and
// ToolName must be set; State/CreatedAt are filled in if zero. Args is stored
// as-is (an empty Args is serialized as "{}").
func (s *ActionStore) Create(a *Action) error {
	if a == nil {
		return errors.New("approval: nil action")
	}
	if strings.TrimSpace(a.ID) == "" {
		return errors.New("approval: action id required")
	}
	if strings.TrimSpace(a.RunID) == "" {
		return errors.New("approval: action run_id required")
	}
	if strings.TrimSpace(a.ToolName) == "" {
		return errors.New("approval: action tool_name required")
	}
	if a.State == "" {
		a.State = StatePending
	}
	if a.RiskLevel == "" {
		a.RiskLevel = RiskMedium // conservative default for a pentest tool
	}
	if a.CreatedAt.IsZero() {
		a.CreatedAt = time.Now().UTC()
	}
	args := a.Args
	if len(strings.TrimSpace(string(args))) == 0 {
		args = json.RawMessage("{}")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	_, err := s.db.Exec(
		`INSERT INTO approval_actions
		    (id, run_id, tool_name, args, state, risk_level, summary, result, created_at, resolved_at, claimed_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, '', ?, NULL, NULL)`,
		a.ID, a.RunID, a.ToolName, string(args), string(a.State),
		string(a.RiskLevel), a.Summary, a.CreatedAt.UTC().Format(time.RFC3339Nano),
	)
	if err != nil {
		return fmt.Errorf("approval: create %s: %w", a.ID, err)
	}
	return nil
}

// Get returns the action with the given id, or ErrNotFound if it does not exist.
func (s *ActionStore) Get(id string) (*Action, error) {
	row := s.db.QueryRow(
		`SELECT id, run_id, tool_name, args, state, risk_level, summary, result, created_at, resolved_at, claimed_at
		   FROM approval_actions WHERE id = ?`, id)
	a, err := scanAction(row.Scan)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("approval: get %s: %w", id, err)
	}
	return &a, nil
}

// ListPending returns all pending actions, optionally filtered by runID (empty
// runID lists pending actions across all runs). Results are ordered oldest-first.
func (s *ActionStore) ListPending(runID string) ([]Action, error) {
	var (
		query string
		args  []any
	)
	if runID == "" {
		query = `SELECT id, run_id, tool_name, args, state, risk_level, summary, result, created_at, resolved_at, claimed_at
		           FROM approval_actions WHERE state = ? ORDER BY created_at ASC`
		args = []any{string(StatePending)}
	} else {
		query = `SELECT id, run_id, tool_name, args, state, risk_level, summary, result, created_at, resolved_at, claimed_at
		           FROM approval_actions WHERE state = ? AND run_id = ? ORDER BY created_at ASC`
		args = []any{string(StatePending), runID}
	}
	return s.queryActions(query, args...)
}

// ListAll returns all actions for a runID (empty runID lists every action),
// newest-first.
func (s *ActionStore) ListAll(runID string) ([]Action, error) {
	var (
		query string
		args  []any
	)
	if runID == "" {
		query = `SELECT id, run_id, tool_name, args, state, risk_level, summary, result, created_at, resolved_at, claimed_at
		           FROM approval_actions ORDER BY created_at DESC`
	} else {
		query = `SELECT id, run_id, tool_name, args, state, risk_level, summary, result, created_at, resolved_at, claimed_at
		           FROM approval_actions WHERE run_id = ? ORDER BY created_at DESC`
		args = []any{runID}
	}
	return s.queryActions(query, args...)
}

// Claim atomically transitions an action from "pending" to "applying" and
// records the claim time. This implements claim-before-dispatch (at-most-once):
// the conditional UPDATE guarantees exactly one caller wins. Returns the
// updated action. ErrNotClaimable is returned if the action is missing or no
// longer pending.
func (s *ActionStore) Claim(id string) (*Action, error) {
	now := time.Now().UTC().Format(time.RFC3339Nano)

	s.mu.Lock()
	res, err := s.db.Exec(
		`UPDATE approval_actions
		    SET state = ?, claimed_at = ?
		  WHERE id = ? AND state = ?`,
		string(StateApplying), now, id, string(StatePending),
	)
	s.mu.Unlock()
	if err != nil {
		return nil, fmt.Errorf("approval: claim %s: %w", id, err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		// Distinguish "missing" from "not pending" for clearer errors.
		if _, err := s.Get(id); err != nil {
			return nil, ErrNotClaimable
		}
		return nil, ErrNotClaimable
	}
	return s.Get(id)
}

// Approve resolves a claimed/pending action as applied, attaching the
// execution result (JSON). ErrAlreadyResolved is returned if the action is
// already terminal.
func (s *ActionStore) Approve(id string, result json.RawMessage) error {
	return s.resolve(id, StateApplied, result)
}

// Reject resolves an action as rejected, recording the human's reason (stored
// in Result as {"reason": "..."}). ErrAlreadyResolved is returned if the
// action is already terminal.
func (s *ActionStore) Reject(id string, reason string) error {
	payload, err := json.Marshal(map[string]string{"reason": strings.TrimSpace(reason)})
	if err != nil {
		return fmt.Errorf("approval: reject %s: marshal reason: %w", id, err)
	}
	return s.resolve(id, StateRejected, payload)
}

// Fail marks a dispatched action as failed with the given error detail. Used
// when execution throws after Claim. ErrAlreadyResolved is returned if the
// action is already terminal.
func (s *ActionStore) Fail(id string, detail string) error {
	payload, err := json.Marshal(map[string]string{"error": detail})
	if err != nil {
		return fmt.Errorf("approval: fail %s: marshal detail: %w", id, err)
	}
	return s.resolve(id, StateFailed, payload)
}

// MarkUnknown is the crash-recovery transition: an action stuck in "applying"
// (claimed but never resolved) is moved to "unknown" so it is never silently
// retried. This makes the outcome explicit and non-retryable. Pending actions
// are also accepted so callers can sweep the whole table on startup.
func (s *ActionStore) MarkUnknown(id string) error {
	now := time.Now().UTC().Format(time.RFC3339Nano)

	s.mu.Lock()
	res, err := s.db.Exec(
		`UPDATE approval_actions
		    SET state = ?, resolved_at = ?
		  WHERE id = ? AND state IN (?, ?)`,
		string(StateUnknown), now, id, string(StateApplying), string(StatePending),
	)
	s.mu.Unlock()
	if err != nil {
		return fmt.Errorf("approval: mark unknown %s: %w", id, err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		if _, gerr := s.Get(id); gerr != nil {
			return ErrNotFound
		}
		return ErrAlreadyResolved
	}
	return nil
}

// resolve is the shared terminal-transition helper for Approve/Reject/Fail. It
// only transitions from a non-terminal state (pending or applying).
func (s *ActionStore) resolve(id string, state ActionState, result json.RawMessage) error {
	now := time.Now().UTC().Format(time.RFC3339Nano)
	resStr := ""
	if len(result) > 0 {
		resStr = string(result)
	}

	s.mu.Lock()
	res, err := s.db.Exec(
		`UPDATE approval_actions
		    SET state = ?, result = ?, resolved_at = ?
		  WHERE id = ? AND state IN (?, ?)`,
		string(state), resStr, now, id, string(StatePending), string(StateApplying),
	)
	s.mu.Unlock()
	if err != nil {
		return fmt.Errorf("approval: resolve %s to %s: %w", id, state, err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		if _, gerr := s.Get(id); gerr != nil {
			return ErrNotFound
		}
		return ErrAlreadyResolved
	}
	return nil
}

// queryActions runs a SELECT and scans all matching rows into Action values.
func (s *ActionStore) queryActions(query string, args ...any) ([]Action, error) {
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("approval: query actions: %w", err)
	}
	defer rows.Close()

	var out []Action
	for rows.Next() {
		a, err := scanAction(rows.Scan)
		if err != nil {
			return nil, fmt.Errorf("approval: scan action: %w", err)
		}
		out = append(out, a)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("approval: iterate actions: %w", err)
	}
	return out, nil
}

// scanAction decodes one row into an Action. scan is the Scan method of either
// *sql.Row or *sql.Rows (both share the func(dest ...any) error signature).
// Timestamps are stored as RFC3339 text and parsed back into time.Time;
// resolved_at/claimed_at are nullable and become nil pointers when absent.
func scanAction(scan func(dest ...any) error) (Action, error) {
	var (
		a          Action
		createdAt  string
		resolvedAt sql.NullString
		claimedAt  sql.NullString
		args       string
		result     string
	)
	err := scan(
		&a.ID, &a.RunID, &a.ToolName, &args, &a.State, &a.RiskLevel,
		&a.Summary, &result, &createdAt, &resolvedAt, &claimedAt,
	)
	if err != nil {
		return Action{}, err
	}
	if args != "" {
		a.Args = json.RawMessage(args)
	}
	if result != "" {
		a.Result = json.RawMessage(result)
	}
	if t, err := time.Parse(time.RFC3339Nano, createdAt); err == nil {
		a.CreatedAt = t
	}
	if resolvedAt.Valid {
		if t, err := time.Parse(time.RFC3339Nano, resolvedAt.String); err == nil {
			a.ResolvedAt = &t
		}
	}
	if claimedAt.Valid {
		if t, err := time.Parse(time.RFC3339Nano, claimedAt.String); err == nil {
			a.ClaimedAt = &t
		}
	}
	return a, nil
}

// dangerousTools is the set of tool names that should require human approval
// before unattended dispatch in a pentest engagement. It includes both the
// project's MCP tool names (e.g. run_exploit, sqlmap_scan, hydra_attack) and
// the underlying CLI binaries (e.g. sqlmap, msfconsole, hydra, nmap).
//
// nmap is included conservatively: although a basic connect scan is benign,
// unattended scans with aggressive flags (-A, -sV -sC, --script, -T4) are
// noisy, can trigger IDS/IPS, and are exactly the kind of thing a HITL gate
// should review. Callers wanting finer control can inspect Args separately.
var dangerousTools = map[string]bool{
	// --- Exploitation / payload frameworks ---
	"run_exploit":           true,
	"run_auxiliary_module":  true,
	"run_post_module":       true,
	"generate_payload":      true,
	"custom_exploit":        true,
	"metasploit":            true,
	"msfconsole":            true,
	"msfvenom":              true,
	"msfrpc":                true,
	"armitage":              true,
	"cobaltstrike":          true,
	"sliver":                true,

	// --- SQL injection ---
	"sqlmap":      true,
	"sqlmap_scan": true,

	// --- Credential brute force / cracking ---
	"hydra":                        true,
	"hydra_attack":                 true,
	"medusa":                       true,
	"ncrack":                       true,
	"john":                         true,
	"johntheripper":                true,
	"hashcat":                      true,
	"responder":                    true,
	"responder_credential_harvest": true,

	// --- Scanners (aggressive / noisy) ---
	"nmap":                  true, // conservative: review unattended scans
	"nmap_scan":             true,
	"masscan":               true,
	"nikto":                 true,
	"nuclei":                true,
	"nuclei_scan":           true,
	"wpscan":                true,
	"dirb":                  true,
	"gobuster":              true,
	"ffuf":                  true,
	"feroxbuster":           true,
	"rustscan_fast_scan":    true,
	"arp_scan_discovery":    true,

	// --- Web active scanners ---
	"zap":        true,
	"owaspzap":   true,
	"burpsuite":  true,
	"burp":       true,
	"w3af":       true,
	"skipfish":   true,
	"arachni":    true,

	// --- Wireless / MITM ---
	"aircrack-ng":  true,
	"aircrack":     true,
	"reaver":       true,
	"bettercap":    true,
	"ettercap":     true,
	"mitm6":        true,
	"evilgrade":    true,

	// --- Post-exploit / destructive ---
	"mimikatz": true,
}

// IsDangerous reports whether a tool invocation requires human approval before
// unattended dispatch. The check is case-insensitive and matches both the
// project's MCP tool names and common CLI binary names. Tools not in the known
// dangerous set are considered safe-by-default; callers may additionally gate
// on RiskLevel or Args for custom rules.
func IsDangerous(toolName string) bool {
	return dangerousTools[strings.ToLower(strings.TrimSpace(toolName))]
}
