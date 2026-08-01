# stage: report
# category: CIS_benchmarks


# CIS Ubuntu Linux 16.04 LTS Benchmark v2.0.0 - 1.7.5

## Profile Applicability

- Level 1 - Server
- Level 1 - Workstation

## Description

The contents of the `/etc/issue` file are displayed to users prior to login for local terminals.

## Rationale

If the `/etc/issue` file does not have the correct ownership it could be modified by unauthorized users with incorrect or misleading information.

## Audit Procedure

### Command Line

Run the following command and verify `Uid` and `Gid` are both `0/root` and `Access` is `644`:

```bash
stat /etc/issue
```

Expected output:

```
Access: (0644/-rw-r--r--) Uid: ( 0/ root) Gid: ( 0/ root)
```

## Expected Result

`/etc/issue` should be owned by root:root with permissions 644.

## Remediation

### Command Line

Run the following commands to set permissions on `/etc/issue`:

```bash
chown root:root /etc/issue
chmod u-x,go-wx /etc/issue
```

## Default Value

Access: (0644/-rw-r--r--) Uid: ( 0/ root) Gid: ( 0/ root)

## References

None.

## CIS Controls

| Controls Version | Control                             |
| ---------------- | ----------------------------------- |
| v7               | 5.1 Establish Secure Configurations |

## Assessment Status

Automated
