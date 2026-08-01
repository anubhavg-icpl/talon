# stage: report
# category: CIS_benchmarks


# 5.2.6 Ensure SSH IgnoreRhosts is enabled (Scored)

## Profile Applicability

- Level 1 - Server
- Level 1 - Workstation

## Description

The `IgnoreRhosts` parameter specifies that `.rhosts` and `.shosts` files will not be used in `RhostsRSAAuthentication` or `HostbasedAuthentication`.

## Rationale

Setting this parameter forces users to enter a password when authenticating with ssh.

## Audit Procedure

Run the following command and verify that output matches:

```bash
grep "^IgnoreRhosts" /etc/ssh/sshd_config
```

## Expected Result

```
IgnoreRhosts yes
```

## Remediation

Edit the `/etc/ssh/sshd_config` file to set the parameter as follows:

```
IgnoreRhosts yes
```

## Default Value

IgnoreRhosts yes

## References

- CIS Controls: 9 - Limitation and Control of Network Ports, Protocols, and Services

## Profile

- Level 1 - Server
- Level 1 - Workstation
