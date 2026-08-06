package core

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/anubhavg-icpl/talon/internal/llm"
)

// detectionToolSpecs returns agent-facing detection pipeline tools.
func detectionToolSpecs() []llm.ToolSpec {
	return []llm.ToolSpec{
		{
			Name:        "detection_create_case",
			Description: "Create a detection case from a security alert. " +
				"Registers the alert with its type, entity, severity, and source data for the triage→investigation→tuning pipeline.",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"alert_type":  map[string]any{"type": "string", "description": "Alert category (e.g. mfa-fatigue, c2-beacon, impossible-travel)"},
					"title":       map[string]any{"type": "string"},
					"entity":      map[string]any{"type": "string", "description": "Affected entity (user, host, IP)"},
					"entity_type": map[string]any{"type": "string", "description": "user, host, ip, etc."},
					"severity":    map[string]any{"type": "string", "enum": []any{"critical", "high", "medium", "low", "info"}},
					"source_data": map[string]any{"type": "object", "additionalProperties": true},
				},
				"required": []any{"alert_type", "title", "entity"},
			},
		},
		{
			Name:        "detection_triage",
			Description: "Run a triage skill on a detection case. " +
				"Searches the detection skill catalog for the matching triage skill, " +
				"applies its checks, and returns an escalate/dismiss verdict with evidence.",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"case_id":   map[string]any{"type": "string"},
					"skill_id":  map[string]any{"type": "string", "description": "Optional: specific triage skill to use (from skill_search)"},
					"checks": map[string]any{
						"type": "array",
						"description": "Triage check results (from running the skill methodology)",
						"items": map[string]any{
							"type": "object",
							"properties": map[string]any{
								"name":    map[string]any{"type": "string"},
								"outcome": map[string]any{"type": "string", "enum": []any{"RISK", "CLEAR"}},
								"label":   map[string]any{"type": "string"},
								"detail":  map[string]any{"type": "string"},
							},
						},
					},
					"evidence":  map[string]any{"type": "string"},
					"reasoning": map[string]any{"type": "string"},
				},
				"required": []any{"case_id", "checks"},
			},
		},
		{
			Name:        "detection_investigate",
			Description: "Run an investigation skill on an escalated case. " +
				"Collects investigation signals, derives a verdict (malicious/suspicious/inconclusive/benign), " +
				"and recommends containment actions.",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"case_id":  map[string]any{"type": "string"},
					"skill_id": map[string]any{"type": "string", "description": "Optional: specific investigation skill"},
					"signals": map[string]any{
						"type": "array",
						"description": "Investigation signal results",
						"items": map[string]any{
							"type": "object",
							"properties": map[string]any{
								"name":   map[string]any{"type": "string"},
								"fired":  map[string]any{"type": "boolean"},
								"detail": map[string]any{"type": "string"},
							},
						},
					},
					"evidence":  map[string]any{"type": "string"},
					"reasoning": map[string]any{"type": "string"},
					"actions":   map[string]any{"type": "array", "items": map[string]any{"type": "string"}},
				},
				"required": []any{"case_id", "signals"},
			},
		},
		{
			Name:        "detection_tune",
			Description: "Propose a detection tuning change after case closure. " +
				"Returns one action (exclude, include, modify, fork, none) for human review.",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"case_id":   map[string]any{"type": "string"},
					"skill_id":  map[string]any{"type": "string"},
					"action":    map[string]any{"type": "string", "enum": []any{"exclude", "include", "modify", "fork", "none"}},
					"target":    map[string]any{"type": "string"},
					"value":     map[string]any{"type": "string"},
					"rationale": map[string]any{"type": "string"},
				},
				"required": []any{"case_id", "action"},
			},
		},
		{
			Name:        "detection_list_cases",
			Description: "List detection cases, optionally filtered by alert type.",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"alert_type": map[string]any{"type": "string"},
				},
			},
		},
	}
}

// detectionCaseStoreKey is used to store the CaseStore in the tracker context.
// Since we can't extend tracker, we use a package-level store per run.
var (
	detectionStores   = make(map[string]*CaseStore)
	detectionStoresMu sync.Mutex
)

