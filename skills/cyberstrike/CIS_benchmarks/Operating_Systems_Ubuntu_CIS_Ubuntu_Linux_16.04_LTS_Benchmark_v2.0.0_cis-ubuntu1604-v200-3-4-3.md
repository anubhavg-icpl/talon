# stage: report
# category: CIS_benchmarks


# CIS Ubuntu Linux 16.04 LTS Benchmark v2.0.0 - Control 3.4.3

## Description

The Reliable Datagram Sockets (RDS) protocol is a transport layer protocol designed to provide low-latency, high-bandwidth communications between cluster nodes. It was developed by the Oracle Corporation.

## Rationale

If the protocol is not being used, it is recommended that kernel module not be loaded, disabling the service to reduce the potential attack surface.

## Impact

None.

## Audit Procedure

### Command Line

Run the following commands and verify the output is as indicated:

```bash
modprobe -n -v rds
```

```bash
lsmod | grep rds
```

## Expected Result

```
install /bin/true
```

The `lsmod` command should return no output.

## Remediation

### Command Line

Edit or create a file in the `/etc/modprobe.d/` directory ending in `.conf`

Example: `vi /etc/modprobe.d/rds.conf`

and add the following line:

```
install rds /bin/true
```

## Default Value

RDS module is available but typically not loaded by default.

## References

1. CIS Ubuntu Linux 16.04 LTS Benchmark v2.0.0

## CIS Controls

Version 7

9.2 Ensure Only Approved Ports, Protocols and Services Are Running - Ensure that only network ports, protocols, and services listening on a system with validated business needs, are running on each system.

## Profile Applicability

- Level 2 - Server
- Level 2 - Workstation

## Assessment Status

Automated
