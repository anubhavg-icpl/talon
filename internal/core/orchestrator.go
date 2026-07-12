// Package core implements the agent orchestrator: a tool-calling loop that
// delegates to recon/exploit/post_exploit/codegen/report subagents over MCP
// tools, with tool tracking, tool-call limits, context trimming, and a
// human-in-the-loop gate on nmap_scan.
package core

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/anubhavg-icpl/talon/internal/llm"
)

// maxToolCalls caps the total number of tool calls across an entire run --
// both the orchestrator's own delegate_* calls and the real MCP/codegen
// calls they trigger -- the conservative reading of "cap total tool calls
// across the whole run".
const maxToolCalls = 30

// maxSubagentModelTurns caps nested subagent model turns (each turn may
// issue one or more tool calls). Recon/exploit are capped tighter so a
// chatty model cannot burn the whole run budget on retries.
const (
	maxSubagentModelTurnsDefault = 10
	maxReconModelTurns           = 4
	maxExploitModelTurns         = 6
	maxPostExploitModelTurns     = 4
	maxCodegenModelTurns         = 4
	// maxOrchestratorTurns caps top-level orchestrator model rounds
	// (delegate_* planning). After this, the run finishes with whatever
	// summary is available instead of hanging on more LLM calls.
	maxOrchestratorTurns = 8
	// llmTurnTimeout bounds a single Converse call inside the agent loop.
	llmTurnTimeout = 90 * time.Second
)

// contextTrimTrigger/contextTrimKeep bound the running conversation size:
// once it exceeds contextTrimTrigger characters, all but the last
// contextTrimKeep tool-result messages are dropped.
const (
	contextTrimTrigger = 100_000
	contextTrimKeep    = 3
)

var errBudgetExhausted = errors.New("agent: tool call budget exhausted")

// tracker is the run-scoped tool-call counter and log, replacing
// ToolCallTrackerMiddleware + ToolCallLimitMiddleware's combined state --
// there's no middleware stack in this port to hang them on, so both live
// here and get threaded explicitly through every loop.
type tracker struct {
	count int
	log   []ToolCallRecord
}

func (t *tracker) allow() error {
	if t.count >= maxToolCalls {
		return errBudgetExhausted
	}
	t.count++
	return nil
}

func (t *tracker) record(name string, args map[string]any, output string) {
	t.log = append(t.log, ToolCallRecord{Index: len(t.log), ToolName: name, Args: args, Output: output})
}

// orchestratorSession is the resumable state for one interrupted run. The
// Orchestrator itself (mu/sessions in types.go) is the only place state can
// ride between a Run() that returns Interrupted and the matching Resume()
// call -- keyed by the caller's RunInput, which the caller is required to
// pass back unchanged.
//
// ponytail: keying sessions by RunInput means two concurrent runs against an
// identical target (same IP/CVE/service/description/context) collide.
// Upgrade by adding an explicit session/thread ID field to RunInput if that
// ever matters (e.g. the HTTP API serving concurrent identical requests).
type orchestratorSession struct {
	orchestratorMessages []llm.Message
	resolvedDelegates    []llm.ToolResult
	remainingDelegates   []llm.ToolCall
	delegateCallID       string
	delegateName         string
	subagentMessages     []llm.Message
	subagentResolved     []llm.ToolResult
	subagentRemaining    []llm.ToolCall
	pendingCallID        string
	pendingToolName      string
	pendingArgs          map[string]any
	toolCallCount        int
	toolLog              []ToolCallRecord
}

// pausedDelegate is what runDelegateBatch returns when one of the delegate
// calls in a batch triggered a nested HITL interrupt.
type pausedDelegate struct {
	callID       string
	name         string
	remaining    []llm.ToolCall
	subInterrupt *subagentInterrupt
}

// delegateBatchResume resumes runDelegateBatch mid-batch: the delegate call
// that was paused (currentCallID/currentName/subResume) plus the calls after
// it in the same orchestrator turn that hadn't started yet.
type delegateBatchResume struct {
	resolvedSoFar []llm.ToolResult
	currentCallID string
	currentName   string
	subResume     *subagentResume
	remaining     []llm.ToolCall
}

