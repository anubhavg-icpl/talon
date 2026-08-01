# stage: report
# category: CIS_benchmarks


# 1.5.1.6 Ensure permissions on /etc/issue.net are configured (Automated)

## Description

The contents of the `/etc/issue.net` file are displayed to users prior to login for remote connections from configured services.

## Rationale

If the `/etc/issue.net` file does not have the correct ownership it could be modified by unauthorized users with incorrect or misleading information.

## Audit Procedure

Run the following command and verify `Uid` and `Gid` are both `0/root` and `Access` is `644`:

```bash
# stat /etc/issue.net
Access: (0644/-rw-r--r--)  Uid: (    0/    root)   Gid: (    0/    root)
```

## Expected Result

The output should show `Access: (0644/-rw-r--r--)`, `Uid: ( 0/ root)`, and `Gid: ( 0/ root)`.

## Remediation

Run the following commands to set permissions on `/etc/issue.net`:

```bash
# chown root:root /etc/issue.net
# chmod 644 /etc/issue.net
```

`/etc` is stateless on Container-Optimized OS. Therefore, the steps mentioned above needs to be performed after every boot.

## CIS Controls

| Controls Version | Control                                 | IG 1 | IG 2 | IG 3 |
| ---------------- | --------------------------------------- | ---- | ---- | ---- |
| v8               | 3.3 Configure Data Access Control Lists | x    | x    | x    |
| v7               | 5.1 Establish Secure Configurations     | x    | x    | x    |

## Profile

Level 2 - Server | Automated
