# stage: report
# category: CIS_benchmarks


# 2.2.5 Ensure LDAP server is not installed (Automated)

## Profile

- Level 1 - Server
- Level 1 - Workstation

## Description

The Lightweight Directory Access Protocol (LDAP) was introduced as a replacement for NIS/YP. It is a service that provides a method for looking up information from a central database.

## Rationale

If the system will not need to act as an LDAP server, it is recommended that the software be removed to reduce the potential attack surface.

## Audit Procedure

### Command Line

Run the following command to verify `slapd` is not installed:

```bash
# dpkg-query -s slapd &>/dev/null && echo "slapd is installed"
```

Nothing should be returned.

## Expected Result

No output (empty result).

## Remediation

### Command Line

Run one of the following commands to remove `slapd`:

```bash
# apt purge slapd
```

## References

1. For more detailed documentation on OpenLDAP, go to the project homepage at http://www.openldap.org.
2. NIST SP 800-53 Rev. 5: CM-7

## CIS Controls

- v8: 4.8 - Uninstall or Disable Unnecessary Services on Enterprise Assets and Software
- v7: 9.2 - Ensure Only Approved Ports, Protocols and Services Are Running
