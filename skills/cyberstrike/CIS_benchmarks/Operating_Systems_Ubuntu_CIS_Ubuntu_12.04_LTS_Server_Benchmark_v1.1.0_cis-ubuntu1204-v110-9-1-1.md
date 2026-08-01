# stage: report
# category: CIS_benchmarks


# 9.1.1 Enable cron Daemon (Scored)

## Profile Applicability

- Level 1

## Description

The `cron` daemon is used to execute batch jobs on the system.

## Rationale

While there may not be user jobs that need to be run on the system, the system does have maintenance jobs that may include security monitoring that have to run and `cron` is used to execute them.

## Audit Procedure

### Using Command Line

Ensure proper start conditions listed for `cron`:

```bash
/sbin/initctl show-config cron
```

## Expected Result

```
cron
  start on runlevel [2345]
  stop on runlevel [!2345]
```

## Remediation

### Using Command Line

Edit start lines in `/etc/init/cron.conf` to match the following:

```bash
start on runlevel [2345]
```

## Default Value

cron is enabled by default.

## References

- CIS Ubuntu 12.04 LTS Server Benchmark v1.1.0

## Profile

Level 1 - Scored
