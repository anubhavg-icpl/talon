# stage: report
# category: CIS_benchmarks


# 3.4.2 Ensure /etc/hosts.allow is configured (Scored)

## Profile Applicability

- Level 1 - Server
- Level 1 - Workstation

## Description

The `/etc/hosts.allow` file specifies which IP addresses are permitted to connect to the host. It is intended to be used in conjunction with the `/etc/hosts.deny` file.

## Rationale

The `/etc/hosts.allow` file supports access control by IP and helps ensure that only authorized systems can connect to the system.

## Audit Procedure

Run the following command and verify the contents of the `/etc/hosts.allow` file:

```bash
cat /etc/hosts.allow
```

## Expected Result

The file should contain appropriate network allow rules for your environment.

## Remediation

Run the following command to create `/etc/hosts.allow`:

```bash
echo "ALL: <net>/<mask>, <net>/<mask>, ..." >/etc/hosts.allow
```

where each `<net>/<mask>` combination (for example, "192.168.1.0/255.255.255.0") represents one network block in use by your organization that requires access to this system.

## Default Value

Contents of `/etc/hosts.allow` file will vary depending on your network configuration.

## References

- CIS Controls: 9.2 - Leverage Host-based Firewalls

## Profile

- Level 1 - Server
- Level 1 - Workstation
