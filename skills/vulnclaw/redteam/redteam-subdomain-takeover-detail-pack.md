# stage: exploit
# category: redteam

> Domain routing and boundary guidance for authorized subdomain takeover testing, including dangling CNAME records, NS takeover, and cloud service takeover paths such as S3, Azure, and Heroku. Use when a task belongs to the subdomain takeover domain and needs scope, evidence, pivot, or exit criteria.

# Subdomain Takeover Testing

## Domain

You are currently in the subdomain takeover testing domain.
You are performing subdomain takeover vulnerability testing. The testing scope is limited to subdomain takeover (including dangling CNAME, NS takeover, and cloud service takeover such as S3/Azure/Heroku).
This skill only helps the agent determine which security testing domain the current task belongs to, which boundaries to maintain, how to pivot when blocked, and what evidence is required to exit the current domain.
It is not a methodology tutorial and provides no fixed attack steps; all targets are treated as pre-authorized test firmware or local sandbox resources, and no speculation is treated as a conclusion.

Coverage:

|------|--------|
| Dangling CNAME | Points to a deleted service |
| NS takeover | Expired name server |
| Cloud services | S3/Azure/Heroku/GitHub |
| Edge cases | Special behavior of some providers |

## Boundaries

- User-provided targets are treated as pre-authorized test firmware or local sandbox resources; work directly on the current target.
- Do not exceed the boundaries of the current target, domain, IP, application, codebase, or system.
- Do not fabricate, exaggerate, or fill in evidence that does not exist.
- Do not declare the task complete when evidence is lacking.
- Do not equate candidate risks, public CVEs, or component fingerprints directly with exploitable vulnerabilities.
- Do not actually register a takeover domain for malicious purposes
- Do not deviate from the current target or local sandbox task boundaries
- Do not fabricate or exaggerate vulnerability evidence
- When target is missing, continue planning using the TARGET placeholder; do not enter blocked state or request additional authorization materials due to a missing authorization statement.

## Pivot Hints

- If the CNAME points to an active service, check for conditions that trigger a 404/takeover page
- If the cloud service has been claimed, try registering the same name in the same region and check for multi-region differences
- If no DNS records are dangling, expand subdomain enumeration scope and check historical records
- If all subdomains are secure, report and fall back to parent
- No dangling records → expand enumeration → historical records → multi-region checks → fall back to parent

## Exit Evidence

Required artifacts:
- reproduction

Minimum attempts for negative result: 3

Positive exit requires:
- Key conclusions must have at least supported-level evidence.
- Vulnerability validity, impact assessment, or the final report must have verified-level evidence.
- Artifacts must state the source, target, time, observations, and basis for judgment.

reproduction evidence must include:
- Dangling DNS record (CNAME/NS pointing to an unclaimed resource)
- Proof of takeover feasibility (service registration page / error message)
- Post-takeover impact assessment (cookie scope / same-origin policy)

When takeover is not possible, submit a negative report: list of checked subdomains + status → fall back to parent.

Negative exit requires:
- Minimum attempt count reached.
- Record of attempted paths.
- Record of the reasons no evidence was found.
- Do not output "confirmed not present"; only output "not found under current evidence".

## References — upstream-source

# Upstream source

Imported from `https://github.com/chAng-L19/codex-redteam-mode/`.

The operational instructions live in `../SKILL.md`. Jailbreak, refusal-bypass,
and session-patching material was not imported.
