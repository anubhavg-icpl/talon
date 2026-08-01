# stage: report
# category: CIS_benchmarks


# Ensure separate partition exists for /home

## Profile

Level 2 - Server / Level 2 - Workstation, Assessment: Automated

## Description

The /home directory is used to support disk storage needs of local users.

## Rationale

If the system is intended to support local users, create a separate partition for the /home directory to protect against resource exhaustion and to set additional mount options to protect the filesystem.

## Impact

Resizing filesystems is a common activity in cloud-hosted servers. Separate filesystem partitions may prevent tying the overall disk management to one large volume. It may be necessary to use Logical Volume Management (LVM) or similar tools to create a flexible partitioning scheme for /home.

## Audit Procedure

### Command Line

Run the following command and verify the output shows that /home is mounted:

```bash
# findmnt -kn /home
```

## Expected Result

```
/home     /dev/sda6  ext4   rw,nosuid,nodev,relatime
```

## Remediation

### Command Line

For specific configuration requirements of the /home mount for your environment, modify /etc/fstab:

```bash
# Configure /etc/fstab as appropriate:
# Example:
# <device> /home   <fstype>   defaults,rw,nosuid,nodev,relatime   0 0
```

Creating a new partition is a destructive operation that typically requires a reinstallation or use of LVM.

## Default Value

By default /home is not a separate partition.

## References

1. http://tldp.org/HOWTO/LVM-HOWTO/
2. NIST SP 800-53 Rev. 5: MP-2

## CIS Controls

v8 - 3.3 Configure Data Access Control Lists (IG 1, IG 2, IG 3)
v7 - 14.6 Protect Information through Access Control Lists (IG 1, IG 2, IG 3)

MITRE ATT&CK Mappings: T1499, T1499.001 (TA0005) - M1038
