# stage: post_exploit
# category: mitre_attack_mobile


# T1437.001 Web Protocols

> **Sub-technique of:** T1437

## High-Level Description

Adversaries may communicate using application layer protocols associated with web protocols traffic to avoid detection/network filtering by blending in with existing traffic. Commands to remote mobile devices, and often the results of those commands, will be embedded within the protocol traffic between the mobile client and server.

Web protocols such as HTTP and HTTPS are used for web traffic as well as well as notification services native to mobile messaging services such as Google Cloud Messaging (GCM) and newly, Firebase Cloud Messaging (FCM), (GCM/FCM: two-way communication) and Apple Push Notification Service (APNS; one-way server-to-device). Such notification services leverage HTTP/S via the respective API and are commonly abused on Android and iOS respectively in order blend in with routine device traffic making it difficult for enterprises to inspect.

## Kill Chain Phase

- Command and Control (TA0037)

**Platforms:** Android, iOS

## What to Check

- [ ] Identify if Web Protocols technique is applicable to target mobile environment
- [ ] Check Android devices for indicators of Web Protocols
- [ ] Check iOS devices for indicators of Web Protocols
- [ ] Assess detection coverage (1 detection strategies)

## How to Test

### Identify Attack Surface

Determine if the target mobile environment is susceptible to Web Protocols by examining the target platforms (Android, iOS).

### Assess Existing Defenses

Review whether mitigations for T1437.001 are in place. If defenses are absent or misconfigured, this technique may be exploitable.

## Remediation Guide

No specific mitigations documented for this technique.

## Detection

### Detection of Web Protocols

## Risk Assessment

| Finding                            | Severity | Impact              |
| ---------------------------------- | -------- | ------------------- |
| Web Protocols technique applicable | Low      | Command And Control |

## CWE Categories

| CWE ID  | Title                              |
| ------- | ---------------------------------- |
| CWE-300 | Channel Accessible by Non-Endpoint |

## References

- [NIST Mobile Threat Catalogue](https://pages.nist.gov/mobile-threat-catalogue/application-threats/APP-29.html)
- [MITRE ATT&CK Mobile - T1437.001](https://attack.mitre.org/techniques/T1437/001)
