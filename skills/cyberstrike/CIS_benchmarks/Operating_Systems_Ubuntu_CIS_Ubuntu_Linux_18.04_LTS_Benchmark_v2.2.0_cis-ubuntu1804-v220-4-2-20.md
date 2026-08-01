# stage: report
# category: CIS_benchmarks


# CIS Ubuntu Linux 18.04 LTS Benchmark v2.2.0 - Control 4.2.20

## Profile Applicability

- Level 1 - Server
- Level 1 - Workstation

## Description

The `PermitUserEnvironment` option allows users to present environment options to the SSH daemon.

## Rationale

Permitting users the ability to set environment variables through the SSH daemon could potentially allow users to bypass security controls (e.g. setting an execution path that has SSH executing trojan'd programs).

## Audit Procedure

### Command Line

Run the following command and verify the output:

```bash
sshd -T | grep -i permituserenvironment
```

### Expected Result

```
permituserenvironment no
```

## Remediation

### Command Line

Edit the `/etc/ssh/sshd_config` file to set the parameter as follows:

```bash
PermitUserEnvironment no
```

## Default Value

PermitUserEnvironment no

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
