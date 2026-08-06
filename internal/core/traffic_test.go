package core

import (
	"strings"
	"testing"
)

func TestTrafficStore_RecordAndGet(t *testing.T) {
	ts := NewTrafficStore(t.TempDir(), "test-run")
	rec := ts.Record("GET", "http://example.com/api",
		map[string]string{"Authorization": "Bearer test"}, "",
		200, "200 OK", map[string]string{"Content-Type": "application/json"},
		`{"status":"ok"}`, "http_probe_batch")

	if rec.ID == "" {
		t.Error("expected non-empty ID")
	}
	if rec.Seq != 1 {
		t.Errorf("expected seq 1, got %d", rec.Seq)
	}

	got, ok := ts.Get(rec.ID)
	if !ok {
		t.Fatal("expected to find traffic record")
	}
	if got.URL != "http://example.com/api" {
		t.Errorf("unexpected URL: %s", got.URL)
	}
}

func TestTrafficStore_Search(t *testing.T) {
	ts := NewTrafficStore(t.TempDir(), "test-run")
	ts.Record("GET", "http://example.com/login", nil, "", 200, "OK", nil, "login page", "probe")
	ts.Record("POST", "http://example.com/api/users", nil, "", 201, "Created", nil, "created", "probe")
	ts.Record("GET", "http://example.com/admin", nil, "", 403, "Forbidden", nil, "forbidden", "probe")

	results := ts.Search("admin", 10)
	if len(results) != 1 {
		t.Fatalf("expected 1 result for 'admin', got %d", len(results))
	}

	results = ts.Search("POST", 10)
	if len(results) != 1 {
		t.Fatalf("expected 1 POST result, got %d", len(results))
	}
}

func TestTrafficStore_List(t *testing.T) {
	ts := NewTrafficStore(t.TempDir(), "test-run")
	for i := 0; i < 5; i++ {
		ts.Record("GET", "http://example.com/"+string(rune('a'+i)), nil, "", 200, "OK", nil, "", "probe")
	}

	all := ts.List(3)
	if len(all) != 3 {
		t.Fatalf("expected 3 results, got %d", len(all))
	}
	// Newest first
	if !strings.Contains(all[0].URL, "e") {
		t.Error("expected newest first")
	}
}

func TestTrafficStore_Persist(t *testing.T) {
	dir := t.TempDir()
	ts := NewTrafficStore(dir, "test-run")
	ts.Record("GET", "http://example.com", nil, "", 200, "OK", nil, "hello", "probe")

	if err := ts.Persist(); err != nil {
		t.Fatalf("Persist: %v", err)
	}
}
