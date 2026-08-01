# stage: report
# category: CIS_benchmarks


# 6.2.15 Ensure all groups in /etc/passwd exist in /etc/group (Automated)

## Description

Over time, system administration errors and changes can lead to groups being defined in `/etc/passwd` but not in `/etc/group`.

## Rationale

Groups defined in the `/etc/passwd` file but not in the `/etc/group` file pose a threat to system security since group permissions are not properly managed.

## Audit Procedure

Run the following script and verify no results are returned:

```bash
#!/bin/bash

for i in $(cut -s -d: -f4 /etc/passwd | sort -u ); do
  grep -q ":$i:" /etc/group
  if [ $? -ne 0 ]; then
    echo "Group $i is referenced by /etc/passwd but does not exist in /etc/group"
  fi
done
```

## Expected Result

No results should be returned.

## Remediation

Analyze the output of the Audit step above and perform the appropriate action to correct any discrepancies found.

## Profile

- Level 2 - Server
