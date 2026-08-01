package core

import (
	"strings"
	"testing"
)

func TestValidateThreeGate(t *testing.T) {
	e, v := ValidateThreeGate(GateEvidence{
		Baseline: "no session",
		Attack:   "Session 1 created",
		Diff:     "session differs from empty baseline",
	})
	if !e.Passed || len(v) != 0 {
		t.Fatalf("expected pass, got passed=%v violations=%v", e.Passed, v)
	}

	e, v = ValidateThreeGate(GateEvidence{Baseline: "x", Attack: "y"})
	if e.Passed || len(v) == 0 {
		t.Fatalf("expected fail without diff")
	}
}

func TestExtractFindingsSessionAndPorts(t *testing.T) {
	input := RunInput{
		TargetIP:    "10.0.0.5",
		CVEID:       "CVE-2011-2523",
		ServiceName: "vsftpd 2.3.4",
	}
	log := []ToolCallRecord{
		{Index: 0, ToolName: "nmap_scan", Output: "21/tcp open ftp\n22/tcp open ssh\n"},
		{Index: 1, ToolName: "run_exploit", Output: "[*] Exploit completed.\n[*] Session 1 created in the background.\n"},
		{Index: 2, ToolName: "send_session_command", Output: "uid=0(root) gid=0(root) groups=0(root)\n"},
	}
	findings := ExtractFindings(input, log, "rooted the box", true, true)
	if len(findings) < 3 {
		t.Fatalf("expected >=3 findings, got %d: %+v", len(findings), findings)
	}
	var hasCritical, hasPort, hasProof bool
	for _, f := range findings {
		if f.Severity == SeverityCritical && f.Evidence.Passed {
			hasCritical = true
		}
		if strings.Contains(f.Title, "Open port") {
			hasPort = true
		}
		if strings.Contains(f.Title, "proof") || strings.Contains(f.Title, "Post-exploit") {
			hasProof = true
		}
		if !f.Evidence.Passed && f.Severity == SeverityCritical {
			t.Errorf("critical finding %s failed 3-gate: %+v", f.ID, f.Evidence)
		}
	}
	if !hasCritical {
		t.Error("expected critical RCE finding")
	}
	if !hasPort {
		t.Error("expected open-port findings")
	}
	if !hasProof {
		t.Error("expected post-exploit proof finding")
	}

	sum := SummarizeFindings(findings)
	if sum.Total != len(findings) || sum.Confirmed == 0 {
		t.Fatalf("bad summary: %+v", sum)
	}
}

func TestExtractFindingsNuclei(t *testing.T) {
	input := RunInput{TargetIP: "1.2.3.4"}
	log := []ToolCallRecord{
		{ToolName: "nuclei_scan", Output: "[critical] CVE-2021-41773 Apache Path Traversal\n[medium] missing-security-headers\n"},
	}
	findings := ExtractFindings(input, log, "", false, false)
	if len(findings) < 2 {
		t.Fatalf("expected nuclei findings, got %d", len(findings))
	}
	if findings[0].Severity != SeverityCritical {
		t.Errorf("first severity=%s", findings[0].Severity)
	}
}

func TestDedupFindings(t *testing.T) {
	in := []Finding{
		{Title: "Open port 21", Endpoint: "1.1.1.1:21", AttackVector: "other"},
		{Title: "Open port 21", Endpoint: "1.1.1.1:21", AttackVector: "other"},
		{Title: "RCE", Endpoint: "1.1.1.1", AttackVector: "other"},
	}
	out := DedupFindings(in)
	if len(out) != 2 {
		t.Fatalf("expected 2 after dedup, got %d", len(out))
	}
}

func TestBuildReport(t *testing.T) {
	input := RunInput{TargetIP: "127.0.0.1", CVEID: "CVE-2011-2523", ServiceName: "vsftpd"}
	log := []ToolCallRecord{
		{Index: 0, ToolName: "nmap_scan", Output: "21/tcp open ftp\n"},
		{Index: 1, ToolName: "run_exploit", Output: "Session 2 created\n"},
	}
	findings := ExtractFindings(input, log, "session open", true, true)
	rep := BuildReport(input, log, "session open", findings, true, true)
	if rep.Markdown == "" {
		t.Fatal("empty markdown")
	}
	for _, sec := range []string{"Executive Summary", "Findings", "Methodology", "Timeline", "Validation"} {
		if !strings.Contains(rep.Markdown, sec) {
			t.Errorf("missing section %q", sec)
		}
	}
	if rep.JudgeVerdict == nil || !*rep.JudgeVerdict {
		t.Error("expected judge true")
	}
	if rep.Summary.Total == 0 {
		t.Error("expected findings in summary")
	}
}

func TestInjectSkills(t *testing.T) {
	p := InjectSkills("base prompt", "recon")
	if !strings.Contains(p, "3-Gate") {
		t.Error("expected 3-gate skill in all stages")
	}
	if !strings.Contains(p, "Recon Methodology") {
		t.Error("expected recon methodology")
	}
	if !strings.Contains(p, "base prompt") {
		t.Error("base preserved")
	}
	skills := ListSkills()
	if len(skills) < 5 {
		t.Fatalf("expected skill catalog, got %d", len(skills))
	}
}

