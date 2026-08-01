# stage: post_exploit
# category: mitre_attack


# T1074.002 Remote Data Staging

> **Sub-technique of:** T1074

## High-Level Description

Adversaries may stage data collected from multiple systems in a central location or directory on one system prior to Exfiltration. Data may be kept in separate files or combined into one file through techniques such as Archive Collected Data. Interactive command shells may be used, and common functionality within cmd and bash may be used to copy data into a staging location.

In cloud environments, adversaries may stage data within a particular instance or virtual machine before exfiltration. An adversary may Create Cloud Instance and stage data in that instance.

By staging data on one system prior to Exfiltration, adversaries can minimize the number of connections made to their C2 server and better evade detection.

## Kill Chain Phase

- Collection (TA0009)

**Platforms:** Windows, IaaS, Linux, macOS, ESXi

## What to Check

- [ ] Identify if Remote Data Staging technique is applicable to target environment
- [ ] Check Windows systems for indicators of Remote Data Staging
- [ ] Check IaaS systems for indicators of Remote Data Staging
- [ ] Check Linux systems for indicators of Remote Data Staging
- [ ] Assess detection coverage (1 detection strategies)

## How to Test

### Manual Testing

1. **Identify Attack Surface**: Determine if the target environment is susceptible to Remote Data Staging by examining the target platforms (Windows, IaaS, Linux).

2. **Assess Existing Defenses**: Review whether mitigations for T1074.002 are in place. If defenses are absent or misconfigured, this technique may be exploitable.

3. **Execute Test**: Use tools and methods described in the MITRE ATT&CK page and external references below.

> **Note**: No Atomic Red Team tests available for this technique. See [Atomic Red Team GitHub](https://github.com/redcanaryco/atomic-red-team) for updates.

## Remediation Guide

No specific mitigations documented for this technique.

## Detection

### Detection of Remote Data Staging Prior to Exfiltration

## Risk Assessment

| Finding                                  | Severity | Impact     |
| ---------------------------------------- | -------- | ---------- |
| Remote Data Staging technique applicable | Low      | Collection |

## CWE Categories

| CWE ID  | Title                             |
| ------- | --------------------------------- |
| CWE-200 | Exposure of Sensitive Information |

## References

- [Mandiant M-Trends 2020](https://www.mandiant.com/sites/default/files/2021-09/mtrends-2020.pdf)
- [Atomic Red Team - T1074.002](https://github.com/redcanaryco/atomic-red-team/tree/master/atomics/T1074.002)
- [MITRE ATT&CK - T1074.002](https://attack.mitre.org/techniques/T1074/002)
