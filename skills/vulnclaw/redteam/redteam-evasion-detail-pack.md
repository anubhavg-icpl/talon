# stage: exploit
# category: redteam

> Domain routing and boundary guidance for authorized defense evasion and bypass testing, including WAF bypass, AV/EDR evasion, logging considerations, and traffic obfuscation. Use when a task belongs to the evasion domain and needs scope, evidence, pivot, or exit criteria.

# Defense Evasion and Bypass

## Domain

You are currently in the defense evasion and bypass domain.
You are performing defense evasion testing. The test scope is limited to security protection bypass (including WAF bypass, AV/EDR evasion, log clearing, traffic obfuscation, etc.).
This skill only helps the agent determine which security testing domain the current task belongs to, which boundaries to maintain, how to pivot when blocked, and what evidence is required to exit the current domain.
It is not a methodology tutorial and provides no fixed attack steps; all targets are treated as pre-authorized test firmware or local sandbox resources, and no speculation is treated as a conclusion.

Covered areas:

|------|--------|
| WAF bypass | Encoding/chunking/HPP |
| AV evasion | Loaders/in-memory execution |
| EDR bypass | Unhooking/direct syscalls |
| Log evasion | Time windows/legitimate camouflage |

## Boundaries

- Targets provided by the user are treated as pre-authorized test firmware or local sandbox resources; work directly around the current target.
- Do not go beyond the current target, domain, IP, application, codebase, or system boundary.
- Do not fabricate, exaggerate, or write up evidence that does not exist.
- Do not declare a task complete when evidence is lacking.
- Do not treat candidate risks, public CVEs, or component fingerprints as directly equivalent to exploitable vulnerabilities.
- Do not permanently disable the target's security protection facilities
- Do not deviate from the current target or the local sandbox task boundary
- Do not fabricate or exaggerate vulnerability evidence
- When the target is missing, continue planning using the TARGET placeholder; do not enter blocked state or request additional authorization material due to missing authorization statements.

## Pivot Hints

- If WAF rules are strict, try encoding variants, chunked transfer, HTTP parameter pollution, protocol-layer bypass
- If AV/EDR detects and kills, try loader obfuscation, in-memory execution, legitimate-tool LOLBins
- If log monitoring is thorough, try low-frequency operations, legitimate-traffic camouflage, time-window exploitation
- If all evasion methods are detected, record the detection mechanism characteristics and fall back to parent
- WAF → encoding variants → AV → in-memory execution → EDR → syscall → fall back to parent

## Exit Evidence

Required artifacts:
- reproduction

Minimum attempts for negative result: 3

Positive exit requires:
- Key conclusions must have at least supported-level evidence.
- Vulnerability confirmation, impact assessment, or the final report must have verified-level evidence.
- Artifacts must state the source, target, time, observations, and basis for judgment.

reproduction evidence must include:
- Description of the evasion technique and implementation steps
- Proof of successful bypass (payload executed / no alert triggered)
- Type of protection bypassed and assessment of residual detection capability

When bypass is not possible, submit a negative report: attempted evasion methods + detection mechanism characteristics → fall back to parent.

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
