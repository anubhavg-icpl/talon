# stage: report
# category: CIS_benchmarks


# 9.3.5 Set SSH MaxAuthTries to 4 or Less (Scored)

## Profile Applicability

- Level 1

## Description

The `MaxAuthTries` parameter specifies the maximum number of authentication attempts permitted per connection. When the login failure count reaches half the number, error messages will be written to the `syslog` file detailing the login failure.

## Rationale

Setting the `MaxAuthTries` parameter to a low number will minimize the risk of successful brute force attacks to the SSH server. While the recommended setting is 4, it is set the number based on site policy.

## Audit Procedure

### Using Command Line

To verify the correct SSH setting, run the following command and verify that the output is as shown:

```bash
grep "^MaxAuthTries" /etc/ssh/sshd_config
```

## Expected Result

```
MaxAuthTries 4
```

## Remediation

### Using Command Line

Edit the `/etc/ssh/sshd_config` file to set the parameter as follows:

```bash
MaxAuthTries 4
```

## Default Value

MaxAuthTries 6

## References

- CIS Ubuntu 12.04 LTS Server Benchmark v1.1.0

## Profile

Level 1 - Scored
