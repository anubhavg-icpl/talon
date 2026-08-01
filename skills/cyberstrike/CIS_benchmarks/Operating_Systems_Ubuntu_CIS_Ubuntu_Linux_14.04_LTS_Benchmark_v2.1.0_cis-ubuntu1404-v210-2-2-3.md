# stage: report
# category: CIS_benchmarks


# 2.2.3 Ensure Avahi Server is not enabled (Scored)

## Profile Applicability

- Level 1 - Server
- Level 1 - Workstation

## Description

Avahi is a free zeroconf implementation, including a system for multicast DNS/DNS-SD service discovery. Avahi allows programs to publish and discover services and hosts running on a local network with no specific configuration. For example, a user can plug a computer into a network and Avahi automatically finds printers to print to, files to look at and people to talk to, as well as network services running on the machine.

## Rationale

Automatic discovery of network services is not normally required for system functionality. It is recommended to disable the service to reduce the potential attack surface.

## Audit Procedure

Run the following commands to verify no start conditions listed for `avahi-daemon`:

```bash
initctl show-config avahi-daemon
```

Verify the output shows `avahi-daemon` with no start conditions.

## Expected Result

The `avahi-daemon` service should have no start conditions listed.

## Remediation

Remove or comment out start lines in `/etc/init/avahi-daemon.conf`:

```bash
#start on runlevel [2345]
```

## Default Value

Avahi may be enabled by default on desktop installations.

## References

- CIS Controls: 9.1 Limit Open Ports, Protocols, and Services

## Profile

- Level 1 - Server
- Level 1 - Workstation
