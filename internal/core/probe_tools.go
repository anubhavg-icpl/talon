package core

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/anubhavg-icpl/talon/internal/llm"
)

// httpProbeBatchToolSpec returns the tool spec for batched HTTP probing.
// This lets the model fire multiple request variants in one tool call,
// returning full bodies and exact request surfaces.
func httpProbeBatchToolSpec() llm.ToolSpec {
	return llm.ToolSpec{
		Name:        "http_probe_batch",
		Description: "Send multiple HTTP request variants in one call. " +
			"Each variant specifies method, path, headers, and body. " +
			"Returns status, headers, and body snippet for each. " +
			"Use for batch-testing endpoints, methods, or payload variations.",
		InputSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"base_url": map[string]any{
					"type":        "string",
					"description": "Base URL (e.g. http://target:8080)",
				},
				"requests": map[string]any{
					"type":        "array",
					"description": "Array of request variants",
					"items": map[string]any{
						"type": "object",
						"properties": map[string]any{
							"label":  map[string]any{"type": "string"},
							"method": map[string]any{"type": "string", "default": "GET"},
							"path":   map[string]any{"type": "string"},
							"headers": map[string]any{
								"type": "object",
								"additionalProperties": map[string]any{"type": "string"},
							},
							"body": map[string]any{"type": "string"},
						},
					},
				},
				"timeout_seconds": map[string]any{
					"type":        "integer",
					"description": "Per-request timeout (default 10, max 30)",
				},
			},
			"required": []any{"base_url", "requests"},
		},
	}
}

// probeRequest represents one request variant in a batch.
type probeRequest struct {
	Label   string            `json:"label"`
	Method  string            `json:"method"`
	Path    string            `json:"path"`
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}

// probeResponse captures the full result of one probe.
type probeResponse struct {
	Label        string            `json:"label"`
	StatusCode   int               `json:"status_code"`
	Status       string            `json:"status"`
	RespHeaders  map[string]string `json:"resp_headers"`
	BodySnippet  string            `json:"body_snippet"`
	BodySize     int               `json:"body_size"`
	Error        string            `json:"error,omitempty"`
	DurationMs   int64             `json:"duration_ms"`
	RequestSent  string            `json:"request_sent"`
}

const (
	maxProbeBatchSize     = 20
	maxProbeBodySnippet   = 2000
	defaultProbeTimeoutSec = 10
	maxProbeTimeoutSec     = 30
)

// handleHTTPProbeBatch sends batched HTTP requests and returns results.
func handleHTTPProbeBatch(ctx context.Context, args map[string]any, tr *tracker, store *EvidenceStore) (string, bool) {
	baseURL, _ := args["base_url"].(string)
	baseURL = strings.TrimRight(baseURL, "/")
	if baseURL == "" {
		return "http_probe_batch: base_url is required", true
	}

	timeoutSec := defaultProbeTimeoutSec
	if v, ok := args["timeout_seconds"].(float64); ok {
		timeoutSec = int(v)
	}
	if timeoutSec <= 0 {
		timeoutSec = defaultProbeTimeoutSec
	}
	if timeoutSec > maxProbeTimeoutSec {
		timeoutSec = maxProbeTimeoutSec
	}

	// Parse requests array
	reqsRaw, ok := args["requests"].([]any)
	if !ok || len(reqsRaw) == 0 {
		return "http_probe_batch: requests array is required", true
	}
	if len(reqsRaw) > maxProbeBatchSize {
		reqsRaw = reqsRaw[:maxProbeBatchSize]
	}

	var reqs []probeRequest
	for _, raw := range reqsRaw {
		m, ok := raw.(map[string]any)
		if !ok {
			continue
		}
		req := probeRequest{
			Label:  strArg(m, "label"),
			Method: strings.ToUpper(strArgOr(m, "method", "GET")),
			Path:   strArg(m, "path"),
			Body:   strArg(m, "body"),
		}
		if req.Path == "" {
			req.Path = "/"
		}
		if h, ok := m["headers"].(map[string]any); ok {
			req.Headers = make(map[string]string)
			for k, v := range h {
				req.Headers[k] = fmt.Sprintf("%v", v)
			}
		}
		reqs = append(reqs, req)
	}

	// Send requests
	results := make([]probeResponse, 0, len(reqs))
	for _, req := range reqs {
		result := sendProbe(ctx, baseURL, req, timeoutSec)
		results = append(results, result)
		// Record each response as evidence
		if store != nil {
			label := req.Label
			if label == "" {
				label = req.Method + " " + req.Path
			}
			store.Record("http_probe_batch", label, result.BodySnippet, result.StatusCode)
		}
	}

	// Build output
	payload := map[string]any{
		"base_url": baseURL,
		"total":    len(results),
		"results":  results,
	}
	raw, _ := json.MarshalIndent(payload, "", "  ")
	if tr != nil {
		tr.record("http_probe_batch", args, fmt.Sprintf("sent=%d", len(results)))
	}
	return string(raw), false
}

