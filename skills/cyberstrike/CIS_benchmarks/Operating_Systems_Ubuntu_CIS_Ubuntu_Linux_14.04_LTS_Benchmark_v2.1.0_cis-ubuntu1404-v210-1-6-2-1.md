# stage: report
# category: CIS_benchmarks


# 1.6.2.1 Ensure AppArmor is not disabled in bootloader configuration (Scored)

## Profile Applicability

- Level 2 - Server
- Level 2 - Workstation

## Description

Configure AppArmor to be enabled at boot time and verify that it has not been overwritten by the bootloader boot parameters.

## Rationale

AppArmor must be enabled at boot time in your bootloader configuration to ensure that the controls it provides are not overridden.

## Audit Procedure

Run the following command and verify that no linux line the `apparmor=0` parameter set:

```bash
grep "^\s*linux" /boot/grub/grub.cfg
```

## Expected Result

No linux lines should contain `apparmor=0`.

## Remediation

Edit `/etc/default/grub` and remove all instances of `apparmor=0` from all CMDLINE_LINUX parameters:

```
GRUB_CMDLINE_LINUX_DEFAULT="quiet"
GRUB_CMDLINE_LINUX=""
```

Run the following command to update the grub2 configuration:

```bash
update-grub
```

## Default Value

AppArmor is enabled by default on Ubuntu systems.

## Notes

This recommendation is designed around the grub bootloader, if LILO or another bootloader is in use in your environment enact equivalent settings.

## References

1. AppArmor Documentation: http://wiki.apparmor.net/index.php/Documentation
2. Ubuntu AppArmor Documentation: https://help.ubuntu.com/community/AppArmor
3. SUSE AppArmor Documentation: https://www.suse.com/documentation/apparmor/

## Profile

- Level 2 - Server
- Level 2 - Workstation
