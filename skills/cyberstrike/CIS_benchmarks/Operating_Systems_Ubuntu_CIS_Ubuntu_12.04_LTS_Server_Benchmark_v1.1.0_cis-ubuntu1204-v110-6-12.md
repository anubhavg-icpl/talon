# stage: report
# category: CIS_benchmarks


# 6.12 Ensure Samba is not enabled (Not Scored)

## Profile Applicability

- Level 1

## Description

The Samba daemon allows system administrators to configure their Linux systems to share file systems and directories with Windows desktops. Samba will advertise the file systems and directories via the Small Message Block (SMB) protocol. Windows desktop users will be able to mount these directories and file systems as letter drives on their systems.

## Rationale

If there is no need to mount directories and file systems to Windows systems, then this service can be deleted to reduce the potential attack surface.

## Audit Procedure

### Using Command Line

Ensure no start conditions listed for `smbd`:

```bash
initctl show-config smbd smbd
```

## Expected Result

No start conditions should be listed for smbd.

## Remediation

### Using Command Line

Remove or comment out start lines in `/etc/init/smbd.conf`:

```bash
sed -i 's/^start/#start/' /etc/init/smbd.conf
```

## Default Value

Not installed by default on Ubuntu 12.04 LTS Server.

## References

- CIS Ubuntu 12.04 LTS Server Benchmark v1.1.0

## Profile

Level 1 - Not Scored
