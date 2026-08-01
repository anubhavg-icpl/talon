package control

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/anubhavg-icpl/talon/internal/core"
	"github.com/anubhavg-icpl/talon/internal/llm"
)

// SLM tools are a curated, UI-safe subset of the Talon codebase surface.
// SmolLM (and other SLMs) call these via a text protocol or native function
// calling — never full MCP exploit tools (nmap/msf) from the dashboard chat.
// Destructive agent work stays on /input/start runs.

// slmToolDef is one catalog entry exposed to the model + GET /llm/tools.
type slmToolDef struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Parameters  map[string]any `json:"parameters"`
	// Safe is always true for this catalog (read-only / soft ops).
	Safe bool `json:"safe"`
}

func slmToolCatalog() []slmToolDef {
	return []slmToolDef{
		{
			Name:        "list_runs",
			Description: "List recent pentest runs (id, target, status, findings count). Use to answer 'what is running' or pick a run_id.",
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"limit":  map[string]any{"type": "integer", "description": "Max runs (default 10, max 50)"},
					"offset": map[string]any{"type": "integer", "description": "Pagination offset"},
				},
			},
			Safe: true,
		},
		{
			Name:        "get_run_status",
			Description: "Get status, output snippet, findings summary, and agent mode for one run_id.",
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"run_id": map[string]any{"type": "string", "description": "Run UUID"},
				},
				"required": []any{"run_id"},
			},
			Safe: true,
		},
		{
			Name:        "get_run_tools",
			Description: "Get the tool call log for a run (names, args summary, short outputs). Use to explain what the agent did.",
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"run_id": map[string]any{"type": "string"},
					"limit":  map[string]any{"type": "integer", "description": "Max tool records (default 20)"},
				},
				"required": []any{"run_id"},
			},
			Safe: true,
		},
		{
			Name:        "get_findings",
			Description: "List structured findings for a run (severity, title, status, 3-gate).",
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"run_id": map[string]any{"type": "string"},
				},
				"required": []any{"run_id"},
			},
			Safe: true,
		},
		{
			Name:        "search_skills",
			Description: "Search the CyberStrike methodology skill catalog (SSRF, JWT, WSTG, MITRE, …). Returns brief hits; use get_skill for full body.",
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"q":        map[string]any{"type": "string", "description": "Search query"},
					"category": map[string]any{"type": "string"},
					"stage":    map[string]any{"type": "string", "description": "recon|exploit|post_exploit|report|all"},
					"limit":    map[string]any{"type": "integer", "description": "Default 8, max 15"},
				},
				"required": []any{"q"},
			},
			Safe: true,
		},
		{
			Name:        "get_skill",
			Description: "Load full methodology text for one skill id from search_skills.",
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"id":        map[string]any{"type": "string"},
					"max_chars": map[string]any{"type": "integer", "description": "Default 4000, max 8000"},
				},
				"required": []any{"id"},
			},
			Safe: true,
		},
		{
			Name:        "list_agents",
			Description: "List Talon agent modes (codename, focus, delegates).",
			Parameters:  map[string]any{"type": "object", "properties": map[string]any{}},
			Safe:        true,
		},
		{
			Name:        "list_playbooks",
			Description: "List playbooks (id, name, agent_mode, tags, description).",
			Parameters:  map[string]any{"type": "object", "properties": map[string]any{}},
			Safe:        true,
		},
		{
			Name:        "list_targets",
			Description: "List engagement targets registered in the ops platform.",
			Parameters:  map[string]any{"type": "object", "properties": map[string]any{}},
			Safe:        true,
		},
		{
			Name:        "service_health",
			Description: "Probe dependency health (core, postgres, ollama, onnx-slm, arsenal, msf, …). Use when user asks if the stack is up.",
			Parameters:  map[string]any{"type": "object", "properties": map[string]any{}},
			Safe:        true,
		},
		{
			Name:        "list_mcp_tools",
			Description: "List connected MCP servers and tool names available to full agent runs (read-only catalog; does not execute exploits).",
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"server": map[string]any{"type": "string", "description": "Optional filter: hexstrike|metasploit|lightpanda"},
				},
			},
			Safe: true,
		},
		{
			Name:        "intel_feed",
			Description: "Recent cross-run intel events (findings, completions).",
			Parameters: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"limit": map[string]any{"type": "integer", "description": "Default 20, max 50"},
				},
			},
			Safe: true,
		},
		{
			Name:        "runs_summary",
			Description: "Aggregate dashboard counters: total runs, by status, findings totals.",
			Parameters:  map[string]any{"type": "object", "properties": map[string]any{}},
			Safe:        true,
		},
	}
}