// subagentSpec is the (model, prompt, tools, gate, executor) tuple for one
// named delegate target.
type subagentSpec struct {
	model        llm.ChatModel
	systemPrompt string
	tools        []llm.ToolSpec
	gate         func(name string) bool
	exec         func(tr *tracker) toolExecFunc
	// maxTurns limits nested model turns; 0 means maxSubagentModelTurnsDefault.
	maxTurns int
}

func (o *Orchestrator) subagentConfig(delegateName string) (subagentSpec, bool) {
	switch delegateName {
	case "delegate_recon":
		return subagentSpec{
			model:        o.model,
			systemPrompt: reconSystemPrompt,
			tools:        o.tools.Subset("nmap_scan", "smbmap_scan", "nuclei_scan"),
			gate:         func(name string) bool { return name == "nmap_scan" },
			exec:         func(tr *tracker) toolExecFunc { return mcpExec(o.tools, tr) },
			maxTurns:     maxReconModelTurns,
		}, true
	case "delegate_exploit":
		return subagentSpec{
			model:        o.model,
			systemPrompt: exploitSystemPrompt,
			tools: o.tools.Subset(
				"list_exploits", "list_payloads", "generate_payload", "run_exploit",
				"run_auxiliary_module", "run_post_module", "sqlmap_scan",
				"arp_scan_discovery", "hydra_attack", "rustscan_fast_scan",
				"responder_credential_harvest",
			),
			exec:     func(tr *tracker) toolExecFunc { return mcpExec(o.tools, tr) },
			maxTurns: maxExploitModelTurns,
		}, true
	case "delegate_post_exploit":
		return subagentSpec{
			model:        o.model,
			systemPrompt: postExploitSystemPrompt,
			tools:        o.tools.Subset("list_active_sessions", "terminate_session", "send_session_command"),
			exec:         func(tr *tracker) toolExecFunc { return mcpExec(o.tools, tr) },
			maxTurns:     maxPostExploitModelTurns,
		}, true
	case "delegate_codegen":
		return subagentSpec{
			model:        o.model,
			systemPrompt: codeGenSystemPrompt,
			maxTurns:     maxCodegenModelTurns,
			tools: []llm.ToolSpec{{
				Name:        o.codegen.Name(),
				Description: o.codegen.Description(),
				InputSchema: map[string]any{
					"type": "object",
					"properties": map[string]any{
						"query": map[string]any{
							"type":        "string",
							"description": "The prompt requesting specific exploit code or helper functions (e.g., 'give code to exploit remote code execution for IP...').",
						},
					},
					"required": []string{"query"},
				},
			}},
			exec: func(tr *tracker) toolExecFunc { return codegenExec(o.codegen, tr) },
		}, true
	case "delegate_report":
		return subagentSpec{
			model:        o.model,
			systemPrompt: reportSystemPrompt,
			exec: func(tr *tracker) toolExecFunc {
				return func(ctx context.Context, call llm.ToolCall) (string, bool) {
					return "agent: the report subagent has no tools available", true
				}
			},
		}, true
	}
	return subagentSpec{}, false
}

// delegateToolSpecs is the synthetic tool surface exposed to the
// orchestrator's own model -- one callable per subagent, taking a single
// free-form "instructions" string, rather than handing the orchestrator
// raw MCP tool access directly.
func delegateToolSpecs() []llm.ToolSpec {
	const desc = "Detailed task instructions for the subagent, including any target/context details it needs."
	schema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"instructions": map[string]any{
				"type":        "string",
				"description": desc,
			},
		},
		"required": []string{"instructions"},
	}
	return []llm.ToolSpec{
		{Name: "delegate_recon", Description: "Verifies if the target service is running on the IP.", InputSchema: schema},
		{Name: "delegate_exploit", Description: "Searches and executes exploits against the verified service.", InputSchema: schema},
		{Name: "delegate_post_exploit", Description: "After an exploit is deployed, uses meterpreter to interact with the session.", InputSchema: schema},
		{Name: "delegate_codegen", Description: "If by using tools you are not able to exploit, then use this agent.", InputSchema: schema},
		{Name: "delegate_report", Description: "Generates final validation report upon confirmed exploit success.", InputSchema: schema},
	}
}

