# stage: post_exploit
# category: mitre_attack_ics


# T0834 Native API

## High-Level Description

Adversaries may directly interact with the native OS application programming interface (API) to access system functions. Native APIs provide a controlled means of calling low-level OS services within the kernel, such as those involving hardware/devices, memory, and processes. These native APIs are leveraged by the OS during system boot (when other system components are not yet initialized) as well as carrying out tasks and requests during routine operations.

Functionality provided by native APIs are often also exposed to user-mode applications via interfaces and libraries. For example, functions such as memcpy and direct operations on memory registers can be used to modify user and system memory space.

## Kill Chain Phase

- Execution (TA0104)

**Platforms:** ICS

## What to Check

- [ ] Identify if Native API technique is applicable to target ICS environment
- [ ] Check ICS/SCADA systems for indicators of Native API
- [ ] Verify mitigations are bypassed or absent (1 known mitigations)
- [ ] Assess detection coverage (1 detection strategies)

## How to Test

### Identify Attack Surface

Determine if the target ICS/SCADA environment is susceptible to Native API by examining operational technology systems and network architecture.

### Assess Existing Defenses

Review whether mitigations for T0834 are in place. If defenses are absent or misconfigured, this technique may be exploitable in the ICS environment.

## Remediation Guide

### M0938 Execution Prevention

Minimize the exposure of API calls that allow the execution of code.

## Detection

### Detection of Native API

## Risk Assessment

| Finding                         | Severity | Impact    |
| ------------------------------- | -------- | --------- |
| Native API technique applicable | Low      | Execution |

## CWE Categories

| CWE ID | Title                                  |
| ------ | -------------------------------------- |
| CWE-94 | Improper Control of Generation of Code |

## References

- [The MITRE Corporation May 2017](https://attack.mitre.org/techniques/T1106/)
- [MITRE ATT&CK ICS - T0834](https://attack.mitre.org/techniques/T0834)
