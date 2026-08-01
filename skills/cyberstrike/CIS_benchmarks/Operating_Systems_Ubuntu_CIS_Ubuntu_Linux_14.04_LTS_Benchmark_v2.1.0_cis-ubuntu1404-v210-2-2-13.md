# stage: report
# category: CIS_benchmarks


# 2.2.13 Ensure HTTP Proxy Server is not enabled (Scored)

## Profile Applicability

- Level 1 - Server
- Level 1 - Workstation

## Description

Squid is a standard proxy server used in many distributions and environments.

## Rationale

If there is no need for a proxy server, it is recommended that the squid proxy be deleted to reduce the potential attack surface.

## Audit Procedure

Ensure no start conditions listed for `squid3`:

```bash
initctl show-config squid3
```

Verify the output shows `squid3` with no start conditions.

## Expected Result

The `squid3` service should have no start conditions listed.

## Remediation

Remove or comment out start lines in `/etc/init/squid3.conf`:

```bash
#start on runlevel [2345]
```

## Default Value

Squid proxy server is not enabled by default.

## References

- CIS Controls: 9.1 Limit Open Ports, Protocols, and Services

## Profile

- Level 1 - Server
- Level 1 - Workstation
