// Package agent is the Go port of pentest_core/final.py: a LangGraph-style
// orchestrator delegating to recon/exploit/post_exploit/codegen/report
// subagents over MCP tools, with tool tracking, tool-call limits, context
// trimming, and a human-in-the-loop gate on nmap_scan.
package core

import (
	"context"
	"sync"

	"github.com/anubhavg-icpl/pentester2/internal/config"
	"github.com/anubhavg-icpl/pentester2/internal/llm"
	"github.com/anubhavg-icpl/pentester2/internal/mcpclient"
)

// RunInput mirrors the TargetRequest/TargetInput payloads from fast_api.py
// and final.py's CLI main().
type RunInput struct {
	TargetIP    string
	CVEID       string
	ServiceName string
	Description string
	Context     config.Context
}

// ToolCallRecord mirrors ToolCallRecord in final.py's ToolCallTrackerMiddleware.
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
	// JudgeVerdict reports whether the judge_output LLM call confirmed the
	// exploitation objective was met. Only meaningful when Interrupted is false.
	JudgeVerdict bool
	Interrupted  bool
	Interrupt    *PendingInterrupt
}

// PendingInterrupt describes a HITL-gated tool call awaiting a decision,
// mirroring the HumanInTheLoopMiddleware interrupt payload in final.py.
type PendingInterrupt struct {
	ToolName string
	Args     map[string]any
}

// Decision resolves a PendingInterrupt, mirroring ResumeRequest in fast_api.py.
type Decision struct {
	Type       string // "approve", "reject", "edit"
	EditedArgs map[string]any
}

// CodegenTool is the "custom_exploit" tool the codegen subagent calls when
// prebuilt Metasploit modules fail -- implemented by internal/codegen,
// injected here to avoid an import cycle (agent <- codegen <- llm, not
// agent -> codegen -> agent).
type CodegenTool interface {
	Name() string
	Description() string
	Call(ctx context.Context, query string) (string, error)
}

// Orchestrator runs one full pentest validation workflow against a live MCP
// tool set. It is stateless between runs except for the ToolCallTracker log,
// mirroring build_agent()'s returned LangGraph app in final.py.
type Orchestrator struct {
	model   llm.ChatModel
	tools   *mcpclient.Multi
	codegen CodegenTool
	judge   llm.ChatModel

	// mu/sessions hold state for interrupted runs. There is no LangGraph
	// checkpointer in this port, so the Orchestrator itself is the one piece
	// of state a paused run rides on between Run() and Resume() -- keyed by
	// the identical RunInput the caller passes to both (see run() in
	// orchestrator.go for how sessions are parked and resumed).
	mu       sync.Mutex
	sessions map[RunInput]*orchestratorSession
}

func New(model llm.ChatModel, judge llm.ChatModel, tools *mcpclient.Multi, codegen CodegenTool) *Orchestrator {
	return &Orchestrator{
		model:    model,
		judge:    judge,
		tools:    tools,
		codegen:  codegen,
		sessions: make(map[RunInput]*orchestratorSession),
	}
}

// Run executes the workflow to completion (or to its first pending
// interrupt). Callers resume an interrupted run by calling Run again with
// the same RunInput plus a non-nil Decision via ResumeWith -- see
// orchestrator.go for the stateful session wrapper used by cmd/api.
func (o *Orchestrator) Run(ctx context.Context, input RunInput) (RunResult, error) {
	return o.run(ctx, input, nil)
}

// Resume continues a previously interrupted run, feeding back the human
// decision for the pending tool call.
func (o *Orchestrator) Resume(ctx context.Context, input RunInput, decision Decision) (RunResult, error) {
	return o.run(ctx, input, &decision)
}