func getOrCreateCaseStore(runID string) *CaseStore {
	detectionStoresMu.Lock()
	defer detectionStoresMu.Unlock()
	if s, ok := detectionStores[runID]; ok {
		return s
	}
	s := NewCaseStore()
	detectionStores[runID] = s
	return s
}

// handleDetectionCreateCase creates a new detection case.
func handleDetectionCreateCase(store *CaseStore, args map[string]any, tr *tracker) (string, bool) {
	if store == nil {
		return "detection: case store not configured", true
	}
	c := DetectionCase{
		AlertType:  strArg(args, "alert_type"),
		Title:      strArg(args, "title"),
		Entity:     strArg(args, "entity"),
		EntityType: strArg(args, "entity_type"),
		Severity:   strArg(args, "severity"),
		SourceData: map[string]any{},
	}
	if sd, ok := args["source_data"].(map[string]any); ok {
		c.SourceData = sd
	}
	if c.Severity == "" {
		c.Severity = "medium"
	}
	created := store.CreateCase(c)
	if tr != nil {
		tr.record("detection_create_case", args, fmt.Sprintf("created %s [%s/%s]", created.ID, created.AlertType, created.Severity))
	}
	raw, _ := json.MarshalIndent(created, "", "  ")
	return string(raw), false
}

// handleDetectionTriage processes a triage result on a case.
func handleDetectionTriage(store *CaseStore, args map[string]any, tr *tracker) (string, bool) {
	if store == nil {
		return "detection: case store not configured", true
	}
	caseID := strArg(args, "case_id")
	skillID := strArg(args, "skill_id")

	// Parse checks
	var checks []TriageCheck
	if checksRaw, ok := args["checks"].([]any); ok {
		for _, raw := range checksRaw {
			m, ok := raw.(map[string]any)
			if !ok {
				continue
			}
			checks = append(checks, TriageCheck{
				Name:    strArg(m, "name"),
				Outcome: DetectionVerdict(strArgOr(m, "outcome", "CLEAR")),
				Label:   strArg(m, "label"),
				Detail:  strArg(m, "detail"),
			})
		}
	}

	verdict := ApplyTriageDecision(checks)
	riskCount := 0
	for _, c := range checks {
		if c.Outcome == "RISK" || c.Outcome == VerdictEscalate {
			riskCount++
		}
	}

	state := TriageState{
		Verdict:   verdict,
		RiskCount: riskCount,
		Checks:    checks,
		Evidence:  strArg(args, "evidence"),
		Reasoning: strArg(args, "reasoning"),
		SkillID:   skillID,
		Timestamp: time.Now().UTC(),
	}

	if err := store.UpdateTriage(caseID, state); err != nil {
		return err.Error(), true
	}

	nextStage := DetermineNextStage(StageTriage, verdict)
	result := PipelineResult{
		CaseID:    caseID,
		Stage:     StageTriage,
		Verdict:   verdict,
		NextStage: nextStage,
		Triage:    &state,
	}

	if tr != nil {
		tr.record("detection_triage", args, fmt.Sprintf("case=%s verdict=%s risk=%d/%d next=%s",
			caseID, verdict, riskCount, len(checks), nextStage))
	}

	raw, _ := json.MarshalIndent(result, "", "  ")
	return string(raw), false
}