func sendProbe(ctx context.Context, baseURL string, req probeRequest, timeoutSec int) probeResponse {
	fullURL := baseURL + req.Path
	method := req.Method
	if method == "" {
		method = "GET"
	}

	var bodyReader io.Reader
	if req.Body != "" {
		bodyReader = strings.NewReader(req.Body)
	}

	httpReq, err := http.NewRequestWithContext(ctx, method, fullURL, bodyReader)
	if err != nil {
		return probeResponse{Label: req.Label, Error: err.Error()}
	}

	// Set headers
	for k, v := range req.Headers {
		httpReq.Header.Set(k, v)
	}

	// Build request representation
	reqSent := fmt.Sprintf("%s %s HTTP/1.1\nHost: %s", method, req.Path, httpReq.URL.Host)
	for k, v := range req.Headers {
		reqSent += fmt.Sprintf("\n%s: %s", k, v)
	}
	if req.Body != "" {
		reqSent += "\n\n" + req.Body
	}

	client := &http.Client{Timeout: time.Duration(timeoutSec) * time.Second}
	start := time.Now()
	resp, err := client.Do(httpReq)
	duration := time.Since(start)

	if err != nil {
		return probeResponse{
			Label:       req.Label,
			Error:       err.Error(),
			DurationMs:  duration.Milliseconds(),
			RequestSent: reqSent,
		}
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(io.LimitReader(resp.Body, 50000))
	bodyStr := string(body)
	snippet := bodyStr
	if len(snippet) > maxProbeBodySnippet {
		snippet = snippet[:maxProbeBodySnippet] + "\n…[truncated]"
	}

	respHeaders := make(map[string]string)
	for k := range resp.Header {
		respHeaders[k] = resp.Header.Get(k)
	}

	return probeResponse{
		Label:       req.Label,
		StatusCode:  resp.StatusCode,
		Status:      resp.Status,
		RespHeaders: respHeaders,
		BodySnippet: snippet,
		BodySize:    len(bodyStr),
		DurationMs:  duration.Milliseconds(),
		RequestSent: reqSent,
	}
}

// --- Web vulnerability tools ---

// webVulnToolSpecs returns web-specific analysis tools.
func webVulnToolSpecs() []llm.ToolSpec {
	return []llm.ToolSpec{
		{
			Name:        "web_headers_audit",
			Description: "Analyze HTTP response headers for security issues. " +
				"Checks for missing security headers (CSP, HSTS, X-Frame-Options, etc.), " +
				"information leakage (Server, X-Powered-By), and misconfiguration.",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"url": map[string]any{
						"type":        "string",
						"description": "URL to audit",
					},
				},
				"required": []any{"url"},
			},
		},
		{
			Name:        "js_endpoint_extract",
			Description: "Extract API endpoints, URLs, and interesting patterns from JavaScript code. " +
				"Fetches a JS file URL or analyzes provided JS source, extracting fetch() calls, " +
				"API paths, AJAX endpoints, and potential secrets.",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"url": map[string]any{
						"type":        "string",
						"description": "URL of the JS file to analyze",
					},
					"source": map[string]any{
						"type":        "string",
						"description": "JS source code to analyze directly (alternative to url)",
					},
				},
			},
		},
	}
}

