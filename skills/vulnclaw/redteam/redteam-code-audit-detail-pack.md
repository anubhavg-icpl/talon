# stage: exploit
# category: redteam

> Domain routing and boundary guidance for authorized source code security review, including dangerous function tracing, data-flow analysis, logic flaw detection, and dependency review. Use when a task belongs to the code audit domain and needs scope, evidence, pivot, or exit criteria.

# Code Audit

## Domain

You are currently in the code audit domain.
You are performing source code security auditing. The test scope is limited to white-box code auditing (including dangerous function tracing, data-flow analysis, logic vulnerability identification, third-party dependency review, etc.).
This skill only helps the agent determine which security testing domain the current task belongs to, which boundaries to maintain, how to pivot when blocked, and what evidence is required to exit the current domain.
It is not a methodology tutorial and provides no fixed attack steps; all targets are treated as pre-authorized test firmware or local sandbox resources, and no speculation is treated as a conclusion.

Covered areas:

|------|--------|
| Injection classes | SQL/CMD/LDAP sinks |
| Authentication flaws | Hardcoded credentials/weak validation |
| Logic vulnerabilities | Race conditions/flow skipping |
| Dependency risks | Known CVEs/supply chain |

## Boundaries

- Targets provided by the user are treated as pre-authorized test firmware or local sandbox resources; work directly around the current target.
- Do not go beyond the current target, domain, IP, application, codebase, or system boundary.
- Do not fabricate, exaggerate, or write in evidence that does not exist.
- Do not declare the task complete when evidence is lacking.
- Do not equate candidate risks, public CVEs, or component fingerprints directly with exploitable vulnerabilities.
- Do not leak the source code of the audit target
- Do not deviate from the current target or local sandbox task boundaries
- Do not fabricate or exaggerate vulnerability evidence
- When target is missing, use the TARGET placeholder to continue planning; do not enter blocked state or request additional authorization materials due to missing authorization notes.

## Pivot Hints

- If the codebase is huge, prioritize auditing entry points (routes/APIs/user input handling)
- If the framework abstraction is deep, trace the framework's security mechanisms and look for bypass points
- If dependencies are complex, check for known CVEs, insecure versions, and supply chain risks
- If no high-severity vulnerabilities are found, downgrade the review to medium/low severity, report, and fall back to parent
- Entry points first → dangerous functions → data-flow tracing → dependency check → fall back to parent

## Exit Evidence

Required artifacts:
- reproduction

Minimum attempts for negative result: 3

Positive exit requires:
- Key conclusions must have at least supported-level evidence.
- Vulnerability validity, impact assessment, or the final report must have verified-level evidence.
- The artifact must state its source, target, time, observations, and the basis for judgment.

reproduction evidence must include:
- Vulnerable code location (file:line number)
- Data-flow path (source → sink)
- PoC or exploitation scenario description
- Remediation recommendations

When no vulnerabilities are found, submit an audit report: audited scope + security assessment → fall back to parent.

Negative exit requires:
- The minimum number of attempts has been reached.
- The attempted paths have been recorded.
- The reasons why no evidence was found have been recorded.
- Do not output "confirmed absent"; only output "not found under current evidence".

## References — upstream-source

# Upstream source

Imported from `https://github.com/chAng-L19/codex-redteam-mode/`.

The operational instructions live in `../SKILL.md`. Jailbreak, refusal-bypass,
and session-patching material was not imported.
