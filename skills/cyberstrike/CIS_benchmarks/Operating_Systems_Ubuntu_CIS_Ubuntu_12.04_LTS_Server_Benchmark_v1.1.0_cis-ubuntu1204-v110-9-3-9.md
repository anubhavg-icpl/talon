# stage: report
# category: CIS_benchmarks


# 9.3.9 Set SSH PermitEmptyPasswords to No (Scored)

## Profile Applicability

- Level 1

## Description

The `PermitEmptyPasswords` parameter specifies if the server allows login to accounts with empty password strings.

## Rationale

Disallowing remote shell access to accounts that have an empty password reduces the probability of unauthorized access to the system.

## Audit Procedure

### Using Command Line

To verify the correct SSH setting, run the following command and verify that the output is as shown:

```bash
grep "^PermitEmptyPasswords" /etc/ssh/sshd_config
```

## Expected Result

```
PermitEmptyPasswords no
```

## Remediation

### Using Command Line

Edit the `/etc/ssh/sshd_config` file to set the parameter as follows:

```bash
PermitEmptyPasswords no
```

## Default Value

PermitEmptyPasswords no

## References

- CIS Ubuntu 12.04 LTS Server Benchmark v1.1.0

## Profile

Level 1 - Scored
