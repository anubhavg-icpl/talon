# stage: report
# category: CIS_benchmarks


# 2.3 Set nosuid option for /tmp Partition (Scored)

## Profile Applicability

- Level 1

## Description

The `nosuid` mount option specifies that the filesystem cannot contain set userid files.

## Rationale

Since the `/tmp` filesystem is only intended for temporary file storage, set this option to ensure that users cannot create set userid files in `/tmp`.

## Audit Procedure

### Using Command Line

```bash
grep /tmp /etc/fstab | grep nosuid
mount | grep /tmp | grep nosuid
```

## Expected Result

Both commands should return output showing `nosuid` is set. If either command emits no output then the system is not configured as recommended.

## Remediation

### Using Command Line

Edit the `/etc/fstab` file and add `nosuid` to the fourth field (mounting options). See the `fstab(5)` manual page for more information.

```bash
mount -o remount,nosuid /tmp
```

## Default Value

By default, the `nosuid` option is not set on the `/tmp` partition.

## References

- CIS Ubuntu 12.04 LTS Server Benchmark v1.1.0

## Profile

Level 1 - Scored
