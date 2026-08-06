# stage: exploit
# category: redteam

> Domain routing and boundary guidance for authorized SSRF testing, including basic SSRF, blind SSRF, protocol smuggling, and cloud metadata access paths. Use when a task belongs to the SSRF domain and needs scope, evidence, pivot, or exit criteria.

# SSRF Server-Side Request Forgery Testing

## Domain

You are currently in the SSRF server-side request forgery testing domain.
You are performing SSRF vulnerability testing. The test scope is limited to server-side request forgery (including basic SSRF, blind SSRF, protocol smuggling, cloud metadata access, etc.).
This skill exists only to help the agent determine which security testing domain the current task belongs to, which boundaries to maintain, how to pivot when blocked, and what evidence is required to exit the current domain.
It is not a methodology tutorial and provides no fixed attack steps; all targets are treated as pre-authorized test firmware or local sandbox resources, and no speculation is treated as a conclusion.

Coverage domains:

| Variant | Typical scenario |
|------|---------|
| Basic SSRF | direct internal network access |
| Blind SSRF | OOB callback confirmation |
| Protocol smuggling | gopher/dict exploitation |
| DNS rebinding | whitelist bypass |
| Cloud metadata | AWS/GCP/Azure IMDSv1 |

## Boundaries

- Targets provided by the user are treated as pre-authorized test firmware or local sandbox resources; work directly around the current target.
- Do not go beyond the current target, domain, IP, application, codebase, or system boundary.
- Do not fabricate, exaggerate, or invent evidence that does not exist.
- Do not declare a task complete when evidence is lacking.
- Do not equate candidate risks, public CVEs, or component fingerprints directly with exploitable vulnerabilities.
- Do not use SSRF to access external systems beyond the current target
- Do not deviate from the current target or local sandbox task boundary
- Do not fabricate or exaggerate vulnerability evidence
- If target is missing, continue planning using the TARGET placeholder; do not enter blocked status or request additional authorization material due to missing authorization statements.

## Pivot Hints

- If a URL whitelist restricts, try DNS rebinding, URL parsing discrepancies, redirect chains, IPv6 mapping
- If protocol-restricted, switch protocols (file://, gopher://, dict://), use redirects to switch protocols
- If the internal network is unreachable, try cloud metadata (169.254.169.254), local service enumeration
- If all request points are not exploitable, fall back to the upper-level knowledge base to reselect a testing direction
- Do not repeatedly retry the same failing bypass method
- Cannot bypass whitelist → DNS rebinding → redirect chains → switch protocol → fall back to upper level

## Exit Evidence

Required artifacts:
- reproduction

Minimum attempts for negative result: 3

Positive exit requires:
- Key conclusions have at least supported-level evidence.
- Vulnerability confirmation, impact assessment, or the final report must have verified-level evidence.
- Artifacts must state the source, target, time, observations, and basis for judgment.

reproduction evidence must include:
- Complete HTTP request (including the SSRF payload URL)
- Evidence that the server initiated a request (internal network response content/DNS callback/timing differences)
- Reachable scope assessment (internal network segments, cloud metadata, local files)

When the vulnerability cannot be proven, submit a negative report: list of tested URL parameters + failure reasons → fall back to upper level.

Negative exit requires:
- The minimum number of attempts is reached.
- Attempted paths are recorded.
- The reasons no evidence was found are recorded.
- Do not output "confirmed not present"; only output "not found under current evidence".

## References — upstream-source

# Upstream source

Imported from `https://github.com/chAng-L19/codex-redteam-mode/`.

The operational instructions live in `../SKILL.md`. Jailbreak, refusal-bypass,
and session-patching material was not imported.
