# stage: report
# category: NIST


# MA-3(4) Restricted Tool Use

> **Enhancement of:** MA-3

## High-Level Description

**Family:** Maintenance (MA)
**Framework:** NIST SP 800-53 Rev 5

Restricting the use of maintenance tools to only authorized personnel applies to systems that are used to carry out maintenance functions.

## What to Check

- [ ] Verify MA-3(4) Restricted Tool Use is documented in SSP
- [ ] Confirm control is operating effectively
- [ ] Review evidence of continuous monitoring for MA-3(4)
- [ ] Verify enhancement builds upon base control MA-3

## How to Test

### Step 1: Review Documentation

Examine the System Security Plan (SSP) and related artifacts for MA-3(4) implementation details. Verify the organization has documented how this control is satisfied.

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

| Tool          | Purpose                           | Usage |
| ------------- | --------------------------------- | ----- |
| Manual Review | Documentation and interview-based | N/A   |

## Remediation Guide

### Control Statement

Restrict the use of maintenance tools to authorized personnel only.

### Implementation Guidance

Restricting the use of maintenance tools to only authorized personnel applies to systems that are used to carry out maintenance functions.

## Risk Assessment

| Finding                                     | Severity | Impact                 |
| ------------------------------------------- | -------- | ---------------------- |
| MA-3(4) Restricted Tool Use not implemented | Medium   | Maintenance            |
| MA-3(4) partially implemented               | Low      | Incomplete Maintenance |

## CWE Categories

| CWE ID | Title                 |
| ------ | --------------------- |
| N/A    | No direct CWE mapping |

## References

- [NIST SP 800-53 Rev 5 - MA-3(4)](https://csrc.nist.gov/projects/cprt/catalog#/cprt/framework/version/SP_800_53_5_1_1/home?element=ma-3.4)
- [NIST SP 800-53A Rev 5 (Assessment Procedures)](https://csrc.nist.gov/pubs/sp/800/53/a/r5/final)
- [NIST SP 800-53 Rev 5 Full Catalog](https://csrc.nist.gov/pubs/sp/800/53/r5/upd1/final)

## Checklist

- [ ] Control documented in SSP
- [ ] Implementation evidence collected
- [ ] Operating effectiveness validated
- [ ] Continuous monitoring in place
- [ ] Related controls (AC-3, AC-5, AC-6) reviewed
