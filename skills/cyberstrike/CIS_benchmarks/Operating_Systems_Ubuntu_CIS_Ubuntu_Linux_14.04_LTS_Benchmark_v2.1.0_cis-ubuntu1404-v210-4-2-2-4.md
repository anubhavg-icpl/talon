# stage: report
# category: CIS_benchmarks


# 4.2.2.4 Ensure syslog-ng is configured to send logs to a remote log host (Not Scored)

## Profile Applicability

- Level 1 - Server
- Level 1 - Workstation

## Description

The `syslog-ng` utility supports the ability to send logs it gathers to a remote log host or to receive messages from remote hosts, reducing administrative overhead.

## Rationale

Storing log data on a remote host protects log integrity from local attacks. If an attacker gains root access on the local system, they could tamper with or remove log data that is stored on the local system.

## Audit Procedure

Review the `/etc/syslog-ng/syslog-ng.conf` file and verify that logs are sent to a central host (where `logfile.example.com` is the name of your central log host):

```bash
destination logserver { tcp("logfile.example.com" port(514)); };
log { source(src); destination(logserver); };
```

## Expected Result

The destination and log directives should be present and configured to forward to a remote log host.

## Remediation

Edit the `/etc/syslog-ng/syslog-ng.conf` file and add the following lines (where `logfile.example.com` is the name of your central log host).

```bash
destination logserver { tcp("logfile.example.com" port(514)); };
log { source(src); destination(logserver); };
```

Run the following command to reload the `syslog-ng` configuration:

```bash
pkill -HUP syslog-ng
```

## Default Value

Not configured by default.

## References

1. See the syslog-ng.conf(5) man page for more information.

## Profile

- Level 1 - Server
- Level 1 - Workstation
