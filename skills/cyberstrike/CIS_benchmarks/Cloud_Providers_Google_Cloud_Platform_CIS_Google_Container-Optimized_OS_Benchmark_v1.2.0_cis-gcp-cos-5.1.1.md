# stage: report
# category: CIS_benchmarks


# 5.1.1 Ensure permissions on /etc/ssh/sshd_config are configured (Automated)

## Description

The `/etc/ssh/sshd_config` file contains configuration specifications for `sshd`. The command below sets the owner and group of the file to root.

## Rationale

The `/etc/ssh/sshd_config` file needs to be protected from unauthorized changes by non-privileged users.

## Audit Procedure

Run the following command and verify `Uid` and `Gid` are both `0/root` and `Access` does not grant permissions to `group` or `other`:

```bash
# stat /etc/ssh/sshd_config
Access: (0600/-rw-------)   Uid: (    0/    root)   Gid: (    0/    root)
```

## Expected Result

- `Uid` is `0/root`
- `Gid` is `0/root`
- Access is `0600/-rw-------` (no permissions for group or other)

## Remediation

Run the following commands to set ownership and permissions on `/etc/ssh/sshd_config`:

```bash
# chown root:root /etc/ssh/sshd_config
# chmod og-rwx /etc/ssh/sshd_config
```

## CIS Controls

| Controls Version | Control                                 | IG 1 | IG 2 | IG 3 |
| ---------------- | --------------------------------------- | ---- | ---- | ---- |
| v8               | 3.3 Configure Data Access Control Lists | x    | x    | x    |
| v7               | 5.1 Establish Secure Configurations     | x    | x    | x    |

## Profile

- Level 1 - Server
