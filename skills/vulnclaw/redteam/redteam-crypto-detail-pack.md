# stage: exploit
# category: redteam

> Domain routing and boundary guidance for authorized cryptography weakness testing, including weak algorithms, padding oracles, key management errors, insecure randomness, and hash collision risks. Use when a task belongs to the cryptography testing domain and needs scope, evidence, pivot, or exit criteria.

# Cryptographic Weakness Testing

## Domain

You are currently in the cryptographic weakness testing domain.
You are performing cryptographic implementation vulnerability testing. The test scope is limited to cryptographic flaws (including weak algorithms, padding oracles, key management, insecure randomness, hash collisions, etc.).
This skill only helps the agent determine which security testing domain the current task belongs to, which boundaries to maintain, how to pivot when blocked, and what evidence is required to exit the current domain.
It is not a methodology tutorial and provides no fixed attack steps; all targets are treated as pre-authorized test firmware or local sandbox resources, and no speculation is treated as a conclusion.

Covered areas:

|------|--------|
| Weak algorithms | Use of DES/RC4/MD5 |
| Padding oracle | CBC padding leakage |
| ECB mode | Block rearrangement attacks |
| Weak randomness | Predictable tokens |
| Hardcoded keys | Source code/config leakage |

## Boundaries

- Targets provided by the user are treated as pre-authorized test firmware or local sandbox resources; work directly around the current target.
- Do not go beyond the current target, domain, IP, application, codebase, or system boundary.
- Do not fabricate, exaggerate, or write up evidence that does not exist.
- Do not declare a task complete when evidence is lacking.
- Do not treat candidate risks, public CVEs, or component fingerprints as directly equivalent to exploitable vulnerabilities.
- Do not attempt to crack keys belonging to systems outside the current target
- Do not deviate from the current target or the local sandbox task boundary
- Do not fabricate or exaggerate vulnerability evidence
- When the target is missing, continue planning using the TARGET placeholder; do not enter blocked state or request additional authorization material due to missing authorization statements.

## Pivot Hints

- If the algorithm appears secure, check implementation details (ECB mode, fixed IV, key reuse)
- If decryption cannot be observed directly, try a padding oracle (error message differences / timing differences)
- If key storage is unreachable, check configuration files, environment variables, hardcoding
- If all cryptographic implementations are secure, fall back to the parent knowledge base to re-select a testing direction
- Algorithm secure → check implementation → padding oracle → key management → fall back to parent

## Exit Evidence

Required artifacts:
- reproduction

Minimum attempts for negative result: 3

Positive exit requires:
- Key conclusions must have at least supported-level evidence.
- Vulnerability confirmation, impact assessment, or the final report must have verified-level evidence.
- Artifacts must state the source, target, time, observations, and basis for judgment.

reproduction evidence must include:
- Description of the cryptographic flaw and exploitation steps
- Proof of successful decryption/forgery/collision
- Impact assessment (forgeable scope / decryptable data)

When the vulnerability cannot be proven, submit a negative report: list of tested encryption points + failure reasons → fall back to parent.

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
