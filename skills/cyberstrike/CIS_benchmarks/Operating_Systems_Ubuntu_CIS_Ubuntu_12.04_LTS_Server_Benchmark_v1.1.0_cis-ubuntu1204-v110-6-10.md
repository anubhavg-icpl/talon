# stage: report
# category: CIS_benchmarks


# 6.10 Ensure HTTP Server is not enabled (Not Scored)

## Profile Applicability

- Level 1

## Description

HTTP or web servers provide the ability to host web site content.

## Rationale

Unless there is a need to run the system as a web server, it is recommended that the package be deleted to reduce the potential attack surface.

## Audit Procedure

### Using Command Line

Run the following to ensure no start links for `apache2` exist in `/etc/rc*.d`:

```bash
ls /etc/rc*.d/S*apache2
```

## Expected Result

No results should be returned.

## Remediation

### Using Command Line

Remove any start links for `apache2` from `/etc/rc*.d`:

```bash
rm /etc/rc*.d/S*apache2
```

## Default Value

Not installed by default on Ubuntu 12.04 LTS Server.

## References

- CIS Ubuntu 12.04 LTS Server Benchmark v1.1.0

## Profile

Level 1 - Not Scored
