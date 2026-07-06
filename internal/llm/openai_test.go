package llm

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"
)

// TestOpenAIConverseLive exercises the real OpenAI.Converse code path against
// a live OpenAI-compatible endpoint (GLM-5.2 via z.ai by default). It is an
// integration test, gated on OPENAI_API_KEY -- skipped, not failed, when the
// key is absent so `go test ./...` stays green offline.
//
//	Run it with:  OPENAI_API_KEY=... go test ./internal/llm/ -run Live -v
func TestOpenAIConverseLive(t *testing.T) {
	key := os.Getenv("OPENAI_API_KEY")
	if key == "" {
		t.Skip("OPENAI_API_KEY not set; skipping live OpenAI provider test")
	}
	base := getenvOr("OPENAI_BASE_URL", "https://api.z.ai/api/coding/paas/v4")
	model := getenvOr("OPENAI_MAIN_MODEL", "glm-5.2")

	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	// 1. Plain chat round-trip.
	m := NewOpenAI(base, key, model)
	out, err := m.Converse(ctx,
		"You are a test echo. Reply with exactly: PONG",
		[]Message{UserMessage("ping")},
		nil)
	if err != nil {
		t.Fatalf("plain Converse: %v", err)
	}
	if !strings.Contains(strings.ToUpper(out.Text), "PONG") {
		t.Fatalf("plain reply missing PONG: %q", out.Text)
	}

	// 2. Function-calling round-trip -- the path the orchestrator loop depends
	// on. Verifies the provider marshals the tool spec, the model emits a
	// tool_call, and our decoder parses the JSON-string arguments back to a
	// map (the GLM/OpenAI wire quirk that breaks naive parsers).
	tools := []ToolSpec{{
		Name:        "nmap_scan",
		Description: "Run an nmap scan against a single target IP",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"target": map[string]any{"type": "string"},
			},
			"required": []string{"target"},
		},
	}}
	out, err = m.Converse(ctx,
		"You are an authorized pentest recon agent. Always call the nmap_scan tool for the given target.",
		[]Message{UserMessage("Scan 192.168.122.250 -- you MUST call nmap_scan.")},
		tools)
	if err != nil {
		t.Fatalf("tool Converse: %v", err)
	}
	if len(out.ToolCalls) == 0 {
		t.Fatalf("expected a tool call, got text %q", out.Text)
	}
	tc := out.ToolCalls[0]
	if tc.Name != "nmap_scan" {
		t.Fatalf("tool call name = %q, want nmap_scan", tc.Name)
	}
	if tc.Args["target"] != "192.168.122.250" {
		t.Fatalf("tool call args = %#v, want target=192.168.122.250", tc.Args)
	}
	if tc.ID == "" {
		t.Fatalf("tool call has no ID; orchestrator bookkeeping needs one")
	}
	t.Logf("OK: tool_call id=%s args=%v", tc.ID, tc.Args)
}

func getenvOr(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
