# stage: report
# category: CIS_benchmarks


# 9.3.14 Set SSH Banner (Scored)

## Profile Applicability

- Level 1

## Description

The `Banner` parameter specifies a file whose contents must be sent to the remote user before authentication is permitted. By default, no banner is displayed.

## Rationale

Banners are used to warn connecting users of the particular site's policy regarding connection. Consult with your legal department for the appropriate warning banner for your site.

## Audit Procedure

### Using Command Line

To verify the correct SSH setting, run the following command and verify that `<bannerfile>` is either `/etc/issue` or `/etc/issue.net`:

```bash
grep "^Banner" /etc/ssh/sshd_config
```

## Expected Result

```
Banner /etc/issue.net
```

## Remediation

### Using Command Line

Edit the `/etc/ssh/sshd_config` file to set the parameter as follows:

```bash
Banner /etc/issue.net
```

## Default Value

No banner is configured by default.

## References

- CIS Ubuntu 12.04 LTS Server Benchmark v1.1.0

## Profile

Level 1 - Scored
