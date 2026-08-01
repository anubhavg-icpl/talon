# stage: post_exploit
# category: mitre_attack


# T1564.011 Ignore Process Interrupts

> **Sub-technique of:** T1564

## High-Level Description

Adversaries may evade defensive mechanisms by executing commands that hide from process interrupt signals. Many operating systems use signals to deliver messages to control process behavior. Command interpreters often include specific commands/flags that ignore errors and other hangups, such as when the user of the active session logs off. These interrupt signals may also be used by defensive tools and/or analysts to pause or terminate specified running processes.

Adversaries may invoke processes using `nohup`, PowerShell `-ErrorAction SilentlyContinue`, or similar commands that may be immune to hangups. This may enable malicious commands and malware to continue execution through system events that would otherwise terminate its execution, such as users logging off or the termination of its C2 network connection.

Hiding from process interrupt signals may allow malware to continue execution, but unlike Trap this does not establish Persistence since the process will not be re-invoked once actually terminated.

## Kill Chain Phase

- Defense Evasion (TA0005)

**Platforms:** Linux, macOS, Windows

## What to Check

- [ ] Identify if Ignore Process Interrupts technique is applicable to target environment
- [ ] Check Linux systems for indicators of Ignore Process Interrupts
- [ ] Check macOS systems for indicators of Ignore Process Interrupts
- [ ] Check Windows systems for indicators of Ignore Process Interrupts
- [ ] Assess detection coverage (1 detection strategies)

## How to Test

### Manual Testing

1. **Identify Attack Surface**: Determine if the target environment is susceptible to Ignore Process Interrupts by examining the target platforms (Linux, macOS, Windows).

2. **Assess Existing Defenses**: Review whether mitigations for T1564.011 are in place. If defenses are absent or misconfigured, this technique may be exploitable.

3. **Execute Test**: Use tools and methods described in the MITRE ATT&CK page and external references below.

> **Note**: No Atomic Red Team tests available for this technique. See [Atomic Red Team GitHub](https://github.com/redcanaryco/atomic-red-team) for updates.

## Remediation Guide

No specific mitigations documented for this technique.

## Detection

### Detection Strategy for Ignore Process Interrupts

## Risk Assessment

| Finding                                        | Severity | Impact          |
| ---------------------------------------------- | -------- | --------------- |
| Ignore Process Interrupts technique applicable | Low      | Defense Evasion |

## CWE Categories

| CWE ID  | Title                        |
| ------- | ---------------------------- |
| CWE-693 | Protection Mechanism Failure |

## References

- [Linux Signal Man](https://man7.org/linux/man-pages/man7/signal.7.html)
- [nohup Linux Man](https://linux.die.net/man/1/nohup)
- [Microsoft PowerShell SilentlyContinue](https://learn.microsoft.com/powershell/module/microsoft.powershell.core/about/about_preference_variables?view=powershell-7.3#debugpreference)
- [Atomic Red Team - T1564.011](https://github.com/redcanaryco/atomic-red-team/tree/master/atomics/T1564.011)
- [MITRE ATT&CK - T1564.011](https://attack.mitre.org/techniques/T1564/011)
