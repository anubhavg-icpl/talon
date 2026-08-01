# stage: post_exploit
# category: mitre_attack


# T1102 Web Service

## High-Level Description

Adversaries may use an existing, legitimate external Web service as a means for relaying data to/from a compromised system. Popular websites, cloud services, and social media acting as a mechanism for C2 may give a significant amount of cover due to the likelihood that hosts within a network are already communicating with them prior to a compromise. Using common services, such as those offered by Google, Microsoft, or Twitter, makes it easier for adversaries to hide in expected noise. Web service providers commonly use SSL/TLS encryption, giving adversaries an added level of protection.

Use of Web services may also protect back-end C2 infrastructure from discovery through malware binary analysis while also enabling operational resiliency (since this infrastructure may be dynamically changed).

## Kill Chain Phase

- Command and Control (TA0011)

**Platforms:** ESXi, Linux, Windows, macOS

## What to Check

- [ ] Identify if Web Service technique is applicable to target environment
- [ ] Check ESXi systems for indicators of Web Service
- [ ] Check Linux systems for indicators of Web Service
- [ ] Check Windows systems for indicators of Web Service
- [ ] Verify mitigations are bypassed or absent (2 known mitigations)
- [ ] Assess detection coverage (1 detection strategies)

## How to Test

### Manual Testing

1. **Identify Attack Surface**: Determine if the target environment is susceptible to Web Service by examining the target platforms (ESXi, Linux, Windows).

2. **Assess Existing Defenses**: Review whether mitigations for T1102 are in place. If defenses are absent or misconfigured, this technique may be exploitable.

3. **Execute Test**: Use tools and methods described in the MITRE ATT&CK page and external references below.

> **Note**: No Atomic Red Team tests available for this technique. See [Atomic Red Team GitHub](https://github.com/redcanaryco/atomic-red-team) for updates.

## Remediation Guide

### M1031 Network Intrusion Prevention

Network intrusion detection and prevention systems that use network signatures to identify traffic for specific adversary malware can be used to mitigate activity at the network level.

### M1021 Restrict Web-Based Content

Web proxies can be used to enforce external network communication policy that prevents use of unauthorized external services.

## Detection

### Suspicious Use of Web Services for C2

## Risk Assessment

| Finding                          | Severity | Impact              |
| -------------------------------- | -------- | ------------------- |
| Web Service technique applicable | Medium   | Command And Control |

## CWE Categories

| CWE ID  | Title                              |
| ------- | ---------------------------------- |
| CWE-300 | Channel Accessible by Non-Endpoint |

## References

- [Broadcom BirdyClient Microsoft Graph API 2024](https://www.broadcom.com/support/security-center/protection-bulletin/birdyclient-malware-leverages-microsoft-graph-api-for-c-c-communication)
- [University of Birmingham C2](https://arxiv.org/ftp/arxiv/papers/1408/1408.1136.pdf)
- [Atomic Red Team - T1102](https://github.com/redcanaryco/atomic-red-team/tree/master/atomics/T1102)
- [MITRE ATT&CK - T1102](https://attack.mitre.org/techniques/T1102)
