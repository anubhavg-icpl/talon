# stage: report
# category: CIS_benchmarks


# 6.1.8 Ensure permissions on /etc/group- are configured (Scored)

## Profile Applicability

- Level 1 - Server
- Level 1 - Workstation

## Description

The `/etc/group-` file contains a backup list of all the valid groups defined in the system.

## Rationale

It is critical to ensure that the `/etc/group-` file is protected from unauthorized access. Although it is protected by default, the file permissions could be changed either inadvertently or through malicious actions.

## Audit Procedure

Run the following command and verify `Uid` and `Gid` are both `0/root` and `Access` is `644` or more restrictive:

```bash
stat /etc/group-
```

## Expected Result

```
Access: (0644/-rw-------)  Uid: (    0/    root)   Gid: (    0/    root)
```

## Remediation

Run the following command to set permissions on `/etc/group-`:

```bash
chown root:root /etc/group-
chmod u-x,go-wx /etc/group-
```

## Default Value

Access: (0644/-rw-------) Uid: ( 0/ root) Gid: ( 0/ root)

## References

None

## CIS Controls

16.14 Encrypt/Hash All Authentication Files And Monitor Their Access - Verify that all authentication files are encrypted or hashed and that these files cannot be accessed without root or administrator privileges. Audit all access to password files in the system.

## Profile

- Level 1
