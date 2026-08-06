# stage: exploit
# category: redteam

> Domain routing and boundary guidance for authorized general injection testing outside SQL injection, including NoSQL, LDAP, XPath, and expression language injection. Use when a task belongs to the general injection domain and needs scope, evidence, pivot, or exit criteria.

# General Injection Testing

## Domain

You are currently in the general injection testing domain.
You are performing general injection vulnerability testing. The test scope covers non-SQL injection categories (including NoSQL injection, LDAP injection, XPath injection, expression language injection, etc.).
This skill only helps the agent determine which security testing domain the current task belongs to, which boundaries to maintain, how to pivot when blocked, and what evidence is required to exit the current domain.
It is not a methodology tutorial and provides no fixed attack steps; all targets are treated as pre-authorized test firmware or local sandbox resources, and no speculation is treated as a conclusion.

Covered areas:

|------|--------|
| NoSQL injection | MongoDB $ne/$regex |
| LDAP injection | )(|(uid=*) |
| XPath injection | ' or '1'='1 |
| EL injection | ${applicationScope} |
| Header injection | CRLF/Host injection |

## Boundaries

- Targets provided by the user are treated as pre-authorized test firmware or local sandbox resources; work directly around the current target.
- Do not go beyond the current target, domain, IP, application, codebase, or system boundary.
- Do not fabricate, exaggerate, or write up evidence that does not exist.
- Do not declare a task complete when evidence is lacking.
- Do not treat candidate risks, public CVEs, or component fingerprints as directly equivalent to exploitable vulnerabilities.
- Do not perform destructive injection operations
- Do not deviate from the current target or the local sandbox task boundary
- Do not fabricate or exaggerate vulnerability evidence
- When the target is missing, continue planning using the TARGET placeholder; do not enter blocked state or request additional authorization material due to missing authorization statements.

## Pivot Hints

- If the backend type is unknown, test in parallel with multiple injection probes ($ne/$gt, *)(|, ' or '1'='1)
- If input is strictly filtered, try encoding bypass, Unicode normalization differences, HPP parameter pollution
- If there is no output, confirm via error-based/time-based blind injection
- If all injection points are secure, fall back to the parent knowledge base to re-select a testing direction
- Backend unknown → multiple probes → encoding bypass → blind injection confirmation → fall back to parent

## Exit Evidence

Required artifacts:
- reproduction

Minimum attempts for negative result: 3

Positive exit requires:
- Key conclusions must have at least supported-level evidence.
- Vulnerability confirmation, impact assessment, or the final report must have verified-level evidence.
- Artifacts must state the source, target, time, observations, and basis for judgment.

reproduction evidence must include:
- Complete request (including the injection payload)
- Proof of successful injection execution
- Injection type annotation and impact assessment

When the vulnerability cannot be proven, submit a negative report: list of tested injection points + failure reasons → fall back to parent.

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
