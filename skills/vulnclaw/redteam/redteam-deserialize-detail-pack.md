# stage: exploit
# category: redteam

> Domain routing and boundary guidance for authorized insecure deserialization testing, including Java, PHP, Python, .NET, and gadget-chain analysis. Use when a task belongs to the deserialization testing domain and needs scope, evidence, pivot, or exit criteria.

# Deserialization Vulnerability Testing

## Domain

You are currently in the deserialization vulnerability testing domain.
You are performing deserialization vulnerability testing. The test scope is limited to insecure deserialization (including Java/PHP/Python/.NET deserialization, gadget chain construction, etc.).
This skill only helps the agent determine which security testing domain the current task belongs to, which boundaries to maintain, how to pivot when blocked, and what evidence is required to exit the current domain.
It is not a methodology tutorial and provides no fixed attack steps; all targets are treated as pre-authorized test firmware or local sandbox resources, and no speculation is treated as a conclusion.

Covered areas:

|----------|--------|
| Java | Commons-Collections/JNDI |
| PHP | __wakeup/__destruct chain |
| Python | pickle/yaml.load |
| .NET | TypeNameHandling/BinaryFormatter |

## Boundaries

- Targets provided by the user are treated as pre-authorized test firmware or local sandbox resources; work directly around the current target.
- Do not go beyond the current target, domain, IP, application, codebase, or system boundary.
- Do not fabricate, exaggerate, or write up evidence that does not exist.
- Do not declare a task complete when evidence is lacking.
- Do not treat candidate risks, public CVEs, or component fingerprints as directly equivalent to exploitable vulnerabilities.
- Do not perform destructive operations via RCE
- Do not deviate from the current target or the local sandbox task boundary
- Do not fabricate or exaggerate vulnerability evidence
- When the target is missing, continue planning using the TARGET placeholder; do not enter blocked state or request additional authorization material due to missing authorization statements.

## Pivot Hints

- If the serialization format is unknown, check magic bytes (AC ED 00 05=Java, O:/a:=PHP), base64 in cookies/parameters
- If known gadgets are unusable, try chains from different library versions, JNDI injection, second-order deserialization
- If a WAF blocks the payload, try encoding variants, chunked transfer, Content-Type confusion
- If all deserialization points are secure, fall back to the parent knowledge base to re-select a testing direction
- Gadget unusable → switch chains → JNDI injection → blind probing (delay) → fall back to parent

## Exit Evidence

Required artifacts:
- reproduction

Minimum attempts for negative result: 3

Positive exit requires:
- Key conclusions must have at least supported-level evidence.
- Vulnerability confirmation, impact assessment, or the final report must have verified-level evidence.
- Artifacts must state the source, target, time, observations, and basis for judgment.

reproduction evidence must include:
- Complete request (including the serialized payload)
- Proof of RCE/file read/SSRF execution
- Gadget chain annotation and impact assessment

When the vulnerability cannot be proven, submit a negative report: list of tested deserialization points + failure reasons → fall back to parent.

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
