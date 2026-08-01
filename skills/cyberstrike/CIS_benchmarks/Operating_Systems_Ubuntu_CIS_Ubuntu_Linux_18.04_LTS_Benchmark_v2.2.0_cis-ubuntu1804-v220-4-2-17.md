# stage: report
# category: CIS_benchmarks


# CIS Ubuntu Linux 18.04 LTS Benchmark v2.2.0 - Control 4.2.17

## Profile Applicability

- Level 1 - Server
- Level 1 - Workstation

## Description

The `MaxStartups` parameter specifies the maximum number of concurrent unauthenticated connections to the SSH daemon.

## Rationale

To protect a system from denial of service due to a large number of pending authentication connection attempts, use the rate limiting function of `MaxStartups` to protect availability of sshd logins and prevent overwhelming the daemon.

## Audit Procedure

### Command Line

Run the following command and verify the output:

```bash
sshd -T | grep -i maxstartups
```

### Expected Result

```
maxstartups 10:30:60
```

OR more restrictive.

## Remediation

### Command Line

Edit the `/etc/ssh/sshd_config` file to set the parameter as follows:

```bash
MaxStartups 10:30:60
```

## Default Value

MaxStartups 10:30:100

## References

1. NIST SP 800-53 Rev. 5: CM-6

## CIS Controls

None

## Profile Applicability

- Level 1 - Server
- Level 1 - Workstation

## Assessment Status

Automated
