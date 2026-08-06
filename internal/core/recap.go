package core

import (
	"fmt"
	"strings"
	"time"
)

// RunRecap is a deterministic (LLM-free) run recap.
// Ported from pentest agent report/solve_report.py.
type RunRecap struct {
	Target      string         `json:"target"`
	RunID       string         `json:"run_id"`
	StartTime   time.Time      `json:"start_time"`
	EndTime     time.Time      `json:"end_time"`
	Duration    string         `json:"duration"`
	SolvePath   []RecapStep    `json:"solve_path"`
	KeyEvidence []RecapEvidence `json:"key_evidence"`
	Reproduction []RecapRepro  `json:"reproduction"`
	FindingCount int           `json:"finding_count"`
	VerifiedCount int          `json:"verified_count"`
}

type RecapStep struct {
	Step      int       `json:"step"`
	Action    string    `json:"action"`
	Result    string    `json:"result"`
	Timestamp time.Time `json:"timestamp"`
}

type RecapEvidence struct {
	ID       string `json:"id"`
	Tool     string `json:"tool"`
	Summary  string `json:"summary"`
	Relevant bool   `json:"relevant"`
}

type RecapRepro struct {
	Label   string `json:"label"`
	Command string `json:"command"`
}

// BuildRecap constructs a deterministic run recap from stored state.
// No LLM needed — pure data assembly from tracker, evidence, and findings.
func BuildRecap(target, runID string, startTime time.Time, tracker *ToolCallTracker, evidence *EvidenceStore, findings []Finding) RunRecap {
	recap := RunRecap{
		Target:    target,
		RunID:     runID,
		StartTime: startTime,
		EndTime:   time.Now().UTC(),
	}

	recap.Duration = recap.EndTime.Sub(startTime).Round(time.Second).String()

	// Build solve path from tracker
	if tracker != nil {
		for i, call := range tracker.Calls {
			recap.SolvePath = append(recap.SolvePath, RecapStep{
				Step:   i + 1,
				Action: fmt.Sprintf("%s(%s)", call.ToolName, truncateStr(formatArgs(call.Args), 100)),
				Result: truncateStr(call.Output, 200),
			})
		}
	}

	// Build key evidence from store
	if evidence != nil {
		all := evidence.List()
		limit := 20
		if len(all) < limit {
			limit = len(all)
		}
		for _, e := range all[:limit] {
			recap.KeyEvidence = append(recap.KeyEvidence, RecapEvidence{
				ID: e.ID, Tool: e.Tool, Summary: e.Summary,
			})
		}
	}

	// Count findings
	recap.FindingCount = len(findings)
	for _, f := range findings {
		if f.Evidence.Passed {
			recap.VerifiedCount++
		}
		// Build reproduction commands from PoC
		if f.PoC != "" {
			recap.Reproduction = append(recap.Reproduction, RecapRepro{
				Label:   f.Title,
				Command: f.PoC,
			})
		}
	}

	return recap
}

// FormatMarkdown renders the recap as a markdown report section.
func (r RunRecap) FormatMarkdown() string {
	var sb strings.Builder

	sb.WriteString("## Run Recap\n\n")
	sb.WriteString(fmt.Sprintf("**Target:** %s  \n", r.Target))
	sb.WriteString(fmt.Sprintf("**Duration:** %s  \n", r.Duration))
	sb.WriteString(fmt.Sprintf("**Findings:** %d total (%d verified by 3-gate)  \n\n", r.FindingCount, r.VerifiedCount))

	// Solve path timeline
	if len(r.SolvePath) > 0 {
		sb.WriteString("### Solve Path\n\n")
		limit := len(r.SolvePath)
		if limit > 30 {
			limit = 30
		}
		for _, step := range r.SolvePath[:limit] {
			sb.WriteString(fmt.Sprintf("%d. **%s** → %s\n", step.Step, step.Action, step.Result))
		}
		if len(r.SolvePath) > 30 {
			sb.WriteString(fmt.Sprintf("\n…(%d more steps omitted)\n", len(r.SolvePath)-30))
		}
		sb.WriteString("\n")
	}

	// Key evidence
	if len(r.KeyEvidence) > 0 {
		sb.WriteString("### Key Evidence\n\n")
		sb.WriteString("| ID | Tool | Summary |\n")
		sb.WriteString("|-----|------|---------|\n")
		for _, e := range r.KeyEvidence {
			summary := strings.ReplaceAll(e.Summary, "|", "\\|")
			if len(summary) > 80 {
				summary = summary[:80] + "..."
			}
			sb.WriteString(fmt.Sprintf("| %s | %s | %s |\n", e.ID, e.Tool, summary))
		}
		sb.WriteString("\n")
	}

	// Reproduction
	if len(r.Reproduction) > 0 {
		sb.WriteString("### Reproduction\n\n")
		for _, r := range r.Reproduction {
			sb.WriteString(fmt.Sprintf("**%s:**\n```\n%s\n```\n\n", r.Label, r.Command))
		}
	}

	return sb.String()
}

// ToolCallTracker is a lightweight tracker for recap purposes.
// In practice, the existing tracker's log is used.
type ToolCallTracker struct {
	Calls []ToolCallRecord
}

func formatArgs(args map[string]any) string {
	if len(args) == 0 {
		return ""
	}
	var parts []string
	for k, v := range args {
		parts = append(parts, fmt.Sprintf("%s=%v", k, v))
	}
	return strings.Join(parts, ", ")
}

// --- Headless presets ---

// RunPreset defines a configuration preset for headless/CI runs.
type RunPreset struct {
	Name        string
	Description string
	AgentMode   string
	MaxTurns    int
	MaxToolCalls int
}

// Presets returns the available run presets for headless CI.
func Presets() map[string]RunPreset {
	return map[string]RunPreset{
		"quick": {
			Name:        "quick",
			Description: "Fast scan — recon + basic vulnerability discovery only",
			AgentMode:   "recon",
			MaxTurns:    15,
			MaxToolCalls: 50,
		},
		"standard": {
			Name:        "standard",
			Description: "Standard run — recon + exploit + report",
			AgentMode:   "auto",
			MaxTurns:    30,
			MaxToolCalls: 150,
		},
		"deep": {
			Name:        "deep",
			Description: "Deep assessment — all stages with maximum tool budget",
			AgentMode:   "auto",
			MaxTurns:    60,
			MaxToolCalls: 400,
		},
	}
}

// ExitCode represents the CI exit code from a headless run.
const (
	ExitClean        = 0 // No confirmed findings
	ExitInfraError   = 1 // Infrastructure/run error
	ExitConfirmedVuln = 2 // Confirmed vulnerability (3-gate passed)
)
