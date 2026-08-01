# stage: post_exploit
# category: mitre_attack


# T1011.001 Exfiltration Over Bluetooth

> **Sub-technique of:** T1011

## High-Level Description

Adversaries may attempt to exfiltrate data over Bluetooth rather than the command and control channel. If the command and control network is a wired Internet connection, an adversary may opt to exfiltrate data using a Bluetooth communication channel.

Adversaries may choose to do this if they have sufficient access and proximity. Bluetooth connections might not be secured or defended as well as the primary Internet-connected channel because it is not routed through the same enterprise network.

## Kill Chain Phase

- Exfiltration (TA0010)

**Platforms:** Linux, macOS, Windows

## What to Check

- [ ] Identify if Exfiltration Over Bluetooth technique is applicable to target environment
- [ ] Check Linux systems for indicators of Exfiltration Over Bluetooth
- [ ] Check macOS systems for indicators of Exfiltration Over Bluetooth
- [ ] Check Windows systems for indicators of Exfiltration Over Bluetooth
- [ ] Verify mitigations are bypassed or absent (2 known mitigations)
- [ ] Assess detection coverage (1 detection strategies)

## How to Test

### Manual Testing

1. **Identify Attack Surface**: Determine if the target environment is susceptible to Exfiltration Over Bluetooth by examining the target platforms (Linux, macOS, Windows).

2. **Assess Existing Defenses**: Review whether mitigations for T1011.001 are in place. If defenses are absent or misconfigured, this technique may be exploitable.

3. **Execute Test**: Use tools and methods described in the MITRE ATT&CK page and external references below.

> **Note**: No Atomic Red Team tests available for this technique. See [Atomic Red Team GitHub](https://github.com/redcanaryco/atomic-red-team) for updates.

## Remediation Guide

### M1042 Disable or Remove Feature or Program

Disable Bluetooth in local computer security settings or by group policy if it is not needed within an environment.

### M1028 Operating System Configuration

Prevent the creation of new network adapters where possible.

## Detection

### Detection of Bluetooth-Based Data Exfiltration

## Risk Assessment

| Finding                                          | Severity | Impact       |
| ------------------------------------------------ | -------- | ------------ |
| Exfiltration Over Bluetooth technique applicable | Low      | Exfiltration |

## CWE Categories

| CWE ID  | Title                             |
| ------- | --------------------------------- |
| CWE-200 | Exposure of Sensitive Information |

## References

- [Atomic Red Team - T1011.001](https://github.com/redcanaryco/atomic-red-team/tree/master/atomics/T1011.001)
- [MITRE ATT&CK - T1011.001](https://attack.mitre.org/techniques/T1011/001)
