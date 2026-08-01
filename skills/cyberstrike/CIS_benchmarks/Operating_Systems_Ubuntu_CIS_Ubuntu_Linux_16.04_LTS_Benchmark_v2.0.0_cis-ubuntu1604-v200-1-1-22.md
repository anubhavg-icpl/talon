# stage: report
# category: CIS_benchmarks


# Ensure sticky bit is set on all world-writable directories

## Description

Setting the sticky bit on world writable directories prevents users from deleting or renaming files in that directory that are not owned by them.

## Rationale

This feature prevents the ability to delete or rename files in world writable directories (such as /tmp) that are owned by another user.

## Impact

None noted.

## Audit Procedure

### Command Line

```bash
# Run the following command to verify no world writable directories exist without the sticky bit set:
df --local -P | awk '(if (NR!=1) print $6}' | xargs -I '{}' find '{}' -xdev -type d \( -perm -0002 -a ! -perm -1000 \) 2>/dev/null
```

## Expected Result

No output should be returned.

## Remediation

### Command Line

```bash
# Run the following command to set the sticky bit on all world writable directories:
df --local -P | awk '{if (NR!=1) print $6}' | xargs -I '{}' find '{}' -xdev -type d \( -perm -0002 -a ! -perm -1000 \) 2>/dev/null | xargs -I '{}' chmod a+t '{}'
```

## Default Value

The sticky bit is set on /tmp by default.

## References

- CIS Controls Version 7 - 5.1 Establish Secure Configurations: Maintain documented, standard security configuration standards for all authorized operating systems and software.

## Profile

Level 1 - Server / Level 1 - Workstation, Assessment: Automated
