# stage: report
# category: CIS_benchmarks


# 5.1.6 Ensure permissions on /etc/cron.monthly are configured (Scored)

## Profile Applicability

- Level 1 - Server
- Level 1 - Workstation

## Description

The `/etc/cron.monthly` directory contains system cron jobs that need to run on a monthly basis. The files in this directory cannot be manipulated by the `crontab` command, but are instead edited by system administrators using a text editor. The commands below restrict read/write and search access to user and group root, preventing regular users from accessing this directory.

## Rationale

Granting write access to this directory for non-privileged users could provide them the means for gaining unauthorized elevated privileges. Granting read access to this directory could give an unprivileged user insight in how to gain elevated privileges or circumvent auditing controls.

## Audit Procedure

Run the following command and verify `Uid` and `Gid` are both `0/root` and `Access` does not grant permissions to `group` or `other`:

```bash
stat /etc/cron.monthly
```

## Expected Result

```
Access: (0700/-rw-------)  Uid: (    0/    root)   Gid: (    0/    root)
```

## Remediation

Run the following commands to set ownership and permissions on `/etc/cron.monthly`:

```bash
chown root:root /etc/cron.monthly
chmod og-rwx /etc/cron.monthly
```

## Default Value

Not configured by default.

## References

- CIS Controls: 5.1 - Minimize And Sparingly Use Administrative Privileges

## Profile

- Level 1 - Server
- Level 1 - Workstation
