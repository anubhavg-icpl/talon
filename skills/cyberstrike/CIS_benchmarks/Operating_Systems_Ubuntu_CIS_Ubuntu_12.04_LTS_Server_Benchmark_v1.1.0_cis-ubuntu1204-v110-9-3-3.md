# stage: report
# category: CIS_benchmarks


# 9.3.3 Set Permissions on /etc/ssh/sshd_config (Scored)

## Profile Applicability

- Level 1

## Description

The `/etc/ssh/sshd_config` file contains configuration specifications for `sshd`. The command below sets the owner and group of the file to root.

## Rationale

The `/etc/ssh/sshd_config` file needs to be protected from unauthorized changes by nonprivileged users, but needs to be readable as this information is used with many nonprivileged programs.

## Audit Procedure

### Using Command Line

Run the following command to determine the user and group ownership on the `/etc/ssh/sshd_config` file:

```bash
/bin/ls -l /etc/ssh/sshd_config
```

## Expected Result

```
-rw------- 1 root root 762 Sep 23 002 /etc/ssh/sshd_config
```

## Remediation

### Using Command Line

If the user and group ownership of the `/etc/ssh/sshd_config` file are incorrect, run the following command to correct them:

```bash
chown root:root /etc/ssh/sshd_config
```

If the permissions are incorrect, run the following command to correct them:

```bash
chmod 600 /etc/ssh/sshd_config
```

## Default Value

The default ownership is root:root with permissions 644.

## References

- CIS Ubuntu 12.04 LTS Server Benchmark v1.1.0

## Profile

Level 1 - Scored
