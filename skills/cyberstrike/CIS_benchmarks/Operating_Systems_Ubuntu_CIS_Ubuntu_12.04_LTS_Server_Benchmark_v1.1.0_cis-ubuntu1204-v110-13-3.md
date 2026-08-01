# stage: report
# category: CIS_benchmarks


# 13.3 Verify No Legacy "+" Entries Exist in /etc/shadow File (Scored)

## Profile Applicability

- Level 1

## Description

The character `+` in various files used to be markers for systems to insert data from NIS maps at a certain point in a system configuration file. These entries are no longer required on most systems, but may exist in files that have been imported from other platforms.

## Rationale

These entries may provide an avenue for attackers to gain privileged access on the system.

## Audit Procedure

### Using Command Line

Run the following command and verify that no output is returned:

```bash
/bin/grep '^+:' /etc/shadow
```

## Expected Result

No output should be returned. Any output indicates legacy `+` entries exist in `/etc/shadow`.

## Remediation

### Using Command Line

Delete these entries if they exist.

## Default Value

No legacy `+` entries exist in `/etc/shadow` by default on Ubuntu 12.04.

## References

- CIS Ubuntu 12.04 LTS Server Benchmark v1.1.0

## Profile

Level 1 - Scored
