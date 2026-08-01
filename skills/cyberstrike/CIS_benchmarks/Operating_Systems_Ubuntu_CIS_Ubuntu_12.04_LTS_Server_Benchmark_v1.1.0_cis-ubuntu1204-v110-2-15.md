# stage: report
# category: CIS_benchmarks


# 2.15 Add nosuid Option to /run/shm Partition (Scored)

## Profile Applicability

- Level 1

## Description

The `nosuid` mount option specifies that the `/run/shm` (temporary filesystem stored in memory) will not execute `setuid` and `setgid` on executable programs as such, but rather execute them with the uid and gid of the user executing the program.

## Rationale

Setting this option on a file system prevents users from introducing privileged programs onto the system and allowing non-root users to execute them.

## Audit Procedure

### Using Command Line

```bash
grep /run/shm /etc/fstab | grep nosuid
mount | grep /run/shm | grep nosuid
```

## Expected Result

Both commands should return output showing `nosuid` is set. If either command emits no output then the system is not configured as recommended.

## Remediation

### Using Command Line

Edit the `/etc/fstab` file and add `nosuid` to the fourth field (mounting options). Look for entries that have mount points that contain `/run/shm`. See the `fstab(5)` manual page for more information.

```bash
mount -o remount,nosuid /run/shm
```

## Default Value

By default, the `nosuid` option is not set on the `/run/shm` partition.

## References

- CIS Ubuntu 12.04 LTS Server Benchmark v1.1.0

## Profile

Level 1 - Scored
