// Package gatekeeper implements capability-based access control for external
// tools and services, adapted from Cloudflare's OS gatekeeper pattern. Each
// external service or tool gets its own Gatekeeper instance that moderates
// access, handles authentication, enforces scope boundaries, and logs every
// privileged action to a durable audit log.
//
// A Gatekeeper mints time-boxed Sessions in exchange for validated
// Capabilities. Every privileged action taken under a session is recorded as
// an ActionLog row, so the full provenance of an engagement can be
// reconstructed after the fact.
//
// A GatekeeperRegistry owns the shared audit database and catalogues the
// configured gatekeepers. Tables are created by Migrate.
//
// The schema targets SQLite (the project's in-process store) and every query
// uses parameterized placeholders to prevent injection.
package gatekeeper

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

// AuthType enumerates the credential schemes a gatekeeper can use to talk to
// its backing service.
type AuthType string

const (
	// AuthTypeOAuth authenticates via OAuth 2.0 bearer tokens. The Credentials
	// map typically holds client_id / client_secret / refresh_token.
	AuthTypeOAuth AuthType = "oauth"
	// AuthTypeAPIKey authenticates via a static API key.
	AuthTypeAPIKey AuthType = "apikey"
	// AuthTypeBasic authenticates via HTTP Basic auth.
	AuthTypeBasic AuthType = "basic"
)

// SessionStatus is the lifecycle state of an access session.
type SessionStatus string

const (
	SessionActive  SessionStatus = "active"
	SessionRevoked SessionStatus = "revoked"
	SessionExpired SessionStatus = "expired"
)

// defaultSessionTTL is applied when a capability does not carry its own
// ExpiresAt.
const defaultSessionTTL = time.Hour

