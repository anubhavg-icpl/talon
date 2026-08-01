# stage: report
# category: CIS_benchmarks


# 9.3.1 Set SSH Protocol to 2 (Scored)

## Profile Applicability

- Level 1

## Description

SSH supports two different and incompatible protocols: SSH1 and SSH2. SSH1 was the original protocol and was subject to security issues. SSH2 is more advanced and secure.

## Rationale

SSH v1 suffers from insecurities that do not affect SSH v2.

## Audit Procedure

### Using Command Line

To verify the correct SSH setting, run the following command and verify that the output is as shown:

```bash
grep "^Protocol" /etc/ssh/sshd_config
```

## Expected Result

```
Protocol 2
```

## Remediation

### Using Command Line

Edit the `/etc/ssh/sshd_config` file to set the parameter as follows:

```bash
Protocol 2
```

## Default Value

Protocol 2

## References

- CIS Ubuntu 12.04 LTS Server Benchmark v1.1.0

## Profile

Level 1 - Scored
