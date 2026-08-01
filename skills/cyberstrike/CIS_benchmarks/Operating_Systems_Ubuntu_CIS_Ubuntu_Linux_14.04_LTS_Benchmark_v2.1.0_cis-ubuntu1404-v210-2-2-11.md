# stage: report
# category: CIS_benchmarks


# 2.2.11 Ensure IMAP and POP3 server is not enabled (Scored)

## Profile Applicability

- Level 1 - Server
- Level 1 - Workstation

## Description

`dovecot` is an open source IMAP and POP3 server for Linux based systems.

## Rationale

Unless POP3 and/or IMAP servers are to be provided by this system, it is recommended that the service be deleted to reduce the potential attack surface.

## Audit Procedure

Run the following commands to verify no start conditions listed for `dovecot`:

```bash
initctl show-config dovecot
```

Verify the output shows `dovecot` with no start conditions.

## Expected Result

The `dovecot` service should have no start conditions listed.

## Remediation

Remove or comment out start lines in `/etc/init/dovecot.conf`:

```bash
#start on runlevel [2345]
```

**Note:** Several IMAP/POP3 servers exist and can use other service names. `exim` and `cyrus-imap` are example services that provide an HTTP server. These and other services should also be audited.

## Default Value

dovecot is not enabled by default.

## References

- CIS Controls: 9.1 Limit Open Ports, Protocols, and Services

## Profile

- Level 1 - Server
- Level 1 - Workstation
