# stage: report
# category: CIS_benchmarks


# CIS Ubuntu Linux 18.04 LTS Benchmark v2.2.0 - Control 3.4.3.1.2

## Description

nftables is a subsystem of the Linux kernel providing filtering and classification of network packets/datagrams/frames and is the successor to iptables.

## Rationale

Running both `iptables` and `nftables` may lead to conflict.

## Impact

None.

## Audit Procedure

### Command Line

Run the following command to verify that `nftables` is not installed:

```bash
dpkg-query -s nftables &>/dev/null && echo "nftables is installed"
```

Nothing should be returned.

## Expected Result

No output should be returned (nftables is not installed).

## Remediation

### Command Line

Run the following command to remove `nftables`:

```bash
apt purge nftables
```

## Default Value

nftables is not installed by default.

## References

1. NIST SP 800-53 Rev. 5: CA-9, CM-7
2. CIS Ubuntu Linux 18.04 LTS Benchmark v2.2.0

## CIS Controls

Version 8

4.4 Implement and Manage a Firewall on Servers - Implement and manage a firewall on servers, where supported.

4.5 Implement and Manage a Firewall on End-User Devices - Implement and manage a host-based firewall or port-filtering tool on end-user devices.

Version 7

9.4 Apply Host-based Firewalls or Port Filtering - Apply host-based firewalls or port filtering tools on end systems, with a default-deny rule that drops all traffic except those services and ports that are explicitly allowed.

## Profile Applicability

- Level 1 - Server
- Level 1 - Workstation

## Assessment Status

Automated
