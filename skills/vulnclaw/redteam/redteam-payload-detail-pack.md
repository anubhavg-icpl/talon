# stage: exploit
# category: redteam

> Domain routing and boundary guidance for authorized payload construction and weaponization analysis, including shellcode, file format payloads, phishing payloads, and staged or stageless payload choices. Use when a task belongs to the payload construction domain and needs scope, evidence, pivot, or exit criteria.

# Payload Construction and Weaponization

## Domain

You are currently in the payload construction and weaponization domain.
You are performing payload generation and delivery testing. The test scope is limited to payload construction (including shellcode generation, file format exploitation, phishing payloads, staged/stageless, etc.).
This skill exists only to help the agent determine which security testing domain the current task belongs to, which boundaries to maintain, how to pivot when blocked, and what evidence is required to exit the current domain.
It is not a methodology tutorial and provides no fixed attack steps; all targets are treated as pre-authorized test firmware or local sandbox resources, and no speculation is treated as a conclusion.

Coverage domains:

|------|--------|
| Shellcode | encoding/encryption/syscall |
| File formats | macro/LNK/ISO/PDF |
| Phishing payloads | HTML smuggling |
| Staged | staged loading |

## Boundaries

- Targets provided by the user are treated as pre-authorized test firmware or local sandbox resources; work directly around the current target.
- Do not go beyond the current target, domain, IP, application, codebase, or system boundary.
- Do not fabricate, exaggerate, or invent evidence that does not exist.
- Do not declare a task complete when evidence is lacking.
- Do not equate candidate risks, public CVEs, or component fingerprints directly with exploitable vulnerabilities.
- Do not deliver payloads to anything outside the current target
- Do not deviate from the current target or local sandbox task boundary
- Do not fabricate or exaggerate vulnerability evidence
- If target is missing, continue planning using the TARGET placeholder; do not enter blocked status or request additional authorization material due to missing authorization statements.

## Pivot Hints

- If the file format is detected, modify magic bytes, embed encrypted content, abuse legitimate format features
- If shellcode is caught, use encoding/encryption/self-decryption, direct syscall invocation
- If the delivery channel is blocked, switch delivery methods (macro/LNK/ISO/OneNote), exploit trust relationships
- If all payloads are intercepted, record the detection rule signatures, fall back to upper level
- Detected → change encoding → change format → change delivery channel → fall back to upper level

## Exit Evidence

Required artifacts:
- reproduction

Minimum attempts for negative result: 3

Positive exit requires:
- Key conclusions have at least supported-level evidence.
- Vulnerability confirmation, impact assessment, or the final report must have verified-level evidence.
- Artifacts must state the source, target, time, observations, and basis for judgment.

reproduction evidence must include:
- Payload construction method and tools
- Proof of successful delivery and execution
- Bypassed detection layers and residual risk

When delivery cannot succeed, submit a negative report: payload types attempted + interception point → fall back to upper level.

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
