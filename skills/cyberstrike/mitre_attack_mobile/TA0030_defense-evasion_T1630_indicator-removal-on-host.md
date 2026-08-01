# stage: post_exploit
# category: mitre_attack_mobile


# T1630 Indicator Removal on Host

## High-Level Description

Adversaries may delete, alter, or hide generated artifacts on a device, including files, jailbreak status, or the malicious application itself. These actions may interfere with event collection, reporting, or other notifications used to detect intrusion activity. This may compromise the integrity of mobile security solutions by causing notable events or information to go unreported.

## Kill Chain Phase

- Defense Evasion (TA0030)

**Platforms:** iOS, Android

## What to Check

- [ ] Identify if Indicator Removal on Host technique is applicable to target mobile environment
- [ ] Check iOS devices for indicators of Indicator Removal on Host
- [ ] Check Android devices for indicators of Indicator Removal on Host
- [ ] Verify mitigations are bypassed or absent (3 known mitigations)
- [ ] Assess detection coverage (1 detection strategies)

## How to Test

### Identify Attack Surface

Determine if the target mobile environment is susceptible to Indicator Removal on Host by examining the target platforms (iOS, Android).

### Assess Existing Defenses

Review whether mitigations for T1630 are in place. If defenses are absent or misconfigured, this technique may be exploitable.

## Remediation Guide

### M1001 Security Updates

Security updates typically provide patches for vulnerabilities that could be abused by malicious applications.

### M1011 User Guidance

Inform users that device rooting or granting unnecessary access to the accessibility service presents security risks that could be taken advantage of without their knowledge.

### M1002 Attestation

Attestation can detect unauthorized modifications to devices. Mobile security software can then use this information and take appropriate mitigation action.

## Detection

### Detection of Indicator Removal on Host

## Risk Assessment

| Finding                                        | Severity | Impact          |
| ---------------------------------------------- | -------- | --------------- |
| Indicator Removal on Host technique applicable | Low      | Defense Evasion |

## CWE Categories

| CWE ID  | Title                        |
| ------- | ---------------------------- |
| CWE-693 | Protection Mechanism Failure |

## References

- [NIST Mobile Threat Catalogue](https://pages.nist.gov/mobile-threat-catalogue/application-threats/APP-43.html)
- [MITRE ATT&CK Mobile - T1630](https://attack.mitre.org/techniques/T1630)
