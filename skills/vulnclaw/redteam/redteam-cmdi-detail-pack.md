# stage: exploit
# category: redteam

> Domain routing and boundary guidance for authorized operating system command injection testing, including direct injection, blind injection, out-of-band callbacks, and argument injection. Use when a task belongs to the command injection domain and needs scope, evidence, pivot, or exit criteria.

# OS Command Injection Testing

## Domain

You are currently in the OS command injection testing domain.
You are performing operating system command injection testing. The test scope is limited to command injection (including direct injection, blind injection, OOB out-of-band, argument injection, etc.).
This skill only helps the agent determine which security testing domain the current task belongs to, which boundaries to maintain, how to pivot when blocked, and what evidence is required to exit the current domain.
It is not a methodology tutorial and provides no fixed attack steps; all targets are treated as pre-authorized test firmware or local sandbox resources, and no speculation is treated as a conclusion.

Covered areas:

| Variant | Typical scenario |
|------|---------|
| Direct injection | Controllable command concatenation |
| Blind injection | Time-delay/OOB confirmation |
| Argument injection | --flag injection |
| Environment variable injection | env override |

## Boundaries

- Targets provided by the user are treated as pre-authorized test firmware or local sandbox resources; work directly around the current target.
- Do not go beyond the current target, domain, IP, application, codebase, or system boundary.
- Do not fabricate, exaggerate, or write in evidence that does not exist.
- Do not declare the task complete when evidence is lacking.
- Do not equate candidate risks, public CVEs, or component fingerprints directly with exploitable vulnerabilities.
- Do not execute destructive system commands such as rm -rf / format
- Do not deviate from the current target or local sandbox task boundaries
- Do not fabricate or exaggerate vulnerability evidence
- When target is missing, use the TARGET placeholder to continue planning; do not enter blocked state or request additional authorization materials due to missing authorization notes.

## Pivot Hints

- If command separators are filtered, try newlines, $() substitution, backticks, or %0a encoding
- If spaces are blocked, use $IFS, {cmd,arg}, or tab as substitutes
- If commands are blacklisted, use wildcards (c?t /etc/p?sswd), variable concatenation, or base64-encoded execution
- If there is no output echo, use DNS out-of-band, time-delay inference, or write files to the web directory
- If all entry points are not injectable, fall back to the parent knowledge base and re-select a testing direction
- Separators blocked → switch encoding → no echo, go OOB → argument injection → fall back to parent

## Exit Evidence

Required artifacts:
- reproduction

Minimum attempts for negative result: 3

Positive exit requires:
- Key conclusions must have at least supported-level evidence.
- Vulnerability validity, impact assessment, or the final report must have verified-level evidence.
- The artifact must state its source, target, time, observations, and the basis for judgment.

reproduction evidence must include:
- Complete HTTP request (including the command injection payload)
- Proof of command execution (echoed output / DNS callback / time-delay difference / file creation)
- Injection type annotation and privilege level assessment

When the vulnerability cannot be proven, submit a negative report: list of tested parameters + reasons for failure → fall back to parent.

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
