# stage: post_exploit
# category: mitre_attack_ics


# T0872 Indicator Removal on Host

## High-Level Description

Adversaries may attempt to remove indicators of their presence on a system in an effort to cover their tracks. In cases where an adversary may feel detection is imminent, they may try to overwrite, delete, or cover up changes they have made to the device.

## Kill Chain Phase

- Evasion (TA0103)

**Platforms:** ICS

## What to Check

- [ ] Identify if Indicator Removal on Host technique is applicable to target ICS environment
- [ ] Check ICS/SCADA systems for indicators of Indicator Removal on Host
- [ ] Verify mitigations are bypassed or absent (1 known mitigations)
- [ ] Assess detection coverage (1 detection strategies)

## How to Test

### Identify Attack Surface

Determine if the target ICS/SCADA environment is susceptible to Indicator Removal on Host by examining operational technology systems and network architecture.

### Assess Existing Defenses

Review whether mitigations for T0872 are in place. If defenses are absent or misconfigured, this technique may be exploitable in the ICS environment.

## Remediation Guide

### M0922 Restrict File and Directory Permissions

Protect files stored locally with proper permissions to limit opportunities for adversaries to remove indicators of their activity on the system.

## Detection

### Detection of Indicator Removal on Host

## Risk Assessment

| Finding                                        | Severity | Impact  |
| ---------------------------------------------- | -------- | ------- |
| Indicator Removal on Host technique applicable | Low      | Evasion |

## CWE Categories

| CWE ID  | Title                        |
| ------- | ---------------------------- |
| CWE-693 | Protection Mechanism Failure |

## References

- [MITRE ATT&CK ICS - T0872](https://attack.mitre.org/techniques/T0872)
