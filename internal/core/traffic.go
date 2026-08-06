package core

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/anubhavg-icpl/talon/internal/llm"
)

// TrafficRecord stores one HTTP request/response pair for replay and analysis.
// Ported from pentest agent traffic/store.py.
type TrafficRecord struct {
	ID         string            `json:"id"`
	Seq        int               `json:"seq"`
	Method     string            `json:"method"`
	URL        string            `json:"url"`
	ReqHeaders map[string]string `json:"req_headers"`
	ReqBody    string            `json:"req_body,omitempty"`
	StatusCode int               `json:"status_code"`
	Status     string            `json:"status"`
	RespHeaders map[string]string `json:"resp_headers"`
	RespBody   string            `json:"resp_body"`
	Timestamp  time.Time         `json:"timestamp"`
	Source     string            `json:"source"` // which tool generated this
}

// TrafficStore is a per-run append-only store of HTTP traffic.
type TrafficStore struct {
	mu      sync.Mutex
	records []TrafficRecord
	seq     int
	dataDir string
	runID   string
}

// NewTrafficStore creates a new traffic store for a run.
func NewTrafficStore(dataDir, runID string) *TrafficStore {
	if dataDir == "" {
		dataDir = "talon-data/traffic"
	}
	return &TrafficStore{
		dataDir: filepath.Join(dataDir, runID),
		runID:   runID,
	}
}

// Record captures an HTTP request/response pair.
// The ID is a deterministic content+sequence hash (resume-safe).
func (ts *TrafficStore) Record(method, url string, reqHeaders map[string]string, reqBody string,
	statusCode int, status string, respHeaders map[string]string, respBody string, source string) TrafficRecord {

	ts.mu.Lock()
	defer ts.mu.Unlock()

	ts.seq++
	// Deterministic hash: seq + method + url + reqBody hash
	h := sha256.New()
	fmt.Fprintf(h, "%d|%s|%s|%s", ts.seq, method, url, reqBody)
	contentHash := hex.EncodeToString(h.Sum(nil))[:12]
	id := "t" + contentHash

	rec := TrafficRecord{
		ID:          id,
		Seq:         ts.seq,
		Method:      method,
		URL:         url,
		ReqHeaders:  reqHeaders,
		ReqBody:     reqBody,
		StatusCode:  statusCode,
		Status:      status,
		RespHeaders: respHeaders,
		RespBody:    truncateStr(respBody, 50000), // Cap stored body
		Timestamp:   time.Now().UTC(),
		Source:      source,
	}
	ts.records = append(ts.records, rec)
	return rec
}

// Get retrieves a traffic record by ID.
func (ts *TrafficStore) Get(id string) (TrafficRecord, bool) {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	for _, r := range ts.records {
		if r.ID == id {
			return r, true
		}
	}
	return TrafficRecord{}, false
}

// List returns all traffic records (newest first).
func (ts *TrafficStore) List(limit int) []TrafficRecord {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	if limit <= 0 || limit > len(ts.records) {
		limit = len(ts.records)
	}
	out := make([]TrafficRecord, limit)
	for i, j := 0, len(ts.records)-1; i < limit; i, j = i+1, j-1 {
		out[i] = ts.records[j]
	}
	return out
}

// Search searches traffic records by URL, method, or response body.
func (ts *TrafficStore) Search(query string, limit int) []TrafficRecord {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	if limit <= 0 {
		limit = 10
	}
	query = strings.ToLower(query)
	var out []TrafficRecord
	for i := len(ts.records) - 1; i >= 0 && len(out) < limit; i-- {
		r := ts.records[i]
		if strings.Contains(strings.ToLower(r.URL), query) ||
			strings.Contains(strings.ToLower(r.Method), query) ||
			strings.Contains(strings.ToLower(r.RespBody), query) {
			out = append(out, r)
		}
	}
	return out
}

// Count returns total stored traffic records.
func (ts *TrafficStore) Count() int {
	ts.mu.Lock()
	defer ts.mu.Unlock()
	return len(ts.records)
}

// Persist writes the traffic index to disk as JSONL.
func (ts *TrafficStore) Persist() error {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	if err := os.MkdirAll(ts.dataDir, 0755); err != nil {
		return err
	}

	path := filepath.Join(ts.dataDir, "requests.jsonl")
	var sb strings.Builder
	for _, r := range ts.records {
		data, err := json.Marshal(r)
		if err != nil {
			continue
		}
		sb.Write(data)
		sb.WriteString("\n")
	}

	tmpPath := path + ".tmp"
	if err := os.WriteFile(tmpPath, []byte(sb.String()), 0644); err != nil {
		return err
	}
	return os.Rename(tmpPath, path)
}

