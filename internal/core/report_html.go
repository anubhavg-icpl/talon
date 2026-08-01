package core

import (
	_ "embed"
	"encoding/base64"
	"fmt"
	"html"
	"strings"
	"time"
)

//go:embed talon-mark-red.webp
var talonMarkRedWebP []byte

// ReportHTMLInput gathers everything needed for a print/PDF-ready engagement report.
type ReportHTMLInput struct {
	RunID       string
	Input       RunInput
	Report      *StructuredReport
	Findings    []Finding
	ToolLog     []ToolCallRecord
	KillChain   *KillChainAnalysis
	Methodology *MethodologyState
	AgentMode   string
	Status      string
	StartedAt   time.Time
	EndedAt     time.Time
	FinalMsg    string
}

// RenderReportHTML builds a self-contained, print-friendly HTML report
// (open in browser → Print → Save as PDF). Uses structured fields first;
// falls back to report markdown narrative only when findings are empty.
func RenderReportHTML(in ReportHTMLInput) string {
	findings := in.Findings
	if in.Report != nil && len(in.Report.Findings) > 0 {
		findings = in.Report.Findings
	}
	summary := SummarizeFindings(findings)
	if in.Report != nil && in.Report.Summary.Total > 0 {
		summary = in.Report.Summary
	}

	target := in.Input.TargetIP
	if in.Report != nil && in.Report.Target != "" {
		target = in.Report.Target
	}
	cve := in.Input.CVEID
	if in.Report != nil && in.Report.CVEID != "" {
		cve = in.Report.CVEID
	}

	logo := ""
	if len(talonMarkRedWebP) > 0 {
		logo = "data:image/webp;base64," + base64.StdEncoding.EncodeToString(talonMarkRedWebP)
	}

	var b strings.Builder
	b.WriteString("<!DOCTYPE html>\n<html lang=\"en\">\n<head>\n")
	b.WriteString("<meta charset=\"utf-8\"/>\n")
	b.WriteString("<meta name=\"viewport\" content=\"width=device-width,initial-scale=1\"/>\n")
	fmt.Fprintf(&b, "<title>Talon Report · %s · %s</title>\n", html.EscapeString(shortID(in.RunID)), html.EscapeString(target))
	b.WriteString(reportHTMLCSS())
	b.WriteString("</head>\n<body>\n")

	// Toolbar (screen only)
	b.WriteString(`<div class="toolbar no-print">
<button type="button" onclick="window.print()">Print / Save PDF</button>
<span class="toolbar-hint">Use the browser print dialog → Save as PDF · paper size A4</span>
</div>
`)

	// Header
	b.WriteString(`<header class="report-header">`)
	if logo != "" {
		fmt.Fprintf(&b, `<img class="logo" src="%s" alt="Talon" width="56" height="56"/>`, logo)
	}
	b.WriteString(`<div class="header-text">
<p class="brand">TALON AI</p>
<h1>Engagement Validation Report</h1>
<p class="tagline">Authorized red-team engagement · CyberStrike-compatible findings · 3-gate evidence</p>
</div>
</header>
`)

	// Meta strip
	b.WriteString(`<section class="meta-strip">`)
	metaRow(&b, "Run", shortID(in.RunID))
	metaRow(&b, "Target", target)
	if cve != "" {
		metaRow(&b, "CVE", cve)
	}
	if in.Input.ServiceName != "" {
		metaRow(&b, "Service", in.Input.ServiceName)
	}
	if in.AgentMode != "" {
		metaRow(&b, "Agent", strings.ToUpper(in.AgentMode))
	}
	if in.Status != "" {
		metaRow(&b, "Status", in.Status)
	}
	if !in.StartedAt.IsZero() {
		metaRow(&b, "Started", in.StartedAt.UTC().Format(time.RFC3339))
	}
	if !in.EndedAt.IsZero() {
		metaRow(&b, "Ended", in.EndedAt.UTC().Format(time.RFC3339))
	}
	genAt := time.Now().UTC()
	if in.Report != nil && !in.Report.GeneratedAt.IsZero() {
		genAt = in.Report.GeneratedAt
	}
	metaRow(&b, "Generated", genAt.UTC().Format(time.RFC3339))
	b.WriteString(`</section>`)

	// Executive summary cards
	b.WriteString(`<section class="section">
<h2>Executive summary</h2>
<div class="cards">`)
	fmt.Fprintf(&b, `<div class="card"><p class="card-label">Findings</p><p class="card-value">%d</p>
<p class="card-sub">%d crit · %d high · %d med · %d low · %d info</p></div>`,
		summary.Total, summary.Critical, summary.High, summary.Medium, summary.Low, summary.Info)

	gatePass, gateFail := 0, 0
	for _, f := range findings {
		if f.Evidence.Passed {
			gatePass++
		} else if f.Evidence.Baseline != "" || f.Evidence.Attack != "" || f.Evidence.Diff != "" {
			gateFail++
		}
	}
	fmt.Fprintf(&b, `<div class="card"><p class="card-label">3-Gate</p><p class="card-value"><span class="ok">%d</span> / <span class="warn">%d</span></p>
<p class="card-sub">pass / fail</p></div>`, gatePass, gateFail)

	judgeLabel := "Not evaluated"
	judgeClass := "muted"
	if in.Report != nil && in.Report.JudgeVerdict != nil {
		if *in.Report.JudgeVerdict {
			judgeLabel = "COMPROMISED"
			judgeClass = "bad"
		} else {
			judgeLabel = "NOT COMPROMISED"
			judgeClass = "ok"
		}
	}
	fmt.Fprintf(&b, `<div class="card"><p class="card-label">Judge</p><p class="card-value %s">%s</p>
<p class="card-sub">objective verification</p></div>`, judgeClass, html.EscapeString(judgeLabel))

	stages := []string{}
	if in.Report != nil {
		stages = in.Report.StagesCovered
	}
	if len(stages) == 0 {
		stages = coveredStages(in.ToolLog)
	}
	stageStr := "—"
	if len(stages) > 0 {
		stageStr = strings.Join(stages, " → ")
	}
	fmt.Fprintf(&b, `<div class="card"><p class="card-label">Stages</p><p class="card-value stages">%s</p>
<p class="card-sub">pipeline coverage</p></div>`, html.EscapeString(stageStr))
	b.WriteString(`</div>`)

	// Narrative
	b.WriteString(`<div class="narrative">`)
	switch {
	case in.Report != nil && in.Report.JudgeVerdict != nil && *in.Report.JudgeVerdict:
		fmt.Fprintf(&b, `<p>The validation pipeline <strong class="bad">confirmed compromise</strong> of <code>%s</code>. %d finding(s) recorded; %d passed 3-gate evidence.</p>`,
			html.EscapeString(target), summary.Total, summary.Confirmed)
	case summary.Critical > 0 || summary.High > 0:
		fmt.Fprintf(&b, `<p>Elevated findings were observed on <code>%s</code> without a confirmed judge compromise verdict. Review high/critical items and remediate.</p>`,
			html.EscapeString(target))
	default:
		fmt.Fprintf(&b, `<p>No confirmed compromise of <code>%s</code> in this run. Stages executed: %s.</p>`,
			html.EscapeString(target), html.EscapeString(stageStr))
	}
	b.WriteString(`</div></section>`)

	// Methodology
	b.WriteString(`<section class="section"><h2>Methodology pipeline</h2><div class="pipeline">`)
	pipeline := []struct{ name, label string }{
		{"recon", "Recon"},
		{"exploit", "Exploit"},
		{"post_exploit", "Post-exploit"},
		{"codegen", "Codegen"},
		{"report", "Report"},
	}
	covered := map[string]bool{}
	for _, s := range stages {
		covered[s] = true
	}
	for _, rec := range in.ToolLog {
		covered[stageForTool(rec.ToolName)] = true
	}
	if in.Methodology != nil {
		for _, it := range in.Methodology.Items {
			if it.Covered {
				covered[strings.ToLower(strings.ReplaceAll(it.Stage, "-", "_"))] = true
			}
		}
	}
	for i, p := range pipeline {
		if i > 0 {
			b.WriteString(`<span class="pipe-arrow">→</span>`)
		}
		cls := "pipe-off"
		mark := "○"
		if covered[p.name] {
			cls = "pipe-on"
			mark = "●"
		}
		fmt.Fprintf(&b, `<span class="pipe %s">%s %s</span>`, cls, mark, html.EscapeString(p.label))
	}
	b.WriteString(`</div>`)
	if in.Methodology != nil && in.Methodology.Percent > 0 {
		fmt.Fprintf(&b, `<p class="muted small">Coverage %d%% · %d/%d stages</p>`,
			in.Methodology.Percent, in.Methodology.CoveredCount, in.Methodology.TotalCount)
	}
	b.WriteString(`</section>`)

	// Findings
	b.WriteString(`<section class="section"><h2>Findings</h2>`)
	if len(findings) == 0 {
		b.WriteString(`<p class="muted">No structured findings in this engagement.</p>`)
	} else {
		fmt.Fprintf(&b, `<p class="muted small">%d finding(s) · %d three-gate confirmed</p>`, summary.Total, summary.Confirmed)
		order := []string{SeverityCritical, SeverityHigh, SeverityMedium, SeverityLow, SeverityInfo}
		for _, sev := range order {
			for _, f := range findings {
				if strings.ToLower(f.Severity) != sev {
					continue
				}
				writeFindingHTML(&b, f)
			}
		}
	}
	b.WriteString(`</section>`)

	// Timeline
	if len(in.ToolLog) > 0 {
		b.WriteString(`<section class="section"><h2>Tool timeline</h2>
<table class="timeline"><thead><tr><th>#</th><th>Tool</th><th>Stage</th><th>Outcome</th></tr></thead><tbody>`)
		for _, rec := range in.ToolLog {
			outcome, ocls := toolOutcome(rec.Output)
			fmt.Fprintf(&b, `<tr><td class="muted">%d</td><td><code>%s</code></td><td>%s</td><td class="%s">%s</td></tr>`,
				rec.Index, html.EscapeString(rec.ToolName), html.EscapeString(stageForTool(rec.ToolName)), ocls, html.EscapeString(outcome))
		}
		b.WriteString(`</tbody></table></section>`)
	}

	// Kill chain
	if in.KillChain != nil {
		b.WriteString(`<section class="section"><h2>Kill chain analysis</h2>`)
		if in.KillChain.MaxSev != "" {
			fmt.Fprintf(&b, `<p class="muted small">Max severity: <span class="sev sev-%s">%s</span></p>`,
				html.EscapeString(strings.ToLower(in.KillChain.MaxSev)), html.EscapeString(in.KillChain.MaxSev))
		}
		if len(in.KillChain.Chains) == 0 {
			b.WriteString(`<p class="muted">No multi-stage chains detected.</p>`)
		} else {
			b.WriteString(`<ul class="chains">`)
			for _, c := range in.KillChain.Chains {
				fmt.Fprintf(&b, `<li><span class="sev sev-%s">%s</span> <code>%s → %s</code>`,
					html.EscapeString(strings.ToLower(c.Severity)), html.EscapeString(c.Severity),
					html.EscapeString(c.From), html.EscapeString(c.To))
				if c.Reason != "" {
					fmt.Fprintf(&b, `<p class="muted small">%s</p>`, html.EscapeString(c.Reason))
				}
				b.WriteString(`</li>`)
			}
			b.WriteString(`</ul>`)
		}
		if len(in.KillChain.NextSteps) > 0 {
			fmt.Fprintf(&b, `<p class="muted small"><strong>Suggested next:</strong> %s</p>`,
				html.EscapeString(strings.Join(in.KillChain.NextSteps, ", ")))
		}
		b.WriteString(`</section>`)
	}

	// Validation
	b.WriteString(`<section class="section"><h2>Validation (3-gate protocol)</h2>
<ol class="gates">
<li><strong>Baseline</strong> — benign / pre-attack observation</li>
<li><strong>Attack</strong> — exploit or modified observation</li>
<li><strong>Diff</strong> — measurable difference between baseline and attack</li>
</ol>
<table class="metrics"><tbody>`)
	fmt.Fprintf(&b, `<tr><td>3-gate PASS</td><td class="ok">%d</td></tr>`, gatePass)
	fmt.Fprintf(&b, `<tr><td>3-gate FAIL</td><td class="warn">%d</td></tr>`, gateFail)
	fmt.Fprintf(&b, `<tr><td>Judge</td><td class="%s">%s</td></tr>`, judgeClass, html.EscapeString(judgeLabel))
	b.WriteString(`</tbody></table></section>`)

	// Agent output (avoid dumping the full markdown report twice)
	final := agentNote(in.FinalMsg, in.Report)
	if final != "" {
		b.WriteString(`<section class="section"><h2>Agent output</h2><pre class="agent-out">`)
		b.WriteString(html.EscapeString(final))
		b.WriteString(`</pre></section>`)
	}

	// Footer
	b.WriteString(`<footer class="report-footer">
<p>Generated by <strong>Talon AI</strong> · structured reporting · CyberStrike-inspired findings + 3-gate validation</p>
<p class="muted small">AUTHORIZED TARGETS ONLY · CLASSIFIED AS OPERATOR ENGAGEMENT ARTIFACT</p>
</footer>
</body></html>`)

	return b.String()
}

