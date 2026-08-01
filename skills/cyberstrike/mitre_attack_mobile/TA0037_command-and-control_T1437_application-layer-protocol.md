# stage: post_exploit
# category: mitre_attack_mobile


# T1437 Application Layer Protocol

## High-Level Description

Adversaries may communicate using application layer protocols to avoid detection/network filtering by blending in with existing traffic. Commands to the mobile device, and often the results of those commands, will be embedded within the protocol traffic between the mobile device and server.

Adversaries may utilize many different protocols, including those used for web browsing, transferring files, electronic mail, or DNS.

## Kill Chain Phase

- Command and Control (TA0037)

**Platforms:** Android, iOS

## What to Check

- [ ] Identify if Application Layer Protocol technique is applicable to target mobile environment
- [ ] Check Android devices for indicators of Application Layer Protocol
- [ ] Check iOS devices for indicators of Application Layer Protocol
- [ ] Assess detection coverage (1 detection strategies)

## How to Test

### Identify Attack Surface

Determine if the target mobile environment is susceptible to Application Layer Protocol by examining the target platforms (Android, iOS).

### Assess Existing Defenses

Review whether mitigations for T1437 are in place. If defenses are absent or misconfigured, this technique may be exploitable.

## Remediation Guide

No specific mitigations documented for this technique.

## Detection

### Detection of Application Layer Protocol

## Risk Assessment

| Finding                                         | Severity | Impact              |
| ----------------------------------------------- | -------- | ------------------- |
| Application Layer Protocol technique applicable | Low      | Command And Control |

## CWE Categories

| CWE ID  | Title                              |
| ------- | ---------------------------------- |
| CWE-300 | Channel Accessible by Non-Endpoint |

## References

- [NIST Mobile Threat Catalogue](https://pages.nist.gov/mobile-threat-catalogue/application-threats/APP-29.html)
- [MITRE ATT&CK Mobile - T1437](https://attack.mitre.org/techniques/T1437)
