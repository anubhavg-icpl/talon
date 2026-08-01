# stage: report
# category: CIS_benchmarks


# 1.1.1.2 Ensure mounting of freevxfs filesystems is disabled (Scored)

## Profile Applicability

- Level 1 - Server
- Level 1 - Workstation

## Description

The freevxfs filesystem type is a free version of the Veritas type filesystem. This is the primary filesystem type for HP-UX operating systems.

## Rationale

Removing support for unneeded filesystem types reduces the local attack surface of the system. If this filesystem type is not needed, disable it.

## Audit Procedure

```bash
modprobe -n -v freevxfs
# Expected output: install /bin/true

lsmod | grep freevxfs
# Expected output: <No output>
```

## Expected Result

`modprobe -n -v freevxfs` should return `install /bin/true` and `lsmod | grep freevxfs` should return no output.

## Remediation

```bash
# Edit or create the file /etc/modprobe.d/CIS.conf and add the following line:
echo "install freevxfs /bin/true" >> /etc/modprobe.d/CIS.conf

# Unload the freevxfs module:
rmmod freevxfs
```

## Default Value

By default, freevxfs filesystem mounting is allowed.

## References

- CIS Ubuntu Linux 14.04 LTS Benchmark v2.1.0
- CIS Controls: 13 Data Protection

## Profile

- Level 1
