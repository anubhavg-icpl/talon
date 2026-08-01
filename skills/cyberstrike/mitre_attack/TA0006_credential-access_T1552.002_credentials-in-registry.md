# stage: post_exploit
# category: mitre_attack


# T1552.002 Credentials in Registry

> **Sub-technique of:** T1552

## High-Level Description

Adversaries may search the Registry on compromised systems for insecurely stored credentials. The Windows Registry stores configuration information that can be used by the system or other programs. Adversaries may query the Registry looking for credentials and passwords that have been stored for use by other programs or services. Sometimes these credentials are used for automatic logons.

Example commands to find Registry keys related to password information:

- Local Machine Hive: <code>reg query HKLM /f password /t REG_SZ /s</code>
- Current User Hive: <code>reg query HKCU /f password /t REG_SZ /s</code>

## Kill Chain Phase

- Credential Access (TA0006)

**Platforms:** Windows

## What to Check

- [ ] Identify if Credentials in Registry technique is applicable to target environment
- [ ] Check Windows systems for indicators of Credentials in Registry
- [ ] Verify mitigations are bypassed or absent (3 known mitigations)
- [ ] Assess detection coverage (1 detection strategies)

## How to Test

### Atomic Red Team Tests

The following tests are from [Atomic Red Team](https://github.com/redcanaryco/atomic-red-team) and provide actionable ways to test this technique:

### Atomic Test 1: Enumeration for Credentials in Registry

Queries to enumerate for credentials in the Registry. Upon execution, any registry key containing the word "password" will be displayed.

**Supported Platforms:** windows

```cmd
reg query HKLM /f password /t REG_SZ /s
reg query HKCU /f password /t REG_SZ /s
```

### Atomic Test 2: Enumeration for PuTTY Credentials in Registry

Queries to enumerate for PuTTY credentials in the Registry. PuTTY must be installed for this test to work. If any registry
entries are found, they will be displayed.

**Supported Platforms:** windows

```cmd
reg query HKCU\Software\SimonTatham\PuTTY\Sessions /t REG_SZ /s
```

### Manual Testing

If Atomic Red Team tests are not applicable, manually verify the technique by:

1. **Identify Attack Surface**: Determine if the target environment is susceptible to Credentials in Registry by examining the target platforms (Windows).

2. **Assess Existing Defenses**: Review whether mitigations for T1552.002 are in place. If defenses are absent or misconfigured, this technique may be exploitable.

3. **Execute Test**: Use tools and methods described in the MITRE ATT&CK page and external references below.

## Remediation Guide

### M1027 Password Policies

Do not store credentials within the Registry.

### M1026 Privileged Account Management

If it is necessary that software must store credentials in the Registry, then ensure the associated accounts have limited permissions so they cannot be abused if obtained by an adversary.

### M1047 Audit

Proactively search for credentials within the Registry and attempt to remediate the risk.

## Detection

### Detect Credential Discovery via Windows Registry Enumeration

## Risk Assessment

| Finding                                      | Severity | Impact            |
| -------------------------------------------- | -------- | ----------------- |
| Credentials in Registry technique applicable | High     | Credential Access |

## CWE Categories

| CWE ID  | Title                                |
| ------- | ------------------------------------ |
| CWE-522 | Insufficiently Protected Credentials |

## References

- [Pentestlab Stored Credentials](https://pentestlab.blog/2017/04/19/stored-credentials/)
- [Atomic Red Team - T1552.002](https://github.com/redcanaryco/atomic-red-team/tree/master/atomics/T1552.002)
- [MITRE ATT&CK - T1552.002](https://attack.mitre.org/techniques/T1552/002)
