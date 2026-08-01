# stage: report
# category: CIS_benchmarks


# 2.1.11 Ensure openbsd-inetd is not installed (Scored)

## Profile Applicability

- Level 1 - Server
- Level 1 - Workstation

## Description

The `inetd` daemon listens for well known services and dispatches the appropriate daemon to properly respond to service requests.

## Rationale

If there are no `inetd` services required, it is recommended that the daemon be removed.

## Audit Procedure

Run the following command and verify `openbsd-inetd` is not installed:

```bash
dpkg -s openbsd-inetd
```

## Expected Result

The command should indicate that the package is not installed (e.g., `dpkg-query: package 'openbsd-inetd' is not installed`).

## Remediation

Run the following command to uninstall `openbsd-inetd`:

```bash
apt-get remove openbsd-inetd
```

## Default Value

openbsd-inetd is not installed by default on Ubuntu 14.04.

## References

- CIS Controls: 9.1 Limit Open Ports, Protocols, and Services

## Profile

- Level 1 - Server
- Level 1 - Workstation
