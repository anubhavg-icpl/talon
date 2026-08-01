# stage: report
# category: CIS_benchmarks


# 2.7 Create Separate Partition for /var/log (Scored)

## Profile Applicability

- Level 1

## Description

The `/var/log` directory is used by system services to store log data.

## Rationale

There are two important reasons to ensure that system logs are stored on a separate partition: protection against resource exhaustion (since logs can grow quite large) and protection of audit data.

## Audit Procedure

### Using Command Line

```bash
grep "[[:space:]]/var/log[[:space:]]" /etc/fstab
```

## Expected Result

A partition entry for `/var/log` should be returned. If the command emits no output then the system is not configured as recommended.

## Remediation

### Using Command Line

For new installations, during installation create a custom partition setup and specify a separate partition for `/var/log`.

For systems that were previously installed, use the Logical Volume Manager (LVM) to create partitions.

## Default Value

By default, Ubuntu does not create a separate partition for `/var/log`.

## References

- CIS Ubuntu 12.04 LTS Server Benchmark v1.1.0
- AJ Lewis, "LVM HOWTO", http://tldp.org/HOWTO/LVM-HOWTO/

## Profile

Level 1 - Scored