// headerFinding represents one security header analysis result.
type headerFinding struct {
	Header   string `json:"header"`
	Status   string `json:"status"` // "missing", "present", "weak"
	Value    string `json:"value,omitempty"`
	Issue    string `json:"issue,omitempty"`
	Severity string `json:"severity,omitempty"`
}

// handleWebHeadersAudit fetches and analyzes security headers.
func handleWebHeadersAudit(ctx context.Context, args map[string]any, tr *tracker) (string, bool) {
	targetURL, _ := args["url"].(string)
	if targetURL == "" {
		return "web_headers_audit: url is required", true
	}

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequestWithContext(ctx, "GET", targetURL, nil)
	if err != nil {
		return fmt.Sprintf("error: %v", err), true
	}
	req.Header.Set("User-Agent", "Talon-Security-Scanner/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Sprintf("fetch error: %v", err), true
	}
	defer resp.Body.Close()

	var findings []headerFinding

	// Security headers that should be present
	securityHeaders := []struct {
		Header  string
		Check   func(string) (bool, string) // returns (ok, issue)
	}{
		{"Content-Security-Policy", func(v string) (bool, string) {
			if v == "" {
				return false, ""
			}
			if strings.Contains(v, "unsafe-inline") || strings.Contains(v, "unsafe-eval") {
				return false, "CSP contains unsafe-inline or unsafe-eval"
			}
			return true, ""
		}},
		{"Strict-Transport-Security", nil},
		{"X-Frame-Options", nil},
		{"X-Content-Type-Options", nil},
		{"Referrer-Policy", nil},
		{"Permissions-Policy", nil},
	}

	for _, sh := range securityHeaders {
		val := resp.Header.Get(sh.Header)
		if val == "" {
			findings = append(findings, headerFinding{
				Header: sh.Header, Status: "missing",
				Issue:    "Security header not set",
				Severity: "medium",
			})
		} else if sh.Check != nil {
			ok, issue := sh.Check(val)
			if !ok {
				findings = append(findings, headerFinding{
					Header: sh.Header, Status: "weak", Value: val,
					Issue: issue, Severity: "medium",
				})
			} else {
				findings = append(findings, headerFinding{
					Header: sh.Header, Status: "present", Value: val,
				})
			}
		} else {
			findings = append(findings, headerFinding{
				Header: sh.Header, Status: "present", Value: val,
			})
		}
	}

	// Information leakage headers
	infoHeaders := []string{"Server", "X-Powered-By", "X-AspNet-Version", "X-AspNetMvc-Version"}
	for _, h := range infoHeaders {
		val := resp.Header.Get(h)
		if val != "" {
			findings = append(findings, headerFinding{
				Header: h, Status: "present", Value: val,
				Issue:    "Information disclosure: reveals technology stack",
				Severity: "low",
			})
		}
	}

	// Cookie analysis
	for _, cookie := range resp.Cookies() {
		issues := []string{}
		if !cookie.Secure {
			issues = append(issues, "missing Secure flag")
		}
		if !cookie.HttpOnly {
			issues = append(issues, "missing HttpOnly flag")
		}
		if cookie.SameSite == 0 {
			issues = append(issues, "missing SameSite attribute")
		}
		if len(issues) > 0 {
			findings = append(findings, headerFinding{
				Header: "Set-Cookie", Status: "weak", Value: cookie.Name,
				Issue:    strings.Join(issues, "; "),
				Severity: "medium",
			})
		}
	}

	payload := map[string]any{
		"url":         targetURL,
		"status_code": resp.StatusCode,
		"findings":    findings,
		"summary": map[string]int{
			"total":    len(findings),
			"missing":  countStatus(findings, "missing"),
			"weak":     countStatus(findings, "weak"),
			"present":  countStatus(findings, "present"),
		},
	}
	raw, _ := json.MarshalIndent(payload, "", "  ")
	if tr != nil {
		tr.record("web_headers_audit", args, fmt.Sprintf("findings=%d", len(findings)))
	}
	return string(raw), false
}

