# stage: post_exploit
# category: mitre_attack_mobile


# T1637 Dynamic Resolution

## High-Level Description

Adversaries may dynamically establish connections to command and control infrastructure to evade common detections and remediations. This may be achieved by using malware that shares a common algorithm with the infrastructure the adversary uses to receive the malware's communications. This algorithm can be used to dynamically adjust parameters such as the domain name, IP address, or port number the malware uses for command and control.

## Kill Chain Phase

- Command and Control (TA0037)

**Platforms:** Android, iOS

## What to Check

- [ ] Identify if Dynamic Resolution technique is applicable to target mobile environment
- [ ] Check Android devices for indicators of Dynamic Resolution
- [ ] Check iOS devices for indicators of Dynamic Resolution
- [ ] Assess detection coverage (1 detection strategies)

## How to Test

### Identify Attack Surface

Determine if the target mobile environment is susceptible to Dynamic Resolution by examining the target platforms (Android, iOS).

### Assess Existing Defenses

Review whether mitigations for T1637 are in place. If defenses are absent or misconfigured, this technique may be exploitable.

## Remediation Guide

No specific mitigations documented for this technique.

## Detection

### Detection of Dynamic Resolution

## Risk Assessment

| Finding                                 | Severity | Impact              |
| --------------------------------------- | -------- | ------------------- |
| Dynamic Resolution technique applicable | Low      | Command And Control |

## CWE Categories

| CWE ID  | Title                              |
| ------- | ---------------------------------- |
| CWE-300 | Channel Accessible by Non-Endpoint |

## References

- [Data Driven Security DGA](https://datadrivensecurity.info/blog/posts/2014/Oct/dga-part2/)
- [MITRE ATT&CK Mobile - T1637](https://attack.mitre.org/techniques/T1637)
