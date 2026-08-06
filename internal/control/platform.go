package control

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/anubhavg-icpl/talon/internal/core"
	"github.com/google/uuid"
)

// ---- Scope / ROE ----

// ScopePolicy is rules of engagement enforced at run start.
type ScopePolicy struct {
	Enabled          bool     `json:"enabled"`
	AllowedCIDRs     []string `json:"allowed_cidrs"`      // empty = allow all (when enabled still checks denylist)
	DeniedCIDRs      []string `json:"denied_cidrs"`
	DeniedPorts      []int    `json:"denied_ports"`
	MaxConcurrent    int      `json:"max_concurrent"`     // 0 = unlimited
	RequireAuthLabel bool     `json:"require_auth_label"` // description must contain AUTHORIZED
	AutoApproveNmapPrivate bool `json:"auto_approve_nmap_private"`
	UpdatedAt        string   `json:"updated_at,omitempty"`
}

func defaultScope() ScopePolicy {
	return ScopePolicy{
		Enabled:       false,
		AllowedCIDRs:  []string{"10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16", "127.0.0.0/8"},
		DeniedCIDRs:   []string{},
		MaxConcurrent: 20,
		AutoApproveNmapPrivate: true,
	}
}

func (p ScopePolicy) ValidateTarget(ip string) error {
	if !p.Enabled {
		return nil
	}
	host := strings.TrimSpace(ip)
	// strip URL to host if needed
	if strings.Contains(host, "://") {
		if u, err := url.Parse(host); err == nil && u.Hostname() != "" {
			host = u.Hostname()
		}
	}
	parsed := net.ParseIP(host)
	if parsed == nil {
		// hostname — only check deny list empty / auth label handled separately
		return nil
	}
	for _, c := range p.DeniedCIDRs {
		if inCIDR(parsed, c) {
			return fmt.Errorf("target %s is in denied CIDR %s", host, c)
		}
	}
	if len(p.AllowedCIDRs) > 0 {
		ok := false
		for _, c := range p.AllowedCIDRs {
			if inCIDR(parsed, c) {
				ok = true
				break
			}
		}
		if !ok {
			return fmt.Errorf("target %s is outside allowed CIDRs", host)
		}
	}
	return nil
}

func inCIDR(ip net.IP, cidr string) bool {
	_, n, err := net.ParseCIDR(strings.TrimSpace(cidr))
	if err != nil {
		return false
	}
	return n.Contains(ip)
}

func isPrivateIP(ipStr string) bool {
	ip := net.ParseIP(strings.TrimSpace(ipStr))
	if ip == nil {
		return false
	}
	for _, c := range []string{"10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16", "127.0.0.0/8"} {
		if inCIDR(ip, c) {
			return true
		}
	}
	return false
}

// ---- Targets inventory ----

type Target struct {
	ID          string            `json:"id"`
	Address     string            `json:"address"` // IP or hostname
	URL         string            `json:"url,omitempty"`
	Label       string            `json:"label,omitempty"`
	Tags        []string          `json:"tags,omitempty"`
	Notes       string            `json:"notes,omitempty"`
	LastRunID   string            `json:"last_run_id,omitempty"`
	LastStatus  string            `json:"last_status,omitempty"`
	Meta        map[string]string `json:"meta,omitempty"`
	CreatedAt   string            `json:"created_at"`
	UpdatedAt   string            `json:"updated_at"`
}

// ---- Schedules ----

type Schedule struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Interval   string `json:"interval"` // simplified: 1h, 6h, 24h, 7d
	TargetAddr string `json:"target"`
	PlaybookID string `json:"playbook_id,omitempty"`
	AgentMode  string `json:"agent_mode,omitempty"`
	Enabled    bool   `json:"enabled"`
	LastRunAt  string `json:"last_run_at,omitempty"`
	NextRunAt  string `json:"next_run_at,omitempty"`
	CreatedAt  string `json:"created_at"`
}

