# stage: post_exploit
# category: mitre_attack_mobile


# T1636.004 SMS Messages

> **Sub-technique of:** T1636

## High-Level Description

Adversaries may utilize standard operating system APIs to gather SMS messages. On Android, this can be accomplished using the SMS Content Provider. iOS provides no standard API to access SMS messages.

If the device has been jailbroken or rooted, an adversary may be able to access SMS Messages without the user’s knowledge or approval.

## Kill Chain Phase

- Collection (TA0035)

**Platforms:** Android, iOS

## What to Check

- [ ] Identify if SMS Messages technique is applicable to target mobile environment
- [ ] Check Android devices for indicators of SMS Messages
- [ ] Check iOS devices for indicators of SMS Messages
- [ ] Verify mitigations are bypassed or absent (1 known mitigations)
- [ ] Assess detection coverage (1 detection strategies)

## How to Test

### Identify Attack Surface

Determine if the target mobile environment is susceptible to SMS Messages by examining the target platforms (Android, iOS).

### Assess Existing Defenses

Review whether mitigations for T1636.004 are in place. If defenses are absent or misconfigured, this technique may be exploitable.

## Remediation Guide

### M1011 User Guidance

Access to SMS messages is an uncommonly needed permission, so users should be instructed to use extra scrutiny when granting access to their SMS messages.

## Detection

### Detection of SMS Messages

## Risk Assessment

| Finding                           | Severity | Impact     |
| --------------------------------- | -------- | ---------- |
| SMS Messages technique applicable | Low      | Collection |

## CWE Categories

| CWE ID  | Title                             |
| ------- | --------------------------------- |
| CWE-200 | Exposure of Sensitive Information |

## References

- [NIST Mobile Threat Catalogue](https://pages.nist.gov/mobile-threat-catalogue/application-threats/APP-13.html)
- [MITRE ATT&CK Mobile - T1636.004](https://attack.mitre.org/techniques/T1636/004)
