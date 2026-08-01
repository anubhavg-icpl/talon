# stage: report
# category: CIS_benchmarks


# 5.1.12 Ensure SSH PermitUserEnvironment is disabled (Automated)

## Description

The `PermitUserEnvironment` option allows users to present environment options to the `ssh` daemon.

## Rationale

Permitting users the ability to set environment variables through the SSH daemon could potentially allow users to bypass security controls (e.g. setting an execution path that has `ssh` executing trojan'd programs).

## Audit Procedure

Run the following command and verify that output matches:

```bash
# sshd -T | grep permituserenvironment

PermitUserEnvironment no
```

## Expected Result

Output should show `PermitUserEnvironment no`.

## Remediation

Edit the `/etc/ssh/sshd_config` file to set the parameter as follows:

```
PermitUserEnvironment no
```

## Default Value

PermitUserEnvironment no

## CIS Controls

| Controls Version | Control                                 | IG 1 | IG 2 | IG 3 |
| ---------------- | --------------------------------------- | ---- | ---- | ---- |
| v8               | 3.3 Configure Data Access Control Lists | x    | x    | x    |
| v7               | 5.1 Establish Secure Configurations     | x    | x    | x    |

## Profile

- Level 1 - Server
