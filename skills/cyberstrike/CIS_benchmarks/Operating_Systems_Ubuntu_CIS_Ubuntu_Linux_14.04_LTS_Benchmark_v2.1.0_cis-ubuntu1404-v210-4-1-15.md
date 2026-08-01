# stage: report
# category: CIS_benchmarks


# 4.1.15 Ensure changes to system administration scope (sudoers) is collected (Scored)

## Profile Applicability

- Level 2 - Server
- Level 2 - Workstation

## Description

Monitor scope changes for system administrations. If the system has been properly configured to force system administrators to log in as themselves first and then use the `sudo` command to execute privileged commands, it is possible to monitor changes in scope. The file `/etc/sudoers` will be written to when the file or its attributes have changed. The audit records will be tagged with the identifier "scope."

## Rationale

Changes in the `/etc/sudoers` file can indicate that an unauthorized change has been made to scope of system administrator activity.

## Audit Procedure

Run the following commands:

```bash
grep scope /etc/audit/audit.rules
auditctl -l | grep scope
```

## Expected Result

Verify output of both matches:

```
-w /etc/sudoers -p wa -k scope
-w /etc/sudoers.d/ -p wa -k scope
```

## Remediation

Add the following line to the `/etc/audit/audit.rules` file:

```bash
-w /etc/sudoers -p wa -k scope
-w /etc/sudoers.d/ -p wa -k scope
```

**Notes:** Reloading the auditd config to set active settings may require a system reboot.

## Default Value

Not configured by default.

## References

1. CIS Controls v6.1 - 5.4 Log Administrative User Addition And Removal

## Profile

- Level 2 - Server
- Level 2 - Workstation
