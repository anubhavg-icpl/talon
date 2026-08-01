package control

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/anubhavg-icpl/talon/internal/config"
	"github.com/anubhavg-icpl/talon/internal/strike"
)

// probeTimeout bounds each individual service probe.
const probeTimeout = 3 * time.Second

// ServiceHealth is one dependency's liveness, served by GET /health/services
// so the dashboard settings page can show REAL per-service status instead of
// static info cards.
type ServiceHealth struct {
	Name      string `json:"name"`
	Endpoint  string `json:"endpoint"`
	Status    string `json:"status"` // online | offline | unconfigured
	Detail    string `json:"detail"`
	LatencyMS int64  `json:"latency_ms"`
}

// handleServiceHealth is GET /health/services — actively probes every
// dependency talon-core orchestrates, concurrently, and reports honest
// per-service status. Results are cached in Redis for a few seconds so the
// dashboard's 10s polling stays instant even with slow probes.
func (s *Server) handleServiceHealth(w http.ResponseWriter, r *http.Request) {
	if s.cache != nil {
		if cached, ok := s.cache.Get(r.Context(), cacheKeyHealthServices); ok {
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("X-Cache", "hit")
			_, _ = w.Write([]byte(cached))
			return
		}
	}

	checks := []func() ServiceHealth{
		checkCore,
		s.checkPostgres,
		s.checkRedis,
		checkArsenal,
		checkMSF,
		checkRabbitMQ,
		checkOllama,
		checkONNXSLM,
	}
	out := make([]ServiceHealth, len(checks))
	var wg sync.WaitGroup
	for i, check := range checks {
		wg.Add(1)
		go func() {
			defer wg.Done()
			out[i] = check()
		}()
	}
	wg.Wait()

	if s.cache != nil {
		if raw, err := json.Marshal(map[string]any{"services": out}); err == nil {
			s.cache.Set(r.Context(), cacheKeyHealthServices, string(raw), healthServicesTTL)
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("X-Cache", "miss")
			_, _ = w.Write(raw)
			return
		}
	}
	writeJSON(w, http.StatusOK, map[string]any{"services": out})
}

// checkRedis pings the cache itself (PING); unconfigured when no cache.
func (s *Server) checkRedis() ServiceHealth {
	h := ServiceHealth{Name: "redis", Endpoint: ":6380"}
	if s.cache == nil {
		h.Status, h.Detail = "unconfigured", "REDIS_URL unset/unreachable — endpoints run uncached"
		return h
	}
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), probeTimeout)
	defer cancel()
	err := s.cache.Ping(ctx)
	h.LatencyMS = time.Since(start).Milliseconds()
	if err != nil {
		h.Status, h.Detail = "offline", "ping: "+err.Error()
		return h
	}
	h.Status, h.Detail = "online", "PING ok"
	return h
}

func checkCore() ServiceHealth {
	return ServiceHealth{
		Name:     "talon-core",
		Endpoint: ":8000",
		Status:   "online",
		Detail:   "serving — this response",
	}
}

// checkPostgres pings the configured database (SELECT 1 via pool ping).
func (s *Server) checkPostgres() ServiceHealth {
	h := ServiceHealth{Name: "postgres", Endpoint: ":5432"}
	if s.db == nil {
		h.Status, h.Detail = "unconfigured", "TALON_DATABASE_URL unset/unreachable — JSON persistence active"
		return h
	}
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), probeTimeout)
	defer cancel()
	err := s.db.pool.Ping(ctx)
	h.LatencyMS = time.Since(start).Milliseconds()
	if err != nil {
		h.Status, h.Detail = "offline", "ping: "+err.Error()
		return h
	}
	h.Status, h.Detail = "online", "SELECT 1 ok"
	return h
}

func checkArsenal() ServiceHealth {
	h := ServiceHealth{Name: "arsenal-engine", Endpoint: config.LoadHexstrikeConfig().ServerURL}
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), probeTimeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, h.Endpoint+"/health", nil)
	if err != nil {
		h.Status, h.Detail = "offline", err.Error()
		return h
	}
	resp, err := http.DefaultClient.Do(req)
	h.LatencyMS = time.Since(start).Milliseconds()
	if err != nil {
		h.Status, h.Detail = "offline", "unreachable: "+err.Error()
		return h
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		h.Status, h.Detail = "offline", fmt.Sprintf("GET /health → %d", resp.StatusCode)
		return h
	}
	h.Status, h.Detail = "online", "GET /health → 200"
	return h
}

