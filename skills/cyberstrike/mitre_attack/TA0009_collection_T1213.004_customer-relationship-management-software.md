# stage: post_exploit
# category: mitre_attack


# T1213.004 Customer Relationship Management Software

> **Sub-technique of:** T1213

## High-Level Description

Adversaries may leverage Customer Relationship Management (CRM) software to mine valuable information. CRM software is used to assist organizations in tracking and managing customer interactions, as well as storing customer data.

Once adversaries gain access to a victim organization, they may mine CRM software for customer data. This may include personally identifiable information (PII) such as full names, emails, phone numbers, and addresses, as well as additional details such as purchase histories and IT support interactions. By collecting this data, an adversary may be able to send personalized Phishing emails, engage in SIM swapping, or otherwise target the organization’s customers in ways that enable financial gain or the compromise of additional organizations.

CRM software may be hosted on-premises or in the cloud. Information stored in these solutions may vary based on the specific instance or environment. Examples of CRM software include Microsoft Dynamics 365, Salesforce, Zoho, Zendesk, and HubSpot.

## Kill Chain Phase

- Collection (TA0009)

**Platforms:** SaaS

## What to Check

- [ ] Identify if Customer Relationship Management Software technique is applicable to target environment
- [ ] Check SaaS systems for indicators of Customer Relationship Management Software
- [ ] Verify mitigations are bypassed or absent (4 known mitigations)
- [ ] Assess detection coverage (1 detection strategies)

## How to Test

### Manual Testing

1. **Identify Attack Surface**: Determine if the target environment is susceptible to Customer Relationship Management Software by examining the target platforms (SaaS).

2. **Assess Existing Defenses**: Review whether mitigations for T1213.004 are in place. If defenses are absent or misconfigured, this technique may be exploitable.

3. **Execute Test**: Use tools and methods described in the MITRE ATT&CK page and external references below.

> **Note**: No Atomic Red Team tests available for this technique. See [Atomic Red Team GitHub](https://github.com/redcanaryco/atomic-red-team) for updates.

## Remediation Guide

### M1018 User Account Management

Enforce the principle of least-privilege. Consider implementing access control mechanisms that include both authentication and authorization.

### M1017 User Training

Develop and publish policies that define acceptable information to be stored in CRM databases and acceptable handling of customer data. Only store customer information required for business operations.

### M1054 Software Configuration

Consider implementing data retention policies to automate periodically archiving and/or deleting data that is no longer needed.

### M1047 Audit

Consider periodic review of accounts and privileges for critical and sensitive CRM data.

## Detection

### Detecting Suspicious Access to CRM Data in SaaS Environments

## Risk Assessment

| Finding                                                        | Severity | Impact     |
| -------------------------------------------------------------- | -------- | ---------- |
| Customer Relationship Management Software technique applicable | High     | Collection |

## CWE Categories

| CWE ID  | Title                             |
| ------- | --------------------------------- |
| CWE-200 | Exposure of Sensitive Information |

## References

- [Bleeping Computer Bank Hack 2020](https://www.bleepingcomputer.com/news/security/customer-owned-bank-informs-100k-of-breach-exposing-account-balance-pii/)
- [Bleeping Computer Mint Mobile Hack 2021](https://www.bleepingcomputer.com/news/security/mint-mobile-hit-by-a-data-breach-after-numbers-ported-data-accessed/)
- [Bleeping Computer US Cellular Hack 2022](https://www.bleepingcomputer.com/news/security/uscellular-discloses-data-breach-after-billing-system-hack/)
- [Atomic Red Team - T1213.004](https://github.com/redcanaryco/atomic-red-team/tree/master/atomics/T1213.004)
- [MITRE ATT&CK - T1213.004](https://attack.mitre.org/techniques/T1213/004)
