# stage: report
# category: CIS_benchmarks


# CIS Ubuntu Linux 18.04 LTS Benchmark v2.2.0 - Control 4.5.3

## Profile Applicability

- Level 1 - Server
- Level 1 - Workstation

## Description

The `usermod` command can be used to specify which group the root user belongs to. This affects permissions of files that are created by the root user.

## Rationale

Using GID 0 for the `root` account helps prevent `root`-owned files from accidentally becoming accessible to non-privileged users.

## Audit Procedure

### Command Line

Run the following command and verify the result is `0`:

```bash
grep "^root:" /etc/passwd | cut -f4 -d:
```

### Expected Result

```
0
```

## Remediation

### Command Line

Run the following command to set the `root` user default group to GID `0`:

```bash
usermod -g 0 root
```

## References

1. NIST SP 800-53 Rev. 5: CM-1, CM-2, CM-6, CM-7, IA-5

## CIS Controls

v8 - 3.3 Configure Data Access Control Lists - Configure data access control lists based on a user's need to know.

v7 - 14.6 Protect Information through Access Control Lists.

## Profile Applicability

- Level 1 - Server
- Level 1 - Workstation

## Assessment Status

Automated
