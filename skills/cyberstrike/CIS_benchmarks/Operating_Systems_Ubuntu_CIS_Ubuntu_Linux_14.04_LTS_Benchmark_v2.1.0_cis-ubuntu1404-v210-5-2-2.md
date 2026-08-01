# stage: report
# category: CIS_benchmarks


# 5.2.2 Ensure SSH Protocol is set to 2 (Scored)

## Profile Applicability

- Level 1 - Server
- Level 1 - Workstation

## Description

SSH supports two different and incompatible protocols: SSH1 and SSH2. SSH1 was the original protocol and was subject to security issues. SSH2 is more advanced and secure.

## Rationale

SSH v1 suffers from insecurities that do not affect SSH v2.

## Audit Procedure

Run the following command and verify that output matches:

```bash
grep "^Protocol" /etc/ssh/sshd_config
```

## Expected Result

```
Protocol 2
```

## Remediation

Edit the `/etc/ssh/sshd_config` file to set the parameter as follows:

```
Protocol 2
```

## Default Value

Protocol 2

## References

- CIS Controls: 3.4 - Use Only Secure Channels For Remote System Administration

## Profile

- Level 1 - Server
- Level 1 - Workstation
