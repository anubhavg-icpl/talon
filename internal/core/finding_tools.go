package core

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/anubhavg-icpl/talon/internal/llm"
)

// FindingBag is a run-scoped bag of agent-reported findings (mid-run).
// Thread-safe; used by report_finding / triage_finding synthetic tools.
type FindingBag struct {
	mu       sync.Mutex
	findings []Finding
	seq      int
}

func NewFindingBag() *FindingBag { return &FindingBag{} }

func (b *FindingBag) Snapshot() []Finding {
	if b == nil {
		return nil
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	return append([]Finding(nil), b.findings...)
}

func (b *FindingBag) Add(f Finding) Finding {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.seq++
	if f.ID == "" {
		f.ID = fmt.Sprintf("FIND-%03d", b.seq)
	}
	if f.Status == "" {
		// Dedup-aware: similar title+endpoint → new (needs triage)
		similar := false
		for _, ex := range b.findings {
			if dedupeKey(ex) == dedupeKey(f) {
				similar = true
				break
			}
		}
		if similar {
			f.Status = FindingStatusNew
		} else {
			f.Status = FindingStatusApproved
		}
	}
	if f.CreatedAt.IsZero() {
		f.CreatedAt = time.Now().UTC()
	}
	// Run 3-gate validation on evidence
	f.Evidence, _ = ValidateThreeGate(f.Evidence)
	b.findings = append(b.findings, f)
	return f
}

func (b *FindingBag) Triage(id, status, duplicateOf string) (Finding, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	status = strings.ToLower(strings.TrimSpace(status))
	switch status {
	case FindingStatusApproved, FindingStatusDup, FindingStatusOpen, "fixed", "ignored", FindingStatusNew:
	default:
		return Finding{}, fmt.Errorf("invalid status %q", status)
	}
	for i := range b.findings {
		if b.findings[i].ID == id {
			b.findings[i].Status = status
			if status == FindingStatusDup && duplicateOf != "" {
				// encode duplicate link in description note
				if !strings.Contains(b.findings[i].Description, "duplicate_of=") {
					b.findings[i].Description += fmt.Sprintf(" [duplicate_of=%s]", duplicateOf)
				}
			}
			return b.findings[i], nil
		}
	}
	return Finding{}, fmt.Errorf("finding %s not found", id)
}

// MergeExtracted merges auto-extracted findings with agent-reported ones,
// preferring agent reports on same dedupe key, then appends unique extracts.
func MergeExtracted(reported, extracted []Finding) []Finding {
	seen := map[string]bool{}
	out := make([]Finding, 0, len(reported)+len(extracted))
	for _, f := range reported {
		k := dedupeKey(f)
		if seen[k] {
			continue
		}
		seen[k] = true
		out = append(out, f)
	}
	for _, f := range extracted {
		k := dedupeKey(f)
		if seen[k] {
			continue
		}
		seen[k] = true
		out = append(out, f)
	}
	return out
}

// reportFindingToolSpec is the CyberStrike-style report_vulnerability surface.
func reportFindingToolSpec() llm.ToolSpec {
	return llm.ToolSpec{
		Name:        "report_finding",
		Description: "Record a structured security finding with 3-gate evidence (baseline, attack, diff). Always record real findings — triage duplicates afterward.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"severity":     map[string]any{"type": "string", "enum": []any{"critical", "high", "medium", "low", "info"}},
				"title":        map[string]any{"type": "string"},
				"description":  map[string]any{"type": "string"},
				"cwe_id":       map[string]any{"type": "string"},
				"endpoint":     map[string]any{"type": "string"},
				"attack_vector": map[string]any{"type": "string"},
				"baseline":     map[string]any{"type": "string", "description": "Gate 1: pre-attack observation"},
				"attack":       map[string]any{"type": "string", "description": "Gate 2: exploit observation"},
				"diff":         map[string]any{"type": "string", "description": "Gate 3: measurable difference"},
				"recommendation": map[string]any{"type": "string"},
				"steps_to_reproduce": map[string]any{"type": "string"},
				"business_impact": map[string]any{"type": "string"},
				"poc":          map[string]any{"type": "string"},
				"stage":        map[string]any{"type": "string"},
			},
			"required": []any{"severity", "title", "description"},
		},
	}
}

