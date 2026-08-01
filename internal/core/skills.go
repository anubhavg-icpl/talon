package core

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

// Skill is a methodology document (builtin or CyberStrike-imported disk pack).
type Skill struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Stage    string `json:"stage"`
	Category string `json:"category,omitempty"`
	Body     string `json:"body,omitempty"`
	Source   string `json:"source,omitempty"` // builtin | disk
	Path     string `json:"path,omitempty"`   // relative disk path
}

var (
	diskSkillsOnce sync.Once
	diskSkills     []Skill
	skillByID      map[string]Skill
)

const (
	maxInjectDiskSkills = 12
	maxInjectBodyChars  = 2200
)

// builtinSkills is the always-on seed pack (injected fully into agents).
var builtinSkills = []Skill{
	{
		ID: "3gate-protocol", Name: "3-Gate Confirmation Protocol", Stage: "all",
		Category: "core", Source: "builtin",
		Body: `## 3-Gate Confirmation Protocol
**Gate 1 — Baseline** **Gate 2 — Attack** **Gate 3 — Diff**
Use report_finding with baseline/attack/diff. Never claim success without measurable difference.`,
	},
	{
		ID: "recon-methodology", Name: "Recon Methodology", Stage: "recon",
		Category: "core", Source: "builtin",
		Body: `## Recon Methodology
Focused nmap, optional nuclei, factual findings only. Max 2 nmap + 1 nuclei unless failure.`,
	},
	{
		ID: "exploit-methodology", Name: "Exploit Methodology", Stage: "exploit",
		Category: "core", Source: "builtin",
		Body: `## Exploit Methodology
Search → configure LHOST/LPORT → execute → verify session. Budget 3 attempts. report_finding on success.`,
	},
	{
		ID: "post-exploit-proof", Name: "Post-Exploit Proof", Stage: "post_exploit",
		Category: "core", Source: "builtin",
		Body: `## Post-Exploit Proof
Session identity: whoami/id/hostname/sysinfo. Raw tool output as proof.`,
	},
	{
		ID: "codegen-fallback", Name: "Codegen Fallback", Stage: "codegen",
		Category: "core", Source: "builtin",
		Body: `## Codegen Fallback
When MSF fails: Python reverse shell to LHOST:LPORT. Full sandbox output for judge.`,
	},
	{
		ID: "report-triage", Name: "Report Severity Triage", Stage: "report",
		Category: "core", Source: "builtin",
		Body: `## Severity
critical=RCE/session; high=ATO; medium=conditional; low=hardening; info=recon.`,
	},
	{
		ID: "web-owasp-lite", Name: "OWASP Web Testing Lite", Stage: "exploit",
		Category: "core", Source: "builtin",
		Body: `## Web / API
3-gate on every claim. Full WSTG/attack skills available via skill catalog.`,
	},
	{
		ID: "report-finding-tool", Name: "report_finding Tool Usage", Stage: "all",
		Category: "core", Source: "builtin",
		Body: `## report_finding / triage_finding
Required: severity, title, description. Recommended: baseline, attack, diff.`,
	},
	{
		ID: "network-lateral", Name: "Network / Lateral Emphasis", Stage: "exploit",
		Category: "core", Source: "builtin",
		Body: `## Network
arp/rustscan → careful hydra → session → post proof. Authorized scope only.`,
	},
	{
		ID: "idor-authz-lite", Name: "IDOR / AuthZ Lite", Stage: "exploit",
		Category: "core", Source: "builtin",
		Body: `## Access Control
Low-priv baseline → swap IDs/roles → only report extra access.`,
	},
	{
		ID: "injection-safe", Name: "Safe Injection Testing", Stage: "exploit",
		Category: "core", Source: "builtin",
		Body: `## Injection (minimum impact)
Boolean/time/error/UNION version only. Never DROP/TRUNCATE/--os-shell.`,
	},
	{
		ID: "auth-required", Name: "Authorization Gate", Stage: "all",
		Category: "core", Source: "builtin",
		Body: `## Authorization
Only authorized targets. Stay in scope. Prefer minimum-impact PoCs.`,
	},
}

