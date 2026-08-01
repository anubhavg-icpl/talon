# stage: report
# category: CIS_benchmarks


# 6.2.6 Ensure no duplicate GIDs exist (Automated)

## Profile

- Level 1 - Server
- Level 1 - Workstation

## Description

Although the `groupadd` program will not let you create a duplicate Group ID (GID), it is possible for an administrator to manually edit the `/etc/group` file and change the GID field.

## Rationale

User groups must be assigned unique GIDs for accountability and to ensure appropriate access protections.

## Audit Procedure

### Command Line

Run the following script and verify no results are returned:

```bash
#!/usr/bin/env bash

{
  while read -r l_count l_gid; do
      if [ "$l_count" -gt 1 ]; then
        echo -e "Duplicate GID: \"$l_gid\" Groups: \"$(awk -F: '($3 == n) { print $1 }' n=$l_gid /etc/group | xargs)\""
      fi
  done < <(cut -f3 -d":" /etc/group | sort -n | uniq -c)
}
```

## Expected Result

No output should be returned.

## Remediation

### Command Line

Based on the results of the audit script, establish unique GIDs and review all files owned by the shared GID to determine which group they are supposed to belong to.

## References

1. NIST SP 800-53 Rev. 5: CM-1, CM-2, CM-6, CM-7, IA-5

## Additional Information

You can also use the `grpck` command to check for other inconsistencies in the `/etc/group` file.

## CIS Controls

- v8: **5.2** Use Unique Passwords
- v7: **4.4** Use Unique Passwords
