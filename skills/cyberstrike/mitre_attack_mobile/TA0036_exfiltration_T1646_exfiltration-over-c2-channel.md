# stage: post_exploit
# category: mitre_attack_mobile


# T1646 Exfiltration Over C2 Channel

## High-Level Description

Adversaries may steal data by exfiltrating it over an existing command and control channel. Stolen data is encoded into the normal communications channel using the same protocol as command and control communications.

## Kill Chain Phase

- Exfiltration (TA0036)

**Platforms:** Android, iOS

## What to Check

- [ ] Identify if Exfiltration Over C2 Channel technique is applicable to target mobile environment
- [ ] Check Android devices for indicators of Exfiltration Over C2 Channel
- [ ] Check iOS devices for indicators of Exfiltration Over C2 Channel
- [ ] Assess detection coverage (1 detection strategies)

## How to Test

### Identify Attack Surface

Determine if the target mobile environment is susceptible to Exfiltration Over C2 Channel by examining the target platforms (Android, iOS).

### Assess Existing Defenses

Review whether mitigations for T1646 are in place. If defenses are absent or misconfigured, this technique may be exploitable.

## Remediation Guide

No specific mitigations documented for this technique.

## Detection

### Detection of Exfiltration Over C2 Channel

## Risk Assessment

| Finding                                           | Severity | Impact       |
| ------------------------------------------------- | -------- | ------------ |
| Exfiltration Over C2 Channel technique applicable | Low      | Exfiltration |

## CWE Categories

| CWE ID  | Title                             |
| ------- | --------------------------------- |
| CWE-200 | Exposure of Sensitive Information |

## References

- [NIST Mobile Threat Catalogue](https://pages.nist.gov/mobile-threat-catalogue/application-threats/APP-29.html)
- [MITRE ATT&CK Mobile - T1646](https://attack.mitre.org/techniques/T1646)
