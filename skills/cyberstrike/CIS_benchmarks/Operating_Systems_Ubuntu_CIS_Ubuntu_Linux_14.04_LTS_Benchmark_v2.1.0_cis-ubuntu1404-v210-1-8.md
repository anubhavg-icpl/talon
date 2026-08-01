# stage: report
# category: CIS_benchmarks


# 1.8 Ensure updates, patches, and additional security software are installed (Not Scored)

## Profile Applicability

- Level 1 - Server
- Level 1 - Workstation

## Description

Periodically patches are released for included software either due to security flaws or to include additional functionality.

## Rationale

Newer patches may contain security enhancements that would not be available through the latest full update. As a result, it is recommended that the latest software patches be used to take advantage of the latest functionality. As with any software installation, organizations need to determine if a given update meets their requirements and verify the compatibility and supportability of any additional software against the update revision that is selected.

## Audit Procedure

Run the following command and verify there are no updates or patches to install:

```bash
apt-get -s upgrade
```

## Expected Result

The output should indicate that there are no packages to upgrade (0 upgraded, 0 newly installed, 0 to remove).

## Remediation

Use your package manager to update all packages on the system according to site policy.

## Default Value

Not applicable.

## Notes

Site policy may mandate a testing period before install onto production systems for available updates.

## References

- CIS Controls: 4.5 Use Automated Patch Management And Software Update Tools

## Profile

- Level 1 - Server
- Level 1 - Workstation
