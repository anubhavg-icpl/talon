package core

import (
	"fmt"
	"strings"
)

// RunSnapshot is a lightweight run view for compare/export.
type RunSnapshot struct {
	Target        string           `json:"target"`
	CVEID         string           `json:"cve_id,omitempty"`
	ServiceName   string           `json:"service_name,omitempty"`
	AgentMode     string           `json:"agent_mode,omitempty"`
	JudgeVerdict  *bool            `json:"judge_verdict,omitempty"`
	FindingsCount int              `json:"findings_count"`
	Summary       FindingsSummary  `json:"findings_summary"`
	Findings      []Finding        `json:"findings,omitempty"`
	Stages        []string         `json:"stages_covered,omitempty"`
	ToolCount     int              `json:"tool_count"`
	KillChain     *KillChainAnalysis `json:"kill_chain,omitempty"`
	Methodology   *MethodologyState  `json:"methodology,omitempty"`
}

// CompareResult is a side-by-side analysis of two runs.
type CompareResult struct {
	A              RunSnapshot `json:"a"`
	B              RunSnapshot `json:"b"`
	OnlyInA        []Finding   `json:"only_in_a"`
	OnlyInB        []Finding   `json:"only_in_b"`
	SharedTitles   []string    `json:"shared_titles"`
	SeverityDelta  map[string]int `json:"severity_delta"` // B-A counts by severity
	VerdictChange  string      `json:"verdict_change"`
	Markdown       string      `json:"markdown"`
}

// BuildSnapshot packs run inputs into a comparable snapshot.
func BuildSnapshot(input RunInput, findings []Finding, toolLog []ToolCallRecord, judge *bool, kc *KillChainAnalysis, meth *MethodologyState) RunSnapshot {
	sum := SummarizeFindings(findings)
	return RunSnapshot{
		Target:        input.TargetIP,
		CVEID:         input.CVEID,
		ServiceName:   input.ServiceName,
		AgentMode:     NormalizeAgentMode(input.AgentMode),
		JudgeVerdict:  judge,
		FindingsCount: sum.Total,
		Summary:       sum,
		Findings:      findings,
		Stages:        coveredStages(toolLog),
		ToolCount:     len(toolLog),
		KillChain:     kc,
		Methodology:   meth,
	}
}

// CompareRuns diffs two snapshots (CyberStrike-inspired dual engagement review).
func CompareRuns(a, b RunSnapshot) CompareResult {
	titlesA := map[string]Finding{}
	titlesB := map[string]Finding{}
	for _, f := range a.Findings {
		titlesA[strings.ToLower(f.Title)] = f
	}
	for _, f := range b.Findings {
		titlesB[strings.ToLower(f.Title)] = f
	}

	var onlyA, onlyB []Finding
	var shared []string
	for k, f := range titlesA {
		if _, ok := titlesB[k]; ok {
			shared = append(shared, f.Title)
		} else {
			onlyA = append(onlyA, f)
		}
	}
	for k, f := range titlesB {
		if _, ok := titlesA[k]; !ok {
			onlyB = append(onlyB, f)
		}
	}

	delta := map[string]int{
		"critical": b.Summary.Critical - a.Summary.Critical,
		"high":     b.Summary.High - a.Summary.High,
		"medium":   b.Summary.Medium - a.Summary.Medium,
		"low":      b.Summary.Low - a.Summary.Low,
		"info":     b.Summary.Info - a.Summary.Info,
		"total":    b.Summary.Total - a.Summary.Total,
	}

	verdict := "unchanged"
	av, bv := false, false
	if a.JudgeVerdict != nil {
		av = *a.JudgeVerdict
	}
	if b.JudgeVerdict != nil {
		bv = *b.JudgeVerdict
	}
	switch {
	case a.JudgeVerdict == nil && b.JudgeVerdict == nil:
		verdict = "both_unevaluated"
	case !av && bv:
		verdict = "gained_compromise"
	case av && !bv:
		verdict = "lost_compromise"
	case av && bv:
		verdict = "both_compromised"
	default:
		verdict = "both_clean"
	}

	var md strings.Builder
	md.WriteString("# Run compare\n\n")
	fmt.Fprintf(&md, "| | **A** `%s` | **B** `%s` |\n| --- | --- | --- |\n", a.Target, b.Target)
	fmt.Fprintf(&md, "| Mode | %s | %s |\n", a.AgentMode, b.AgentMode)
	fmt.Fprintf(&md, "| Findings | %d | %d |\n", a.FindingsCount, b.FindingsCount)
	fmt.Fprintf(&md, "| Tools | %d | %d |\n", a.ToolCount, b.ToolCount)
	fmt.Fprintf(&md, "| Verdict | %s | %s |\n", boolStr(a.JudgeVerdict), boolStr(b.JudgeVerdict))
	fmt.Fprintf(&md, "\n**Verdict change:** %s\n", verdict)
	fmt.Fprintf(&md, "\n**Shared titles:** %d · **Only A:** %d · **Only B:** %d\n", len(shared), len(onlyA), len(onlyB))

	return CompareResult{
		A: a, B: b,
		OnlyInA: onlyA, OnlyInB: onlyB, SharedTitles: shared,
		SeverityDelta: delta, VerdictChange: verdict, Markdown: md.String(),
	}
}

func boolStr(p *bool) string {
	if p == nil {
		return "n/a"
	}
	if *p {
		return "compromised"
	}
	return "clean"
}

// ExportBundle is a portable JSON package for a single run.
type ExportBundle struct {
	Version     string              `json:"version"`
	RunID       string              `json:"run_id"`
	Snapshot    RunSnapshot         `json:"snapshot"`
	ReportMD    string              `json:"report_markdown,omitempty"`
	History     []string            `json:"history,omitempty"`
	ToolLog     []ToolCallRecord    `json:"tool_log,omitempty"`
	Notes       []OperatorNote      `json:"notes,omitempty"`
}

// OperatorNote is a free-form operator annotation on a run.
type OperatorNote struct {
	ID        string `json:"id"`
	Author    string `json:"author,omitempty"`
	Body      string `json:"body"`
	CreatedAt string `json:"created_at"`
}

// TimelineEvent is one attack-timeline row derived from tools + findings.
type TimelineEvent struct {
	Index   int    `json:"index"`
	Kind    string `json:"kind"` // tool | finding | status
	Label   string `json:"label"`
	Stage   string `json:"stage,omitempty"`
	Detail  string `json:"detail,omitempty"`
	Severity string `json:"severity,omitempty"`
}

// BuildTimeline merges tool log and findings into a chronological operator view.
func BuildTimeline(toolLog []ToolCallRecord, findings []Finding) []TimelineEvent {
	events := make([]TimelineEvent, 0, len(toolLog)+len(findings))
	for _, rec := range toolLog {
		detail := rec.Output
		if len(detail) > 240 {
			detail = detail[:240] + "…"
		}
		events = append(events, TimelineEvent{
			Index:  rec.Index,
			Kind:   "tool",
			Label:  rec.ToolName,
			Stage:  stageForTool(rec.ToolName),
			Detail: detail,
		})
	}
	// Findings after tools (extraction order), index continues
	base := len(toolLog)
	for i, f := range findings {
		events = append(events, TimelineEvent{
			Index:    base + i,
			Kind:     "finding",
			Label:    f.Title,
			Stage:    f.Stage,
			Detail:   f.Description,
			Severity: f.Severity,
		})
	}
	return events
}
