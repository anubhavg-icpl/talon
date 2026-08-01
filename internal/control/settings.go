package control

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

// configKey describes one dashboard-manageable setting.
type configKey struct {
	Label  string
	Secret bool   // masked on read, write-only
	Hot    bool   // applies without restart (feature flags)
	EnvVar string // env fallback
}

// configKeys is the whitelist for GET/PUT /config. Anything not listed is
// rejected — the dashboard can never mutate arbitrary env.
var configKeys = map[string]configKey{
	"LLM_PROVIDER":        {Label: "LLM provider", EnvVar: "LLM_PROVIDER"},
	"OPENAI_BASE_URL":     {Label: "OpenAI-compatible base URL", EnvVar: "OPENAI_BASE_URL"},
	"OPENAI_API_KEY":      {Label: "OpenAI API key", Secret: true, EnvVar: "OPENAI_API_KEY"},
	"OPENAI_MAIN_MODEL":   {Label: "OpenAI main model", EnvVar: "OPENAI_MAIN_MODEL"},
	"OPENAI_JUDGE_MODEL":  {Label: "OpenAI judge model", EnvVar: "OPENAI_JUDGE_MODEL"},
	"OPENAI_CODE_MODEL":   {Label: "OpenAI code model", EnvVar: "OPENAI_CODE_MODEL"},
	"OLLAMA_URL":          {Label: "Ollama URL", EnvVar: "OLLAMA_URL"},
	"OLLAMA_MAIN_MODEL":   {Label: "Ollama main model", EnvVar: "OLLAMA_MAIN_MODEL"},
	"ONNX_BASE_URL":       {Label: "ONNX/SmolLM base URL", EnvVar: "ONNX_BASE_URL"},
	"ONNX_API_KEY":        {Label: "ONNX/SmolLM API key", Secret: true, EnvVar: "ONNX_API_KEY"},
	"ONNX_MAIN_MODEL":     {Label: "ONNX/SmolLM main model", EnvVar: "ONNX_MAIN_MODEL"},
	"ONNX_JUDGE_MODEL":    {Label: "ONNX/SmolLM judge model", EnvVar: "ONNX_JUDGE_MODEL"},
	"ONNX_CODE_MODEL":     {Label: "ONNX/SmolLM code model", EnvVar: "ONNX_CODE_MODEL"},
	"AGENT_MODEL_ID":      {Label: "Bedrock agent model", EnvVar: "AGENT_MODEL_ID"},
	"JUDGE_MODEL_ID":      {Label: "Bedrock judge model", EnvVar: "JUDGE_MODEL_ID"},
	"CODE_MODEL_ID":       {Label: "Bedrock code model", EnvVar: "CODE_MODEL_ID"},
	"LHOST":               {Label: "Reverse-shell LHOST", EnvVar: "LHOST"},
	"LPORT":               {Label: "Reverse-shell LPORT", EnvVar: "LPORT"},
	"FEATURE_AI_ANALYSIS": {Label: "AI analysis feature", Hot: true, EnvVar: "FEATURE_AI_ANALYSIS"},
}

// secretMask is shown instead of secret values.
const secretMask = "••••••••"

// Settings resolves config: Postgres overrides → env → defaults. Values are
// cached briefly so per-request feature-flag checks stay off the DB.
type Settings struct {
	db       *DB
	mu       sync.RWMutex
	cached   map[string]json.RawMessage
	cachedAt time.Time
}

// NewSettings builds a resolver; db may be nil (env-only, read-only mode).
func NewSettings(db *DB) *Settings { return &Settings{db: db} }

const settingsCacheTTL = 10 * time.Second

// overrides returns the DB-stored overrides (fresh at most settingsCacheTTL old).
func (st *Settings) overrides(ctx context.Context) map[string]json.RawMessage {
	st.mu.RLock()
	if st.cached != nil && time.Since(st.cachedAt) < settingsCacheTTL {
		defer st.mu.RUnlock()
		return st.cached
	}
	st.mu.RUnlock()

	st.mu.Lock()
	defer st.mu.Unlock()
	if st.cached != nil && time.Since(st.cachedAt) < settingsCacheTTL {
		return st.cached
	}
	if st.db == nil {
		st.cached = map[string]json.RawMessage{}
	} else {
		m, err := st.db.allSettings(ctx)
		if err != nil {
			m = map[string]json.RawMessage{}
		}
		st.cached = m
	}
	st.cachedAt = time.Now()
	return st.cached
}

