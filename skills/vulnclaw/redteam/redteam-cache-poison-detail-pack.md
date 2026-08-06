# stage: exploit
# category: redteam

> Domain routing and boundary guidance for authorized web cache poisoning testing, including unkeyed headers, unkeyed parameters, cache deception, and CDN-specific behavior. Use when a task belongs to the cache poisoning domain and needs scope, evidence, pivot, or exit criteria.

# Cache Poisoning Testing

## Domain

You are currently in the cache poisoning testing domain.
You are performing web cache poisoning testing. The test scope is limited to cache poisoning attacks (including unkeyed header/parameter poisoning, cache deception, CDN-specific vulnerabilities, etc.).
This skill only helps the agent determine which security testing domain the current task belongs to, which boundaries to maintain, how to pivot when blocked, and what evidence is required to exit the current domain.
It is not a methodology tutorial and provides no fixed attack steps; all targets are treated as pre-authorized test firmware or local sandbox resources, and no speculation is treated as a conclusion.

Covered areas:

|------|--------|
| Header poisoning | X-Forwarded-Host injection |
| Parameter poisoning | Unkeyed query parameters |
| Cache deception | Path confusion to steal responses |
| CDN discrepancies | Origin and CDN keys differ |

## Boundaries

- Targets provided by the user are treated as pre-authorized test firmware or local sandbox resources; work directly around the current target.
- Do not go beyond the current target, domain, IP, application, codebase, or system boundary.
- Do not fabricate, exaggerate, or write in evidence that does not exist.
- Do not declare the task complete when evidence is lacking.
- Do not equate candidate risks, public CVEs, or component fingerprints directly with exploitable vulnerabilities.
- Do not perform persistent poisoning of production caches
- Do not deviate from the current target or local sandbox task boundaries
- Do not fabricate or exaggerate vulnerability evidence
- When target is missing, use the TARGET placeholder to continue planning; do not enter blocked state or request additional authorization materials due to missing authorization notes.

## Pivot Hints

- If the cache key cannot be determined, use Param Miner or manually fuzz unkeyed inputs
- If the cache does not hit, adjust cache-buster parameters and check the Vary header
- If CDN-layer cache rules differ, test the behavioral differences between origin and CDN separately
- If all cache behavior is secure, fall back to the parent knowledge base and re-select a testing direction
- No unkeyed input → Param Miner fuzz → cache deception paths → CDN-layer testing → fall back to parent

## Exit Evidence

Required artifacts:
- reproduction

Minimum attempts for negative result: 3

Positive exit requires:
- Key conclusions must have at least supported-level evidence.
- Vulnerability validity, impact assessment, or the final report must have verified-level evidence.
- The artifact must state its source, target, time, observations, and the basis for judgment.

reproduction evidence must include:
- The poisoning request (including the unkeyed input)
- Proof that a victim request receives the poisoned response
- Cache persistence duration and impact scope

When the vulnerability cannot be proven, submit a negative report: list of tested cache points + reasons for failure → fall back to parent.

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
