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

// Ollama implements ChatModel against a local Ollama server's /api/chat
// endpoint (https://github.com/ollama/ollama), so Talon can run with zero
// AWS dependency for inference -- swap in wherever a *Bedrock is used.
type Ollama struct {
	baseURL string
	model   string
	http    *http.Client
}

func NewOllama(baseURL, model string) *Ollama {
	return &Ollama{
		baseURL: strings.TrimRight(baseURL, "/"),
		model:   model,
		http:    &http.Client{},
	}
}

type ollamaMessage struct {
	Role      string          `json:"role"`
	Content   string          `json:"content,omitempty"`
	ToolCalls []ollamaToolUse `json:"tool_calls,omitempty"`
}

type ollamaToolUse struct {
	Function ollamaFunctionCall `json:"function"`
}

type ollamaFunctionCall struct {
	Name      string         `json:"name"`
	Arguments map[string]any `json:"arguments"`
}

type ollamaTool struct {
	Type     string             `json:"type"`
	Function ollamaFunctionSpec `json:"function"`
}

type ollamaFunctionSpec struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Parameters  map[string]any `json:"parameters"`
}

type ollamaChatRequest struct {
	Model    string          `json:"model"`
	Messages []ollamaMessage `json:"messages"`
	Tools    []ollamaTool    `json:"tools,omitempty"`
	Stream   bool            `json:"stream"`
}

type ollamaChatResponse struct {
	Message ollamaMessage `json:"message"`
	Done    bool          `json:"done"`
	Error   string        `json:"error"`
}

func (o *Ollama) Converse(ctx context.Context, systemPrompt string, messages []Message, tools []ToolSpec) (Message, error) {
	var ollamaMessages []ollamaMessage
	if systemPrompt != "" {
		ollamaMessages = append(ollamaMessages, ollamaMessage{Role: "system", Content: systemPrompt})
	}
	for _, m := range messages {
		ollamaMessages = append(ollamaMessages, toOllamaMessage(m))
	}

	req := ollamaChatRequest{
		Model:    o.model,
		Messages: ollamaMessages,
		Stream:   false,
	}
	for _, t := range tools {
		params := t.InputSchema
		if params == nil {
			params = map[string]any{"type": "object", "properties": map[string]any{}}
		}
		req.Tools = append(req.Tools, ollamaTool{
			Type: "function",
			Function: ollamaFunctionSpec{
				Name:        t.Name,
				Description: t.Description,
				Parameters:  params,
			},
		})
	}

	body, err := json.Marshal(req)
	if err != nil {
		return Message{}, fmt.Errorf("llm: marshal ollama request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, o.baseURL+"/api/chat", bytes.NewReader(body))
	if err != nil {
		return Message{}, fmt.Errorf("llm: build ollama request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := o.http.Do(httpReq)
	if err != nil {
		return Message{}, fmt.Errorf("llm: ollama request failed: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return Message{}, fmt.Errorf("llm: read ollama response: %w", err)
	}
	if resp.StatusCode >= 400 {
		return Message{}, fmt.Errorf("llm: ollama returned HTTP %d: %s", resp.StatusCode, string(raw))
	}

	var out ollamaChatResponse
	if err := json.Unmarshal(raw, &out); err != nil {
		return Message{}, fmt.Errorf("llm: decode ollama response: %w", err)
	}
	if out.Error != "" {
		return Message{}, fmt.Errorf("llm: ollama error: %s", out.Error)
	}

	result := Message{Role: RoleAssistant, Text: out.Message.Content}
	for i, tc := range out.Message.ToolCalls {
		result.ToolCalls = append(result.ToolCalls, ToolCall{
			// Ollama doesn't assign tool-call IDs the way Bedrock/OpenAI do;
			// synthesize a stable one so the rest of the orchestrator's
			// ID-keyed bookkeeping (matching ToolResult back to ToolCall)
			// works unchanged.
			ID:   fmt.Sprintf("call_%d", i),
			Name: tc.Function.Name,
			Args: tc.Function.Arguments,
		})
	}
	return result, nil
}

func toOllamaMessage(m Message) ollamaMessage {
	switch m.Role {
	case RoleAssistant:
		out := ollamaMessage{Role: "assistant", Content: m.Text}
		for _, tc := range m.ToolCalls {
			out.ToolCalls = append(out.ToolCalls, ollamaToolUse{
				Function: ollamaFunctionCall{Name: tc.Name, Arguments: tc.Args},
			})
		}
		return out
	case RoleTool:
		// Ollama has no tool_call_id linkage -- each tool result is just a
		// "tool" role message with the output text, matched by order.
		var content strings.Builder
		for i, tr := range m.ToolResults {
			if i > 0 {
				content.WriteString("\n")
			}
			content.WriteString(tr.Content)
		}
		return ollamaMessage{Role: "tool", Content: content.String()}
	default:
		return ollamaMessage{Role: "user", Content: m.Text}
	}
}
