# stage: report
# category: CIS_benchmarks


# 12.6 Verify User/Group Ownership on /etc/group (Scored)

## Profile Applicability

- Level 1

## Description

The `/etc/group` file contains a list of all the valid groups defined in the system. The command below allows read/write access for root and read access for everyone else.

## Rationale

The `/etc/group` file needs to be protected from unauthorized changes by non-privileged users, but needs to be readable as this information is used with many non-privileged programs.

## Audit Procedure

### Using Command Line

Run the following command to determine the permissions on the `/etc/group` file:

```bash
/bin/ls -l /etc/group
```

## Expected Result

```
-rw-r--r-- 1 root root <size> <date> /etc/group
```

Owner should be `root` and group should be `root`.

## Remediation

### Using Command Line

If the ownership of the `/etc/group` file are incorrect, run the following command to correct them:

```bash
/bin/chown root:root /etc/group
```

## Default Value

root:root

## References

- CIS Ubuntu 12.04 LTS Server Benchmark v1.1.0

## Profile

Level 1 - Scored