func writeFindingHTML(b *strings.Builder, f Finding) {
	sev := strings.ToLower(f.Severity)
	gate := "none"
	gateLabel := "—"
	if f.Evidence.Passed {
		gate = "pass"
		gateLabel = "3-GATE PASS"
	} else if f.Evidence.Baseline != "" || f.Evidence.Attack != "" || f.Evidence.Diff != "" {
		gate = "fail"
		gateLabel = "3-GATE FAIL"
	}

	fmt.Fprintf(b, `<article class="finding sev-%s">
<div class="finding-head">
<span class="sev sev-%s">%s</span>
<span class="gate gate-%s">%s</span>
<span class="muted small">%s</span>
<span class="muted small stage">%s</span>
<span class="id">%s</span>
</div>
<h3>%s</h3>
`, sev, sev, html.EscapeString(strings.ToUpper(f.Severity)), gate, gateLabel,
		html.EscapeString(f.Status), html.EscapeString(f.Stage), html.EscapeString(f.ID),
		html.EscapeString(f.Title))

	if f.Endpoint != "" {
		fmt.Fprintf(b, `<p class="endpoint"><code>%s</code></p>`, html.EscapeString(f.Endpoint))
	}
	if f.Description != "" {
		fmt.Fprintf(b, `<p class="desc">%s</p>`, html.EscapeString(f.Description))
	}
	if f.Source != "" || f.CWEID != "" {
		fmt.Fprintf(b, `<p class="muted small">`)
		if f.Source != "" {
			fmt.Fprintf(b, `Source: <code>%s</code> `, html.EscapeString(f.Source))
		}
		if f.CWEID != "" {
			fmt.Fprintf(b, `· Ref: %s`, html.EscapeString(f.CWEID))
		}
		b.WriteString(`</p>`)
	}

	if f.Evidence.Baseline != "" || f.Evidence.Attack != "" || f.Evidence.Diff != "" {
		b.WriteString(`<div class="evidence"><p class="ev-title">3-gate evidence</p><ol>`)
		if f.Evidence.Baseline != "" {
			fmt.Fprintf(b, `<li><strong>Baseline</strong> — %s</li>`, html.EscapeString(f.Evidence.Baseline))
		} else {
			b.WriteString(`<li class="missing"><strong>Baseline</strong> — <em>not provided</em></li>`)
		}
		if f.Evidence.Attack != "" {
			fmt.Fprintf(b, `<li><strong>Attack</strong> — %s</li>`, html.EscapeString(truncate(f.Evidence.Attack, 1200)))
		} else {
			b.WriteString(`<li class="missing"><strong>Attack</strong> — <em>not provided</em></li>`)
		}
		if f.Evidence.Diff != "" {
			fmt.Fprintf(b, `<li><strong>Diff</strong> — %s</li>`, html.EscapeString(f.Evidence.Diff))
		} else {
			b.WriteString(`<li class="missing"><strong>Diff</strong> — <em>not provided</em></li>`)
		}
		b.WriteString(`</ol></div>`)
	}

	if steps := formatSteps(f.StepsToReproduce); steps != "" {
		fmt.Fprintf(b, `<div class="steps"><p class="ev-title">Steps to reproduce</p>%s</div>`, steps)
	}
	if f.Recommendation != "" {
		fmt.Fprintf(b, `<p class="remedy"><strong>Recommendation:</strong> %s</p>`, html.EscapeString(f.Recommendation))
	}
	if f.BusinessImpact != "" {
		fmt.Fprintf(b, `<p class="muted small"><strong>Business impact:</strong> %s</p>`, html.EscapeString(f.BusinessImpact))
	}
	b.WriteString(`</article>`)
}

