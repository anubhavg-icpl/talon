package core

import (
	"fmt"
	"strings"
	"time"
)

// StructuredReport is the portable pentest report produced at end of run
// (CyberStrike generate_report sections, adapted for Talon's pipeline).
type StructuredReport struct {
	// Markdown is the full human-readable report.
	Markdown string `json:"markdown"`
	// GeneratedAt is when the report was built.
	GeneratedAt time.Time `json:"generated_at"`
	// Sections lists included section names.
	Sections []string `json:"sections"`
	// Findings is the structured finding list embedded in the report.
	Findings []Finding `json:"findings"`
	// Summary is severity roll-up.
	Summary FindingsSummary `json:"summary"`
	// JudgeVerdict mirrors the run judge (nil if not set).
	JudgeVerdict *bool `json:"judge_verdict,omitempty"`
	// Target is a short target label for UI headers.
	Target string `json:"target"`
	// CVEID is the CVE under test, if any.
	CVEID string `json:"cve_id,omitempty"`
	// Duration tools / stages covered.
	StagesCovered []string `json:"stages_covered,omitempty"`
}

// BuildReport assembles a CyberStrike-style multi-section validation report
// from run input, tool log, final message, findings, and judge result.
func BuildReport(input RunInput, toolLog []ToolCallRecord, finalMessage string, findings []Finding, judgeVerdict bool, judgeSet bool) StructuredReport {
	now := time.Now().UTC()
	summary := SummarizeFindings(findings)
	stages := coveredStages(toolLog)

	var judgePtr *bool
	if judgeSet {
		v := judgeVerdict
		judgePtr = &v
	}

	sections := []string{
		"executive_summary",
		"findings",
		"methodology",
		"timeline",
		"validation",
		"agent_output",
	}

	var b strings.Builder
	writeExecutiveSummary(&b, input, summary, judgePtr, stages)
	writeFindingsSection(&b, findings, summary)
	writeMethodologySection(&b, stages, toolLog)
	writeTimelineSection(&b, toolLog)
	writeValidationSection(&b, findings, judgePtr)
	writeAgentOutputSection(&b, finalMessage)

	return StructuredReport{
		Markdown:      b.String(),
		GeneratedAt:   now,
		Sections:      sections,
		Findings:      findings,
		Summary:       summary,
		JudgeVerdict:  judgePtr,
		Target:        input.TargetIP,
		CVEID:         input.CVEID,
		StagesCovered: stages,
	}
}

func writeExecutiveSummary(b *strings.Builder, input RunInput, summary FindingsSummary, judge *bool, stages []string) {
	b.WriteString("# Talon Validation Report\n\n")
	b.WriteString("## Executive Summary\n\n")
	fmt.Fprintf(b, "| Field | Value |\n| --- | --- |\n")
	fmt.Fprintf(b, "| **Target** | `%s` |\n", input.TargetIP)
	if input.CVEID != "" {
		fmt.Fprintf(b, "| **CVE** | `%s` |\n", input.CVEID)
	}
	if input.ServiceName != "" {
		fmt.Fprintf(b, "| **Service** | %s |\n", input.ServiceName)
	}
	fmt.Fprintf(b, "| **Findings** | %d total (%d critical, %d high, %d medium, %d low, %d info) |\n",
		summary.Total, summary.Critical, summary.High, summary.Medium, summary.Low, summary.Info)
	fmt.Fprintf(b, "| **3-Gate Confirmed** | %d |\n", summary.Confirmed)
	if judge != nil {
		status := "NOT COMPROMISED"
		if *judge {
			status = "COMPROMISED"
		}
		fmt.Fprintf(b, "| **Judge Verdict** | **%s** |\n", status)
	} else {
		b.WriteString("| **Judge Verdict** | _not evaluated_ |\n")
	}
	if len(stages) > 0 {
		fmt.Fprintf(b, "| **Stages covered** | %s |\n", strings.Join(stages, " → "))
	}
	b.WriteString("\n")

	// One-paragraph narrative.
	switch {
	case judge != nil && *judge:
		fmt.Fprintf(b, "The validation pipeline **confirmed compromise** of `%s`", input.TargetIP)
		if input.CVEID != "" {
			fmt.Fprintf(b, " under %s", input.CVEID)
		}
		fmt.Fprintf(b, ". %d finding(s) were recorded; %d passed the 3-gate evidence protocol.\n\n",
			summary.Total, summary.Confirmed)
	case summary.Critical > 0:
		fmt.Fprintf(b, "Critical-severity evidence was observed on `%s`, but the judge did not confirm full objective completion. Review findings carefully.\n\n", input.TargetIP)
	default:
		fmt.Fprintf(b, "No confirmed compromise of `%s` in this run. Methodology stages executed: %s.\n\n",
			input.TargetIP, strings.Join(stages, ", "))
	}
}

