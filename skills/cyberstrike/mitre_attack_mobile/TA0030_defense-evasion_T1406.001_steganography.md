# stage: post_exploit
# category: mitre_attack_mobile


# T1406.001 Steganography

> **Sub-technique of:** T1406

## High-Level Description

Adversaries may use steganography techniques in order to prevent the detection of hidden information. Steganographic techniques can be used to hide data in digital media such as images, audio tracks, video clips, or text files.

## Kill Chain Phase

- Defense Evasion (TA0030)

**Platforms:** Android

## What to Check

- [ ] Identify if Steganography technique is applicable to target mobile environment
- [ ] Check Android devices for indicators of Steganography
- [ ] Assess detection coverage (1 detection strategies)

## How to Test

### Identify Attack Surface

Determine if the target mobile environment is susceptible to Steganography by examining the target platforms (Android).

### Assess Existing Defenses

Review whether mitigations for T1406.001 are in place. If defenses are absent or misconfigured, this technique may be exploitable.

## Remediation Guide

No specific mitigations documented for this technique.

## Detection

### Detection of Steganography

## Risk Assessment

| Finding                            | Severity | Impact          |
| ---------------------------------- | -------- | --------------- |
| Steganography technique applicable | Low      | Defense Evasion |

## CWE Categories

| CWE ID  | Title                        |
| ------- | ---------------------------- |
| CWE-693 | Protection Mechanism Failure |

## References

- [MITRE ATT&CK Mobile - T1406.001](https://attack.mitre.org/techniques/T1406/001)
