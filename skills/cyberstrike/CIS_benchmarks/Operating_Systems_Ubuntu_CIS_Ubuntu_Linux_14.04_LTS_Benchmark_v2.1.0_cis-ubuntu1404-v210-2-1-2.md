# stage: report
# category: CIS_benchmarks


# 2.1.2 Ensure daytime services are not enabled (Scored)

## Profile Applicability

- Level 1 - Server
- Level 1 - Workstation

## Description

`daytime` is a network service that responds with the server's current date and time. This service is intended for debugging and testing purposes. It is recommended that this service be disabled.

## Rationale

Disabling this service will reduce the remote attack surface of the system.

## Audit Procedure

Verify the `daytime` service is not enabled. Run the following command and verify results are as indicated:

```bash
grep -R "^daytime" /etc/inetd.*
```

No results should be returned.

Check `/etc/xinetd.conf` and `/etc/xinetd.d/*` and verify all `daytime` services have `disable = yes` set.

## Expected Result

No output should be returned from the grep command. All daytime services in xinetd should have `disable = yes`.

## Remediation

Comment out or remove any lines starting with `daytime` from `/etc/inetd.conf` and `/etc/inetd.d/*`.

Set `disable = yes` on all `daytime` services in `/etc/xinetd.conf` and `/etc/xinetd.d/*`.

## Default Value

daytime services are not enabled by default.

## References

- CIS Controls: 9.1 Limit Open Ports, Protocols, and Services

## Profile

- Level 1 - Server
- Level 1 - Workstation
