# stage: report
# category: CIS_benchmarks


# 6.4 Ensure DHCP Server is not enabled (Scored)

## Profile Applicability

- Level 1

## Description

The Dynamic Host Configuration Protocol (DHCP) is a service that allows machines to be dynamically assigned IP addresses.

## Rationale

Unless a server is specifically set up to act as a DHCP server, it is recommended that this service be deleted to reduce the potential attack surface.

## Audit Procedure

### Using Command Line

Ensure no start conditions listed for `isc-dhcp-server` or `isc-dhcp-server6`:

```bash
initctl show-config isc-dhcp-server isc-dhcp-server
initctl show-config isc-dhcp-server6 isc-dhcp-server6
```

## Expected Result

No start conditions should be listed for isc-dhcp-server or isc-dhcp-server6.

## Remediation

### Using Command Line

Remove or comment out start lines in `/etc/init/isc-dhcp-server.conf` and `/etc/init/isc-dhcp-server6.conf`:

```bash
sed -i 's/^start/#start/' /etc/init/isc-dhcp-server.conf
sed -i 's/^start/#start/' /etc/init/isc-dhcp-server6.conf
```

## Default Value

Not enabled by default on Ubuntu 12.04 LTS Server.

## References

- CIS Ubuntu 12.04 LTS Server Benchmark v1.1.0
- http://www.isc.org/software/dhcp

## Profile

Level 1 - Scored
