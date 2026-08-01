# stage: post_exploit
# category: mitre_attack


# T1055.008 Ptrace System Calls

> **Sub-technique of:** T1055

## High-Level Description

Adversaries may inject malicious code into processes via ptrace (process trace) system calls in order to evade process-based defenses as well as possibly elevate privileges. Ptrace system call injection is a method of executing arbitrary code in the address space of a separate live process.

Ptrace system call injection involves attaching to and modifying a running process. The ptrace system call enables a debugging process to observe and control another process (and each individual thread), including changing memory and register values. Ptrace system call injection is commonly performed by writing arbitrary code into a running process (ex: <code>malloc</code>) then invoking that memory with <code>PTRACE_SETREGS</code> to set the register containing the next instruction to execute. Ptrace system call injection can also be done with <code>PTRACE_POKETEXT</code>/<code>PTRACE_POKEDATA</code>, which copy data to a specific address in the target processes’ memory (ex: the current address of the next instruction).

Ptrace system call injection may not be possible targeting processes that are non-child processes and/or have higher-privileges.

Running code in the context of another process may allow access to the process's memory, system/network resources, and possibly elevated privileges. Execution via ptrace system call injection may also evade detection from security products since the execution is masked under a legitimate process.

## Kill Chain Phase

- Defense Evasion (TA0005)
- Privilege Escalation (TA0004)

**Platforms:** Linux

## What to Check

- [ ] Identify if Ptrace System Calls technique is applicable to target environment
- [ ] Check Linux systems for indicators of Ptrace System Calls
- [ ] Verify mitigations are bypassed or absent (2 known mitigations)
- [ ] Assess detection coverage (1 detection strategies)

## How to Test

### Manual Testing

1. **Identify Attack Surface**: Determine if the target environment is susceptible to Ptrace System Calls by examining the target platforms (Linux).

2. **Assess Existing Defenses**: Review whether mitigations for T1055.008 are in place. If defenses are absent or misconfigured, this technique may be exploitable.

3. **Execute Test**: Use tools and methods described in the MITRE ATT&CK page and external references below.

> **Note**: No Atomic Red Team tests available for this technique. See [Atomic Red Team GitHub](https://github.com/redcanaryco/atomic-red-team) for updates.

## Remediation Guide

### M1040 Behavior Prevention on Endpoint

Some endpoint security solutions can be configured to block some types of process injection based on common sequences of behavior that occur during the injection process.

### M1026 Privileged Account Management

Utilize Yama (ex: /proc/sys/kernel/yama/ptrace_scope) to mitigate ptrace based process injection by restricting the use of ptrace to privileged users only. Other mitigation controls involve the deployment of security kernel modules that provide advanced access control and process restrictions such as SELinux, grsecurity, and AppArmor.

## Detection

### Detection Strategy for Ptrace-Based Process Injection on Linux

## Risk Assessment

| Finding                                  | Severity | Impact          |
| ---------------------------------------- | -------- | --------------- |
| Ptrace System Calls technique applicable | High     | Defense Evasion |

## CWE Categories

| CWE ID  | Title                        |
| ------- | ---------------------------- |
| CWE-693 | Protection Mechanism Failure |

## References

- [PTRACE man](http://man7.org/linux/man-pages/man2/ptrace.2.html)
- [Medium Ptrace JUL 2018](https://medium.com/@jain.sm/code-injection-in-running-process-using-ptrace-d3ea7191a4be)
- [BH Linux Inject](https://github.com/gaffe23/linux-inject/blob/master/slides_BHArsenal2015.pdf)
- [GNU Acct](https://www.gnu.org/software/acct/)
- [RHEL auditd](https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux/6/html/security_guide/chap-system_auditing)
- [Chokepoint preload rootkits](http://www.chokepoint.net/2014/02/detecting-userland-preload-rootkits.html)
- [Atomic Red Team - T1055.008](https://github.com/redcanaryco/atomic-red-team/tree/master/atomics/T1055.008)
- [MITRE ATT&CK - T1055.008](https://attack.mitre.org/techniques/T1055/008)
