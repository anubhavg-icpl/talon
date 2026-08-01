package core

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// Severity levels for structured findings (CyberStrike-compatible scale).
const (
	SeverityCritical = "critical"
	SeverityHigh     = "high"
	SeverityMedium   = "medium"
	SeverityLow      = "low"
	SeverityInfo     = "info"
)

// FindingStatus tracks triage lifecycle (inspired by CyberStrike vulnerability statuses).
const (
	FindingStatusNew      = "new"
	FindingStatusApproved = "approved"
	FindingStatusOpen     = "open"
	FindingStatusDup      = "duplicate"
)

// GateEvidence is one arm of the 3-gate confirmation protocol:
// baseline request/response → exploit request/response → measurable diff.
type GateEvidence struct {
	// Baseline is the benign / unmodified observation (Gate 1).
	Baseline string `json:"baseline,omitempty"`
	// Attack is the exploit / modified observation (Gate 2).
	Attack string `json:"attack,omitempty"`
	// Diff is a one-sentence comparison proving a measurable difference (Gate 3).
	Diff string `json:"diff,omitempty"`
	// Passed is true when all three gates have non-empty evidence.
	Passed bool `json:"passed"`
}

// Finding is a structured security finding extracted or reported during a run.
// Mirrors CyberStrike's report_vulnerability schema at a portable subset.
type Finding struct {
	ID                string       `json:"id"`
	Severity          string       `json:"severity"`
	Title             string       `json:"title"`
	Description       string       `json:"description"`
	CWEID             string       `json:"cwe_id,omitempty"`
	Endpoint          string       `json:"endpoint,omitempty"`
	AttackVector      string       `json:"attack_vector,omitempty"`
	StepsToReproduce  string       `json:"steps_to_reproduce,omitempty"`
	BusinessImpact    string       `json:"business_impact,omitempty"`
	Recommendation    string       `json:"recommendation,omitempty"`
	PoC               string       `json:"poc,omitempty"`
	Evidence          GateEvidence `json:"evidence"`
	Status            string       `json:"status"`
	Source            string       `json:"source"` // tool name, "judge", "orchestrator", etc.
	Stage             string       `json:"stage,omitempty"` // recon | exploit | post_exploit | codegen | report
	CreatedAt         time.Time    `json:"created_at"`
}

// FindingsSummary is a compact roll-up for status endpoints and list views.
type FindingsSummary struct {
	Total    int            `json:"total"`
	Critical int            `json:"critical"`
	High     int            `json:"high"`
	Medium   int            `json:"medium"`
	Low      int            `json:"low"`
	Info     int            `json:"info"`
	Confirmed int           `json:"confirmed"` // 3-gate passed
	BySeverity map[string]int `json:"by_severity,omitempty"`
}

// SummarizeFindings counts findings by severity and 3-gate confirmation.
func SummarizeFindings(findings []Finding) FindingsSummary {
	s := FindingsSummary{BySeverity: map[string]int{}}
	for _, f := range findings {
		s.Total++
		sev := strings.ToLower(f.Severity)
		s.BySeverity[sev]++
		switch sev {
		case SeverityCritical:
			s.Critical++
		case SeverityHigh:
			s.High++
		case SeverityMedium:
			s.Medium++
		case SeverityLow:
			s.Low++
		default:
			s.Info++
		}
		if f.Evidence.Passed {
			s.Confirmed++
		}
	}
	return s
}

// ValidateThreeGate checks the CyberStrike-style confirmation protocol:
// Gate 1 baseline, Gate 2 attack, Gate 3 measurable diff.
func ValidateThreeGate(e GateEvidence) (GateEvidence, []string) {
	var violations []string
	e.Baseline = strings.TrimSpace(e.Baseline)
	e.Attack = strings.TrimSpace(e.Attack)
	e.Diff = strings.TrimSpace(e.Diff)
	if e.Baseline == "" {
		violations = append(violations, "gate1_baseline: missing baseline observation")
	}
	if e.Attack == "" {
		violations = append(violations, "gate2_attack: missing attack observation")
	}
	if e.Diff == "" {
		violations = append(violations, "gate3_diff: missing measurable difference")
	} else if sameObservation(e.Baseline, e.Attack) && !strings.Contains(strings.ToLower(e.Diff), "differ") {
		// Soft check: if baseline≈attack and diff doesn't claim difference, fail gate 3.
		if len(e.Baseline) > 20 && e.Baseline == e.Attack {
			violations = append(violations, "gate3_diff: baseline and attack observations are identical")
		}
	}
	e.Passed = len(violations) == 0
	return e, violations
}