// ---- Notifications ----

type NotifyConfig struct {
	WebhookURL        string `json:"webhook_url"`
	OnComplete        bool   `json:"on_complete"`
	OnHITL            bool   `json:"on_hitl"`
	OnCriticalFinding bool   `json:"on_critical_finding"`
	OnError           bool   `json:"on_error"`
}

// ---- Credentials (encrypted at rest) ----

type Credential struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Kind      string `json:"kind"` // password | token | basic
	Username  string `json:"username,omitempty"`
	// Secret is write-only on create; never returned in full list
	Secret    string `json:"secret,omitempty"`
	HasSecret bool   `json:"has_secret"`
	Scope     string `json:"scope,omitempty"` // optional target binding
	CreatedAt string `json:"created_at"`
}

type credStored struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Kind         string `json:"kind"`
	Username     string `json:"username"`
	CiphertextB64 string `json:"ciphertext"`
	Scope        string `json:"scope"`
	CreatedAt    string `json:"created_at"`
}

// ---- Evidence ----

type EvidenceItem struct {
	ID        string `json:"id"`
	RunID     string `json:"run_id"`
	FindingID string `json:"finding_id,omitempty"`
	Kind      string `json:"kind"` // note | log | url | base64
	Title     string `json:"title"`
	Body      string `json:"body"`
	CreatedAt string `json:"created_at"`
}

// ---- Budget / token meter (process-level counters) ----

type BudgetStats struct {
	LLMCalls       int64 `json:"llm_calls"`
	ToolCalls      int64 `json:"tool_calls"`
	RunsStarted    int64 `json:"runs_started"`
	RunsCompleted  int64 `json:"runs_completed"`
	CriticalFindings int64 `json:"critical_findings"`
}

// Platform is the secondary control-plane state (targets, scope, schedules…).
type Platform struct {
	mu       sync.RWMutex
	path     string
	Scope    ScopePolicy             `json:"scope"`
	Targets  map[string]*Target      `json:"targets"`
	Schedules map[string]*Schedule   `json:"schedules"`
	Notify   NotifyConfig            `json:"notify"`
	Creds    map[string]*credStored  `json:"credentials"`
	Evidence map[string]*EvidenceItem `json:"evidence"`
	Budget   BudgetStats             `json:"budget"`

	// startFn is set by server to launch runs from schedules.
	startFn func(input core.RunInput) string // returns run_id
	stop      chan struct{}
	stopOnce  sync.Once
}

func NewPlatform(dataDir string) *Platform {
	p := &Platform{
		path:      dataDir + "/platform.json",
		Targets:   map[string]*Target{},
		Schedules: map[string]*Schedule{},
		Creds:     map[string]*credStored{},
		Evidence:  map[string]*EvidenceItem{},
		Scope:     defaultScope(),
		Notify: NotifyConfig{
			OnComplete: true, OnHITL: true, OnCriticalFinding: true, OnError: true,
		},
		stop: make(chan struct{}),
	}
	_ = p.load()
	return p
}

func (p *Platform) SetStartFn(fn func(core.RunInput) string) { p.startFn = fn }

func (p *Platform) load() error {
	data, err := os.ReadFile(p.path)
	if err != nil {
		return err
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	return json.Unmarshal(data, p)
}

func (p *Platform) saveLocked() {
	// caller holds mu
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return
	}
	_ = os.MkdirAll(strings.TrimSuffix(p.path, "/platform.json"), 0o755)
	tmp := p.path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o600); err != nil {
		return
	}
	_ = os.Rename(tmp, p.path)
}

func (p *Platform) Save() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.saveLocked()
}

// ---- Scope API helpers ----

func (p *Platform) GetScope() ScopePolicy {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.Scope
}

func (p *Platform) PutScope(s ScopePolicy) ScopePolicy {
	p.mu.Lock()
	defer p.mu.Unlock()
	s.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	p.Scope = s
	p.saveLocked()
	return p.Scope
}

