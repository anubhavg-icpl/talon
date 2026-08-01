# stage: report
# category: CIS_benchmarks


# CIS Ubuntu Linux 16.04 LTS Benchmark v2.0.0 - 1.6.1.2

## Profile Applicability

- Level 1 - Server
- Level 1 - Workstation

## Description

Configure AppArmor to be enabled at boot time and verify that it has not been overwritten by the bootloader boot parameters.

_Note: This recommendation is designed around the grub bootloader, if LILO or another bootloader is in use in your environment enact equivalent settings._

## Rationale

AppArmor must be enabled at boot time in your bootloader configuration to ensure that the controls it provides are not overridden.

## Audit Procedure

### Command Line

Run the following commands to verify that all `linux` lines have the `apparmor=1` and `security=apparmor` parameters set:

```bash
grep "^\s*linux" /boot/grub/grub.cfg | grep -v "apparmor=1"
```

Nothing should be returned.

```bash
grep "^\s*linux" /boot/grub/grub.cfg | grep -v "security=apparmor"
```

Nothing should be returned.

## Expected Result

Both commands should return no output, indicating all linux boot lines include `apparmor=1` and `security=apparmor`.

## Remediation

### Command Line

Edit `/etc/default/grub` and add the `apparmor=1` and `security=apparmor` parameters to the `GRUB_CMDLINE_LINUX=` line:

```
GRUB_CMDLINE_LINUX="apparmor=1 security=apparmor"
```

Run the following command to update the grub2 configuration:

```bash
update-grub
```

## Default Value

Not applicable.

## References

None.

## CIS Controls

| Controls Version | Control                                               |
| ---------------- | ----------------------------------------------------- |
| v7               | 14.6 Protect Information through Access Control Lists |

## Assessment Status

Automated
