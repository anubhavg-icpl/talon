// Package llm defines the model-agnostic chat/tool-use contract every agent
// (orchestrator, subagents, judge) is built against, and a Bedrock Converse
// implementation of it. Mirrors the role ChatBedrockConverse/ChatOllama play
// in pentest_core/final.py.
package llm

import "context"

type Role string

const (
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
	RoleTool      Role = "tool"
)

// ToolCall is a single tool invocation requested by the model.
type ToolCall struct {
	ID   string
	Name string
	Args map[string]any
}

// ToolResult is fed back to the model as the output of a prior ToolCall.
type ToolResult struct {
	ToolCallID string
	Name       string
	Content    string
	IsError    bool
}

// Message is one turn in the conversation. Exactly one of Text/ToolCalls is
// populated on assistant turns produced by the model; ToolResults is
// populated on turns fed back to the model after executing tool calls.
type Message struct {
	Role        Role
	Text        string
	ToolCalls   []ToolCall
	ToolResults []ToolResult
}

func UserMessage(text string) Message   { return Message{Role: RoleUser, Text: text} }
func AssistantText(text string) Message { return Message{Role: RoleAssistant, Text: text} }
func ToolResultMessage(r ToolResult) Message {
	return Message{Role: RoleTool, ToolResults: []ToolResult{r}}
}

// ToolSpec describes a callable tool to the model, equivalent to an MCP
// tool's name/description/inputSchema triple.
type ToolSpec struct {
	Name        string
	Description string
	InputSchema map[string]any // raw JSON Schema
}

// ChatModel is implemented by internal/llm's Bedrock client and by any test
// double. Converse returns the model's next turn: either assistant text
// (done) or one or more tool calls to execute before calling Converse again.
type ChatModel interface {
	Converse(ctx context.Context, systemPrompt string, messages []Message, tools []ToolSpec) (Message, error)
}
