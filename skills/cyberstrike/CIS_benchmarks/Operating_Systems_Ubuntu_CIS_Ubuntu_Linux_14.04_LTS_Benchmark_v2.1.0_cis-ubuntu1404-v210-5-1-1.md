# stage: report
# category: CIS_benchmarks


# 5.1.1 Ensure cron daemon is enabled (Scored)

## Profile Applicability

- Level 1 - Server
- Level 1 - Workstation

## Description

The `cron` daemon is used to execute batch jobs on the system.

## Rationale

While there may not be user jobs that need to be run on the system, the system does have maintenance jobs that may include security monitoring that have to run, and `cron` is used to execute them.

## Audit Procedure

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

Edit start lines in `/etc/init/cron.conf` to match the following:

```bash
start on runlevel [2345]
```

## Default Value

Enabled by default on Ubuntu 14.04.

## References

- CIS Controls: 6 - Maintenance, Monitoring, and Analysis of Audit Logs

## Profile

- Level 1 - Server
- Level 1 - Workstation
