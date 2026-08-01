# stage: report
# category: CIS_benchmarks


# 2.2.7 Ensure DNS Server is not installed (Automated)

## Profile

- Level 1 - Server
- Level 1 - Workstation

## Description

The Domain Name System (DNS) is a hierarchical naming system that maps names to IP addresses for computers, services and other resources connected to a network.

## Rationale

Unless a system is specifically designated to act as a DNS server, it is recommended that the package be deleted to reduce the potential attack surface.

## Audit Procedure

### Command Line

Run the following command to verify `DNS server` is not installed:

```bash
# dpkg-query -s bind9 &>/dev/null && echo "bind9 is installed"
```

Nothing should be returned.

## Expected Result

No output (empty result).

## Remediation

### Command Line

Run the following commands to disable `DNS server`:

```bash
# apt purge bind9
```

## References

1. NIST SP 800-53 Rev. 5: CM-7

## CIS Controls

- v8: 4.8 - Uninstall or Disable Unnecessary Services on Enterprise Assets and Software
- v7: 9.2 - Ensure Only Approved Ports, Protocols and Services Are Running
