# stage: report
# category: CIS_benchmarks


# CIS Ubuntu Linux 18.04 LTS Benchmark v2.2.0 - Control 4.2.4

## Profile Applicability

- Level 1 - Server
- Level 1 - Workstation

## Description

The `Banner` parameter specifies a file whose contents must be sent to the remote user before authentication is permitted. By default, no banner is displayed.

## Rationale

Banners are used to warn connecting users of the particular site's policy regarding connection. Presenting a warning message prior to the normal user login may assist the prosecution of trespassers on the computer system.

## Audit Procedure

### Command Line

Run the following command and verify that output matches:

```bash
sshd -T | grep -i banner
```

### Expected Result

```
banner /etc/issue.net
```

## Remediation

### Command Line

Edit the `/etc/ssh/sshd_config` file to set the parameter as follows:

```bash
Banner /etc/issue.net
```

## Default Value

Banner none

## References

1. NIST SP 800-53 Rev. 5: AC-8

## CIS Controls

None

## Profile Applicability

- Level 1 - Server
- Level 1 - Workstation

## Assessment Status

Automated
