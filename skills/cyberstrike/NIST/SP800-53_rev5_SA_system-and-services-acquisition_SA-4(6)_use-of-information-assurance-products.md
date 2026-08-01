# stage: recon
# category: NIST


# SA-4(6) Use of Information Assurance Products

> **Enhancement of:** SA-4

## High-Level Description

**Family:** System and Services Acquisition (SA)
**Framework:** NIST SP 800-53 Rev 5

Commercial off-the-shelf IA or IA-enabled information technology products used to protect classified information by cryptographic means may be required to use NSA-approved key management. See [NSA CSFC](#3d575737-98cb-459d-b41c-d7e82b73ad78).

## What to Check

- [ ] Verify SA-4(6) Use of Information Assurance Products is documented in SSP
- [ ] Validate all 2 control requirements are implemented
- [ ] Confirm control is operating effectively
- [ ] Review evidence of continuous monitoring for SA-4(6)
- [ ] Verify enhancement builds upon base control SA-4

## How to Test

### Step 1: Review Documentation

Examine the System Security Plan (SSP) and related artifacts for SA-4(6) implementation details. Verify the organization has documented how this control is satisfied.

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

Employ only government off-the-shelf or commercial off-the-shelf information assurance and information assurance-enabled information technology products that compose an NSA-approved solution to protect classified information when the networks used to transmit the information are at a lower classification level than the information being transmitted; and
Ensure that these products have been evaluated and/or validated by NSA or in accordance with NSA-approved procedures.

### Implementation Guidance

Commercial off-the-shelf IA or IA-enabled information technology products used to protect classified information by cryptographic means may be required to use NSA-approved key management. See [NSA CSFC](#3d575737-98cb-459d-b41c-d7e82b73ad78).

## Risk Assessment

| Finding                                                       | Severity | Impact                                     |
| ------------------------------------------------------------- | -------- | ------------------------------------------ |
| SA-4(6) Use of Information Assurance Products not implemented | Medium   | System and Services Acquisition            |
| SA-4(6) partially implemented                                 | Low      | Incomplete System and Services Acquisition |

## CWE Categories

| CWE ID | Title         |
| ------ | ------------- |
| CWE-16 | Configuration |

## References

- [NIST SP 800-53 Rev 5 - SA-4(6)](https://csrc.nist.gov/projects/cprt/catalog#/cprt/framework/version/SP_800_53_5_1_1/home?element=sa-4.6)
- [NIST SP 800-53A Rev 5 (Assessment Procedures)](https://csrc.nist.gov/pubs/sp/800/53/a/r5/final)
- [NIST SP 800-53 Rev 5 Full Catalog](https://csrc.nist.gov/pubs/sp/800/53/r5/upd1/final)

## Checklist

- [ ] Control documented in SSP
- [ ] Implementation evidence collected
- [ ] Operating effectiveness validated
- [ ] Continuous monitoring in place
- [ ] Related controls (SC-8, SC-12, SC-13) reviewed
