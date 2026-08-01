package llm

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// TokenHandler is invoked for each content delta during streaming. Return a
// non-nil error to abort the stream (context cancel is preferred).
type TokenHandler func(token string) error

// Streamer is implemented by backends that can emit tokens as they generate.
// OpenAI-compatible endpoints (openai, onnx-slm, vLLM, …) implement this via
// ConverseStream on *OpenAI. Ollama has its own path if needed later.
type Streamer interface {
	// ConverseStream is chat-only (no tools). onToken receives each content
	// delta; the full assistant text is returned when the stream completes.
	ConverseStream(ctx context.Context, systemPrompt string, messages []Message, onToken TokenHandler) (Message, error)
}

// ConverseStream streams an OpenAI-compatible chat completion (stream=true)
// and fans content deltas to onToken. Tools are not sent — use Converse for
// agent tool loops. Designed for millisecond-class token delivery to SSE
// clients (dashboard UI / WASM consumers).
func (o *OpenAI) ConverseStream(ctx context.Context, systemPrompt string, messages []Message, onToken TokenHandler) (Message, error) {
	var oaiMsgs []oaiMessage
	if systemPrompt != "" {
		oaiMsgs = append(oaiMsgs, oaiMessage{Role: "system", Content: systemPrompt})
	}
	for _, m := range messages {
		oaiMsgs = append(oaiMsgs, toOAIMessage(m))
	}

	reqBody := oaiChatRequest{
		Model:    o.model,
		Messages: oaiMsgs,
		Stream:   true,
	}
	body, err := json.Marshal(reqBody)
	if err != nil {
		return Message{}, fmt.Errorf("llm: marshal stream request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, o.baseURL+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return Message{}, fmt.Errorf("llm: build stream request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "text/event-stream")
	if o.apiKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+o.apiKey)
	}

	// Streaming must not use the bounded http.Client.Timeout as a total
	// deadline — that would cut long generations. Clone transport, Timeout=0.
	client := &http.Client{Transport: o.http.Transport}
	resp, err := client.Do(httpReq)
	if err != nil {
		return Message{}, fmt.Errorf("llm: stream request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		raw, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return Message{}, fmt.Errorf("llm: stream HTTP %d: %s", resp.StatusCode, string(raw))
	}

	var full strings.Builder
	sc := bufio.NewScanner(resp.Body)
	// Token lines are small; bump buffer for safety on long reasoning chunks.
	sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	for sc.Scan() {
		if err := ctx.Err(); err != nil {
			return Message{}, err
		}
		line := sc.Text()
		if line == "" || strings.HasPrefix(line, ":") {
			continue
		}
		if !strings.HasPrefix(line, "data:") {
			continue
		}
		payload := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		if payload == "[DONE]" {
			break
		}

		var chunk oaiStreamChunk
		if err := json.Unmarshal([]byte(payload), &chunk); err != nil {
			// Tolerate non-JSON heartbeats / partial frames.
			continue
		}
		if chunk.Error != nil {
			return Message{}, fmt.Errorf("llm: stream error: %s", chunk.Error.Message)
		}
		for _, ch := range chunk.Choices {
			tok := ch.Delta.Content
			if tok == "" {
				continue
			}
			full.WriteString(tok)
			if onToken != nil {
				if err := onToken(tok); err != nil {
					return Message{}, err
				}
			}
		}
	}
	if err := sc.Err(); err != nil {
		return Message{}, fmt.Errorf("llm: read stream: %w", err)
	}
	return Message{Role: RoleAssistant, Text: full.String()}, nil
}

// oaiStreamChunk is one SSE data payload from OpenAI-compatible stream=true.
type oaiStreamChunk struct {
	Choices []struct {
		Delta struct {
			Role    string `json:"role"`
			Content string `json:"content"`
		} `json:"delta"`
		FinishReason *string `json:"finish_reason"`
	} `json:"choices"`
	Error *oaiError `json:"error,omitempty"`
}

// AsStreamer returns a Streamer if m supports token streaming.
func AsStreamer(m ChatModel) (Streamer, bool) {
	if s, ok := m.(Streamer); ok {
		return s, true
	}
	return nil, false
}
