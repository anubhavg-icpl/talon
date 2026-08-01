# stage: report
# category: NIST


# CP-2(2) Capacity Planning

> **Enhancement of:** CP-2

## High-Level Description

**Family:** Contingency Planning (CP)
**Framework:** NIST SP 800-53 Rev 5

Capacity planning is needed because different threats can result in a reduction of the available processing, telecommunications, and support services intended to support essential mission and business functions. Organizations anticipate degraded operations during contingency operations and factor the degradation into capacity planning. For capacity planning, environmental support refers to any environmental factor for which the organization determines that it needs to provide support in a contingency situation, even if in a degraded state. Such determinations are based on an organizational assessment of risk, system categorization (impact level), and organizational risk tolerance.

## What to Check

- [ ] Verify CP-2(2) Capacity Planning is documented in SSP
- [ ] Confirm control is operating effectively
- [ ] Review evidence of continuous monitoring for CP-2(2)
- [ ] Verify enhancement builds upon base control CP-2

## How to Test

### Step 1: Review Documentation

Examine the System Security Plan (SSP) and related artifacts for CP-2(2) implementation details. Verify the organization has documented how this control is satisfied.

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

Conduct capacity planning so that necessary capacity for information processing, telecommunications, and environmental support exists during contingency operations.

### Implementation Guidance

Capacity planning is needed because different threats can result in a reduction of the available processing, telecommunications, and support services intended to support essential mission and business functions. Organizations anticipate degraded operations during contingency operations and factor the degradation into capacity planning. For capacity planning, environmental support refers to any environmental factor for which the organization determines that it needs to provide support in a contingency situation, even if in a degraded state. Such determinations are based on an organizational assessment of risk, system categorization (impact level), and organizational risk tolerance.

## Risk Assessment

| Finding                                   | Severity | Impact                          |
| ----------------------------------------- | -------- | ------------------------------- |
| CP-2(2) Capacity Planning not implemented | Medium   | Contingency Planning            |
| CP-2(2) partially implemented             | Low      | Incomplete Contingency Planning |

## CWE Categories

| CWE ID | Title                 |
| ------ | --------------------- |
| N/A    | No direct CWE mapping |

## References

- [NIST SP 800-53 Rev 5 - CP-2(2)](https://csrc.nist.gov/projects/cprt/catalog#/cprt/framework/version/SP_800_53_5_1_1/home?element=cp-2.2)
- [NIST SP 800-53A Rev 5 (Assessment Procedures)](https://csrc.nist.gov/pubs/sp/800/53/a/r5/final)
- [NIST SP 800-53 Rev 5 Full Catalog](https://csrc.nist.gov/pubs/sp/800/53/r5/upd1/final)

## Checklist

- [ ] Control documented in SSP
- [ ] Implementation evidence collected
- [ ] Operating effectiveness validated
- [ ] Continuous monitoring in place
- [ ] Related controls (PE-11, PE-12, PE-13, PE-14, PE-18) reviewed
