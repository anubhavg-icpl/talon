# stage: post_exploit
# category: mitre_attack


# T1583.006 Web Services

> **Sub-technique of:** T1583

## High-Level Description

Adversaries may register for web services that can be used during targeting. A variety of popular websites exist for adversaries to register for a web-based service that can be abused during later stages of the adversary lifecycle, such as during Command and Control (Web Service), Exfiltration Over Web Service, or Phishing. Using common services, such as those offered by Google, GitHub, or Twitter, makes it easier for adversaries to hide in expected noise. By utilizing a web service, adversaries can make it difficult to physically tie back operations to them.

## Kill Chain Phase

- Resource Development (TA0042)

**Platforms:** PRE

## What to Check

- [ ] Identify if Web Services technique is applicable to target environment
- [ ] Check PRE systems for indicators of Web Services
- [ ] Verify mitigations are bypassed or absent (1 known mitigations)
- [ ] Assess detection coverage (1 detection strategies)

## How to Test

### Manual Testing

1. **Identify Attack Surface**: Determine if the target environment is susceptible to Web Services by examining the target platforms (PRE).

2. **Assess Existing Defenses**: Review whether mitigations for T1583.006 are in place. If defenses are absent or misconfigured, this technique may be exploitable.

3. **Execute Test**: Use tools and methods described in the MITRE ATT&CK page and external references below.

> **Note**: No Atomic Red Team tests available for this technique. See [Atomic Red Team GitHub](https://github.com/redcanaryco/atomic-red-team) for updates.

## Remediation Guide

### M1056 Pre-compromise

This technique cannot be easily mitigated with preventive controls since it is based on behaviors performed outside of the scope of enterprise defenses and controls.

## Detection

### Detection of Web Services

## Risk Assessment

| Finding                           | Severity | Impact               |
| --------------------------------- | -------- | -------------------- |
| Web Services technique applicable | High     | Resource Development |

## CWE Categories

| CWE ID | Title                 |
| ------ | --------------------- |
| N/A    | No direct CWE mapping |

## References

- [Hacker News GitHub Abuse 2024](https://thehackernews.com/expert-insights/2024/05/github-abuse-flaw-shows-why-we-cant.html)
- [FireEye APT29](https://services.google.com/fh/files/misc/rpt-apt29-hammertoss-stealthy-tactics-define-en.pdf)
- [ThreatConnect Infrastructure Dec 2020](https://threatconnect.com/blog/infrastructure-research-hunting/)
- [Atomic Red Team - T1583.006](https://github.com/redcanaryco/atomic-red-team/tree/master/atomics/T1583.006)
- [MITRE ATT&CK - T1583.006](https://attack.mitre.org/techniques/T1583/006)
