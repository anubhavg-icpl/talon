package core

import "strings"

// CoverageItem is one methodology checklist cell for a run.
type CoverageItem struct {
	Stage   string `json:"stage"`
	Label   string `json:"label"`
	Covered bool   `json:"covered"`
	Tools   []string `json:"tools,omitempty"`
	Notes   string `json:"notes,omitempty"`
}

// MethodologyState is run-level methodology coverage (CyberStrike-inspired).
type MethodologyState struct {
	Items          []CoverageItem `json:"items"`
	CoveredCount   int            `json:"covered_count"`
	TotalCount     int            `json:"total_count"`
	Percent        int            `json:"percent"`
	AgentMode      string         `json:"agent_mode,omitempty"`
}

// ComputeMethodology derives coverage from the tool log + agent mode.
func ComputeMethodology(toolLog []ToolCallRecord, agentMode string) MethodologyState {
	toolsByStage := map[string][]string{}
	for _, rec := range toolLog {
		st := stageForTool(rec.ToolName)
		toolsByStage[st] = appendUnique(toolsByStage[st], rec.ToolName)
	}

	checklist := []struct {
		stage string
		label string
	}{
		{"recon", "Reconnaissance / service verification"},
		{"exploit", "Exploitation / module execution"},
		{"post_exploit", "Post-exploitation / proof of compromise"},
		{"codegen", "Custom exploit (codegen fallback)"},
		{"report", "Reporting / validation"},
	}

	// report is "covered" if we have any tools or will be covered at finalize
	items := make([]CoverageItem, 0, len(checklist))
	covered := 0
	for _, c := range checklist {
		tools := toolsByStage[c.stage]
		ok := len(tools) > 0
		// For modes that skip stages, mark N/A as covered for percentage fairness
		mode := NormalizeAgentMode(agentMode)
		skipped := false
		allowed := AllowedDelegates(mode)
		switch c.stage {
		case "recon":
			skipped = !allowed["delegate_recon"]
		case "exploit":
			skipped = !allowed["delegate_exploit"]
		case "post_exploit":
			skipped = !allowed["delegate_post_exploit"]
		case "codegen":
			skipped = !allowed["delegate_codegen"]
		case "report":
			// always expected
		}
		note := ""
		if skipped {
			note = "skipped by agent mode"
			ok = true // don't penalize
		}
		if ok {
			covered++
		}
		items = append(items, CoverageItem{
			Stage: c.stage, Label: c.label, Covered: ok, Tools: tools, Notes: note,
		})
	}
	total := len(items)
	pct := 0
	if total > 0 {
		pct = (covered * 100) / total
	}
	return MethodologyState{
		Items: items, CoveredCount: covered, TotalCount: total,
		Percent: pct, AgentMode: NormalizeAgentMode(agentMode),
	}
}

func appendUnique(ss []string, v string) []string {
	for _, s := range ss {
		if s == v {
			return ss
		}
	}
	return append(ss, v)
}

// FormatMethodologyMarkdown for embedding in reports.
func FormatMethodologyMarkdown(m MethodologyState) string {
	var b strings.Builder
	b.WriteString("## Methodology Coverage\n\n")
	b.WriteString("| Stage | Covered | Tools |\n| --- | --- | --- |\n")
	for _, it := range m.Items {
		mark := "☐"
		if it.Covered {
			mark = "☑"
		}
		tools := "—"
		if len(it.Tools) > 0 {
			tools = "`" + strings.Join(it.Tools, "`, `") + "`"
		}
		note := it.Label
		if it.Notes != "" {
			note += " (" + it.Notes + ")"
		}
		b.WriteString("| " + mark + " " + note + " | " + mark + " | " + tools + " |\n")
	}
	b.WriteString("\n**Coverage:** ")
	b.WriteString(itoa(m.CoveredCount) + "/" + itoa(m.TotalCount) + " (" + itoa(m.Percent) + "%)\n")
	return b.String()
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var digits []byte
	neg := n < 0
	if neg {
		n = -n
	}
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	if neg {
		digits = append([]byte{'-'}, digits...)
	}
	return string(digits)
}