func formatSteps(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	// Normalize literal \n sequences from LLM output
	raw = strings.ReplaceAll(raw, `\n`, "\n")
	lines := strings.Split(raw, "\n")
	var items []string
	for _, ln := range lines {
		ln = strings.TrimSpace(ln)
		if ln == "" {
			continue
		}
		// strip leading "1. " numbering if present — re-number in <ol>
		ln = strings.TrimLeft(ln, "0123456789. )-")
		ln = strings.TrimSpace(ln)
		if ln != "" {
			items = append(items, ln)
		}
	}
	if len(items) == 0 {
		return "<pre>" + html.EscapeString(raw) + "</pre>"
	}
	var b strings.Builder
	b.WriteString("<ol>")
	for _, it := range items {
		fmt.Fprintf(&b, "<li>%s</li>", html.EscapeString(it))
	}
	b.WriteString("</ol>")
	return b.String()
}

func toolOutcome(out string) (string, string) {
	lower := strings.ToLower(out)
	switch {
	case strings.Contains(lower, "error") || strings.Contains(lower, "failed") || strings.Contains(lower, "timeout") || strings.Contains(lower, `"success":false`):
		return "error/fail", "bad"
	case reSessionCreated.MatchString(out):
		return "session created", "ok"
	case extractProof(out) != "":
		return "proof captured", "ok"
	case len(out) == 0:
		return "empty", "muted"
	default:
		return "ok", "ok"
	}
}

