package core

import (
	"fmt"
	"strings"
)

// KillChainLink is one edge in a derived attack path.
type KillChainLink struct {
	From     string `json:"from"`
	To       string `json:"to"`
	Severity string `json:"severity"`
	Reason   string `json:"reason"`
}

// KillChainAnalysis summarizes attack paths from structured findings
// (portable subset of CyberStrike killchain analysis).
type KillChainAnalysis struct {
	Chains    []KillChainLink `json:"chains"`
	NextSteps []string        `json:"next_steps"`
	Summary   string          `json:"summary"`
	MaxSev    string          `json:"max_severity"`
}

var stageOrder = []string{"recon", "exploit", "post_exploit", "codegen", "report"}

var stageNext = map[string][]string{
	"recon":        {"exploit", "web-testing", "credential-attack"},
	"exploit":      {"post_exploit", "privilege-escalation", "persistence"},
	"post_exploit": {"lateral-movement", "data-exfil", "cleanup"},
	"codegen":      {"post_exploit", "verify-rce"},
	"report":       {},
}

// AnalyzeKillChain derives chains from finding stages + severities.
func AnalyzeKillChain(findings []Finding) KillChainAnalysis {
	stagesSeen := map[string]string{} // stage → max severity
	var criticalTitles []string
	for _, f := range findings {
		st := f.Stage
		if st == "" {
			st = "other"
		}
		sev := strings.ToLower(f.Severity)
		if prev, ok := stagesSeen[st]; !ok || severityRank(sev) > severityRank(prev) {
			stagesSeen[st] = sev
		}
		if sev == SeverityCritical || sev == SeverityHigh {
			criticalTitles = append(criticalTitles, f.Title)
		}
	}

	var chains []KillChainLink
	for i := 0; i < len(stageOrder)-1; i++ {
		from, to := stageOrder[i], stageOrder[i+1]
		if _, ok := stagesSeen[from]; !ok {
			continue
		}
		if _, ok := stagesSeen[to]; !ok {
			continue
		}
		sev := maxSev(stagesSeen[from], stagesSeen[to])
		chains = append(chains, KillChainLink{
			From: from, To: to, Severity: sev,
			Reason: fmt.Sprintf("%s evidence led to %s activity", from, to),
		})
	}

	// Cross-links: recon open-port + RCE critical → high chain
	if _, hasRecon := stagesSeen["recon"]; hasRecon {
		if sev, hasExp := stagesSeen["exploit"]; hasExp && severityRank(sev) >= severityRank(SeverityHigh) {
			// ensure link exists
			found := false
			for _, c := range chains {
				if c.From == "recon" && c.To == "exploit" {
					found = true
					break
				}
			}
			if !found {
				chains = append(chains, KillChainLink{
					From: "recon", To: "exploit", Severity: sev,
					Reason: "recon surface enabled successful exploitation",
				})
			}
		}
	}

	next := map[string]bool{}
	for st := range stagesSeen {
		for _, n := range stageNext[st] {
			if _, done := stagesSeen[n]; !done {
				next[n] = true
			}
		}
	}
	// Heuristics from titles
	for _, t := range criticalTitles {
		tl := strings.ToLower(t)
		if strings.Contains(tl, "session") || strings.Contains(tl, "rce") {
			next["privilege-escalation"] = true
			next["lateral-movement"] = true
		}
		if strings.Contains(tl, "open port") {
			next["service-fingerprint"] = true
			next["exploit"] = true
		}
	}
	var nextSteps []string
	for n := range next {
		nextSteps = append(nextSteps, n)
	}

	maxS := SeverityInfo
	for _, s := range stagesSeen {
		if severityRank(s) > severityRank(maxS) {
			maxS = s
		}
	}

	var b strings.Builder
	b.WriteString("## Kill Chain Analysis\n\n")
	if len(chains) == 0 {
		b.WriteString("No multi-stage kill chains detected from current findings.\n")
	} else {
		for _, c := range chains {
			fmt.Fprintf(&b, "- **%s:** `%s` → `%s` — %s\n", strings.ToUpper(c.Severity), c.From, c.To, c.Reason)
		}
	}
	if len(nextSteps) > 0 {
		b.WriteString("\n**Suggested next steps:** ")
		b.WriteString(strings.Join(nextSteps, ", "))
		b.WriteString("\n")
	}

	return KillChainAnalysis{
		Chains:    chains,
		NextSteps: nextSteps,
		Summary:   b.String(),
		MaxSev:    maxS,
	}
}

func severityRank(s string) int {
	switch strings.ToLower(s) {
	case SeverityCritical:
		return 4
	case SeverityHigh:
		return 3
	case SeverityMedium:
		return 2
	case SeverityLow:
		return 1
	default:
		return 0
	}
}

func maxSev(a, b string) string {
	if severityRank(a) >= severityRank(b) {
		return a
	}
	return b
}
