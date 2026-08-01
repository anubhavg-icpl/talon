# stage: report
# category: CIS_benchmarks


# 2.1.10 Ensure xinetd is not enabled (Scored)

## Profile Applicability

- Level 1 - Server
- Level 1 - Workstation

## Description

The eXtended InterNET Daemon (`xinetd`) is an open source super daemon that replaced the original `inetd` daemon. The `xinetd` daemon listens for well known services and dispatches the appropriate daemon to properly respond to service requests.

## Rationale

If there are no `xinetd` services required, it is recommended that the daemon be disabled.

## Audit Procedure

Run the following commands to verify no start conditions listed for `xinetd`:

```bash
initctl show-config xinetd
```

Verify the output shows `xinetd` with no start conditions.

## Expected Result

The `xinetd` service should have no start conditions listed.

## Remediation

Remove or comment out start lines in `/etc/init/xinetd.conf`:

```bash
#start on runlevel [2345]
```

## Default Value

xinetd is not enabled by default on Ubuntu 14.04.

## References

- CIS Controls: 9.1 Limit Open Ports, Protocols, and Services

## Profile

- Level 1 - Server
- Level 1 - Workstation
