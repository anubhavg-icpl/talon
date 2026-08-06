// Package audit implements a tamper-evident compliance audit trail for the
// Talon pentest platform. Every significant action taken during an engagement —
// by a human operator, the autonomous agent, or the system itself — is recorded
// as an immutable AuditEntry keyed to a run.
//
// This is adapted from Cloudflare OS observers (structured observers that
// record agent/tool activity for compliance and replay), reinterpreted as a
// dedicated, queryable compliance log. The store is SQL-backed and safe for
// concurrent use, supporting filtered queries by actor and severity, JSON
// export for incident hand-off, and severity roll-up statistics.
package audit

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Actor identifies who or what performed an audited action.
type Actor string

const (
	// ActorUser marks actions performed by a human operator.
	ActorUser Actor = "user"
	// ActorAgent marks actions performed by the autonomous agent/LLM.
	ActorAgent Actor = "agent"
	// ActorSystem marks platform-initiated actions (scheduler, hooks, etc.).
	ActorSystem Actor = "system"
)

// Severity classifies the importance/risk of an audit event.
type Severity string

const (
	// SeverityInfo is routine informational activity (e.g. run started).
	SeverityInfo Severity = "info"
	// SeverityLow is minor noteworthy events (e.g. config change).
	SeverityLow Severity = "low"
	// SeverityMedium is events warranting review (e.g. scope boundary approach).
	SeverityMedium Severity = "medium"
	// SeverityHigh is significant security-relevant events (e.g. exploit fired).
	SeverityHigh Severity = "high"
	// SeverityCritical is severe events requiring immediate attention.
	SeverityCritical Severity = "critical"
)

// AuditEntry is a single immutable record in the compliance audit trail.
type AuditEntry struct {
	// ID is the unique identifier for the entry (UUID when auto-generated).
	ID string `json:"id"`
	// RunID scopes the entry to a specific engagement/run.
	RunID string `json:"run_id"`
	// Actor is who performed the action: user, agent, or system.
	Actor Actor `json:"actor"`
	// Action is a short verb describing what happened (e.g. "tool.invoke",
	// "finding.approve", "run.resume").
	Action string `json:"action"`
	// ResourceType is the kind of object acted upon (e.g. "tool", "finding",
	// "run", "blueprint").
	ResourceType string `json:"resource_type"`
	// ResourceID is the identifier of the affected resource, if any.
	ResourceID string `json:"resource_id,omitempty"`
	// Details is arbitrary structured context for the event, stored as JSON.
	// Use json.RawMessage for pre-encoded payloads or a marshalable value with
	// AuditStore.LogWithDetails.
	Details json.RawMessage `json:"details,omitempty"`
	// IPAddress is the originating IP of a user-initiated action, when known.
	IPAddress string `json:"ip_address,omitempty"`
	// Timestamp is the UTC time the event occurred.
	Timestamp time.Time `json:"timestamp"`
	// Severity classifies the event's importance (info/low/medium/high/critical).
	Severity Severity `json:"severity"`
}

// AuditStore is a SQL-backed append-only audit log. It is safe for concurrent
// use. Entries are never mutated or deleted in the normal course of operation;
// the Log method is the only write path.
type AuditStore struct {
	db *sql.DB
	mu sync.Mutex
}

// auditSchema creates the audit_entries table. Severity and actor are
// constrained by CHECK clauses to the documented vocabularies. Indexes cover
// the common filtered queries (by run, by run+actor, by run+severity). Uses
// only portable SQLite-compatible DDL.
const auditSchema = `
CREATE TABLE IF NOT EXISTS audit_entries (
	id            TEXT PRIMARY KEY,
	run_id        TEXT NOT NULL,
	actor         TEXT NOT NULL CHECK(actor IN ('user','agent','system')),
	action        TEXT NOT NULL,
	resource_type TEXT NOT NULL DEFAULT '',
	resource_id   TEXT NOT NULL DEFAULT '',
	details       TEXT NOT NULL DEFAULT '',
	ip_address    TEXT NOT NULL DEFAULT '',
	timestamp     TEXT NOT NULL,
	severity      TEXT NOT NULL CHECK(severity IN ('info','low','medium','high','critical'))
);
CREATE INDEX IF NOT EXISTS idx_audit_run        ON audit_entries(run_id);
CREATE INDEX IF NOT EXISTS idx_audit_run_actor  ON audit_entries(run_id, actor);
CREATE INDEX IF NOT EXISTS idx_audit_run_sev    ON audit_entries(run_id, severity);
CREATE INDEX IF NOT EXISTS idx_audit_timestamp  ON audit_entries(timestamp);
`

// NewAuditStore wraps an already-open *sql.DB. The caller is responsible for
// registering the SQL driver (e.g. sqlite3) and for the connection's lifetime.
// Call Migrate once before first use to create the schema.
func NewAuditStore(db *sql.DB) *AuditStore {
	return &AuditStore{db: db}
}

// Migrate creates the audit_entries table and indexes if they do not already
// exist. It is idempotent and safe to call on every startup.
func (s *AuditStore) Migrate() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, err := s.db.Exec(auditSchema)
	if err != nil {
		return fmt.Errorf("audit: migrate: %w", err)
	}
	return nil
}

