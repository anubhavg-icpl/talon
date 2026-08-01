# stage: report
# category: CIS_benchmarks


# 4.2.1.3 Ensure rsyslog default file permissions configured (Scored)

## Profile Applicability

- Level 1 - Server
- Level 1 - Workstation

## Description

rsyslog will create logfiles that do not already exist on the system. This setting controls what permissions will be applied to these newly created files.

## Rationale

It is important to ensure that log files have the correct permissions to ensure that sensitive data is archived and protected.

## Audit Procedure

Run the following command and verify that `$FileCreateMode` is `0640` or more restrictive:

```bash
grep ^\$FileCreateMode /etc/rsyslog.conf /etc/rsyslog.d/*.conf
```

## Expected Result

```
$FileCreateMode 0640
```

## Remediation

Edit the `/etc/rsyslog.conf` and `/etc/rsyslog.d/*.conf` files and set `$FileCreateMode` to `0640` or more restrictive:

```bash
$FileCreateMode 0640
```

**Notes:** You should also ensure this is not overridden with less restrictive settings in any `/etc/rsyslog.d/*` conf file.

## Default Value

Not explicitly set by default.

## References

1. See the rsyslog.conf(5) man page for more information.

## Profile

- Level 1 - Server
- Level 1 - Workstation
