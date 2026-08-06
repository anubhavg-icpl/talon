package core

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestTargetStore_CRUD(t *testing.T) {
	dir := t.TempDir()
	ts := NewTargetStore(dir)

	// Create
	state, err := ts.GetOrCreate("10.0.0.1")
	if err != nil {
		t.Fatalf("GetOrCreate: %v", err)
	}
	if state.Target != "10.0.0.1" {
		t.Errorf("expected target 10.0.0.1, got %s", state.Target)
	}

	// Modify
	state.AddFinding(TargetFinding{
		Title: "SQL Injection", Severity: "high", Status: "verified",
	})
	state.AddReconDim(ReconDimension{Name: "port_scan", Status: "complete", Summary: "22,80,443 open"})

	// Save
	if err := ts.Save(state); err != nil {
		t.Fatalf("Save: %v", err)
	}

	// Verify file exists
	path := filepath.Join(dir, "10.0.0.1.json")
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected state file at %s: %v", path, err)
	}

	// Reload
	state2, err := ts.GetOrCreate("10.0.0.1")
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	if len(state2.Findings) != 1 {
		t.Errorf("expected 1 finding, got %d", len(state2.Findings))
	}
	if len(state2.ReconDims) != 1 {
		t.Errorf("expected 1 recon dim, got %d", len(state2.ReconDims))
	}
}

func TestTargetStore_Snapshot(t *testing.T) {
	dir := t.TempDir()
	ts := NewTargetStore(dir)

	state, _ := ts.GetOrCreate("example.com")
	state.AddFinding(TargetFinding{Title: "XSS", Severity: "medium", Status: "verified"})
	ts.Save(state)

	// Snapshot before modification
	ts.Snapshot(state, "before-expand")

	state.AddFinding(TargetFinding{Title: "RCE", Severity: "critical", Status: "verified"})
	ts.Save(state)

	state2, _ := ts.GetOrCreate("example.com")
	if len(state2.Snapshots) != 1 {
		t.Fatalf("expected 1 snapshot, got %d", len(state2.Snapshots))
	}
	if len(state2.Findings) != 2 {
		t.Errorf("expected 2 findings, got %d", len(state2.Findings))
	}
}

func TestTargetState_ResumePlan(t *testing.T) {
	state := &TargetState{
		Target:        "10.0.0.1",
		Slug:          "10-0-0-1",
		SchemaVersion: targetSchemaVersion,
	}

	// Empty state → should suggest initial recon
	steps := state.ResumePlan()
	if len(steps) == 0 {
		t.Fatal("expected steps for empty state")
	}
	if steps[0].Category != "recon" {
		t.Errorf("expected recon category for empty state, got %s", steps[0].Category)
	}

	// Add complete port scan
	state.AddReconDim(ReconDimension{Name: "port_scan", Status: "complete"})

	// Add a failed vector
	state.AddFailedVector(FailedVector{Vector: "sql_injection", Target: "/login", Reason: "WAF blocked"})
	state.AddFailedVector(FailedVector{Vector: "sql_injection", Target: "/search", Reason: "WAF blocked"})

	steps = state.ResumePlan()
	// Should have: DNS recon, web_enum, service_id (port_scan complete), sql_injection retry
	foundRetry := false
	for _, step := range steps {
		if step.Category == "exploit_retry" && strings.Contains(step.Action, "sql_injection") {
			foundRetry = true
		}
	}
	if !foundRetry {
		t.Error("expected exploit_retry step for failed SQL injection")
	}

	// Add verified finding → should suggest post-exploit
	state.AddFinding(TargetFinding{Title: "RCE", Severity: "critical", Status: "verified"})
	steps = state.ResumePlan()
	foundPostExploit := false
	for _, step := range steps {
		if step.Category == "post_exploit" {
			foundPostExploit = true
		}
	}
	if !foundPostExploit {
		t.Error("expected post_exploit step with verified findings")
	}
}

func TestTargetState_Summary(t *testing.T) {
	state := &TargetState{
		Target: "example.com",
		Slug:   "example-com",
	}
	state.AddFinding(TargetFinding{Title: "XSS", Severity: "medium", Status: "verified"})
	state.AddReconDim(ReconDimension{Name: "port_scan", Status: "complete"})

	summary := state.Summary()
	if !strings.Contains(summary, "example.com") {
		t.Error("summary should contain target")
	}
	if !strings.Contains(summary, "1 verified") {
		t.Error("summary should show 1 verified finding")
	}
	if !strings.Contains(summary, "port_scan=complete") {
		t.Error("summary should show recon dim status")
	}
}

func TestTargetSlug(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"10.0.0.1", "10.0.0.1"},
		{"http://example.com:8080", "http-example.com-8080"},
		{"192.168.1.100", "192.168.1.100"},
		{"test target!", "test-target"},
	}
	for _, tt := range tests {
		got := targetSlug(tt.input)
		if got != tt.want {
			t.Errorf("targetSlug(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
