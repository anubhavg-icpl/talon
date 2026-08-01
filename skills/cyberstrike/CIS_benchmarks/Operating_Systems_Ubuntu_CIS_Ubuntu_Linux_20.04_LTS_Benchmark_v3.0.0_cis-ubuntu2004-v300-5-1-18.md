# stage: report
# category: CIS_benchmarks


# Ensure sshd MaxStartups is configured (Automated)

## Profile Applicability

- Level 1 - Server
- Level 1 - Workstation

## Description

The `MaxStartups` parameter specifies the maximum number of concurrent unauthenticated connections to the SSH daemon.

## Rationale

To protect a system from denial of service due to a large number of pending authentication connection attempts, use the rate limiting function of MaxStartups to protect availability of sshd logins and prevent overwhelming the daemon.

## Audit Procedure

### Command Line

Run the following command to verify `MaxStartups` is `10:30:60` or more restrictive:

```bash
# sshd -T | awk '$1 ~ /^\s*maxstartups/{split($2, a, ":");if(a[1] > 10 || a[2] > 30 || a[3] > 60) print $0}'
```

Nothing should be returned.

## Expected Result

Nothing should be returned (MaxStartups is 10:30:60 or more restrictive).

## Remediation

### Command Line

Edit the `/etc/ssh/sshd_config` file to set the `MaxStartups` parameter to `10:30:60` or more restrictive above any `Include` entries as follows:

```
MaxStartups 10:30:60
```

Note: First occurrence of a option takes precedence. If Include locations are enabled, used, and order of precedence is understood in your environment, the entry may be created in a file in Include location.

## Default Value

MaxStartups 10:30:100

## References

1. SSHD_CONFIG(5)
2. NIST SP 800-53 Rev. 5: CM-1, CM-2, CM-6, CM-7, IA-5

## CIS Controls

v8 - 0.0 Explicitly Not Mapped

v7 - 0.0 Explicitly Not Mapped

MITRE ATT&CK Mappings: T1499, T1499.002 | TA0040
