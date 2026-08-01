# stage: post_exploit
# category: mitre_attack


# T1606.001 Web Cookies

> **Sub-technique of:** T1606

## High-Level Description

Adversaries may forge web cookies that can be used to gain access to web applications or Internet services. Web applications and services (hosted in cloud SaaS environments or on-premise servers) often use session cookies to authenticate and authorize user access.

Adversaries may generate these cookies in order to gain access to web resources. This differs from Steal Web Session Cookie and other similar behaviors in that the cookies are new and forged by the adversary, rather than stolen or intercepted from legitimate users. Most common web applications have standardized and documented cookie values that can be generated using provided tools or interfaces. The generation of web cookies often requires secret values, such as passwords, Private Keys, or other cryptographic seed values.

Once forged, adversaries may use these web cookies to access resources (Web Session Cookie), which may bypass multi-factor and other authentication protection mechanisms.

## Kill Chain Phase

- Credential Access (TA0006)

**Platforms:** Linux, macOS, Windows, SaaS, IaaS

## What to Check

- [ ] Identify if Web Cookies technique is applicable to target environment
- [ ] Check Linux systems for indicators of Web Cookies
- [ ] Check macOS systems for indicators of Web Cookies
- [ ] Check Windows systems for indicators of Web Cookies
- [ ] Verify mitigations are bypassed or absent (2 known mitigations)
- [ ] Assess detection coverage (1 detection strategies)

## How to Test

### Manual Testing

1. **Identify Attack Surface**: Determine if the target environment is susceptible to Web Cookies by examining the target platforms (Linux, macOS, Windows).

2. **Assess Existing Defenses**: Review whether mitigations for T1606.001 are in place. If defenses are absent or misconfigured, this technique may be exploitable.

3. **Execute Test**: Use tools and methods described in the MITRE ATT&CK page and external references below.

> **Note**: No Atomic Red Team tests available for this technique. See [Atomic Red Team GitHub](https://github.com/redcanaryco/atomic-red-team) for updates.

## Remediation Guide

### M1047 Audit

Administrators should perform an audit of all access lists and the permissions they have been granted to access web applications and services. This should be done extensively on all resources in order to establish a baseline, followed up on with periodic audits of new or updated resources. Suspicious accounts/credentials should be investigated and removed.

### M1054 Software Configuration

Configure browsers/applications to regularly delete persistent web cookies.

## Detection

### Detection Strategy for Forged Web Cookies

## Risk Assessment

| Finding                          | Severity | Impact            |
| -------------------------------- | -------- | ----------------- |
| Web Cookies technique applicable | Low      | Credential Access |

## CWE Categories

| CWE ID  | Title                                |
| ------- | ------------------------------------ |
| CWE-522 | Insufficiently Protected Credentials |

## References

- [Volexity SolarWinds](https://www.volexity.com/blog/2020/12/14/dark-halo-leverages-solarwinds-compromise-to-breach-organizations/)
- [Unit 42 Mac Crypto Cookies January 2019](https://unit42.paloaltonetworks.com/mac-malware-steals-cryptocurrency-exchanges-cookies/)
- [Pass The Cookie](https://wunderwuzzi23.github.io/blog/passthecookie.html)
- [Atomic Red Team - T1606.001](https://github.com/redcanaryco/atomic-red-team/tree/master/atomics/T1606.001)
- [MITRE ATT&CK - T1606.001](https://attack.mitre.org/techniques/T1606/001)
