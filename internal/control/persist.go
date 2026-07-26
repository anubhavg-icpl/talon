package control

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

// persistedFile is the on-disk shape of the run table (runs.json).
type persistedFile struct {
	Sessions map[string]*Session `json:"sessions"`
}

// EnablePersistence loads runs from <dir>/runs.json (if present) and flushes
// every subsequent store mutation back to that file. Runs that were left in
// a non-terminal state by a shutdown cannot be resumed (orchestrator state
// is in-memory), so they are marked as errors on load.
func (s *Store) EnablePersistence(dir string) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	path := filepath.Join(dir, "runs.json")

	s.mu.Lock()
	defer s.mu.Unlock()
	s.persistPath = path

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	var pf persistedFile
	if err := json.Unmarshal(data, &pf); err != nil {
		log.Printf("control: ignoring unreadable %s: %v", path, err)
		return nil
	}
	for id, sess := range pf.Sessions {
		if sess == nil {
			continue
		}
		switch sess.Status {
		case "initializing", "running", "awaiting_approval":
			old := sess.Status
			sess.Status = "error"
			sess.Output = "run interrupted by server shutdown; cannot resume"
			sess.PendingInterrupt = nil
			sess.History = append(sess.History, historyLine("marked error on reload (was %s)", old))
		}
		s.sessions[id] = sess
	}
	log.Printf("control: loaded %d persisted runs from %s", len(pf.Sessions), path)
	return nil
}

// EnablePostgres switches persistence to Postgres: imports runs.json once if
// the table is empty, marks interrupted runs as error, and flushes every
// mutation via upserts. Runs are NOT loaded into memory — the in-memory
// table holds active runs only; Get/PaginatedList/Summary fall through to
// SQL, so the registry scales to any size.
func (s *Store) EnablePostgres(ctx context.Context, db *DB, jsonDir string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.pg = db
	s.pgCtx = context.Background()

	// One-time import: legacy JSON store → empty PG table.
	if empty, err := db.runsEmpty(ctx); err == nil && empty && jsonDir != "" {
		if data, err := os.ReadFile(filepath.Join(jsonDir, "runs.json")); err == nil {
			var pf persistedFile
			if err := json.Unmarshal(data, &pf); err == nil {
				n := 0
				for id, sess := range pf.Sessions {
					if sess != nil {
						if err := db.upsertRun(ctx, id, sess); err == nil {
							n++
						}
					}
				}
				log.Printf("control: imported %d runs from runs.json into postgres", n)
			}
		}
	}

	if err := db.markInterruptedRunsError(ctx); err != nil {
		log.Printf("control: pg stale-run sweep: %v", err)
	}
	return nil
}

// saveLocked flushes one run's state. Callers must hold s.mu.
// Postgres upsert when enabled, else full JSON file flush, else no-op.
// Errors are logged, never fatal: the in-memory store remains authoritative.
func (s *Store) saveLocked(runID string) {
	if s.pg != nil {
		if sess, ok := s.sessions[runID]; ok {
			if err := s.pg.upsertRun(s.pgCtx, runID, sess); err != nil {
				log.Printf("control: pg upsert run %s: %v", runID, err)
			}
		}
		return
	}
	if s.persistPath == "" {
		return
	}
	data, err := json.Marshal(persistedFile{Sessions: s.sessions})
	if err != nil {
		log.Printf("control: marshal runs: %v", err)
		return
	}
	tmp := s.persistPath + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		log.Printf("control: write %s: %v", tmp, err)
		return
	}
	if err := os.Rename(tmp, s.persistPath); err != nil {
		log.Printf("control: rename %s: %v", tmp, err)
	}
}
