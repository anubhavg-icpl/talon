# stage: post_exploit
# category: mitre_attack_mobile


# T1481.002 Bidirectional Communication

> **Sub-technique of:** T1481

## High-Level Description

Adversaries may use an existing, legitimate external Web service channel as a means for sending commands to and receiving output from a compromised system. Compromised systems may leverage popular websites and social media to host command and control (C2) instructions. Those infected systems can then send the output from those commands back over that Web service channel. The return traffic may occur in a variety of ways, depending on the Web service being utilized. For example, the return traffic may take the form of the compromised system posting a comment on a forum, issuing a pull request to development project, updating a document hosted on a Web service, or by sending a Tweet.

Popular websites and social media, acting as a mechanism for C2, may give a significant amount of cover. This is due to the likelihood that hosts within a network are already communicating with them prior to a compromise. Using common services, such as those offered by Google or Twitter, makes it easier for adversaries to hide in expected noise. Web service providers commonly use SSL/TLS encryption, giving adversaries an added level of protection.

## Kill Chain Phase

- Command and Control (TA0037)

**Platforms:** Android, iOS

## What to Check

- [ ] Identify if Bidirectional Communication technique is applicable to target mobile environment
- [ ] Check Android devices for indicators of Bidirectional Communication
- [ ] Check iOS devices for indicators of Bidirectional Communication
- [ ] Assess detection coverage (1 detection strategies)

## How to Test

### Identify Attack Surface

Determine if the target mobile environment is susceptible to Bidirectional Communication by examining the target platforms (Android, iOS).

### Assess Existing Defenses

Review whether mitigations for T1481.002 are in place. If defenses are absent or misconfigured, this technique may be exploitable.

## Remediation Guide

No specific mitigations documented for this technique.

## Detection

### Detection of Bidirectional Communication

## Risk Assessment

| Finding                                          | Severity | Impact              |
| ------------------------------------------------ | -------- | ------------------- |
| Bidirectional Communication technique applicable | Low      | Command And Control |

## CWE Categories

| CWE ID  | Title                              |
| ------- | ---------------------------------- |
| CWE-300 | Channel Accessible by Non-Endpoint |

## References

- [MITRE ATT&CK Mobile - T1481.002](https://attack.mitre.org/techniques/T1481/002)
