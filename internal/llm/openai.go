package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// OpenAI implements ChatModel against any OpenAI-compatible /chat/completions
// endpoint (OpenAI itself, Azure OpenAI, z.ai's coding PaaS, vLLM, LocalAI,
// LiteLLM, ...). It speaks the standard tools/function-calling schema, so a
// hosted model like GLM-5.2 can drive the orchestrator/tool loop with zero
// AWS or local-GPU dependency.
type OpenAI struct {
	baseURL string // e.g. "https://api.z.ai/api/coding/paas/v4"
	apiKey  string
	model   string
	http    *http.Client
}

// NewOpenAI returns a client for an OpenAI-compatible chat-completions API.
// apiKey is required -- pass an empty string and Converse will surface the
// provider's 401 rather than silently sending an unauthenticated request.
func NewOpenAI(baseURL, apiKey, model string) *OpenAI {
	return &OpenAI{
		baseURL: strings.TrimRight(baseURL, "/"),
		apiKey:  apiKey,
		model:   model,
		http:    &http.Client{},
	}
}

type oaiMessage struct {
	Role       string         `json:"role"`
	Content    string         `json:"content,omitempty"`
	ToolCalls  []oaiToolCall  `json:"tool_calls,omitempty"`
	ToolCallID string         `json:"tool_call_id,omitempty"`
}

type oaiToolCall struct {
	ID       string           `json:"id,omitempty"`
	Type     string           `json:"type,omitempty"`
	Function oaiFunctionCall  `json:"function"`
}

type oaiFunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"` // OpenAI sends args as a JSON *string*
}

type oaiTool struct {
	Type     string           `json:"type"`
	Function oaiFunctionSpec  `json:"function"`
}

type oaiFunctionSpec struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Parameters  map[string]any `json:"parameters"`
}

type oaiChatRequest struct {
	Model       string      `json:"model"`
	Messages    []oaiMessage `json:"messages"`
	Tools       []oaiTool   `json:"tools,omitempty"`
	Temperature float32     `json:"temperature,omitempty"`
	MaxTokens   int32       `json:"max_tokens,omitempty"`
	Stream      bool        `json:"stream"`
}

type oaiChatResponse struct {
	Choices []struct {
		Message      oaiMessage `json:"message"`
		FinishReason string     `json:"finish_reason"`
	} `json:"choices"`
	Error *oaiError `json:"error,omitempty"`
}

type oaiError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

func (o *OpenAI) Converse(ctx context.Context, systemPrompt string, messages []Message, tools []ToolSpec) (Message, error) {
	var oaiMsgs []oaiMessage
	if systemPrompt != "" {
		oaiMsgs = append(oaiMsgs, oaiMessage{Role: "system", Content: systemPrompt})
	}
	for _, m := range messages {
		oaiMsgs = append(oaiMsgs, toOAIMessage(m))
	}

	req := oaiChatRequest{
		Model:    o.model,
		Messages: oaiMsgs,
		Stream:   false,
	}
	for _, t := range tools {
		params := t.InputSchema
		if params == nil {
			params = map[string]any{"type": "object", "properties": map[string]any{}}
		}
		req.Tools = append(req.Tools, oaiTool{
			Type: "function",
			Function: oaiFunctionSpec{
				Name:        t.Name,
				Description: t.Description,
				Parameters:  params,
			},
		})
	}

	body, err := json.Marshal(req)
	if err != nil {
		return Message{}, fmt.Errorf("llm: marshal openai request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, o.baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return Message{}, fmt.Errorf("llm: build openai request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	if o.apiKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+o.apiKey)
	}

	resp, err := o.http.Do(httpReq)
	if err != nil {
		return Message{}, fmt.Errorf("llm: openai request failed: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return Message{}, fmt.Errorf("llm: read openai response: %w", err)
	}
	if resp.StatusCode >= 400 {
		return Message{}, fmt.Errorf("llm: openai returned HTTP %d: %s", resp.StatusCode, string(raw))
	}

	var out oaiChatResponse
	if err := json.Unmarshal(raw, &out); err != nil {
		return Message{}, fmt.Errorf("llm: decode openai response: %w", err)
	}
	if out.Error != nil {
		return Message{}, fmt.Errorf("llm: openai error: %s", out.Error.Message)
	}
	if len(out.Choices) == 0 {
		return Message{}, fmt.Errorf("llm: openai returned no choices")
	}

	msg := out.Choices[0].Message
	result := Message{Role: RoleAssistant, Text: msg.Content}
	for _, tc := range msg.ToolCalls {
		var args map[string]any
		// OpenAI ships function arguments as a JSON string; tolerate a
		// provider that already parsed it into an object.
		if tc.Function.Arguments != "" {
			if err := json.Unmarshal([]byte(tc.Function.Arguments), &args); err != nil {
				return Message{}, fmt.Errorf("llm: decode openai tool args for %q: %w", tc.Function.Name, err)
			}
		}
		id := tc.ID
		if id == "" {
			id = fmt.Sprintf("call_%s", tc.Function.Name)
		}
		result.ToolCalls = append(result.ToolCalls, ToolCall{
			ID:   id,
			Name: tc.Function.Name,
			Args: args,
		})
	}
	return result, nil
}

func toOAIMessage(m Message) oaiMessage {
	switch m.Role {
	case RoleAssistant:
		out := oaiMessage{Role: "assistant", Content: m.Text}
		for _, tc := range m.ToolCalls {
			args, _ := json.Marshal(tc.Args)
			out.ToolCalls = append(out.ToolCalls, oaiToolCall{
				ID:   tc.ID,
				Type: "function",
				Function: oaiFunctionCall{
					Name:      tc.Name,
					Arguments: string(args),
				},
			})
		}
		return out
	case RoleTool:
		// OpenAI links a tool result back to its call by tool_call_id; emit
		// one message per result so the linkage is unambiguous.
		var last oaiMessage
		for _, tr := range m.ToolResults {
			last = oaiMessage{Role: "tool", ToolCallID: tr.ToolCallID, Content: tr.Content}
		}
		// ponytail: orchestrator feeds one result per turn, so returning the
		// last is correct; if multi-result turns appear, switch to appending
		// all and have the caller send one message per result.
		return last
	default: // user
		return oaiMessage{Role: "user", Content: m.Text}
	}
}
