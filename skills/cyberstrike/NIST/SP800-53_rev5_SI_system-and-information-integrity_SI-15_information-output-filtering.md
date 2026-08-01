# stage: recon
# category: NIST


# SI-15 Information Output Filtering

## High-Level Description

**Family:** System and Information Integrity (SI)
**Framework:** NIST SP 800-53 Rev 5

Certain types of attacks, including SQL injections, produce output results that are unexpected or inconsistent with the output results that would be expected from software programs or applications. Information output filtering focuses on detecting extraneous content, preventing such extraneous content from being displayed, and then alerting monitoring tools that anomalous behavior has been discovered.

## What to Check

- [ ] Verify SI-15 Information Output Filtering is documented in SSP
- [ ] Confirm control is operating effectively
- [ ] Review evidence of continuous monitoring for SI-15

## How to Test

### Step 1: Review Documentation

Examine the System Security Plan (SSP) and related artifacts for SI-15 implementation details. Verify the organization has documented how this control is satisfied.

### Step 2: Validate Implementation

```
# For cloud environments, use cloud-audit-mcp tools
# For on-premises, review system configurations directly

# Example: Check if account management policies exist
grep -r "account.management\|access.control" /etc/security/ 2>/dev/null
```

### Step 3: Test Operating Effectiveness

Verify the control is actively functioning, not just documented. Check logs, configurations, and operational evidence.

## Tools

| Tool            | Purpose                    | Usage                          |
| --------------- | -------------------------- | ------------------------------ |
| cloud-audit-mcp | Check integrity monitoring | `cloud_audit_monitoring`       |
| AWS CLI         | Review GuardDuty/Inspector | `aws guardduty list-detectors` |

## Remediation Guide

### Control Statement

Validate information output from the following software programs and/or applications to ensure that the information is consistent with the expected content: [organization-defined].

### Implementation Guidance

Certain types of attacks, including SQL injections, produce output results that are unexpected or inconsistent with the output results that would be expected from software programs or applications. Information output filtering focuses on detecting extraneous content, preventing such extraneous content from being displayed, and then alerting monitoring tools that anomalous behavior has been discovered.

## Risk Assessment

| Finding                                            | Severity | Impact                                      |
| -------------------------------------------------- | -------- | ------------------------------------------- |
| SI-15 Information Output Filtering not implemented | High     | System and Information Integrity            |
| SI-15 partially implemented                        | Medium   | Incomplete System and Information Integrity |

## CWE Categories

| CWE ID | Title                     |
| ------ | ------------------------- |
| CWE-20 | Improper Input Validation |

## References

- [NIST SP 800-53 Rev 5 - SI-15](https://csrc.nist.gov/projects/cprt/catalog#/cprt/framework/version/SP_800_53_5_1_1/home?element=si-15)
- [NIST SP 800-53A Rev 5 (Assessment Procedures)](https://csrc.nist.gov/pubs/sp/800/53/a/r5/final)
- [NIST SP 800-53 Rev 5 Full Catalog](https://csrc.nist.gov/pubs/sp/800/53/r5/upd1/final)

## Checklist

- [ ] Control documented in SSP
- [ ] Implementation evidence collected
- [ ] Operating effectiveness validated
- [ ] Continuous monitoring in place
- [ ] Related controls (SI-3, SI-4, SI-11) reviewed
