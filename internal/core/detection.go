package core

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

// DetectionVerdict is the outcome of a triage or investigation skill.
type DetectionVerdict string

const (
	VerdictEscalate     DetectionVerdict = "escalate"
	VerdictDismiss      DetectionVerdict = "dismiss"
	VerdictMalicious    DetectionVerdict = "malicious"
	VerdictSuspicious   DetectionVerdict = "suspicious"
	VerdictInconclusive DetectionVerdict = "inconclusive"
	VerdictBenign       DetectionVerdict = "benign"
)

// DetectionCase represents one alert/detection event being processed.
type DetectionCase struct {
	ID          string                 `json:"id"`
	AlertType   string                 `json:"alert_type"`
	Title       string                 `json:"title"`
	Entity      string                 `json:"entity"`
	EntityType  string                 `json:"entity_type"` // user, host, ip, etc.
	Severity    string                 `json:"severity"`
	SourceData  map[string]any         `json:"source_data"`
	TriageState *TriageState           `json:"triage_state,omitempty"`
	InvestigationState *InvestigationState `json:"investigation_state,omitempty"`
	TuningState *TuningState           `json:"tuning_state,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// TriageState holds the result of running a triage skill on a case.
type TriageState struct {
	Verdict   DetectionVerdict `json:"verdict"`
	RiskCount int              `json:"risk_count"`
	Checks    []TriageCheck    `json:"checks"`
	Evidence  string           `json:"evidence"`
	Reasoning string           `json:"reasoning"`
	SkillID   string           `json:"skill_id"`
	Timestamp time.Time        `json:"timestamp"`
}

// TriageCheck is one individual check within a triage skill.
type TriageCheck struct {
	Name     string           `json:"name"`
	Outcome  DetectionVerdict `json:"outcome"` // "RISK" or "CLEAR"
	Label    string           `json:"label"`   // outcome label e.g. "vip-target", "burst"
	Detail   string           `json:"detail"`
}

// InvestigationState holds the result of running an investigation skill.
type InvestigationState struct {
	Verdict       DetectionVerdict   `json:"verdict"`
	Signals       []InvestigationSignal `json:"signals"`
	Evidence      string             `json:"evidence"`
	Reasoning     string             `json:"reasoning"`
	Actions       []string           `json:"recommended_actions"`
	SkillID       string             `json:"skill_id"`
	Timestamp     time.Time          `json:"timestamp"`
}

// InvestigationSignal is one signal (A/B/C/D pattern) from an investigation.
type InvestigationSignal struct {
	Name    string `json:"name"`     // e.g. "Beacon confirmation", "Destination attribution"
	Fired   bool   `json:"fired"`
	Detail  string `json:"detail"`
}

// TuningState holds a proposed detection tuning change.
type TuningState struct {
	Action    string   `json:"action"`    // exclude, include, modify, fork, none
	Target    string   `json:"target"`    // what to change
	Value     string   `json:"value"`     // proposed value
	Rationale string   `json:"rationale"`
	SkillID   string   `json:"skill_id"`
	Timestamp time.Time `json:"timestamp"`
}

// CaseStore manages detection cases in-memory (run-scoped).
type CaseStore struct {
	mu    sync.Mutex
	cases map[string]*DetectionCase
}

func NewCaseStore() *CaseStore {
	return &CaseStore{cases: make(map[string]*DetectionCase)}
}

// CreateCase registers a new detection case.
func (cs *CaseStore) CreateCase(c DetectionCase) *DetectionCase {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	if c.ID == "" {
		c.ID = fmt.Sprintf("case-%d", len(cs.cases)+1)
	}
	if c.CreatedAt.IsZero() {
		c.CreatedAt = time.Now().UTC()
	}
	c.UpdatedAt = time.Now().UTC()
	cs.cases[c.ID] = &c
	return &c
}

// GetCase retrieves a case by ID.
func (cs *CaseStore) GetCase(id string) (*DetectionCase, bool) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	c, ok := cs.cases[id]
	return c, ok
}

// ListCases returns all cases (optionally filtered by alert type).
func (cs *CaseStore) ListCases(alertType string) []*DetectionCase {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	var out []*DetectionCase
	for _, c := range cs.cases {
		if alertType == "" || c.AlertType == alertType {
			out = append(out, c)
		}
	}
	return out
}

// UpdateTriage sets the triage result on a case.
func (cs *CaseStore) UpdateTriage(id string, state TriageState) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	c, ok := cs.cases[id]
	if !ok {
		return fmt.Errorf("case %s not found", id)
	}
	c.TriageState = &state
	c.UpdatedAt = time.Now().UTC()
	return nil
}

// UpdateInvestigation sets the investigation result on a case.
func (cs *CaseStore) UpdateInvestigation(id string, state InvestigationState) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	c, ok := cs.cases[id]
	if !ok {
		return fmt.Errorf("case %s not found", id)
	}
	c.InvestigationState = &state
	c.UpdatedAt = time.Now().UTC()
	return nil
}

// UpdateTuning sets the tuning result on a case.
func (cs *CaseStore) UpdateTuning(id string, state TuningState) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	c, ok := cs.cases[id]
	if !ok {
		return fmt.Errorf("case %s not found", id)
	}
	c.TuningState = &state
	c.UpdatedAt = time.Now().UTC()
	return nil
}

// ApplyTriageDecision computes the verdict from triage checks using majority rule.
func ApplyTriageDecision(checks []TriageCheck) DetectionVerdict {
	riskCount := 0
	for _, c := range checks {
		if c.Outcome == "RISK" || c.Outcome == VerdictEscalate {
			riskCount++
		}
	}
	// Majority rule: 2/3 or more → escalate
	if riskCount*2 >= len(checks) && len(checks) > 0 {
		return VerdictEscalate
	}
	return VerdictDismiss
}

// ApplyInvestigationVerdict derives the verdict from signals.
// Signal pattern: B or C plus A or D → malicious; D alone with A → malicious.
func ApplyInvestigationVerdict(signals []InvestigationSignal) DetectionVerdict {
	fired := make(map[string]bool)
	for _, s := range signals {
		fired[s.Name] = s.Fired
	}

	// Simplified verdict logic from the C2 beacon pattern:
	// malicious: B or C plus A or D, or D alone with A
	if (fired["B"] || fired["C"]) && (fired["A"] || fired["D"]) {
		return VerdictMalicious
	}
	if fired["D"] && fired["A"] {
		return VerdictMalicious
	}
	// suspicious: confirmed beacon (A) but no clear attribution
	if fired["A"] && !fired["B"] && !fired["C"] && !fired["D"] {
		return VerdictSuspicious
	}
	// benign: all signals cleared and A resolved to benign
	if !fired["A"] && !fired["B"] && !fired["C"] && !fired["D"] {
		return VerdictBenign
	}
	return VerdictInconclusive
}

// PipelineStage represents the detection triage→investigation→tuning pipeline.
type PipelineStage string

const (
	StageTriage        PipelineStage = "triage"
	StageInvestigation PipelineStage = "investigation"
	StageTuning        PipelineStage = "tuning"
)

// PipelineResult is the output of running the detection pipeline on one case.
type PipelineResult struct {
	CaseID    string            `json:"case_id"`
	Stage     PipelineStage     `json:"stage"`
	Verdict   DetectionVerdict  `json:"verdict"`
	Summary   string            `json:"summary"`
	NextStage PipelineStage     `json:"next_stage,omitempty"`
	Triage    *TriageState      `json:"triage,omitempty"`
	Investigation *InvestigationState `json:"investigation,omitempty"`
	Tuning    *TuningState      `json:"tuning,omitempty"`
}

// DetermineNextStage decides which pipeline stage runs next based on the verdict.
func DetermineNextStage(stage PipelineStage, verdict DetectionVerdict) PipelineStage {
	switch stage {
	case StageTriage:
		if verdict == VerdictEscalate {
			return StageInvestigation
		}
		return StageTuning // dismissed cases get tuning review
	case StageInvestigation:
		if verdict == VerdictMalicious || verdict == VerdictSuspicious {
			return StageTuning
		}
		return "" // benign/inconclusive → pipeline ends
	case StageTuning:
		return "" // tuning is the last stage
	}
	return ""
}

// FormatVerdictSummary produces a human-readable one-line summary.
func FormatVerdictSummary(r PipelineResult) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Case %s [%s]: %s", r.CaseID, r.Stage, r.Verdict))
	if r.Triage != nil && len(r.Triage.Checks) > 0 {
		sb.WriteString(fmt.Sprintf(" (%d checks, %d risk)", len(r.Triage.Checks), r.Triage.RiskCount))
	}
	if r.Investigation != nil && len(r.Investigation.Signals) > 0 {
		fired := 0
		for _, s := range r.Investigation.Signals {
			if s.Fired {
				fired++
			}
		}
		sb.WriteString(fmt.Sprintf(" (%d/%d signals fired)", fired, len(r.Investigation.Signals)))
	}
	if r.Tuning != nil && r.Tuning.Action != "none" {
		sb.WriteString(fmt.Sprintf(" → tuning: %s %s", r.Tuning.Action, r.Tuning.Target))
	}
	if r.NextStage != "" {
		sb.WriteString(fmt.Sprintf(" → next: %s", r.NextStage))
	}
	return sb.String()
}
