# stage: post_exploit
# category: mitre_attack_mobile


# T1636.005 Accounts

> **Sub-technique of:** T1636

## High-Level Description

Adversaries may utilize standard operating system APIs to gather account data. On Android, this can be accomplished by using the AccountManager API. For example, adversaries may use the `getAccounts()` method to list all accounts. On iOS, this can be accomplished by using the Keychain services.

If the device has been jailbroken or rooted, adversaries may be able to access Accounts without the users’ knowledge or approval.

## Kill Chain Phase

- Collection (TA0035)

**Platforms:** Android, iOS

## What to Check

- [ ] Identify if Accounts technique is applicable to target mobile environment
- [ ] Check Android devices for indicators of Accounts
- [ ] Check iOS devices for indicators of Accounts
- [ ] Verify mitigations are bypassed or absent (2 known mitigations)
- [ ] Assess detection coverage (1 detection strategies)

## How to Test

### Identify Attack Surface

Determine if the target mobile environment is susceptible to Accounts by examining the target platforms (Android, iOS).

### Assess Existing Defenses

Review whether mitigations for T1636.005 are in place. If defenses are absent or misconfigured, this technique may be exploitable.

## Remediation Guide

### M1006 Use Recent OS Version

OS feature updates often enhance security and privacy around permissions.

### M1011 User Guidance

Access to accounts is an uncommonly needed permission, so users should be instructed to use extra scrutiny when granting access to their accounts.

## Detection

### Detection of Accounts

## Risk Assessment

| Finding                       | Severity | Impact     |
| ----------------------------- | -------- | ---------- |
| Accounts technique applicable | Low      | Collection |

## CWE Categories

| CWE ID  | Title                             |
| ------- | --------------------------------- |
| CWE-200 | Exposure of Sensitive Information |

## References

- [Android_AccountManager_Feb2025](https://developer.android.com/reference/android/accounts/AccountManager)
- [MITRE ATT&CK Mobile - T1636.005](https://attack.mitre.org/techniques/T1636/005)