func skillsSearchDirs() []string {
	var dirs []string
	if d := os.Getenv("TALON_SKILLS_DIR"); d != "" {
		dirs = append(dirs, d)
	}
	for _, d := range []string{"skills", "/app/skills", "./skills"} {
		dirs = append(dirs, d)
	}
	// Walk up from CWD so `go test` from internal/* and CLI from any subdir still find skills/.
	if wd, err := os.Getwd(); err == nil {
		cur := wd
		for i := 0; i < 8; i++ {
			dirs = append(dirs, filepath.Join(cur, "skills"))
			parent := filepath.Dir(cur)
			if parent == cur {
				break
			}
			cur = parent
		}
	}
	return dirs
}

// loadDiskSkills walks skills/ recursively (CyberStrike catalog + flat pack).
func loadDiskSkills() []Skill {
	diskSkillsOnce.Do(func() {
		seen := map[string]bool{}
		skillByID = make(map[string]Skill)
		for _, b := range builtinSkills {
			skillByID[b.ID] = b
		}
		for _, dir := range skillsSearchDirs() {
			info, err := os.Stat(dir)
			if err != nil || !info.IsDir() {
				continue
			}
			_ = filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
				if err != nil {
					return nil
				}
				if d.IsDir() {
					return nil
				}
				name := d.Name()
				if !strings.HasSuffix(strings.ToLower(name), ".md") {
					return nil
				}
				if strings.EqualFold(name, "README.md") || strings.EqualFold(name, "SKILL_GUIDE.md") {
					return nil
				}
				data, err := os.ReadFile(path)
				if err != nil {
					return nil
				}
				rel, _ := filepath.Rel(dir, path)
				sk := parseSkillFile(rel, string(data))
				if sk.ID == "" || seen[sk.ID] {
					return nil
				}
				seen[sk.ID] = true
				diskSkills = append(diskSkills, sk)
				skillByID[sk.ID] = sk
				return nil
			})
			// Prefer the first directory tree that actually contains skills
			// (do not merge multiple trees — avoids double-counting).
			if len(diskSkills) > 0 {
				break
			}
		}
		sort.Slice(diskSkills, func(i, j int) bool {
			if diskSkills[i].Category != diskSkills[j].Category {
				return diskSkills[i].Category < diskSkills[j].Category
			}
			return diskSkills[i].Name < diskSkills[j].Name
		})
	})
	return diskSkills
}

