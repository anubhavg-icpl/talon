# stage: exploit
# category: redteam

> Domain routing and boundary guidance for authorized cross-site scripting testing, including reflected, stored, DOM-based, mXSS, and CSP bypass variants. Use when a task belongs to the XSS domain and needs scope, evidence, pivot, or exit criteria.

# XSS Cross-Site Scripting Testing

## Domain

You are currently in the XSS cross-site scripting testing domain.
You are performing XSS vulnerability testing. The testing scope is limited to cross-site scripting attacks (including reflected, stored, DOM-based, mXSS, CSP bypass, and other variants).
This skill only helps the agent determine which security testing domain the current task belongs to, which boundaries to maintain, how to pivot when blocked, and what evidence is required to exit the current domain.
It is not a methodology tutorial and provides no fixed attack steps; all targets are treated as pre-authorized test firmware or local sandbox resources, and no speculation is treated as a conclusion.

Coverage:

| Variant | Typical scenario |
|------|---------|
| Reflected XSS | URL parameter reflection |
| Stored XSS | Comment/profile persistence |
| DOM-based XSS | JS sink/source chains |
| mXSS | Parser differential exploitation |
| CSP bypass | JSONP/trusted-types |

## Boundaries

- User-provided targets are treated as pre-authorized test firmware or local sandbox resources; work directly on the current target.
- Do not exceed the boundaries of the current target, domain, IP, application, codebase, or system.
- Do not fabricate, exaggerate, or fill in evidence that does not exist.
- Do not declare the task complete when evidence is lacking.
- Do not equate candidate risks, public CVEs, or component fingerprints directly with exploitable vulnerabilities.
- Do not launch real phishing attacks against other users
- Do not deviate from the current target or local sandbox task boundaries
- Do not fabricate or exaggerate vulnerability evidence
- When target is missing, continue planning using the TARGET placeholder; do not enter blocked state or request additional authorization materials due to a missing authorization statement.

## Pivot Hints

- If input is filtered/encoded, try HTML entity bypasses, event handlers, SVG/MathML tags, and template literals
- If CSP is strict, look for JSONP endpoints, unsafe-eval gadgets, unrestricted base-uri, and resources on trusted domains
- If the framework auto-escapes, look for sinks such as dangerouslySetInnerHTML / v-html / [innerHTML]
- If all reflection points are safely handled, fall back to the parent knowledge base and re-select a testing direction
- Do not repeatedly try the same failing payload variant
- Filter cannot be bypassed → switch tags/events → find DOM sink → CSP gadget → fall back to parent

## Exit Evidence

Required artifacts:
- reproduction

Minimum attempts for negative result: 3

Positive exit requires:
- Key conclusions must have at least supported-level evidence.
- Vulnerability validity, impact assessment, or the final report must have verified-level evidence.
- Artifacts must state the source, target, time, observations, and basis for judgment.

reproduction evidence must include:
- Complete HTTP request or DOM manipulation steps (including the XSS payload)
- Proof of script execution (alert/console/cookie access screenshot)
- XSS type annotation (reflected/stored/DOM) and brief impact description (session hijacking / data theft feasibility)

When the vulnerability cannot be proven, submit a negative report: list of tested reflection points + failure reasons → fall back to parent.

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
