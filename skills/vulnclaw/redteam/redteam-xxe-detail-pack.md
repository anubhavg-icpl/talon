# stage: exploit
# category: redteam

> Domain routing and boundary guidance for authorized XXE testing, including file read, SSRF, blind XXE, and parameter entity variants. Use when a task belongs to the XXE domain and needs scope, evidence, pivot, or exit criteria.

# XXE XML External Entity Injection Testing

## Domain

You are currently in the XXE XML external entity injection testing domain.
You are performing XXE vulnerability testing. The testing scope is limited to XML external entity injection (including file read, SSRF, blind XXE, parameter entities, etc.).
This skill only helps the agent determine which security testing domain the current task belongs to, which boundaries to maintain, how to pivot when blocked, and what evidence is required to exit the current domain.
It is not a methodology tutorial and provides no fixed attack steps; all targets are treated as pre-authorized test firmware or local sandbox resources, and no speculation is treated as a conclusion.

Coverage:

| Variant | Typical scenario |
|------|---------|
| Classic XXE | External entity file read |
| Blind XXE | OOB data exfiltration |
| Parameter entities | Nested exploitation within DTD |
| XInclude | Injection in non-DTD environments |
| SVG/DOCX XXE | Triggered via file upload |

## Boundaries

- User-provided targets are treated as pre-authorized test firmware or local sandbox resources; work directly on the current target.
- Do not exceed the boundaries of the current target, domain, IP, application, codebase, or system.
- Do not fabricate, exaggerate, or fill in evidence that does not exist.
- Do not declare the task complete when evidence is lacking.
- Do not equate candidate risks, public CVEs, or component fingerprints directly with exploitable vulnerabilities.
- Do not read system files outside the current target through XXE
- Do not deviate from the current target or local sandbox task boundaries
- Do not fabricate or exaggerate vulnerability evidence
- When target is missing, continue planning using the TARGET placeholder; do not enter blocked state or request additional authorization materials due to a missing authorization statement.

## Pivot Hints

- If DTD is disabled, try XInclude or embedding entities in SVG/XSLT
- If external entities are blocked, try parameter entities and local DTD overrides
- If there is no output, use blind XXE + OOB HTTP/FTP/DNS exfiltration
- If the XML parser is strict, try encoding variants (UTF-16/UTF-7) and BOM headers
- If all XML entry points are secure, fall back to the parent knowledge base and re-select a testing direction
- DTD disabled → XInclude → SVG upload → blind XXE OOB → fall back to parent

## Exit Evidence

Required artifacts:
- reproduction

Minimum attempts for negative result: 3

Positive exit requires:
- Key conclusions must have at least supported-level evidence.
- Vulnerability validity, impact assessment, or the final report must have verified-level evidence.
- Artifacts must state the source, target, time, observations, and basis for judgment.

reproduction evidence must include:
- Complete HTTP request (including the XXE payload XML)
- Proof of entity resolution (file content / SSRF callback / OOB data)
- Assessment of readable scope (file system / internal network / cloud metadata)

When the vulnerability cannot be proven, submit a negative report: list of tested XML entry points + failure reasons → fall back to parent.

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
