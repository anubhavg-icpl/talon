# stage: report
# category: CIS_benchmarks


# CIS OKE Benchmark v1.7.0 - Control 5.1.4

## Profile Applicability

- **Level:** 1

## Description

Use approved container registries.

## Rationale

Allowing unrestricted access to external container registries provides the opportunity for malicious or unapproved containers to be deployed into the cluster. Allow listing only approved container registries reduces this risk.

## Impact

All container images to be deployed to the cluster must be hosted within an approved container image registry.

## Audit Procedure

No specific audit procedure provided. Review organizational policies to ensure only approved container registries are in use.

## Remediation

No specific remediation procedure provided. Implement organizational policies to restrict container registries to approved sources only.

## Default Value

Not applicable.

## References

1. https://docs.cloud.oracle.com/en-us/iaas/Content/ContEng/Concepts/contengoverview.htm

## CIS Controls

| Controls Version | Control                          | IG 1 | IG 2 | IG 3 |
| ---------------- | -------------------------------- | ---- | ---- | ---- |
| v7               | 5.2 Maintain Secure Images       |      | \*   | \*   |
| v7               | 5.3 Securely Store Master Images |      | \*   | \*   |

## MITRE ATT&CK Mappings

| Techniques / Sub-techniques | Tactics        | Mitigations  |
| --------------------------- | -------------- | ------------ |
| T1133, T1195                | TA0001, TA0003 | M1016, M1042 |

## Profile

**Level 1** (Manual)
