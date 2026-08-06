# stage: exploit
# category: specialized

> HackerOne bounty program scope-guard workflow — read program scope, enforce scope and program rules, then delegate each in-scope asset to pentest-flow

# HackerOne Bounty Scope-Guard Skill

You are currently executing the HackerOne bug bounty workflow. This skill is a **scope-guard wrapper**:
it first parses and enforces program scope and program rules, then delegates each **in-scope asset** to
the `pentest-flow` skill for the actual penetration testing. **Never touch out-of-scope assets at any stage.**

The launch parameter is a HackerOne program link (`<SCOPE LINK>`), for example
`hackerone.com/<handle>` or `.../policy_scopes`. This skill has no preset target
(frontmatter `requires_target: false`) — targets are discovered from the scope.

## Phase 1: Read scope

1. **Prefer fetching the link**
   - Use the fetch tool to access the provided `<SCOPE LINK>`.
   - HackerOne is a JavaScript SPA: fetching `hackerone.com/*` typically returns only
     an **empty app shell** with no rendered scope rows (scope data comes from `hackerone.com/graphql`
     or authenticated `api.hackerone.com`, which naive fetch cannot reach).
   - **Detect empty shell**: if the response has no recognizable scope table / asset list
     (only a skeleton like `<div id="app">`), treat it as a fetch failure.

2. **Fallback: ask user to paste scope**
   - When fetch returns empty / is login-walled / is JS-rendered only, **this is the norm, not the exception**.
   - Ask the user to paste the in-scope and out-of-scope tables directly from the program page's
     **Scope** tab. Provide a paste format example for reference:

     ```
     In scope:
     https://api.example.com        | URL       | Eligible for bounty
     *.example.com                  | WILDCARD  | Eligible for bounty
     app.example.com                | URL       | In scope, NOT bounty-eligible
     com.example.android            | GOOGLE_PLAY_APP_ID | Eligible for bounty

     Out of scope:
     blog.example.com               | URL
     *.corp.example.com             | WILDCARD
     ```

3. **Lenient parse**
   - Extract two lists from the pasted or fetched result: **in-scope** and **out-of-scope**.
   - Identify asset type (human label or API enum are both fine):
     `URL`, `WILDCARD` (`*.x.com`), `CIDR`/IP, `SOURCE_CODE`,
     `GOOGLE_PLAY_APP_ID`/`APPLE_STORE_APP_ID`/`TESTFLIGHT`/`OTHER_APK`/`OTHER_IPA`,
     `HARDWARE`, `AI_MODEL`, `SMART_CONTRACT`, `OTHER`, etc.
   - Identify **eligibility tri-state** (submission and bounty are two independent booleans):
     - `submission=true, bounty=true` → in scope, testable, bounty-eligible.
     - `submission=true, bounty=false` → **in scope, testable, no bounty** (do not confuse with out-of-scope).
     - `submission=false` → **out of scope, never test**.
   - When parsing is uncertain, confirm with the user. **Never default an asset to in-scope.**

4. **Output**
   - In-scope asset list (with type + eligibility).
   - Out-of-scope **deny-list** (enforced throughout).

## Phase 2: Enforce boundaries

Before entering any testing, explicitly write down and always obey these **hard rules**:

1. **Scope hard boundaries**
   - Only test assets in the in-scope list.
   - Assets in the out-of-scope deny-list must **never be touched** — no fetch, no scan, no payload.
   - `pentest-flow` can only directly act on `URL` and `WILDCARD` types; other types
     (mobile app / source / CIDR / hardware, etc.) are not automated yet — confirm the handling approach with the user first.

