# stage: report
# category: CIS_benchmarks


# 3.5.1 Ensure DCCP is disabled (Not Scored)

## Profile Applicability

- Level 1 - Server
- Level 1 - Workstation

## Description

The Datagram Congestion Control Protocol (DCCP) is a transport layer protocol that supports streaming media and telephony. DCCP provides a way to gain access to congestion control, without having to do it at the application layer, but does not provide in-sequence delivery.

## Rationale

If the protocol is not required, it is recommended that the drivers not be installed to reduce the potential attack surface.

## Audit Procedure

Run the following commands and verify the output is as indicated:

```bash
modprobe -n -v dccp
# Expected: install /bin/true

lsmod | grep dccp
# Expected: <No output>
```

## Expected Result

```
install /bin/true
```

And no output from `lsmod | grep dccp`.

## Remediation

Edit or create the file `/etc/modprobe.d/CIS.conf` and add the following line:

```
install dccp /bin/true
```

## Default Value

Not disabled by default.

## References

- CIS Controls: 9.1 - Limit Open Ports, Protocols, and Services

## Profile

- Level 1 - Server
- Level 1 - Workstation
