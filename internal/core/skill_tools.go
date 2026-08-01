package core

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/anubhavg-icpl/talon/internal/llm"
)

// skillToolSpecs exposes the CyberStrike skill catalog to agents at runtime
// (lazy load — full pack is not dumped into the system prompt).
func skillToolSpecs() []llm.ToolSpec {
	return []llm.ToolSpec{
		{
			Name: "skill_search",
			Description: "Search the CyberStrike methodology skill catalog (~7600 skills: WEB/WSTG, MITRE ATT&CK, CIS, NIST, attack-ssrf/jwt/idor, post-exploit, …). " +
				"Use before specialized testing to load the right methodology. Returns brief hits (id, name, category, stage, preview). Then call skill_get for full text.",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"q": map[string]any{
						"type":        "string",
						"description": "Search query (e.g. ssrf, jwt, kerberos, wstg-athn, privilege escalation)",
					},
					"category": map[string]any{
						"type":        "string",
						"description": "Optional category filter: WEB, mitre_attack, CIS_benchmarks, NIST, attack, cyberstrike, …",
					},
					"stage": map[string]any{
						"type":        "string",
						"description": "Optional stage: recon|exploit|post_exploit|report|all",
					},
					"limit": map[string]any{
						"type":        "integer",
						"description": "Max results (default 8, max 20)",
					},
				},
				"required": []any{"q"},
			},
		},
		{
			Name: "skill_get",
			Description: "Load the full methodology body for one skill by id (from skill_search). " +
				"Use this to inject precise testing steps into your plan before tool execution.",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"id": map[string]any{
						"type":        "string",
						"description": "Skill id returned by skill_search (e.g. disk-cyberstrike-attack-ssrf)",
					},
					"max_chars": map[string]any{
						"type":        "integer",
						"description": "Optional truncate body length (default 6000, max 12000)",
					},
				},
				"required": []any{"id"},
			},
		},
	}
}

func handleSkillSearch(args map[string]any, tr *tracker) (string, bool) {
	q, _ := args["q"].(string)
	q = strings.TrimSpace(q)
	if q == "" {
		return "skill_search: q is required", true
	}
	cat, _ := args["category"].(string)
	stage, _ := args["stage"].(string)
	limit := 8
	switch v := args["limit"].(type) {
	case float64:
		limit = int(v)
	case int:
		limit = v
	}
	if limit <= 0 {
		limit = 8
	}
	if limit > 20 {
		limit = 20
	}

	res := QuerySkills(SkillQuery{
		Q:        q,
		Category: cat,
		Stage:    stage,
		Brief:    true,
		Limit:    limit,
		Offset:   0,
	})

	type hit struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		Category string `json:"category"`
		Stage    string `json:"stage"`
		Path     string `json:"path,omitempty"`
	}
	hits := make([]hit, 0, len(res.Skills))
	for _, s := range res.Skills {
		hits = append(hits, hit{
			ID: s.ID, Name: s.Name, Category: s.Category, Stage: s.Stage, Path: s.Path,
		})
	}
	payload := map[string]any{
		"query":   q,
		"total":   res.Total,
		"showing": len(hits),
		"hits":    hits,
		"hint":    "Call skill_get with a hit id to load full methodology text.",
	}
	raw, _ := json.MarshalIndent(payload, "", "  ")
	out := string(raw)
	if tr != nil {
		tr.record("skill_search", args, fmt.Sprintf("total=%d showing=%d q=%q", res.Total, len(hits), q))
	}
	return out, false
}

func handleSkillGet(args map[string]any, tr *tracker) (string, bool) {
	id, _ := args["id"].(string)
	id = strings.TrimSpace(id)
	if id == "" {
		return "skill_get: id is required", true
	}
	maxChars := 6000
	switch v := args["max_chars"].(type) {
	case float64:
		maxChars = int(v)
	case int:
		maxChars = v
	}
	if maxChars <= 0 {
		maxChars = 6000
	}
	if maxChars > 12000 {
		maxChars = 12000
	}

	sk, ok := GetSkill(id)
	if !ok {
		// Fuzzy: search by id substring
		res := QuerySkills(SkillQuery{Q: id, Brief: false, Limit: 1})
		if res.Total == 0 || len(res.Skills) == 0 {
			return fmt.Sprintf("skill_get: skill %q not found — run skill_search first", id), true
		}
		// Prefer exact id match if any
		found := false
		for _, s := range ListSkills() {
			if s.ID == id {
				sk = s
				found = true
				break
			}
		}
		if !found {
			// Get full body for first search hit
			sk, ok = GetSkill(res.Skills[0].ID)
			if !ok {
				sk = res.Skills[0]
			}
		}
	}

	body := sk.Body
	truncated := false
	if len(body) > maxChars {
		body = body[:maxChars] + "\n…[truncated]"
		truncated = true
	}
	payload := map[string]any{
		"id":        sk.ID,
		"name":      sk.Name,
		"category":  sk.Category,
		"stage":     sk.Stage,
		"path":      sk.Path,
		"truncated": truncated,
		"body":      body,
	}
	raw, _ := json.MarshalIndent(payload, "", "  ")
	if tr != nil {
		tr.record("skill_get", args, fmt.Sprintf("id=%s bytes=%d", sk.ID, len(body)))
	}
	return string(raw), false
}

// withAgentTools adds findings + skill tools to a subagent tool list.
func withAgentTools(base []llm.ToolSpec) []llm.ToolSpec {
	out := append([]llm.ToolSpec{}, base...)
	out = append(out, findingToolSpecs()...)
	out = append(out, skillToolSpecs()...)
	return out
}
