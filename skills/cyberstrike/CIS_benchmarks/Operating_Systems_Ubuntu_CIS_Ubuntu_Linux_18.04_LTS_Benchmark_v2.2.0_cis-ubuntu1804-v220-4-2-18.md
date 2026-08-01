# stage: report
# category: CIS_benchmarks


# CIS Ubuntu Linux 18.04 LTS Benchmark v2.2.0 - Control 4.2.18

## Profile Applicability

- Level 1 - Server
- Level 1 - Workstation

## Description

The `PermitEmptyPasswords` parameter specifies if the SSH server allows login to accounts with empty password strings.

## Rationale

Disallowing remote shell access to accounts that have an empty password reduces the probability of unauthorized access to the system.

## Audit Procedure

### Command Line

Run the following command and verify the output:

```bash
sshd -T | grep -i permitemptypasswords
```

### Expected Result

```
permitemptypasswords no
```

## Remediation

### Command Line

Edit the `/etc/ssh/sshd_config` file to set the parameter as follows:

```bash
PermitEmptyPasswords no
```

## Default Value

PermitEmptyPasswords no

## References

1. NIST SP 800-53 Rev. 5: CM-7

## CIS Controls

v8 - 4.8 Uninstall or Disable Unnecessary Services on Enterprise Assets and Software.

v7 - 9.2 Ensure Only Approved Ports, Protocols, and Services Are Running.

## Profile Applicability

- Level 1 - Server
- Level 1 - Workstation

## Assessment Status

Automated