func slmToolSpecs() []llm.ToolSpec {
	cat := slmToolCatalog()
	out := make([]llm.ToolSpec, 0, len(cat))
	for _, t := range cat {
		out = append(out, llm.ToolSpec{
			Name:        t.Name,
			Description: t.Description,
			InputSchema: t.Parameters,
		})
	}
	return out
}

func slmSystemPrompt(extra string) string {
	var b strings.Builder
	b.WriteString(`You are Talon Assist — a local SLM copilot for the Talon pentest operations console.
You help operators understand runs, findings, skills, agents, and stack health.

TOOL USE (required when answering about live data):
Output exactly one tool call on its own line, then stop:
TOOL_CALL {"name":"<tool_name>","arguments":{...}}

Rules:
- Prefer tools over guessing run IDs, statuses, skill content, or health.
- After tool results appear, give a clear concise operator answer.
- Do NOT invent tool results. Do NOT claim exploits ran unless tool output says so.
- You cannot start attacks or call MSF/nmap from this chat — only the read tools below.
- Keep answers short and actionable for a terminal-style UI.

Available tools:
`)
	for _, t := range slmToolCatalog() {
		params, _ := json.Marshal(t.Parameters)
		fmt.Fprintf(&b, "- %s: %s\n  schema: %s\n", t.Name, t.Description, string(params))
	}
	if extra != "" {
		b.WriteString("\n")
		b.WriteString(extra)
	}
	return b.String()
}

// parseSLMToolCall extracts a TOOL_CALL {...} or <tool_call>...</tool_call> from model text.
func parseSLMToolCall(text string) (name string, args map[string]any, ok bool) {
	text = strings.TrimSpace(text)
	// Prefer explicit TOOL_CALL line.
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "TOOL_CALL") {
			raw := strings.TrimSpace(strings.TrimPrefix(line, "TOOL_CALL"))
			name, args, ok = decodeToolJSON(raw)
			if ok {
				return name, args, true
			}
		}
	}
	// XML-ish block.
	if i := strings.Index(text, "<tool_call>"); i >= 0 {
		j := strings.Index(text, "</tool_call>")
		if j > i {
			raw := strings.TrimSpace(text[i+len("<tool_call>") : j])
			return decodeToolJSON(raw)
		}
	}
	// Bare JSON object with name+arguments anywhere.
	if i := strings.Index(text, `{"name"`); i >= 0 {
		return decodeToolJSON(text[i:])
	}
	if i := strings.Index(text, `{"name":`); i >= 0 {
		return decodeToolJSON(text[i:])
	}
	return "", nil, false
}

func decodeToolJSON(raw string) (string, map[string]any, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", nil, false
	}
	// Truncate to first complete JSON object.
	if !strings.HasPrefix(raw, "{") {
		if i := strings.Index(raw, "{"); i >= 0 {
			raw = raw[i:]
		}
	}
	var envelope struct {
		Name      string         `json:"name"`
		Arguments map[string]any `json:"arguments"`
		Args      map[string]any `json:"args"`
	}
	if err := json.Unmarshal([]byte(raw), &envelope); err != nil {
		// Try to find matching braces.
		end := findJSONObjectEnd(raw)
		if end <= 0 {
			return "", nil, false
		}
		if err := json.Unmarshal([]byte(raw[:end]), &envelope); err != nil {
			return "", nil, false
		}
	}
	if envelope.Name == "" {
		return "", nil, false
	}
	args := envelope.Arguments
	if args == nil {
		args = envelope.Args
	}
	if args == nil {
		args = map[string]any{}
	}
	return envelope.Name, args, true
}

func findJSONObjectEnd(s string) int {
	depth := 0
	inStr := false
	esc := false
	for i, r := range s {
		if inStr {
			if esc {
				esc = false
				continue
			}
			if r == '\\' {
				esc = true
				continue
			}
			if r == '"' {
				inStr = false
			}
			continue
		}
		switch r {
		case '"':
			inStr = true
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return i + 1
			}
		}
	}
	return -1
}

func argInt(args map[string]any, key string, def, max int) int {
	v, ok := args[key]
	if !ok {
		return def
	}
	n := def
	switch t := v.(type) {
	case float64:
		n = int(t)
	case int:
		n = t
	case json.Number:
		if i, err := t.Int64(); err == nil {
			n = int(i)
		}
	case string:
		var i int
		if _, err := fmt.Sscanf(t, "%d", &i); err == nil {
			n = i
		}
	}
	if n <= 0 {
		n = def
	}
	if max > 0 && n > max {
		n = max
	}
	return n
}

