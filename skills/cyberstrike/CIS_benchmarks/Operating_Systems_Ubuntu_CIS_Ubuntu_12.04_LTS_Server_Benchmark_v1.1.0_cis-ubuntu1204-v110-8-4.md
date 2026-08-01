# stage: report
# category: CIS_benchmarks


# 8.4 Configure logrotate (Not Scored)

## Profile Applicability

- Level 1

## Description

The system includes the capability of rotating log files regularly to avoid filling up the system with logs or making the logs unmanageable large. The file `/etc/logrotate.d/rsyslog` is the configuration file used to rotate log files created by `rsyslog`.

## Rationale

By keeping the log files smaller and more manageable, a system administrator can easily archive these files to another system and spend less time looking through inordinately large log files.

## Audit Procedure

### Using Command Line

Review the `/etc/logrotate.d/rsyslog` file to determine if the appropriate system logs are rotated according to your site policy.

```bash
cat /etc/logrotate.d/rsyslog
```

## Expected Result

The logrotate configuration should include appropriate system logs rotated according to site policy.

## Remediation

### Using Command Line

Edit `/etc/logrotate.d/rsyslog` file to include appropriate system logs according to your site policy.

```bash
vi /etc/logrotate.d/rsyslog
```

## Default Value

By default, logrotate is configured with basic rotation settings.

## References

- CIS Ubuntu 12.04 LTS Server Benchmark v1.1.0

## Profile

Level 1 - Not Scored
