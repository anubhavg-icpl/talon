# stage: report
# category: CIS_benchmarks


# 1.3.2 Ensure filesystem integrity is regularly checked (Scored)

## Profile Applicability

- Level 1 - Server
- Level 1 - Workstation

## Description

Periodic checking of the filesystem integrity is needed to detect changes to the filesystem.

## Rationale

Periodic file checking allows the system administrator to determine on a regular basis if critical files have been changed in an unauthorized fashion.

## Audit Procedure

Run the following commands to determine if there is a cron job scheduled to run the aide check:

```bash
crontab -u root -l | grep aide
grep -r aide /etc/cron.* /etc/crontab
```

Ensure a cron job in compliance with site policy is returned.

## Expected Result

A cron job entry for AIDE should be returned showing scheduled integrity checks.

## Remediation

Run the following command:

```bash
crontab -u root -e
```

Add the following line to the crontab:

```
0 5 * * * /usr/bin/aide --config /etc/aide/aide.conf --check
```

## Default Value

No AIDE cron job is configured by default.

## Notes

The checking in this recommendation occurs every day at 5am. Alter the frequency and time of the checks in compliance with site policy.

## References

- CIS Controls: 3.5 Use File Integrity Tools For Critical System Files

## Profile

- Level 1 - Server
- Level 1 - Workstation