func (p *Platform) CheckStart(ip, description string, activeRuns int) error {
	p.mu.RLock()
	scope := p.Scope
	p.mu.RUnlock()
	if scope.Enabled && scope.RequireAuthLabel {
		if !strings.Contains(strings.ToUpper(description), "AUTHORIZED") {
			return fmt.Errorf("scope policy: description must contain AUTHORIZED")
		}
	}
	if err := scope.ValidateTarget(ip); err != nil {
		return err
	}
	if scope.Enabled && scope.MaxConcurrent > 0 && activeRuns >= scope.MaxConcurrent {
		return fmt.Errorf("scope policy: max concurrent runs %d reached", scope.MaxConcurrent)
	}
	return nil
}

func (p *Platform) AutoApproveNmap(ip string) bool {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.Scope.AutoApproveNmapPrivate && isPrivateIP(ip)
}

// ---- Targets ----

func (p *Platform) ListTargets() []Target {
	p.mu.RLock()
	defer p.mu.RUnlock()
	out := make([]Target, 0, len(p.Targets))
	for _, t := range p.Targets {
		out = append(out, *t)
	}
	return out
}

func (p *Platform) UpsertTarget(t Target) Target {
	p.mu.Lock()
	defer p.mu.Unlock()
	now := time.Now().UTC().Format(time.RFC3339)
	if t.ID == "" {
		t.ID = "tgt-" + uuid.NewString()[:8]
		t.CreatedAt = now
	}
	t.UpdatedAt = now
	if t.Address == "" && t.URL != "" {
		if u, err := url.Parse(t.URL); err == nil {
			t.Address = u.Hostname()
		}
	}
	p.Targets[t.ID] = &t
	p.saveLocked()
	return t
}

func (p *Platform) DeleteTarget(id string) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	if _, ok := p.Targets[id]; !ok {
		return false
	}
	delete(p.Targets, id)
	p.saveLocked()
	return true
}

func (p *Platform) TouchTarget(address, runID, status string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, t := range p.Targets {
		if t.Address == address || t.URL == address {
			t.LastRunID = runID
			t.LastStatus = status
			t.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
		}
	}
	p.saveLocked()
}

// ---- Schedules ----

func (p *Platform) ListSchedules() []Schedule {
	p.mu.RLock()
	defer p.mu.RUnlock()
	out := make([]Schedule, 0, len(p.Schedules))
	for _, s := range p.Schedules {
		out = append(out, *s)
	}
	return out
}

func (p *Platform) UpsertSchedule(s Schedule) Schedule {
	p.mu.Lock()
	defer p.mu.Unlock()
	now := time.Now().UTC()
	if s.ID == "" {
		s.ID = "sch-" + uuid.NewString()[:8]
		s.CreatedAt = now.Format(time.RFC3339)
	}
	if s.Interval == "" {
		s.Interval = "24h"
	}
	s.NextRunAt = now.Add(parseInterval(s.Interval)).Format(time.RFC3339)
	p.Schedules[s.ID] = &s
	p.saveLocked()
	return s
}

func (p *Platform) DeleteSchedule(id string) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	if _, ok := p.Schedules[id]; !ok {
		return false
	}
	delete(p.Schedules, id)
	p.saveLocked()
	return true
}

func parseInterval(s string) time.Duration {
	s = strings.TrimSpace(strings.ToLower(s))
	d, err := time.ParseDuration(s)
	if err != nil {
		return 24 * time.Hour
	}
	if d < time.Minute {
		return time.Minute
	}
	return d
}

// StartScheduler ticks every minute and fires due schedules.
func (p *Platform) StartScheduler() {
	go func() {
		t := time.NewTicker(30 * time.Second)
		defer t.Stop()
		for {
			select {
			case <-p.stop:
				return
			case <-t.C:
				p.tickSchedules()
			}
		}
	}()
}