2. **Program rules (layered on top of Talon's existing `BLOCKED_PATTERNS` / `RESERVED_IP_RANGES`)**
   - **No DoS / no availability impact**: no stress testing, resource exhaustion, or mass concurrency (the primary red line for automated scanners).
   - **Respect rate limits / automation limits**: low speed, serial; obey program-declared "no automated scanning" clauses.
   - **No social engineering**: do not target personnel, no phishing.
   - **Minimal impact / no PII exfiltration**: verify the vulnerability and stop; do not export real user data or perform destructive operations.

3. **Exception handling**
   - If any step may cross scope boundaries or trigger the above rules, **stop and ask the user.**

## Phase 3: Enumerate & confirm

1. List all parsed in-scope assets to the user (numbered, with type and eligibility).
2. Ask which asset to start with (or all).
3. **Process one asset at a time**, confirming each, to avoid concurrency-induced boundary violations or rate limit triggers.

## Phase 4: Delegate to pentest-flow

For a **single in-scope asset** selected by the user:

1. Pass that asset as the target to the `pentest-flow` skill, executing the full
   recon → vuln-discovery → exploitation workflow.
2. Stay within scope throughout: new subdomains / endpoints discovered by `pentest-flow` that fall outside
   the in-scope definition (especially those not matching any in-scope `WILDCARD`) must be **excluded** and flagged to the user.
3. Continuously check against the Phase 2 program rules.

## Phase 5: Report — HackerOne submission format

For each confirmed finding, produce a report in **HackerOne submission format**:

1. **Title** — concise vulnerability description (type + affected asset).
2. **Asset** — the affected in-scope asset (URL / identifier).
3. **Severity (CVSS)** — CVSS vector and score (Critical/High/Medium/Low).
4. **Steps to Reproduce** — reproducible step-by-step actions (including requests / responses / payloads).
5. **Impact** — exploitability and business impact.
6. **Remediation** — fix recommendations.

For multiple findings, each gets its own section; attach a parameterizable Python PoC (requests library).
Remind the user: the report is for **manual submission** on HackerOne; this skill does not auto-submit.

## Reference Documents

- `references/hackerone-report-and-scope.md` — scope parsing reference (asset type ↔ API enum,
  eligibility tri-state, paste table shape), program rules enforcement checklist, HackerOne submission format report template.

## References — hackerone-report-and-scope

# HackerOne Report Template & Scope Parsing Reference

This reference is for the `hackerone` skill: below is the target shape for **scope parsing**
and the report template for **HackerOne submission format**. Technical terms are kept in English.

## 1. Scope Parsing Reference

### 1.1 Asset Type (human label ↔ API enum)

| Human label            | API enum                                   | pentest-flow can handle directly |
| ---------------------- | ------------------------------------------ | -------------------------------- |
| Domain / URL           | `URL`                                      | ✅                               |
| Wildcard `*.x.com`     | `WILDCARD`                                 | ✅                               |
| IP range / CIDR        | `CIDR`                                     | ⚠️ Needs confirmation (may be restricted) |
| Source code            | `SOURCE_CODE`                              | ❌ Manual                        |
| Android app            | `GOOGLE_PLAY_APP_ID` / `OTHER_APK`         | ❌ Needs specialized workflow    |
| iOS app                | `APPLE_STORE_APP_ID` / `TESTFLIGHT` / `OTHER_IPA` | ❌ Needs specialized workflow |
| Hardware               | `HARDWARE`                                 | ❌ Manual                        |
| AI Model               | `AI_MODEL`                                 | ❌ Manual                        |
| Smart Contract         | `SMART_CONTRACT`                           | ❌ Manual                        |
| Other / ASN            | `OTHER`                                    | ⚠️ Needs confirmation            |

### 1.2 Eligibility Tri-state (submission × bounty, two independent booleans)

| submission | bounty | Meaning                                    | Action             |
| ---------- | ------ | ------------------------------------------ | ------------------ |
| true       | true   | in scope, testable, bounty-eligible        | Normal testing     |
| true       | false  | in scope, testable, **no bounty**          | Normal testing (do not skip) |
| false      | —      | **out of scope**                           | **Never touch**    |

### 1.3 Target Shape for Pasted Scope Table (lenient parse)

```
In scope:
https://api.example.com        | URL       | Eligible for bounty
*.example.com                  | WILDCARD  | Eligible for bounty
app.example.com                | URL       | In scope, NOT bounty-eligible
com.example.android            | GOOGLE_PLAY_APP_ID | Eligible for bounty

Out of scope:
blog.example.com               | URL
*.corp.example.com             | WILDCARD
```

Parsing notes:
- Each line must extract at least the **asset identifier** and **in/out classification**; type and eligibility should be identified when possible.
- Column order and separators (`|`, tab, multiple spaces) may vary — match tokens leniently.
- Any line whose classification cannot be determined must be **confirmed with the user; never default to in-scope**.

## 2. Program Rules Enforcement Checklist

Check each item before testing any asset:

- **No DoS / no availability impact** — no stress testing, resource exhaustion, or mass concurrency.
- **Rate limit / automation limit** — low-speed serial; obey "no automated scanning" clauses.
- **No social engineering** — do not target personnel, no phishing.
- **Minimal impact / no PII exfiltration** — verify and stop; do not export real user data.
- Layered on top of Talon's existing `BLOCKED_PATTERNS` / `RESERVED_IP_RANGES`.

## 3. HackerOne Submission Format Report Template

Each finding uses the following structure:

```markdown
### [Title] <vulnerability type> on <asset>

**Asset:** <affected in-scope asset (URL / identifier)>

**Severity (CVSS):** <Critical | High | Medium | Low> —
`CVSS:3.1/AV:N/AC:L/PR:N/UI:N/S:U/C:H/I:H/A:N` (score: X.X)

**Steps to Reproduce:**
1. <steps, including requests / responses / payloads>
2. ...

**Impact:**
<exploitability and business impact>

**Remediation:**
<fix recommendations>

**Proof of Concept:**
(attach a parameterizable Python PoC, requests library; for verification only, non-destructive)
```

Reminder: reports are for the user to **manually submit** on HackerOne; the skill does not auto-submit.
