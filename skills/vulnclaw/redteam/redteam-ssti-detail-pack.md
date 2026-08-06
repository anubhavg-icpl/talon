# stage: exploit
# category: redteam

> Domain routing and boundary guidance for authorized server-side template injection testing, including Jinja2, Twig, Freemarker, Velocity, and Thymeleaf engines. Use when a task belongs to the SSTI domain and needs scope, evidence, pivot, or exit criteria.

# SSTI Server-Side Template Injection Testing

## Domain

You are currently in the SSTI server-side template injection testing domain.
You are performing server-side template injection testing. The testing scope is limited to SSTI (including template engines such as Jinja2, Twig, Freemarker, Velocity, and Thymeleaf).
This skill only helps the agent determine which security testing domain the current task belongs to, which boundaries to maintain, how to pivot when blocked, and what evidence is required to exit the current domain.
It is not a methodology tutorial and provides no fixed attack steps; all targets are treated as pre-authorized test firmware or local sandbox resources, and no speculation is treated as a conclusion.

Coverage:

|------|------------|
| Jinja2 | {{config.__class__.__init__.__globals__}} |
| Twig | {{_self.env.registerUndefinedFilterCallback}} |
| Freemarker | <#assign ex="freemarker.template.utility.Execute"?new()> |
| Velocity | #set($x='')#set($rt=$x.class.forName('java.lang.Runtime')) |
| Thymeleaf | __${T(java.lang.Runtime).getRuntime().exec('id')}__ |

## Boundaries

- User-provided targets are treated as pre-authorized test firmware or local sandbox resources; work directly on the current target.
- Do not exceed the boundaries of the current target, domain, IP, application, codebase, or system.
- Do not fabricate, exaggerate, or fill in evidence that does not exist.
- Do not declare the task complete when evidence is lacking.
- Do not equate candidate risks, public CVEs, or component fingerprints directly with exploitable vulnerabilities.
- Do not execute destructive system commands through RCE
- Do not deviate from the current target or local sandbox task boundaries
- Do not fabricate or exaggerate vulnerability evidence
- When target is missing, continue planning using the TARGET placeholder; do not enter blocked state or request additional authorization materials due to a missing authorization statement.

## Pivot Hints

- If the template engine is unknown, use detection payloads ({{7*7}}, ${7*7}, <%= 7*7 %>) to identify the engine
- If sandbox restrictions apply, try MRO chain traversal, built-in object escape, and known gadget chains
- If input is filtered, try encoding bypasses, string concatenation, and alternative attribute-access syntax
- If all template points are secure, fall back to the parent knowledge base and re-select a testing direction
- Do not repeatedly try the same failing payload
- Probe returns no output → blind time-based injection → switch detection syntax → try other parameters → fall back to parent

## Exit Evidence

Required artifacts:
- reproduction

Minimum attempts for negative result: 3

Positive exit requires:
- Key conclusions must have at least supported-level evidence.
- Vulnerability validity, impact assessment, or the final report must have verified-level evidence.
- Artifacts must state the source, target, time, observations, and basis for judgment.

reproduction evidence must include:
- Complete HTTP request (including the template injection payload)
- Proof of template execution (computation result / command output / file read)
- Engine type annotation and RCE feasibility assessment

When the vulnerability cannot be proven, submit a negative report: list of tested template points + failure reasons → fall back to parent.

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
