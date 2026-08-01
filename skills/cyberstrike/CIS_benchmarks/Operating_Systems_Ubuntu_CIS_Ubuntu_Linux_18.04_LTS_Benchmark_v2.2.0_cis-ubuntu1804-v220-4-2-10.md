# stage: report
# category: CIS_benchmarks


# CIS Ubuntu Linux 18.04 LTS Benchmark v2.2.0 - Control 4.2.10

## Profile Applicability

- Level 1 - Server
- Level 1 - Workstation

## Description

The `IgnoreRhosts` parameter specifies that `.rhosts` and `.shosts` files will not be used in `RhostsRSAAuthentication` or `HostbasedAuthentication`.

## Rationale

Setting this parameter forces users to enter a password when authenticating with SSH.

## Audit Procedure

### Command Line

Run the following command and verify the output:

```bash
sshd -T | grep -i ignorerhosts
```

### Expected Result

```
ignorerhosts yes
```

## Remediation

### Command Line

Edit the `/etc/ssh/sshd_config` file to set the parameter as follows:

```bash
IgnoreRhosts yes
```

## Default Value

IgnoreRhosts yes

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
