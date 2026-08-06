package core

import (
	"testing"
)

func TestCaseStore_CreateAndGet(t *testing.T) {
	cs := NewCaseStore()
	c := cs.CreateCase(DetectionCase{
		AlertType: "mfa-fatigue",
		Title:     "Multiple MFA denials",
		Entity:    "user@example.com",
		Severity:  "high",
	})
	if c.ID == "" {
		t.Error("expected non-empty ID")
	}
	got, ok := cs.GetCase(c.ID)
	if !ok {
		t.Fatal("expected to find case")
	}
	if got.AlertType != "mfa-fatigue" {
		t.Errorf("expected alert_type mfa-fatigue, got %s", got.AlertType)
	}
}

func TestApplyTriageDecision_Escalate(t *testing.T) {
	checks := []TriageCheck{
		{Name: "VIP", Outcome: "RISK"},
		{Name: "Velocity", Outcome: "RISK"},
		{Name: "Baseline", Outcome: "CLEAR"},
	}
	verdict := ApplyTriageDecision(checks)
	if verdict != VerdictEscalate {
		t.Errorf("expected escalate for 2/3 risk, got %s", verdict)
	}
}

func TestApplyTriageDecision_Dismiss(t *testing.T) {
	checks := []TriageCheck{
		{Name: "VIP", Outcome: "CLEAR"},
		{Name: "Velocity", Outcome: "CLEAR"},
		{Name: "Baseline", Outcome: "CLEAR"},
	}
	verdict := ApplyTriageDecision(checks)
	if verdict != VerdictDismiss {
		t.Errorf("expected dismiss for 0/3 risk, got %s", verdict)
	}
}

func TestApplyInvestigationVerdict_Malicious(t *testing.T) {
	signals := []InvestigationSignal{
		{Name: "A", Fired: true}, // Beacon confirmed
		{Name: "B", Fired: true}, // Bad infrastructure
	}
	verdict := ApplyInvestigationVerdict(signals)
	if verdict != VerdictMalicious {
		t.Errorf("expected malicious for A+B, got %s", verdict)
	}
}

func TestApplyInvestigationVerdict_Benign(t *testing.T) {
	signals := []InvestigationSignal{
		{Name: "A", Fired: false},
		{Name: "B", Fired: false},
		{Name: "C", Fired: false},
		{Name: "D", Fired: false},
	}
	verdict := ApplyInvestigationVerdict(signals)
	if verdict != VerdictBenign {
		t.Errorf("expected benign for all clear, got %s", verdict)
	}
}

func TestApplyInvestigationVerdict_Suspicious(t *testing.T) {
	signals := []InvestigationSignal{
		{Name: "A", Fired: true}, // Beacon only
		{Name: "B", Fired: false},
		{Name: "C", Fired: false},
		{Name: "D", Fired: false},
	}
	verdict := ApplyInvestigationVerdict(signals)
	if verdict != VerdictSuspicious {
		t.Errorf("expected suspicious for A-only, got %s", verdict)
	}
}

func TestDetermineNextStage(t *testing.T) {
	tests := []struct {
		stage   PipelineStage
		verdict DetectionVerdict
		want    PipelineStage
	}{
		{StageTriage, VerdictEscalate, StageInvestigation},
		{StageTriage, VerdictDismiss, StageTuning},
		{StageInvestigation, VerdictMalicious, StageTuning},
		{StageInvestigation, VerdictBenign, ""},
		{StageTuning, "", ""},
	}
	for _, tt := range tests {
		got := DetermineNextStage(tt.stage, tt.verdict)
		if got != tt.want {
			t.Errorf("DetermineNextStage(%s, %s) = %s, want %s", tt.stage, tt.verdict, got, tt.want)
		}
	}
}

func TestCaseStore_UpdateTriage(t *testing.T) {
	cs := NewCaseStore()
	c := cs.CreateCase(DetectionCase{AlertType: "c2-beacon", Title: "Test", Entity: "host1"})

	err := cs.UpdateTriage(c.ID, TriageState{
		Verdict: VerdictEscalate,
		Checks:  []TriageCheck{{Name: "test", Outcome: "RISK"}},
	})
	if err != nil {
		t.Fatalf("UpdateTriage: %v", err)
	}

	got, _ := cs.GetCase(c.ID)
	if got.TriageState == nil {
		t.Fatal("expected triage state set")
	}
	if got.TriageState.Verdict != VerdictEscalate {
		t.Errorf("expected escalate, got %s", got.TriageState.Verdict)
	}
}

func TestFormatVerdictSummary(t *testing.T) {
	r := PipelineResult{
		CaseID:  "case-1",
		Stage:   StageTriage,
		Verdict: VerdictEscalate,
		Triage: &TriageState{
			Checks:    []TriageCheck{{}, {}, {}},
			RiskCount: 2,
		},
		NextStage: StageInvestigation,
	}
	summary := FormatVerdictSummary(r)
	if summary == "" {
		t.Error("expected non-empty summary")
	}
}
