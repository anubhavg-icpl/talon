# stage: exploit
# category: redteam

> Domain routing and boundary guidance for authorized open redirect testing, including parameter redirects, meta or JavaScript redirects, and OAuth redirect_uri abuse. Use when a task belongs to the open redirect domain and needs scope, evidence, pivot, or exit criteria.

# Open Redirect Testing

## Domain

You are currently in the open redirect testing domain.
You are performing open redirect vulnerability testing. The test scope is limited to URL redirect bypasses (including parameter redirects, Meta/JS redirects, OAuth redirect_uri exploitation, etc.).
This skill exists only to help the agent determine which security testing domain the current task belongs to, which boundaries to maintain, how to pivot when blocked, and what evidence is required to exit the current domain.
It is not a methodology tutorial and provides no fixed attack steps; all targets are treated as pre-authorized test firmware or local sandbox resources, and no speculation is treated as a conclusion.

Coverage domains:

|------|--------|
| Parameter redirect | ?url=//evil.com |
| OAuth redirect | redirect_uri tampering |
| Meta refresh | HTML meta tag |
| JS redirect | controllable location.href |

## Boundaries

- Targets provided by the user are treated as pre-authorized test firmware or local sandbox resources; work directly around the current target.
- Do not go beyond the current target, domain, IP, application, codebase, or system boundary.
- Do not fabricate, exaggerate, or invent evidence that does not exist.
- Do not declare a task complete when evidence is lacking.
- Do not equate candidate risks, public CVEs, or component fingerprints directly with exploitable vulnerabilities.
- Do not use redirects to phish real users
- Do not deviate from the current target or local sandbox task boundary
- Do not fabricate or exaggerate vulnerability evidence
- If target is missing, continue planning using the TARGET placeholder; do not enter blocked status or request additional authorization material due to missing authorization statements.

## Pivot Hints

- If a URL whitelist is enforced, try @ symbols, backslashes, URL encoding, double encoding
- If restricted to relative paths, try protocol-relative URLs (//evil.com), path traversal
- If only same-domain is allowed, look for subdomain takeover or open redirect chains
- If all redirect points are secure, fall back to the upper-level knowledge base to reselect a testing direction
- Strict whitelist → URL encoding → path confusion → chained redirects → fall back to upper level

## Exit Evidence

Required artifacts:
- reproduction

Minimum attempts for negative result: 3

Positive exit requires:
- Key conclusions have at least supported-level evidence.
- Vulnerability confirmation, impact assessment, or the final report must have verified-level evidence.
- Artifacts must state the source, target, time, observations, and basis for judgment.

reproduction evidence must include:
- Complete request URL (including the redirect payload)
- Proof of successful redirect to an external domain
- Exploitation scenario assessment (OAuth token leakage / phishing aid)

When the vulnerability cannot be proven, submit a negative report: list of tested redirect points + failure reasons → fall back to upper level.

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
