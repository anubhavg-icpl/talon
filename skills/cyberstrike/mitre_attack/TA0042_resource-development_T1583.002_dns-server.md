# stage: post_exploit
# category: mitre_attack


# T1583.002 DNS Server

> **Sub-technique of:** T1583

## High-Level Description

Adversaries may set up their own Domain Name System (DNS) servers that can be used during targeting. During post-compromise activity, adversaries may utilize DNS traffic for various tasks, including for Command and Control (ex: Application Layer Protocol). Instead of hijacking existing DNS servers, adversaries may opt to configure and run their own DNS servers in support of operations.

By running their own DNS servers, adversaries can have more control over how they administer server-side DNS C2 traffic (DNS). With control over a DNS server, adversaries can configure DNS applications to provide conditional responses to malware and, generally, have more flexibility in the structure of the DNS-based C2 channel.

## Kill Chain Phase

- Resource Development (TA0042)

**Platforms:** PRE

## What to Check

- [ ] Identify if DNS Server technique is applicable to target environment
- [ ] Check PRE systems for indicators of DNS Server
- [ ] Verify mitigations are bypassed or absent (1 known mitigations)
- [ ] Assess detection coverage (1 detection strategies)

## How to Test

### Manual Testing

1. **Identify Attack Surface**: Determine if the target environment is susceptible to DNS Server by examining the target platforms (PRE).

2. **Assess Existing Defenses**: Review whether mitigations for T1583.002 are in place. If defenses are absent or misconfigured, this technique may be exploitable.

3. **Execute Test**: Use tools and methods described in the MITRE ATT&CK page and external references below.

> **Note**: No Atomic Red Team tests available for this technique. See [Atomic Red Team GitHub](https://github.com/redcanaryco/atomic-red-team) for updates.

## Remediation Guide

### M1056 Pre-compromise

This technique cannot be easily mitigated with preventive controls since it is based on behaviors performed outside of the scope of enterprise defenses and controls.

## Detection

### Detection of DNS Server

## Risk Assessment

| Finding                         | Severity | Impact               |
| ------------------------------- | -------- | -------------------- |
| DNS Server technique applicable | High     | Resource Development |

## CWE Categories

| CWE ID | Title                 |
| ------ | --------------------- |
| N/A    | No direct CWE mapping |

## References

- [Unit42 DNS Mar 2019](https://unit42.paloaltonetworks.com/dns-tunneling-how-dns-can-be-abused-by-malicious-actors/)
- [Atomic Red Team - T1583.002](https://github.com/redcanaryco/atomic-red-team/tree/master/atomics/T1583.002)
- [MITRE ATT&CK - T1583.002](https://attack.mitre.org/techniques/T1583/002)
