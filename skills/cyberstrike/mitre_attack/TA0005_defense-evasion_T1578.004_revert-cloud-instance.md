# stage: post_exploit
# category: mitre_attack


# T1578.004 Revert Cloud Instance

> **Sub-technique of:** T1578

## High-Level Description

An adversary may revert changes made to a cloud instance after they have performed malicious activities in attempt to evade detection and remove evidence of their presence. In highly virtualized environments, such as cloud-based infrastructure, this may be accomplished by restoring virtual machine (VM) or data storage snapshots through the cloud management dashboard or cloud APIs.

Another variation of this technique is to utilize temporary storage attached to the compute instance. Most cloud providers provide various types of storage including persistent, local, and/or ephemeral, with the ephemeral types often reset upon stop/restart of the VM.

## Kill Chain Phase

- Defense Evasion (TA0005)

**Platforms:** IaaS

## What to Check

- [ ] Identify if Revert Cloud Instance technique is applicable to target environment
- [ ] Check IaaS systems for indicators of Revert Cloud Instance
- [ ] Assess detection coverage (1 detection strategies)

## How to Test

### Manual Testing

1. **Identify Attack Surface**: Determine if the target environment is susceptible to Revert Cloud Instance by examining the target platforms (IaaS).

2. **Assess Existing Defenses**: Review whether mitigations for T1578.004 are in place. If defenses are absent or misconfigured, this technique may be exploitable.

3. **Execute Test**: Use tools and methods described in the MITRE ATT&CK page and external references below.

> **Note**: No Atomic Red Team tests available for this technique. See [Atomic Red Team GitHub](https://github.com/redcanaryco/atomic-red-team) for updates.

## Remediation Guide

No specific mitigations documented for this technique.

## Detection

### Detection Strategy for Modify Cloud Compute Infrastructure: Revert Cloud Instance

## Risk Assessment

| Finding                                    | Severity | Impact          |
| ------------------------------------------ | -------- | --------------- |
| Revert Cloud Instance technique applicable | Low      | Defense Evasion |

## CWE Categories

| CWE ID  | Title                        |
| ------- | ---------------------------- |
| CWE-693 | Protection Mechanism Failure |

## References

- [Tech Republic - Restore AWS Snapshots](https://www.techrepublic.com/blog/the-enterprise-cloud/backing-up-and-restoring-snapshots-on-amazon-ec2-machines/)
- [Google - Restore Cloud Snapshot](https://cloud.google.com/compute/docs/disks/restore-and-delete-snapshots)
- [Atomic Red Team - T1578.004](https://github.com/redcanaryco/atomic-red-team/tree/master/atomics/T1578.004)
- [MITRE ATT&CK - T1578.004](https://attack.mitre.org/techniques/T1578/004)