// --- Traffic agent tools ---

// trafficToolSpecs returns the agent-facing traffic inspection tools.
func trafficToolSpecs() []llm.ToolSpec {
	return []llm.ToolSpec{
		{
			Name:        "traffic_list",
			Description: "List stored HTTP traffic from previous requests. Each entry has an id, method, URL, status, and size.",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"limit": map[string]any{"type": "integer", "description": "Max results (default 20)"},
				},
			},
		},
		{
			Name:        "traffic_view",
			Description: "View full request/response details of a stored traffic record by id.",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"id": map[string]any{"type": "string", "description": "Traffic record id"},
				},
				"required": []any{"id"},
			},
		},
		{
			Name:        "traffic_search",
			Description: "Search stored HTTP traffic by URL, method, or response body content.",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"q":     map[string]any{"type": "string"},
					"limit": map[string]any{"type": "integer"},
				},
				"required": []any{"q"},
			},
		},
	}
}

// handleTrafficList handles the traffic_list tool call.
func handleTrafficList(store *TrafficStore, args map[string]any, tr *tracker) (string, bool) {
	if store == nil {
		return "traffic store not configured", true
	}
	limit := 20
	if v, ok := args["limit"].(float64); ok {
		limit = int(v)
	}
	records := store.List(limit)

	type item struct {
		ID        string `json:"id"`
		Method    string `json:"method"`
		URL       string `json:"url"`
		StatusCode int   `json:"status_code"`
		Size      int    `json:"size"`
	}
	items := make([]item, 0, len(records))
	for _, r := range records {
		items = append(items, item{
			ID: r.ID, Method: r.Method, URL: r.URL,
			StatusCode: r.StatusCode, Size: len(r.RespBody),
		})
	}
	payload := map[string]any{
		"total":   store.Count(),
		"showing": len(items),
		"items":   items,
	}
	raw, _ := json.MarshalIndent(payload, "", "  ")
	if tr != nil {
		tr.record("traffic_list", args, fmt.Sprintf("total=%d", store.Count()))
	}
	return string(raw), false
}

// handleTrafficView handles the traffic_view tool call.
func handleTrafficView(store *TrafficStore, args map[string]any, tr *tracker) (string, bool) {
	if store == nil {
		return "traffic store not configured", true
	}
	id, _ := args["id"].(string)
	rec, ok := store.Get(id)
	if !ok {
		return fmt.Sprintf("traffic %q not found", id), true
	}
	payload := map[string]any{
		"id":          rec.ID,
		"method":      rec.Method,
		"url":         rec.URL,
		"req_headers": rec.ReqHeaders,
		"req_body":    rec.ReqBody,
		"status_code": rec.StatusCode,
		"status":      rec.Status,
		"resp_headers": rec.RespHeaders,
		"resp_body":   rec.RespBody,
		"timestamp":   rec.Timestamp,
		"source":      rec.Source,
	}
	raw, _ := json.MarshalIndent(payload, "", "  ")
	if tr != nil {
		tr.record("traffic_view", args, fmt.Sprintf("id=%s", rec.ID))
	}
	return string(raw), false
}

// handleTrafficSearch handles the traffic_search tool call.
func handleTrafficSearch(store *TrafficStore, args map[string]any, tr *tracker) (string, bool) {
	if store == nil {
		return "traffic store not configured", true
	}
	q, _ := args["q"].(string)
	limit := 10
	if v, ok := args["limit"].(float64); ok {
		limit = int(v)
	}
	results := store.Search(q, limit)

	type item struct {
		ID        string `json:"id"`
		Method    string `json:"method"`
		URL       string `json:"url"`
		StatusCode int   `json:"status_code"`
	}
	items := make([]item, 0, len(results))
	for _, r := range results {
		items = append(items, item{
			ID: r.ID, Method: r.Method, URL: r.URL, StatusCode: r.StatusCode,
		})
	}
	payload := map[string]any{
		"query":   q,
		"total":   len(results),
		"items":   items,
	}
	raw, _ := json.MarshalIndent(payload, "", "  ")
	if tr != nil {
		tr.record("traffic_search", args, fmt.Sprintf("q=%q hits=%d", q, len(results)))
	}
	return string(raw), false
}
