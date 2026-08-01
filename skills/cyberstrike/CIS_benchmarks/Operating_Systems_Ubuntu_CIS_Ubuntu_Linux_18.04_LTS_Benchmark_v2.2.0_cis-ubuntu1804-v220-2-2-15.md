# stage: report
# category: CIS_benchmarks


# 2.2.15 Ensure dnsmasq is not installed (Automated)

## Profile

- Level 1 - Server
- Level 1 - Workstation

## Description

`dnsmasq` is a lightweight tool that provides DNS caching, DNS forwarding and DHCP (Dynamic Host Configuration Protocol) services.

## Rationale

Unless a system is specifically designated to act as a DNS caching, DNS forwarding and/or DHCP server, it is recommended that the package be removed to reduce the potential attack surface.

## Audit Procedure

### Command Line

Run one of the following commands to verify `dnsmasq` is not installed:

```bash
# dpkg-query -s dnsmasq &>/dev/null && echo "dnsmasq is installed"
```

Nothing should be returned.

## Expected Result

No output (empty result).

## Remediation

### Command Line

Run the following command to remove `dnsmasq`:

```bash
# apt purge dnsmasq
```

## References

1. NIST SP 800-53 Rev. 5: CM-7

## CIS Controls

- v8: 4.8 - Uninstall or Disable Unnecessary Services on Enterprise Assets and Software
- v7: 9.2 - Ensure Only Approved Ports, Protocols and Services Are Running
