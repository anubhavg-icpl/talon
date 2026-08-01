# stage: report
# category: CIS_benchmarks


# 5.1.4 Ensure talk server is not enabled (Scored)

## Profile Applicability

- Level 1

## Description

The talk software makes it possible for users to send and receive messages across systems through a terminal session. The talk client (allows initiate of talk sessions) is installed by default.

## Rationale

The software presents a security risk as it uses unencrypted protocols for communication.

## Audit Procedure

### Using Command Line

Ensure the `talk` services are not enabled:

```bash
grep ^talk /etc/inetd.conf
grep ^ntalk /etc/inetd.conf
```

## Expected Result

No results should be returned.

## Remediation

### Using Command Line

Remove or comment out any `talk` or `ntalk` lines in `/etc/inetd.conf`:

```bash
sed -i 's/^talk/#talk/' /etc/inetd.conf
sed -i 's/^ntalk/#ntalk/' /etc/inetd.conf
```

## Default Value

Not enabled by default on Ubuntu 12.04 LTS Server.

## References

- CIS Ubuntu 12.04 LTS Server Benchmark v1.1.0

## Profile

Level 1 - Scored