func checkMSF() ServiceHealth {
	cfg := config.LoadMSFConfig()
	h := ServiceHealth{
		Name:     "msfrpcd",
		Endpoint: fmt.Sprintf("%s:%s", cfg.Server, cfg.Port),
	}
	if cfg.Password == "" {
		h.Status, h.Detail = "unconfigured", "MSF_PASSWORD not set"
		return h
	}
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), probeTimeout)
	defer cancel()
	client, err := strike.NewClient(ctx, cfg)
	if err != nil {
		h.Status, h.Detail = "offline", "auth: "+err.Error()
		return h
	}
	res, err := client.Call(ctx, "core.version")
	h.LatencyMS = time.Since(start).Milliseconds()
	if err != nil {
		h.Status, h.Detail = "offline", "core.version: "+err.Error()
		return h
	}
	h.Status = "online"
	if v, ok := res["version"].(string); ok {
		h.Detail = "RPC auth ok — framework " + v
	} else {
		h.Detail = "RPC auth ok"
	}
	return h
}

func checkRabbitMQ() ServiceHealth {
	rawURL := config.LoadAMQPConfig().URL
	h := ServiceHealth{Name: "rabbitmq", Endpoint: "amqp"}
	if rawURL == "" {
		h.Status, h.Detail = "unconfigured", "AMQP_URL not set (relay disabled)"
		return h
	}
	u, err := url.Parse(rawURL)
	if err != nil || u.Host == "" {
		h.Status, h.Detail = "offline", "unparseable AMQP_URL"
		return h
	}
	host := u.Host
	if u.Port() == "" {
		host = net.JoinHostPort(u.Host, "5672")
	}
	h.Endpoint = host
	start := time.Now()
	conn, err := net.DialTimeout("tcp", host, probeTimeout)
	h.LatencyMS = time.Since(start).Milliseconds()
	if err != nil {
		h.Status, h.Detail = "offline", "dial: "+err.Error()
		return h
	}
	_ = conn.Close()
	h.Status, h.Detail = "online", "AMQP port reachable"
	return h
}

func checkOllama() ServiceHealth {
	llmCfg := config.LoadLLMConfig()
	h := ServiceHealth{Name: "ollama", Endpoint: llmCfg.OllamaURL}
	optional := ""
	if llmCfg.Provider != "ollama" {
		optional = fmt.Sprintf(" (optional — LLM_PROVIDER=%s)", llmCfg.Provider)
	}
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), probeTimeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, h.Endpoint+"/api/version", nil)
	if err != nil {
		h.Status, h.Detail = "offline", err.Error()
		return h
	}
	resp, err := http.DefaultClient.Do(req)
	h.LatencyMS = time.Since(start).Milliseconds()
	if err != nil {
		h.Status, h.Detail = "offline", "unreachable"+optional
		return h
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		h.Status, h.Detail = "offline", fmt.Sprintf("GET /api/version → %d%s", resp.StatusCode, optional)
		return h
	}
	h.Status, h.Detail = "online", "GET /api/version → 200"+optional
	return h
}

// checkONNXSLM probes the local SmolLM / ONNX Runtime OpenAI-compatible
// service (compose profile "slm", default :8090). Optional unless
// LLM_PROVIDER=onnx.
func checkONNXSLM() ServiceHealth {
	llmCfg := config.LoadLLMConfig()
	// ONNXBaseURL is like http://localhost:8090/v1 — health lives at /health.
	base := strings.TrimRight(llmCfg.ONNXBaseURL, "/")
	endpoint := strings.TrimSuffix(base, "/v1")
	if endpoint == "" {
		endpoint = "http://localhost:8090"
	}
	h := ServiceHealth{Name: "onnx-slm", Endpoint: endpoint}
	optional := ""
	if llmCfg.Provider != "onnx" {
		optional = fmt.Sprintf(" (optional — LLM_PROVIDER=%s)", llmCfg.Provider)
	}
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), probeTimeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint+"/health", nil)
	if err != nil {
		h.Status, h.Detail = "offline", err.Error()
		return h
	}
	resp, err := http.DefaultClient.Do(req)
	h.LatencyMS = time.Since(start).Milliseconds()
	if err != nil {
		h.Status, h.Detail = "offline", "unreachable"+optional
		return h
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		h.Status, h.Detail = "offline", fmt.Sprintf("GET /health → %d%s", resp.StatusCode, optional)
		return h
	}
	var body struct {
		Backend string `json:"backend"`
		ModelID string `json:"model_id"`
		Ready   bool   `json:"ready"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&body)
	detail := "GET /health → 200" + optional
	if body.Backend != "" {
		detail = fmt.Sprintf("%s backend=%s model=%s", detail, body.Backend, body.ModelID)
	}
	h.Status, h.Detail = "online", detail
	return h
}
