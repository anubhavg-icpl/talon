# stage: report
# category: CIS_benchmarks


# 3.3.3 Ensure IPv6 is disabled (Not Scored)

## Profile Applicability

- Level 1 - Server
- Level 1 - Workstation

## Description

Although IPv6 has many advantages over IPv4, few organizations have implemented IPv6.

## Rationale

If IPv6 is not to be used, it is recommended that it be disabled to reduce the attack surface of the system.

## Audit Procedure

Run the following command and verify that each linux line has the `ipv6.disable=1` parameter set:

```bash
grep "^\s*linux" /boot/grub/grub.cfg
```

## Expected Result

Each linux line should contain `ipv6.disable=1`.

## Remediation

Edit `/etc/default/grub` and add `ipv6.disable=1` to GRUB_CMDLINE_LINUX:

```
GRUB_CMDLINE_LINUX="ipv6.disable=1"
```

Run the following command to update the grub2 configuration:

```bash
update-grub
```

## Default Value

IPv6 is enabled by default.

## References

- CIS Controls: 3 - Secure Configurations for Hardware and Software on Mobile Devices, Laptops, Workstations, and Servers
- CIS Controls: 11 - Secure Configurations for Network Devices such as Firewalls, Routers and switches

## Profile

- Level 1 - Server
- Level 1 - Workstation
