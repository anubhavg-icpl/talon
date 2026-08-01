# stage: post_exploit
# category: mitre_attack_ics


# T0801 Monitor Process State

## High-Level Description

Adversaries may gather information about the physical process state. This information may be used to gain more information about the process itself or used as a trigger for malicious actions. The sources of process state information may vary such as, OPC tags, historian data, specific PLC block information, or network traffic.

## Kill Chain Phase

- Collection (TA0100)

**Platforms:** ICS

## What to Check

- [ ] Identify if Monitor Process State technique is applicable to target ICS environment
- [ ] Check ICS/SCADA systems for indicators of Monitor Process State
- [ ] Verify mitigations are bypassed or absent (1 known mitigations)
- [ ] Assess detection coverage (1 detection strategies)

## How to Test

### Identify Attack Surface

Determine if the target ICS/SCADA environment is susceptible to Monitor Process State by examining operational technology systems and network architecture.

### Assess Existing Defenses

Review whether mitigations for T0801 are in place. If defenses are absent or misconfigured, this technique may be exploitable in the ICS environment.

## Remediation Guide

### M0816 Mitigation Limited or Not Effective

This type of attack technique cannot be easily mitigated with preventive controls since it is based on the abuse of system features.

## Detection

### Detection of Monitor Process State

## Risk Assessment

| Finding                                    | Severity | Impact     |
| ------------------------------------------ | -------- | ---------- |
| Monitor Process State technique applicable | Medium   | Collection |

## CWE Categories

| CWE ID  | Title                             |
| ------- | --------------------------------- |
| CWE-200 | Exposure of Sensitive Information |

## References

- [MITRE ATT&CK ICS - T0801](https://attack.mitre.org/techniques/T0801)
