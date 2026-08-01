# stage: report
# category: CIS_benchmarks


# 13.18 Check for Duplicate Group Names (Scored)

## Profile Applicability

- Level 1

## Description

Although the `groupadd` program will not let you create a duplicate group name, it is possible for an administrator to manually edit the `/etc/group` file and change the group name.

## Rationale

If a group is assigned a duplicate group name, it will create and have access to files with the first GID for that group in `/etc/group`. Effectively, the GID is shared, which is a security problem.

## Audit Procedure

### Using Command Line

This script checks to make sure all group names in the `/etc/group` file are unique.

```bash
#!/bin/bash
cat /etc/group | /usr/bin/cut -f1 -d":" | /usr/bin/sort -n | /usr/bin/uniq -c |\
while read x ; do
    [ -z "${x}" ] && break
    set - $x
    if [ $1 -gt 1 ]; then
        gids=`/usr/bin/awk -F: '($1 == n) { print $3 }' n=$2 \
            /etc/group | xargs`
        echo "Duplicate Group Name ($2): ${gids}"
    fi
done
```

## Expected Result

No output should be returned. Any output indicates duplicate group names exist in `/etc/group`.

## Remediation

### Using Command Line

Based on the results of the script, establish unique names for the user groups. File group ownerships will automatically reflect the change as long as the groups have unique GIDs.

## Default Value

The `groupadd` command prevents creation of duplicate group names by default.

## References

- CIS Ubuntu 12.04 LTS Server Benchmark v1.1.0

## Profile

Level 1 - Scored
