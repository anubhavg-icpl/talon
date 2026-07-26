package control

import (
	"bufio"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/anubhavg-icpl/talon/internal/core"
	"github.com/anubhavg-icpl/talon/internal/llm"
	"github.com/coder/websocket"
)

func TestListRunsNewestFirst(t *testing.T) {
	t.Parallel()
	store := NewStore()
	store.Create("r-old", core.RunInput{SessionID: "r-old", TargetIP: "1.1.1.1"})
	store.Create("r-new", core.RunInput{SessionID: "r-new", TargetIP: "2.2.2.2", CVEID: "CVE-2011-2523"})
	store.SetStatus("r-new", "running")
	store.SetResult("r-new", core.RunResult{FinalMessage: "done", JudgeVerdict: true})

	runs := store.List()
	if len(runs) != 2 {
		t.Fatalf("expected 2 runs, got %d", len(runs))
	}
	if runs[0].RunID != "r-new" {
		t.Fatalf("expected newest first, got %s", runs[0].RunID)
	}
	if runs[0].JudgeVerdict == nil || !*runs[0].JudgeVerdict {
		t.Fatal("expected judge verdict true on completed run")
	}
	if runs[1].JudgeVerdict != nil {
		t.Fatal("expected nil judge verdict on run without result")
	}
	if runs[1].Target != "1.1.1.1" {
		t.Fatalf("target=%q want 1.1.1.1", runs[1].Target)
	}
}

func TestHandleListRuns(t *testing.T) {
	t.Parallel()
	store := NewStore()
	store.Create("r1", core.RunInput{SessionID: "r1", TargetIP: "10.0.0.5"})
	srv := NewServer(nil, store)

	req := httptest.NewRequest(http.MethodGet, "/runs", nil)
	rec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d want 200", rec.Code)
	}
	var body struct {
		Runs []RunSummary `json:"runs"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(body.Runs) != 1 || body.Runs[0].RunID != "r1" || body.Runs[0].Target != "10.0.0.5" {
		t.Fatalf("unexpected body: %+v", body.Runs)
	}
}

func TestServiceHealthShape(t *testing.T) {
	t.Parallel()
	srv := NewServer(nil, NewStore())
	req := httptest.NewRequest(http.MethodGet, "/health/services", nil)
	rec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d want 200", rec.Code)
	}
	var body struct {
		Services []ServiceHealth `json:"services"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(body.Services) != 7 {
		t.Fatalf("expected 7 services, got %d", len(body.Services))
	}
	names := map[string]string{}
	for _, svc := range body.Services {
		names[svc.Name] = svc.Status
	}
	if names["talon-core"] != "online" {
		t.Fatalf("talon-core status=%q want online", names["talon-core"])
	}
	if names["postgres"] != "unconfigured" {
		t.Fatalf("postgres status=%q want unconfigured (no DB in test)", names["postgres"])
	}
	if names["redis"] != "unconfigured" {
		t.Fatalf("redis status=%q want unconfigured (no cache in test)", names["redis"])
	}
	for _, want := range []string{"arsenal-engine", "msfrpcd", "rabbitmq", "ollama"} {
		if _, ok := names[want]; !ok {
			t.Fatalf("missing service %s", want)
		}
	}
}

func TestCORSHeaders(t *testing.T) {
	t.Parallel()
	srv := NewServer(nil, NewStore())

	req := httptest.NewRequest(http.MethodOptions, "/runs", nil)
	rec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("OPTIONS status=%d want 204", rec.Code)
	}
	if rec.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Fatal("missing CORS origin header")
	}
}

func TestWebSocketTerminalRun(t *testing.T) {
	t.Parallel()
	store := NewStore()
	store.Create("r1", core.RunInput{SessionID: "r1", TargetIP: "10.0.0.5"})
	store.SetResult("r1", core.RunResult{
		FinalMessage: "rooted",
		JudgeVerdict: true,
		ToolLog: []core.ToolCallRecord{
			{Index: 0, ToolName: "nmap_scan", Args: map[string]any{"target": "10.0.0.5"}},
		},
	})
	srv := NewServer(nil, store)
	ts := httptest.NewServer(srv.Handler())
	defer ts.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	conn, _, err := websocket.Dial(ctx, "ws"+strings.TrimPrefix(ts.URL, "http")+"/monitor/ws/r1", nil)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close(websocket.StatusNormalClosure, "")

	var types []string
	for {
		_, data, err := conn.Read(ctx)
		if err != nil {
			t.Fatalf("read: %v", err)
		}
		var msg wsMessage
		if err := json.Unmarshal(data, &msg); err != nil {
			t.Fatalf("decode: %v", err)
		}
		types = append(types, msg.Type)
		if msg.Type == "status" {
			var body struct {
				Status string `json:"status"`
			}
			raw, _ := json.Marshal(msg.Data)
			_ = json.Unmarshal(raw, &body)
			if terminalStatus(body.Status) {
				break
			}
		}
	}
	if len(types) != 2 || types[0] != "tool" || types[1] != "status" {
		t.Fatalf("messages=%v want [tool status]", types)
	}
}

