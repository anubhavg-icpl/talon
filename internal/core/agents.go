package core

import "strings"

// Agent modes (CyberStrike-inspired specialists). Operators pick one on start;
// the orchestrator adjusts available delegates and prompt emphasis.
const (
	AgentModeFull     = "full"     // recon → exploit → post → codegen → report (default)
	AgentModeRecon    = "recon"    // recon + report only
	AgentModeExploit  = "exploit"  // exploit-focused (light recon + exploit + post)
	AgentModeWeb      = "web"      // web/OWASP emphasis + nuclei/sqlmap path
	AgentModeNetwork  = "network"  // internal network / AD-style emphasis
	AgentModePost     = "post"     // post-exploit only (existing session assumed)
)

// AgentInfo describes a selectable specialist for the dashboard / API.
type AgentInfo struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Codename    string   `json:"codename"`
	Focus       string   `json:"focus"`
	Description string   `json:"description"`
	Delegates   []string `json:"delegates"`
}

// ListAgents returns the builtin specialist catalog.
func ListAgents() []AgentInfo {
	return []AgentInfo{
		{
			ID: AgentModeFull, Name: "Full Pipeline", Codename: "COMMANDER",
			Focus: "general", Description: "Complete recon → exploit → post-exploit → codegen → judge pipeline.",
			Delegates: []string{"delegate_recon", "delegate_exploit", "delegate_post_exploit", "delegate_codegen", "delegate_report"},
		},
		{
			ID: AgentModeRecon, Name: "Recon Specialist", Codename: "GHOST",
			Focus: "recon", Description: "Service/CVE verification only — nmap, nuclei, smbmap. No exploitation.",
			Delegates: []string{"delegate_recon", "delegate_report"},
		},
		{
			ID: AgentModeExploit, Name: "Exploit Specialist", Codename: "STRIKER",
			Focus: "exploit", Description: "Light recon then aggressive Metasploit/module exploitation + post-proof.",
			Delegates: []string{"delegate_recon", "delegate_exploit", "delegate_post_exploit", "delegate_codegen", "delegate_report"},
		},
		{
			ID: AgentModeWeb, Name: "Web Application", Codename: "STRIKER-WEB",
			Focus: "web", Description: "OWASP-oriented: nuclei, sqlmap, web-focused exploit modules + 3-gate reporting.",
			Delegates: []string{"delegate_recon", "delegate_exploit", "delegate_post_exploit", "delegate_report"},
		},
		{
			ID: AgentModeNetwork, Name: "Internal Network", Codename: "PHANTOM",
			Focus: "network", Description: "Network-first: arp/rustscan/hydra/responder emphasis + lateral-style post-exploit.",
			Delegates: []string{"delegate_recon", "delegate_exploit", "delegate_post_exploit", "delegate_report"},
		},
		{
			ID: AgentModePost, Name: "Post-Exploit", Codename: "CIPHER",
			Focus: "post_exploit", Description: "Proof-of-compromise only — assume a session already exists.",
			Delegates: []string{"delegate_post_exploit", "delegate_report"},
		},
	}
}

// NormalizeAgentMode maps free-form input to a known mode (default: full).
func NormalizeAgentMode(mode string) string {
	m := strings.ToLower(strings.TrimSpace(mode))
	switch m {
	case "", AgentModeFull, "default", "commander", "cyberstrike":
		return AgentModeFull
	case AgentModeRecon, "ghost", "explore":
		return AgentModeRecon
	case AgentModeExploit, "striker":
		return AgentModeExploit
	case AgentModeWeb, "web-application", "webapp", "owasp":
		return AgentModeWeb
	case AgentModeNetwork, "internal", "internal-network", "ad", "phantom":
		return AgentModeNetwork
	case AgentModePost, "post_exploit", "post-exploit", "cipher":
		return AgentModePost
	default:
		return AgentModeFull
	}
}

// AgentModePrompt returns extra system guidance for the orchestrator for a mode.
func AgentModePrompt(mode string) string {
	switch NormalizeAgentMode(mode) {
	case AgentModeRecon:
		return "AGENT MODE: RECON ONLY. Run delegate_recon, then delegate_report. Do NOT exploit."
	case AgentModeExploit:
		return "AGENT MODE: EXPLOIT. Quick recon confirm then prioritize delegate_exploit and post_exploit. Use codegen if modules fail."
	case AgentModeWeb:
		return "AGENT MODE: WEB APPLICATION. Prefer nuclei_scan, sqlmap_scan, and web-relevant exploits. Apply OWASP thinking and 3-gate evidence for every claimed vuln."
	case AgentModeNetwork:
		return "AGENT MODE: INTERNAL NETWORK. Prefer arp_scan_discovery, rustscan_fast_scan, hydra_attack, responder_credential_harvest, and lateral-style post-exploit proof."
	case AgentModePost:
		return "AGENT MODE: POST-EXPLOIT ONLY. Skip recon/exploit. Use delegate_post_exploit for proof, then report."
	default:
		return "AGENT MODE: FULL PIPELINE. Follow recon → exploit → post_exploit → (codegen if needed) → report strictly."
	}
}

// AllowedDelegates returns the set of delegate_* tools permitted for a mode.
func AllowedDelegates(mode string) map[string]bool {
	mode = NormalizeAgentMode(mode)
	for _, a := range ListAgents() {
		if a.ID == mode {
			out := make(map[string]bool, len(a.Delegates))
			for _, d := range a.Delegates {
				out[d] = true
			}
			return out
		}
	}
	return map[string]bool{
		"delegate_recon": true, "delegate_exploit": true,
		"delegate_post_exploit": true, "delegate_codegen": true, "delegate_report": true,
	}
}
