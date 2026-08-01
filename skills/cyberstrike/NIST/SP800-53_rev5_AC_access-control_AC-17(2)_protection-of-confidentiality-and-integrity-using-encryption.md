# stage: report
# category: NIST


# AC-17(2) Protection of Confidentiality and Integrity Using Encryption

> **Enhancement of:** AC-17

## High-Level Description

**Family:** Access Control (AC)
**Framework:** NIST SP 800-53 Rev 5

Virtual private networks can be used to protect the confidentiality and integrity of remote access sessions. Transport Layer Security (TLS) is an example of a cryptographic protocol that provides end-to-end communications security over networks and is used for Internet communications and online transactions.

## What to Check

- [ ] Verify AC-17(2) Protection of Confidentiality and Integrity Using Encryption is documented in SSP
- [ ] Confirm control is operating effectively
- [ ] Review evidence of continuous monitoring for AC-17(2)
- [ ] Verify enhancement builds upon base control AC-17

## How to Test

### Step 1: Review Documentation

Examine the System Security Plan (SSP) and related artifacts for AC-17(2) implementation details. Verify the organization has documented how this control is satisfied.

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

| Tool            | Purpose                                | Usage                      |
| --------------- | -------------------------------------- | -------------------------- |
| cloud-audit-mcp | Check IAM policies and access controls | `cloud_audit_iam_policies` |
| hackbrowser-mcp | Test web application access controls   | `browser_auth_test`        |

## Remediation Guide

### Control Statement

Implement cryptographic mechanisms to protect the confidentiality and integrity of remote access sessions.

### Implementation Guidance

Virtual private networks can be used to protect the confidentiality and integrity of remote access sessions. Transport Layer Security (TLS) is an example of a cryptographic protocol that provides end-to-end communications security over networks and is used for Internet communications and online transactions.

## Risk Assessment

| Finding                                                                               | Severity | Impact                    |
| ------------------------------------------------------------------------------------- | -------- | ------------------------- |
| AC-17(2) Protection of Confidentiality and Integrity Using Encryption not implemented | High     | Access Control            |
| AC-17(2) partially implemented                                                        | Medium   | Incomplete Access Control |

## CWE Categories

| CWE ID  | Title                   |
| ------- | ----------------------- |
| CWE-284 | Improper Access Control |

## References

- [NIST SP 800-53 Rev 5 - AC-17(2)](https://csrc.nist.gov/projects/cprt/catalog#/cprt/framework/version/SP_800_53_5_1_1/home?element=ac-17.2)
- [NIST SP 800-53A Rev 5 (Assessment Procedures)](https://csrc.nist.gov/pubs/sp/800/53/a/r5/final)
- [NIST SP 800-53 Rev 5 Full Catalog](https://csrc.nist.gov/pubs/sp/800/53/r5/upd1/final)

## Checklist

- [ ] Control documented in SSP
- [ ] Implementation evidence collected
- [ ] Operating effectiveness validated
- [ ] Continuous monitoring in place
- [ ] Related controls (SC-8, SC-12, SC-13) reviewed
