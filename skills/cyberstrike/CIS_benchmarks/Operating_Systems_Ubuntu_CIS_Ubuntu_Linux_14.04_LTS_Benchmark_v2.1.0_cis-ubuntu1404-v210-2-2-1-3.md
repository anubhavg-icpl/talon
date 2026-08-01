# stage: report
# category: CIS_benchmarks


# 2.2.1.3 Ensure chrony is configured (Scored)

## Profile Applicability

- Level 1 - Server
- Level 1 - Workstation

## Description

`chrony` is a daemon which implements the Network Time Protocol (NTP) is designed to synchronize system clocks across a variety of systems and use a source that is highly accurate. More information on `chrony` can be found at http://chrony.tuxfamily.org/. `chrony` can be configured to be a client and/or a server.

This recommendation only applies if chrony is in use on the system.

## Rationale

If chrony is in use on the system proper configuration is vital to ensuring time synchronization is working properly.

## Audit Procedure

Run the following command and verify remote server is configured properly:

```bash
grep "^(server|pool)" /etc/chrony/chrony.conf
```

Expected output:

```
server <remote-server>
```

Multiple servers may be configured.

## Expected Result

At least one remote server or pool should be configured in `/etc/chrony/chrony.conf`.

## Remediation

Add or edit server or pool lines to `/etc/chrony/chrony.conf` as appropriate:

```
server <remote-server>
```

## Default Value

chrony server settings are not configured by default.

## References

- http://chrony.tuxfamily.org/
- CIS Controls: 6.1 Use At Least Two Synchronized Time Sources For All Servers And Network Equipment

## Profile

- Level 1 - Server
- Level 1 - Workstation