func writeFindingsSection(b *strings.Builder, findings []Finding, summary FindingsSummary) {
	b.WriteString("## Findings\n\n")
	if len(findings) == 0 {
		b.WriteString("_No structured findings extracted from this run._\n\n")
		return
	}
	fmt.Fprintf(b, "**%d** finding(s) — **%d** 3-gate confirmed.\n\n", summary.Total, summary.Confirmed)

	// Order: critical → high → medium → low → info
	order := []string{SeverityCritical, SeverityHigh, SeverityMedium, SeverityLow, SeverityInfo}
	for _, sev := range order {
		for _, f := range findings {
			if strings.ToLower(f.Severity) != sev {
				continue
			}
			gate := "FAIL"
			if f.Evidence.Passed {
				gate = "PASS"
			}
			fmt.Fprintf(b, "### [%s] %s (`%s`)\n\n", strings.ToUpper(f.Severity), f.Title, f.ID)
			fmt.Fprintf(b, "- **Status:** %s | **Stage:** %s | **Source:** `%s` | **3-Gate:** %s\n",
				f.Status, f.Stage, f.Source, gate)
			if f.CWEID != "" {
				fmt.Fprintf(b, "- **CWE/Ref:** %s\n", f.CWEID)
			}
			if f.Endpoint != "" {
				fmt.Fprintf(b, "- **Endpoint:** `%s`\n", f.Endpoint)
			}
			if f.Description != "" {
				fmt.Fprintf(b, "\n%s\n", f.Description)
			}
			if f.Evidence.Passed || f.Evidence.Baseline != "" {
				b.WriteString("\n**Evidence (3-gate):**\n")
				if f.Evidence.Baseline != "" {
					fmt.Fprintf(b, "1. **Baseline:** %s\n", f.Evidence.Baseline)
				}
				if f.Evidence.Attack != "" {
					fmt.Fprintf(b, "2. **Attack:** %s\n", f.Evidence.Attack)
				}
				if f.Evidence.Diff != "" {
					fmt.Fprintf(b, "3. **Diff:** %s\n", f.Evidence.Diff)
				}
			}
			if f.StepsToReproduce != "" {
				b.WriteString("\n**Steps to reproduce:**\n")
				// Normalize literal \n from models into real line breaks / list items
				steps := strings.ReplaceAll(f.StepsToReproduce, `\n`, "\n")
				for _, ln := range strings.Split(steps, "\n") {
					ln = strings.TrimSpace(ln)
					if ln == "" {
						continue
					}
					fmt.Fprintf(b, "%s\n", ln)
				}
			}
			if f.Recommendation != "" {
				fmt.Fprintf(b, "\n**Recommendation:** %s\n", f.Recommendation)
			}
			if f.BusinessImpact != "" {
				fmt.Fprintf(b, "\n**Business impact:** %s\n", f.BusinessImpact)
			}
			b.WriteString("\n")
		}
	}
}