func sameObservation(a, b string) bool {
	return strings.TrimSpace(a) == strings.TrimSpace(b)
}

// dedupeKey builds a stable key for near-duplicate suppression
// (endpoint + attack_vector + title class), inspired by CyberStrike.
func dedupeKey(f Finding) string {
	ep := strings.ToLower(strings.TrimSpace(f.Endpoint))
	av := strings.ToLower(strings.TrimSpace(f.AttackVector))
	title := strings.ToLower(strings.TrimSpace(f.Title))
	if ep == "" {
		ep = "no-endpoint"
	}
	if av == "" {
		av = "other"
	}
	return ep + "|" + av + "|" + title
}

// DedupFindings keeps the first of each dedupe key; later duplicates are dropped.
func DedupFindings(in []Finding) []Finding {
	seen := make(map[string]bool, len(in))
	out := make([]Finding, 0, len(in))
	for _, f := range in {
		k := dedupeKey(f)
		if seen[k] {
			continue
		}
		seen[k] = true
		out = append(out, f)
	}
	return out
}

var (
	reSessionCreated = regexp.MustCompile(`(?i)session\s+(\d+)\s+created|meterpreter\s+session\s+(\d+)|command\s+shell\s+session\s+(\d+)`)
	reOpenPorts      = regexp.MustCompile(`(?i)(\d+)/tcp\s+open\s+(\S+)`)
	reNucleiHit      = regexp.MustCompile(`(?i)\[(critical|high|medium|low|info)\]\s+(.+)`)
	reWhoami         = regexp.MustCompile(`(?im)^(uid=\d+|root|www-data|administrator|nt authority)\\?\S*`)
	reCVE            = regexp.MustCompile(`(?i)CVE-\d{4}-\d{4,}`)
)

