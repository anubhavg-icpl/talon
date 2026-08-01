# stage: report
# category: CIS_benchmarks


# 2.22 Disable Mounting of hfsplus Filesystems (Not Scored)

## Profile Applicability

- Level 2

## Description

The `hfsplus` filesystem type is a hierarchical filesystem designed to replace `hfs` that allows you to mount Mac OS filesystems.

## Rationale

Removing support for unneeded filesystem types reduces the local attack surface of the server. If this filesystem type is not needed, disable it.

## Audit Procedure

### Using Command Line

```bash
/sbin/modprobe -n -v hfsplus
/sbin/lsmod | grep hfsplus
```

## Expected Result

The first command should output `install /bin/true`. The second command should produce no output (module not loaded).

## Remediation

### Using Command Line

Edit or create the file `/etc/modprobe.d/CIS.conf` and add the following line:

```
install hfsplus /bin/true
```

## Default Value

By default, the `hfsplus` module is available for loading.

## References

- CIS Ubuntu 12.04 LTS Server Benchmark v1.1.0

## Profile

Level 2 - Not Scored
