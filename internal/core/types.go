// Package core implements the agent orchestrator: a tool-calling loop that
// delegates to recon/exploit/post_exploit/codegen/report subagents over MCP
// tools, with tool tracking, tool-call limits, context trimming, and a
// human-in-the-loop gate on nmap_scan.
package core

import (
	"context"
	"errors"
	"strings"
	"sync"

	"github.com/anubhavg-icpl/talon/internal/config"
	"github.com/anubhavg-icpl/talon/internal/llm"
	"github.com/anubhavg-icpl/talon/internal/mcpclient"
)

// RunInput describes the target and attacker context for one validation run.
type RunInput struct {
	// SessionID uniquely identifies this run for interrupt parking/resume.
	// Prefer setting this (HTTP run_id or AMQP task run_id) so concurrent
	// runs against identical targets do not collide in the session map.
	// When empty, a composite of the target fields is used as a fallback.
	SessionID   string
	TargetIP    string
	CVEID       string
	ServiceName string
	Description string
	// AgentMode selects specialist pipeline (full|recon|exploit|web|network|post).
	AgentMode string
	Context   config.Context
}

// ToolCallRecord is one logged tool invocation within a run.
type ToolCallRecord struct {
	Index    int
	ToolName string
	Args     map[string]any
	Output   string
}

// RunResult is what a completed (or interrupted) run produces.
type RunResult struct {
	FinalMessage string
	ToolLog      []ToolCallRecord
	// JudgeVerdict reports whether the judge model confirmed the
	// exploitation objective was met. Only meaningful when Interrupted is false.
	JudgeVerdict bool
	// JudgeSet is true when JudgeVerdict was populated (false means judge
	// was skipped, not "judge said false").
	JudgeSet bool
	// Findings are structured security findings extracted from the run
	// (CyberStrike-inspired). Empty while Interrupted.
	Findings []Finding
	// Report is the structured multi-section validation report. Nil while Interrupted.
	Report *StructuredReport
	// KillChain is derived attack-path analysis. Nil while Interrupted.
	KillChain *KillChainAnalysis
	// Methodology is stage coverage for the run.
	Methodology *MethodologyState
	Interrupted bool
	Interrupt   *PendingInterrupt
}

// PendingInterrupt describes a HITL-gated tool call awaiting a decision.
type PendingInterrupt struct {
	ToolName string
	Args     map[string]any
}

// Decision resolves a PendingInterrupt.
type Decision struct {
	// Type is one of "approve", "reject", "edit" (lowercase). Callers should
	// normalize with NormalizeDecision before Resume.
	Type       string
	EditedArgs map[string]any
}

// NormalizeDecision lowercases Type and validates it is approve|reject|edit.
// For "edit", EditedArgs must be non-nil.
func NormalizeDecision(d Decision) (Decision, error) {
	t := strings.ToLower(strings.TrimSpace(d.Type))
	switch t {
	case "approve", "reject":
		return Decision{Type: t, EditedArgs: d.EditedArgs}, nil
	case "edit":
		if d.EditedArgs == nil {
			return Decision{}, errors.New("agent: decision edit requires edited_args")
		}
		return Decision{Type: t, EditedArgs: d.EditedArgs}, nil
	default:
		return Decision{}, errors.New("agent: decision must be approve, reject, or edit")
	}
}

// CodegenTool is the "custom_exploit" tool the codegen subagent calls when
// prebuilt Metasploit modules fail -- implemented by internal/forge,
// injected here to avoid an import cycle (core <- forge <- llm, not
// core -> forge -> core).
type CodegenTool interface {
	Name() string
	Description() string
	Call(ctx context.Context, query string) (string, error)
}

// ProgressFunc is invoked when the tool log grows (after each subagent/tool
// batch) so the control plane can surface live progress without waiting for
// HITL or completion. Wired via context (WithProgress) so concurrent runs
// do not share a process-global hook.
type ProgressFunc func(toolLog []ToolCallRecord)

// FindingsProgressFunc is invoked when mid-run findings change (report_finding).
type FindingsProgressFunc func(findings []Finding)

type progressCtxKey struct{}
type findingsProgressCtxKey struct{}

// WithProgress attaches a progress callback to ctx for one Run/Resume.
func WithProgress(ctx context.Context, fn ProgressFunc) context.Context {
	if fn == nil {
		return ctx
	}
	return context.WithValue(ctx, progressCtxKey{}, fn)
}

// WithFindingsProgress attaches a mid-run findings callback.
func WithFindingsProgress(ctx context.Context, fn FindingsProgressFunc) context.Context {
	if fn == nil {
		return ctx
	}
	return context.WithValue(ctx, findingsProgressCtxKey{}, fn)
}

func progressFrom(ctx context.Context) ProgressFunc {
	fn, _ := ctx.Value(progressCtxKey{}).(ProgressFunc)
	return fn
}

func findingsProgressFrom(ctx context.Context) FindingsProgressFunc {
	fn, _ := ctx.Value(findingsProgressCtxKey{}).(FindingsProgressFunc)
	return fn
}

func reportProgress(ctx context.Context, tr *tracker) {
	fn := progressFrom(ctx)
	if fn == nil || tr == nil {
		return
	}
	fn(append([]ToolCallRecord(nil), tr.log...))
}

func reportFindingsProgress(ctx context.Context, bag *FindingBag) {
	fn := findingsProgressFrom(ctx)
	if fn == nil || bag == nil {
		return
	}
	fn(bag.Snapshot())
}

// Orchestrator runs one full pentest validation workflow against a live MCP
// tool set. It is stateless between runs except for its tool-call log and
// any parked interrupted-run sessions.
type Orchestrator struct {
	model   llm.ChatModel
	tools   *mcpclient.Multi
	codegen CodegenTool
	judge   llm.ChatModel

	// mu/sessions hold state for interrupted runs -- the Orchestrator is
	// the one place a paused run's state rides between Run() and Resume(),
	// keyed by sessionKey(input) (prefer SessionID; see parkSession).
	mu       sync.Mutex
	sessions map[string]*orchestratorSession
}

func New(model llm.ChatModel, judge llm.ChatModel, tools *mcpclient.Multi, codegen CodegenTool) *Orchestrator {
	return &Orchestrator{
		model:    model,
		judge:    judge,
		tools:    tools,
		codegen:  codegen,
		sessions: make(map[string]*orchestratorSession),
	}
}

// Run executes the workflow to completion (or to its first pending
// interrupt). Callers resume an interrupted run via Resume, passing the
// identical RunInput plus the human's Decision.
func (o *Orchestrator) Run(ctx context.Context, input RunInput) (RunResult, error) {
	return o.run(ctx, input, nil)
}

// Resume continues a previously interrupted run, feeding back the human
// decision for the pending tool call.
func (o *Orchestrator) Resume(ctx context.Context, input RunInput, decision Decision) (RunResult, error) {
	norm, err := NormalizeDecision(decision)
	if err != nil {
		return RunResult{}, err
	}
	return o.run(ctx, input, &norm)
}
