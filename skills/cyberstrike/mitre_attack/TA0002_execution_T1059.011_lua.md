# stage: post_exploit
# category: mitre_attack


# T1059.011 Lua

> **Sub-technique of:** T1059

## High-Level Description

Adversaries may abuse Lua commands and scripts for execution. Lua is a cross-platform scripting and programming language primarily designed for embedded use in applications. Lua can be executed on the command-line (through the stand-alone lua interpreter), via scripts (<code>.lua</code>), or from Lua-embedded programs (through the <code>struct lua_State</code>).

Lua scripts may be executed by adversaries for malicious purposes. Adversaries may incorporate, abuse, or replace existing Lua interpreters to allow for malicious Lua command execution at runtime.

## Kill Chain Phase

- Execution (TA0002)

**Platforms:** Linux, Network Devices, Windows, macOS

## What to Check

- [ ] Identify if Lua technique is applicable to target environment
- [ ] Check Linux systems for indicators of Lua
- [ ] Check Network Devices systems for indicators of Lua
- [ ] Check Windows systems for indicators of Lua
- [ ] Verify mitigations are bypassed or absent (3 known mitigations)
- [ ] Assess detection coverage (1 detection strategies)

## How to Test

### Manual Testing

1. **Identify Attack Surface**: Determine if the target environment is susceptible to Lua by examining the target platforms (Linux, Network Devices, Windows).

2. **Assess Existing Defenses**: Review whether mitigations for T1059.011 are in place. If defenses are absent or misconfigured, this technique may be exploitable.

3. **Execute Test**: Use tools and methods described in the MITRE ATT&CK page and external references below.

> **Note**: No Atomic Red Team tests available for this technique. See [Atomic Red Team GitHub](https://github.com/redcanaryco/atomic-red-team) for updates.

## Remediation Guide

### M1033 Limit Software Installation

Prevent users from installing Lua where not required.

### M1047 Audit

Inventory systems for unauthorized Lua installations.

### M1038 Execution Prevention

Denylist Lua interpreters where appropriate.

## Detection

### Detection Strategy for Lua Scripting Abuse

## Risk Assessment

| Finding                  | Severity | Impact    |
| ------------------------ | -------- | --------- |
| Lua technique applicable | Low      | Execution |

## CWE Categories

| CWE ID | Title                                  |
| ------ | -------------------------------------- |
| CWE-94 | Improper Control of Generation of Code |

## References

- [Kaspersky Lua](https://media.kasperskycontenthub.com/wp-content/uploads/sites/43/2018/03/07190154/The-ProjectSauron-APT_research_KL.pdf)
- [Lua main page](https://www.lua.org/start.html)
- [Lua state](https://pgl.yoyo.org/luai/i/lua_State)
- [Cyphort EvilBunny](https://web.archive.org/web/20150311013500/http:/www.cyphort.com/evilbunny-malware-instrumented-lua/)
- [PoetRat Lua](https://blog.talosintelligence.com/poetrat-update/)
- [Lua Proofpoint Sunseed](https://www.proofpoint.com/us/blog/threat-insight/asylum-ambuscade-state-actor-uses-compromised-private-ukrainian-military-emails)
- [Atomic Red Team - T1059.011](https://github.com/redcanaryco/atomic-red-team/tree/master/atomics/T1059.011)
- [MITRE ATT&CK - T1059.011](https://attack.mitre.org/techniques/T1059/011)