// ExtractFindings builds structured findings from a completed run's tool log,
// final message, and judge verdict. This is the portable equivalent of
// CyberStrike's report_vulnerability pipeline without requiring live HTTP
// proxy testers — it synthesizes findings from Talon's recon→exploit→post
// tool evidence.
func ExtractFindings(input RunInput, toolLog []ToolCallRecord, finalMessage string, judgeVerdict bool, judgeSet bool) []Finding {
	now := time.Now().UTC()
	var findings []Finding
	idSeq := 0
	nextID := func() string {
		idSeq++
		return fmt.Sprintf("FIND-%03d", idSeq)
	}

	// --- Recon findings: open ports / services ---
	for _, rec := range toolLog {
		stage := stageForTool(rec.ToolName)
		out := rec.Output
		if out == "" {
			continue
		}

		switch {
		case rec.ToolName == "nmap_scan" || rec.ToolName == "rustscan_fast_scan":
			matches := reOpenPorts.FindAllStringSubmatch(out, -1)
			for _, m := range matches {
				if len(m) < 3 {
					continue
				}
				port, svc := m[1], m[2]
				f := Finding{
					ID:          nextID(),
					Severity:    SeverityInfo,
					Title:       fmt.Sprintf("Open port %s/tcp (%s)", port, svc),
					Description: fmt.Sprintf("Port scan of %s reported %s/tcp open running %s.", input.TargetIP, port, svc),
					Endpoint:    fmt.Sprintf("%s:%s", input.TargetIP, port),
					AttackVector: "other",
					Source:      rec.ToolName,
					Stage:       stage,
					Status:      FindingStatusApproved,
					CreatedAt:   now,
					Evidence: GateEvidence{
						Baseline: "host presumed filtered/closed before scan",
						Attack:   fmt.Sprintf("%s/tcp open %s", port, svc),
						Diff:     "service is reachable and responds on the scanned port",
						Passed:   true,
					},
					Recommendation: "Restrict exposure of unnecessary services; verify patch level.",
				}
				findings = append(findings, f)
			}

		case rec.ToolName == "nuclei_scan":
			matches := reNucleiHit.FindAllStringSubmatch(out, -1)
			for _, m := range matches {
				if len(m) < 3 {
					continue
				}
				sev, title := strings.ToLower(m[1]), strings.TrimSpace(m[2])
				if sev == "" {
					sev = SeverityMedium
				}
				f := Finding{
					ID:           nextID(),
					Severity:     sev,
					Title:        "Nuclei: " + truncate(title, 120),
					Description:  title,
					Endpoint:     input.TargetIP,
					AttackVector: "other",
					Source:       rec.ToolName,
					Stage:        stage,
					Status:       FindingStatusApproved,
					CreatedAt:    now,
					Evidence: GateEvidence{
						Baseline: "template match requires live target response",
						Attack:   truncate(out, 500),
						Diff:     "nuclei template matched on target response",
						Passed:   true,
					},
					PoC: truncate(out, 800),
				}
				if cve := reCVE.FindString(title); cve != "" {
					f.CWEID = cve // store CVE reference in CWEID field when no CWE known
				}
				findings = append(findings, f)
			}

		case rec.ToolName == "run_exploit" || rec.ToolName == "custom_exploit":
			if reSessionCreated.MatchString(out) {
				m := reSessionCreated.FindStringSubmatch(out)
				sid := firstNonEmpty(m[1:]...)
				f := Finding{
					ID:             nextID(),
					Severity:       SeverityCritical,
					Title:          "Remote code execution / session established",
					Description:    fmt.Sprintf("Exploit against %s established interactive session %s.", input.TargetIP, sid),
					CWEID:          "CWE-94",
					Endpoint:       input.TargetIP,
					AttackVector:   "other",
					Source:         rec.ToolName,
					Stage:          stage,
					Status:         FindingStatusApproved,
					CreatedAt:      now,
					BusinessImpact: "Full remote code execution on the target host.",
					Recommendation: "Patch the vulnerable service immediately; isolate the host; rotate credentials.",
					StepsToReproduce: fmt.Sprintf(
						"1. Target %s (CVE=%s service=%s)\n2. Run exploit module via %s\n3. Observe session %s created",
						input.TargetIP, input.CVEID, input.ServiceName, rec.ToolName, sid,
					),
					Evidence: GateEvidence{
						Baseline: "no active session before exploit",
						Attack:   truncate(out, 600),
						Diff:     fmt.Sprintf("session %s created — interactive shell/meterpreter available", sid),
						Passed:   true,
					},
					PoC: truncate(out, 1000),
				}
				if input.CVEID != "" {
					f.Title = fmt.Sprintf("RCE via %s — session established", input.CVEID)
				}
				findings = append(findings, f)
			} else if strings.Contains(strings.ToLower(out), "exploit failed") ||
				strings.Contains(strings.ToLower(out), "no session") {
				// Failed exploit attempt — info only, helps report methodology coverage.
				f := Finding{
					ID:          nextID(),
					Severity:    SeverityInfo,
					Title:       "Exploit attempt failed (no session)",
					Description: fmt.Sprintf("Module via %s did not produce a session on %s.", rec.ToolName, input.TargetIP),
					Endpoint:    input.TargetIP,
					Source:      rec.ToolName,
					Stage:       stage,
					Status:      FindingStatusOpen,
					CreatedAt:   now,
					Evidence: GateEvidence{
						Baseline: "pre-exploit state: no session",
						Attack:   truncate(out, 400),
						Diff:     "no session created — exploit unsuccessful",
						Passed:   true,
					},
				}
				findings = append(findings, f)
			}

		case rec.ToolName == "send_session_command":
			// Proof-of-compromise command output elevates / confirms RCE.
			proof := extractProof(out)
			if proof != "" {
				f := Finding{
					ID:             nextID(),
					Severity:       SeverityCritical,
					Title:          "Post-exploit proof of compromise",
					Description:    "Command execution on established session returned host identity evidence.",
					CWEID:          "CWE-78",
					Endpoint:       input.TargetIP,
					AttackVector:   "other",
					Source:         rec.ToolName,
					Stage:          stage,
					Status:         FindingStatusApproved,
					CreatedAt:      now,
					BusinessImpact: "Attacker can execute arbitrary commands as the session user.",
					Recommendation: "Treat host as compromised; rebuild from known-good image; audit lateral movement.",
					Evidence: GateEvidence{
						Baseline: "session exists but identity unconfirmed",
						Attack:   proof,
						Diff:     "command output proves code execution on the target",
						Passed:   true,
					},
					PoC: proof,
				}
				findings = append(findings, f)
			}
		}
	}

	// --- Judge-confirmed compromise when tools were ambiguous ---
	if judgeSet && judgeVerdict {
		hasCritical := false
		for _, f := range findings {
			if f.Severity == SeverityCritical && f.Evidence.Passed {
				hasCritical = true
				break
			}
		}
		if !hasCritical {
			findings = append(findings, Finding{
				ID:             nextID(),
				Severity:       SeverityCritical,
				Title:          "Judge-confirmed compromise",
				Description:    "Independent judge model verified that the run output demonstrates successful exploitation.",
				Endpoint:       input.TargetIP,
				Source:         "judge",
				Stage:          "report",
				Status:         FindingStatusApproved,
				CreatedAt:      now,
				BusinessImpact: "Compromise confirmed by secondary model review.",
				Evidence: GateEvidence{
					Baseline: "objective: verify RCE/session compromise",
					Attack:   truncate(finalMessage, 600),
					Diff:     "judge returned TRUE — objective met",
					Passed:   true,
				},
				PoC: truncate(finalMessage, 800),
			})
		}
	}

	// Tag CVE on primary critical findings when known.
	if input.CVEID != "" {
		for i := range findings {
			if findings[i].Severity == SeverityCritical && findings[i].CWEID == "" {
				// Keep CWE if set; otherwise store CVE in description context.
				if !strings.Contains(findings[i].Description, input.CVEID) {
					findings[i].Description += " CVE: " + input.CVEID + "."
				}
			}
		}
	}

	return DedupFindings(findings)
}

