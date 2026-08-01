# stage: report
# category: CIS_benchmarks


# 7.4.1 Install TCP Wrappers (Scored)

## Profile Applicability

- Level 1

## Description

TCP Wrappers provides a simple access list and standardized logging method for services capable of supporting it. In the past, services that were called from `inetd` and `xinetd` supported the use of tcp wrappers. As `inetd` and `xinetd` have been falling in disuse, any service that can support tcp wrappers will have the `libwrap.so` library attached to it.

## Rationale

TCP Wrappers provide a good simple access list mechanism to services that may not have that support built in. It is recommended that all services that can support TCP Wrappers, use it.

## Audit Procedure

### Using Command Line

Run the following to ensure `tcpd` is installed:

```bash
dpkg -s tcpd
```

Ensure package status is `installed ok installed`.

To verify if a service supports TCP Wrappers, run the following command:

```bash
ldd <path-to-daemon> | grep libwrap.so
```

If there is any output, then the service supports TCP Wrappers.

## Expected Result

```
Status: install ok installed
```

## Remediation

### Using Command Line

Install `tcpd`:

```bash
apt-get install tcpd
```

## Default Value

TCP Wrappers is not installed by default.

## References

- CIS Ubuntu 12.04 LTS Server Benchmark v1.1.0

## Profile

Level 1 - Scored
