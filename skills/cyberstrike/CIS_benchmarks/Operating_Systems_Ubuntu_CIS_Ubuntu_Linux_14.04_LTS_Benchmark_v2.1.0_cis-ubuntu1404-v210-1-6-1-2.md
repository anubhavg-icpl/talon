# stage: report
# category: CIS_benchmarks


# 1.6.1.2 Ensure the SELinux state is enforcing (Scored)

## Profile Applicability

- Level 2 - Server
- Level 2 - Workstation

## Description

Set SELinux to enable when the system is booted.

## Rationale

SELinux must be enabled at boot time in to ensure that the controls it provides are in effect at all times.

## Audit Procedure

Run the following commands and ensure output matches:

```bash
grep SELINUX=enforcing /etc/selinux/config
sestatus
```

## Expected Result

```
SELINUX=enforcing

SELinux status: enabled
Current mode: enforcing
Mode from config file: enforcing
```

## Remediation

Edit the `/etc/selinux/config` file to set the SELINUX parameter:

```
SELINUX=enforcing
```

## Default Value

Not applicable.

## References

- CIS Controls: 14.4 Protect Information With Access Control Lists

## Profile

- Level 2 - Server
- Level 2 - Workstation
