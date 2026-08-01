package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// TestOpenAIConverseStream parses OpenAI-compatible SSE and fans tokens out.
func TestOpenAIConverseStream(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/chat/completions" {
			http.NotFound(w, r)
			return
		}
		var body map[string]any
		_ = json.NewDecoder(r.Body).Decode(&body)
		if body["stream"] != true {
			t.Errorf("expected stream=true, got %v", body["stream"])
		}
		w.Header().Set("Content-Type", "text/event-stream")
		fl := w.(http.Flusher)
		chunks := []string{"Hello", " ", "Smol", "LM"}
		for _, c := range chunks {
			payload := fmt.Sprintf(
				`data: {"id":"x","object":"chat.completion.chunk","choices":[{"index":0,"delta":{"content":%q}}]}`+"\n\n",
				c,
			)
			_, _ = w.Write([]byte(payload))
			fl.Flush()
		}
		_, _ = w.Write([]byte("data: [DONE]\n\n"))
		fl.Flush()
	}))
	defer ts.Close()

	m := NewOpenAI(ts.URL+"/v1", "test-key", "smollm")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var got []string
	out, err := m.ConverseStream(ctx, "sys", []Message{UserMessage("hi")}, func(tok string) error {
		got = append(got, tok)
		return nil
	})
	if err != nil {
		t.Fatalf("ConverseStream: %v", err)
	}
	if full := strings.Join(got, ""); full != "Hello SmolLM" {
		t.Fatalf("tokens=%q full=%q want Hello SmolLM", got, full)
	}
	if out.Text != "Hello SmolLM" {
		t.Fatalf("out.Text=%q", out.Text)
	}
	if _, ok := AsStreamer(m); !ok {
		t.Fatal("AsStreamer(*OpenAI) should be true")
	}
}
