# stage: report
# category: CIS_benchmarks


# CIS 2.6 — Set Group named or root for BIND Directories and Files

## Profile Applicability

- Level 1 - Authoritative Name Server
- Level 1 - Caching Only Name Server

## Description

All the BIND directories and files should have a group of either named or root.

## Rationale

In general, the BIND directories and files default to a group of named, however some system files may have a group of root. Examples of system files include `chroot`'ed system device files. Either group root or named is accepted, as the intent is to prevent unexpected group ids, from getting inappropriate access to BIND files. Run time directories to which BIND will need write access should have a group of named, so that write access may be granted via the group permissions.

## Impact

Not Applicable

## Audit Procedure

Ensure that the BIND benchmark variables used below are set as described in the benchmark overview. Run the command below to ensure that all BIND directories and files have a group of either named or root.

```bash
# find $BIND_HOME $RUNDIR \! \( -group root -o -group named \) -ls
```

There should be no files listed in the output from the find command.

## Remediation

Run the command below to change all BIND directories and files to the group named.

```bash
chgrp -R named $BIND_HOME $RUNDIR
```

## Default Value

The default rpm install has all directories and files in the BIND home and the run time directory with a group of named.

## References

Not Applicable

## CIS Controls

| Controls Version | Control                                              | IG 1 | IG 2 | IG 3 |
| ---------------- | ---------------------------------------------------- | ---- | ---- | ---- |
| v6               | 14.4 - Protect Information with Access Control Lists | Y    | Y    | Y    |

## MITRE ATT&CK Mappings

| Tactic          | Technique                                           |
| --------------- | --------------------------------------------------- |
| Defense Evasion | T1222 - File and Directory Permissions Modification |
| Persistence     | T1546 - Event Triggered Execution                   |

## Profile

- Level 1 - Authoritative Name Server
- Level 1 - Caching Only Name Server
