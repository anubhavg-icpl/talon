# stage: report
# category: CIS_benchmarks


# 6.1.5 Ensure permissions on /etc/gshadow are configured (Scored)

## Profile Applicability

- Level 1 - Server
- Level 1 - Workstation

## Description

The `/etc/gshadow` file is used to store the information about groups that is critical to the security of those accounts, such as the hashed password and other security information.

## Rationale

If attackers can gain read access to the `/etc/gshadow` file, they can easily run a password cracking program against the hashed password to break it. Other security information that is stored in the `/etc/gshadow` file (such as group administrators) could also be useful to subvert the group.

## Audit Procedure

Run the following command and verify `Uid` is `0/root`, `Gid` is `<gid>/shadow`, and `Access` is `640` or more restrictive:

```bash
stat /etc/gshadow
```

## Expected Result

```
Access: (0640/-rw-r-----)  Uid: (    0/    root)   Gid: (   42/  shadow)
```

## Remediation

Run the following commands to set permissions on `/etc/gshadow`:

```bash
chown root:shadow /etc/gshadow
chmod o-rwx,g-rw /etc/gshadow
```

## Default Value

Access: (0640/-rw-r-----) Uid: ( 0/ root) Gid: ( 42/ shadow)

## References

None

## CIS Controls

16.14 Encrypt/Hash All Authentication Files And Monitor Their Access - Verify that all authentication files are encrypted or hashed and that these files cannot be accessed without root or administrator privileges. Audit all access to password files in the system.

## Profile

- Level 1
