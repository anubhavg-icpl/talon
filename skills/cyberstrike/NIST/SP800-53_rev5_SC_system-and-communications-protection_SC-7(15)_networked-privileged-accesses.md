# stage: report
# category: NIST


# SC-7(15) Networked Privileged Accesses

> **Enhancement of:** SC-7

## High-Level Description

**Family:** System and Communications Protection (SC)
**Framework:** NIST SP 800-53 Rev 5

Privileged access provides greater accessibility to system functions, including security functions. Adversaries attempt to gain privileged access to systems through remote access to cause adverse mission or business impacts, such as by exfiltrating information or bringing down a critical system capability. Routing networked, privileged access requests through a dedicated, managed interface further restricts privileged access for increased access control and auditing.

## What to Check

- [ ] Verify SC-7(15) Networked Privileged Accesses is documented in SSP
- [ ] Confirm control is operating effectively
- [ ] Review evidence of continuous monitoring for SC-7(15)
- [ ] Verify enhancement builds upon base control SC-7

## How to Test

### Step 1: Review Documentation

Examine the System Security Plan (SSP) and related artifacts for SC-7(15) implementation details. Verify the organization has documented how this control is satisfied.

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

| Tool            | Purpose                               | Usage                                |
| --------------- | ------------------------------------- | ------------------------------------ |
| cloud-audit-mcp | Check encryption and network controls | `cloud_audit_encryption`             |
| nmap            | Network scanning                      | `nmap -sV --script ssl-enum-ciphers` |

## Remediation Guide

### Control Statement

Route networked, privileged accesses through a dedicated, managed interface for purposes of access control and auditing.

### Implementation Guidance

Privileged access provides greater accessibility to system functions, including security functions. Adversaries attempt to gain privileged access to systems through remote access to cause adverse mission or business impacts, such as by exfiltrating information or bringing down a critical system capability. Routing networked, privileged access requests through a dedicated, managed interface further restricts privileged access for increased access control and auditing.

## Risk Assessment

| Finding                                                | Severity | Impact                                          |
| ------------------------------------------------------ | -------- | ----------------------------------------------- |
| SC-7(15) Networked Privileged Accesses not implemented | High     | System and Communications Protection            |
| SC-7(15) partially implemented                         | Medium   | Incomplete System and Communications Protection |

## CWE Categories

| CWE ID  | Title                                |
| ------- | ------------------------------------ |
| CWE-311 | Missing Encryption of Sensitive Data |

## References

- [NIST SP 800-53 Rev 5 - SC-7(15)](https://csrc.nist.gov/projects/cprt/catalog#/cprt/framework/version/SP_800_53_5_1_1/home?element=sc-7.15)
- [NIST SP 800-53A Rev 5 (Assessment Procedures)](https://csrc.nist.gov/pubs/sp/800/53/a/r5/final)
- [NIST SP 800-53 Rev 5 Full Catalog](https://csrc.nist.gov/pubs/sp/800/53/r5/upd1/final)

## Checklist

- [ ] Control documented in SSP
- [ ] Implementation evidence collected
- [ ] Operating effectiveness validated
- [ ] Continuous monitoring in place
- [ ] Related controls (AC-2, AC-3, AU-2, SI-4) reviewed
