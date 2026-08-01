# stage: report
# category: CIS_benchmarks


# 6.6 Ensure LDAP is not enabled (Not Scored)

## Profile Applicability

- Level 1

## Description

The Lightweight Directory Access Protocol (LDAP) was introduced as a replacement for NIS/YP. It is a service that provides a method for looking up information from a central database.

## Rationale

If the server will not need to act as an LDAP client or server, it is recommended that the software be disabled to reduce the potential attack surface.

## Audit Procedure

### Using Command Line

Run the following command:

```bash
dpkg -s slapd
```

## Expected Result

Ensure package status is not-installed or dpkg returns no info is available.

## Remediation

### Using Command Line

Uninstall the `slapd` package:

```bash
apt-get purge slapd
```

## Default Value

Not installed by default on Ubuntu 12.04 LTS Server.

## References

- CIS Ubuntu 12.04 LTS Server Benchmark v1.1.0
- http://www.openldap.org

## Profile

Level 1 - Not Scored