func metaRow(b *strings.Builder, k, v string) {
	fmt.Fprintf(b, `<div class="meta-item"><span class="meta-k">%s</span><span class="meta-v">%s</span></div>`,
		html.EscapeString(k), html.EscapeString(v))
}

func shortID(id string) string {
	if len(id) > 8 {
		return id[:8]
	}
	return id
}

// agentNote returns a short operator-facing note, not the full structured markdown.
func agentNote(output string, report *StructuredReport) string {
	out := strings.TrimSpace(output)
	if out == "" {
		return ""
	}
	// Full report accidentally stored as Output — pull Agent Output section only.
	if strings.HasPrefix(out, "# Talon") || strings.Contains(out, "## Executive Summary") {
		if i := strings.Index(out, "## Agent Output"); i >= 0 {
			sec := out[i+len("## Agent Output"):]
			if j := strings.Index(sec, "\n## "); j >= 0 {
				sec = sec[:j]
			}
			if k := strings.Index(sec, "\n---"); k >= 0 {
				sec = sec[:k]
			}
			sec = strings.TrimSpace(sec)
			sec = strings.TrimPrefix(sec, "```")
			sec = strings.TrimSuffix(sec, "```")
			return strings.TrimSpace(sec)
		}
		return ""
	}
	// Prefer non-markdown short messages
	if report != nil && len(out) > 4000 {
		return truncate(out, 2000)
	}
	return out
}