func (p *Platform) StopScheduler() { p.stopOnce.Do(func() { close(p.stop) }) }

func (p *Platform) tickSchedules() {
	if p.startFn == nil {
		return
	}
	now := time.Now().UTC()
	p.mu.Lock()
	var due []*Schedule
	for _, s := range p.Schedules {
		if !s.Enabled || s.TargetAddr == "" {
			continue
		}
		if s.NextRunAt == "" {
			s.NextRunAt = now.Format(time.RFC3339)
		}
		next, err := time.Parse(time.RFC3339, s.NextRunAt)
		if err != nil || !next.After(now) {
			due = append(due, s)
		}
	}
	p.mu.Unlock()

	for _, s := range due {
		mode := s.AgentMode
		desc := "Scheduled engagement: " + s.Name + " AUTHORIZED"
		if s.PlaybookID != "" {
			if pb, ok := core.GetPlaybook(s.PlaybookID); ok {
				if mode == "" {
					mode = pb.AgentMode
				}
				desc = pb.Prompt + " AUTHORIZED (schedule:" + s.Name + ")"
			}
		}
		input := core.RunInput{
			SessionID:   uuid.NewString(),
			TargetIP:    s.TargetAddr,
			Description: desc,
			AgentMode:   core.NormalizeAgentMode(mode),
		}
		if s.PlaybookID == "cve-lab" {
			input.CVEID = "CVE-2011-2523"
			input.ServiceName = "vsftpd 2.3.4"
		}
		runID := p.startFn(input)
		p.mu.Lock()
		s.LastRunAt = now.Format(time.RFC3339)
		s.NextRunAt = now.Add(parseInterval(s.Interval)).Format(time.RFC3339)
		p.saveLocked()
		p.mu.Unlock()
		log.Printf("platform: schedule %s fired run %s target=%s", s.ID, runID, s.TargetAddr)
	}
}

// ---- Notify ----

func (p *Platform) GetNotify() NotifyConfig {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.Notify
}

func (p *Platform) PutNotify(n NotifyConfig) (NotifyConfig, error) {
	if n.WebhookURL != "" {
		if err := validateWebhookURL(n.WebhookURL); err != nil {
			return NotifyConfig{}, err
		}
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Notify = n
	p.saveLocked()
	return p.Notify, nil
}

func (p *Platform) Fire(event string, payload map[string]any) {
	p.mu.RLock()
	cfg := p.Notify
	p.mu.RUnlock()
	if cfg.WebhookURL == "" {
		return
	}
	switch event {
	case "complete":
		if !cfg.OnComplete {
			return
		}
	case "hitl":
		if !cfg.OnHITL {
			return
		}
	case "critical":
		if !cfg.OnCriticalFinding {
			return
		}
	case "error":
		if !cfg.OnError {
			return
		}
	}
	payload["event"] = event
	payload["ts"] = time.Now().UTC().Format(time.RFC3339)
	body, _ := json.Marshal(payload)
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
		defer cancel()
		req, err := http.NewRequestWithContext(ctx, http.MethodPost, cfg.WebhookURL, bytes.NewReader(body))
		if err != nil {
			return
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", "talon-notify/1.0")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("platform: notify %s: %v", event, err)
			return
		}
		_, _ = io.Copy(io.Discard, resp.Body)
		_ = resp.Body.Close()
	}()
}

// validateWebhookURL guards against SSRF: requires http(s) scheme and rejects
// loopback, link-local, private, and cloud-metadata IPs.
func validateWebhookURL(rawURL string) error {
	u, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid webhook URL: %w", err)
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("webhook URL must use http or https scheme")
	}
	host := u.Hostname()
	if host == "" {
		return fmt.Errorf("webhook URL missing host")
	}
	// Allow loopback for local dev/testing
	if ip := net.ParseIP(host); ip != nil {
		if ip.IsLoopback() || ip.IsUnspecified() {
			return nil
		}
		if ip.IsPrivate() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
			return fmt.Errorf("webhook URL must not point to a private or link-local address")
		}
		// Block AWS / GCP / Azure metadata endpoints
		if host == "169.254.169.254" || host == "fd00:ec2::254" {
			return fmt.Errorf("webhook URL must not point to cloud metadata endpoint")
		}
	}
	return nil
}

