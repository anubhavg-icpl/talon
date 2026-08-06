# stage: recon
# category: redteam

> Recon intake skill for first contact with a bare domain, URL, or IP address. Use to build an initial recon_profile and provide factual inputs for CVE lookup and attack-path routing.

# Recon Intake

## Domain

The reconnaissance entry point when a bare domain/URL/IP first enters security assessment.
Responsible for building a target asset profile (recon_profile) from scratch, providing factual grounds for subsequent CVE Lookup and attack path assignment.

Phase order:
1. DNS resolution + liveness detection
2. Port scanning + service fingerprinting
3. Subdomain enumeration
4. Web directory/sensitive file discovery
5. WAF/CDN identification
6. Technology stack fingerprinting (CMS/framework/middleware versions)

## Boundaries

- Do not perform any active vulnerability exploitation
- Do not send destructive requests (DELETE/DROP/shutdown)
- Probe only the current target; do not automatically expand to related domains/IPs not present in the task
- Do not perform brute forcing or password spraying
- Do not bypass rate limits (if rate-limited, slow down or pause)
- Reconnaissance depth stops at information gathering; do not enter the vulnerability verification phase

## Pivot Hints

- WAF/CDN blocks direct connection → try historical DNS, email headers, certificate search to obtain the real IP
- Subdomain enumeration limited → certificate transparency logs, DNS zone transfer, reverse lookup of related domains
- Directory scanning blocked → slow down, change User-Agent, use custom wordlists
- All ports filtered → check IPv6, try common high ports, confirm target is alive
- Information gathering saturated → organize the collected information, output recon_profile and advance to the next phase

## Exit Evidence

- Required: recon_profile, port_scan_result, service_fingerprint
- min_attempts: 4

## References — upstream-source

# Upstream source

Imported from `https://github.com/chAng-L19/codex-redteam-mode/`.

The operational instructions live in `../SKILL.md`. Jailbreak, refusal-bypass,
and session-patching material was not imported.
