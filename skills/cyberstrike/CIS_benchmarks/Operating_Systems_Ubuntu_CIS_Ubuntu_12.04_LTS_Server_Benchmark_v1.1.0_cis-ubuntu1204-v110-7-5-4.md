# stage: report
# category: CIS_benchmarks


# 7.5.4 Disable TIPC (Not Scored)

## Profile Applicability

- Level 1

## Description

The Transparent Inter-Process Communication (TIPC) protocol is designed to provide communication between cluster nodes.

## Rationale

If the protocol is not being used, it is recommended that kernel module not be loaded, disabling the service to reduce the potential attack surface.

## Audit Procedure

### Using Command Line

Perform the following to determine if TIPC is disabled.

```bash
grep "install tipc /bin/true" /etc/modprobe.d/CIS.conf
```

## Expected Result

```
install tipc /bin/true
```

## Remediation

### Using Command Line

```bash
echo "install tipc /bin/true" >> /etc/modprobe.d/CIS.conf
```

## Default Value

TIPC kernel module is available by default.

## References

- CIS Ubuntu 12.04 LTS Server Benchmark v1.1.0

## Profile

Level 1 - Not Scored