func argString(args map[string]any, key string) string {
	v, _ := args[key].(string)
	return strings.TrimSpace(v)
}

// execSLMTool runs one curated tool against server state. Always returns a
// string result (error messages are content, not Go errors, so the model can react).
func (s *Server) execSLMTool(name string, args map[string]any) string {
	if args == nil {
		args = map[string]any{}
	}
	switch name {
	case "list_runs":
		limit := argInt(args, "limit", 10, 50)
		offset := argInt(args, "offset", 0, 0)
		runs, total, err := s.store.PaginatedList(limit, offset)
		if err != nil {
			return "error: " + err.Error()
		}
		type row struct {
			RunID         string `json:"run_id"`
			Target        string `json:"target"`
			Status        string `json:"status"`
			ToolCalls     int    `json:"tool_calls"`
			FindingsCount int    `json:"findings_count"`
			AgentMode     string `json:"agent_mode,omitempty"`
			StartedAt     string `json:"started_at,omitempty"`
		}
		out := make([]row, 0, len(runs))
		for _, r := range runs {
			started := ""
			if !r.StartedAt.IsZero() {
				started = r.StartedAt.UTC().Format(time.RFC3339)
			}
			out = append(out, row{
				RunID: r.RunID, Target: r.Target, Status: r.Status,
				ToolCalls: r.ToolCalls, FindingsCount: r.FindingsCount,
				AgentMode: r.AgentMode, StartedAt: started,
			})
		}
		return mustJSON(map[string]any{"total": total, "runs": out})

	case "get_run_status":
		id := argString(args, "run_id")
		if id == "" {
			return "error: run_id required"
		}
		sess, ok := s.store.Get(id)
		if !ok {
			return "error: run not found"
		}
		out := map[string]any{
			"run_id": id, "status": sess.Status, "target": sess.RunInput.TargetIP,
			"agent_mode": sess.RunInput.AgentMode, "tool_calls": len(sess.ToolLog),
		}
		if sess.Output != "" {
			o := sess.Output
			if len(o) > 1500 {
				o = o[:1500] + "…"
			}
			out["output_preview"] = o
		}
		if sess.JudgeSet {
			out["judge_verdict"] = sess.JudgeVerdict
		}
		if len(sess.Findings) > 0 {
			out["findings_count"] = len(sess.Findings)
		}
		return mustJSON(out)

	case "get_run_tools":
		id := argString(args, "run_id")
		if id == "" {
			return "error: run_id required"
		}
		log, ok := s.store.ToolLog(id)
		if !ok {
			return "error: run not found"
		}
		limit := argInt(args, "limit", 20, 100)
		if len(log) > limit {
			log = log[len(log)-limit:]
		}
		type row struct {
			Index    int    `json:"index"`
			ToolName string `json:"tool"`
			Args     any    `json:"args,omitempty"`
			Output   string `json:"output_preview,omitempty"`
		}
		rows := make([]row, 0, len(log))
		for _, t := range log {
			prev := t.Output
			if len(prev) > 400 {
				prev = prev[:400] + "…"
			}
			rows = append(rows, row{Index: t.Index, ToolName: t.ToolName, Args: t.Args, Output: prev})
		}
		return mustJSON(map[string]any{"run_id": id, "tools": rows, "count": len(rows)})

	case "get_findings":
		id := argString(args, "run_id")
		if id == "" {
			return "error: run_id required"
		}
		sess, ok := s.store.Get(id)
		if !ok {
			return "error: run not found"
		}
		type row struct {
			ID       string `json:"id"`
			Title    string `json:"title"`
			Severity string `json:"severity"`
			Status   string `json:"status"`
			Passed   bool   `json:"gate_passed"`
		}
		rows := make([]row, 0, len(sess.Findings))
		for _, f := range sess.Findings {
			rows = append(rows, row{
				ID: f.ID, Title: f.Title, Severity: f.Severity,
				Status: f.Status, Passed: f.Evidence.Passed,
			})
		}
		return mustJSON(map[string]any{"run_id": id, "findings": rows, "count": len(rows)})

	case "search_skills":
		q := argString(args, "q")
		if q == "" {
			return "error: q required"
		}
		limit := argInt(args, "limit", 8, 15)
		res := core.QuerySkills(core.SkillQuery{
			Q: q, Category: argString(args, "category"), Stage: argString(args, "stage"),
			Brief: true, Limit: limit,
		})
		type hit struct {
			ID, Name, Category, Stage string
		}
		hits := make([]hit, 0, len(res.Skills))
		for _, sk := range res.Skills {
			hits = append(hits, hit{ID: sk.ID, Name: sk.Name, Category: sk.Category, Stage: sk.Stage})
		}
		return mustJSON(map[string]any{"total": res.Total, "hits": hits})

	case "get_skill":
		id := argString(args, "id")
		if id == "" {
			return "error: id required"
		}
		maxChars := argInt(args, "max_chars", 4000, 8000)
		sk, ok := core.GetSkill(id)
		if !ok {
			return "error: skill not found"
		}
		body := sk.Body
		if len(body) > maxChars {
			body = body[:maxChars] + "\n…[truncated]"
		}
		return mustJSON(map[string]any{
			"id": sk.ID, "name": sk.Name, "category": sk.Category,
			"stage": sk.Stage, "body": body,
		})

	case "list_agents":
		agents := core.ListAgents()
		return mustJSON(map[string]any{"agents": agents, "count": len(agents)})

	case "list_playbooks":
		pbs := core.ListPlaybooks()
		type row struct {
			ID, Name, Codename, Description, AgentMode string
			Tags                                       []string
		}
		rows := make([]row, 0, len(pbs))
		for _, p := range pbs {
			rows = append(rows, row{
				ID: p.ID, Name: p.Name, Codename: p.Codename,
				Description: p.Description, AgentMode: p.AgentMode, Tags: p.Tags,
			})
		}
		return mustJSON(map[string]any{"playbooks": rows, "count": len(rows)})

	case "list_targets":
		if s.platform == nil {
			return mustJSON(map[string]any{"targets": []any{}, "note": "platform not configured"})
		}
		return mustJSON(map[string]any{"targets": s.platform.ListTargets()})

	case "service_health":
		// Same probes as /health/services — run concurrent so assist stays snappy.
		checks := []func() ServiceHealth{
			checkCore, s.checkPostgres, s.checkRedis, checkArsenal, checkMSF,
			checkRabbitMQ, checkOllama, checkONNXSLM,
		}
		type row struct {
			Name, Status, Detail, Endpoint string
			LatencyMS                      int64
		}
		rows := make([]row, len(checks))
		var wg sync.WaitGroup
		for i, c := range checks {
			wg.Add(1)
			go func(i int, c func() ServiceHealth) {
				defer wg.Done()
				h := c()
				rows[i] = row{
					Name: h.Name, Status: h.Status, Detail: h.Detail,
					Endpoint: h.Endpoint, LatencyMS: h.LatencyMS,
				}
			}(i, c)
		}
		wg.Wait()
		out := make([]map[string]any, 0, len(rows))
		for _, h := range rows {
			out = append(out, map[string]any{
				"name": h.Name, "status": h.Status, "detail": h.Detail,
				"latency_ms": h.LatencyMS, "endpoint": h.Endpoint,
			})
		}
		return mustJSON(map[string]any{"services": out, "probed_at": time.Now().UTC().Format(time.RFC3339)})

	case "list_mcp_tools":
		filter := argString(args, "server")
		if s.tools == nil {
			return mustJSON(map[string]any{"servers": []any{}, "note": "no MCP servers connected"})
		}
		servers := s.tools.Servers()
		type toolRow struct {
			Name string `json:"name"`
		}
		type srvRow struct {
			Name  string    `json:"name"`
			Tools []toolRow `json:"tools"`
		}
		var rows []srvRow
		for _, srv := range servers {
			if filter != "" && !strings.EqualFold(srv.Name, filter) {
				continue
			}
			tr := make([]toolRow, 0, len(srv.Tools))
			for _, tn := range srv.Tools {
				tr = append(tr, toolRow{Name: tn})
			}
			rows = append(rows, srvRow{Name: srv.Name, Tools: tr})
		}
		return mustJSON(map[string]any{"servers": rows})

	case "intel_feed":
		limit := argInt(args, "limit", 20, 50)
		events := s.store.IntelFeed(limit)
		return mustJSON(map[string]any{"events": events, "count": len(events)})

	case "runs_summary":
		sum := s.store.Summary()
		return mustJSON(sum)

	default:
		return fmt.Sprintf("error: unknown tool %q — use only catalog tools", name)
	}
}

func mustJSON(v any) string {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Sprintf("error: marshal: %v", err)
	}
	return string(b)
}

// knownSLMTool reports whether name is in the curated catalog.
func knownSLMTool(name string) bool {
	for _, t := range slmToolCatalog() {
		if t.Name == name {
			return true
		}
	}
	return false
}
