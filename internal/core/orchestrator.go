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
	"strings"
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
	// toolExecTimeout bounds a single real tool call (MCP tool / codegen).
	// Legitimate tools (nmap/nuclei/sqlmap/run_exploit) run 30-90s, but an
	// LLM-crafted exploit (custom_exploit) or a hung upstream MCP server can
	// block forever and stall the whole run -- the failure that hung the DVWA
	// run. A bounded call is aborted at the deadline and surfaces a clear
	// message the agent can react to (report partial / try another path).
	toolExecTimeout = 3 * time.Minute
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
	// bag holds mid-run agent-reported findings (report_finding tool).
	bag *FindingBag
	// store holds all tool outputs as evidence for the agent to inspect.
	store *EvidenceStore
	// cases holds detection pipeline cases for SOC triage/investigation/tuning.
	cases *CaseStore
	// correction observes tool lifecycle and emits anti-hallucination hints.
	correction *CorrectionLayer
	// traffic records HTTP request/response pairs for replay and analysis.
	traffic *TrafficStore
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

// sessionKey returns the map key used to park/resume interrupted runs.
// Prefer SessionID (HTTP run_id / AMQP task run_id). Fall back to a composite
// of target fields so unit tests without an ID still work — that fallback can
// collide under concurrent identical targets, so production callers must set
// SessionID.
func sessionKey(input RunInput) string {
	if input.SessionID != "" {
		return input.SessionID
	}
	return fmt.Sprintf("%s|%s|%s|%s|%s|%d",
		input.TargetIP, input.CVEID, input.ServiceName, input.Description,
		input.Context.LHOST, input.Context.LPORT)
}

// orchestratorSession is the resumable state for one interrupted run. The
// Orchestrator itself (mu/sessions in types.go) is the only place state can
// ride between a Run() that returns Interrupted and the matching Resume()
// call -- keyed by sessionKey(input).
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
	findingBag          *FindingBag
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
			systemPrompt: InjectSkills(reconSystemPrompt, "recon"),
			tools:        withAgentTools(o.tools.Subset("nmap_scan", "smbmap_scan", "nuclei_scan")),
			gate:         func(name string) bool { return name == "nmap_scan" },
			exec: func(tr *tracker) toolExecFunc {
				return hybridExec(mcpExec(o.tools, tr), tr.bag, tr.store, "recon", tr)
			},
			maxTurns: maxReconModelTurns,
		}, true
	case "delegate_exploit":
		return subagentSpec{
			model:        o.model,
			systemPrompt: InjectSkills(exploitSystemPrompt, "exploit"),
			tools: withAgentTools(o.tools.Subset(
				"list_exploits", "list_payloads", "generate_payload", "run_exploit",
				"run_auxiliary_module", "run_post_module", "sqlmap_scan",
				"arp_scan_discovery", "hydra_attack", "rustscan_fast_scan",
				"responder_credential_harvest",
			)),
			exec: func(tr *tracker) toolExecFunc {
				return hybridExec(mcpExec(o.tools, tr), tr.bag, tr.store, "exploit", tr)
			},
			maxTurns: maxExploitModelTurns,
		}, true
	case "delegate_post_exploit":
		return subagentSpec{
			model:        o.model,
			systemPrompt: InjectSkills(postExploitSystemPrompt, "post_exploit"),
			tools:        withAgentTools(o.tools.Subset("list_active_sessions", "terminate_session", "send_session_command")),
			exec: func(tr *tracker) toolExecFunc {
				return hybridExec(mcpExec(o.tools, tr), tr.bag, tr.store, "post_exploit", tr)
			},
			maxTurns: maxPostExploitModelTurns,
		}, true
	case "delegate_codegen":
		if o.codegen == nil {
			// Avoid panicking on Name()/Call if wiring omitted the tool.
			return subagentSpec{}, false
		}
		return subagentSpec{
			model:        o.model,
			systemPrompt: InjectSkills(codeGenSystemPrompt, "codegen"),
			maxTurns:     maxCodegenModelTurns,
			tools: withAgentTools([]llm.ToolSpec{{
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
			}}),
			exec: func(tr *tracker) toolExecFunc {
				return hybridExec(codegenExec(o.codegen, tr), tr.bag, tr.store, "codegen", tr)
			},
		}, true
	case "delegate_report":
		return subagentSpec{
			model:        o.model,
			systemPrompt: InjectSkills(reportSystemPrompt, "report"),
			tools:        withAgentTools(nil),
			exec: func(tr *tracker) toolExecFunc {
				return hybridExec(func(ctx context.Context, call llm.ToolCall) (string, bool) {
					return "agent: report subagent supports report_finding, triage_finding, skill_search, skill_get, evidence_list, evidence_view, evidence_search only", true
				}, tr.bag, tr.store, "report", tr)
			},
		}, true
	}
	return subagentSpec{}, false
}

