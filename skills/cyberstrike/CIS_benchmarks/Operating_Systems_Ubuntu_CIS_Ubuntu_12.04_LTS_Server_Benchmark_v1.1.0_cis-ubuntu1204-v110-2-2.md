# stage: report
# category: CIS_benchmarks


# 2.2 Set nodev option for /tmp Partition (Scored)

## Profile Applicability

- Level 1

## Description

The `nodev` mount option specifies that the filesystem cannot contain special devices.

## Rationale

Since the `/tmp` filesystem is not intended to support devices, set this option to ensure that users cannot attempt to create block or character special devices in `/tmp`.

## Audit Procedure

### Using Command Line

```bash
grep /tmp /etc/fstab | grep nodev
mount | grep /tmp | grep nodev
```

## Expected Result

Both commands should return output showing `nodev` is set. If either command emits no output then the system is not configured as recommended.

## Remediation

### Using Command Line

Edit the `/etc/fstab` file and add `nodev` to the fourth field (mounting options). See the `fstab(5)` manual page for more information.

```bash
mount -o remount,nodev /tmp
```

## Default Value

By default, the `nodev` option is not set on the `/tmp` partition.

## References

- CIS Ubuntu 12.04 LTS Server Benchmark v1.1.0

## Profile

Level 1 - Scored
