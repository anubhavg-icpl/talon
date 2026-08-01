# stage: post_exploit
# category: mitre_attack


# T1564.009 Resource Forking

> **Sub-technique of:** T1564

## High-Level Description

Adversaries may abuse resource forks to hide malicious code or executables to evade detection and bypass security applications. A resource fork provides applications a structured way to store resources such as thumbnail images, menu definitions, icons, dialog boxes, and code. Usage of a resource fork is identifiable when displaying a file’s extended attributes, using <code>ls -l@</code> or <code>xattr -l</code> commands. Resource forks have been deprecated and replaced with the application bundle structure. Non-localized resources are placed at the top level directory of an application bundle, while localized resources are placed in the <code>/Resources</code> folder.

Adversaries can use resource forks to hide malicious data that may otherwise be stored directly in files. Adversaries can execute content with an attached resource fork, at a specified offset, that is moved to an executable location then invoked. Resource fork content may also be obfuscated/encrypted until execution.

## Kill Chain Phase

- Defense Evasion (TA0005)

**Platforms:** macOS

## What to Check

- [ ] Identify if Resource Forking technique is applicable to target environment
- [ ] Check macOS systems for indicators of Resource Forking
- [ ] Verify mitigations are bypassed or absent (1 known mitigations)
- [ ] Assess detection coverage (1 detection strategies)

## How to Test

### Manual Testing

1. **Identify Attack Surface**: Determine if the target environment is susceptible to Resource Forking by examining the target platforms (macOS).

2. **Assess Existing Defenses**: Review whether mitigations for T1564.009 are in place. If defenses are absent or misconfigured, this technique may be exploitable.

3. **Execute Test**: Use tools and methods described in the MITRE ATT&CK page and external references below.

> **Note**: No Atomic Red Team tests available for this technique. See [Atomic Red Team GitHub](https://github.com/redcanaryco/atomic-red-team) for updates.

## Remediation Guide

### M1013 Application Developer Guidance

Configure applications to use the application bundle structure which leverages the <code>/Resources</code> folder location.

## Detection

### Detection Strategy for Resource Forking on macOS

## Risk Assessment

| Finding                               | Severity | Impact          |
| ------------------------------------- | -------- | --------------- |
| Resource Forking technique applicable | Medium   | Defense Evasion |

## CWE Categories

| CWE ID  | Title                        |
| ------- | ---------------------------- |
| CWE-693 | Protection Mechanism Failure |

## References

- [tau bundlore erika noerenberg 2020](https://blogs.vmware.com/security/2020/06/tau-threat-analysis-bundlore-macos-mm-install-macos.html)
- [Resource and Data Forks](https://flylib.com/books/en/4.395.1.192/1/)
- [ELC Extended Attributes](https://eclecticlight.co/2020/10/24/theres-more-to-files-than-data-extended-attributes/)
- [sentinellabs resource named fork 2020](https://www.sentinelone.com/labs/resourceful-macos-malware-hides-in-named-fork/)
- [macOS Hierarchical File System Overview](http://tenon.com/products/codebuilder/User_Guide/6_File_Systems.html#anchor520553)
- [Atomic Red Team - T1564.009](https://github.com/redcanaryco/atomic-red-team/tree/master/atomics/T1564.009)
- [MITRE ATT&CK - T1564.009](https://attack.mitre.org/techniques/T1564/009)
