package core

import (
	"strings"
	"testing"
)

func TestEvidenceStore_RecordAndGet(t *testing.T) {
	s := NewEvidenceStore()
	rec := s.Record("nmap_scan", "target=10.0.0.1", "Host is up. Port 22 open.", 200)

	if rec.ID != "e001" {
		t.Errorf("expected ID e001, got %s", rec.ID)
	}
	if rec.Tool != "nmap_scan" {
		t.Errorf("expected tool nmap_scan, got %s", rec.Tool)
	}
	if rec.ContentHash == "" {
		t.Error("expected non-empty content hash")
	}

	got, ok := s.Get("e001")
	if !ok {
		t.Fatal("expected to find evidence e001")
	}
	if got.Content != "Host is up. Port 22 open." {
		t.Errorf("unexpected content: %s", got.Content)
	}
}

func TestEvidenceStore_Dedup(t *testing.T) {
	s := NewEvidenceStore()
	r1 := s.Record("fetch", "url=http://x.com", "same content here", 200)
	r2 := s.Record("fetch", "url=http://y.com", "same content here", 200)

	if r1.DuplicateOf != "" {
		t.Error("first record should not be a duplicate")
	}
	if r2.DuplicateOf != r1.ID {
		t.Errorf("expected duplicate_of=%s, got %s", r1.ID, r2.DuplicateOf)
	}
	if s.Count() != 2 {
		t.Errorf("expected 2 records, got %d", s.Count())
	}
}

func TestEvidenceStore_Search(t *testing.T) {
	s := NewEvidenceStore()
	s.Record("fetch", "", "The flag{test_flag} was found here", 200)
	s.Record("nmap", "", "Port 8080 open", 200)
	s.Record("curl", "", "HTTP/1.1 200 OK", 200)

	results := s.Search("flag", 10)
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if !strings.Contains(results[0].Content, "flag{test_flag}") {
		t.Error("expected flag in result content")
	}

	results = s.Search("200", 10)
	if len(results) != 1 {
		t.Fatalf("expected 1 result for '200', got %d", len(results))
	}
}

func TestEvidenceStore_Cap(t *testing.T) {
	s := NewEvidenceStore()
	s.maxStore = 5
	for i := 0; i < 10; i++ {
		s.Record("tool", "", strings.Repeat("x", i*100+1), 200)
	}
	if s.Count() != 5 {
		t.Errorf("expected 5 records after cap, got %d", s.Count())
	}
}

func TestMakeEvidencePreview_Small(t *testing.T) {
	small := "short content"
	p := makeEvidencePreview(small, 6000)
	if p != small {
		t.Error("small content should pass through unchanged")
	}
}

func TestMakeEvidencePreview_Large(t *testing.T) {
	large := strings.Repeat("a", 10000)
	large += "\nflag{hidden_flag}\n"
	large += strings.Repeat("b", 5000)

	p := makeEvidencePreview(large, 1000)
	if len(p) > 1200 {
		t.Errorf("preview too large: %d", len(p))
	}
	if !strings.Contains(p, "flag{hidden_flag}") {
		t.Error("preview should contain signal line with flag")
	}
	if !strings.Contains(p, "active-context") {
		t.Error("preview should contain header")
	}
}

func TestParseStatusCode(t *testing.T) {
	tests := []struct {
		input string
		want  int
	}{
		{"HTTP/1.1 200 OK", 200},
		{"Status: 404 Not Found", 404},
		{"no status here", 0},
		{"HTTP/2 500 Internal Server Error", 500},
	}
	for _, tt := range tests {
		got := ParseStatusCode(tt.input)
		if got != tt.want {
			t.Errorf("ParseStatusCode(%q) = %d, want %d", tt.input, got, tt.want)
		}
	}
}

func TestExtractFlags(t *testing.T) {
	text := "some text flag{abc123} more text CTF{xyz} end"
	flags := ExtractFlags(text)
	if len(flags) != 2 {
		t.Fatalf("expected 2 flags, got %d", len(flags))
	}
	if flags[0] != "flag{abc123}" {
		t.Errorf("expected flag{abc123}, got %s", flags[0])
	}
}