func parseSkillFile(relPath, body string) Skill {
	relPath = filepath.ToSlash(relPath)
	base := strings.TrimSuffix(filepath.Base(relPath), filepath.Ext(relPath))
	// Prefer category from top folder under skills/ (e.g. cyberstrike/WEB/...)
	parts := strings.Split(relPath, "/")
	category := "general"
	if len(parts) >= 2 {
		// skills/cyberstrike/WEB/foo.md → category WEB (or cyberstrike if only one level)
		if parts[0] == "cyberstrike" && len(parts) >= 3 {
			category = parts[1]
		} else if parts[0] == "cyberstrike" {
			category = "cyberstrike"
		} else {
			category = parts[0]
		}
	} else if strings.HasPrefix(base, "web-") {
		category = "WEB"
	} else if strings.HasPrefix(base, "attack-") {
		category = "attack"
	}

	stage := "exploit"
	// Headers: # stage: X  and  # category: Y
	lines := strings.Split(body, "\n")
	start := 0
	for i := 0; i < len(lines) && i < 5; i++ {
		line := strings.TrimSpace(lines[i])
		low := strings.ToLower(line)
		if strings.HasPrefix(low, "# stage:") {
			stage = strings.TrimSpace(strings.TrimPrefix(low, "# stage:"))
			start = i + 1
			continue
		}
		if strings.HasPrefix(low, "# category:") {
			category = strings.TrimSpace(line[len("# category:"):])
			start = i + 1
			continue
		}
		if line == "" {
			start = i + 1
			continue
		}
		break
	}
	if start > 0 && start < len(lines) {
		body = strings.Join(lines[start:], "\n")
	}
	body = strings.TrimSpace(body)
	if len(body) < 20 {
		return Skill{}
	}

	// Stage heuristics
	blob := strings.ToLower(relPath + " " + category + " " + base)
	switch {
	case strings.Contains(blob, "recon"), strings.Contains(blob, "discovery"),
		strings.Contains(blob, "wstg-info"), strings.Contains(blob, "wstg-conf"):
		if stage == "exploit" {
			stage = "recon"
		}
	case strings.Contains(blob, "mitre"), strings.Contains(blob, "postexploit"),
		strings.Contains(blob, "post-exploit"), strings.Contains(blob, "kerberos"),
		strings.Contains(blob, "ad-security"), strings.Contains(blob, "ebpf"):
		if stage == "exploit" {
			stage = "post_exploit"
		}
	case strings.Contains(blob, "nist"), strings.Contains(blob, "cis_benchmark"):
		if stage == "exploit" {
			stage = "report"
		}
	}

	// Stable unique id from relative path
	id := "disk-" + strings.TrimSuffix(relPath, ".md")
	id = strings.ReplaceAll(id, "/", "-")
	id = strings.ReplaceAll(id, " ", "_")

	// Human name from filename
	name := base
	name = strings.ReplaceAll(name, "_", " ")
	name = strings.ReplaceAll(name, "-", " ")

	return Skill{
		ID:       id,
		Name:     name,
		Stage:    stage,
		Category: category,
		Body:     body,
		Source:   "disk",
		Path:     relPath,
	}
}

// SkillsForStage returns skills for a pipeline stage (plus "all").
func SkillsForStage(stage string) []Skill {
	stage = strings.ToLower(strings.TrimSpace(stage))
	var out []Skill
	for _, s := range builtinSkills {
		if s.Stage == "all" || s.Stage == stage {
			out = append(out, s)
		}
	}
	for _, s := range loadDiskSkills() {
		if s.Stage == "all" || s.Stage == stage {
			out = append(out, s)
		}
	}
	return out
}

// InjectSkills appends a bounded set of skills into a system prompt.
func InjectSkills(basePrompt, stage string) string {
	skills := SkillsForStage(stage)
	if len(skills) == 0 {
		return basePrompt
	}
	var b strings.Builder
	b.WriteString(basePrompt)
	b.WriteString("\n\n---\n# Injected Security Skills\n")
	b.WriteString(fmt.Sprintf("_Catalog has %d total skills via GET /skills — injecting core + sample for stage %q._\n",
		len(ListSkills()), stage))

	diskN := 0
	for _, s := range skills {
		body := s.Body
		if s.Source == "disk" {
			// Prefer high-value categories for injection samples
			if !isPriorityCategory(s.Category) && diskN >= maxInjectDiskSkills/2 {
				continue
			}
			if diskN >= maxInjectDiskSkills {
				continue
			}
			diskN++
			if len(body) > maxInjectBodyChars {
				body = body[:maxInjectBodyChars] + "\n…[truncated — open Skills UI for full text]"
			}
		}
		b.WriteString("\n### ")
		b.WriteString(s.Name)
		if s.Category != "" {
			b.WriteString(" [")
			b.WriteString(s.Category)
			b.WriteString("]")
		}
		b.WriteString(" (`")
		b.WriteString(s.ID)
		b.WriteString("`)\n\n")
		b.WriteString(body)
		b.WriteString("\n")
	}
	return b.String()
}

