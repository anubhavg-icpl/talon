# stage: report
# category: CIS_benchmarks


# 2.1.4 Ensure echo services are not enabled (Scored)

## Profile Applicability

- Level 1 - Server
- Level 1 - Workstation

## Description

`echo` is a network service that responds to clients with the data sent to it by the client. This service is intended for debugging and testing purposes. It is recommended that this service be disabled.

## Rationale

Disabling this service will reduce the remote attack surface of the system.

## Audit Procedure

Verify the `echo` service is not enabled. Run the following command and verify results are as indicated:

```bash
grep -R "^echo" /etc/inetd.*
```

No results should be returned.

Check `/etc/xinetd.conf` and `/etc/xinetd.d/*` and verify all `echo` services have `disable = yes` set.

## Expected Result

No output should be returned from the grep command. All echo services in xinetd should have `disable = yes`.

## Remediation

Comment out or remove any lines starting with `echo` from `/etc/inetd.conf` and `/etc/inetd.d/*`.

Set `disable = yes` on all `echo` services in `/etc/xinetd.conf` and `/etc/xinetd.d/*`.

## Default Value

echo services are not enabled by default.

## References

- CIS Controls: 9.1 Limit Open Ports, Protocols, and Services

## Profile

- Level 1 - Server
- Level 1 - Workstation
