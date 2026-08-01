# stage: post_exploit
# category: mitre_attack


# T1546.016 Installer Packages

> **Sub-technique of:** T1546

## High-Level Description

Adversaries may establish persistence and elevate privileges by using an installer to trigger the execution of malicious content. Installer packages are OS specific and contain the resources an operating system needs to install applications on a system. Installer packages can include scripts that run prior to installation as well as after installation is complete. Installer scripts may inherit elevated permissions when executed. Developers often use these scripts to prepare the environment for installation, check requirements, download dependencies, and remove files after installation.

Using legitimate applications, adversaries have distributed applications with modified installer scripts to execute malicious content. When a user installs the application, they may be required to grant administrative permissions to allow the installation. At the end of the installation process of the legitimate application, content such as macOS `postinstall` scripts can be executed with the inherited elevated permissions. Adversaries can use these scripts to execute a malicious executable or install other malicious components (such as a Launch Daemon) with the elevated permissions.

Depending on the distribution, Linux versions of package installer scripts are sometimes called maintainer scripts or post installation scripts. These scripts can include `preinst`, `postinst`, `prerm`, `postrm` scripts and run as root when executed.

For Windows, the Microsoft Installer services uses `.msi` files to manage the installing, updating, and uninstalling of applications. These installation routines may also include instructions to perform additional actions that may be abused by adversaries.

## Kill Chain Phase

- Privilege Escalation (TA0004)
- Persistence (TA0003)

**Platforms:** Linux, Windows, macOS

## What to Check

- [ ] Identify if Installer Packages technique is applicable to target environment
- [ ] Check Linux systems for indicators of Installer Packages
- [ ] Check Windows systems for indicators of Installer Packages
- [ ] Check macOS systems for indicators of Installer Packages
- [ ] Assess detection coverage (1 detection strategies)

## How to Test

### Manual Testing

1. **Identify Attack Surface**: Determine if the target environment is susceptible to Installer Packages by examining the target platforms (Linux, Windows, macOS).

2. **Assess Existing Defenses**: Review whether mitigations for T1546.016 are in place. If defenses are absent or misconfigured, this technique may be exploitable.

3. **Execute Test**: Use tools and methods described in the MITRE ATT&CK page and external references below.

> **Note**: No Atomic Red Team tests available for this technique. See [Atomic Red Team GitHub](https://github.com/redcanaryco/atomic-red-team) for updates.

## Remediation Guide

No specific mitigations documented for this technique.

## Detection

### Detection Strategy for T1546.016 - Event Triggered Execution via Installer Packages

## Risk Assessment

| Finding                                 | Severity | Impact               |
| --------------------------------------- | -------- | -------------------- |
| Installer Packages technique applicable | High     | Privilege Escalation |

## CWE Categories

| CWE ID  | Title                         |
| ------- | ----------------------------- |
| CWE-269 | Improper Privilege Management |

## References

- [Application Bundle Manipulation Brandon Dalton](https://redcanary.com/blog/mac-application-bundles/)
- [Debian Manual Maintainer Scripts](https://www.debian.org/doc/debian-policy/ch-maintainerscripts.html#s-mscriptsinstact)
- [Windows AppleJeus GReAT](https://securelist.com/operation-applejeus/87553/)
- [Microsoft Installation Procedures](https://learn.microsoft.com/windows/win32/msi/installation-procedure-tables-group)
- [wardle evilquest parti](https://objective-see.com/blog/blog_0x59.html)
- [Installer Package Scripting Rich Trouton](https://cpb-us-e1.wpmucdn.com/sites.psu.edu/dist/4/24696/files/2019/07/psumac2019-345-Installer-Package-Scripting-Making-your-deployments-easier-one-at-a-time.pdf)
- [Atomic Red Team - T1546.016](https://github.com/redcanaryco/atomic-red-team/tree/master/atomics/T1546.016)
- [MITRE ATT&CK - T1546.016](https://attack.mitre.org/techniques/T1546/016)
