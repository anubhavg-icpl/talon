# stage: report
# category: CIS_benchmarks


# 7.4.5 Verify Permissions on /etc/hosts.deny (Scored)

## Profile Applicability

- Level 1

## Description

The `/etc/hosts.deny` file contains network information that is used by many system applications and therefore must be readable for these applications to operate.

## Rationale

It is critical to ensure that the `/etc/hosts.deny` file is protected from unauthorized write access. Although it is protected by default, the file permissions could be changed either inadvertently or through malicious actions.

## Audit Procedure

### Using Command Line

Run the following command to determine the permissions on the `/etc/hosts.deny` file.

```bash
/bin/ls -l /etc/hosts.deny
```

## Expected Result

```
-rw-r--r-- 1 root root 2055 Jan 30 16:30 /etc/hosts.deny
```

## Remediation

### Using Command Line

If the permissions of the `/etc/hosts.deny` file are incorrect, run the following command to correct them:

```bash
/bin/chmod 644 /etc/hosts.deny
```

## Default Value

Permissions are 644 by default.

## References

- CIS Ubuntu 12.04 LTS Server Benchmark v1.1.0

## Profile

Level 1 - Scored
