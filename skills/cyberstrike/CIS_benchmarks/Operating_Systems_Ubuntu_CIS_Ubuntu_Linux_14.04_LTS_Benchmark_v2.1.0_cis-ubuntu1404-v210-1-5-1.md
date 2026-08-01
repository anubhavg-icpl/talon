# stage: report
# category: CIS_benchmarks


# 1.5.1 Ensure core dumps are restricted (Scored)

## Profile Applicability

- Level 1 - Server
- Level 1 - Workstation

## Description

A core dump is the memory of an executable program. It is generally used to determine why a program aborted. It can also be used to glean confidential information from a core file. The system provides the ability to set a soft limit for core dumps, but this can be overridden by the user.

## Rationale

Setting a hard limit on core dumps prevents users from overriding the soft variable. If core dumps are required, consider setting limits for user groups (see `limits.conf(5)`). In addition, setting the `fs.suid_dumpable` variable to 0 will prevent setuid programs from dumping core.

## Audit Procedure

Run the following commands and verify output matches:

```bash
grep "hard core" /etc/security/limits.conf /etc/security/limits.d/*
# Expected: * hard core 0
sysctl fs.suid_dumpable
# Expected: fs.suid_dumpable = 0
grep "fs\.suid_dumpable" /etc/sysctl.conf /etc/sysctl.d/*
# Expected: fs.suid_dumpable = 0
```

## Expected Result

```
* hard core 0
fs.suid_dumpable = 0
fs.suid_dumpable = 0
```

## Remediation

Add the following line to `/etc/security/limits.conf` or a `/etc/security/limits.d/*` file:

```
* hard core 0
```

Set the following parameter in `/etc/sysctl.conf` or a `/etc/sysctl.d/*` file:

```
fs.suid_dumpable = 0
```

Run the following command to set the active kernel parameter:

```bash
sysctl -w fs.suid_dumpable=0
```

## Default Value

Not applicable.

## Notes

It has been reported that due to Ubuntu bug #50093 this setting (and some others) can fail to apply properly on reboot requiring it to be manually re-applied. One method of accomplishing this is to add `sysctl -p` to run on reboot to your systems crontab.

## References

- CIS Controls: 13 Data Protection

## Profile

- Level 1 - Server
- Level 1 - Workstation
