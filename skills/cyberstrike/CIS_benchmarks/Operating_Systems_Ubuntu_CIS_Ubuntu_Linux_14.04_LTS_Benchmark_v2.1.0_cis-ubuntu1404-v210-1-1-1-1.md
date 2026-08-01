# stage: report
# category: CIS_benchmarks


# 1.1.1.1 Ensure mounting of cramfs filesystems is disabled (Scored)

## Profile Applicability

- Level 1 - Server
- Level 1 - Workstation

## Description

The cramfs filesystem type is a compressed read-only Linux filesystem embedded in small footprint systems. A cramfs image can be used without having to first decompress the image.

## Rationale

Removing support for unneeded filesystem types reduces the local attack surface of the server. If this filesystem type is not needed, disable it.

## Audit Procedure

```bash
modprobe -n -v cramfs
# Expected output: install /bin/true

lsmod | grep cramfs
# Expected output: <No output>
```

## Expected Result

`modprobe -n -v cramfs` should return `install /bin/true` and `lsmod | grep cramfs` should return no output.

## Remediation

```bash
# Edit or create the file /etc/modprobe.d/CIS.conf and add the following line:
echo "install cramfs /bin/true" >> /etc/modprobe.d/CIS.conf

# Unload the cramfs module:
rmmod cramfs
```

## Default Value

By default, cramfs filesystem mounting is allowed.

## References

- CIS Ubuntu Linux 14.04 LTS Benchmark v2.1.0
- CIS Controls: 13 Data Protection

## Profile

- Level 1
