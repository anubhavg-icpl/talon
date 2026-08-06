# stage: exploit
# category: redteam

> Domain routing and boundary guidance for authorized business logic vulnerability testing, including race conditions, flow bypass, price tampering, permission logic errors, and bulk operation abuse. Use when a task belongs to the logic testing domain and needs scope, evidence, pivot, or exit criteria.

# Business Logic Vulnerability Testing

## Domain

You are currently in the business logic vulnerability testing domain.
You are performing business logic vulnerability testing. The test scope is limited to logic flaws (including race conditions, flow skipping, price tampering, permission logic errors, bulk operation abuse, etc.).
This skill only helps the agent determine which security testing domain the current task belongs to, which boundaries to maintain, how to pivot when blocked, and what evidence is required to exit the current domain.
It is not a methodology tutorial and provides no fixed attack steps; all targets are treated as pre-authorized test firmware or local sandbox resources, and no speculation is treated as a conclusion.

Covered areas:

|------|--------|
| Race conditions | Double submission/concurrent consumption |
| Flow skipping | Payment step bypass |
| Price tampering | Client-side amount modification |
| Permission logic | Incomplete role checks |
| IDOR | Unauthorized object reference access |

## Boundaries

- Targets provided by the user are treated as pre-authorized test firmware or local sandbox resources; work directly around the current target.
- Do not go beyond the current target, domain, IP, application, codebase, or system boundary.
- Do not fabricate, exaggerate, or write up evidence that does not exist.
- Do not declare a task complete when evidence is lacking.
- Do not treat candidate risks, public CVEs, or component fingerprints as directly equivalent to exploitable vulnerabilities.
- Do not cause real financial loss
- Do not deviate from the current target or the local sandbox task boundary
- Do not fabricate or exaggerate vulnerability evidence
- When the target is missing, continue planning using the TARGET placeholder; do not enter blocked state or request additional authorization material due to missing authorization statements.

## Pivot Hints

- If the flow has strict validation, analyze the state machine and look for skippable intermediate steps
- If the race window is extremely small, increase concurrent thread count, exploit HTTP/2 single-packet multi-request
- If amounts/quantities have server-side validation, try negative numbers, extremely large numbers, floating-point precision, currency unit confusion
- If all business logic is secure, fall back to the parent knowledge base to re-select a testing direction
- Strict validation → race attack → parameter tampering → flow analysis → fall back to parent

## Exit Evidence

Required artifacts:
- reproduction

Minimum attempts for negative result: 3

Positive exit requires:
- Key conclusions must have at least supported-level evidence.
- Vulnerability confirmation, impact assessment, or the final report must have verified-level evidence.
- Artifacts must state the source, target, time, observations, and basis for judgment.

reproduction evidence must include:
- Complete sequence of operation steps
- Proof of successful logic flaw exploitation (balance change/privilege escalation/flow bypass)
- Business impact assessment

When the vulnerability cannot be proven, submit a negative report: list of tested logic points + failure reasons → fall back to parent.

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
