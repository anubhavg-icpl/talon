# stage: report
# category: CIS_benchmarks


# 3.5.3 Ensure RDS is disabled (Not Scored)

## Profile Applicability

- Level 1 - Server
- Level 1 - Workstation

## Description

The Reliable Datagram Sockets (RDS) protocol is a transport layer protocol designed to provide low-latency, high-bandwidth communications between cluster nodes. It was developed by the Oracle Corporation.

## Rationale

If the protocol is not being used, it is recommended that kernel module not be loaded, disabling the service to reduce the potential attack surface.

## Audit Procedure

Run the following commands and verify the output is as indicated:

```bash
modprobe -n -v rds
# Expected: install /bin/true

lsmod | grep rds
# Expected: <No output>
```

## Expected Result

```
install /bin/true
```

And no output from `lsmod | grep rds`.

## Remediation

Edit or create the file `/etc/modprobe.d/CIS.conf` and add the following line:

```
install rds /bin/true
```

## Default Value

Not disabled by default.

## References

- CIS Controls: 9.1 - Limit Open Ports, Protocols, and Services

## Profile

- Level 1 - Server
- Level 1 - Workstation
