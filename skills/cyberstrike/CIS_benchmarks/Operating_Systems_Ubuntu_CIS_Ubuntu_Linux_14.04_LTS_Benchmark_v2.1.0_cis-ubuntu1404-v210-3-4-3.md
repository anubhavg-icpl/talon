# stage: report
# category: CIS_benchmarks


# 3.4.3 Ensure /etc/hosts.deny is configured (Scored)

## Profile Applicability

- Level 1 - Server
- Level 1 - Workstation

## Description

The `/etc/hosts.deny` file specifies which IP addresses are **not** permitted to connect to the host. It is intended to be used in conjunction with the `/etc/hosts.allow` file.

## Rationale

The `/etc/hosts.deny` file serves as a failsafe so that any host not specified in `/etc/hosts.allow` is denied access to the system.

## Audit Procedure

Run the following command and verify the contents of the `/etc/hosts.deny` file:

```bash
cat /etc/hosts.deny
```

## Expected Result

```
ALL: ALL
```

## Remediation

Run the following command to create `/etc/hosts.deny`:

```bash
echo "ALL: ALL" >> /etc/hosts.deny
```

## Default Value

Contents of `/etc/hosts.deny` file may include additional options depending on your network configuration.

## References

- CIS Controls: 9.2 - Leverage Host-based Firewalls

## Profile

- Level 1 - Server
- Level 1 - Workstation
