# stage: report
# category: CIS_benchmarks


# 2.2.14 Ensure SNMP Server is not enabled (Scored)

## Profile Applicability

- Level 1 - Server
- Level 1 - Workstation

## Description

The Simple Network Management Protocol (SNMP) server is used to listen for SNMP commands from an SNMP management system, execute the commands or collect the information and then send results back to the requesting system.

## Rationale

The SNMP server can communicate using SNMP v1, which transmits data in the clear and does not require authentication to execute commands. Unless absolutely necessary, it is recommended that the SNMP service not be used. If SNMP is required the server should be configured to disallow SNMP v1.

## Audit Procedure

Run the following to ensure no start links for `snmpd` exist in `/etc/rc*.d`:

```bash
ls /etc/rc*.d/S*snmpd
```

No results should be returned.

## Expected Result

No output should be returned, indicating that `snmpd` has no start links.

## Remediation

Run the following command to disable `snmpd`:

```bash
update-rc.d snmpd disable
```

**Note:** Additional methods of disabling a service exist. Consult your distribution documentation for appropriate methods.

## Default Value

SNMP server is not enabled by default.

## References

- CIS Controls: 9.1 Limit Open Ports, Protocols, and Services

## Profile

- Level 1 - Server
- Level 1 - Workstation
