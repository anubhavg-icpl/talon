# stage: report
# category: CIS_benchmarks


# CIS Ubuntu Linux 16.04 LTS Benchmark v2.0.0 - Control 3.4.2

## Description

The Stream Control Transmission Protocol (SCTP) is a transport layer protocol used to support message oriented communication, with several streams of messages in one connection. It serves a similar function as TCP and UDP, incorporating features of both. It is message-oriented like UDP, and ensures reliable in-sequence transport of messages with congestion control like TCP.

## Rationale

If the protocol is not being used, it is recommended that kernel module not be loaded, disabling the service to reduce the potential attack surface.

## Impact

None.

## Audit Procedure

### Command Line

Run the following commands and verify the output is as indicated:

```bash
modprobe -n -v sctp | grep -E '(sctp|install)'
```

```bash
lsmod | grep sctp
```

## Expected Result

```
install /bin/true
```

The `lsmod` command should return no output.

## Remediation

### Command Line

Edit or create a file in the `/etc/modprobe.d/` directory ending in `.conf`

Example: `vi /etc/modprobe.d/sctp.conf`

and add the following line:

```
install sctp /bin/true
```

## Default Value

SCTP module is available but typically not loaded by default.

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