// finalizeResult attaches structured findings + multi-section report to a
// completed (non-interrupted) RunResult. Interrupted results are returned as-is.
func finalizeResult(input RunInput, result RunResult, tr *tracker) RunResult {
	if result.Interrupted {
		return result
	}
	bag := tr.bag
	extracted := ExtractFindings(input, result.ToolLog, result.FinalMessage, result.JudgeVerdict, result.JudgeSet)
	reported := bag.Snapshot()
	findings := MergeExtracted(reported, extracted)
	report := BuildReport(input, result.ToolLog, result.FinalMessage, findings, result.JudgeVerdict, result.JudgeSet)
	kc := AnalyzeKillChain(findings)
	meth := ComputeMethodology(result.ToolLog, input.AgentMode)
	// Embed kill chain + methodology into report markdown.
	if report.Markdown != "" {
		report.Markdown = report.Markdown + "\n" + kc.Summary + "\n" + FormatMethodologyMarkdown(meth)
	}
	// Append deterministic recap to the report.
	recapTracker := &ToolCallTracker{Calls: result.ToolLog}
	recap := BuildRecap(input.TargetIP, input.SessionID, time.Now().Add(-time.Hour), recapTracker, tr.store, findings)
	if recap.FormatMarkdown() != "" {
		report.Markdown = report.Markdown + "\n" + recap.FormatMarkdown()
	}
	// Persist findings to target state store for cross-run continuity.
	if input.TargetIP != "" {
		go persistTargetFindings(input.TargetIP, findings)
	}
	result.Findings = findings
	result.Report = &report
	result.KillChain = &kc
	result.Methodology = &meth
	// Prefer structured markdown as the operator-facing final message when
	// the agent left only a thin summary — still keep agent text inside the report.
	if strings.TrimSpace(result.FinalMessage) == "" && report.Markdown != "" {
		result.FinalMessage = report.Markdown
	}
	return result
}

// delegateToolSpecs is the synthetic tool surface exposed to the
// orchestrator's own model -- one callable per subagent, taking a single
// free-form "instructions" string, rather than handing the orchestrator
// raw MCP tool access directly. Filtered by agent mode when mode is set.
func delegateToolSpecs(mode string) []llm.ToolSpec {
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
	all := []llm.ToolSpec{
		{Name: "delegate_recon", Description: "Verifies if the target service is running on the IP.", InputSchema: schema},
		{Name: "delegate_exploit", Description: "Searches and executes exploits against the verified service.", InputSchema: schema},
		{Name: "delegate_post_exploit", Description: "After an exploit is deployed, uses meterpreter to interact with the session.", InputSchema: schema},
		{Name: "delegate_codegen", Description: "If by using tools you are not able to exploit, then use this agent.", InputSchema: schema},
		{Name: "delegate_report", Description: "Generates final validation report upon confirmed exploit success. May use report_finding.", InputSchema: schema},
	}
	allowed := AllowedDelegates(mode)
	out := make([]llm.ToolSpec, 0, len(all))
	for _, t := range all {
		if allowed[t.Name] {
			out = append(out, t)
		}
	}
	if len(out) == 0 {
		return all
	}
	return out
}