func TestFindingBagAndMerge(t *testing.T) {
	bag := NewFindingBag()
	f := bag.Add(Finding{
		Severity: SeverityCritical, Title: "RCE", Description: "session",
		Endpoint: "1.1.1.1", AttackVector: "other",
		Evidence: GateEvidence{Baseline: "none", Attack: "session 1", Diff: "created"},
	})
	if f.ID == "" || !f.Evidence.Passed {
		t.Fatalf("bad add: %+v", f)
	}
	triaged, err := bag.Triage(f.ID, FindingStatusApproved, "")
	if err != nil || triaged.Status != FindingStatusApproved {
		t.Fatalf("triage: %v %+v", err, triaged)
	}
	extracted := []Finding{{Title: "RCE", Endpoint: "1.1.1.1", AttackVector: "other", Severity: SeverityHigh}}
	merged := MergeExtracted(bag.Snapshot(), extracted)
	if len(merged) != 1 {
		t.Fatalf("dedupe merge expected 1, got %d", len(merged))
	}
}

func TestAgentModeAndKillChain(t *testing.T) {
	if NormalizeAgentMode("web-application") != AgentModeWeb {
		t.Fatal("web normalize")
	}
	allowed := AllowedDelegates(AgentModeRecon)
	if allowed["delegate_exploit"] {
		t.Fatal("recon should not allow exploit")
	}
	findings := []Finding{
		{Stage: "recon", Severity: SeverityInfo, Title: "port"},
		{Stage: "exploit", Severity: SeverityCritical, Title: "session"},
		{Stage: "post_exploit", Severity: SeverityCritical, Title: "proof"},
	}
	kc := AnalyzeKillChain(findings)
	if len(kc.Chains) == 0 {
		t.Fatal("expected chains")
	}
	m := ComputeMethodology([]ToolCallRecord{{ToolName: "nmap_scan"}, {ToolName: "run_exploit"}}, AgentModeFull)
	if m.CoveredCount < 2 {
		t.Fatalf("coverage %+v", m)
	}
}

func TestSkillCatalogLoaded(t *testing.T) {
	all := ListSkills()
	if len(all) < 100 {
		t.Fatalf("expected large catalog, got %d (cwd skills/ present?)", len(all))
	}
	r := QuerySkills(SkillQuery{Brief: true, Limit: 10, Category: "WEB"})
	if r.Total == 0 {
		// category may be named differently
		r = QuerySkills(SkillQuery{Brief: true, Limit: 10, Q: "ssrf"})
	}
	if r.Total == 0 {
		t.Fatal("expected search hits for ssrf or WEB")
	}
	if len(r.Categories) == 0 {
		t.Fatal("expected categories")
	}
}

func TestSkillSearchAndGetTools(t *testing.T) {
	out, isErr := handleSkillSearch(map[string]any{"q": "ssrf", "limit": float64(5)}, nil)
	if isErr {
		t.Fatalf("search error: %s", out)
	}
	if !strings.Contains(out, "hits") && !strings.Contains(out, "total") {
		t.Fatalf("unexpected search payload: %s", out[:minInt(200, len(out))])
	}
	// pick an id if any
	res := QuerySkills(SkillQuery{Q: "ssrf", Brief: true, Limit: 1})
	if res.Total == 0 || len(res.Skills) == 0 {
		t.Skip("no ssrf skills in catalog (skills/ missing?)")
	}
	id := res.Skills[0].ID
	body, isErr := handleSkillGet(map[string]any{"id": id}, nil)
	if isErr {
		t.Fatalf("get error: %s", body)
	}
	if !strings.Contains(body, "body") {
		t.Fatalf("expected body in skill_get: %s", body[:minInt(200, len(body))])
	}
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func TestCompareAndTimelineAndPlaybooks(t *testing.T) {
	if len(ListPlaybooks()) < 3 {
		t.Fatal("expected playbooks")
	}
	a := BuildSnapshot(RunInput{TargetIP: "1.1.1.1", AgentMode: "full"}, []Finding{
		{Title: "RCE", Severity: SeverityCritical},
	}, []ToolCallRecord{{Index: 0, ToolName: "nmap_scan", Output: "open"}}, boolPtr(true), nil, nil)
	b := BuildSnapshot(RunInput{TargetIP: "1.1.1.1", AgentMode: "web"}, []Finding{
		{Title: "RCE", Severity: SeverityCritical},
		{Title: "Open port", Severity: SeverityInfo},
	}, []ToolCallRecord{{Index: 0, ToolName: "run_exploit", Output: "session"}}, boolPtr(true), nil, nil)
	cmp := CompareRuns(a, b)
	if cmp.SeverityDelta["total"] != 1 {
		t.Fatalf("delta total=%v", cmp.SeverityDelta)
	}
	tl := BuildTimeline([]ToolCallRecord{{Index: 0, ToolName: "nmap_scan", Output: "open"}}, a.Findings)
	if len(tl) < 2 {
		t.Fatalf("timeline len=%d", len(tl))
	}
}

func boolPtr(v bool) *bool { return &v }
