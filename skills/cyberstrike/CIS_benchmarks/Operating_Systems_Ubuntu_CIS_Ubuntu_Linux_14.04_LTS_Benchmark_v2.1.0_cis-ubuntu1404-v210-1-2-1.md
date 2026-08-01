# stage: report
# category: CIS_benchmarks


# 1.2.1 Ensure package manager repositories are configured (Not Scored)

## Profile Applicability

- Level 1 - Server
- Level 1 - Workstation

## Description

Systems need to have package manager repositories configured to ensure they receive the latest patches and updates.

## Rationale

If a system's package repositories are misconfigured important patches may not be identified or a rogue repository could introduce compromised software.

## Audit Procedure

Run the following command and verify package repositories are configured correctly:

```bash
apt-cache policy
```

## Expected Result

Verify that repositories are configured according to site policy.

## Remediation

Configure your package manager repositories according to site policy.

## Default Value

Not applicable. Configuration varies by site policy.

## References

- CIS Controls: 4.5 Use Automated Patch Management And Software Update Tools

## Profile

- Level 1 - Server
- Level 1 - Workstation
