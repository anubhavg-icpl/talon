# stage: report
# category: CIS_benchmarks


# 4.2.3 Ensure rsyslog or syslog-ng is installed (Scored)

## Profile Applicability

- Level 1 - Server
- Level 1 - Workstation

## Description

The `rsyslog` and `syslog-ng` software are recommended replacements to the original `syslogd` daemon which provide improvements over `syslogd`, such as connection-oriented (i.e. TCP) transmission of logs, the option to log to database formats, and the encryption of log data en route to a central logging server.

## Rationale

The security enhancements of `rsyslog` and `syslog-ng` such as connection-oriented (i.e. TCP) transmission of logs, the option to log to database formats, and the encryption of log data en route to a central logging server) justify installing and configuring the package.

## Audit Procedure

Verify either rsyslog or syslog-ng is installed. Depending on the package management in use one of the following command groups may provide the needed information:

```bash
dpkg -s rsyslog
dpkg -s syslog-ng
```

## Expected Result

One of the packages should show as installed (`Status: install ok installed`).

## Remediation

Install rsyslog or `syslog-ng` using one of the following commands:

```bash
apt-get install rsyslog
apt-get install syslog-ng
```

## Default Value

rsyslog is typically installed by default on Ubuntu 14.04.

## References

1. CIS Controls v6.1 - 6.2 Ensure Audit Log Settings Support Appropriate Log Entry Formatting

## Profile

- Level 1 - Server
- Level 1 - Workstation