func reportHTMLCSS() string {
	return `<style>
:root{
  --bg:#0a0608;--fg:#f2e8ea;--muted:#a8989c;--card:#120a0e;--border:#3a1a22;
  --primary:#e11d48;--ok:#4ade80;--warn:#fbbf24;--bad:#f43f5e;
  --crit:#ef4444;--high:#f97316;--med:#eab308;--low:#fca5a5;--info:#94a3b8;
  --font:ui-monospace,SFMono-Regular,Menlo,Monaco,Consolas,"Liberation Mono",monospace;
}
*{box-sizing:border-box}
body{margin:0;background:var(--bg);color:var(--fg);font-family:var(--font);font-size:13px;line-height:1.55;padding:0 0 3rem}
.toolbar{position:sticky;top:0;z-index:20;display:flex;align-items:center;gap:1rem;padding:.75rem 1.5rem;background:rgba(10,6,8,.92);border-bottom:1px solid var(--border);backdrop-filter:blur(8px)}
.toolbar button{background:var(--primary);color:#fff;border:0;border-radius:4px;padding:.5rem 1rem;font:inherit;font-size:11px;letter-spacing:.12em;text-transform:uppercase;cursor:pointer}
.toolbar button:hover{filter:brightness(1.1)}
.toolbar-hint{color:var(--muted);font-size:11px}
.report-header{display:flex;gap:1.25rem;align-items:center;padding:2rem 1.75rem 1.25rem;border-bottom:1px solid var(--border);background:linear-gradient(180deg,rgba(225,29,72,.08),transparent)}
.logo{width:56px;height:56px;object-fit:contain;border-radius:8px;border:1px solid var(--border);background:#000}
.brand{margin:0;color:var(--primary);font-size:11px;letter-spacing:.28em;font-weight:700}
.report-header h1{margin:.2rem 0;font-size:1.45rem;letter-spacing:.04em;font-weight:700}
.tagline{margin:0;color:var(--muted);font-size:11px}
.meta-strip{display:grid;grid-template-columns:repeat(auto-fill,minmax(160px,1fr));gap:.5rem;padding:1rem 1.75rem;border-bottom:1px solid var(--border);background:var(--card)}
.meta-item{display:flex;flex-direction:column;gap:.15rem}
.meta-k{font-size:9px;letter-spacing:.16em;text-transform:uppercase;color:var(--muted)}
.meta-v{font-size:12px;word-break:break-all}
.section{padding:1.5rem 1.75rem;border-bottom:1px solid var(--border)}
.section h2{margin:0 0 1rem;font-size:12px;letter-spacing:.2em;text-transform:uppercase;color:var(--primary)}
.cards{display:grid;grid-template-columns:repeat(auto-fit,minmax(160px,1fr));gap:.75rem;margin-bottom:1rem}
.card{background:var(--card);border:1px solid var(--border);border-radius:6px;padding:1rem}
.card-label{margin:0;font-size:9px;letter-spacing:.16em;text-transform:uppercase;color:var(--muted)}
.card-value{margin:.35rem 0 0;font-size:1.35rem;font-weight:700}
.card-value.stages{font-size:12px;font-weight:600;line-height:1.4}
.card-sub{margin:.25rem 0 0;font-size:10px;color:var(--muted)}
.narrative p{margin:0;color:var(--fg)}
.pipeline{display:flex;flex-wrap:wrap;align-items:center;gap:.5rem}
.pipe{padding:.4rem .7rem;border-radius:4px;border:1px solid var(--border);font-size:11px;letter-spacing:.08em;text-transform:uppercase}
.pipe-on{border-color:rgba(225,29,72,.45);background:rgba(225,29,72,.12);color:var(--primary)}
.pipe-off{color:var(--muted)}
.pipe-arrow{color:var(--muted)}
.finding{background:var(--card);border:1px solid var(--border);border-left:3px solid var(--info);border-radius:6px;padding:1rem 1.1rem;margin:0 0 .85rem}
.finding.sev-critical{border-left-color:var(--crit)}
.finding.sev-high{border-left-color:var(--high)}
.finding.sev-medium{border-left-color:var(--med)}
.finding.sev-low{border-left-color:var(--low)}
.finding.sev-info{border-left-color:var(--info)}
.finding-head{display:flex;flex-wrap:wrap;gap:.45rem;align-items:center;margin-bottom:.5rem}
.finding h3{margin:0 0 .4rem;font-size:14px;font-weight:650;letter-spacing:.02em}
.endpoint{margin:.2rem 0 .5rem}
.endpoint code,.finding code,code{font-family:inherit;color:var(--primary);font-size:11px;word-break:break-all}
.desc{margin:.4rem 0;color:var(--fg);opacity:.92}
.sev{display:inline-block;padding:.15rem .45rem;border-radius:3px;font-size:9px;letter-spacing:.12em;font-weight:700;border:1px solid}
.sev-critical{color:var(--crit);border-color:rgba(239,68,68,.4);background:rgba(239,68,68,.12)}
.sev-high{color:var(--high);border-color:rgba(249,115,22,.4);background:rgba(249,115,22,.12)}
.sev-medium{color:var(--med);border-color:rgba(234,179,8,.4);background:rgba(234,179,8,.12)}
.sev-low{color:var(--low);border-color:rgba(252,165,165,.35);background:rgba(252,165,165,.1)}
.sev-info{color:var(--info);border-color:rgba(148,163,184,.35);background:rgba(148,163,184,.1)}
.gate{font-size:9px;letter-spacing:.1em;padding:.15rem .4rem;border-radius:3px;border:1px solid}
.gate-pass{color:var(--ok);border-color:rgba(74,222,128,.4)}
.gate-fail{color:var(--warn);border-color:rgba(251,191,36,.4)}
.gate-none{color:var(--muted);border-color:var(--border)}
.id{margin-left:auto;font-size:10px;color:var(--muted)}
.stage{text-transform:uppercase;letter-spacing:.1em}
.evidence,.steps{margin:.75rem 0;padding:.75rem;background:rgba(0,0,0,.35);border:1px solid var(--border);border-radius:4px}
.ev-title{margin:0 0 .5rem;font-size:9px;letter-spacing:.16em;text-transform:uppercase;color:var(--primary)}
.evidence ol,.steps ol{margin:0;padding-left:1.2rem}
.evidence li,.steps li{margin:.35rem 0}
.evidence li.missing{color:var(--muted)}
.remedy{margin:.6rem 0 0}
.timeline{width:100%;border-collapse:collapse;font-size:11px}
.timeline th{text-align:left;color:var(--muted);font-weight:600;letter-spacing:.1em;text-transform:uppercase;font-size:9px;padding:.5rem;border-bottom:1px solid var(--border)}
.timeline td{padding:.4rem .5rem;border-bottom:1px solid rgba(58,26,34,.6);vertical-align:top}
.metrics{border-collapse:collapse;font-size:12px}
.metrics td{padding:.4rem .75rem .4rem 0;border-bottom:1px solid var(--border)}
.gates{margin:0 0 1rem;padding-left:1.2rem}
.chains{margin:0;padding-left:1.1rem}
.agent-out{margin:0;padding:1rem;background:#000;border:1px solid var(--border);border-radius:4px;white-space:pre-wrap;word-break:break-word;font-size:11px;color:#e8dce0;max-height:20rem;overflow:auto}
.report-footer{padding:1.5rem 1.75rem;text-align:center;color:var(--muted);font-size:11px;border-top:1px solid var(--border)}
.report-footer p{margin:.35rem 0}
.muted{color:var(--muted)}
.small{font-size:11px}
.ok{color:var(--ok)}
.warn{color:var(--warn)}
.bad{color:var(--bad)}
@media print{
  body{background:#fff;color:#111;font-size:10.5pt;padding:0}
  .no-print{display:none!important}
  .toolbar{display:none!important}
  .report-header{background:none;border-color:#ddd;padding:0 0 12pt}
  .brand,.section h2,.ev-title,.endpoint code,code{color:#b91c1c}
  .meta-strip,.card,.finding,.evidence,.steps,.agent-out{background:#fafafa;border-color:#ddd;color:#111}
  .meta-k,.muted,.card-sub,.card-label,.pipe-off,.timeline th{color:#555}
  .pipe-on{color:#b91c1c;border-color:#fca5a5;background:#fef2f2}
  .finding{break-inside:avoid}
  .agent-out{max-height:none;color:#222;background:#f5f5f5}
  .report-footer{border-color:#ddd;color:#555}
  a{color:inherit;text-decoration:none}
}
@page{margin:14mm 12mm}
</style>
`
}
