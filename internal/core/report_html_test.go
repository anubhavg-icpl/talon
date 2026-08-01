package core

import (
	"strings"
	"testing"
	"time"
)

func TestRenderReportHTML(t *testing.T) {
	j := false
	r := StructuredReport{
		Target:      "10.170.0.3",
		GeneratedAt: time.Now().UTC(),
		Findings: []Finding{{
			ID: "FIND-001", Severity: SeverityHigh, Title: "DVWA Accessible",
			Status: "approved", Stage: "recon", Source: "report_finding",
			Endpoint: "http://10.170.0.3/login.php", Description: "DVWA v1.10 low",
			StepsToReproduce: "1. nmap\\n2. open login.php", Recommendation: "Isolate",
			Evidence: GateEvidence{Baseline: "filtered", Diff: "HTTP 200", Passed: false},
		}},
		Summary:       FindingsSummary{Total: 1, High: 1},
		JudgeVerdict:  &j,
		StagesCovered: []string{"recon", "exploit"},
	}
	html := RenderReportHTML(ReportHTMLInput{
		RunID:     "4e5da981-33fc-411f-a27a-08d29019f918",
		Input:     RunInput{TargetIP: "10.170.0.3", ServiceName: "http", AgentMode: "full"},
		Report:    &r,
		Findings:  r.Findings,
		ToolLog:   []ToolCallRecord{{Index: 0, ToolName: "nmap_scan", Output: "80/tcp open"}},
		FinalMsg:  "⚠️ exploit timeout",
		Status:    "completed",
		StartedAt: time.Now().Add(-time.Hour),
		EndedAt:   time.Now(),
	})
	checks := []string{
		"data:image/webp;base64,",
		"Print / Save PDF",
		"TALON AI",
		"DVWA Accessible",
		"10.170.0.3",
		"3-GATE FAIL",
		"@media print",
		"nmap_scan",
	}
	for _, c := range checks {
		if !strings.Contains(html, c) {
			t.Fatalf("missing %q in HTML (len=%d)", c, len(html))
		}
	}
	if strings.Contains(html, "<pre>%s</pre>") {
		t.Fatal("raw pre dump template still present")
	}
}