func triageFindingToolSpec() llm.ToolSpec {
	return llm.ToolSpec{
		Name:        "triage_finding",
		Description: "Triage a previously reported finding: approved, duplicate, open, ignored, fixed.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"id":           map[string]any{"type": "string"},
				"status":       map[string]any{"type": "string", "enum": []any{"approved", "duplicate", "open", "ignored", "fixed", "new"}},
				"duplicate_of": map[string]any{"type": "string"},
			},
			"required": []any{"id", "status"},
		},
	}
}

func findingToolSpecs() []llm.ToolSpec {
	return []llm.ToolSpec{reportFindingToolSpec(), triageFindingToolSpec()}
}

// hybridExec routes local agent tools (findings + CyberStrike skills) then MCP/codegen.
func hybridExec(base toolExecFunc, bag *FindingBag, stage string, tr *tracker) toolExecFunc {
	return func(ctx context.Context, call llm.ToolCall) (string, bool) {
		switch call.Name {
		case "report_finding":
			out, err := handleReportFinding(bag, stage, call.Args, tr)
			reportFindingsProgress(ctx, bag)
			return out, err
		case "triage_finding":
			out, err := handleTriageFinding(bag, call.Args, tr)
			reportFindingsProgress(ctx, bag)
			return out, err
		case "skill_search":
			return handleSkillSearch(call.Args, tr)
		case "skill_get":
			return handleSkillGet(call.Args, tr)
		default:
			return base(ctx, call)
		}
	}
}

func handleReportFinding(bag *FindingBag, stage string, args map[string]any, tr *tracker) (string, bool) {
	if bag == nil {
		return "finding bag not configured", true
	}
	sev, _ := args["severity"].(string)
	title, _ := args["title"].(string)
	desc, _ := args["description"].(string)
	if strings.TrimSpace(title) == "" {
		return "title is required", true
	}
	if strings.TrimSpace(sev) == "" {
		sev = SeverityMedium
	}
	st, _ := args["stage"].(string)
	if st == "" {
		st = stage
	}
	f := Finding{
		Severity:         strings.ToLower(sev),
		Title:            title,
		Description:      desc,
		CWEID:            strArg(args, "cwe_id"),
		Endpoint:         strArg(args, "endpoint"),
		AttackVector:     strArg(args, "attack_vector"),
		Recommendation:   strArg(args, "recommendation"),
		StepsToReproduce: strArg(args, "steps_to_reproduce"),
		BusinessImpact:   strArg(args, "business_impact"),
		PoC:              strArg(args, "poc"),
		Stage:            st,
		Source:           "report_finding",
		Evidence: GateEvidence{
			Baseline: strArg(args, "baseline"),
			Attack:   strArg(args, "attack"),
			Diff:     strArg(args, "diff"),
		},
	}
	// If poc has structure but gates empty, use poc as attack
	if f.Evidence.Attack == "" && f.PoC != "" {
		f.Evidence.Attack = f.PoC
	}
	added := bag.Add(f)
	tr.record("report_finding", args, fmt.Sprintf("recorded %s [%s] 3-gate=%v status=%s", added.ID, added.Severity, added.Evidence.Passed, added.Status))
	msg := fmt.Sprintf("Recorded as %s (status: %s, 3-gate: %v).", added.ID, added.Status, added.Evidence.Passed)
	if !added.Evidence.Passed {
		msg += " WARNING: 3-gate incomplete — add baseline/attack/diff for confirmation."
	}
	if added.Status == FindingStatusNew {
		msg += " Similar finding may exist — call triage_finding to approve or mark duplicate."
	}
	return msg, false
}

func handleTriageFinding(bag *FindingBag, args map[string]any, tr *tracker) (string, bool) {
	if bag == nil {
		return "finding bag not configured", true
	}
	id, _ := args["id"].(string)
	status, _ := args["status"].(string)
	dup, _ := args["duplicate_of"].(string)
	f, err := bag.Triage(id, status, dup)
	if err != nil {
		tr.record("triage_finding", args, err.Error())
		return err.Error(), true
	}
	out := fmt.Sprintf("Triaged %s → status=%s", f.ID, f.Status)
	tr.record("triage_finding", args, out)
	return out, false
}

func strArg(args map[string]any, key string) string {
	v, _ := args[key].(string)
	return strings.TrimSpace(v)
}

// encodeFindingJSON helper for tests/debug
func encodeFindingJSON(f Finding) string {
	b, _ := json.Marshal(f)
	return string(b)
}
