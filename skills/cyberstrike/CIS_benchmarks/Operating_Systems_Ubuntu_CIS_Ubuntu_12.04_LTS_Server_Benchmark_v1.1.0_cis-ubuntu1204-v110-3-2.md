# stage: report
# category: CIS_benchmarks


# 3.2 Set Permissions on bootloader config (Scored)

## Profile Applicability

- Level 1

## Description

Set permission on the your boot loaders config file to read and write for root only.

## Rationale

Setting the permissions to read and write for root only prevents non-root users from seeing the boot parameters or changing them. Non-root users who read the boot parameters may be able to identify weaknesses in security upon boot and be able to exploit them.

## Audit Procedure

### Using Command Line

Perform the following to determine if the `/boot/grub/grub.cfg` file permissions are correct:

```bash
stat -L -c "%a" /boot/grub/grub.cfg | egrep ".00"
```

## Expected Result

The command should return a permissions value ending in `00` (e.g., `400` or `600`). If the above command emits no output then the system is not configured as recommended.

## Remediation

### Using Command Line

Run the following to set the permissions for `/boot/grub/grub.cfg`:

```bash
chmod og-rwx /boot/grub/grub.cfg
```

## Default Value

By default, `/boot/grub/grub.cfg` may have permissions more permissive than recommended.

## References

- CIS Ubuntu 12.04 LTS Server Benchmark v1.1.0

## Profile

Level 1 - Scored