// Capability describes a single grant: which tools may be invoked, under which
// scopes, whether the grant is read-only, and when it lapses.
type Capability struct {
	// Tool is the set of tool names this capability authorizes (e.g.
	// "search_code", "create_issue").
	Tool []string `json:"tool"`
	// Scope is the permission scopes granted (e.g. "repo:read", "dns:write").
	Scope []string `json:"scope"`
	// ReadOnly, when true, restricts the capability to non-mutating actions.
	ReadOnly bool `json:"read_only"`
	// ExpiresAt optionally bounds the capability in time. A nil value defers
	// expiry to the gatekeeper default (defaultSessionTTL).
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

// GatekeeperConfig is the declarative description of a gatekeeper. It is what
// operators register with a GatekeeperRegistry and what List() returns.
type GatekeeperConfig struct {
	// Name uniquely identifies the gatekeeper (e.g. "github", "slack").
	Name string `json:"name"`
	// Type is the backing service category (e.g. "vcs", "messaging", "cloud").
	Type string `json:"type"`
	// AuthType selects the credential scheme used against the backing service.
	AuthType AuthType `json:"auth_type"`
	// Credentials holds the secret material keyed by the auth scheme. It is
	// omitted from List() output by the registry to avoid leaking secrets.
	Credentials map[string]string `json:"-"`
	// AllowedTools is the allow-list of tool names a capability may request. An
	// empty slice means "any tool this gatekeeper fronts".
	AllowedTools []string `json:"allowed_tools,omitempty"`
	// Scopes is the superset of scopes a capability may be granted. An empty
	// slice means "no scope restriction".
	Scopes []string `json:"scopes,omitempty"`
	// RequireApproval forces human-in-the-loop approval before privileged
	// actions are recorded as approved (see LogAction / ApproveAction).
	RequireApproval bool `json:"require_approval"`
}

// Session is an issued, time-boxed grant. It is created by RequestAccess and
// referenced by ID from every action taken under it.
type Session struct {
	ID             string        `json:"id"`
	GatekeeperName string        `json:"gatekeeper_name"`
	Capabilities   []Capability  `json:"capabilities"`
	CreatedAt      time.Time     `json:"created_at"`
	ExpiresAt      *time.Time    `json:"expires_at,omitempty"`
	Status         SessionStatus `json:"status"`
}

// ActionLog is one row in the audit trail of a session.
type ActionLog struct {
	ID         string     `json:"id"`
	SessionID  string     `json:"session_id"`
	Action     string     `json:"action"`
	Resource   string     `json:"resource"`
	Result     string     `json:"result"`
	Approved   bool       `json:"approved"`
	Timestamp  time.Time  `json:"timestamp"`
	ApprovalID *string    `json:"approval_id,omitempty"`
}

// Gatekeeper moderates access to a single backing service. It validates
// requested capabilities against its config, mints sessions, and persists
// every action to an audit log. All methods are safe for concurrent use.
type Gatekeeper struct {
	config GatekeeperConfig
	db     *sql.DB
	mu     sync.Mutex
}

// NewGatekeeper binds a config to an audit-log database. The database must
// already have been migrated (see Migrate).
func NewGatekeeper(config GatekeeperConfig, db *sql.DB) (*Gatekeeper, error) {
	if config.Name == "" {
		return nil, errors.New("gatekeeper: config name is required")
	}
	if db == nil {
		return nil, errors.New("gatekeeper: db is required")
	}
	return &Gatekeeper{config: config, db: db}, nil
}

// Config returns the gatekeeper's configuration.
func (g *Gatekeeper) Config() GatekeeperConfig { return g.config }

// RequireApproval reports whether this gatekeeper gates privileged actions
// behind human approval.
func (g *Gatekeeper) RequireApproval() bool { return g.config.RequireApproval }

// RequestAccess validates the requested capability against this gatekeeper's
// allow-lists and, on success, mints a new active Session recorded in the
// audit log. The session inherits the capability's expiry, or
// defaultSessionTTL when none is set.
func (g *Gatekeeper) RequestAccess(cap Capability) (*Session, error) {
	if len(cap.Tool) == 0 {
		return nil, errors.New("gatekeeper: capability must name at least one tool")
	}
	if err := g.validateCapability(cap); err != nil {
		return nil, err
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	now := time.Now().UTC()
	expiresAt := cap.ExpiresAt
	if expiresAt == nil {
		e := now.Add(defaultSessionTTL)
		expiresAt = &e
	} else {
		v := expiresAt.UTC()
		expiresAt = &v
	}
	// Reject already-expired capabilities up front.
	if !expiresAt.After(now) {
		return nil, errors.New("gatekeeper: capability has already expired")
	}

	capsJSON, err := json.Marshal([]Capability{cap})
	if err != nil {
		return nil, fmt.Errorf("gatekeeper: marshal capabilities: %w", err)
	}

	sess := &Session{
		ID:             "gk-" + uuid.NewString(),
		GatekeeperName: g.config.Name,
		Capabilities:   []Capability{cap},
		CreatedAt:      now,
		ExpiresAt:      expiresAt,
		Status:         SessionActive,
	}

	_, err = g.db.Exec(
		`INSERT INTO gatekeeper_sessions (id, gatekeeper_name, capabilities, created_at, expires_at, status)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		sess.ID, sess.GatekeeperName, string(capsJSON),
		now.Unix(), expiresAt.Unix(), string(sess.Status),
	)
	if err != nil {
		return nil, fmt.Errorf("gatekeeper: persist session: %w", err)
	}
	return sess, nil
}

// validateCapability enforces that the requested tools and scopes fall within
// the gatekeeper's configured allow-lists.
func (g *Gatekeeper) validateCapability(cap Capability) error {
	if len(g.config.AllowedTools) > 0 {
		allowed := make(map[string]struct{}, len(g.config.AllowedTools))
		for _, t := range g.config.AllowedTools {
			allowed[t] = struct{}{}
		}
		for _, t := range cap.Tool {
			if _, ok := allowed[t]; !ok {
				return fmt.Errorf("gatekeeper: tool %q is not allowed for %q", t, g.config.Name)
			}
		}
	}
	if len(g.config.Scopes) > 0 {
		permitted := make(map[string]struct{}, len(g.config.Scopes))
		for _, s := range g.config.Scopes {
			permitted[s] = struct{}{}
		}
		for _, s := range cap.Scope {
			if _, ok := permitted[s]; !ok {
				return fmt.Errorf("gatekeeper: scope %q is not permitted for %q", s, g.config.Name)
			}
		}
	}
	return nil
}

// ValidateSession reports whether a session is currently usable: it must
// exist, be in the active state, and not have lapsed. Expired sessions are
// lazily transitioned to the "expired" status on first detection so callers
// never see a stale "active" row.
func (g *Gatekeeper) ValidateSession(sessionID string) bool {
	g.mu.Lock()
	defer g.mu.Unlock()

	var status string
	var expiresUnix sql.NullInt64
	err := g.db.QueryRow(
		`SELECT status, expires_at FROM gatekeeper_sessions WHERE id = ?`,
		sessionID,
	).Scan(&status, &expiresUnix)
	if err != nil {
		return false
	}
	if SessionStatus(status) != SessionActive {
		return false
	}
	if expiresUnix.Valid && expiresUnix.Int64 <= time.Now().Unix() {
		_, _ = g.db.Exec(
			`UPDATE gatekeeper_sessions SET status = ? WHERE id = ?`,
			string(SessionExpired), sessionID,
		)
		return false
	}
	return true
}

// LogAction appends an audit record for an action taken under a session. It
// returns an error if the session does not exist or is no longer active. When
// the gatekeeper is configured with RequireApproval, callers should log with
// approved=false first and flip it later via ApproveAction once a human
// approves.
func (g *Gatekeeper) LogAction(sessionID, action, resource, result string, approved bool) error {
	if action == "" {
		return errors.New("gatekeeper: action is required")
	}
	g.mu.Lock()
	defer g.mu.Unlock()

	if !g.sessionActiveLocked(sessionID) {
		return fmt.Errorf("gatekeeper: session %q is not active", sessionID)
	}

	id := "act-" + uuid.NewString()
	now := time.Now().UTC()
	_, err := g.db.Exec(
		`INSERT INTO gatekeeper_actions (id, session_id, action, resource, result, approved, timestamp)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		id, sessionID, action, resource, result, approved, now.Unix(),
	)
	if err != nil {
		return fmt.Errorf("gatekeeper: log action: %w", err)
	}
	return nil
}

// ApproveAction stamps an existing action with an approval identifier and
// flips its approved flag. It is a convenience for the RequireApproval flow
// and is a no-op error-wise if the action is already approved.
func (g *Gatekeeper) ApproveAction(actionID, approvalID string) error {
	if approvalID == "" {
		return errors.New("gatekeeper: approval id is required")
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	res, err := g.db.Exec(
		`UPDATE gatekeeper_actions SET approved = 1, approval_id = ? WHERE id = ?`,
		approvalID, actionID,
	)
	if err != nil {
		return fmt.Errorf("gatekeeper: approve action: %w", err)
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return fmt.Errorf("gatekeeper: action %q not found", actionID)
	}
	return nil
}

// GetActions returns the audit trail for a session, oldest first. Returns a
// nil slice (with no error) for a session that has taken no actions.
func (g *Gatekeeper) GetActions(sessionID string) ([]ActionLog, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	rows, err := g.db.Query(
		`SELECT id, session_id, action, resource, result, approved, timestamp, approval_id
		 FROM gatekeeper_actions
		 WHERE session_id = ?
		 ORDER BY timestamp ASC`,
		sessionID,
	)
	if err != nil {
		return nil, fmt.Errorf("gatekeeper: query actions: %w", err)
	}
	defer rows.Close()

	var out []ActionLog
	for rows.Next() {
		var a ActionLog
		var approved int
		var ts int64
		var approvalID sql.NullString
		if err := rows.Scan(&a.ID, &a.SessionID, &a.Action, &a.Resource,
			&a.Result, &approved, &ts, &approvalID); err != nil {
			return nil, fmt.Errorf("gatekeeper: scan action: %w", err)
		}
		a.Approved = approved != 0
		a.Timestamp = time.Unix(ts, 0).UTC()
		if approvalID.Valid {
			s := approvalID.String
			a.ApprovalID = &s
		}
		out = append(out, a)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("gatekeeper: iterate actions: %w", err)
	}
	return out, nil
}

// Revoke transitions a session to the revoked state, immediately invalidating
// it. Further actions under the session are rejected. Revoking an unknown
// session returns an error; revoking an already-revoked session is idempotent.
func (g *Gatekeeper) Revoke(sessionID string) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	res, err := g.db.Exec(
		`UPDATE gatekeeper_sessions SET status = ? WHERE id = ?`,
		string(SessionRevoked), sessionID,
	)
	if err != nil {
		return fmt.Errorf("gatekeeper: revoke session: %w", err)
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return fmt.Errorf("gatekeeper: session %q not found", sessionID)
	}
	return nil
}

// sessionActiveLocked is a non-locking helper used by methods that already
// hold g.mu. It lazily expires lapsed sessions and reports active status.
func (g *Gatekeeper) sessionActiveLocked(sessionID string) bool {
	var status string
	var expiresUnix sql.NullInt64
	err := g.db.QueryRow(
		`SELECT status, expires_at FROM gatekeeper_sessions WHERE id = ?`,
		sessionID,
	).Scan(&status, &expiresUnix)
	if err != nil {
		return false
	}
	if SessionStatus(status) != SessionActive {
		return false
	}
	if expiresUnix.Valid && expiresUnix.Int64 <= time.Now().Unix() {
		_, _ = g.db.Exec(
			`UPDATE gatekeeper_sessions SET status = ? WHERE id = ?`,
			string(SessionExpired), sessionID,
		)
		return false
	}
	return true
}

// GatekeeperRegistry is the top-level catalogue of configured gatekeepers. It
// owns the shared audit database and constructs Gatekeeper instances on
// demand. All methods are safe for concurrent use.
type GatekeeperRegistry struct {
	mu          sync.RWMutex
	db          *sql.DB
	gatekeepers map[string]*Gatekeeper
	configs     map[string]GatekeeperConfig
}

// NewGatekeeperRegistry binds a registry to an audit-log database. The database
// must already have been migrated (see Migrate).
func NewGatekeeperRegistry(db *sql.DB) (*GatekeeperRegistry, error) {
	if db == nil {
		return nil, errors.New("gatekeeper: db is required")
	}
	return &GatekeeperRegistry{
		db:          db,
		gatekeepers: make(map[string]*Gatekeeper),
		configs:     make(map[string]GatekeeperConfig),
	}, nil
}

// Register adds (or replaces) a gatekeeper by config name. The gatekeeper
// shares the registry's audit database.
func (r *GatekeeperRegistry) Register(config GatekeeperConfig) {
	gk := &Gatekeeper{config: config, db: r.db}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.gatekeepers[config.Name] = gk
	r.configs[config.Name] = config
}

// Get returns the gatekeeper registered under name, or an error if no
// gatekeeper is registered with that name.
func (r *GatekeeperRegistry) Get(name string) (*Gatekeeper, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	gk, ok := r.gatekeepers[name]
	if !ok {
		return nil, fmt.Errorf("gatekeeper: %q is not registered", name)
	}
	return gk, nil
}

// List returns the configs of all registered gatekeepers. Credentials are
// stripped to avoid leaking secrets.
func (r *GatekeeperRegistry) List() []GatekeeperConfig {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]GatekeeperConfig, 0, len(r.configs))
	for _, c := range r.configs {
		// Strip credentials before returning — never expose secrets in a
		// listing. c is a copy, so this does not mutate the stored config.
		c.Credentials = nil
		out = append(out, c)
	}
	return out
}

