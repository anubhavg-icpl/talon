# stage: report
# category: CIS_benchmarks


# 4.1.1.2 Ensure system is disabled when audit logs are full (Scored)

## Profile Applicability

- Level 2 - Server
- Level 2 - Workstation

## Description

The `auditd` daemon can be configured to halt the system when the audit logs are full.

## Rationale

In high security contexts, the risk of detecting unauthorized access or nonrepudiation exceeds the benefit of the system's availability.

## Audit Procedure

Run the following commands and verify output matches:

```bash
grep space_left_action /etc/audit/auditd.conf
grep action_mail_acct /etc/audit/auditd.conf
grep admin_space_left_action /etc/audit/auditd.conf
```

## Expected Result

```
space_left_action = email
action_mail_acct = root
admin_space_left_action = halt
```

## Remediation

Set the following parameters in `/etc/audit/auditd.conf`:

```bash
space_left_action = email
action_mail_acct = root
admin_space_left_action = halt
```

## Default Value

Not configured by default.

## References

1. CIS Controls v6.1 - 6.3 Ensure Audit Logging Systems Are Not Subject To Loss

## Profile

- Level 2 - Server
- Level 2 - Workstation
