# stage: report
# category: CIS_benchmarks


# 2.1.1 Ensure chargen services are not enabled (Scored)

## Profile Applicability

- Level 1 - Server
- Level 1 - Workstation

## Description

`chargen` is a network service that responds with 0 to 512 ASCII characters for each connection it receives. This service is intended for debugging and testing purposes. It is recommended that this service be disabled.

## Rationale

Disabling this service will reduce the remote attack surface of the system.

## Audit Procedure

Verify the `chargen` service is not enabled. Run the following command and verify results are as indicated:

```bash
grep -R "^chargen" /etc/inetd.*
```

No results should be returned.

Check `/etc/xinetd.conf` and `/etc/xinetd.d/*` and verify all `chargen` services have `disable = yes` set.

## Expected Result

No output should be returned from the grep command. All chargen services in xinetd should have `disable = yes`.

## Remediation

Comment out or remove any lines starting with `chargen` from `/etc/inetd.conf` and `/etc/inetd.d/*`.

Set `disable = yes` on all `chargen` services in `/etc/xinetd.conf` and `/etc/xinetd.d/*`.

## Default Value

chargen services are not enabled by default.

## References

- CIS Controls: 9.1 Limit Open Ports, Protocols, and Services

## Profile

- Level 1 - Server
- Level 1 - Workstation