func TestStreamTerminalRun(t *testing.T) {
	t.Parallel()
	store := NewStore()
	store.Create("r1", core.RunInput{SessionID: "r1", TargetIP: "10.0.0.5"})
	store.SetResult("r1", core.RunResult{
		FinalMessage: "rooted",
		JudgeVerdict: true,
		ToolLog: []core.ToolCallRecord{
			{Index: 0, ToolName: "nmap_scan", Args: map[string]any{"target": "10.0.0.5"}},
		},
	})
	srv := NewServer(nil, store)
	ts := httptest.NewServer(srv.Handler())
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/monitor/stream/r1")
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer resp.Body.Close()
	if ct := resp.Header.Get("Content-Type"); ct != "text/event-stream" {
		t.Fatalf("content-type=%q", ct)
	}

	var events []string
	sc := bufio.NewScanner(resp.Body)
	for sc.Scan() {
		line := sc.Text()
		if strings.HasPrefix(line, "event: ") {
			events = append(events, strings.TrimPrefix(line, "event: "))
		}
	}
	if err := sc.Err(); err != nil {
		t.Fatalf("scan: %v", err)
	}
	// Terminal run: stream must deliver the tool + terminal status, then close.
	if len(events) != 2 || events[0] != "tool" || events[1] != "status" {
		t.Fatalf("events=%v want [tool status]", events)
	}
}

type fakeAnalyzer struct {
	text string
	err  error
}

func (f fakeAnalyzer) Converse(_ context.Context, _ string, _ []llm.Message, _ []llm.ToolSpec) (llm.Message, error) {
	if f.err != nil {
		return llm.Message{}, f.err
	}
	return llm.AssistantText(f.text), nil
}

func TestAnalyzeRun(t *testing.T) {
	t.Parallel()
	store := NewStore()
	store.Create("r1", core.RunInput{SessionID: "r1", TargetIP: "10.0.0.5", CVEID: "CVE-2011-2523"})
	store.SetResult("r1", core.RunResult{
		FinalMessage: "vsftpd backdoor confirmed",
		JudgeVerdict: true,
		ToolLog:      []core.ToolCallRecord{{Index: 0, ToolName: "nmap_scan", Output: "21/tcp open vsftpd 2.3.4"}},
	})
	srv := NewServer(nil, store, WithAnalyzer(fakeAnalyzer{text: "EXECUTIVE SUMMARY\nrooted."}))

	req := httptest.NewRequest(http.MethodPost, "/analyze/r1", nil)
	rec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d want 200 (%s)", rec.Code, rec.Body.String())
	}
	var body struct {
		RunID    string `json:"run_id"`
		Analysis string `json:"analysis"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if body.RunID != "r1" || !strings.Contains(body.Analysis, "EXECUTIVE SUMMARY") {
		t.Fatalf("unexpected body: %+v", body)
	}
}

type fakeCache struct{ m map[string]string }

func (f fakeCache) Get(_ context.Context, key string) (string, bool) {
	v, ok := f.m[key]
	return v, ok
}
func (f fakeCache) Set(_ context.Context, key, val string, _ time.Duration) { f.m[key] = val }
func (f fakeCache) Del(_ context.Context, key string)                         { delete(f.m, key) }
func (f fakeCache) Ping(_ context.Context) error                              { return nil }

func TestAnalyzeCaching(t *testing.T) {
	t.Parallel()
	cache := fakeCache{m: map[string]string{}}
	store := NewStore()
	store.Create("r1", core.RunInput{SessionID: "r1", TargetIP: "10.0.0.5"})
	store.SetResult("r1", core.RunResult{FinalMessage: "done", JudgeVerdict: true})
	srv := NewServer(nil, store, WithAnalyzer(fakeAnalyzer{text: "BRIEF"}), WithCache(cache))

	// First call: miss → LLM → cached.
	req := httptest.NewRequest(http.MethodPost, "/analyze/r1", nil)
	rec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK || rec.Header().Get("X-Cache") != "miss" {
		t.Fatalf("first call: code=%d x-cache=%q", rec.Code, rec.Header().Get("X-Cache"))
	}
	if _, ok := cache.m[cacheKeyAnalyzePrefix+"r1"]; !ok {
		t.Fatal("analysis not cached")
	}

	// Second call: hit — served from cache (analyzer would still answer, but
	// the header proves the cache path).
	req = httptest.NewRequest(http.MethodPost, "/analyze/r1", nil)
	rec = httptest.NewRecorder()
	srv.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK || rec.Header().Get("X-Cache") != "hit" {
		t.Fatalf("second call: code=%d x-cache=%q", rec.Code, rec.Header().Get("X-Cache"))
	}
}

func TestConfigEndpointsWithoutDB(t *testing.T) {
	t.Setenv("FEATURE_AI_ANALYSIS", "true")
	srv := NewServer(nil, NewStore(), WithSettings(NewSettings(nil)))

	// GET works env-only, keys are present and masked flags correct.
	req := httptest.NewRequest(http.MethodGet, "/config", nil)
	rec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /config status=%d", rec.Code)
	}
	var body struct {
		Config []ConfigEntry `json:"config"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode: %v", err)
	}
	byKey := map[string]ConfigEntry{}
	for _, e := range body.Config {
		byKey[e.Key] = e
	}
	if !byKey["OPENAI_API_KEY"].Secret {
		t.Fatal("OPENAI_API_KEY must be marked secret")
	}
	if byKey["OPENAI_API_KEY"].Writable {
		t.Fatal("writable must be false without postgres")
	}
	if _, ok := byKey["FEATURE_AI_ANALYSIS"]; !ok {
		t.Fatal("missing FEATURE_AI_ANALYSIS entry")
	}

	// PUT without postgres → 503.
	req = httptest.NewRequest(http.MethodPut, "/config", strings.NewReader(`{"LLM_PROVIDER":"ollama"}`))
	rec = httptest.NewRecorder()
	srv.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("PUT /config status=%d want 503", rec.Code)
	}
}

