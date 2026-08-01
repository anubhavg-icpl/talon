# stage: post_exploit
# category: mitre_attack_mobile


# T1661 Application Versioning

## High-Level Description

An adversary may push an update to a previously benign application to add malicious code. This can be accomplished by pushing an initially benign, functional application to a trusted application store, such as the Google Play Store or the Apple App Store. This allows the adversary to establish a trusted userbase that may grant permissions to the application prior to the introduction of malicious code. Then, an application update could be pushed to introduce malicious code.

This technique could also be accomplished by compromising a developer’s account. This would allow an adversary to take advantage of an existing userbase without having to establish the userbase themselves.

## Kill Chain Phase

- Initial Access (TA0027)
- Defense Evasion (TA0030)

**Platforms:** Android, iOS

## What to Check

- [ ] Identify if Application Versioning technique is applicable to target mobile environment
- [ ] Check Android devices for indicators of Application Versioning
- [ ] Check iOS devices for indicators of Application Versioning
- [ ] Verify mitigations are bypassed or absent (2 known mitigations)
- [ ] Assess detection coverage (1 detection strategies)

## How to Test

### Identify Attack Surface

Determine if the target mobile environment is susceptible to Application Versioning by examining the target platforms (Android, iOS).

### Assess Existing Defenses

Review whether mitigations for T1661 are in place. If defenses are absent or misconfigured, this technique may be exploitable.

## Remediation Guide

### M1012 Enterprise Policy

Enterprises can provision policies to mobile devices for application allow-listing, ensuring only approved applications are installed onto mobile devices.

### M1006 Use Recent OS Version

Android 11 and above implement application hibernation, which can hibernate an application that has not been used for a few months and can reset the application’s permission requests.

## Detection

### Detection of Application Versioning

## Risk Assessment

| Finding                                     | Severity | Impact         |
| ------------------------------------------- | -------- | -------------- |
| Application Versioning technique applicable | Low      | Initial Access |

## CWE Categories

| CWE ID | Title                     |
| ------ | ------------------------- |
| CWE-20 | Improper Input Validation |

## References

- [android_app_breaking_bad](https://www.welivesecurity.com/2023/05/23/android-app-breaking-bad-legitimate-screen-recording-file-exfiltration/)
- [NIST Mobile Threat Catalogue](https://pages.nist.gov/mobile-threat-catalogue/supply-chain-threats/SPC-20.html)
- [MITRE ATT&CK Mobile - T1661](https://attack.mitre.org/techniques/T1661)
