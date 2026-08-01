package control

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/anubhavg-icpl/talon/internal/config"
	"github.com/anubhavg-icpl/talon/internal/llm"
)

// llmStreamRequest is the body for POST /llm/stream — plain chat, no tools.
// Tokens are pushed as Server-Sent Events so the dashboard (and any WASM
// client) can render partial output in milliseconds.
type llmStreamRequest struct {
	// Messages are OpenAI-style chat turns. system may also be set separately.
	Messages []llmStreamMsg `json:"messages"`
	System   string         `json:"system,omitempty"`
	// Role selects which provider/model pair to use (main|judge|code).
	Role string `json:"role,omitempty"`
	// Model overrides the resolved model id when non-empty.
	Model string `json:"model,omitempty"`
	// MaxTokens caps generation (0 → LLMConfig.MaxTokens).
	MaxTokens int32 `json:"max_tokens,omitempty"`
	// Temperature; negative means "use config default".
	Temperature *float32 `json:"temperature,omitempty"`
}

type llmStreamMsg struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// handleLLMStream is POST /llm/stream — SSE of token deltas from the active
// LLM provider. Prefer LLM_PROVIDER=onnx (SmolLM/ONNX) or openai for true
// streaming; non-streaming backends emit a single content event.
//
// Event types:
//
//	event: meta     data: {"provider","model","backend"}
//	event: token    data: {"content":"..."}
//	event: done     data: {"text":"<full>","ms":N}
//	event: error    data: {"error":"..."}
func (s *Server) handleLLMStream(w http.ResponseWriter, r *http.Request) {
	var req llmStreamRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON: "+err.Error())
		return
	}
	if len(req.Messages) == 0 {
		writeError(w, http.StatusBadRequest, "messages required")
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		writeError(w, http.StatusInternalServerError, "streaming unsupported")
		return
	}

	role := req.Role
	if role == "" {
		role = config.RoleMain
	}
	llmCfg := config.LoadLLMConfig()
	provider, modelID := config.ResolveModel(llmCfg, role)
	if req.Model != "" {
		modelID = req.Model
	}

	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Minute)
	defer cancel()

	model, err := llm.NewModel(ctx, llmCfg, provider, modelID)
	if err != nil {
		writeError(w, http.StatusBadGateway, "llm init: "+err.Error())
		return
	}

	// SSE headers — disable proxy buffering for millisecond delivery.
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")
	w.WriteHeader(http.StatusOK)
	flusher.Flush()

	writeSSE := func(event string, v any) {
		raw, _ := json.Marshal(v)
		_, _ = fmt.Fprintf(w, "event: %s\ndata: %s\n\n", event, raw)
		flusher.Flush()
	}

	writeSSE("meta", map[string]any{
		"provider": provider,
		"model":    modelID,
		"role":     role,
	})

	system := req.System
	var messages []llm.Message
	for _, m := range req.Messages {
		switch m.Role {
		case "system":
			// Prefer first system in messages if top-level system empty.
			if system == "" {
				system = m.Content
			}
		case "assistant":
			messages = append(messages, llm.AssistantText(m.Content))
		default:
			messages = append(messages, llm.UserMessage(m.Content))
		}
	}

	start := time.Now()
	if streamer, ok := llm.AsStreamer(model); ok {
		out, err := streamer.ConverseStream(ctx, system, messages, func(tok string) error {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			writeSSE("token", map[string]string{"content": tok})
			return nil
		})
		if err != nil {
			writeSSE("error", map[string]string{"error": err.Error()})
			return
		}
		writeSSE("done", map[string]any{
			"text": out.Text,
			"ms":   time.Since(start).Milliseconds(),
		})
		return
	}

	// Fallback: non-streaming backends (Bedrock, Ollama today) — one shot.
	out, err := model.Converse(ctx, system, messages, nil)
	if err != nil {
		writeSSE("error", map[string]string{"error": err.Error()})
		return
	}
	if out.Text != "" {
		writeSSE("token", map[string]string{"content": out.Text})
	}
	writeSSE("done", map[string]any{
		"text": out.Text,
		"ms":   time.Since(start).Milliseconds(),
	})
	log.Printf("llm/stream provider=%s model=%s ms=%d (non-stream backend)", provider, modelID, time.Since(start).Milliseconds())
}

// handleLLMInfo is GET /llm/info — which provider/model the stream path would use.
func (s *Server) handleLLMInfo(w http.ResponseWriter, r *http.Request) {
	llmCfg := config.LoadLLMConfig()
	role := r.URL.Query().Get("role")
	if role == "" {
		role = config.RoleMain
	}
	provider, modelID := config.ResolveModel(llmCfg, role)
	writeJSON(w, http.StatusOK, map[string]any{
		"provider":      provider,
		"model":         modelID,
		"role":          role,
		"onnx_base_url": llmCfg.ONNXBaseURL,
		"ollama_url":    llmCfg.OllamaURL,
		"openai_base":   llmCfg.OpenAIBaseURL,
		"stream_path":   "/llm/stream",
		"assist_path":   "/llm/assist",
		"tools_path":    "/llm/tools",
		"tool_count":    len(slmToolCatalog()),
	})
}
