# stage: report
# category: CIS_benchmarks


# CIS Ubuntu Linux 16.04 LTS Benchmark v2.0.0 - 1.7.4

## Profile Applicability

- Level 1 - Server
- Level 1 - Workstation

## Description

The contents of the `/etc/motd` file are displayed to users after login and function as a message of the day for authenticated users.

## Rationale

If the `/etc/motd` file does not have the correct ownership it could be modified by unauthorized users with incorrect or misleading information.

## Audit Procedure

### Command Line

Run the following command and verify `Uid` and `Gid` are both `0/root` and `Access` is `644`, or the file doesn't exist:

```bash
stat /etc/motd
```

Expected output:

```
Access: (0644/-rw-r--r--) Uid: ( 0/ root) Gid: ( 0/ root)
```

OR

```
stat: cannot stat '/etc/motd': No such file or directory
```

## Expected Result

`/etc/motd` should be owned by root:root with permissions 644, or the file should not exist.

## Remediation

### Command Line

Run the following commands to set permissions on `/etc/motd`:

```bash
chown root:root /etc/motd
chmod u-x,go-wx /etc/motd
```

**OR** run the following command to remove the `/etc/motd` file:

```bash
rm /etc/motd
```

## Additional Information

If Message of the day is not needing, this file can be removed.

## Default Value

File doesn't exist.

## References

None.

## CIS Controls

| Controls Version | Control                             |
| ---------------- | ----------------------------------- |
| v7               | 5.1 Establish Secure Configurations |

## Assessment Status

Automated
