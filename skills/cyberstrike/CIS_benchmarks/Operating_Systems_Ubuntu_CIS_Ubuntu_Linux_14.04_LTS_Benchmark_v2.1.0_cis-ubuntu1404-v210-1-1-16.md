# stage: report
# category: CIS_benchmarks


# 1.1.16 Ensure noexec option set on /run/shm partition (Scored)

## Profile Applicability

- Level 1 - Server
- Level 1 - Workstation

## Description

The noexec mount option specifies that the filesystem cannot contain executable binaries.

## Rationale

Setting this option on a file system prevents users from executing programs from shared memory. This deters users from introducing potentially malicious software on the system.

## Audit Procedure

```bash
mount | grep /run/shm
# Verify that the noexec option is set on /run/shm
# Expected output: tmpfs on /run/shm type tmpfs (rw,nosuid,nodev,noexec,relatime)
```

## Expected Result

The output should show the `noexec` option is set for the /run/shm partition.

## Remediation

```bash
# Edit the /etc/fstab file and add noexec to the fourth field (mounting options)
# for the /run/shm partition. See the fstab(5) manual page for more information.

# Run the following command to remount /run/shm:
mount -o remount,noexec /run/shm
```

## Default Value

By default, the noexec option is not set on /run/shm.

## References

- CIS Ubuntu Linux 14.04 LTS Benchmark v2.1.0

## Profile

- Level 1
