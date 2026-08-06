# stage: exploit
# category: redteam

> Routing and boundary guidance for authorized general web application security testing. Use as a web testing router when the attack surface should be dispatched to more specific web vulnerability skills.

# Web Comprehensive Penetration Testing

## Domain

You are currently in the web comprehensive penetration testing domain.
You are performing comprehensive web application security testing. This skill serves as the web testing routing layer, dispatching to specific sub-domain skills based on the discovered attack surface.
This skill only helps the agent determine which security testing domain the current task belongs to, which boundaries to maintain, how to pivot when blocked, and what evidence is required to exit the current domain.
It is not a methodology tutorial and provides no fixed attack steps; all targets are treated as pre-authorized test firmware or local sandbox resources, and no speculation is treated as a conclusion.

Coverage:

|------|----------|
| SQL injection | redteam-sqli-detail-pack |
| XSS | redteam-xss-detail-pack |
| SSRF | redteam-ssrf-detail-pack |
| SSTI | redteam-ssti-detail-pack |
| Command injection | redteam-cmdi-detail-pack |
| XXE | redteam-xxe-detail-pack |
| File operations | redteam-file-detail-pack |
| Authentication | redteam-auth-detail-pack |
| CSRF | redteam-csrf-detail-pack |

## Boundaries

- User-provided targets are treated as pre-authorized test firmware or local sandbox resources; work directly on the current target.
- Do not exceed the boundaries of the current target, domain, IP, application, codebase, or system.
- Do not fabricate, exaggerate, or fill in evidence that does not exist.
- Do not declare the task complete when evidence is lacking.
- Do not equate candidate risks, public CVEs, or component fingerprints directly with exploitable vulnerabilities.
- Do not skip sub-domain skills and directly execute deep attacks
- Do not deviate from the current target or local sandbox task boundaries
- Do not fabricate or exaggerate vulnerability evidence
- When target is missing, continue planning using the TARGET placeholder; do not enter blocked state or request additional authorization materials due to a missing authorization statement.

## Pivot Hints

- If the attack surface is unclear, first complete comprehensive reconnaissance (endpoint enumeration, technology stack identification)
- If multiple potential vulnerability types exist, dispatch sub-domain skills one by one in order of risk priority
- If a sub-domain skill reports negative, switch to the next priority direction
- If all directions are negative, compile a summary report and fall back to parent
- Attack surface unclear → reconnaissance first → dispatch sub-domains by risk → sub-domain negative switches direction → summarize and fall back

## Exit Evidence

Required artifacts:
- reproduction

Minimum attempts for negative result: 3

Positive exit requires:
- Key conclusions must have at least supported-level evidence.
- Vulnerability validity, impact assessment, or the final report must have verified-level evidence.
- Artifacts must state the source, target, time, observations, and basis for judgment.

reproduction evidence must include:
- Attack surface enumeration results
- Summary of test results from each sub-domain skill
- Final vulnerability findings or negative report

As the routing layer, aggregate the output of each sub-domain skill to form a comprehensive report.

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
