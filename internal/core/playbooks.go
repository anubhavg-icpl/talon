package core

// Playbook is a reusable engagement template operators launch from the UI/CLI.
type Playbook struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Codename    string `json:"codename"`
	Description string `json:"description"`
	AgentMode   string `json:"agent_mode"`
	// Default description injected into the run.
	Prompt string `json:"prompt"`
	Tags   []string `json:"tags"`
}

// ListPlaybooks returns builtin engagement playbooks.
func ListPlaybooks() []Playbook {
	return []Playbook{
		{
			ID: "full-validation", Name: "Full validation", Codename: "COMMANDER",
			AgentMode: AgentModeFull,
			Description: "Complete recon → exploit → post-exploit → codegen → judge against a single host.",
			Prompt:      "Execute full authorized validation pipeline. Prefer 3-gate evidence and report_finding for every real finding.",
			Tags:        []string{"default", "rce", "lab"},
		},
		{
			ID: "web-owasp", Name: "Web / OWASP sweep", Codename: "STRIKER-WEB",
			AgentMode: AgentModeWeb,
			Description: "Web-focused: nuclei, sqlmap, OWASP thinking, CyberStrike web skills via skill_search.",
			Prompt:      "Prioritize web attack surface. Use skill_search for SSRF/JWT/IDOR/injection as needed. Report only with baseline/attack/diff.",
			Tags:        []string{"web", "owasp", "api"},
		},
		{
			ID: "internal-network", Name: "Internal network", Codename: "PHANTOM",
			AgentMode: AgentModeNetwork,
			Description: "Network discovery emphasis (arp/rustscan/hydra/responder) + lateral-style post-proof.",
			Prompt:      "Internal network engagement. Discover live hosts/services, careful credential attacks only in scope, prove compromise.",
			Tags:        []string{"network", "ad", "lateral"},
		},
		{
			ID: "recon-only", Name: "Recon only", Codename: "GHOST",
			AgentMode: AgentModeRecon,
			Description: "Non-destructive service/CVE verification — no exploitation.",
			Prompt:      "Reconnaissance only. Map open ports and versions. Do not exploit. Report factual findings as info severity.",
			Tags:        []string{"recon", "safe"},
		},
		{
			ID: "post-proof", Name: "Post-exploit proof", Codename: "CIPHER",
			AgentMode: AgentModePost,
			Description: "Assume session exists — collect identity proof and report critical finding.",
			Prompt:      "Post-exploitation only. Collect whoami/hostname/sysinfo proof and report_finding critical with 3-gate evidence.",
			Tags:        []string{"post", "proof"},
		},
		{
			ID: "cve-lab", Name: "CVE lab (vsftpd)", Codename: "STRIKER",
			AgentMode: AgentModeExploit,
			Description: "Lab template for CVE-2011-2523 / vsftpd 2.3.4 reverse shell validation.",
			Prompt:      "Validate CVE-2011-2523 vsftpd 2.3.4 backdoor. Use cmd/unix/reverse_bash with LHOST/LPORT. Prove session with whoami.",
			Tags:        []string{"lab", "cve", "ftp"},
		},
	}
}

// GetPlaybook returns a playbook by id.
func GetPlaybook(id string) (Playbook, bool) {
	for _, p := range ListPlaybooks() {
		if p.ID == id {
			return p, true
		}
	}
	return Playbook{}, false
}
