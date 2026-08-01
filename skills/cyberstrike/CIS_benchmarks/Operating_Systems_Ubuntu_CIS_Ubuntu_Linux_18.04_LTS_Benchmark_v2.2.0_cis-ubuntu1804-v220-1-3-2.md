# stage: report
# category: CIS_benchmarks


# 1.3.2 Ensure package manager repositories are configured (Manual)

## Profile Applicability

- Level 1 - Server
- Level 1 - Workstation

## Description

Systems need to have package manager repositories configured to ensure they receive the latest patches and updates.

## Rationale

If a system's package repositories are misconfigured important patches may not be identified or a rogue repository could introduce compromised software.

## Audit Procedure

### Command Line

Run the following command and verify package repositories are configured correctly:

```bash
apt-cache policy
```

## Expected Result

Package repositories should be properly configured according to site policy.

## Remediation

### Command Line

Configure your package manager repositories according to site policy.

## References

1. NIST SP 800-53 Rev. 5: SI-2

## CIS Controls

| Controls Version | Control                                                      | IG 1 | IG 2 | IG 3 |
| ---------------- | ------------------------------------------------------------ | ---- | ---- | ---- |
| v8               | 7.3 Perform Automated Operating System Patch Management      | X    | X    | X    |
| v7               | 3.4 Deploy Automated Operating System Patch Management Tools | X    | X    | X    |
| v7               | 3.5 Deploy Automated Software Patch Management Tools         | X    | X    | X    |

## MITRE ATT&CK Mappings

| Techniques / Sub-techniques                                                                                           | Tactics | Mitigations |
| --------------------------------------------------------------------------------------------------------------------- | ------- | ----------- |
| T1068, T1068.000, T1195, T1195.001, T1195.002, T1203, T1203.000, T1210, T1210.000, T1211, T1211.000, T1212, T1212.000 | TA0001  | M1051       |
