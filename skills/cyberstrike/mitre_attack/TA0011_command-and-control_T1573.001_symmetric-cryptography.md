# stage: post_exploit
# category: mitre_attack


# T1573.001 Symmetric Cryptography

> **Sub-technique of:** T1573

## High-Level Description

Adversaries may employ a known symmetric encryption algorithm to conceal command and control traffic rather than relying on any inherent protections provided by a communication protocol. Symmetric encryption algorithms use the same key for plaintext encryption and ciphertext decryption. Common symmetric encryption algorithms include AES, DES, 3DES, Blowfish, and RC4.

## Kill Chain Phase

- Command and Control (TA0011)

**Platforms:** ESXi, Linux, macOS, Network Devices, Windows

## What to Check

- [ ] Identify if Symmetric Cryptography technique is applicable to target environment
- [ ] Check ESXi systems for indicators of Symmetric Cryptography
- [ ] Check Linux systems for indicators of Symmetric Cryptography
- [ ] Check macOS systems for indicators of Symmetric Cryptography
- [ ] Verify mitigations are bypassed or absent (1 known mitigations)
- [ ] Assess detection coverage (1 detection strategies)

## How to Test

### Manual Testing

1. **Identify Attack Surface**: Determine if the target environment is susceptible to Symmetric Cryptography by examining the target platforms (ESXi, Linux, macOS).

2. **Assess Existing Defenses**: Review whether mitigations for T1573.001 are in place. If defenses are absent or misconfigured, this technique may be exploitable.

3. **Execute Test**: Use tools and methods described in the MITRE ATT&CK page and external references below.

> **Note**: No Atomic Red Team tests available for this technique. See [Atomic Red Team GitHub](https://github.com/redcanaryco/atomic-red-team) for updates.

## Remediation Guide

### M1031 Network Intrusion Prevention

Network intrusion detection and prevention systems that use network signatures to identify traffic for specific adversary malware can be used to mitigate activity at the network level.

## Detection

### Detection Strategy for Encrypted Channel via Symmetric Cryptography across OS Platforms

## Risk Assessment

| Finding                                     | Severity | Impact              |
| ------------------------------------------- | -------- | ------------------- |
| Symmetric Cryptography technique applicable | Low      | Command And Control |

## CWE Categories

| CWE ID  | Title                              |
| ------- | ---------------------------------- |
| CWE-300 | Channel Accessible by Non-Endpoint |

## References

- [University of Birmingham C2](https://arxiv.org/ftp/arxiv/papers/1408/1408.1136.pdf)
- [Atomic Red Team - T1573.001](https://github.com/redcanaryco/atomic-red-team/tree/master/atomics/T1573.001)
- [MITRE ATT&CK - T1573.001](https://attack.mitre.org/techniques/T1573/001)
