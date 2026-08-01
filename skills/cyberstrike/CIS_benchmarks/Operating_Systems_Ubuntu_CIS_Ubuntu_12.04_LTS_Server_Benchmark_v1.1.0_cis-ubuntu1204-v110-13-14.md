# stage: report
# category: CIS_benchmarks


# 13.14 Check for Duplicate UIDs (Scored)

## Profile Applicability

- Level 1

## Description

Although the `useradd` program will not let you create a duplicate User ID (UID), it is possible for an administrator to manually edit the `/etc/passwd` file and change the UID field.

## Rationale

Users must be assigned unique UIDs for accountability and to ensure appropriate access protections.

## Audit Procedure

### Using Command Line

This script checks to make sure all UIDs in the `/etc/passwd` file are unique.

```bash
#!/bin/bash
/bin/cat /etc/passwd | /usr/bin/cut -f3 -d":" | /usr/bin/sort -n | /usr/bin/uniq -c |\
while read x ; do
    [ -z "${x}" ] && break
    set - $x
    if [ $1 -gt 1 ]; then
        users=`/usr/bin/awk -F: '($3 == n) { print $1 }' n=$2 \
            /etc/passwd | /usr/bin/xargs`
        echo "Duplicate UID ($2): ${users}"
    fi
done
```

## Expected Result

No output should be returned. Any output indicates duplicate UIDs exist in `/etc/passwd`.

## Remediation

### Using Command Line

Based on the results of the script, establish unique UIDs and review all files owned by the shared UID to determine which UID they are supposed to belong to.

## Default Value

The `useradd` command assigns unique UIDs by default.

## References

- CIS Ubuntu 12.04 LTS Server Benchmark v1.1.0

## Profile

Level 1 - Scored
