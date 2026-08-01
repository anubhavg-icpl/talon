# stage: report
# category: CIS_benchmarks


# Ensure /home partition includes the nodev option

## Description

The nodev mount option specifies that the filesystem cannot contain special devices.

## Rationale

Since the user partitions are not intended to support devices, set this option to ensure that users cannot attempt to create block or character special devices.

## Impact

None noted.

## Audit Procedure

### Command Line

```bash
# If a /home partition exists, verify that the nodev option is set.
# Run the following command and verify that nothing is returned:
findmnt -n /home | grep -v nodev
```

## Expected Result

No output should be returned.

## Remediation

### Command Line

```bash
# Edit the /etc/fstab file and add nodev to the fourth field (mounting options) for the /home partition.
# See the fstab(5) manual page for more information.
mount -o remount,nodev /home
```

## Default Value

Not set by default.

## References

- CIS Controls Version 7 - 5.1 Establish Secure Configurations: Maintain documented, standard security configuration standards for all authorized operating systems and software.

## Profile

Level 1 - Server / Level 1 - Workstation, Assessment: Automated
