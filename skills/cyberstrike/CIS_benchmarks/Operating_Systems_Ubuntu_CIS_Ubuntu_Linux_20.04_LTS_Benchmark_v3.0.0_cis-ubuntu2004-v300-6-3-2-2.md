# stage: report
# category: CIS_benchmarks


# 6.3.2.2 Ensure audit logs are not automatically deleted (Automated)

## Profile

- Level 2 - Server
- Level 2 - Workstation

## Description

The `max_log_file_action` setting determines how to handle the audit log file reaching the max file size. A value of `keep_logs` will rotate the logs but never delete old logs.

## Rationale

In high security contexts, the benefits of maintaining a long audit history exceed the cost of storing the audit history.

## Audit Procedure

### Command Line

Run the following command and verify output matches:

```bash
# grep max_log_file_action /etc/audit/auditd.conf
max_log_file_action = keep_logs
```

## Expected Result

The output should show `max_log_file_action = keep_logs`.

## Remediation

### Command Line

Set the following parameter in `/etc/audit/auditd.conf`:

```
max_log_file_action = keep_logs
```

## References

1. NIST SP 800-53 Rev. 5: AU-8

## CIS Controls

| Controls Version | Control                               | IG 1 | IG 2 | IG 3 |
| ---------------- | ------------------------------------- | ---- | ---- | ---- |
| v8               | 8.3 Ensure Adequate Audit Log Storage | X    | X    | X    |
| v7               | 6.4 Ensure adequate storage for logs  |      | X    | X    |

MITRE ATT&CK Mappings: T1562, T1562.006 / TA0005 / M1053
