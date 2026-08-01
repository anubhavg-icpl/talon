# stage: report
# category: CIS_benchmarks


# 4.3 Ensure logrotate is configured (Not Scored)

## Profile Applicability

- Level 1 - Server
- Level 1 - Workstation

## Description

The system includes the capability of rotating log files regularly to avoid filling up the system with logs or making the logs unmanageable. The file `/etc/logrotate.d/rsyslog` is the configuration file used to rotate log files created by `rsyslog`.

## Rationale

By keeping the log files smaller and more manageable, a system administrator can easily archive these files to another system and spend less time looking through inordinately large log files.

## Audit Procedure

Review `/etc/logrotate.conf` and `/etc/logrotate.d/*` and verify logs are rotated according to site policy.

```bash
cat /etc/logrotate.conf
ls /etc/logrotate.d/
```

## Expected Result

Log rotation should be configured with appropriate rotation intervals, retention policies, and compression settings according to site policy.

## Remediation

Edit `/etc/logrotate.conf` and `/etc/logrotate.d/*` to ensure logs are rotated according to site policy. The following is an example configuration:

```
/var/log/syslog
{
    rotate 7
    daily
    missingok
    notifempty
    delaycompress
    compress
    postrotate
        reload rsyslog >/dev/null 2>&1 || true
    endscript
}
```

## Default Value

logrotate is installed and configured with basic rotation settings by default.

## References

1. See the logrotate(8) man page for more information.

## Profile

- Level 1 - Server
- Level 1 - Workstation
