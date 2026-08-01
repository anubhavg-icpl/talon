# stage: post_exploit
# category: mitre_attack


# T1036.011 Overwrite Process Arguments

> **Sub-technique of:** T1036

## High-Level Description

Adversaries may modify a process's in-memory arguments to change its name in order to appear as a legitimate or benign process. On Linux, the operating system stores command-line arguments in the process’s stack and passes them to the `main()` function as the `argv` array. The first element, `argv[0]`, typically contains the process name or path - by default, the command used to actually start the process (e.g., `cat /etc/passwd`). By default, the Linux `/proc` filesystem uses this value to represent the process name. The `/proc/<PID>/cmdline` file reflects the contents of this memory, and tools like `ps` use it to display process information. Since arguments are stored in user-space memory at launch, this modification can be performed without elevated privileges.

During runtime, adversaries can erase the memory used by all command-line arguments for a process, overwriting each argument string with null bytes. This removes evidence of how the process was originally launched. They can then write a spoofed string into the memory region previously occupied by `argv[0]` to mimic a benign command, such as `cat resolv.conf`. The new command-line string is reflected in `/proc/<PID>/cmdline` and displayed by tools like `ps`.

## Kill Chain Phase

- Defense Evasion (TA0005)

**Platforms:** Linux

## What to Check

- [ ] Identify if Overwrite Process Arguments technique is applicable to target environment
- [ ] Check Linux systems for indicators of Overwrite Process Arguments
- [ ] Assess detection coverage (1 detection strategies)

## How to Test

### Manual Testing

1. **Identify Attack Surface**: Determine if the target environment is susceptible to Overwrite Process Arguments by examining the target platforms (Linux).

2. **Assess Existing Defenses**: Review whether mitigations for T1036.011 are in place. If defenses are absent or misconfigured, this technique may be exploitable.

3. **Execute Test**: Use tools and methods described in the MITRE ATT&CK page and external references below.

> **Note**: No Atomic Red Team tests available for this technique. See [Atomic Red Team GitHub](https://github.com/redcanaryco/atomic-red-team) for updates.

## Remediation Guide

No specific mitigations documented for this technique.

## Detection

### Detection Strategy for Overwritten Process Arguments Masquerading

## Risk Assessment

| Finding                                          | Severity | Impact          |
| ------------------------------------------------ | -------- | --------------- |
| Overwrite Process Arguments technique applicable | High     | Defense Evasion |

## CWE Categories

| CWE ID  | Title                        |
| ------- | ---------------------------- |
| CWE-693 | Protection Mechanism Failure |

## References

- [Microsoft XorDdos Linux Stealth 2022](https://www.microsoft.com/en-us/security/blog/2022/05/19/rise-in-xorddos-a-deeper-look-at-the-stealthy-ddos-malware-targeting-linux-devices/)
- [Sandfly BPFDoor 2022](https://sandflysecurity.com/blog/bpfdoor-an-evasive-linux-backdoor-technical-analysis/)
- [Atomic Red Team - T1036.011](https://github.com/redcanaryco/atomic-red-team/tree/master/atomics/T1036.011)
- [MITRE ATT&CK - T1036.011](https://attack.mitre.org/techniques/T1036/011)
