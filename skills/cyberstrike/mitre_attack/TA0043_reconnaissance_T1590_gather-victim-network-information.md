# stage: recon
# category: mitre_attack


# T1590 Gather Victim Network Information

## High-Level Description

Adversaries may gather information about the victim's networks that can be used during targeting. Information about networks may include a variety of details, including administrative data (ex: IP ranges, domain names, etc.) as well as specifics regarding its topology and operations.

Adversaries may gather this information in various ways, such as direct collection actions via Active Scanning or Phishing for Information. Information about networks may also be exposed to adversaries via online or other accessible data sets (ex: Search Open Technical Databases). Gathering this information may reveal opportunities for other forms of reconnaissance (ex: Active Scanning or Search Open Websites/Domains), establishing operational resources (ex: Acquire Infrastructure or Compromise Infrastructure), and/or initial access (ex: Trusted Relationship).

## Kill Chain Phase

- Reconnaissance (TA0043)

**Platforms:** PRE

## What to Check

- [ ] Identify if Gather Victim Network Information technique is applicable to target environment
- [ ] Check PRE systems for indicators of Gather Victim Network Information
- [ ] Verify mitigations are bypassed or absent (1 known mitigations)
- [ ] Assess detection coverage (1 detection strategies)

## How to Test

### Manual Testing

1. **Identify Attack Surface**: Determine if the target environment is susceptible to Gather Victim Network Information by examining the target platforms (PRE).

2. **Assess Existing Defenses**: Review whether mitigations for T1590 are in place. If defenses are absent or misconfigured, this technique may be exploitable.

3. **Execute Test**: Use tools and methods described in the MITRE ATT&CK page and external references below.

> **Note**: No Atomic Red Team tests available for this technique. See [Atomic Red Team GitHub](https://github.com/redcanaryco/atomic-red-team) for updates.

## Remediation Guide

### M1056 Pre-compromise

This technique cannot be easily mitigated with preventive controls since it is based on behaviors performed outside of the scope of enterprise defenses and controls. Efforts should focus on minimizing the amount and sensitivity of data available to external parties.

## Detection

### Detection of Gather Victim Network Information

## Risk Assessment

| Finding                                                | Severity | Impact         |
| ------------------------------------------------------ | -------- | -------------- |
| Gather Victim Network Information technique applicable | High     | Reconnaissance |

## CWE Categories

| CWE ID  | Title                             |
| ------- | --------------------------------- |
| CWE-200 | Exposure of Sensitive Information |

## References

- [Circl Passive DNS](https://www.circl.lu/services/passive-dns/)
- [DNS Dumpster](https://dnsdumpster.com/)
- [WHOIS](https://who.is/)
- [Atomic Red Team - T1590](https://github.com/redcanaryco/atomic-red-team/tree/master/atomics/T1590)
- [MITRE ATT&CK - T1590](https://attack.mitre.org/techniques/T1590)
