# stage: report
# category: CIS_benchmarks


# 2.2.8 Ensure DNS Server is not enabled (Scored)

## Profile Applicability

- Level 1 - Server
- Level 1 - Workstation

## Description

The Domain Name System (DNS) is a hierarchical naming system that maps names to IP addresses for computers, services and other resources connected to a network.

## Rationale

Unless a system is specifically designated to act as a DNS server, it is recommended that the package be deleted to reduce the potential attack surface.

## Audit Procedure

Run the following to ensure no start links for `bind9` exist in `/etc/rc*.d`:

```bash
ls /etc/rc*.d/S*bind9
```

No results should be returned.

## Expected Result

No output should be returned, indicating that `bind9` has no start links.

## Remediation

Run the following command to disable `bind9`:

```bash
update-rc.d bind9 disable
```

## Default Value

DNS server is not enabled by default.

## References

- CIS Controls: 9.1 Limit Open Ports, Protocols, and Services

## Profile

- Level 1 - Server
- Level 1 - Workstation
