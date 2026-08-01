# stage: report
# category: CIS_benchmarks


# Ensure nodev option set on removable media partitions

## Description

The nodev mount option specifies that the filesystem cannot contain special devices.

## Rationale

Removable media containing character and block special devices could be used to circumvent security controls by allowing non-root users to access sensitive device files such as /dev/kmem or the raw disk partitions.

## Impact

None noted.

## Audit Procedure

### Command Line

```bash
# Run the following command and verify that the nodev option is set on all removable media partitions.
mount
```

## Expected Result

Verify that the nodev option is set on all removable media partitions.

## Remediation

### Command Line

```bash
# Edit the /etc/fstab file and add nodev to the fourth field (mounting options) of all removable media partitions.
# Look for entries that have mount points that contain words such as floppy or cdrom.
# See the fstab(5) manual page for more information.
```

## Default Value

Not set by default.

## References

- CIS Controls Version 7 - 5.1 Establish Secure Configurations: Maintain documented, standard security configuration standards for all authorized operating systems and software.

## Profile

Level 1 - Server / Level 1 - Workstation, Assessment: Manual
