# stage: post_exploit
# category: mitre_attack


# T1567.001 Exfiltration to Code Repository

> **Sub-technique of:** T1567

## High-Level Description

Adversaries may exfiltrate data to a code repository rather than over their primary command and control channel. Code repositories are often accessible via an API (ex: https://api.github.com). Access to these APIs are often over HTTPS, which gives the adversary an additional level of protection.

Exfiltration to a code repository can also provide a significant amount of cover to the adversary if it is a popular service already used by hosts within the network.

## Kill Chain Phase

- Exfiltration (TA0010)

**Platforms:** Linux, macOS, Windows, ESXi

## What to Check

- [ ] Identify if Exfiltration to Code Repository technique is applicable to target environment
- [ ] Check Linux systems for indicators of Exfiltration to Code Repository
- [ ] Check macOS systems for indicators of Exfiltration to Code Repository
- [ ] Check Windows systems for indicators of Exfiltration to Code Repository
- [ ] Verify mitigations are bypassed or absent (1 known mitigations)
- [ ] Assess detection coverage (1 detection strategies)

## How to Test

### Manual Testing

1. **Identify Attack Surface**: Determine if the target environment is susceptible to Exfiltration to Code Repository by examining the target platforms (Linux, macOS, Windows).

2. **Assess Existing Defenses**: Review whether mitigations for T1567.001 are in place. If defenses are absent or misconfigured, this technique may be exploitable.

3. **Execute Test**: Use tools and methods described in the MITRE ATT&CK page and external references below.

> **Note**: No Atomic Red Team tests available for this technique. See [Atomic Red Team GitHub](https://github.com/redcanaryco/atomic-red-team) for updates.

## Remediation Guide

### M1021 Restrict Web-Based Content

Web proxies can be used to enforce an external network communication policy that prevents use of unauthorized external services.

## Detection

### Detection Strategy for Exfiltration to Code Repository

## Risk Assessment

| Finding                                              | Severity | Impact       |
| ---------------------------------------------------- | -------- | ------------ |
| Exfiltration to Code Repository technique applicable | Low      | Exfiltration |

## CWE Categories

| CWE ID  | Title                             |
| ------- | --------------------------------- |
| CWE-200 | Exposure of Sensitive Information |

## References

- [Atomic Red Team - T1567.001](https://github.com/redcanaryco/atomic-red-team/tree/master/atomics/T1567.001)
- [MITRE ATT&CK - T1567.001](https://attack.mitre.org/techniques/T1567/001)
