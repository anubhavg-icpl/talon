# stage: post_exploit
# category: mitre_attack


# T1491 Defacement

## High-Level Description

Adversaries may modify visual content available internally or externally to an enterprise network, thus affecting the integrity of the original content. Reasons for Defacement include delivering messaging, intimidation, or claiming (possibly false) credit for an intrusion. Disturbing or offensive images may be used as a part of Defacement in order to cause user discomfort, or to pressure compliance with accompanying messages.

## Kill Chain Phase

- Impact (TA0040)

**Platforms:** Windows, IaaS, Linux, macOS, ESXi

## What to Check

- [ ] Identify if Defacement technique is applicable to target environment
- [ ] Check Windows systems for indicators of Defacement
- [ ] Check IaaS systems for indicators of Defacement
- [ ] Check Linux systems for indicators of Defacement
- [ ] Verify mitigations are bypassed or absent (1 known mitigations)
- [ ] Assess detection coverage (1 detection strategies)

## How to Test

### Manual Testing

1. **Identify Attack Surface**: Determine if the target environment is susceptible to Defacement by examining the target platforms (Windows, IaaS, Linux).

2. **Assess Existing Defenses**: Review whether mitigations for T1491 are in place. If defenses are absent or misconfigured, this technique may be exploitable.

3. **Execute Test**: Use tools and methods described in the MITRE ATT&CK page and external references below.

> **Note**: No Atomic Red Team tests available for this technique. See [Atomic Red Team GitHub](https://github.com/redcanaryco/atomic-red-team) for updates.

## Remediation Guide

### M1053 Data Backup

Consider implementing IT disaster recovery plans that contain procedures for taking regular data backups that can be used to restore organizational data. Ensure backups are stored off system and is protected from common methods adversaries may use to gain access and destroy the backups to prevent recovery.

## Detection

### Defacement via File and Web Content Modification Across Platforms

## Risk Assessment

| Finding                         | Severity | Impact |
| ------------------------------- | -------- | ------ |
| Defacement technique applicable | Low      | Impact |

## CWE Categories

| CWE ID  | Title                             |
| ------- | --------------------------------- |
| CWE-400 | Uncontrolled Resource Consumption |

## References

- [Atomic Red Team - T1491](https://github.com/redcanaryco/atomic-red-team/tree/master/atomics/T1491)
- [MITRE ATT&CK - T1491](https://attack.mitre.org/techniques/T1491)