func TestAnalyzeFeatureFlagDisabled(t *testing.T) {
	t.Setenv("FEATURE_AI_ANALYSIS", "false")
	store := NewStore()
	store.Create("r1", core.RunInput{SessionID: "r1", TargetIP: "10.0.0.5"})
	srv := NewServer(nil, store, WithAnalyzer(fakeAnalyzer{text: "x"}), WithSettings(NewSettings(nil)))

	req := httptest.NewRequest(http.MethodPost, "/analyze/r1", nil)
	rec := httptest.NewRecorder()
	srv.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusForbidden {
		t.Fatalf("disabled analysis status=%d want 403", rec.Code)
	}
}

func TestAnalyzeRunErrors(t *testing.T) {
	t.Parallel()
	store := NewStore()
	store.Create("r1", core.RunInput{SessionID: "r1", TargetIP: "10.0.0.5"})

	// No analyzer configured → 503.
	noModel := NewServer(nil, store)
	req := httptest.NewRequest(http.MethodPost, "/analyze/r1", nil)
	rec := httptest.NewRecorder()
	noModel.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("no-analyzer status=%d want 503", rec.Code)
	}

	// Unknown run → 404.
	srv := NewServer(nil, store, WithAnalyzer(fakeAnalyzer{text: "x"}))
	req = httptest.NewRequest(http.MethodPost, "/analyze/nope", nil)
	rec = httptest.NewRecorder()
	srv.Handler().ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("unknown-run status=%d want 404", rec.Code)
	}
}

func TestPersistenceRoundTrip(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()

	store := NewStore()
	if err := store.EnablePersistence(dir); err != nil {
		t.Fatalf("enable: %v", err)
	}
	store.Create("r1", core.RunInput{SessionID: "r1", TargetIP: "10.0.0.9", CVEID: "CVE-2011-2523"})
	store.SetResult("r1", core.RunResult{FinalMessage: "done", JudgeVerdict: true})
	store.Create("r2", core.RunInput{SessionID: "r2", TargetIP: "10.0.0.10"})
	store.SetStatus("r2", "running")

	fresh := NewStore()
	if err := fresh.EnablePersistence(dir); err != nil {
		t.Fatalf("reload: %v", err)
	}
	runs := fresh.List()
	if len(runs) != 2 {
		t.Fatalf("expected 2 persisted runs, got %d", len(runs))
	}
	byID := map[string]RunSummary{}
	for _, r := range runs {
		byID[r.RunID] = r
	}
	if byID["r1"].Status != "completed" || byID["r1"].CVEID != "CVE-2011-2523" {
		t.Fatalf("r1 after reload: %+v", byID["r1"])
	}
	// Non-terminal runs must be converted to error on reload.
	if byID["r2"].Status != "error" {
		t.Fatalf("r2 status=%q want error", byID["r2"].Status)
	}
}
