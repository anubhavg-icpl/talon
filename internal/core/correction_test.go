package core

import (
	"strings"
	"testing"
)

func TestCorrection_DuplicateDetection(t *testing.T) {
	c := NewCorrectionLayer()
	args := map[string]any{"target": "10.0.0.1"}

	// First call - no hint
	hint := c.BeforeToolCall("nmap_scan", args)
	if hint != "" {
		t.Errorf("expected empty hint on first call, got: %s", hint)
	}

	// Record two calls
	c.AfterToolCall("nmap_scan", args, "Host is up", false, 0)
	c.AfterToolCall("nmap_scan", args, "Host is up", false, 0)

	// Third call should warn about duplication
	hint = c.BeforeToolCall("nmap_scan", args)
	if !strings.Contains(hint, "already appeared") {
		t.Errorf("expected duplicate warning, got: %s", hint)
	}
}

func TestCorrection_DegradedHealth(t *testing.T) {
	c := NewCorrectionLayer()
	args := map[string]any{"ip": "10.0.0.1"}

	// Record 3 failures
	for i := 0; i < 3; i++ {
		c.AfterToolCall("nmap_scan", args, "connection refused", false, 0)
	}

	// Next call should warn about degraded state
	hint := c.BeforeToolCall("nmap_scan", args)
	if !strings.Contains(hint, "degraded") {
		t.Errorf("expected degraded warning, got: %s", hint)
	}
}

func TestCorrection_HealthyRecovery(t *testing.T) {
	c := NewCorrectionLayer()
	args := map[string]any{"ip": "10.0.0.1"}

	// Record 2 failures then a success
	c.AfterToolCall("nmap_scan", args, "timeout", false, 0)
	c.AfterToolCall("nmap_scan", args, "error", false, 0)
	c.AfterToolCall("nmap_scan", args, "Host is up", false, 0)

	// Should NOT be degraded anymore
	hint := c.BeforeToolCall("nmap_scan", args)
	if strings.Contains(hint, "degraded") {
		t.Errorf("should not be degraded after success, got: %s", hint)
	}
}

func TestCorrection_StallDetection(t *testing.T) {
	c := NewCorrectionLayer()

	// First, record some active tool calls
	c.AfterToolCall("nmap_scan", nil, "Host is up", false, 0)
	c.AfterToolCall("fetch", nil, "OK", false, 0)

	// Then only evidence reads
	c.AfterToolCall("evidence_list", nil, "{}", false, 0)
	c.AfterToolCall("evidence_view", map[string]any{"id": "e001"}, "content", false, 0)
	c.AfterToolCall("evidence_view", map[string]any{"id": "e002"}, "content", false, 0)
	c.AfterToolCall("evidence_search", map[string]any{"q": "flag"}, "results", false, 0)

	hint := c.AfterToolCall("evidence_list", nil, "{}", false, 0)
	if !strings.Contains(hint, "Stall detected") {
		t.Errorf("expected stall detection, got: %s", hint)
	}
}

func TestCorrection_CompletionGate_FlagInEvidence(t *testing.T) {
	c := NewCorrectionLayer()
	store := NewEvidenceStore()
	store.Record("fetch", "", "Found flag{test_flag_123} in response", 200)

	// Claim that includes the flag should pass
	ok, feedback := c.ValidateCompletion("I found the flag{test_flag_123}", store)
	if !ok {
		t.Errorf("expected validation to pass, got feedback: %s", feedback)
	}
	if feedback != "" {
		t.Errorf("expected empty feedback on success, got: %s", feedback)
	}
}

func TestCorrection_CompletionGate_FlagNotInEvidence(t *testing.T) {
	c := NewCorrectionLayer()
	store := NewEvidenceStore()
	store.Record("fetch", "", "No flags here", 200)

	// Claim with a flag not in evidence should be rejected
	ok, feedback := c.ValidateCompletion("I found the flag{fake_flag}", store)
	if ok {
		t.Error("expected validation to fail for unbacked flag claim")
	}
	if !strings.Contains(feedback, "REJECTION") {
		t.Errorf("expected rejection feedback, got: %s", feedback)
	}
}

func TestCorrection_CompletionGate_AfterMaxRetries(t *testing.T) {
	c := NewCorrectionLayer()
	store := NewEvidenceStore()

	// Exhaust retries
	for i := 0; i < maxRejectionRetries+1; i++ {
		c.ValidateCompletion("flag{fake}", store)
	}

	// Should now accept with unverified flag
	ok, feedback := c.ValidateCompletion("flag{fake}", store)
	if !ok {
		t.Error("expected acceptance after max retries")
	}
	if !strings.Contains(feedback, "UNVERIFIED") {
		t.Errorf("expected unverified flag in feedback, got: %s", feedback)
	}
}

func TestCorrection_EvidenceViewRedundancy(t *testing.T) {
	c := NewCorrectionLayer()

	// View same evidence multiple times
	c.AfterToolCall("evidence_view", map[string]any{"id": "e001"}, "content", false, 0)
	c.AfterToolCall("evidence_view", map[string]any{"id": "e001"}, "content", false, 0)

	hint := c.BeforeToolCall("evidence_view", map[string]any{"id": "e001"})
	if !strings.Contains(hint, "already called") {
		t.Errorf("expected redundancy warning, got: %s", hint)
	}
}