func countStatus(findings []headerFinding, status string) int {
	count := 0
	for _, f := range findings {
		if f.Status == status {
			count++
		}
	}
	return count
}

// handleJSEndpointExtract extracts endpoints from JavaScript.
func handleJSEndpointExtract(ctx context.Context, args map[string]any, tr *tracker) (string, bool) {
	source, _ := args["source"].(string)
	urlStr, _ := args["url"].(string)

	if source == "" && urlStr == "" {
		return "js_endpoint_extract: either url or source is required", true
	}

	if source == "" {
		client := &http.Client{Timeout: 10 * time.Second}
		req, err := http.NewRequestWithContext(ctx, "GET", urlStr, nil)
		if err != nil {
			return fmt.Sprintf("js_endpoint_extract: invalid url %q: %v", urlStr, err), true
		}
		req.Header.Set("User-Agent", "Talon-Security-Scanner/1.0")
		resp, err := client.Do(req)
		if err != nil {
			return fmt.Sprintf("fetch error: %v", err), true
		}
		defer resp.Body.Close()
		data, _ := io.ReadAll(io.LimitReader(resp.Body, 500000))
		source = string(data)
	}

	// Extract patterns
	endpoints := extractJSEndpoints(source)
	secrets := extractJSSecrets(source)

	payload := map[string]any{
		"source_size": len(source),
		"endpoints":   endpoints,
		"secrets":     secrets,
		"total_endpoints": len(endpoints),
		"total_secrets":   len(secrets),
	}
	raw, _ := json.MarshalIndent(payload, "", "  ")
	if tr != nil {
		tr.record("js_endpoint_extract", args, fmt.Sprintf("endpoints=%d secrets=%d", len(endpoints), len(secrets)))
	}
	return string(raw), false
}

var (
	apiPathRe   = regexp.MustCompile(`["'` + "`" + `](/(?:api|v[0-9]+|graphql|rest|service|backend|admin|internal)/[A-Za-z0-9_./?=&-]{1,200})["'` + "`" + `]`)
	urlLitRe    = regexp.MustCompile(`["'` + "`" + `](https?://[^\s"'` + "`" + `<>{})\]]{5,200})["'` + "`" + `]`)
	fetchCallRe = regexp.MustCompile(`(?:fetch|axios\.\w+|XMLHttpRequest|\.open)\s*\(\s*["'` + "`" + `]([^"'` + "`" + `]{1,200})["'` + "`" + `]`)
	secretRe    = regexp.MustCompile(`(?i)(?:api[_-]?key|secret|token|password|passwd|auth|bearer)\s*[:=]\s*["'` + "`" + `]([A-Za-z0-9+/=_-]{8,80})["'` + "`" + `]`)
)

func extractJSEndpoints(source string) []string {
	seen := make(map[string]bool)
	var out []string

	extract := func(re *regexp.Regexp) {
		matches := re.FindAllStringSubmatch(source, -1)
		for _, m := range matches {
			if len(m) > 1 && !seen[m[1]] {
				seen[m[1]] = true
				out = append(out, m[1])
			}
		}
	}

	extract(apiPathRe)
	extract(urlLitRe)
	extract(fetchCallRe)

	return out
}

func extractJSSecrets(source string) []string {
	seen := make(map[string]bool)
	var out []string

	matches := secretRe.FindAllStringSubmatch(source, -1)
	for _, m := range matches {
		if len(m) > 1 && !seen[m[1]] {
			seen[m[1]] = true
			out = append(out, m[1])
		}
	}
	return out
}

func strArgOr(args map[string]any, key, def string) string {
	v, ok := args[key].(string)
	if !ok || v == "" {
		return def
	}
	return v
}
