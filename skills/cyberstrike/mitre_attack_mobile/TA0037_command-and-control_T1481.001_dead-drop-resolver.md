# stage: post_exploit
# category: mitre_attack_mobile


# T1481.001 Dead Drop Resolver

> **Sub-technique of:** T1481

## High-Level Description

Adversaries may use an existing, legitimate external Web service to host information that points to additional command and control (C2) infrastructure. Adversaries may post content, known as a dead drop resolver, on Web services with embedded (and often obfuscated/encoded) domains or IP addresses. Once infected, victims will reach out to and be redirected by these resolvers.

Popular websites and social media, acting as a mechanism for C2, may give a significant amount of cover. This is due to the likelihood that hosts within a network are already communicating with them prior to a compromise. Using common services, such as those offered by Google or Twitter, makes it easier for adversaries to hide in expected noise. Web service providers commonly use SSL/TLS encryption, giving adversaries an added level of protection.

Use of a dead drop resolver may also protect back-end C2 infrastructure from discovery through malware binary analysis, or enable operational resiliency (since this infrastructure may be dynamically changed).

## Kill Chain Phase

- Command and Control (TA0037)

**Platforms:** Android, iOS

## What to Check

- [ ] Identify if Dead Drop Resolver technique is applicable to target mobile environment
- [ ] Check Android devices for indicators of Dead Drop Resolver
- [ ] Check iOS devices for indicators of Dead Drop Resolver
- [ ] Assess detection coverage (1 detection strategies)

## How to Test

### Identify Attack Surface

Determine if the target mobile environment is susceptible to Dead Drop Resolver by examining the target platforms (Android, iOS).

### Assess Existing Defenses

Review whether mitigations for T1481.001 are in place. If defenses are absent or misconfigured, this technique may be exploitable.

## Remediation Guide

No specific mitigations documented for this technique.

## Detection

### Detection of Dead Drop Resolver

## Risk Assessment

| Finding                                 | Severity | Impact              |
| --------------------------------------- | -------- | ------------------- |
| Dead Drop Resolver technique applicable | Medium   | Command And Control |

## CWE Categories

| CWE ID  | Title                              |
| ------- | ---------------------------------- |
| CWE-300 | Channel Accessible by Non-Endpoint |

## References

- [MITRE ATT&CK Mobile - T1481.001](https://attack.mitre.org/techniques/T1481/001)