func writeMethodologySection(b *strings.Builder, stages []string, toolLog []ToolCallRecord) {
	b.WriteString("## Methodology\n\n")
	b.WriteString("Talon follows a fixed sequential pipeline (recon → exploit → post-exploit → optional codegen → judge):\n\n")
	pipeline := []struct {
		name string
		desc string
	}{
		{"recon", "Service/CVE verification (nmap HITL, nuclei, smbmap)"},
		{"exploit", "Metasploit module search + run_exploit with session poll"},
		{"post_exploit", "Session interaction for proof of compromise"},
		{"codegen", "LLM-written custom exploit in Docker sandbox (fallback)"},
		{"report", "Structured findings + judge verdict"},
	}
	covered := map[string]bool{}
	for _, s := range stages {
		covered[s] = true
	}
	// Also mark from tool log stages.
	for _, rec := range toolLog {
		covered[stageForTool(rec.ToolName)] = true
	}
	for _, p := range pipeline {
		mark := "☐"
		if covered[p.name] {
			mark = "☑"
		}
		fmt.Fprintf(b, "- %s **%s** — %s\n", mark, p.name, p.desc)
	}
	b.WriteString("\n")
}

func writeTimelineSection(b *strings.Builder, toolLog []ToolCallRecord) {
	b.WriteString("## Timeline\n\n")
	if len(toolLog) == 0 {
		b.WriteString("_No tool calls recorded._\n\n")
		return
	}
	b.WriteString("| # | Tool | Stage | Outcome |\n| --- | --- | --- | --- |\n")
	for _, rec := range toolLog {
		outcome := "ok"
		lower := strings.ToLower(rec.Output)
		if strings.Contains(lower, "error") || strings.Contains(lower, "failed") || strings.Contains(lower, "timeout") {
			outcome = "error/fail"
		} else if reSessionCreated.MatchString(rec.Output) {
			outcome = "session created"
		} else if extractProof(rec.Output) != "" {
			outcome = "proof captured"
		} else if len(rec.Output) == 0 {
			outcome = "empty"
		}
		fmt.Fprintf(b, "| %d | `%s` | %s | %s |\n", rec.Index, rec.ToolName, stageForTool(rec.ToolName), outcome)
	}
	b.WriteString("\n")
}

func writeValidationSection(b *strings.Builder, findings []Finding, judge *bool) {
	b.WriteString("## Validation (3-Gate Protocol)\n\n")
	b.WriteString("Findings that claim exploitation must satisfy CyberStrike-compatible gates:\n\n")
	b.WriteString("1. **Gate 1 — Baseline:** benign/pre-attack observation\n")
	b.WriteString("2. **Gate 2 — Attack:** exploit/modified observation\n")
	b.WriteString("3. **Gate 3 — Diff:** measurable difference between baseline and attack\n\n")

	passed, failed := 0, 0
	for _, f := range findings {
		if f.Evidence.Passed {
			passed++
		} else {
			failed++
		}
	}
	fmt.Fprintf(b, "| Metric | Count |\n| --- | --- |\n")
	fmt.Fprintf(b, "| 3-gate PASS | %d |\n", passed)
	fmt.Fprintf(b, "| 3-gate FAIL | %d |\n", failed)
	if judge != nil {
		j := "FALSE"
		if *judge {
			j = "TRUE"
		}
		fmt.Fprintf(b, "| Judge objective met | %s |\n", j)
	}
	b.WriteString("\n")
}

func writeAgentOutputSection(b *strings.Builder, finalMessage string) {
	b.WriteString("## Agent Output\n\n")
	if strings.TrimSpace(finalMessage) == "" {
		b.WriteString("_No final agent message._\n\n")
		return
	}
	b.WriteString("```\n")
	b.WriteString(finalMessage)
	if !strings.HasSuffix(finalMessage, "\n") {
		b.WriteString("\n")
	}
	b.WriteString("```\n\n")
	b.WriteString("---\n_Generated by Talon structured reporting (CyberStrike-inspired findings + 3-gate validation)._\n")
}

func coveredStages(toolLog []ToolCallRecord) []string {
	order := []string{"recon", "exploit", "post_exploit", "codegen", "report"}
	seen := map[string]bool{}
	for _, rec := range toolLog {
		seen[stageForTool(rec.ToolName)] = true
	}
	var out []string
	for _, s := range order {
		if seen[s] {
			out = append(out, s)
		}
	}
	return out
}