// Get resolves one key: DB override → env → "".
func (st *Settings) Get(ctx context.Context, key string) (val string, source string) {
	if raw, ok := st.overrides(ctx)[key]; ok {
		var s string
		if err := json.Unmarshal(raw, &s); err == nil {
			return s, "database"
		}
	}
	if v := os.Getenv(key); v != "" {
		return v, "env"
	}
	return "", "default"
}

// FeatureEnabled resolves a FEATURE_* key as a bool (default def).
func (st *Settings) FeatureEnabled(ctx context.Context, key string, def bool) bool {
	if raw, ok := st.overrides(ctx)[key]; ok {
		var b bool
		if err := json.Unmarshal(raw, &b); err == nil {
			return b
		}
		var s string
		if err := json.Unmarshal(raw, &s); err == nil {
			return s == "true" || s == "1"
		}
	}
	if v := os.Getenv(key); v != "" {
		return v == "true" || v == "1"
	}
	return def
}

// ConfigEntry is one key in GET /config's response.
type ConfigEntry struct {
	Key     string `json:"key"`
	Label   string `json:"label"`
	Value   string `json:"value"`            // masked for secrets
	Set     bool   `json:"set"`              // has any value
	Secret  bool   `json:"secret"`           // write-only
	Hot     bool   `json:"hot"`              // applies without restart
	Source  string `json:"source"`           // database | env | default
	Writable bool  `json:"writable"`         // false without Postgres
}

// handleGetConfig is GET /config — the effective config with masked secrets.
func (s *Server) handleGetConfig(w http.ResponseWriter, r *http.Request) {
	if s.settings == nil {
		writeError(w, http.StatusServiceUnavailable, "config unavailable")
		return
	}
	entries := make([]ConfigEntry, 0, len(configKeys))
	for key, meta := range configKeys {
		val, source := s.settings.Get(r.Context(), key)
		entry := ConfigEntry{
			Key:    key,
			Label:  meta.Label,
			Set:    val != "",
			Secret: meta.Secret,
			Hot:    meta.Hot,
			Source: source,
			Writable: s.db != nil,
		}
		if meta.Secret && val != "" {
			entry.Value = secretMask
		} else {
			entry.Value = val
		}
		entries = append(entries, entry)
	}
	writeJSON(w, http.StatusOK, map[string]any{"config": entries})
}

// handlePutConfig is PUT /config — body {key: value, ...} (whitelisted keys).
// String values and bool feature flags are both accepted. LLM/connection
// changes persist and apply on core restart; FEATURE_* keys are hot.
func (s *Server) handlePutConfig(w http.ResponseWriter, r *http.Request) {
	if s.db == nil || s.settings == nil {
		writeError(w, http.StatusServiceUnavailable, "config editing requires postgres")
		return
	}
	var body map[string]any
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if len(body) == 0 {
		writeError(w, http.StatusBadRequest, "no keys provided")
		return
	}
	updated := []string{}
	for key, val := range body {
		meta, ok := configKeys[key]
		if !ok {
			writeError(w, http.StatusBadRequest, "unknown or unmanaged key: "+key)
			return
		}
		// Empty string on a secret means "leave unchanged" (masked round-trip).
		if meta.Secret {
			if sv, ok := val.(string); ok && (sv == "" || sv == secretMask || strings.TrimSpace(sv) == "") {
				continue
			}
		}
		raw, err := json.Marshal(val)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid value for "+key)
			return
		}
		if err := s.db.setSetting(r.Context(), key, raw); err != nil {
			writeError(w, http.StatusInternalServerError, "failed to save "+key)
			return
		}
		updated = append(updated, key)
	}
	// Bust the resolver cache so hot flags apply immediately.
	s.settings.mu.Lock()
	s.settings.cached = nil
	s.settings.mu.Unlock()
	writeJSON(w, http.StatusOK, map[string]any{"updated": updated, "note": "LLM/connection changes apply on core restart; FEATURE_* flags apply immediately"})
}
