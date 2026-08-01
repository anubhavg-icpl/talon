# stage: recon
# category: mitre_attack_mobile


# T1422.002 Wi-Fi Discovery

> **Sub-technique of:** T1422

## High-Level Description

Adversaries may search for information about Wi-Fi networks, such as network names and passwords, on compromised systems. Adversaries may use Wi-Fi information as part of Discovery or Credential Access activity to support both ongoing and future campaigns.

## Kill Chain Phase

- Discovery (TA0032)

**Platforms:** Android, iOS

## What to Check

- [ ] Identify if Wi-Fi Discovery technique is applicable to target mobile environment
- [ ] Check Android devices for indicators of Wi-Fi Discovery
- [ ] Check iOS devices for indicators of Wi-Fi Discovery
- [ ] Verify mitigations are bypassed or absent (1 known mitigations)
- [ ] Assess detection coverage (1 detection strategies)

## How to Test

### Identify Attack Surface

Determine if the target mobile environment is susceptible to Wi-Fi Discovery by examining the target platforms (Android, iOS).

### Assess Existing Defenses

Review whether mitigations for T1422.002 are in place. If defenses are absent or misconfigured, this technique may be exploitable.

## Remediation Guide

### M1006 Use Recent OS Version

Android 10 introduced changes that prevent normal applications from accessing sensitive device identifiers.

## Detection

### Detection of Wi-Fi Discovery

## Risk Assessment

| Finding                              | Severity | Impact    |
| ------------------------------------ | -------- | --------- |
| Wi-Fi Discovery technique applicable | High     | Discovery |

## CWE Categories

| CWE ID  | Title                             |
| ------- | --------------------------------- |
| CWE-200 | Exposure of Sensitive Information |

## References

- [MITRE ATT&CK Mobile - T1422.002](https://attack.mitre.org/techniques/T1422/002)
