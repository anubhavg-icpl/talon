# stage: report
# category: CIS_benchmarks


# 5.3 Ensure daytime is not enabled (Scored)

## Profile Applicability

- Level 1

## Description

`daytime` is a network service that responds with the server's current date and time. This service is intended for debugging and testing purposes. It is recommended that this service be disabled.

## Rationale

Disabling this service will reduce the remote attack surface of the system.

## Audit Procedure

### Using Command Line

Ensure the `daytime` services are not enabled:

```bash
grep ^daytime /etc/inetd.conf
```

## Expected Result

No results should be returned.

## Remediation

### Using Command Line

Remove or comment out any `daytime` lines in `/etc/inetd.conf`:

```bash
sed -i 's/^daytime/#daytime/' /etc/inetd.conf
```

## Default Value

Not enabled by default on Ubuntu 12.04 LTS Server.

## References

- CIS Ubuntu 12.04 LTS Server Benchmark v1.1.0

## Profile

Level 1 - Scored
