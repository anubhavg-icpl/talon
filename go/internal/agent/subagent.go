package agent

import (
	"context"

	"github.com/anubhavg-icpl/pentester2/go/internal/llm"
	"github.com/anubhavg-icpl/pentester2/go/internal/mcpclient"
)

// toolExecFunc runs one real tool call (an MCP tool or the injected codegen
// tool) and records it into the tracker. Tool-level failures never surface
// as a Go error here -- they come back as (output, isErr=true) so the model
// sees them and can react, mirroring how a failed ToolMessage still flows
// back into the conversation in final.py.
type toolExecFunc func(ctx context.Context, call llm.ToolCall) (output string, isErr bool)

func mcpExec(tools *mcpclient.Multi, tr *tracker) toolExecFunc {
	return func(ctx context.Context, call llm.ToolCall) (string, bool) {
		out, err := tools.Call(ctx, call.Name, call.Args)
		isErr := err != nil
		if isErr && out == "" {
			out = err.Error()
		}
		tr.record(call.Name, call.Args, out)
		return out, isErr
	}
}

func codegenExec(codegen CodegenTool, tr *tracker) toolExecFunc {
	return func(ctx context.Context, call llm.ToolCall) (string, bool) {
		query, _ := call.Args["query"].(string)
		out, err := codegen.Call(ctx, query)
		isErr := err != nil
		if isErr && out == "" {
			out = err.Error()
		}
		tr.record(call.Name, call.Args, out)
		return out, isErr
	}
}

// subagentResume carries a paused nested subagent loop across the
// Run()/Resume() boundary -- see orchestrator.go's orchestratorSession for
// how this rides along on the Orchestrator between calls.
type subagentResume struct {
	messages        []llm.Message
	resolvedResults []llm.ToolResult
	remainingCalls  []llm.ToolCall
	gatedCall       llm.ToolCall
	decision        Decision
}

// subagentInterrupt is returned by runSubagent when it hits a HITL-gated
// tool call (nmap_scan) instead of executing it.
type subagentInterrupt struct {
	callID          string
	toolName        string
	args            map[string]any
	messages        []llm.Message
	resolvedResults []llm.ToolResult
	remainingCalls  []llm.ToolCall
}

// applyDecision executes (or rejects) a gated tool call per a human
// decision, mirroring HumanInTheLoopMiddleware's approve/edit/reject
// handling in final.py's main() loop.
func applyDecision(ctx context.Context, call llm.ToolCall, decision Decision, exec toolExecFunc, tr *tracker) (llm.ToolResult, error) {
	switch decision.Type {
	case "approve", "edit":
		args := call.Args
		if decision.Type == "edit" {
			args = decision.EditedArgs
		}
		if err := tr.allow(); err != nil {
			return llm.ToolResult{}, err
		}
		out, isErr := exec(ctx, llm.ToolCall{ID: call.ID, Name: call.Name, Args: args})
		return llm.ToolResult{ToolCallID: call.ID, Name: call.Name, Content: out, IsError: isErr}, nil
	case "reject":
		return llm.ToolResult{ToolCallID: call.ID, Name: call.Name, Content: "Human reviewer rejected this tool call.", IsError: true}, nil
	default:
		return llm.ToolResult{ToolCallID: call.ID, Name: call.Name, Content: "agent: unknown decision type " + decision.Type, IsError: true}, nil
	}
}

// runSubagent drives one subagent's nested tool-calling loop to completion
// (final text, no more tool calls requested) or until it hits a HITL-gated
// tool call, mirroring the isolated inner create_agent() loop
// SubAgentMiddleware runs for each configured subagent in final.py.
//
// Pass resume to continue a loop previously paused by an interrupt:
// resume.gatedCall is executed per resume.decision, then the loop carries on
// with resume.remainingCalls before asking the model for its next turn.
func runSubagent(ctx context.Context, model llm.ChatModel, systemPrompt string, tools []llm.ToolSpec, task string, gate func(name string) bool, exec toolExecFunc, resume *subagentResume, tr *tracker) (string, *subagentInterrupt, error) {
	var messages []llm.Message
	var pendingResults []llm.ToolResult
	var pendingCalls []llm.ToolCall
	skipModelCall := resume != nil

	if resume != nil {
		messages = resume.messages
		result, err := applyDecision(ctx, resume.gatedCall, resume.decision, exec, tr)
		if err != nil {
			return "", nil, err
		}
		pendingResults = append(append([]llm.ToolResult{}, resume.resolvedResults...), result)
		pendingCalls = resume.remainingCalls
	} else {
		messages = []llm.Message{llm.UserMessage(task)}
	}

	for {
		if !skipModelCall {
			msg, err := model.Converse(ctx, systemPrompt, messages, tools)
			if err != nil {
				return "", nil, err
			}
			if len(msg.ToolCalls) == 0 {
				return msg.Text, nil, nil
			}
			messages = append(messages, msg)
			pendingCalls = msg.ToolCalls
			pendingResults = nil
		}
		skipModelCall = false

		for i, tc := range pendingCalls {
			if gate != nil && gate(tc.Name) {
				return "", &subagentInterrupt{
					callID:          tc.ID,
					toolName:        tc.Name,
					args:            tc.Args,
					messages:        messages,
					resolvedResults: pendingResults,
					remainingCalls:  append([]llm.ToolCall{}, pendingCalls[i+1:]...),
				}, nil
			}
			if err := tr.allow(); err != nil {
				return "", nil, err
			}
			out, isErr := exec(ctx, tc)
			pendingResults = append(pendingResults, llm.ToolResult{ToolCallID: tc.ID, Name: tc.Name, Content: out, IsError: isErr})
		}

		messages = append(messages, llm.Message{Role: llm.RoleTool, ToolResults: pendingResults})
		pendingCalls = nil
		pendingResults = nil
	}
}
