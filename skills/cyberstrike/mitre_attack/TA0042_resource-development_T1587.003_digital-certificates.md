# stage: post_exploit
# category: mitre_attack


# T1587.003 Digital Certificates

> **Sub-technique of:** T1587

## High-Level Description

Adversaries may create self-signed SSL/TLS certificates that can be used during targeting. SSL/TLS certificates are designed to instill trust. They include information about the key, information about its owner's identity, and the digital signature of an entity that has verified the certificate's contents are correct. If the signature is valid, and the person examining the certificate trusts the signer, then they know they can use that key to communicate with its owner. In the case of self-signing, digital certificates will lack the element of trust associated with the signature of a third-party certificate authority (CA).

Adversaries may create self-signed SSL/TLS certificates that can be used to further their operations, such as encrypting C2 traffic (ex: Asymmetric Cryptography with Web Protocols) or even enabling Adversary-in-the-Middle if added to the root of trust (i.e. Install Root Certificate).

After creating a digital certificate, an adversary may then install that certificate (see Install Digital Certificate) on infrastructure under their control.

## Kill Chain Phase

- Resource Development (TA0042)

**Platforms:** PRE

## What to Check

- [ ] Identify if Digital Certificates technique is applicable to target environment
- [ ] Check PRE systems for indicators of Digital Certificates
- [ ] Verify mitigations are bypassed or absent (1 known mitigations)
- [ ] Assess detection coverage (1 detection strategies)

## How to Test

### Manual Testing

1. **Identify Attack Surface**: Determine if the target environment is susceptible to Digital Certificates by examining the target platforms (PRE).

2. **Assess Existing Defenses**: Review whether mitigations for T1587.003 are in place. If defenses are absent or misconfigured, this technique may be exploitable.

3. **Execute Test**: Use tools and methods described in the MITRE ATT&CK page and external references below.

> **Note**: No Atomic Red Team tests available for this technique. See [Atomic Red Team GitHub](https://github.com/redcanaryco/atomic-red-team) for updates.

## Remediation Guide

### M1056 Pre-compromise

This technique cannot be easily mitigated with preventive controls since it is based on behaviors performed outside of the scope of enterprise defenses and controls.

## Detection

### Detection of Digital Certificates

## Risk Assessment

| Finding                                   | Severity | Impact               |
| ----------------------------------------- | -------- | -------------------- |
| Digital Certificates technique applicable | Low      | Resource Development |

## CWE Categories

| CWE ID | Title                 |
| ------ | --------------------- |
| N/A    | No direct CWE mapping |

## References

- [Splunk Kovar Certificates 2017](https://www.splunk.com/en_us/blog/security/tall-tales-of-hunting-with-tls-ssl-certificates.html)
- [Atomic Red Team - T1587.003](https://github.com/redcanaryco/atomic-red-team/tree/master/atomics/T1587.003)
- [MITRE ATT&CK - T1587.003](https://attack.mitre.org/techniques/T1587/003)
