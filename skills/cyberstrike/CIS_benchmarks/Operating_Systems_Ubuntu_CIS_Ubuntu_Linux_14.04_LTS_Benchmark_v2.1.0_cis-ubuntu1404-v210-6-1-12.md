# stage: report
# category: CIS_benchmarks


# 6.1.12 Ensure no ungrouped files or directories exist (Scored)

## Profile Applicability

- Level 1 - Server
- Level 1 - Workstation

## Description

Sometimes when administrators delete users or groups from the system they neglect to remove all files owned by those users or groups.

## Rationale

A new user who is assigned the deleted user's user ID or group ID may then end up "owning" these files, and thus have more access on the system than was intended.

## Audit Procedure

Run the following command and verify no files are returned:

```bash
df --local -P | awk {'if (NR!=1) print $6'} | xargs -I '{}' find '{}' -xdev -nogroup
```

The command above only searches local filesystems, there may still be compromised items on network mounted partitions. Additionally the `--local` option to `df` is not universal to all versions, it can be omitted to search all filesystems on a system including network mounted filesystems or the following command can be run manually for each partition:

```bash
find <partition> -xdev -nogroup
```

## Expected Result

No output should be returned.

## Remediation

Locate files that are owned by users or groups not listed in the system configuration files, and reset the ownership of these files to some active user on the system as appropriate.

## Default Value

Not applicable.

## References

None

## CIS Controls

14 Controlled Access Based on the Need to Know - Controlled Access Based on the Need to Know

## Profile

- Level 1
