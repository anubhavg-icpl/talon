# stage: report
# category: CIS_benchmarks


# 6.13 Ensure HTTP Proxy Server is not enabled (Not Scored)

## Profile Applicability

- Level 1

## Description

Squid is a standard proxy server used in many distributions and environments.

## Rationale

If there is no need for a proxy server, it is recommended that the squid proxy be deleted to reduce the potential attack surface.

## Audit Procedure

### Using Command Line

Ensure no start conditions listed for `squid3`:

```bash
initctl show-config squid3 squid3
```

## Expected Result

No start conditions should be listed for squid3.

## Remediation

### Using Command Line

Remove or comment out start lines in `/etc/init/squid3.conf`:

```bash
sed -i 's/^start/#start/' /etc/init/squid3.conf
```

## Default Value

Not installed by default on Ubuntu 12.04 LTS Server.

## References

- CIS Ubuntu 12.04 LTS Server Benchmark v1.1.0

## Profile

Level 1 - Not Scored
