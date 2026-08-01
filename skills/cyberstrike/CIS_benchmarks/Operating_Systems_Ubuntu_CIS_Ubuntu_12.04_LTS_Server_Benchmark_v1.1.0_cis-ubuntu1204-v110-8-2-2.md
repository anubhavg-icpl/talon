# stage: report
# category: CIS_benchmarks


# 8.2.2 Ensure the rsyslog Service is activated (Scored)

## Profile Applicability

- Level 1

## Description

Once the rsyslog package is installed it needs to be activated.

## Rationale

If the rsyslog service is not activated the system will not have a syslog service running.

## Audit Procedure

### Using Command Line

Ensure that the `rsyslog` service is active:

```bash
initctl show-config rsyslog
```

## Expected Result

```
rsyslog start on filesystem
stop on runlevel [06]
```

## Remediation

### Using Command Line

Set the proper start conditions in `/etc/init/rsyslog.conf`:

```bash
start on filesystem
```

## Default Value

The rsyslog service may not be activated by default if it was manually installed.

## References

- CIS Ubuntu 12.04 LTS Server Benchmark v1.1.0

## Profile

Level 1 - Scored
