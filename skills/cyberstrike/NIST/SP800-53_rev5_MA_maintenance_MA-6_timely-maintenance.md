# stage: report
# category: NIST


# MA-6 Timely Maintenance

## High-Level Description

**Family:** Maintenance (MA)
**Framework:** NIST SP 800-53 Rev 5

Organizations specify the system components that result in increased risk to organizational operations and assets, individuals, other organizations, or the Nation when the functionality provided by those components is not operational. Organizational actions to obtain maintenance support include having appropriate contracts in place.

## What to Check

- [ ] Verify MA-6 Timely Maintenance is documented in SSP
- [ ] Confirm control is operating effectively
- [ ] Review evidence of continuous monitoring for MA-6

## How to Test

### Step 1: Review Documentation

Examine the System Security Plan (SSP) and related artifacts for MA-6 implementation details. Verify the organization has documented how this control is satisfied.

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

Obtain maintenance support and/or spare parts for [organization-defined] within [organization-defined] of failure.

### Implementation Guidance

Organizations specify the system components that result in increased risk to organizational operations and assets, individuals, other organizations, or the Nation when the functionality provided by those components is not operational. Organizational actions to obtain maintenance support include having appropriate contracts in place.

## Risk Assessment

| Finding                                 | Severity | Impact                 |
| --------------------------------------- | -------- | ---------------------- |
| MA-6 Timely Maintenance not implemented | Medium   | Maintenance            |
| MA-6 partially implemented              | Low      | Incomplete Maintenance |

## CWE Categories

| CWE ID | Title                 |
| ------ | --------------------- |
| N/A    | No direct CWE mapping |

## References

- [NIST SP 800-53 Rev 5 - MA-6](https://csrc.nist.gov/projects/cprt/catalog#/cprt/framework/version/SP_800_53_5_1_1/home?element=ma-6)
- [NIST SP 800-53A Rev 5 (Assessment Procedures)](https://csrc.nist.gov/pubs/sp/800/53/a/r5/final)
- [NIST SP 800-53 Rev 5 Full Catalog](https://csrc.nist.gov/pubs/sp/800/53/r5/upd1/final)

## Checklist

- [ ] Control documented in SSP
- [ ] Implementation evidence collected
- [ ] Operating effectiveness validated
- [ ] Continuous monitoring in place
- [ ] Related controls (CM-8, CP-2, CP-7, RA-7, SA-15) reviewed
