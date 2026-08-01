# stage: report
# category: CIS_benchmarks


# 2.12 Add noexec Option to Removable Media Partitions (Not Scored)

## Profile Applicability

- Level 1

## Description

Set `noexec` on removable media to prevent programs from executing from the removable media.

## Rationale

Setting this option on a file system prevents users from executing programs from the removable media. This deters users from being able to introduce potentially malicious software on the system.

## Audit Procedure

### Using Command Line

```bash
grep <each removable media mountpoint> /etc/fstab
```

Verify that `noexec` is an option.

## Expected Result

The `noexec` option should be present in the mount options for each removable media mountpoint.

## Remediation

### Using Command Line

Edit the `/etc/fstab` file and add `noexec` to the fourth field (mounting options). Look for entries that have mount points that contain words such as floppy or cdrom. See the `fstab(5)` manual page for more information.

## Default Value

By default, the `noexec` option is not set on removable media partitions.

## References

- CIS Ubuntu 12.04 LTS Server Benchmark v1.1.0

## Profile

Level 1 - Not Scored
