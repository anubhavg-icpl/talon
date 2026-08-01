# stage: report
# category: NIST


# PE-23 Facility Location

## High-Level Description

**Family:** Physical and Environmental Protection (PE)
**Framework:** NIST SP 800-53 Rev 5

Physical and environmental hazards include floods, fires, tornadoes, earthquakes, hurricanes, terrorism, vandalism, an electromagnetic pulse, electrical interference, and other forms of incoming electromagnetic radiation. The location of system components within the facility is addressed in [PE-18](#pe-18).

## What to Check

- [ ] Verify PE-23 Facility Location is documented in SSP
- [ ] Validate all 2 control requirements are implemented
- [ ] Confirm control is operating effectively
- [ ] Review evidence of continuous monitoring for PE-23

## How to Test

### Step 1: Review Documentation

Examine the System Security Plan (SSP) and related artifacts for PE-23 implementation details. Verify the organization has documented how this control is satisfied.

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

Plan the location or site of the facility where the system resides considering physical and environmental hazards; and
For existing facilities, consider the physical and environmental hazards in the organizational risk management strategy.

### Implementation Guidance

Physical and environmental hazards include floods, fires, tornadoes, earthquakes, hurricanes, terrorism, vandalism, an electromagnetic pulse, electrical interference, and other forms of incoming electromagnetic radiation. The location of system components within the facility is addressed in [PE-18](#pe-18).

## Risk Assessment

| Finding                                 | Severity | Impact                                           |
| --------------------------------------- | -------- | ------------------------------------------------ |
| PE-23 Facility Location not implemented | Medium   | Physical and Environmental Protection            |
| PE-23 partially implemented             | Low      | Incomplete Physical and Environmental Protection |

## CWE Categories

| CWE ID | Title                 |
| ------ | --------------------- |
| N/A    | No direct CWE mapping |

## References

- [NIST SP 800-53 Rev 5 - PE-23](https://csrc.nist.gov/projects/cprt/catalog#/cprt/framework/version/SP_800_53_5_1_1/home?element=pe-23)
- [NIST SP 800-53A Rev 5 (Assessment Procedures)](https://csrc.nist.gov/pubs/sp/800/53/a/r5/final)
- [NIST SP 800-53 Rev 5 Full Catalog](https://csrc.nist.gov/pubs/sp/800/53/r5/upd1/final)

## Checklist

- [ ] Control documented in SSP
- [ ] Implementation evidence collected
- [ ] Operating effectiveness validated
- [ ] Continuous monitoring in place
- [ ] Related controls (CP-2, PE-18, PE-19, PM-8, PM-9) reviewed
