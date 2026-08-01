# stage: report
# category: CIS_benchmarks


# 3.1.1 Ensure IP forwarding is disabled (Scored)

## Profile Applicability

- Level 1 - Server
- Level 1 - Workstation

## Description

The `net.ipv4.ip_forward` flag is used to tell the system whether it can forward packets or not.

## Rationale

Setting the flag to 0 ensures that a system with multiple interfaces (for example, a hard proxy), will never be able to forward packets, and therefore, never serve as a router.

## Audit Procedure

Run the following command and verify output matches:

```bash
sysctl net.ipv4.ip_forward
# Expected: net.ipv4.ip_forward = 0

grep "net\.ipv4\.ip_forward" /etc/sysctl.conf /etc/sysctl.d/*
# Expected: net.ipv4.ip_forward = 0
```

## Expected Result

```
net.ipv4.ip_forward = 0
```

## Remediation

Set the following parameter in `/etc/sysctl.conf` or a `/etc/sysctl.d/*` file:

```
net.ipv4.ip_forward = 0
```

Run the following commands to set the active kernel parameters:

```bash
sysctl -w net.ipv4.ip_forward=0
sysctl -w net.ipv4.route.flush=1
```

## Default Value

Not specified.

## References

- CIS Controls: 3 - Secure Configurations for Hardware and Software on Mobile Devices, Laptops, Workstations, and Servers
- CIS Controls: 11 - Secure Configurations for Network Devices such as Firewalls, Routers and switches

## Profile

- Level 1 - Server
- Level 1 - Workstation
