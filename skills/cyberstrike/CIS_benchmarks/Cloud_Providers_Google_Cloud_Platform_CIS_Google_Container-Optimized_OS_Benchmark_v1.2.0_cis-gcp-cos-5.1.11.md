# stage: report
# category: CIS_benchmarks


# 5.1.11 Ensure SSH PermitEmptyPasswords is disabled (Automated)

## Description

The `PermitEmptyPasswords` parameter specifies if the SSH server allows login to accounts with empty password strings.

## Rationale

Disallowing remote shell access to accounts that have an empty password reduces the probability of unauthorized access to the system.

## Audit Procedure

Run the following command and verify that output matches:

```bash
# sshd -T | grep permitemptypasswords

PermitEmptyPasswords no
```

## Expected Result

Output should show `PermitEmptyPasswords no`.

## Remediation

Edit the `/etc/ssh/sshd_config` file to set the parameter as follows:

```
PermitEmptyPasswords no
```

## Default Value

PermitEmptyPasswords no

## CIS Controls

| Controls Version | Control                                             | IG 1 | IG 2 | IG 3 |
| ---------------- | --------------------------------------------------- | ---- | ---- | ---- |
| v8               | 6.3 Require MFA for Externally-Exposed Applications |      | x    | x    |
| v7               | 16.3 Require Multi-factor Authentication            |      | x    | x    |

## Profile

- Level 1 - Server