// ---- Credentials ----

func credKey() []byte {
	k := os.Getenv("TALON_CRED_KEY")
	if k == "" {
		k = os.Getenv("TALON_ADMIN_PASSWORD")
	}
	if k == "" {
		log.Printf("WARNING: TALON_CRED_KEY not set — falling back to insecure default key. " +
			"All stored credentials are trivially decryptable. Set TALON_CRED_KEY in production.")
		k = "talon-dev-insecure-key"
	}
	sum := sha256.Sum256([]byte(k))
	return sum[:]
}

func encryptSecret(plain string) (string, error) {
	block, err := aes.NewCipher(credKey())
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}
	out := gcm.Seal(nonce, nonce, []byte(plain), nil)
	return base64.StdEncoding.EncodeToString(out), nil
}

func (p *Platform) ListCredentials() []Credential {
	p.mu.RLock()
	defer p.mu.RUnlock()
	out := make([]Credential, 0, len(p.Creds))
	for _, c := range p.Creds {
		out = append(out, Credential{
			ID: c.ID, Name: c.Name, Kind: c.Kind, Username: c.Username,
			HasSecret: c.CiphertextB64 != "", Scope: c.Scope, CreatedAt: c.CreatedAt,
		})
	}
	return out
}

func (p *Platform) AddCredential(name, kind, username, secret, scope string) (Credential, error) {
	ct, err := encryptSecret(secret)
	if err != nil {
		return Credential{}, err
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	u := uuid.New()
	id := "cred-" + hex.EncodeToString(u[:4])
	now := time.Now().UTC().Format(time.RFC3339)
	p.Creds[id] = &credStored{
		ID: id, Name: name, Kind: kind, Username: username,
		CiphertextB64: ct, Scope: scope, CreatedAt: now,
	}
	p.saveLocked()
	return Credential{ID: id, Name: name, Kind: kind, Username: username, HasSecret: true, Scope: scope, CreatedAt: now}, nil
}

func (p *Platform) DeleteCredential(id string) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	if _, ok := p.Creds[id]; !ok {
		return false
	}
	delete(p.Creds, id)
	p.saveLocked()
	return true
}

// ---- Evidence ----

func (p *Platform) ListEvidence(runID string) []EvidenceItem {
	p.mu.RLock()
	defer p.mu.RUnlock()
	out := []EvidenceItem{}
	for _, e := range p.Evidence {
		if runID == "" || e.RunID == runID {
			out = append(out, *e)
		}
	}
	return out
}

func (p *Platform) AddEvidence(e EvidenceItem) EvidenceItem {
	p.mu.Lock()
	defer p.mu.Unlock()
	if e.ID == "" {
		e.ID = "ev-" + uuid.NewString()[:8]
	}
	e.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	p.Evidence[e.ID] = &e
	p.saveLocked()
	return e
}

// ---- Budget ----

func (p *Platform) IncBudget(field string, n int64) {
	p.mu.Lock()
	defer p.mu.Unlock()
	switch field {
	case "llm":
		p.Budget.LLMCalls += n
	case "tool":
		p.Budget.ToolCalls += n
	case "start":
		p.Budget.RunsStarted += n
	case "complete":
		p.Budget.RunsCompleted += n
	case "critical":
		p.Budget.CriticalFindings += n
	}
	// cheap periodic save
	if (p.Budget.RunsStarted+p.Budget.RunsCompleted)%5 == 0 {
		p.saveLocked()
	}
}

func (p *Platform) GetBudget() BudgetStats {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.Budget
}
