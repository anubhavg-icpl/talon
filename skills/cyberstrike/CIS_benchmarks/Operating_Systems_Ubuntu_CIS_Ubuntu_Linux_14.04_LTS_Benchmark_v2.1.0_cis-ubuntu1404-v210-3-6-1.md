# stage: report
# category: CIS_benchmarks


# 3.6.1 Ensure iptables is installed (Scored)

## Profile Applicability

- Level 1 - Server
- Level 1 - Workstation

## Description

`iptables` allows configuration of the IPv4 tables in the linux kernel and the rules stored within them. Most firewall configuration utilities operate as a front end to `iptables`.

## Rationale

iptables is required for firewall management and configuration.

## Audit Procedure

Run the following command and verify iptables is installed:

```bash
dpkg -s iptables
```

## Expected Result

The package `iptables` should be installed (status: install ok installed).

## Remediation

Run the following command to install `iptables`:

```bash
apt-get install iptables
```

## Default Value

Installed by default on Ubuntu.

## References

- CIS Controls: 9.2 - Leverage Host-based Firewalls

## Profile

- Level 1 - Server
- Level 1 - Workstation
