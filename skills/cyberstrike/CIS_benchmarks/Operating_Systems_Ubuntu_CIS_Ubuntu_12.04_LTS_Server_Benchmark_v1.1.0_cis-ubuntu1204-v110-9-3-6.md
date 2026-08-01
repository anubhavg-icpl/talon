# stage: report
# category: CIS_benchmarks


# 9.3.6 Set SSH IgnoreRhosts to Yes (Scored)

## Profile Applicability

- Level 1

## Description

The `IgnoreRhosts` parameter specifies that `.rhosts` and `.shosts` files will not be used in `RhostsRSAAuthentication` or `HostbasedAuthentication`.

## Rationale

Setting this parameter forces users to enter a password when authenticating with ssh.

## Audit Procedure

### Using Command Line

To verify the correct SSH setting, run the following command and verify that the output is as shown:

```bash
grep "^IgnoreRhosts" /etc/ssh/sshd_config
```

## Expected Result

```
IgnoreRhosts yes
```

## Remediation

### Using Command Line

Edit the `/etc/ssh/sshd_config` file to set the parameter as follows:

```bash
IgnoreRhosts yes
```

## Default Value

IgnoreRhosts yes

## References

- CIS Ubuntu 12.04 LTS Server Benchmark v1.1.0

## Profile

Level 1 - Scored
