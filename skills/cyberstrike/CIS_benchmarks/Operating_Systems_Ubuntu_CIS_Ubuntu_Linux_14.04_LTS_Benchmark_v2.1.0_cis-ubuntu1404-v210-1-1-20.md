# stage: report
# category: CIS_benchmarks


# 1.1.20 Ensure sticky bit is set on all world-writable directories (Scored)

## Profile Applicability

- Level 1 - Server
- Level 1 - Workstation

## Description

Setting the sticky bit on world writable directories prevents users from deleting or renaming files in that directory that are not owned by them.

## Rationale

This feature prevents the ability to delete or rename files in world writable directories (such as /tmp) that are owned by another user.

## Audit Procedure

```bash
# Run the following command to verify no world writable directories exist
# without the sticky bit set:
df --local -P | awk {'if (NR!=1) print $6'} | xargs -I '{}' find '{}' -xdev -type d \( -perm -0002 -a ! -perm -1000 \) 2>/dev/null
# No output should be returned.
```

## Expected Result

No output should be returned. Any output indicates world-writable directories without the sticky bit set.

## Remediation

```bash
# Run the following command to set the sticky bit on all world writable directories:
df --local -P | awk {'if (NR!=1) print $6'} | xargs -I '{}' find '{}' -xdev -type d -perm -0002 2>/dev/null | xargs chmod a+t
```

## Default Value

By default, not all world-writable directories have the sticky bit set.

## Notes

Some distributions may not support the --local option to df.

## References

- CIS Ubuntu Linux 14.04 LTS Benchmark v2.1.0
- CIS Controls: 13 Data Protection

## Profile

- Level 1
