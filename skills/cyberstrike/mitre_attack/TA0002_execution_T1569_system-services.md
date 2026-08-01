# stage: post_exploit
# category: mitre_attack


# T1569 System Services

## High-Level Description

Adversaries may abuse system services or daemons to execute commands or programs. Adversaries can execute malicious content by interacting with or creating services either locally or remotely. Many services are set to run at boot, which can aid in achieving persistence (Create or Modify System Process), but adversaries can also abuse services for one-time or temporary execution.

## Kill Chain Phase

- Execution (TA0002)

**Platforms:** Windows, macOS, Linux

## What to Check

- [ ] Identify if System Services technique is applicable to target environment
- [ ] Check Windows systems for indicators of System Services
- [ ] Check macOS systems for indicators of System Services
- [ ] Check Linux systems for indicators of System Services
- [ ] Verify mitigations are bypassed or absent (4 known mitigations)
- [ ] Assess detection coverage (1 detection strategies)

## How to Test

### Manual Testing

1. **Identify Attack Surface**: Determine if the target environment is susceptible to System Services by examining the target platforms (Windows, macOS, Linux).

2. **Assess Existing Defenses**: Review whether mitigations for T1569 are in place. If defenses are absent or misconfigured, this technique may be exploitable.

3. **Execute Test**: Use tools and methods described in the MITRE ATT&CK page and external references below.

> **Note**: No Atomic Red Team tests available for this technique. See [Atomic Red Team GitHub](https://github.com/redcanaryco/atomic-red-team) for updates.

## Remediation Guide

### M1026 Privileged Account Management

Ensure that permissions disallow services that run at a higher permissions level from being created or interacted with by a user with a lower permission level.

### M1018 User Account Management

Prevent users from installing their own launch agents or launch daemons.

### M1040 Behavior Prevention on Endpoint

On Windows 10, enable Attack Surface Reduction (ASR) rules to block processes created by PsExec from running.

### M1022 Restrict File and Directory Permissions

Ensure that high permission level service binaries cannot be replaced or modified by users with a lower permission level.

## Detection

### Detection Strategy for System Services across OS platforms.

## Risk Assessment

| Finding                              | Severity | Impact    |
| ------------------------------------ | -------- | --------- |
| System Services technique applicable | Low      | Execution |

## CWE Categories

| CWE ID | Title                                  |
| ------ | -------------------------------------- |
| CWE-94 | Improper Control of Generation of Code |

## References

- [Atomic Red Team - T1569](https://github.com/redcanaryco/atomic-red-team/tree/master/atomics/T1569)
- [MITRE ATT&CK - T1569](https://attack.mitre.org/techniques/T1569)
