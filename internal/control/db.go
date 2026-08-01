package control

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DB is the Postgres backing for run persistence and auth (users/sessions).
// Core runs fine without it: when TALON_DATABASE_URL is unset or unreachable,
// the store falls back to JSON persistence and auth is disabled.
type DB struct {
	pool *pgxpool.Pool
}

// ConnectDB opens a pool, verifies it, and applies the (idempotent) schema.
func ConnectDB(ctx context.Context, url string) (*DB, error) {
	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("pg pool: %w", err)
	}
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := pool.Ping(pingCtx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("pg ping: %w", err)
	}
	db := &DB{pool: pool}
	if err := db.migrate(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("pg migrate: %w", err)
	}
	return db, nil
}

func (db *DB) Close() { db.pool.Close() }

func (db *DB) migrate(ctx context.Context) error {
	ddl := `
CREATE EXTENSION IF NOT EXISTS citext;
CREATE TABLE IF NOT EXISTS users (
    id            uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    username      citext UNIQUE NOT NULL,
    password_hash text NOT NULL,
    created_at    timestamptz NOT NULL DEFAULT now()
);
CREATE TABLE IF NOT EXISTS sessions (
    token      text PRIMARY KEY,
    user_id    uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at timestamptz NOT NULL DEFAULT now(),
    expires_at timestamptz NOT NULL
);
CREATE TABLE IF NOT EXISTS runs (
    run_id     text PRIMARY KEY,
    data       jsonb NOT NULL,
    updated_at timestamptz NOT NULL DEFAULT now()
);
-- Queryable summary columns so /runs paginates in SQL instead of scanning jsonb.
ALTER TABLE runs ADD COLUMN IF NOT EXISTS target       text;
ALTER TABLE runs ADD COLUMN IF NOT EXISTS cve_id       text;
ALTER TABLE runs ADD COLUMN IF NOT EXISTS service_name text;
ALTER TABLE runs ADD COLUMN IF NOT EXISTS status       text;
ALTER TABLE runs ADD COLUMN IF NOT EXISTS judge_verdict boolean;
ALTER TABLE runs ADD COLUMN IF NOT EXISTS tool_calls   integer;
ALTER TABLE runs ADD COLUMN IF NOT EXISTS started_at   timestamptz;
CREATE INDEX IF NOT EXISTS runs_started_at_idx ON runs (started_at DESC);
CREATE INDEX IF NOT EXISTS runs_status_idx ON runs (status);
CREATE TABLE IF NOT EXISTS settings (
    key        text PRIMARY KEY,
    value      jsonb NOT NULL,
    updated_at timestamptz NOT NULL DEFAULT now()
);`
	if _, err := db.pool.Exec(ctx, ddl); err != nil {
		return err
	}
	// Backfill summary columns for rows written before they existed.
	_, err := db.pool.Exec(ctx, `
UPDATE runs SET
    target        = data->'RunInput'->>'TargetIP',
    cve_id        = NULLIF(data->'RunInput'->>'CVEID', ''),
    service_name  = NULLIF(data->'RunInput'->>'ServiceName', ''),
    status        = data->>'Status',
    judge_verdict = CASE WHEN data->>'JudgeSet' = 'true' THEN (data->>'JudgeVerdict')::boolean ELSE NULL END,
    tool_calls    = COALESCE(jsonb_array_length(data->'ToolLog'), 0),
    started_at    = COALESCE((data->>'StartedAt')::timestamptz, updated_at)
WHERE target IS NULL`)
	return err
}

// ---- run persistence ----

