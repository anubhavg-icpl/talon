package control

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/anubhavg-icpl/talon/internal/llm"
)

// stubModel returns fixed text; no Streamer so /llm/stream uses one-shot path.
type stubModel struct{ text string }

func (s stubModel) Converse(ctx context.Context, system string, messages []llm.Message, tools []llm.ToolSpec) (llm.Message, error) {
	return llm.AssistantText(s.text), nil
}

func TestLLMStreamSSE(t *testing.T) {
	// Force a known provider path that NewModel can build without AWS/keys.
	// Use onnx → OpenAI client; but our handler builds via NewModel. For unit
	// test we hit the non-stream fallback by mocking via env that won't be
	// used if we inject analyzer — handleLLMStream always NewModel's though.
	// So spin a fake OpenAI stream server and point ONNX_BASE_URL at it.
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasSuffix(r.URL.Path, "/chat/completions") {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/event-stream")
		fl := w.(http.Flusher)
		_, _ = w.Write([]byte(`data: {"choices":[{"delta":{"content":"pong"}}]}` + "\n\n"))
		fl.Flush()
		_, _ = w.Write([]byte("data: [DONE]\n\n"))
		fl.Flush()
	}))
	defer upstream.Close()

	os.Setenv("LLM_PROVIDER", "onnx")
	os.Setenv("ONNX_BASE_URL", upstream.URL+"/v1")
	os.Setenv("ONNX_MAIN_MODEL", "smollm")
	t.Cleanup(func() {
		os.Unsetenv("LLM_PROVIDER")
		os.Unsetenv("ONNX_BASE_URL")
		os.Unsetenv("ONNX_MAIN_MODEL")
	})

	s := NewServer(nil, NewStore())
	body := `{"messages":[{"role":"user","content":"ping"}]}`
	req := httptest.NewRequest(http.MethodPost, "/llm/stream", bytes.NewReader([]byte(body)))
	rr := httptest.NewRecorder()
	s.Mux().ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rr.Code, rr.Body.String())
	}
	ct := rr.Header().Get("Content-Type")
	if !strings.Contains(ct, "text/event-stream") {
		t.Fatalf("Content-Type=%q", ct)
	}

	var events []string
	var token string
	sc := bufio.NewScanner(rr.Body)
	var curEvent string
	for sc.Scan() {
		line := sc.Text()
		if strings.HasPrefix(line, "event:") {
			curEvent = strings.TrimSpace(strings.TrimPrefix(line, "event:"))
			events = append(events, curEvent)
		} else if strings.HasPrefix(line, "data:") && curEvent == "token" {
			var m map[string]string
			_ = json.Unmarshal([]byte(strings.TrimSpace(strings.TrimPrefix(line, "data:"))), &m)
			token += m["content"]
		}
	}
	if token != "pong" {
		t.Fatalf("token=%q events=%v body=%s", token, events, rr.Body.String())
	}
	joined := strings.Join(events, ",")
	if !strings.Contains(joined, "meta") || !strings.Contains(joined, "token") || !strings.Contains(joined, "done") {
		t.Fatalf("events=%v", events)
	}
}

func TestLLMInfo(t *testing.T) {
	os.Setenv("LLM_PROVIDER", "onnx")
	t.Cleanup(func() { os.Unsetenv("LLM_PROVIDER") })
	s := NewServer(nil, NewStore())
	req := httptest.NewRequest(http.MethodGet, "/llm/info", nil)
	rr := httptest.NewRecorder()
	s.Mux().ServeHTTP(rr, req)
	if rr.Code != 200 {
		t.Fatalf("status=%d", rr.Code)
	}
	var out map[string]any
	if err := json.NewDecoder(rr.Body).Decode(&out); err != nil {
		t.Fatal(err)
	}
	if out["provider"] != "onnx" {
		t.Fatalf("provider=%v", out["provider"])
	}
	if out["stream_path"] != "/llm/stream" {
		t.Fatalf("stream_path=%v", out["stream_path"])
	}
	_ = stubModel{text: "x"} // silence unused if compiler is picky in some builds
}
