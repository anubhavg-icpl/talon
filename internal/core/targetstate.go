package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// TargetState stores per-target pentest state for resumption and continuity.
// Ported from pentest agent target_state/store.py + planner.py.
type TargetState struct {
	mu sync.Mutex

	Target    string             `json:"target"`
	Slug      string             `json:"slug"`
	Findings  []TargetFinding    `json:"findings"`
	ReconDims []ReconDimension   `json:"recon_dimensions"`
	FailedVectors []FailedVector `json:"failed_vectors"`
	AttackPath   []AttackStep    `json:"attack_path"`
	Runtime      TargetRuntime   `json:"runtime"`
	Snapshots    []Snapshot      `json:"snapshots"`
	SchemaVersion int            `json:"schema_version"`
	UpdatedAt    time.Time       `json:"updated_at"`
	CreatedAt    time.Time       `json:"created_at"`
}

type TargetFinding struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Severity    string    `json:"severity"`
	Status      string    `json:"status"` // "verified", "unverified", "false_positive"
	EvidenceIDs []string  `json:"evidence_ids,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

type ReconDimension struct {
	Name      string    `json:"name"`     // "port_scan", "dns", "web_enum", etc.
	Status    string    `json:"status"`   // "complete", "partial", "not_started"
	Summary   string    `json:"summary"`
	UpdatedAt time.Time `json:"updated_at"`
}

type FailedVector struct {
	Vector    string    `json:"vector"`   // "sql_injection", "xss", etc.
	Target    string    `json:"target"`   // endpoint/URL
	Reason    string    `json:"reason"`   // why it failed
	TriedAt   time.Time `json:"tried_at"`
}

type AttackStep struct {
	Step      int       `json:"step"`
	Action    string    `json:"action"`
	Result    string    `json:"result"`
	Timestamp time.Time `json:"timestamp"`
}

type TargetRuntime struct {
	OS         string `json:"os,omitempty"`
	Services   string `json:"services,omitempty"`
	Credentials string `json:"credentials,omitempty"`
	Notes      string `json:"notes,omitempty"`
}

type Snapshot struct {
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Label     string    `json:"label"`
	Data      string    `json:"data"` // JSON snapshot of TargetState
}

const targetSchemaVersion = 1

// TargetStore manages per-target state files on disk.
type TargetStore struct {
	dataDir string
	mu      sync.Mutex
}

// NewTargetStore creates a new target state store.
// dataDir is the base directory (default: "talon-data/targets").
func NewTargetStore(dataDir string) *TargetStore {
	if dataDir == "" {
		dataDir = "talon-data/targets"
	}
	return &TargetStore{dataDir: dataDir}
}

// targetSlug converts a target string to a filesystem-safe slug.
func targetSlug(target string) string {
	slug := strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' || r == '.' {
			return r
		}
		return '-'
	}, strings.ToLower(target))
	// Collapse multiple dashes
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}
	return strings.Trim(slug, "-")
}

// statePath returns the filesystem path for a target's state file.
func (ts *TargetStore) statePath(slug string) string {
	return filepath.Join(ts.dataDir, slug+".json")
}

// GetOrCreate retrieves existing state for a target, or creates new state.
func (ts *TargetStore) GetOrCreate(target string) (*TargetState, error) {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	slug := targetSlug(target)
	path := ts.statePath(slug)

	data, err := os.ReadFile(path)
	if err == nil {
		var state TargetState
		if err := json.Unmarshal(data, &state); err != nil {
			return nil, fmt.Errorf("corrupt state file %s: %v", path, err)
		}
		return &state, nil
	}

	// Create new
	state := &TargetState{
		Target:        target,
		Slug:          slug,
		SchemaVersion: targetSchemaVersion,
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
	}
	return state, nil
}

// Save persists target state atomically (write tmp + rename).
func (ts *TargetStore) Save(state *TargetState) error {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	state.UpdatedAt = time.Now().UTC()

	if err := os.MkdirAll(ts.dataDir, 0755); err != nil {
		return fmt.Errorf("create targets dir: %v", err)
	}

	path := ts.statePath(state.Slug)
	data, err := json.MarshalIndent(state, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal state: %v", err)
	}

	// Atomic write: tmp file + rename
	tmpPath := path + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0644); err != nil {
		return fmt.Errorf("write state: %v", err)
	}
	if err := os.Rename(tmpPath, path); err != nil {
		return fmt.Errorf("rename state: %v", err)
	}

	return nil
}

// Snapshot creates a named snapshot of the current state.
func (ts *TargetStore) Snapshot(state *TargetState, label string) error {
	state.mu.Lock()
	defer state.mu.Unlock()

	data, _ := json.Marshal(state)
	snap := Snapshot{
		ID:        fmt.Sprintf("snap-%d", len(state.Snapshots)+1),
		Timestamp: time.Now().UTC(),
		Label:     label,
		Data:      string(data),
	}
	state.Snapshots = append(state.Snapshots, snap)
	return ts.Save(state)
}

// AddFinding adds a verified finding to target state.
func (s *TargetState) AddFinding(f TargetFinding) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if f.CreatedAt.IsZero() {
		f.CreatedAt = time.Now().UTC()
	}
	s.Findings = append(s.Findings, f)
}

// AddReconDim updates or adds a recon dimension.
func (s *TargetState) AddReconDim(dim ReconDimension) {
	s.mu.Lock()
	defer s.mu.Unlock()
	dim.UpdatedAt = time.Now().UTC()
	for i, existing := range s.ReconDims {
		if existing.Name == dim.Name {
			s.ReconDims[i] = dim
			return
		}
	}
	s.ReconDims = append(s.ReconDims, dim)
}

// AddFailedVector records a failed attack vector.
func (s *TargetState) AddFailedVector(fv FailedVector) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if fv.TriedAt.IsZero() {
		fv.TriedAt = time.Now().UTC()
	}
	s.FailedVectors = append(s.FailedVectors, fv)
}

// AddAttackStep appends a step to the attack path.
func (s *TargetState) AddAttackStep(action, result string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.AttackPath = append(s.AttackPath, AttackStep{
		Step:      len(s.AttackPath) + 1,
		Action:    action,
		Result:    result,
		Timestamp: time.Now().UTC(),
	})
}

// ResumePlan builds a deterministic next-step plan from stored state.
// No LLM needed — pure heuristic from prior findings, recon dims, and failures.
func (s *TargetState) ResumePlan() []PlanStep {
	s.mu.Lock()
	defer s.mu.Unlock()

	var steps []PlanStep

	// Check incomplete recon dimensions
	reconComplete := make(map[string]bool)
	for _, dim := range s.ReconDims {
		if dim.Status == "complete" {
			reconComplete[dim.Name] = true
		}
	}

	// Required recon dimensions in order
	requiredRecon := []struct {
		Name string
		Desc string
	}{
		{"port_scan", "Complete port scan to identify open services"},
		{"dns", "DNS enumeration and subdomain discovery"},
		{"web_enum", "Web directory and file enumeration"},
		{"service_id", "Service version identification"},
	}

	for _, req := range requiredRecon {
		if !reconComplete[req.Name] {
			steps = append(steps, PlanStep{
				Priority:   1,
				Category:   "recon",
				Action:     req.Desc,
				Dimension:  req.Name,
				Confidence: 0.9,
			})
		}
	}

	// Check failed vectors for retry or alternative
	failedVectors := make(map[string][]string)
	for _, fv := range s.FailedVectors {
		failedVectors[fv.Vector] = append(failedVectors[fv.Vector], fv.Target)
	}

	for vector, targets := range failedVectors {
		if len(targets) >= 2 {
			steps = append(steps, PlanStep{
				Priority:   2,
				Category:   "exploit_retry",
				Action:     fmt.Sprintf("Retry %s on %d targets with alternative payloads", vector, len(targets)),
				Confidence: 0.5,
			})
		}
	}

	// If we have verified findings, suggest post-exploitation
	verifiedCount := 0
	for _, f := range s.Findings {
		if f.Status == "verified" {
			verifiedCount++
		}
	}
	if verifiedCount > 0 {
		steps = append(steps, PlanStep{
			Priority:   3,
			Category:   "post_exploit",
			Action:     fmt.Sprintf("Pivot from %d verified findings — explore post-exploitation", verifiedCount),
			Confidence: 0.7,
		})
	}

	// If no steps, suggest initial recon
	if len(steps) == 0 {
		steps = append(steps, PlanStep{
			Priority:   1,
			Category:   "recon",
			Action:     "Start fresh reconnaissance — no prior state found",
			Confidence: 0.5,
		})
	}

	return steps
}

// PlanStep represents one step in a resume plan.
type PlanStep struct {
	Priority   int     `json:"priority"`
	Category   string  `json:"category"`
	Action     string  `json:"action"`
	Dimension  string  `json:"dimension,omitempty"`
	Confidence float64 `json:"confidence"`
}

// Summary returns a brief text summary of the target state for prompt injection.
func (s *TargetState) Summary() string {
	s.mu.Lock()
	defer s.mu.Unlock()

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Target: %s\n", s.Target))
	sb.WriteString(fmt.Sprintf("Findings: %d total", len(s.Findings)))
	verified := 0
	for _, f := range s.Findings {
		if f.Status == "verified" {
			verified++
		}
	}
	if verified > 0 {
		sb.WriteString(fmt.Sprintf(" (%d verified)", verified))
	}
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("Recon dimensions: %d (", len(s.ReconDims)))
	for i, dim := range s.ReconDims {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("%s=%s", dim.Name, dim.Status))
	}
	sb.WriteString(")\n")
	sb.WriteString(fmt.Sprintf("Failed vectors: %d\n", len(s.FailedVectors)))
	sb.WriteString(fmt.Sprintf("Attack path steps: %d\n", len(s.AttackPath)))
	return sb.String()
}
