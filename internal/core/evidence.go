package core

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"
)

// EvidenceRecord stores raw tool output plus a bounded prompt-facing preview.
// Ported from pentest agent agent/agent_state.py EvidenceRecord + make_evidence_preview.
type EvidenceRecord struct {
	ID          string    `json:"id"`
	Tool        string    `json:"tool"`
	KeyArgs     string    `json:"key_args"`
	Status      int       `json:"status"`
	Summary     string    `json:"summary"`
	Preview     string    `json:"preview"`
	Content     string    `json:"-"`
	Truncated   bool      `json:"truncated"`
	ContentHash string    `json:"content_hash"`
	DuplicateOf string    `json:"duplicate_of,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

func (e EvidenceRecord) Size() int { return len(e.Content) }

// EvidenceStore is a run-scoped append-only store of tool evidence.
// It deduplicates identical outputs via content hash, caps total records,
// and builds bounded previews so the model doesn't carry full blobs.
type EvidenceStore struct {
	mu       sync.Mutex
	records  []EvidenceRecord
	hashIdx  map[string]string // content_hash -> evidence_id
	seq      int
	maxStore int
	preview  int
}

const (
	defaultMaxEvidence   = 240
	defaultEvidencePreview = 6000
)

func NewEvidenceStore() *EvidenceStore {
	return &EvidenceStore{
		hashIdx:  make(map[string]string),
		maxStore: defaultMaxEvidence,
		preview:  defaultEvidencePreview,
	}
}

// Record adds a tool output as evidence. If the content hash matches an
// existing record, the new record references the original via DuplicateOf
// instead of storing a second copy.
func (s *EvidenceStore) Record(tool string, keyArgs, content string, status int) EvidenceRecord {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.seq++
	id := fmt.Sprintf("e%03d", s.seq)
	hash := contentHash(content)

	// Dedup: if identical content already stored, link to it
	if existingID, ok := s.hashIdx[hash]; ok {
		rec := EvidenceRecord{
			ID:          id,
			Tool:        tool,
			KeyArgs:     truncateStr(keyArgs, 200),
			Status:      status,
			Summary:     oneLine(content, 200),
			ContentHash: hash,
			DuplicateOf: existingID,
			CreatedAt:   time.Now().UTC(),
		}
		s.appendCapped(rec)
		return rec
	}

	s.hashIdx[hash] = id
	rec := EvidenceRecord{
		ID:          id,
		Tool:        tool,
		KeyArgs:     truncateStr(keyArgs, 200),
		Status:      status,
		Summary:     oneLine(content, 200),
		Preview:     makeEvidencePreview(content, s.preview),
		Content:     content,
		ContentHash: hash,
		Truncated:   len(content) > s.preview,
		CreatedAt:   time.Now().UTC(),
	}
	s.appendCapped(rec)
	return rec
}

func (s *EvidenceStore) appendCapped(rec EvidenceRecord) {
	s.records = append(s.records, rec)
	// Enforce cap: drop oldest records
	if len(s.records) > s.maxStore {
		s.records = s.records[len(s.records)-s.maxStore:]
		// Rebuild hashIdx so DuplicateOf links don't dangle after eviction
		s.hashIdx = make(map[string]string, len(s.records))
		for _, r := range s.records {
			if r.ContentHash != "" && r.DuplicateOf == "" {
				s.hashIdx[r.ContentHash] = r.ID
			}
		}
	}
}

// Get retrieves a single evidence record by ID.
func (s *EvidenceStore) Get(id string) (EvidenceRecord, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, r := range s.records {
		if r.ID == id {
			return r, true
		}
	}
	return EvidenceRecord{}, false
}

// List returns all evidence records (newest first).
func (s *EvidenceStore) List() []EvidenceRecord {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]EvidenceRecord, len(s.records))
	copy(out, s.records)
	// Reverse for newest-first
	for i, j := 0, len(out)-1; i < j; i, j = i+1, j-1 {
		out[i], out[j] = out[j], out[i]
	}
	return out
}

// Search performs a case-insensitive substring search over evidence content.
func (s *EvidenceStore) Search(query string, limit int) []EvidenceRecord {
	s.mu.Lock()
	defer s.mu.Unlock()
	if limit <= 0 {
		limit = 10
	}
	query = strings.ToLower(query)
	var out []EvidenceRecord
	for i := len(s.records) - 1; i >= 0 && len(out) < limit; i-- {
		r := s.records[i]
		if strings.Contains(strings.ToLower(r.Content), query) ||
			strings.Contains(strings.ToLower(r.Summary), query) {
			out = append(out, r)
		}
	}
	return out
}

// Count returns the total number of stored evidence records.
func (s *EvidenceStore) Count() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.records)
}

// --- Preview builder ---

var (
	flagRe   = regexp.MustCompile(`[A-Za-z_][A-Za-z0-9_]{1,20}\{[^{}\n]{1,200}\}`)
	statusRe = regexp.MustCompile(`(?:Status|HTTP/\d(?:\.\d)?)\s*:?\s*(\d{3})`)
)

// importantLines picks lines likely useful for a preview.
func importantLines(text string, limit int) []string {
	markers := []string{
		"flag", "ctf{", "status:", "headers:", "location:",
		"set-cookie", "error", "exception", "warning", "failed",
		"unauthorized", "forbidden", "not found", "root:", "uid=",
		"password", "secret", "token", "key=", "admin",
	}
	var out []string
	for _, line := range strings.Split(text, "\n") {
		lower := strings.ToLower(line)
		for _, m := range markers {
			if strings.Contains(lower, m) {
				out = append(out, strings.TrimSpace(line))
				break
			}
		}
		if len(out) >= limit {
			break
		}
	}
	return out
}

// makeEvidencePreview builds a bounded prompt-facing preview that preserves
// high-signal lines while omitting large raw blobs.
func makeEvidencePreview(content string, limit int) string {
	if limit <= 0 || len(content) <= limit {
		return content
	}
	lines := importantLines(content, 36)
	header := fmt.Sprintf(
		"[active-context high-signal preview]\n"+
			"raw_size=%d chars; full raw evidence is stored separately; "+
			"use evidence_view/evidence_search if a missing byte or wider span matters.",
		len(content))

	signalBlock := "[signal lines]\n" + strings.Join(lines, "\n")
	remaining := limit - len(header) - len(signalBlock) - 80
	if remaining < 1000 {
		remaining = 1000
	}

	preview := header + "\n\n" + signalBlock + "\n\n[head/tail context]\n" +
		clipText(content, remaining, "...[raw omitted from active context; inspect saved evidence for full body]...")

	return clipText(preview, limit, "...[preview clipped; inspect saved evidence]...")
}

// clipText returns a deterministic head/tail preview for oversized text.
func clipText(text string, limit int, marker string) string {
	if limit <= 0 || len(text) <= limit {
		return text
	}
	head := limit / 2
	if head < 1 {
		head = 1
	}
	tail := limit - head - len(marker) - 2
	if tail < 1 {
		tail = 1
	}
	return text[:head] + "\n" + marker + "\n" + text[len(text)-tail:]
}

func truncateStr(s string, n int) string {
	runes := []rune(s)
	if len(runes) <= n {
		return s
	}
	return string(runes[:n])
}

func oneLine(s string, limit int) string {
	s = strings.Join(strings.Fields(s), " ")
	return truncateStr(s, limit)
}

func contentHash(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:16])
}

// ParseStatusCode extracts an HTTP status from tool output.
func ParseStatusCode(output string) int {
	m := statusRe.FindStringSubmatch(output)
	if len(m) < 2 {
		return 0
	}
	var code int
	fmt.Sscanf(m[1], "%d", &code)
	return code
}

// ExtractFlags extracts CTF flag-looking tokens from text.
func ExtractFlags(text string) []string {
	matches := flagRe.FindAllString(text, -1)
	seen := make(map[string]bool)
	var out []string
	for _, m := range matches {
		if !seen[m] {
			seen[m] = true
			out = append(out, m)
		}
	}
	return out
}
