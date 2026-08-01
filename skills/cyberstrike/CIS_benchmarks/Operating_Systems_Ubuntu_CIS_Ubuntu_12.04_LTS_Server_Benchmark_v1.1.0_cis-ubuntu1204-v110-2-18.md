# stage: report
# category: CIS_benchmarks


# 2.18 Disable Mounting of cramfs Filesystems (Not Scored)

## Profile Applicability

- Level 2

## Description

The `cramfs` filesystem type is a compressed read-only Linux filesystem embedded in small footprint systems. A `cramfs` image can be used without having to first decompress the image.

## Rationale

Removing support for unneeded filesystem types reduces the local attack surface of the server. If this filesystem type is not needed, disable it.

## Audit Procedure

### Using Command Line

```bash
/sbin/modprobe -n -v cramfs
/sbin/lsmod | grep cramfs
```

## Expected Result

The first command should output `install /bin/true`. The second command should produce no output (module not loaded).

## Remediation

### Using Command Line

Edit or create the file `/etc/modprobe.d/CIS.conf` and add the following line:

```
install cramfs /bin/true
```

## Default Value

By default, the `cramfs` module is available for loading.

## References

- CIS Ubuntu 12.04 LTS Server Benchmark v1.1.0

## Profile

Level 2 - Not Scored