// Remove unregisters a gatekeeper. Existing sessions remain in the audit log
// for forensic completeness; they are not revoked automatically. Returns an
// error if no gatekeeper is registered with the given name.
func (r *GatekeeperRegistry) Remove(name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.gatekeepers[name]; !ok {
		return fmt.Errorf("gatekeeper: %q is not registered", name)
	}
	delete(r.gatekeepers, name)
	delete(r.configs, name)
	return nil
}

// Migrate creates the gatekeeper audit-log tables on the given database. It is
// idempotent and safe to call on every startup.
func Migrate(db *sql.DB) error {
	if db == nil {
		return errors.New("gatekeeper: db is required")
	}
	if _, err := db.Exec(schema); err != nil {
		return fmt.Errorf("gatekeeper: migrate: %w", err)
	}
	return nil
}

// schema creates the session and action audit tables. Timestamps are stored as
// Unix seconds (matching the rest of the project). Capabilities are serialized
// to JSON in a TEXT column since SQLite has no native array type.
const schema = `
CREATE TABLE IF NOT EXISTS gatekeeper_sessions (
	id              TEXT PRIMARY KEY,
	gatekeeper_name TEXT NOT NULL,
	capabilities    TEXT NOT NULL,
	created_at      INTEGER NOT NULL,
	expires_at      INTEGER,
	status          TEXT NOT NULL DEFAULT 'active'
);
CREATE INDEX IF NOT EXISTS idx_gatekeeper_sessions_gatekeeper
	ON gatekeeper_sessions(gatekeeper_name);
CREATE INDEX IF NOT EXISTS idx_gatekeeper_sessions_status
	ON gatekeeper_sessions(status);

CREATE TABLE IF NOT EXISTS gatekeeper_actions (
	id          TEXT PRIMARY KEY,
	session_id  TEXT NOT NULL,
	action      TEXT NOT NULL,
	resource    TEXT NOT NULL DEFAULT '',
	result      TEXT NOT NULL DEFAULT '',
	approved    INTEGER NOT NULL DEFAULT 0,
	timestamp   INTEGER NOT NULL,
	approval_id TEXT,
	FOREIGN KEY (session_id) REFERENCES gatekeeper_sessions(id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_gatekeeper_actions_session
	ON gatekeeper_actions(session_id);
CREATE INDEX IF NOT EXISTS idx_gatekeeper_actions_timestamp
	ON gatekeeper_actions(timestamp);
`
