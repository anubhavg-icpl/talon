# stage: report
# category: CIS_benchmarks


# 5.1.4 Ensure SSH Protocol is set to 2 (Automated)

## Description

Older versions of SSH support two different and incompatible protocols: SSH1 and SSH2. SSH1 was the original protocol and was subject to security issues. SSH2 is more advanced and secure.

## Rationale

SSH v1 suffers from insecurities that do not affect SSH v2.

## Audit Procedure

Run the following command and verify that output matches:

```bash
# grep ^Protocol /etc/ssh/sshd_config

Protocol 2
```

## Expected Result

Output should show `Protocol 2`.

## Remediation

Edit the `/etc/ssh/sshd_config` file to set the parameter as follows:

```
Protocol 2
```

## Default Value

This command no longer exists in newer versions of SSH. This check is still being included for systems that may be running an older version of SSH. As of openSSH version 7.4 this parameter will not cause an issue when included.

## CIS Controls

| Controls Version | Control                                                          | IG 1 | IG 2 | IG 3 |
| ---------------- | ---------------------------------------------------------------- | ---- | ---- | ---- |
| v8               | 3.10 Encrypt Sensitive Data in Transit                           |      | x    | x    |
| v8               | 4.6 Securely Manage Enterprise Assets and Software               | x    | x    | x    |
| v7               | 4.5 Use Multifactor Authentication For All Administrative Access |      | x    | x    |
| v7               | 14.4 Encrypt All Sensitive Information in Transit                |      | x    | x    |

## Profile

- Level 1 - Server
