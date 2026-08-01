# stage: report
# category: CIS_benchmarks


# 4.1.1.1 Ensure audit log storage size is configured (Not Scored)

## Profile Applicability

- Level 2 - Server
- Level 2 - Workstation

## Description

Configure the maximum size of the audit log file. Once the log reaches the maximum size, it will be rotated and a new log file will be started.

## Rationale

It is important that an appropriate size is determined for log files so that they do not impact the system and audit data is not lost.

## Audit Procedure

Run the following command and ensure output is in compliance with site policy:

```bash
grep max_log_file /etc/audit/auditd.conf
```

## Expected Result

```
max_log_file = <MB>
```

The output should show `max_log_file` set to a value in accordance with site policy.

## Remediation

Set the following parameter in `/etc/audit/auditd.conf` in accordance with site policy:

```bash
max_log_file = <MB>
```

**Notes:** The max_log_file parameter is measured in megabytes.

## Default Value

By default, auditd will max out the log files at 5MB and retain only 4 copies of them.

## References

1. CIS Controls v6.1 - 6.3 Ensure Audit Logging Systems Are Not Subject To Loss

## Profile

- Level 2 - Server
- Level 2 - Workstation