// seedPrompt builds the initial user message describing the target and
// attacker context.
func seedPrompt(input RunInput) string {
	mode := NormalizeAgentMode(input.AgentMode)

	// Check for prior target state (resume context)
	var priorState string
	if input.TargetIP != "" {
		ts := NewTargetStore("talon-data/targets")
		if state, err := ts.GetOrCreate(input.TargetIP); err == nil && (len(state.Findings) > 0 || len(state.ReconDims) > 0) {
			priorState = "\n\nPrior Target State:\n" + state.Summary()
		}
	}

	return fmt.Sprintf(
		"Target Info:\n"+
			"- IP: %s\n"+
			"- CVE ID: %s\n"+
			"- Service Name: %s\n"+
			"- Description: %s\n"+
			"- Agent Mode: %s\n\n"+
			"Attacker Context:\n"+
			"- LHOST: %s\n"+
			"- LPORT: %d\n\n"+
			"%s%s\n\n"+
			"Begin the validation workflow now. Use report_finding with 3-gate evidence for real vulns.",
		input.TargetIP, input.CVEID, input.ServiceName, input.Description, mode,
		input.Context.LHOST, input.Context.LPORT,
		AgentModePrompt(mode),
		priorState,
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
	input.AgentMode = NormalizeAgentMode(input.AgentMode)
	if resume == nil {
		messages := []llm.Message{llm.UserMessage(seedPrompt(input))}
		return o.orchestrateLoop(ctx, input, messages, &tracker{bag: NewFindingBag(), store: NewEvidenceStore(), cases: NewCaseStore(), correction: NewCorrectionLayer(), traffic: NewTrafficStore("talon-data/traffic", input.SessionID)})
	}
	return o.resumeRun(ctx, input, *resume)
}

func (o *Orchestrator) resumeRun(ctx context.Context, input RunInput, decision Decision) (RunResult, error) {
	input.AgentMode = NormalizeAgentMode(input.AgentMode)
	key := sessionKey(input)
	o.mu.Lock()
	sess, ok := o.sessions[key]
	if ok {
		delete(o.sessions, key)
	}
	o.mu.Unlock()
	if !ok {
		return RunResult{}, errors.New("agent: no pending interrupt to resume for this session")
	}

	bag := sess.findingBag
	if bag == nil {
		bag = NewFindingBag()
	}
	tr := &tracker{count: sess.toolCallCount, log: sess.toolLog, bag: bag, store: NewEvidenceStore(), cases: NewCaseStore(), correction: NewCorrectionLayer(), traffic: NewTrafficStore("talon-data/traffic", input.SessionID)}
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
			return finalizeResult(input, RunResult{FinalMessage: lastAssistantText(sess.orchestratorMessages), ToolLog: tr.log}, tr), nil
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
	if tr.bag == nil {
		tr.bag = NewFindingBag()
	}
	specs := delegateToolSpecs(input.AgentMode)
	sys := orchestratorSystemPrompt + "\n\n" + AgentModePrompt(input.AgentMode)
	for turn := 0; turn < maxOrchestratorTurns; turn++ {
		log.Printf("talon-core: orchestrator turn %d/%d tools_so_far=%d target=%s mode=%s",
			turn+1, maxOrchestratorTurns, tr.count, input.TargetIP, NormalizeAgentMode(input.AgentMode))

		msg, err := converseWithTimeout(ctx, o.model, sys, messages, specs)
		if err != nil {
			// Soft-fail on LLM timeout/errors after some progress: return
			// what we have rather than leaving the run stuck "running".
			if tr.count > 0 && (errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled)) {
				log.Printf("talon-core: orchestrator LLM timeout after progress: %v", err)
				return finalizeResult(input, RunResult{
					FinalMessage: lastAssistantText(messages) + "\n[orchestrator stopped: LLM timeout]",
					ToolLog:      tr.log,
				}, tr), nil
			}
			return RunResult{}, err
		}
		messages = append(messages, msg)

		if len(msg.ToolCalls) == 0 {
			log.Printf("talon-core: orchestrator final text; judging (tools=%d)", tr.count)
			verdict, err := judgeOutput(ctx, o.judge, msg.Text)
			if err != nil {
				log.Printf("talon-core: judge failed (returning without verdict): %v", err)
				return finalizeResult(input, RunResult{FinalMessage: msg.Text, ToolLog: tr.log}, tr), nil
			}
			return finalizeResult(input, RunResult{FinalMessage: msg.Text, ToolLog: tr.log, JudgeVerdict: verdict, JudgeSet: true}, tr), nil
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
				return finalizeResult(input, RunResult{FinalMessage: lastAssistantText(messages), ToolLog: tr.log}, tr), nil
			}
			if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
				log.Printf("talon-core: delegate batch timeout: %v", err)
				return finalizeResult(input, RunResult{
					FinalMessage: lastAssistantText(messages) + "\n[orchestrator stopped: delegate timeout]",
					ToolLog:      tr.log,
				}, tr), nil
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
	return finalizeResult(input, RunResult{
		FinalMessage: lastAssistantText(messages) + "\n[orchestrator stopped: turn budget reached]",
		ToolLog:      tr.log,
	}, tr), nil
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

	// agentMode is not on the batch; mode filtering is enforced via tool specs.
	// Still reject unknown names.
	for i, tc := range calls {
		sub, ok := o.subagentConfig(tc.Name)
		if !ok {
			resolved = append(resolved, llm.ToolResult{ToolCallID: tc.ID, Name: tc.Name, Content: fmt.Sprintf("agent: unknown or disabled delegate tool %q", tc.Name), IsError: true})
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
		o.sessions = make(map[string]*orchestratorSession)
	}
	o.sessions[sessionKey(input)] = &orchestratorSession{
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
		findingBag:          tr.bag,
	}
}

// persistTargetFindings stores verified findings to the target state store.
func persistTargetFindings(target string, findings []Finding) {
	ts := NewTargetStore("talon-data/targets")
	state, err := ts.GetOrCreate(target)
	if err != nil {
		return
	}
	for _, f := range findings {
		status := "unverified"
		if f.Evidence.Passed {
			status = "verified"
		}
		state.AddFinding(TargetFinding{
			Title:    f.Title,
			Severity: f.Severity,
			Status:   status,
		})
	}
	_ = ts.Save(state)
}