func stageForTool(name string) string {
	switch name {
	case "nmap_scan", "smbmap_scan", "nuclei_scan", "rustscan_fast_scan", "arp_scan_discovery":
		return "recon"
	case "list_exploits", "list_payloads", "generate_payload", "run_exploit",
		"run_auxiliary_module", "sqlmap_scan", "hydra_attack", "responder_credential_harvest":
		return "exploit"
	case "list_active_sessions", "terminate_session", "send_session_command", "run_post_module":
		return "post_exploit"
	case "custom_exploit":
		return "codegen"
	default:
		if strings.HasPrefix(name, "delegate_") {
			return strings.TrimPrefix(name, "delegate_")
		}
		return "other"
	}
}

func extractProof(out string) string {
	out = strings.TrimSpace(out)
	if out == "" {
		return ""
	}
	lower := strings.ToLower(out)
	// Positive indicators of command execution proof.
	indicators := []string{"uid=", "gid=", "www-data", "nt authority", "hostname", "linux", "windows", "darwin", "microsoft windows"}
	for _, ind := range indicators {
		if strings.Contains(lower, ind) {
			return truncate(out, 500)
		}
	}
	if reWhoami.MatchString(out) {
		return truncate(out, 500)
	}
	// Short non-error shell responses (hostname, whoami username alone).
	if len(out) < 200 && !strings.Contains(lower, "error") && !strings.Contains(lower, "failed") &&
		!strings.Contains(lower, "unknown") && strings.Count(out, "\n") < 5 {
		// Likely a single-line hostname/user response.
		if len(strings.Fields(out)) <= 4 && len(out) > 0 {
			return truncate(out, 200)
		}
	}
	return ""
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return "?"
}

func truncate(s string, n int) string {
	s = strings.TrimSpace(s)
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}
