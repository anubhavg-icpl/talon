# stage: recon
# category: mitre_attack_mobile


# T1422.001 Internet Connection Discovery

> **Sub-technique of:** T1422

## High-Level Description

Adversaries may check for Internet connectivity on compromised systems. This may be performed during automated discovery and can be accomplished in numerous ways such as using `adb shell netstat` for Android.

Adversaries may use the results and responses from these requests to determine if the mobile devices are capable of communicating with adversary-owned C2 servers before attempting to connect to them. The results may also be used to identify routes, redirectors, and proxy servers.

## Kill Chain Phase

- Discovery (TA0032)

**Platforms:** Android, iOS

## What to Check

- [ ] Identify if Internet Connection Discovery technique is applicable to target mobile environment
- [ ] Check Android devices for indicators of Internet Connection Discovery
- [ ] Check iOS devices for indicators of Internet Connection Discovery
- [ ] Verify mitigations are bypassed or absent (1 known mitigations)
- [ ] Assess detection coverage (1 detection strategies)

## How to Test

### Identify Attack Surface

Determine if the target mobile environment is susceptible to Internet Connection Discovery by examining the target platforms (Android, iOS).

### Assess Existing Defenses

Review whether mitigations for T1422.001 are in place. If defenses are absent or misconfigured, this technique may be exploitable.

## Remediation Guide

### M1009 Encrypt Network Traffic

Ensure that traffic is encrypted to reduce adversaries’ ability to intercept, decrypt and manipulate traffic.

## Detection

### Detection of Internet Connection Discovery

## Risk Assessment

| Finding                                            | Severity | Impact    |
| -------------------------------------------------- | -------- | --------- |
| Internet Connection Discovery technique applicable | Medium   | Discovery |

## CWE Categories

| CWE ID  | Title                             |
| ------- | --------------------------------- |
| CWE-200 | Exposure of Sensitive Information |

## References

- [adb_commands](https://gist.github.com/Pulimet/5013acf2cd5b28e55036c82c91bd56d8)
- [MITRE ATT&CK Mobile - T1422.001](https://attack.mitre.org/techniques/T1422/001)
