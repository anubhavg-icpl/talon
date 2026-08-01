# stage: report
# category: CIS_benchmarks


# 6.1.13 Audit SUID executables (Not Scored)

## Profile Applicability

- Level 1 - Server
- Level 1 - Workstation

## Description

The owner of a file can set the file's permissions to run with the owner's or group's permissions, even if the user running the program is not the owner or a member of the group. The most common reason for a SUID program is to enable users to perform functions (such as changing their password) that require root privileges.

## Rationale

There are valid reasons for SUID programs, but it is important to identify and review such programs to ensure they are legitimate.

## Audit Procedure

Run the following command to list SUID files:

```bash
df --local -P | awk {'if (NR!=1) print $6'} | xargs -I '{}' find '{}' -xdev -type f -perm -4000
```

The command above only searches local filesystems, there may still be compromised items on network mounted partitions. Additionally the `--local` option to `df` is not universal to all versions, it can be omitted to search all filesystems on a system including network mounted filesystems or the following command can be run manually for each partition:

```bash
find <partition> -xdev -type f -perm -4000
```

## Expected Result

Review the list of SUID files and ensure no rogue SUID programs have been introduced into the system.

## Remediation

Ensure that no rogue SUID programs have been introduced into the system. Review the files returned by the action in the Audit section and confirm the integrity of these binaries.

## Default Value

Not applicable.

## References

None

## CIS Controls

5.1 Minimize And Sparingly Use Administrative Privileges - Minimize administrative privileges and only use administrative accounts when they are required. Implement focused auditing on the use of administrative privileged functions and monitor for anomalous behavior.

## Profile

- Level 1
