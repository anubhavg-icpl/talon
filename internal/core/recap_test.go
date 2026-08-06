package core

import (
	"strings"
	"testing"
	"time"
)

func TestBuildRecap(t *testing.T) {
	start := time.Now().Add(-5 * time.Minute)
	tracker := &ToolCallTracker{
		Calls: []ToolCallRecord{
			{Index: 0, ToolName: "nmap_scan", Args: map[string]any{"target": "10.0.0.1"}, Output: "Port 80 open"},
			{Index: 1, ToolName: "fetch", Args: map[string]any{"url": "http://10.0.0.1"}, Output: "Web page"},
		},
	}
	store := NewEvidenceStore()
	store.Record("nmap_scan", "", "Port 80 open", 200)

	findings := []Finding{
		{Title: "SQL Injection", Severity: "high", PoC: "sqlmap -u http://target/login",
			Evidence: GateEvidence{Passed: true, Baseline: "b", Attack: "a", Diff: "d"}},
	}

	recap := BuildRecap("10.0.0.1", "run-001", start, tracker, store, findings)

	if recap.Target != "10.0.0.1" {
		t.Errorf("expected target 10.0.0.1, got %s", recap.Target)
	}
	if len(recap.SolvePath) != 2 {
		t.Errorf("expected 2 solve steps, got %d", len(recap.SolvePath))
	}
	if recap.FindingCount != 1 {
		t.Errorf("expected 1 finding, got %d", recap.FindingCount)
	}
	if recap.VerifiedCount != 1 {
		t.Errorf("expected 1 verified, got %d", recap.VerifiedCount)
	}
	if len(recap.Reproduction) != 1 {
		t.Errorf("expected 1 reproduction, got %d", len(recap.Reproduction))
	}
}

func TestRecap_FormatMarkdown(t *testing.T) {
	recap := RunRecap{
		Target:    "example.com",
		StartTime: time.Now().Add(-1 * time.Minute),
		EndTime:   time.Now(),
		Duration:  "1m0s",
		SolvePath: []RecapStep{
			{Step: 1, Action: "nmap_scan", Result: "Port 80 open"},
		},
		KeyEvidence: []RecapEvidence{
			{ID: "e001", Tool: "nmap", Summary: "Port 80 open"},
		},
		FindingCount:  2,
		VerifiedCount: 1,
	}

	md := recap.FormatMarkdown()
	if !strings.Contains(md, "Run Recap") {
		t.Error("expected 'Run Recap' in markdown")
	}
	if !strings.Contains(md, "example.com") {
		t.Error("expected target in markdown")
	}
	if !strings.Contains(md, "Solve Path") {
		t.Error("expected solve path section")
	}
	if !strings.Contains(md, "Key Evidence") {
		t.Error("expected key evidence section")
	}
}

func TestPresets(t *testing.T) {
	presets := Presets()
	if len(presets) != 3 {
		t.Fatalf("expected 3 presets, got %d", len(presets))
	}
	quick := presets["quick"]
	if quick.MaxTurns != 15 {
		t.Errorf("quick preset: expected 15 turns, got %d", quick.MaxTurns)
	}
	deep := presets["deep"]
	if deep.MaxTurns != 60 {
		t.Errorf("deep preset: expected 60 turns, got %d", deep.MaxTurns)
	}
}
