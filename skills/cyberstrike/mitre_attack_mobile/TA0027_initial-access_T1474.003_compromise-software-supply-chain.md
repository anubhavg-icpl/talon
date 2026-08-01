# stage: post_exploit
# category: mitre_attack_mobile


# T1474.003 Compromise Software Supply Chain

> **Sub-technique of:** T1474

## High-Level Description

Adversaries may manipulate application software prior to receipt by a final consumer for the purpose of data or system compromise. Supply chain compromise of software can take place in a number of ways, including manipulation of the application source code, manipulation of the update/distribution mechanism for that software, or replacing compiled releases with a modified version.

## Kill Chain Phase

- Initial Access (TA0027)

**Platforms:** Android, iOS

## What to Check

- [ ] Identify if Compromise Software Supply Chain technique is applicable to target mobile environment
- [ ] Check Android devices for indicators of Compromise Software Supply Chain
- [ ] Check iOS devices for indicators of Compromise Software Supply Chain
- [ ] Verify mitigations are bypassed or absent (2 known mitigations)
- [ ] Assess detection coverage (1 detection strategies)

## How to Test

### Identify Attack Surface

Determine if the target mobile environment is susceptible to Compromise Software Supply Chain by examining the target platforms (Android, iOS).

### Assess Existing Defenses

Review whether mitigations for T1474.003 are in place. If defenses are absent or misconfigured, this technique may be exploitable.

## Remediation Guide

### M1004 System Partition Integrity

Ensure Verified Boot is enabled on devices with that capability.

### M1001 Security Updates

Security updates may contain patches that inhibit system software compromises.

## Detection

### Detection of Compromise Software Supply Chain

## Risk Assessment

| Finding                                               | Severity | Impact         |
| ----------------------------------------------------- | -------- | -------------- |
| Compromise Software Supply Chain technique applicable | Low      | Initial Access |

## CWE Categories

| CWE ID | Title                     |
| ------ | ------------------------- |
| CWE-20 | Improper Input Validation |

## References

- [NIST Mobile Threat Catalogue](https://pages.nist.gov/mobile-threat-catalogue/supply-chain-threats/SPC-4.html)
- [NIST Mobile Threat Catalogue](https://pages.nist.gov/mobile-threat-catalogue/supply-chain-threats/SPC-11.html)
- [NIST Mobile Threat Catalogue](https://pages.nist.gov/mobile-threat-catalogue/supply-chain-threats/SPC-12.html)
- [NIST Mobile Threat Catalogue](https://pages.nist.gov/mobile-threat-catalogue/supply-chain-threats/SPC-18.html)
- [NIST Mobile Threat Catalogue](https://pages.nist.gov/mobile-threat-catalogue/supply-chain-threats/SPC-20.html)
- [MITRE ATT&CK Mobile - T1474.003](https://attack.mitre.org/techniques/T1474/003)
