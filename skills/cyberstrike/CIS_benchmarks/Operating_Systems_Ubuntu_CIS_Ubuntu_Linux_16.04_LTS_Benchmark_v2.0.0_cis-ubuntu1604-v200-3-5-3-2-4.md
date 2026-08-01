# stage: report
# category: CIS_benchmarks


# CIS Ubuntu Linux 16.04 LTS Benchmark v2.0.0 - Control 3.5.3.2.4

## Profile

- **Level:** Level 1 - Server, Level 1 - Workstation
- **Assessment Status:** Automated

## Description

Any ports that have been opened on non-loopback addresses need firewall rules to govern traffic.

Note:

- Changing firewall settings while connected over network can result in being locked out of the system.
- Remediation will only affect the active system firewall, be sure to configure the default policy in your firewall management to apply on boot as well.
- The remediation command opens up the port to traffic from all sources. Consult iptables documentation and set any restrictions in compliance with site policy.

## Rationale

Without a firewall rule configured for open ports default firewall policy will drop all packets to these ports.

## Audit Procedure

### Command Line

Run the following command to determine open ports:

```bash
ss -4tuln
```

Run the following command to determine firewall rules:

```bash
iptables -L INPUT -v -n
```

Verify all open ports listening on non-localhost addresses have at least one firewall rule.

## Remediation

### Command Line

For each port identified in the audit which does not have a firewall rule establish a proper rule for accepting inbound connections:

```bash
iptables -A INPUT -p <protocol> --dport <port> -m state --state NEW -j ACCEPT
```

## References

None

## CIS Controls

Version 7

9.2 Ensure Only Approved Ports, Protocols and Services Are Running - Ensure that only network ports, protocols, and services listening on a system with validated business needs, are running on each system.

9.4 Apply Host-based Firewalls or Port Filtering - Apply host-based firewalls or port filtering tools on end systems, with a default-deny rule that drops all traffic except those services and ports that are explicitly allowed.
