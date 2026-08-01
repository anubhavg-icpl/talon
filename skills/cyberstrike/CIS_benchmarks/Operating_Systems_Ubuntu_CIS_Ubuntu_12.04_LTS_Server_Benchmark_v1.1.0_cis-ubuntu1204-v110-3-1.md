# stage: report
# category: CIS_benchmarks


# 3.1 Set User/Group Owner on bootloader config (Scored)

## Profile Applicability

- Level 1

## Description

Set the owner and group of your boot loaders config file to the root user. These instructions default to GRUB stored at `/boot/grub/grub.cfg`.

## Rationale

Setting the owner and group to root prevents non-root users from changing the file.

## Audit Procedure

### Using Command Line

Perform the following to determine if the `/boot/grub/grub.cfg` file has the correct ownership:

```bash
stat -c "%u %g" /boot/grub/grub.cfg | egrep "^0 0"
```

## Expected Result

The command should return `0 0` indicating root:root ownership. If the above command emits no output then the system is not configured as recommended.

## Remediation

### Using Command Line

Run the following to change ownership of `/boot/grub/grub.cfg`:

```bash
chown root:root /boot/grub/grub.cfg
```

## Default Value

By default, `/boot/grub/grub.cfg` is owned by root:root.

## References

- CIS Ubuntu 12.04 LTS Server Benchmark v1.1.0

## Profile

Level 1 - Scored