// handleDetectionInvestigate processes an investigation result.
func handleDetectionInvestigate(store *CaseStore, args map[string]any, tr *tracker) (string, bool) {
	if store == nil {
		return "detection: case store not configured", true
	}
	caseID := strArg(args, "case_id")
	skillID := strArg(args, "skill_id")

	var signals []InvestigationSignal
	if sigsRaw, ok := args["signals"].([]any); ok {
		for _, raw := range sigsRaw {
			m, ok := raw.(map[string]any)
			if !ok {
				continue
			}
			fired, _ := m["fired"].(bool)
			signals = append(signals, InvestigationSignal{
				Name:   strArg(m, "name"),
				Fired:  fired,
				Detail: strArg(m, "detail"),
			})
		}
	}

	verdict := ApplyInvestigationVerdict(signals)
	var actions []string
	if actionsRaw, ok := args["actions"].([]any); ok {
		for _, raw := range actionsRaw {
			if s, ok := raw.(string); ok {
				actions = append(actions, s)
			}
		}
	}

	state := InvestigationState{
		Verdict:   verdict,
		Signals:   signals,
		Evidence:  strArg(args, "evidence"),
		Reasoning: strArg(args, "reasoning"),
		Actions:   actions,
		SkillID:   skillID,
		Timestamp: time.Now().UTC(),
	}

	if err := store.UpdateInvestigation(caseID, state); err != nil {
		return err.Error(), true
	}

	nextStage := DetermineNextStage(StageInvestigation, verdict)
	result := PipelineResult{
		CaseID:         caseID,
		Stage:          StageInvestigation,
		Verdict:        verdict,
		NextStage:      nextStage,
		Investigation:  &state,
	}

	if tr != nil {
		tr.record("detection_investigate", args, fmt.Sprintf("case=%s verdict=%s signals=%d next=%s",
			caseID, verdict, len(signals), nextStage))
	}

	raw, _ := json.MarshalIndent(result, "", "  ")
	return string(raw), false
}

// handleDetectionTune processes a tuning proposal.
func handleDetectionTune(store *CaseStore, args map[string]any, tr *tracker) (string, bool) {
	if store == nil {
		return "detection: case store not configured", true
	}
	caseID := strArg(args, "case_id")

	state := TuningState{
		Action:    strArg(args, "action"),
		Target:    strArg(args, "target"),
		Value:     strArg(args, "value"),
		Rationale: strArg(args, "rationale"),
		SkillID:   strArg(args, "skill_id"),
		Timestamp: time.Now().UTC(),
	}

	if err := store.UpdateTuning(caseID, state); err != nil {
		return err.Error(), true
	}

	result := PipelineResult{
		CaseID: caseID,
		Stage:  StageTuning,
		Tuning: &state,
	}

	if tr != nil {
		tr.record("detection_tune", args, fmt.Sprintf("case=%s action=%s target=%s",
			caseID, state.Action, state.Target))
	}

	raw, _ := json.MarshalIndent(result, "", "  ")
	return string(raw), false
}

// handleDetectionListCases lists all detection cases.
func handleDetectionListCases(store *CaseStore, args map[string]any, tr *tracker) (string, bool) {
	if store == nil {
		return "detection: case store not configured", true
	}
	alertType := strArg(args, "alert_type")
	cases := store.ListCases(alertType)

	type caseSummary struct {
		ID         string            `json:"id"`
		AlertType  string            `json:"alert_type"`
		Title      string            `json:"title"`
		Entity     string            `json:"entity"`
		Severity   string            `json:"severity"`
		TriageVerdict   *DetectionVerdict `json:"triage_verdict,omitempty"`
		InvestigationVerdict *DetectionVerdict `json:"investigation_verdict,omitempty"`
		TuningAction     *string          `json:"tuning_action,omitempty"`
	}

	summaries := make([]caseSummary, 0, len(cases))
	for _, c := range cases {
		s := caseSummary{
			ID: c.ID, AlertType: c.AlertType, Title: c.Title,
			Entity: c.Entity, Severity: c.Severity,
		}
		if c.TriageState != nil {
			v := c.TriageState.Verdict
			s.TriageVerdict = &v
		}
		if c.InvestigationState != nil {
			v := c.InvestigationState.Verdict
			s.InvestigationVerdict = &v
		}
		if c.TuningState != nil {
			s.TuningAction = &c.TuningState.Action
		}
		summaries = append(summaries, s)
	}

	payload := map[string]any{
		"total":   len(summaries),
		"cases":   summaries,
	}
	if tr != nil {
		tr.record("detection_list_cases", args, fmt.Sprintf("total=%d", len(summaries)))
	}
	raw, _ := json.MarshalIndent(payload, "", "  ")
	return string(raw), false
}