// seedPrompt builds the initial user message describing the target and
// attacker context.
func seedPrompt(input RunInput) string {
	return fmt.Sprintf(
		"Target Info:\n"+
			"- IP: %s\n"+
			"- CVE ID: %s\n"+
			"- Service Name: %s\n"+
			"- Description: %s\n\n"+
			"Attacker Context:\n"+
			"- LHOST: %s\n"+
			"- LPORT: %d\n\n"+
			"Begin the validation workflow now.",
		input.TargetIP, input.CVEID, input.ServiceName, input.Description,
		input.Context.LHOST, input.Context.LPORT,
	)
}

func lastAssistantText(messages []llm.Message) string {
	for i := len(messages) - 1; i >= 0; i-- {
		if messages[i].Role == llm.RoleAssistant && messages[i].Text != "" {
			return messages[i].Text
		}
	}
	return "Workflow stopped: tool call limit reached before a final summary was produced."
}

// messageLen is a cheap length-based stand-in for token counting, used only
// to decide when trimContext should fire.
func messageLen(messages []llm.Message) int {
	n := 0
	for _, m := range messages {
		n += len(m.Text)
		for _, tc := range m.ToolCalls {
			n += len(tc.Name)
			for k, v := range tc.Args {
				n += len(k) + len(fmt.Sprint(v))
			}
		}
		for _, tr := range m.ToolResults {
			n += len(tr.Name) + len(tr.Content)
		}
	}
	return n
}

// trimContext keeps the running conversation from growing unbounded: once
// the transcript grows past contextTrimTrigger, every tool-result message
// is dropped except the most recent few, keeping all other messages
// intact. This is a simple length-based trim, not token-accurate.
func trimContext(messages []llm.Message) []llm.Message {
	if messageLen(messages) <= contextTrimTrigger {
		return messages
	}
	var toolIdx []int
	for i, m := range messages {
		if m.Role == llm.RoleTool {
			toolIdx = append(toolIdx, i)
		}
	}
	if len(toolIdx) <= contextTrimKeep {
		return messages
	}
	drop := make(map[int]bool, len(toolIdx)-contextTrimKeep)
	for _, i := range toolIdx[:len(toolIdx)-contextTrimKeep] {
		drop[i] = true
	}
	out := make([]llm.Message, 0, len(messages))
	for i, m := range messages {
		if drop[i] {
			continue
		}
		out = append(out, m)
	}
	return out
}

// run is the single entry point behind both Run and Resume (see types.go).
// A fresh call (resume == nil) seeds the conversation and drives the
// orchestrator loop; a resumed call looks up the parked orchestratorSession
// for this exact RunInput and continues from wherever it paused.
func (o *Orchestrator) run(ctx context.Context, input RunInput, resume *Decision) (RunResult, error) {
	if resume == nil {
		messages := []llm.Message{llm.UserMessage(seedPrompt(input))}
		return o.orchestrateLoop(ctx, input, messages, &tracker{})
	}
	return o.resumeRun(ctx, input, *resume)
}

func (o *Orchestrator) resumeRun(ctx context.Context, input RunInput, decision Decision) (RunResult, error) {
	o.mu.Lock()
	sess, ok := o.sessions[input]
	if ok {
		delete(o.sessions, input)
	}
	o.mu.Unlock()
	if !ok {
		return RunResult{}, errors.New("agent: no pending interrupt to resume for this RunInput")
	}

	tr := &tracker{count: sess.toolCallCount, log: sess.toolLog}
	resumeState := &delegateBatchResume{
		resolvedSoFar: sess.resolvedDelegates,
		currentCallID: sess.delegateCallID,
		currentName:   sess.delegateName,
		remaining:     sess.remainingDelegates,
		subResume: &subagentResume{
			messages:        sess.subagentMessages,
			resolvedResults: sess.subagentResolved,
			remainingCalls:  sess.subagentRemaining,
			gatedCall:       llm.ToolCall{ID: sess.pendingCallID, Name: sess.pendingToolName, Args: sess.pendingArgs},
			decision:        decision,
		},
	}

	resolved, paused, err := o.runDelegateBatch(ctx, nil, tr, resumeState)
	if err != nil {
		if errors.Is(err, errBudgetExhausted) {
			return RunResult{FinalMessage: lastAssistantText(sess.orchestratorMessages), ToolLog: tr.log}, nil
		}
		return RunResult{}, err
	}
	if paused != nil {
		o.parkSession(input, sess.orchestratorMessages, resolved, paused, tr)
		return RunResult{
			Interrupted: true,
			Interrupt:   &PendingInterrupt{ToolName: paused.subInterrupt.toolName, Args: paused.subInterrupt.args},
			ToolLog:     tr.log,
		}, nil
	}

	messages := append(append([]llm.Message{}, sess.orchestratorMessages...), llm.Message{Role: llm.RoleTool, ToolResults: resolved})
	messages = trimContext(messages)
	return o.orchestrateLoop(ctx, input, messages, tr)
}