func (db *DB) upsertRun(ctx context.Context, runID string, sess *Session) error {
	data, err := json.Marshal(sess)
	if err != nil {
		return err
	}
	var verdict *bool
	if sess.JudgeSet {
		v := sess.JudgeVerdict
		verdict = &v
	}
	started := sess.StartedAt
	if started.IsZero() {
		started = time.Now().UTC()
	}
	_, err = db.pool.Exec(ctx,
		`INSERT INTO runs (run_id, data, updated_at, target, cve_id, service_name, status, judge_verdict, tool_calls, started_at)
		 VALUES ($1, $2, now(), $3, $4, $5, $6, $7, $8, $9)
		 ON CONFLICT (run_id) DO UPDATE SET
		   data = $2, updated_at = now(), target = $3, cve_id = $4, service_name = $5,
		   status = $6, judge_verdict = $7, tool_calls = $8, started_at = $9`,
		runID, data,
		sess.RunInput.TargetIP, nilIfEmpty(sess.RunInput.CVEID), nilIfEmpty(sess.RunInput.ServiceName),
		sess.Status, verdict, len(sess.ToolLog), started)
	return err
}

func nilIfEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// getRun loads one full session row; ok=false when unknown.
func (db *DB) getRun(ctx context.Context, runID string) (*Session, bool) {
	var raw []byte
	err := db.pool.QueryRow(ctx, `SELECT data FROM runs WHERE run_id = $1`, runID).Scan(&raw)
	if err != nil {
		return nil, false
	}
	sess := &Session{}
	if err := json.Unmarshal(raw, sess); err != nil {
		return nil, false
	}
	return sess, true
}

