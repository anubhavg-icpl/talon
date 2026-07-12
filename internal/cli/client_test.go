package cli

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSanitizeJSONStringsEscapesControlChars(t *testing.T) {
	// Raw newline inside a string value (invalid JSON); structural newlines OK.
	raw := []byte("{\n  \"status\": \"completed\",\n  \"output\": \"line1\nline2\",\n  \"interrupt\": null\n}")
	clean := sanitizeJSONStrings(raw)
	var out StatusResponse
	if err := json.Unmarshal(clean, &out); err != nil {
		t.Fatalf("unmarshal sanitized: %v\nraw=%q\nclean=%q", err, raw, clean)
	}
	if out.Status != "completed" {
		t.Fatalf("status=%q", out.Status)
	}
	if out.Output != "line1\nline2" {
		t.Fatalf("output=%q", out.Output)
	}
}

func TestClientStartAndStatus(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /input/start", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{
			"run_id":  "run-123",
			"message": "Agent execution started",
		})
	})
	mux.HandleFunc("GET /output/status/{run_id}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"status":    "awaiting_approval",
			"output":    "",
			"interrupt": map[string]any{"ToolName": "nmap_scan", "Args": map[string]any{"target": "127.0.0.1"}},
		})
	})
	mux.HandleFunc("POST /output/resume/{run_id}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"message": "Decision received, resuming orchestrator..."})
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	c, err := NewClient(srv.URL, 5*time.Second)
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()

	start, err := c.Start(ctx, StartRequest{IP: "127.0.0.1", CVEID: "CVE-2011-2523"})
	if err != nil {
		t.Fatalf("start: %v", err)
	}
	if start.RunID != "run-123" {
		t.Fatalf("run_id=%q", start.RunID)
	}

	st, err := c.Status(ctx, start.RunID)
	if err != nil {
		t.Fatalf("status: %v", err)
	}
	if st.Status != "awaiting_approval" {
		t.Fatalf("status=%q", st.Status)
	}
	if st.Interrupt == nil || st.Interrupt.ToolName != "nmap_scan" {
		t.Fatalf("interrupt=%+v", st.Interrupt)
	}

	res, err := c.Resume(ctx, start.RunID, ResumeRequest{Decision: "approve"})
	if err != nil {
		t.Fatalf("resume: %v", err)
	}
	if res.Message == "" {
		t.Fatal("expected resume message")
	}
}

func TestExitCodeForStatus(t *testing.T) {
	cases := map[string]int{
		"completed":         ExitOK,
		"awaiting_approval": ExitAwaitingApproval,
		"not_found":         ExitNotFound,
		"error":             ExitError,
		"running":           ExitOK,
	}
	for st, want := range cases {
		if got := ExitCodeForStatus(st); got != want {
			t.Errorf("ExitCodeForStatus(%q)=%d want %d", st, got, want)
		}
	}
}

func TestParseOutputFormat(t *testing.T) {
	for _, s := range []string{"table", "json", "raw", "TABLE"} {
		if _, err := ParseOutputFormat(s); err != nil {
			t.Errorf("ParseOutputFormat(%q): %v", s, err)
		}
	}
	if _, err := ParseOutputFormat("yaml"); err == nil {
		t.Fatal("expected error for yaml")
	}
}
