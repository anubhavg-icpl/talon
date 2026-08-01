package control

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/anubhavg-icpl/talon/internal/config"
	"github.com/anubhavg-icpl/talon/internal/llm"
)

const (
	// slmMaxToolRounds bounds the assist tool loop (SLMs can thrash).
	slmMaxToolRounds = 5
	// slmAssistTimeout is the wall clock for one assist request.
	slmAssistTimeout = 4 * time.Minute
)

// handleLLMTools is GET /llm/tools — curated UI-safe tool catalog for the
// dashboard assist UI (and clients that build prompts client-side).
func (s *Server) handleLLMTools(w http.ResponseWriter, r *http.Request) {
	cat := slmToolCatalog()
	writeJSON(w, http.StatusOK, map[string]any{
		"tools":       cat,
		"count":       len(cat),
		"protocol":    "TOOL_CALL {\"name\":\"…\",\"arguments\":{…}}",
		"assist_path": "/llm/assist",
		"stream_path": "/llm/stream",
		"note":        "Read-only platform tools for SmolLM/UI assist. Full MCP exploit tools remain on agent runs only.",
	})
}

// llmAssistRequest is POST /llm/assist — multi-turn tool loop with SSE.
type llmAssistRequest struct {
	Messages    []llmStreamMsg `json:"messages"`
	System      string         `json:"system,omitempty"`
	Role        string         `json:"role,omitempty"`
	Model       string         `json:"model,omitempty"`
	MaxTokens   int32          `json:"max_tokens,omitempty"`
	MaxRounds   int            `json:"max_rounds,omitempty"`
	Temperature *float32       `json:"temperature,omitempty"`
	// DisableTools forces plain chat (same as /llm/stream).
	DisableTools bool `json:"disable_tools,omitempty"`
}

// handleLLMAssist is POST /llm/assist — end-to-end SLM assist with codebase tools.
//
// SSE events (aligned with /llm/stream plus tool events):
//
//	event: meta         provider, model, tools_enabled, tool_count
//	event: token        content delta (assistant prose)
//	event: tool_start   name, arguments
//	event: tool_result  name, result (truncated), ms
//	event: round        n, max
//	event: done         text, ms, tool_calls
//	event: error        error
func (s *Server) handleLLMAssist(w http.ResponseWriter, r *http.Request) {
	var req llmAssistRequest
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

	ctx, cancel := context.WithTimeout(r.Context(), slmAssistTimeout)
	defer cancel()

	model, err := llm.NewModel(ctx, llmCfg, provider, modelID)
	if err != nil {
		writeError(w, http.StatusBadGateway, "llm init: "+err.Error())
		return
	}

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
	// Keepalive comments so proxies/browsers don't treat a quiet Converse as a hung stream.
	writePing := func() {
		_, _ = fmt.Fprintf(w, ": ping %d\n\n", time.Now().Unix())
		flusher.Flush()
	}

	toolsOn := !req.DisableTools
	writeSSE("meta", map[string]any{
		"provider":      provider,
		"model":         modelID,
		"role":          role,
		"tools_enabled": toolsOn,
		"tool_count":    len(slmToolCatalog()),
		"protocol":      "text_tool_call",
		"max_rounds":    clampRounds(req.MaxRounds),
	})

	system := slmSystemPrompt(req.System)
	var messages []llm.Message
	for _, m := range req.Messages {
		switch m.Role {
		case "system":
			if req.System == "" && m.Content != "" {
				system = slmSystemPrompt(m.Content)
			}
		case "assistant":
			messages = append(messages, llm.AssistantText(m.Content))
		case "tool":
			// Client may replay prior tool results as user-context.
			messages = append(messages, llm.UserMessage("[tool result]\n"+m.Content))
		default:
			messages = append(messages, llm.UserMessage(m.Content))
		}
	}

	start := time.Now()
	if !toolsOn {
		s.streamPlain(ctx, model, system, messages, writeSSE, start)
		return
	}

	maxRounds := clampRounds(req.MaxRounds)
	var finalText strings.Builder
	toolCalls := 0

	// Prefer native function-calling when the provider supports tools well
	// (openai/ollama). onnx/SmolLM uses text TOOL_CALL protocol via Converse.
	useNative := provider == "openai" || provider == "ollama"

	for round := 1; round <= maxRounds; round++ {
		if ctx.Err() != nil {
			writeSSE("error", map[string]string{"error": ctx.Err().Error()})
			return
		}
		writeSSE("round", map[string]any{"n": round, "max": maxRounds})
		writeSSE("status", map[string]any{
			"phase":   "thinking",
			"message": fmt.Sprintf("round %d/%d — waiting on %s/%s…", round, maxRounds, provider, modelID),
		})

		var out llm.Message
		var genErr error

		// Heartbeat while the model call blocks (OpenAI native path has no mid-call tokens).
		stopHB := startSSEHeartbeat(ctx, writePing)
		if useNative {
			out, genErr = model.Converse(ctx, system, messages, slmToolSpecs())
		} else {
			// Stream tokens for SLM, then parse TOOL_CALL from full text.
			if streamer, ok := llm.AsStreamer(model); ok {
				var buf strings.Builder
				out, genErr = streamer.ConverseStream(ctx, system, messages, func(tok string) error {
					buf.WriteString(tok)
					// Stream only non-tool prose to the UI; still buffer all.
					writeSSE("token", map[string]string{"content": tok})
					return nil
				})
				if genErr == nil && out.Text == "" {
					out.Text = buf.String()
				}
			} else {
				out, genErr = model.Converse(ctx, system, messages, nil)
				if genErr == nil && out.Text != "" {
					writeSSE("token", map[string]string{"content": out.Text})
				}
			}
		}
		stopHB()
		if genErr != nil {
			writeSSE("error", map[string]string{"error": genErr.Error()})
			return
		}

		// Native tool calls.
		if useNative && len(out.ToolCalls) > 0 {
			messages = append(messages, out)
			for _, tc := range out.ToolCalls {
				if !knownSLMTool(tc.Name) {
					writeSSE("tool_result", map[string]any{
						"name": tc.Name, "result": "error: tool not in SLM catalog", "ms": 0,
					})
					messages = append(messages, llm.ToolResultMessage(llm.ToolResult{
						ToolCallID: tc.ID, Name: tc.Name,
						Content: "error: tool not allowed in assist mode", IsError: true,
					}))
					continue
				}
				writeSSE("status", map[string]any{"phase": "tool", "message": "running " + tc.Name})
				writeSSE("tool_start", map[string]any{"name": tc.Name, "arguments": tc.Args})
				t0 := time.Now()
				result := s.execSLMTool(tc.Name, tc.Args)
				ms := time.Since(t0).Milliseconds()
				toolCalls++
				writeSSE("tool_result", map[string]any{
					"name": tc.Name, "result": clipStr(result, 6000), "ms": ms,
				})
				messages = append(messages, llm.ToolResultMessage(llm.ToolResult{
					ToolCallID: tc.ID, Name: tc.Name, Content: result,
				}))
			}
			continue
		}

		// Text protocol TOOL_CALL (SmolLM / onnx / plain models).
		if name, args, ok := parseSLMToolCall(out.Text); ok && knownSLMTool(name) {
			// Strip tool line from visible answer for this round.
			prose := stripToolCallLines(out.Text)
			if prose != "" {
				finalText.WriteString(prose)
			}
			messages = append(messages, llm.AssistantText(out.Text))
			writeSSE("status", map[string]any{"phase": "tool", "message": "running " + name})
			writeSSE("tool_start", map[string]any{"name": name, "arguments": args})
			t0 := time.Now()
			result := s.execSLMTool(name, args)
			ms := time.Since(t0).Milliseconds()
			toolCalls++
			writeSSE("tool_result", map[string]any{
				"name": name, "result": clipStr(result, 6000), "ms": ms,
			})
			// Feed result as a user turn (OpenAI tool role optional for SLMs).
			messages = append(messages, llm.UserMessage(
				fmt.Sprintf("TOOL_RESULT %s:\n%s\n\nUse this data to answer the operator. Call another tool if needed, else reply with the final answer (no TOOL_CALL).", name, result),
			))
			continue
		}

		// Final assistant answer (no more tools).
		writeSSE("status", map[string]any{"phase": "answering", "message": "composing answer"})
		if useNative && out.Text != "" {
			writeSSE("token", map[string]string{"content": out.Text})
		}
		// Prefer the model’s final prose when tools already flushed partial text.
		answer := strings.TrimSpace(out.Text)
		if answer == "" {
			answer = strings.TrimSpace(finalText.String())
		} else {
			finalText.Reset()
			finalText.WriteString(answer)
		}
		writeSSE("done", map[string]any{
			"text":       answer,
			"ms":         time.Since(start).Milliseconds(),
			"tool_calls": toolCalls,
			"rounds":     round,
		})
		log.Printf("llm/assist provider=%s model=%s ms=%d tools=%d rounds=%d",
			provider, modelID, time.Since(start).Milliseconds(), toolCalls, round)
		return
	}

	writeSSE("done", map[string]any{
		"text":       strings.TrimSpace(finalText.String()),
		"ms":         time.Since(start).Milliseconds(),
		"tool_calls": toolCalls,
		"rounds":     maxRounds,
		"note":       "max tool rounds reached",
	})
}

