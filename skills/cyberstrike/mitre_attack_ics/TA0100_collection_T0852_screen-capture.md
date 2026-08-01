# stage: post_exploit
# category: mitre_attack_ics


# T0852 Screen Capture

## High-Level Description

Adversaries may attempt to perform screen capture of devices in the control system environment. Screenshots may be taken of workstations, HMIs, or other devices that display environment-relevant process, device, reporting, alarm, or related data. These device displays may reveal information regarding the ICS process, layout, control, and related schematics. In particular, an HMI can provide a lot of important industrial process information. Analysis of screen captures may provide the adversary with an understanding of intended operations and interactions between critical devices.

## Kill Chain Phase

- Collection (TA0100)

**Platforms:** ICS

## What to Check

- [ ] Identify if Screen Capture technique is applicable to target ICS environment
- [ ] Check ICS/SCADA systems for indicators of Screen Capture
- [ ] Verify mitigations are bypassed or absent (1 known mitigations)
- [ ] Assess detection coverage (1 detection strategies)

## How to Test

### Identify Attack Surface

Determine if the target ICS/SCADA environment is susceptible to Screen Capture by examining operational technology systems and network architecture.

### Assess Existing Defenses

Review whether mitigations for T0852 are in place. If defenses are absent or misconfigured, this technique may be exploitable in the ICS environment.

## Remediation Guide

### M0816 Mitigation Limited or Not Effective

Preventing screen capture on a device may require disabling various system calls supported by the operating systems (e.g., Microsoft WindowsGraphicsCaputer APIs), however, these may be needed for other critical applications.

## Detection

### Detection of Screen Capture

## Risk Assessment

| Finding                             | Severity | Impact     |
| ----------------------------------- | -------- | ---------- |
| Screen Capture technique applicable | Low      | Collection |

## CWE Categories

| CWE ID  | Title                             |
| ------- | --------------------------------- |
| CWE-200 | Exposure of Sensitive Information |

## References

- [ICS-CERT October 2017](https://www.us-cert.gov/ncas/alerts/TA17-293A)
- [MITRE ATT&CK ICS - T0852](https://attack.mitre.org/techniques/T0852)
