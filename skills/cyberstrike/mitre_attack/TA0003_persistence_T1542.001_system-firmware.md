# stage: post_exploit
# category: mitre_attack


# T1542.001 System Firmware

> **Sub-technique of:** T1542

## High-Level Description

Adversaries may modify system firmware to persist on systems.The BIOS (Basic Input/Output System) and The Unified Extensible Firmware Interface (UEFI) or Extensible Firmware Interface (EFI) are examples of system firmware that operate as the software interface between the operating system and hardware of a computer.

System firmware like BIOS and (U)EFI underly the functionality of a computer and may be modified by an adversary to perform or assist in malicious activity. Capabilities exist to overwrite the system firmware, which may give sophisticated adversaries a means to install malicious firmware updates as a means of persistence on a system that may be difficult to detect.

## Kill Chain Phase

- Persistence (TA0003)
- Defense Evasion (TA0005)

**Platforms:** Windows, Network Devices

## What to Check

- [ ] Identify if System Firmware technique is applicable to target environment
- [ ] Check Windows systems for indicators of System Firmware
- [ ] Check Network Devices systems for indicators of System Firmware
- [ ] Verify mitigations are bypassed or absent (3 known mitigations)
- [ ] Assess detection coverage (1 detection strategies)

## How to Test

### Atomic Red Team Tests

The following tests are from [Atomic Red Team](https://github.com/redcanaryco/atomic-red-team) and provide actionable ways to test this technique:

### Atomic Test 1: UEFI Persistence via Wpbbin.exe File Creation

Creates Wpbbin.exe in %systemroot%. This technique can be used for UEFI-based pre-OS boot persistence mechanisms.

- https://grzegorztworek.medium.com/using-uefi-to-inject-executable-files-into-bitlocker-protected-drives-8ff4ca59c94c
- http://download.microsoft.com/download/8/a/2/8a2fb72d-9b96-4e2d-a559-4a27cf905a80/windows-platform-binary-table.docx
- https://github.com/tandasat/WPBT-Builder

**Supported Platforms:** windows
**Elevation Required:** Yes

```powershell
echo "Creating %systemroot%\wpbbin.exe"
New-Item -ItemType File -Path "$env:SystemRoot\System32\wpbbin.exe"
```

### Manual Testing

If Atomic Red Team tests are not applicable, manually verify the technique by:

1. **Identify Attack Surface**: Determine if the target environment is susceptible to System Firmware by examining the target platforms (Windows, Network Devices).

2. **Assess Existing Defenses**: Review whether mitigations for T1542.001 are in place. If defenses are absent or misconfigured, this technique may be exploitable.

3. **Execute Test**: Use tools and methods described in the MITRE ATT&CK page and external references below.

## Remediation Guide

### M1046 Boot Integrity

Check the integrity of the existing BIOS or EFI to determine if it is vulnerable to modification. Use Trusted Platform Module technology. Move system's root of trust to hardware to prevent tampering with the SPI flash memory. Technologies such as Intel Boot Guard can assist with this.

### M1051 Update Software

Patch the BIOS and EFI as necessary.

### M1026 Privileged Account Management

Prevent adversary access to privileged accounts or access necessary to perform this technique.

## Detection

### Detection Strategy for T1542.001 Pre-OS Boot: System Firmware

## Risk Assessment

| Finding                              | Severity | Impact      |
| ------------------------------------ | -------- | ----------- |
| System Firmware technique applicable | Low      | Persistence |

## CWE Categories

| CWE ID  | Title                         |
| ------- | ----------------------------- |
| CWE-276 | Incorrect Default Permissions |

## References

- [McAfee CHIPSEC Blog](https://securingtomorrow.mcafee.com/business/chipsec-support-vault-7-disclosure-scanning/)
- [MITRE Copernicus](http://www.mitre.org/capabilities/cybersecurity/overview/cybersecurity-blog/copernicus-question-your-assumptions-about)
- [Intel HackingTeam UEFI Rootkit](https://web.archive.org/web/20170313124421/http://www.intelsecurity.com/advanced-threat-research/content/data/HT-UEFI-rootkit.html)
- [Github CHIPSEC](https://github.com/chipsec/chipsec)
- [About UEFI](http://www.uefi.org/about)
- [MITRE Trustworthy Firmware Measurement](http://www.mitre.org/publications/project-stories/going-deep-into-the-bios-with-mitre-firmware-security-research)
- [Wikipedia UEFI](https://en.wikipedia.org/wiki/Unified_Extensible_Firmware_Interface)
- [Wikipedia BIOS](https://en.wikipedia.org/wiki/BIOS)
- [Atomic Red Team - T1542.001](https://github.com/redcanaryco/atomic-red-team/tree/master/atomics/T1542.001)
- [MITRE ATT&CK - T1542.001](https://attack.mitre.org/techniques/T1542/001)