// orchestrateLoop drives the orchestrator's own tool-calling loop: seed/
// resumed messages in, delegate_* tool calls out, until the model returns
// final text with no more tool calls (then the judge runs) or the run gets
// interrupted or exhausts its tool-call budget.
func (o *Orchestrator) orchestrateLoop(ctx context.Context, input RunInput, messages []llm.Message, tr *tracker) (RunResult, error) {
	specs := delegateToolSpecs()
	for turn := 0; turn < maxOrchestratorTurns; turn++ {
		log.Printf("talon-core: orchestrator turn %d/%d tools_so_far=%d target=%s",
			turn+1, maxOrchestratorTurns, tr.count, input.TargetIP)

		msg, err := converseWithTimeout(ctx, o.model, orchestratorSystemPrompt, messages, specs)
		if err != nil {
			// Soft-fail on LLM timeout/errors after some progress: return
			// what we have rather than leaving the run stuck "running".
			if tr.count > 0 && (errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled)) {
				log.Printf("talon-core: orchestrator LLM timeout after progress: %v", err)
				return RunResult{
					FinalMessage: lastAssistantText(messages) + "\n[orchestrator stopped: LLM timeout]",
					ToolLog:      tr.log,
				}, nil
			}
			return RunResult{}, err
		}
		messages = append(messages, msg)

		if len(msg.ToolCalls) == 0 {
			log.Printf("talon-core: orchestrator final text; judging (tools=%d)", tr.count)
			verdict, err := judgeOutput(ctx, o.judge, msg.Text)
			if err != nil {
				log.Printf("talon-core: judge failed (returning without verdict): %v", err)
				return RunResult{FinalMessage: msg.Text, ToolLog: tr.log}, nil
			}
			return RunResult{FinalMessage: msg.Text, ToolLog: tr.log, JudgeVerdict: verdict}, nil
		}

		names := make([]string, 0, len(msg.ToolCalls))
		for _, tc := range msg.ToolCalls {
			names = append(names, tc.Name)
		}
		log.Printf("talon-core: orchestrator delegates %v", names)

		resolved, paused, err := o.runDelegateBatch(ctx, msg.ToolCalls, tr, nil)
		reportProgress(ctx, tr)
		if err != nil {
			if errors.Is(err, errBudgetExhausted) {
				log.Printf("talon-core: tool budget exhausted (tools=%d)", tr.count)
				return RunResult{FinalMessage: lastAssistantText(messages), ToolLog: tr.log}, nil
			}
			if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
				log.Printf("talon-core: delegate batch timeout: %v", err)
				return RunResult{
					FinalMessage: lastAssistantText(messages) + "\n[orchestrator stopped: delegate timeout]",
					ToolLog:      tr.log,
				}, nil
			}
			return RunResult{}, err
		}
		if paused != nil {
			o.parkSession(input, messages, resolved, paused, tr)
			return RunResult{
				Interrupted: true,
				Interrupt:   &PendingInterrupt{ToolName: paused.subInterrupt.toolName, Args: paused.subInterrupt.args},
				ToolLog:     tr.log,
			}, nil
		}

		messages = append(messages, llm.Message{Role: llm.RoleTool, ToolResults: resolved})
		messages = trimContext(messages)
	}
	log.Printf("talon-core: orchestrator turn budget exhausted (tools=%d)", tr.count)
	return RunResult{
		FinalMessage: lastAssistantText(messages) + "\n[orchestrator stopped: turn budget reached]",
		ToolLog:      tr.log,
	}, nil
}

