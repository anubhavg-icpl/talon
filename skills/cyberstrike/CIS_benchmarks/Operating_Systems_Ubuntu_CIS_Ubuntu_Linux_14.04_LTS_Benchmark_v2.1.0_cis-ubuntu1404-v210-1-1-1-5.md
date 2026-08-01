# stage: report
# category: CIS_benchmarks


# 1.1.1.5 Ensure mounting of hfsplus filesystems is disabled (Scored)

## Profile Applicability

- Level 1 - Server
- Level 1 - Workstation

## Description

The hfsplus filesystem type is a hierarchical filesystem designed to replace hfs that allows you to mount Mac OS filesystems.

## Rationale

Removing support for unneeded filesystem types reduces the local attack surface of the system. If this filesystem type is not needed, disable it.

## Audit Procedure

```bash
modprobe -n -v hfsplus
# Expected output: install /bin/true

lsmod | grep hfsplus
# Expected output: <No output>
```

## Expected Result

`modprobe -n -v hfsplus` should return `install /bin/true` and `lsmod | grep hfsplus` should return no output.

## Remediation

```bash
# Edit or create the file /etc/modprobe.d/CIS.conf and add the following line:
echo "install hfsplus /bin/true" >> /etc/modprobe.d/CIS.conf

# Unload the hfsplus module:
rmmod hfsplus
```

## Default Value

By default, hfsplus filesystem mounting is allowed.

## References

- CIS Ubuntu Linux 14.04 LTS Benchmark v2.1.0
- CIS Controls: 13 Data Protection

## Profile

- Level 1
