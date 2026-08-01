# stage: report
# category: CIS_benchmarks


# 1.1.17 Ensure nodev option set on removable media partitions (Not Scored)

## Profile Applicability

- Level 1 - Server
- Level 1 - Workstation

## Description

The nodev mount option specifies that the filesystem cannot contain special devices.

## Rationale

Removable media containing character and block special devices could be used to circumvent security controls by allowing non-root users to access sensitive device files such as /dev/kmem or the raw disk partitions.

## Audit Procedure

```bash
# Run the following command and verify that the nodev option is set on all
# removable media partitions:
mount
```

## Expected Result

All removable media partitions should have the `nodev` option set.

## Remediation

```bash
# Edit the /etc/fstab file and add nodev to the fourth field (mounting options)
# of all removable media partitions. Look for entries that have mount points
# that contain words such as floppy or cdrom.
# See the fstab(5) manual page for more information.
```

## Default Value

By default, the nodev option is not set on removable media partitions.

## References

- CIS Ubuntu Linux 14.04 LTS Benchmark v2.1.0

## Profile

- Level 1
