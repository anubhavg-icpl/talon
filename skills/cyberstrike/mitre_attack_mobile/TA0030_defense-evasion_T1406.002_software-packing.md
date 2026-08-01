# stage: post_exploit
# category: mitre_attack_mobile


# T1406.002 Software Packing

> **Sub-technique of:** T1406

## High-Level Description

Adversaries may perform software packing to conceal their code. Software packing is a method of compressing or encrypting an executable. Packing an executable changes the file signature in an attempt to avoid signature-based detection. Most decompression techniques decompress the executable code in memory.

Utilities used to perform software packing are called packers. An example packer is FTT. A more comprehensive list of known packers is available, but adversaries may create their own packing techniques that do not leave the same artifacts as well-known packers to evade defenses.

## Kill Chain Phase

- Defense Evasion (TA0030)

**Platforms:** iOS, Android

## What to Check

- [ ] Identify if Software Packing technique is applicable to target mobile environment
- [ ] Check iOS devices for indicators of Software Packing
- [ ] Check Android devices for indicators of Software Packing
- [ ] Assess detection coverage (1 detection strategies)

## How to Test

### Identify Attack Surface

Determine if the target mobile environment is susceptible to Software Packing by examining the target platforms (iOS, Android).

### Assess Existing Defenses

Review whether mitigations for T1406.002 are in place. If defenses are absent or misconfigured, this technique may be exploitable.

## Remediation Guide

No specific mitigations documented for this technique.

## Detection

### Detection of Software Packing

## Risk Assessment

| Finding                               | Severity | Impact          |
| ------------------------------------- | -------- | --------------- |
| Software Packing technique applicable | Low      | Defense Evasion |

## CWE Categories

| CWE ID  | Title                        |
| ------- | ---------------------------- |
| CWE-693 | Protection Mechanism Failure |

## References

- [MITRE ATT&CK Mobile - T1406.002](https://attack.mitre.org/techniques/T1406/002)
