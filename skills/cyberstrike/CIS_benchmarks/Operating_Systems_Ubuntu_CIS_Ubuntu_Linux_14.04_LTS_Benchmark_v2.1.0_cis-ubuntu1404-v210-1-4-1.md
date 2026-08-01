# stage: report
# category: CIS_benchmarks


# 1.4.1 Ensure permissions on bootloader config are configured (Scored)

## Profile Applicability

- Level 1 - Server
- Level 1 - Workstation

## Description

The grub configuration file contains information on boot settings and passwords for unlocking boot options. The grub configuration is usually `grub.cfg` stored in `/boot/grub`.

## Rationale

Setting the permissions to read and write for root only prevents non-root users from seeing the boot parameters or changing them. Non-root users who read the boot parameters may be able to identify weaknesses in security upon boot and be able to exploit them.

## Audit Procedure

Run the following command and verify `Uid` and `Gid` are both `0/root` and `Access` does not grant permissions to `group` or `other`:

```bash
stat /boot/grub/grub.cfg
```

## Expected Result

```
Access: (0400/-r--------) Uid: ( 0/ root) Gid: ( 0/ root)
```

## Remediation

Run the following commands to set permissions on your grub configuration:

```bash
chown root:root /boot/grub/grub.cfg
chmod og-rwx /boot/grub/grub.cfg
```

## Default Value

Not applicable.

## Notes

This recommendation is designed around the grub bootloader, if LILO or another bootloader is in use in your environment enact equivalent settings.

## References

- CIS Controls: 5.1 Minimize And Sparingly Use Administrative Privileges

## Profile

- Level 1 - Server
- Level 1 - Workstation
