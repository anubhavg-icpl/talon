# stage: report
# category: CIS_benchmarks


# 7.6 Deactivate Wireless Interfaces (Not Scored)

## Profile Applicability

- Level 1

## Description

Wireless networking is used when wired networks are unavailable. Ubuntu provides the nmcli interface which allows system administrators to configure and use wireless networks.

## Rationale

If wireless is not to be used, wireless devices can be disabled to reduce the potential attack surface.

## Audit Procedure

### Using Command Line

Perform the following to determine if wireless interfaces are active.

```bash
ifconfig -a
```

Validate that all interfaces using wireless are down.

## Expected Result

No active wireless interfaces should be present.

## Remediation

### Using Command Line

Use the following command to disable wireless:

```bash
nmcli nm wifi off
```

## Default Value

Wireless interfaces may be active by default if hardware is present.

## References

- CIS Ubuntu 12.04 LTS Server Benchmark v1.1.0

## Profile

Level 1 - Not Scored