func (s *Server) streamPlain(ctx context.Context, model llm.ChatModel, system string, messages []llm.Message, writeSSE func(string, any), start time.Time) {
	if streamer, ok := llm.AsStreamer(model); ok {
		out, err := streamer.ConverseStream(ctx, system, messages, func(tok string) error {
			writeSSE("token", map[string]string{"content": tok})
			return nil
		})
		if err != nil {
			writeSSE("error", map[string]string{"error": err.Error()})
			return
		}
		writeSSE("done", map[string]any{"text": out.Text, "ms": time.Since(start).Milliseconds(), "tool_calls": 0})
		return
	}
	out, err := model.Converse(ctx, system, messages, nil)
	if err != nil {
		writeSSE("error", map[string]string{"error": err.Error()})
		return
	}
	if out.Text != "" {
		writeSSE("token", map[string]string{"content": out.Text})
	}
	writeSSE("done", map[string]any{"text": out.Text, "ms": time.Since(start).Milliseconds(), "tool_calls": 0})
}

func clampRounds(n int) int {
	if n <= 0 {
		return slmMaxToolRounds
	}
	if n > 8 {
		return 8
	}
	return n
}

// startSSEHeartbeat writes SSE comments every 2s until stop is called.
func startSSEHeartbeat(ctx context.Context, ping func()) (stop func()) {
	done := make(chan struct{})
	var once sync.Once
	go func() {
		t := time.NewTicker(2 * time.Second)
		defer t.Stop()
		for {
			select {
			case <-done:
				return
			case <-ctx.Done():
				return
			case <-t.C:
				ping()
			}
		}
	}()
	return func() { once.Do(func() { close(done) }) }
}

func clipStr(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}

func stripToolCallLines(text string) string {
	var b strings.Builder
	for _, line := range strings.Split(text, "\n") {
		trim := strings.TrimSpace(line)
		if strings.HasPrefix(trim, "TOOL_CALL") {
			continue
		}
		if strings.Contains(trim, "<tool_call>") {
			continue
		}
		b.WriteString(line)
		b.WriteByte('\n')
	}
	return strings.TrimSpace(b.String())
}
