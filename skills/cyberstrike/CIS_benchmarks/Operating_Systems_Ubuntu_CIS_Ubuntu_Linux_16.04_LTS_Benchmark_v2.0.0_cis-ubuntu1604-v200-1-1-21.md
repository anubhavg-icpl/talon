# stage: report
# category: CIS_benchmarks


# Ensure noexec option set on removable media partitions

## Description

The noexec mount option specifies that the filesystem cannot contain executable binaries.

## Rationale

Setting this option on a file system prevents users from executing programs from the removable media. This deters users from being able to introduce potentially malicious software on the system.

## Impact

None noted.

## Audit Procedure

### Command Line

```bash
# Run the following command and verify that the noexec option is set on all removable media partitions.
mount
```

## Expected Result

Verify that the noexec option is set on all removable media partitions.

## Remediation

### Command Line

```bash
# Edit the /etc/fstab file and add noexec to the fourth field (mounting options) of all removable media partitions.
# Look for entries that have mount points that contain words such as floppy or cdrom.
# See the fstab(5) manual page for more information.
```

## Default Value

Not set by default.

## References

- CIS Controls Version 7 - 2.6 Address unapproved software: Ensure that unauthorized software is either removed or the inventory is updated in a timely manner.

## Profile

Level 1 - Server / Level 1 - Workstation, Assessment: Manual
