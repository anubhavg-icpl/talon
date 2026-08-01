# stage: post_exploit
# category: mitre_attack_mobile


# T1521.001 Symmetric Cryptography

> **Sub-technique of:** T1521

## High-Level Description

Adversaries may employ a known symmetric encryption algorithm to conceal command and control traffic, rather than relying on any inherent protections provided by a communication protocol. Symmetric encryption algorithms use the same key for plaintext encryption and ciphertext decryption. Common symmetric encryption algorithms include AES, Blowfish, and RC4.

## Kill Chain Phase

- Command and Control (TA0037)

**Platforms:** Android, iOS

## What to Check

- [ ] Identify if Symmetric Cryptography technique is applicable to target mobile environment
- [ ] Check Android devices for indicators of Symmetric Cryptography
- [ ] Check iOS devices for indicators of Symmetric Cryptography
- [ ] Assess detection coverage (1 detection strategies)

## How to Test

### Identify Attack Surface

Determine if the target mobile environment is susceptible to Symmetric Cryptography by examining the target platforms (Android, iOS).

### Assess Existing Defenses

Review whether mitigations for T1521.001 are in place. If defenses are absent or misconfigured, this technique may be exploitable.

## Remediation Guide

No specific mitigations documented for this technique.

## Detection

### Detection of Symmetric Cryptography

## Risk Assessment

| Finding                                     | Severity | Impact              |
| ------------------------------------------- | -------- | ------------------- |
| Symmetric Cryptography technique applicable | Low      | Command And Control |

## CWE Categories

| CWE ID  | Title                              |
| ------- | ---------------------------------- |
| CWE-300 | Channel Accessible by Non-Endpoint |

## References

- [MITRE ATT&CK Mobile - T1521.001](https://attack.mitre.org/techniques/T1521/001)