// converseWithTimeout wraps ChatModel.Converse with llmTurnTimeout so a hung
// provider cannot leave a run stuck in "running" forever.
func converseWithTimeout(ctx context.Context, model llm.ChatModel, system string, messages []llm.Message, tools []llm.ToolSpec) (llm.Message, error) {
	cctx, cancel := context.WithTimeout(ctx, llmTurnTimeout)
	defer cancel()
	return model.Converse(cctx, system, messages, tools)
}

// runDelegateBatch executes one orchestrator turn's delegate_* tool calls in
// order (or continues a batch previously paused mid-way through, via
// resume), stopping at the first one whose nested subagent loop hits a
// HITL-gated tool.
func (o *Orchestrator) runDelegateBatch(ctx context.Context, calls []llm.ToolCall, tr *tracker, resume *delegateBatchResume) ([]llm.ToolResult, *pausedDelegate, error) {
	var resolved []llm.ToolResult

	if resume != nil {
		resolved = append(resolved, resume.resolvedSoFar...)
		sub, ok := o.subagentConfig(resume.currentName)
		if !ok {
			return nil, nil, fmt.Errorf("agent: unknown delegate %q in resumed session", resume.currentName)
		}
		log.Printf("talon-core: resume subagent %s", resume.currentName)
		text, subI, err := runSubagent(ctx, sub.model, sub.systemPrompt, sub.tools, "", sub.gate, sub.exec(tr), resume.subResume, tr, sub.maxTurns)
		if err != nil {
			return nil, nil, err
		}
		reportProgress(ctx, tr)
		if subI != nil {
			return resolved, &pausedDelegate{callID: resume.currentCallID, name: resume.currentName, remaining: resume.remaining, subInterrupt: subI}, nil
		}
		resolved = append(resolved, llm.ToolResult{ToolCallID: resume.currentCallID, Name: resume.currentName, Content: text})
		calls = resume.remaining
	}

	for i, tc := range calls {
		sub, ok := o.subagentConfig(tc.Name)
		if !ok {
			resolved = append(resolved, llm.ToolResult{ToolCallID: tc.ID, Name: tc.Name, Content: fmt.Sprintf("agent: unknown delegate tool %q", tc.Name), IsError: true})
			continue
		}
		if err := tr.allow(); err != nil {
			return nil, nil, err
		}
		instructions, _ := tc.Args["instructions"].(string)
		log.Printf("talon-core: run subagent %s (tools_so_far=%d)", tc.Name, tr.count)
		text, subI, err := runSubagent(ctx, sub.model, sub.systemPrompt, sub.tools, instructions, sub.gate, sub.exec(tr), nil, tr, sub.maxTurns)
		if err != nil {
			return nil, nil, err
		}
		reportProgress(ctx, tr)
		if subI != nil {
			return resolved, &pausedDelegate{callID: tc.ID, name: tc.Name, remaining: append([]llm.ToolCall{}, calls[i+1:]...), subInterrupt: subI}, nil
		}
		resolved = append(resolved, llm.ToolResult{ToolCallID: tc.ID, Name: tc.Name, Content: text})
	}
	return resolved, nil, nil
}

func (o *Orchestrator) parkSession(input RunInput, messages []llm.Message, resolved []llm.ToolResult, paused *pausedDelegate, tr *tracker) {
	o.mu.Lock()
	defer o.mu.Unlock()
	if o.sessions == nil {
		o.sessions = make(map[RunInput]*orchestratorSession)
	}
	o.sessions[input] = &orchestratorSession{
		orchestratorMessages: messages,
		resolvedDelegates:    resolved,
		remainingDelegates:   paused.remaining,
		delegateCallID:       paused.callID,
		delegateName:         paused.name,
		subagentMessages:     paused.subInterrupt.messages,
		subagentResolved:     paused.subInterrupt.resolvedResults,
		subagentRemaining:    paused.subInterrupt.remainingCalls,
		pendingCallID:        paused.subInterrupt.callID,
		pendingToolName:      paused.subInterrupt.toolName,
		pendingArgs:          paused.subInterrupt.args,
		toolCallCount:        tr.count,
		toolLog:              tr.log,
	}
}