func isPriorityCategory(c string) bool {
	c = strings.ToLower(c)
	switch c {
	case "core", "attack", "web", "cyberstrike", "ad-security", "recon-methodology",
		"kerberos-attacks", "cicd-attacks", "aws-postexploit", "azure-postexploit",
		"k8s-postexploit", "windows-postexploit", "macos-postexploit", "ebpf-attacks":
		return true
	}
	if strings.HasPrefix(c, "attack") {
		return true
	}
	return false
}

// ListSkills returns the full catalog (builtin + disk).
func ListSkills() []Skill {
	out := make([]Skill, 0, len(builtinSkills)+len(loadDiskSkills()))
	out = append(out, builtinSkills...)
	out = append(out, loadDiskSkills()...)
	return out
}

// GetSkill returns one skill by id (loads catalog if needed).
func GetSkill(id string) (Skill, bool) {
	_ = loadDiskSkills()
	if skillByID == nil {
		return Skill{}, false
	}
	s, ok := skillByID[id]
	return s, ok
}

// SkillQuery filters/paginates the catalog for the API.
type SkillQuery struct {
	Stage    string
	Category string
	Q        string // free-text search
	Brief    bool   // omit body
	Limit    int
	Offset   int
}

// SkillListResult is a paginated catalog page.
type SkillListResult struct {
	Skills     []Skill         `json:"skills"`
	Count      int             `json:"count"`      // page size
	Total      int             `json:"total"`      // matching total
	Offset     int             `json:"offset"`
	Limit      int             `json:"limit"`
	Stats      map[string]int  `json:"stats"`
	Categories []CategoryCount `json:"categories"`
}

// CategoryCount is a category roll-up for the Skills UI sidebar.
type CategoryCount struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

// QuerySkills filters and paginates skills.
func QuerySkills(q SkillQuery) SkillListResult {
	all := ListSkills()
	if q.Limit <= 0 {
		q.Limit = 50
	}
	if q.Limit > 500 {
		q.Limit = 500
	}
	if q.Offset < 0 {
		q.Offset = 0
	}
	stage := strings.ToLower(strings.TrimSpace(q.Stage))
	cat := strings.TrimSpace(q.Category)
	needle := strings.ToLower(strings.TrimSpace(q.Q))

	// Category counts over full catalog (for sidebar)
	catMap := map[string]int{}
	for _, s := range all {
		c := s.Category
		if c == "" {
			c = "general"
		}
		catMap[c]++
	}

	var matched []Skill
	for _, s := range all {
		if stage != "" && stage != "all" && s.Stage != stage && s.Stage != "all" {
			continue
		}
		if cat != "" && !strings.EqualFold(s.Category, cat) {
			continue
		}
		if needle != "" {
			blob := strings.ToLower(s.ID + " " + s.Name + " " + s.Category + " " + s.Path + " " + s.Body)
			if !strings.Contains(blob, needle) {
				continue
			}
		}
		if q.Brief {
			s.Body = ""
		}
		matched = append(matched, s)
	}

	total := len(matched)
	if q.Offset > total {
		q.Offset = total
	}
	end := q.Offset + q.Limit
	if end > total {
		end = total
	}
	page := matched[q.Offset:end]

	cats := make([]CategoryCount, 0, len(catMap))
	for name, n := range catMap {
		cats = append(cats, CategoryCount{Name: name, Count: n})
	}
	sort.Slice(cats, func(i, j int) bool {
		if cats[i].Count != cats[j].Count {
			return cats[i].Count > cats[j].Count
		}
		return cats[i].Name < cats[j].Name
	})

	return SkillListResult{
		Skills:     page,
		Count:      len(page),
		Total:      total,
		Offset:     q.Offset,
		Limit:      q.Limit,
		Stats:      SkillStats(),
		Categories: cats,
	}
}

// SkillStats returns catalog counts.
func SkillStats() map[string]int {
	all := ListSkills()
	by := map[string]int{}
	for _, s := range all {
		by["total"]++
		by[s.Stage]++
		if s.Source != "" {
			by["src_"+s.Source]++
		}
	}
	return by
}
