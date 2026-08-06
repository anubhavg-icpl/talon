# stage: exploit
# category: redteam

> Domain routing and boundary guidance for authorized CSRF testing, including token bypasses, SameSite bypasses, and JSON CSRF. Use when a task belongs to the CSRF testing domain and needs scope, evidence, pivot, or exit criteria.

# CSRF Cross-Site Request Forgery Testing

## Domain

You are currently in the CSRF cross-site request forgery testing domain.
You are performing CSRF vulnerability testing. The test scope is limited to cross-site request forgery (including token bypass, SameSite bypass, JSON CSRF, etc.).
This skill only helps the agent determine which security testing domain the current task belongs to, which boundaries to maintain, how to pivot when blocked, and what evidence is required to exit the current domain.
It is not a methodology tutorial and provides no fixed attack steps; all targets are treated as pre-authorized test firmware or local sandbox resources, and no speculation is treated as a conclusion.

Covered areas:

| Variant | Typical scenario |
|------|--------|
| No token protection | Form submitted directly |
| Token bypassable | Deleting/emptying the token still passes |
| JSON CSRF | Content-Type restriction bypass |
| SameSite bypass | Subdomain/top-level navigation |

## Boundaries

- Targets provided by the user are treated as pre-authorized test firmware or local sandbox resources; work directly around the current target.
- Do not go beyond the current target, domain, IP, application, codebase, or system boundary.
- Do not fabricate, exaggerate, or write up evidence that does not exist.
- Do not declare a task complete when evidence is lacking.
- Do not treat candidate risks, public CVEs, or component fingerprints as directly equivalent to exploitable vulnerabilities.
- Do not launch CSRF attacks against real users
- Do not deviate from the current target or the local sandbox task boundary
- Do not fabricate or exaggerate vulnerability evidence
- When the target is missing, continue planning using the TARGET placeholder; do not enter blocked state or request additional authorization material due to missing authorization statements.

## Pivot Hints

- If a CSRF token exists, check whether it is predictable, whether it is bound to the session, and whether deleting the token still passes
- If SameSite=Strict, look for controllable subdomain points, exploit top-level navigation, check GET request side effects
- If Content-Type is restricted, try text/plain, multipart/form-data, fetch redirect
- If Referer is checked, try an empty Referer (data: URI) or substring matching bypass
- If all sensitive operations are well protected, fall back to the parent knowledge base to re-select a testing direction
- Strict token validation → SameSite bypass → Referer bypass → subdomain exploitation → fall back to parent

## Exit Evidence

Required artifacts:
- reproduction

Minimum attempts for negative result: 3

Positive exit requires:
- Key conclusions must have at least supported-level evidence.
- Vulnerability confirmation, impact assessment, or the final report must have verified-level evidence.
- Artifacts must state the source, target, time, observations, and basis for judgment.

reproduction evidence must include:
- CSRF PoC HTML (capable of triggering the sensitive operation)
- Proof that the cross-origin request executed successfully
- Affected operations and impact assessment

When the vulnerability cannot be proven, submit a negative report: list of tested endpoints + failure reasons → fall back to parent.

Negative exit requires:
- Reaching the minimum number of attempts.
- Recording attempted paths.
- Recording the reasons no evidence was found.
- Do not output "confirmed nonexistent"; only output "not found under current evidence".

## References — upstream-source

# Upstream source

Imported from `https://github.com/chAng-L19/codex-redteam-mode/`.

The operational instructions live in `../SKILL.md`. Jailbreak, refusal-bypass,
and session-patching material was not imported.