// Log appends entry to the audit trail. If ID is empty a UUID is generated;
// if Timestamp is zero it is set to the current UTC time. Details is stored as
// the raw JSON bytes supplied (empty when none provided).
func (s *AuditStore) Log(entry AuditEntry) error {
	if strings.TrimSpace(entry.ID) == "" {
		entry.ID = uuid.NewString()
	}
	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now().UTC()
	}
	if entry.Severity == "" {
		entry.Severity = SeverityInfo
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	_, err := s.db.Exec(
		`INSERT INTO audit_entries
		    (id, run_id, actor, action, resource_type, resource_id, details, ip_address, timestamp, severity)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		entry.ID, entry.RunID, string(entry.Actor), entry.Action,
		entry.ResourceType, entry.ResourceID, string(entry.Details),
		entry.IPAddress, entry.Timestamp.UTC().Format(time.RFC3339),
		string(entry.Severity),
	)
	if err != nil {
		return fmt.Errorf("audit: log: %w", err)
	}
	return nil
}

// List returns every audit entry for a run, oldest first.
func (s *AuditStore) List(runID string) ([]AuditEntry, error) {
	return s.queryEntries(
		`SELECT id, run_id, actor, action, resource_type, resource_id, details, ip_address, timestamp, severity
		   FROM audit_entries WHERE run_id = ? ORDER BY timestamp ASC`, runID)
}

// ListByActor returns audit entries for a run filtered by actor, oldest first.
func (s *AuditStore) ListByActor(runID, actor string) ([]AuditEntry, error) {
	return s.queryEntries(
		`SELECT id, run_id, actor, action, resource_type, resource_id, details, ip_address, timestamp, severity
		   FROM audit_entries WHERE run_id = ? AND actor = ? ORDER BY timestamp ASC`,
		runID, actor)
}

// ListBySeverity returns audit entries for a run filtered by severity, oldest
// first.
func (s *AuditStore) ListBySeverity(runID, severity string) ([]AuditEntry, error) {
	return s.queryEntries(
		`SELECT id, run_id, actor, action, resource_type, resource_id, details, ip_address, timestamp, severity
		   FROM audit_entries WHERE run_id = ? AND severity = ? ORDER BY timestamp ASC`,
		runID, severity)
}

// Export serializes all audit entries for a run as pretty-printed JSON, suitable
// for incident hand-off or long-term archival.
func (s *AuditStore) Export(runID string) ([]byte, error) {
	entries, err := s.List(runID)
	if err != nil {
		return nil, fmt.Errorf("audit: export %q: %w", runID, err)
	}
	payload := struct {
		RunID   string       `json:"run_id"`
		Count   int          `json:"count"`
		Entries []AuditEntry `json:"entries"`
	}{
		RunID:   runID,
		Count:   len(entries),
		Entries: entries,
	}
	data, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("audit: export %q: marshal: %w", runID, err)
	}
	return data, nil
}

// Stats returns a count of audit entries for a run grouped by severity. The
// returned map always contains keys for every defined severity (info/low/
// medium/high/critical) plus a "total" key, so callers can render a complete
// table even when some severities have zero events.
func (s *AuditStore) Stats(runID string) (map[string]int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	stats := map[string]int{
		"info":     0,
		"low":      0,
		"medium":   0,
		"high":     0,
		"critical": 0,
		"total":    0,
	}
	rows, err := s.db.Query(
		`SELECT severity, COUNT(*) FROM audit_entries WHERE run_id = ? GROUP BY severity`,
		runID)
	if err != nil {
		return nil, fmt.Errorf("audit: stats %q: %w", runID, err)
	}
	defer rows.Close()

	for rows.Next() {
		var sev string
		var n int
		if err := rows.Scan(&sev, &n); err != nil {
			return nil, fmt.Errorf("audit: stats %q: %w", runID, err)
		}
		stats[sev] = n
		stats["total"] += n
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("audit: stats %q: %w", runID, err)
	}
	return stats, nil
}

// queryEntries runs a parameterized SELECT and decodes all matching rows. It
// centralizes row decoding for the List/ListByActor/ListBySeverity methods.
func (s *AuditStore) queryEntries(query string, args ...any) ([]AuditEntry, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []AuditEntry
	for rows.Next() {
		e, err := scanEntry(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

// rowScanner is the subset of *sql.Rows/*sql.Row used for decoding a single row.
type rowScanner interface {
	Scan(dest ...any) error
}

// scanEntry decodes one audit row into an AuditEntry. The details column is
// scanned into a string and converted back to json.RawMessage; the timestamp
// (stored as RFC3339 text) is parsed back into time.Time.
func scanEntry(r rowScanner) (AuditEntry, error) {
	var (
		e           AuditEntry
		actor       string
		details     string
		severity    string
		timestamp   string
	)
	if err := r.Scan(
		&e.ID, &e.RunID, &actor, &e.Action,
		&e.ResourceType, &e.ResourceID, &details, &e.IPAddress,
		&timestamp, &severity,
	); err != nil {
		return AuditEntry{}, err
	}
	e.Actor = Actor(actor)
	e.Severity = Severity(severity)
	if details != "" {
		e.Details = json.RawMessage(details)
	}
	parsed, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return AuditEntry{}, fmt.Errorf("decode timestamp: %w", err)
	}
	e.Timestamp = parsed
	return e, nil
}
