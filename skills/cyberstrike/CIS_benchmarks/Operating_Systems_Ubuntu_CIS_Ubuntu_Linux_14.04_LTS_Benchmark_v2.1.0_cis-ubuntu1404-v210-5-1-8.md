# stage: report
# category: CIS_benchmarks


# 5.1.8 Ensure at/cron is restricted to authorized users (Scored)

## Profile Applicability

- Level 1 - Server
- Level 1 - Workstation

## Description

Configure `/etc/cron.allow` and `/etc/at.allow` to allow specific users to use these services. If `/etc/cron.allow` or `/etc/at.allow` do not exist, then `/etc/at.deny` and `/etc/cron.deny` are checked. Any user not specifically defined in those files is allowed to use at and cron. By removing the files, only users in `/etc/cron.allow` and `/etc/at.allow` are allowed to use at and cron. Note that even though a given user is not listed in `cron.allow`, cron jobs can still be run as that user. The `cron.allow` file only controls administrative access to the crontab command for scheduling and modifying cron jobs.

## Rationale

On many systems, only the system administrator is authorized to schedule `cron` jobs. Using the `cron.allow` file to control who can run `cron` jobs enforces this policy. It is easier to manage an allow list than a deny list. In a deny list, you could potentially add a user ID to the system and forget to add it to the deny files.

## Audit Procedure

Run the following commands and ensure `/etc/cron.deny` and `/etc/at.deny` do not exist:

```bash
stat /etc/cron.deny
stat /etc/at.deny
```

Run the following command and verify `Uid` and `Gid` are both `0/root` and `Access` does not grant permissions to `group` or `other` for both `/etc/cron.allow` and `/etc/at.allow`:

```bash
stat /etc/cron.allow
stat /etc/at.allow
```

## Expected Result

```
stat: cannot stat '/etc/cron.deny': No such file or directory
stat: cannot stat '/etc/at.deny': No such file or directory

Access: (0600/-rw-------)  Uid: (    0/    root)   Gid: (    0/    root)  # for cron.allow
Access: (0600/-rw-------)  Uid: (    0/    root)   Gid: (    0/    root)  # for at.allow
```

## Remediation

Run the following commands to remove `/etc/cron.deny` and `/etc/at.deny` and create and set permissions and ownership for `/etc/cron.allow` and `/etc/at.allow`:

```bash
rm /etc/cron.deny
rm /etc/at.deny
touch /etc/cron.allow
touch /etc/at.allow
chmod og-rwx /etc/cron.allow
chmod og-rwx /etc/at.allow
chown root:root /etc/cron.allow
chown root:root /etc/at.allow
```

## Default Value

Not configured by default.

## References

- CIS Controls: 16 - Account Monitoring and Control

## Profile

- Level 1 - Server
- Level 1 - Workstation
