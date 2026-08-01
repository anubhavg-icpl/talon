# stage: report
# category: CIS_benchmarks


# 6.2.4 Ensure no legacy "+" entries exist in /etc/group (Automated)

## Description

The character + in various files used to be markers for systems to insert data from NIS maps at a certain point in a system configuration file. These entries are no longer required on most systems, but may exist in files that have been imported from other platforms.

## Rationale

These entries may provide an avenue for attackers to gain privileged access on the system.

## Audit Procedure

Run the following command and verify that no output is returned:

```bash
# grep '^\+:' /etc/group
```

## Expected Result

No output should be returned.

## Remediation

Remove any legacy '+' entries from `/etc/group` if they exist.

## CIS Controls

| Controls Version | Control                                            | IG 1 | IG 2 | IG 3 |
| ---------------- | -------------------------------------------------- | ---- | ---- | ---- |
| v8               | 5.6 Centralize Account Management                  |      | x    | x    |
| v7               | 16.2 Configure Centralized Point of Authentication |      | x    | x    |

## Profile

- Level 1 - Server
