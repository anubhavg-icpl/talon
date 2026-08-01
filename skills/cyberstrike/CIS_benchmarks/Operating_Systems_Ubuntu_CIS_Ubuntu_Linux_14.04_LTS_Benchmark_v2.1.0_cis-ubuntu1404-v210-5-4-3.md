# stage: report
# category: CIS_benchmarks


# 5.4.3 Ensure default group for the root account is GID 0 (Scored)

## Profile Applicability

- Level 1 - Server
- Level 1 - Workstation

## Description

The usermod command can be used to specify which group the root user belongs to. This affects permissions of files that are created by the root user.

## Rationale

Using GID 0 for the `root` account helps prevent `root`-owned files from accidentally becoming accessible to non-privileged users.

## Audit Procedure

Run the following command and verify the result is `0`:

```bash
grep "^root:" /etc/passwd | cut -f4 -d:
```

## Expected Result

```
0
```

## Remediation

Run the following command to set the `root` user default group to GID `0`:

```bash
usermod -g 0 root
```

## Default Value

GID 0

## References

- CIS Controls: 5 - Controlled Use of Administration Privileges

## Profile

- Level 1 - Server
- Level 1 - Workstation
