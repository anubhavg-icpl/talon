# stage: report
# category: CIS_benchmarks


# 10.1.3 Set Password Expiring Warning Days (Scored)

## Profile Applicability

- Level 1

## Description

The `PASS_WARN_AGE` parameter in `/etc/login.defs` allows an administrator to notify users that their password will expire in a defined number of days. It is recommended that the `PASS_WARN_AGE` parameter be set to 7 or more days.

## Rationale

Providing an advance warning that a password will be expiring gives users time to think of a secure password. Users caught unaware may choose a simple password or write it down where it may be discovered.

## Audit Procedure

### Using Command Line

```bash
grep PASS_WARN_AGE /etc/login.defs
# Expected output: PASS_WARN_AGE 7

chage --list <user>
# Verify: Number of days of warning before password expires: 7
```

## Expected Result

`PASS_WARN_AGE` is set to 7 or more in `/etc/login.defs`, and all active users have a password warning age of 7 days or more.

## Remediation

### Using Command Line

Set the `PASS_WARN_AGE` parameter to 7 in `/etc/login.defs`:

```bash
# Edit /etc/login.defs and set:
PASS_WARN_AGE 7

# Modify active user parameters to match:
chage --warndays 7 <user>
```

## Default Value

PASS_WARN_AGE 7

## References

- CIS Ubuntu 12.04 LTS Server Benchmark v1.1.0

## Profile

Level 1 - Scored
