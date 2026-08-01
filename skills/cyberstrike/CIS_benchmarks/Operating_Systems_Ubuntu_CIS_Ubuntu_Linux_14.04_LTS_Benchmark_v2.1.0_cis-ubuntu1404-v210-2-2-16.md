# stage: report
# category: CIS_benchmarks


# 2.2.16 Ensure rsync service is not enabled (Scored)

## Profile Applicability

- Level 1 - Server
- Level 1 - Workstation

## Description

The `rsyncd` service can be used to synchronize files between systems over network links.

## Rationale

The `rsyncd` service presents a security risk as it uses unencrypted protocols for communication.

## Audit Procedure

Run the following command to verify that the `rsync` service is not enabled:

```bash
grep ^RSYNC_ENABLE /etc/default/rsync
```

Expected output:

```
RSYNC_ENABLE=false
```

## Expected Result

The output should show `RSYNC_ENABLE=false`.

## Remediation

Edit the `/etc/default/rsync` file and set `RSYNC_ENABLE` to `false`:

```
RSYNC_ENABLE=false
```

## Default Value

rsync service is not enabled by default (`RSYNC_ENABLE=false`).

## References

- CIS Controls: 9.1 Limit Open Ports, Protocols, and Services

## Profile

- Level 1 - Server
- Level 1 - Workstation
