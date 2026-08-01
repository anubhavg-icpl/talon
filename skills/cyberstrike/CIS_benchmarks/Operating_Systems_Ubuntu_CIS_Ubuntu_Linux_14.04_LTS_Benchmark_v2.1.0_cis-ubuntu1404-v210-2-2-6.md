# stage: report
# category: CIS_benchmarks


# 2.2.6 Ensure LDAP server is not enabled (Scored)

## Profile Applicability

- Level 1 - Server
- Level 1 - Workstation

## Description

The Lightweight Directory Access Protocol (LDAP) was introduced as a replacement for NIS/YP. It is a service that provides a method for looking up information from a central database.

## Rationale

If the system will not need to act as an LDAP server, it is recommended that the software be disabled to reduce the potential attack surface.

## Audit Procedure

Run the following to ensure no start links for `slapd` exist in `/etc/rc*.d`:

```bash
ls /etc/rc*.d/S*slapd
```

No results should be returned.

## Expected Result

No output should be returned, indicating that `slapd` has no start links.

## Remediation

Run the following command to disable `slapd`:

```bash
update-rc.d slapd disable
```

## Default Value

slapd is not enabled by default.

## References

1. For more detailed documentation on OpenLDAP, go to the project homepage at http://www.openldap.org.

- CIS Controls: 9.1 Limit Open Ports, Protocols, and Services

## Profile

- Level 1 - Server
- Level 1 - Workstation
