# stage: post_exploit
# category: mitre_attack_mobile


# T1471 Data Encrypted for Impact

## High-Level Description

An adversary may encrypt files stored on a mobile device to prevent the user from accessing them. This may be done in order to extract monetary compensation from a victim in exchange for decryption or a decryption key (ransomware) or to render data permanently inaccessible in cases where the key is not saved or transmitted.

## Kill Chain Phase

- Impact (TA0034)

**Platforms:** Android

## What to Check

- [ ] Identify if Data Encrypted for Impact technique is applicable to target mobile environment
- [ ] Check Android devices for indicators of Data Encrypted for Impact
- [ ] Assess detection coverage (1 detection strategies)

## How to Test

### Identify Attack Surface

Determine if the target mobile environment is susceptible to Data Encrypted for Impact by examining the target platforms (Android).

### Assess Existing Defenses

Review whether mitigations for T1471 are in place. If defenses are absent or misconfigured, this technique may be exploitable.

## Remediation Guide

No specific mitigations documented for this technique.

## Detection

### Detection of Data Encrypted for Impact

## Risk Assessment

| Finding                                        | Severity | Impact |
| ---------------------------------------------- | -------- | ------ |
| Data Encrypted for Impact technique applicable | Low      | Impact |

## CWE Categories

| CWE ID  | Title                             |
| ------- | --------------------------------- |
| CWE-400 | Uncontrolled Resource Consumption |

## References

- [NIST Mobile Threat Catalogue](https://pages.nist.gov/mobile-threat-catalogue/application-threats/APP-28.html)
- [MITRE ATT&CK Mobile - T1471](https://attack.mitre.org/techniques/T1471)