// listRunsPaginated returns one page of summaries (newest first) plus the
// total row count — scales to any table size.
func (db *DB) listRunsPaginated(ctx context.Context, limit, offset int) ([]RunSummary, int, error) {
	var total int
	if err := db.pool.QueryRow(ctx, `SELECT count(1) FROM runs`).Scan(&total); err != nil {
		return nil, 0, err
	}
	rows, err := db.pool.Query(ctx,
		`SELECT run_id, COALESCE(target,''), COALESCE(cve_id,''), COALESCE(service_name,''),
		        COALESCE(status,''), judge_verdict, COALESCE(tool_calls,0),
		        COALESCE(started_at, updated_at), updated_at
		 FROM runs ORDER BY started_at DESC NULLS LAST, updated_at DESC LIMIT $1 OFFSET $2`,
		limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	out := []RunSummary{}
	for rows.Next() {
		var s RunSummary
		var updatedAt time.Time
		if err := rows.Scan(&s.RunID, &s.Target, &s.CVEID, &s.ServiceName, &s.Status, &s.JudgeVerdict, &s.ToolCalls, &s.StartedAt, &updatedAt); err != nil {
			return nil, 0, err
		}
		// For terminal runs the last write is the completion write, so updated_at
		// is the end time — lets the UI freeze the elapsed timer.
		if s.Status == "completed" || s.Status == "error" {
			s.EndedAt = updatedAt
		}
		out = append(out, s)
	}
	return out, total, rows.Err()
}

// RunsSummary is the aggregate served by GET /runs/summary (overview stats).
type RunsSummary struct {
	Total            int `json:"total"`
	Active           int `json:"active"`
	Compromised      int `json:"compromised"`
	AwaitingApproval int `json:"awaiting_approval"`
	Completed        int `json:"completed"`
	Errored          int `json:"errored"`
}

func (db *DB) runsSummary(ctx context.Context) (RunsSummary, error) {
	var s RunsSummary
	err := db.pool.QueryRow(ctx, `
SELECT count(1),
       count(1) FILTER (WHERE status IN ('initializing','running','awaiting_approval')),
       count(1) FILTER (WHERE judge_verdict IS TRUE),
       count(1) FILTER (WHERE status = 'awaiting_approval'),
       count(1) FILTER (WHERE status = 'completed'),
       count(1) FILTER (WHERE status = 'error')
FROM runs`).Scan(&s.Total, &s.Active, &s.Compromised, &s.AwaitingApproval, &s.Completed, &s.Errored)
	return s, err
}

// markInterruptedRunsError flips non-terminal rows to error at boot
// (orchestrator state is in-memory; those runs can never resume).
func (db *DB) markInterruptedRunsError(ctx context.Context) error {
	rows, err := db.pool.Query(ctx,
		`SELECT run_id, data FROM runs WHERE status IN ('initializing','running','awaiting_approval')`)
	if err != nil {
		return err
	}
	defer rows.Close()
	type staleRun struct {
		id   string
		sess *Session
	}
	var stale []staleRun
	for rows.Next() {
		var id string
		var raw []byte
		if err := rows.Scan(&id, &raw); err != nil {
			return err
		}
		sess := &Session{}
		if err := json.Unmarshal(raw, sess); err != nil {
			continue
		}
		stale = append(stale, staleRun{id, sess})
	}
	if err := rows.Err(); err != nil {
		return err
	}
	for _, sr := range stale {
		old := sr.sess.Status
		sr.sess.Status = "error"
		sr.sess.Output = "run interrupted by server shutdown; cannot resume"
		sr.sess.PendingInterrupt = nil
		sr.sess.History = append(sr.sess.History, historyLine("marked error on reload (was %s)", old))
		if err := db.upsertRun(ctx, sr.id, sr.sess); err != nil {
			log.Printf("control: pg mark stale run %s: %v", sr.id, err)
		}
	}
	if len(stale) > 0 {
		log.Printf("control: marked %d interrupted runs as error", len(stale))
	}
	return nil
}

func (db *DB) runsEmpty(ctx context.Context) (bool, error) {
	var n int
	err := db.pool.QueryRow(ctx, `SELECT count(1) FROM runs`).Scan(&n)
	return n == 0, err
}

// ---- auth ----

func (db *DB) seedUser(ctx context.Context, username, passwordHash string) error {
	_, err := db.pool.Exec(ctx,
		`INSERT INTO users (username, password_hash) VALUES ($1, $2)
		 ON CONFLICT (username) DO NOTHING`, username, passwordHash)
	return err
}

func (db *DB) passwordHashFor(ctx context.Context, username string) (userID string, hash string, err error) {
	err = db.pool.QueryRow(ctx,
		`SELECT id, password_hash FROM users WHERE username = $1`, username).Scan(&userID, &hash)
	return
}

func (db *DB) createSession(ctx context.Context, userID, token string, ttl time.Duration) error {
	_, err := db.pool.Exec(ctx,
		`INSERT INTO sessions (token, user_id, expires_at) VALUES ($1, $2, $3)`,
		token, userID, time.Now().Add(ttl))
	return err
}

// sessionUser resolves a token to its username, or "" when unknown/expired.
func (db *DB) sessionUser(ctx context.Context, token string) string {
	var username string
	err := db.pool.QueryRow(ctx,
		`SELECT u.username FROM sessions s JOIN users u ON u.id = s.user_id
		 WHERE s.token = $1 AND s.expires_at > now()`, token).Scan(&username)
	if err != nil {
		return ""
	}
	return username
}

func (db *DB) deleteSession(ctx context.Context, token string) {
	_, _ = db.pool.Exec(ctx, `DELETE FROM sessions WHERE token = $1`, token)
}

// ---- runtime settings (dashboard-managed config) ----

func (db *DB) allSettings(ctx context.Context) (map[string]json.RawMessage, error) {
	rows, err := db.pool.Query(ctx, `SELECT key, value FROM settings`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := map[string]json.RawMessage{}
	for rows.Next() {
		var k string
		var v json.RawMessage
		if err := rows.Scan(&k, &v); err != nil {
			return nil, err
		}
		out[k] = v
	}
	return out, rows.Err()
}

func (db *DB) setSetting(ctx context.Context, key string, value json.RawMessage) error {
	_, err := db.pool.Exec(ctx,
		`INSERT INTO settings (key, value, updated_at) VALUES ($1, $2, now())
		 ON CONFLICT (key) DO UPDATE SET value = $2, updated_at = now()`, key, value)
	return err
}
