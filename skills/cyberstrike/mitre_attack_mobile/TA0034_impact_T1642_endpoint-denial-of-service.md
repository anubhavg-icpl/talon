# stage: post_exploit
# category: mitre_attack_mobile


# T1642 Endpoint Denial of Service

## High-Level Description

Adversaries may perform Endpoint Denial of Service (DoS) attacks to degrade or block the availability of services to users.

On Android versions prior to 7, apps can abuse Device Administrator access to reset the device lock passcode, preventing the user from unlocking the device. After Android 7, only device or profile owners (e.g. MDMs) can reset the device’s passcode.

On iOS devices, this technique does not work because mobile device management servers can only remove the screen lock passcode; they cannot set a new passcode. However, on jailbroken devices, malware has been discovered that can lock the user out of the device.

## Kill Chain Phase

- Impact (TA0034)

**Platforms:** Android, iOS

## What to Check

- [ ] Identify if Endpoint Denial of Service technique is applicable to target mobile environment
- [ ] Check Android devices for indicators of Endpoint Denial of Service
- [ ] Check iOS devices for indicators of Endpoint Denial of Service
- [ ] Verify mitigations are bypassed or absent (2 known mitigations)
- [ ] Assess detection coverage (1 detection strategies)

## How to Test

### Identify Attack Surface

Determine if the target mobile environment is susceptible to Endpoint Denial of Service by examining the target platforms (Android, iOS).

### Assess Existing Defenses

Review whether mitigations for T1642 are in place. If defenses are absent or misconfigured, this technique may be exploitable.

## Remediation Guide

### M1006 Use Recent OS Version

Android 7 changed how the Device Administrator password APIs function.

### M1011 User Guidance

Users should be cautioned against granting administrative access to applications.

## Detection

### Detection of Endpoint Denial of Service

## Risk Assessment

| Finding                                         | Severity | Impact |
| ----------------------------------------------- | -------- | ------ |
| Endpoint Denial of Service technique applicable | Low      | Impact |

## CWE Categories

| CWE ID  | Title                             |
| ------- | --------------------------------- |
| CWE-400 | Uncontrolled Resource Consumption |

## References

- [Xiao-KeyRaider](http://researchcenter.paloaltonetworks.com/2015/08/keyraider-ios-malware-steals-over-225000-apple-accounts-to-create-free-app-utopia/)
- [Android resetPassword](<https://developer.android.com/reference/android/app/admin/DevicePolicyManager.html#resetPassword(java.lang.String,%20int)>)
- [MITRE ATT&CK Mobile - T1642](https://attack.mitre.org/techniques/T1642)
