# stage: report
# category: CIS_benchmarks


# 1.2.2 Ensure GPG keys are configured (Not Scored)

## Profile Applicability

- Level 1 - Server
- Level 1 - Workstation

## Description

Most packages managers implement GPG key signing to verify package integrity during installation.

## Rationale

It is important to ensure that updates are obtained from a valid source to protect against spoofing that could lead to the inadvertent installation of malware on the system.

## Audit Procedure

Run the following command and verify GPG keys are configured correctly for your package manager:

```bash
apt-key list
```

## Expected Result

Verify that GPG keys are configured according to site policy.

## Remediation

Update your package manager GPG keys in accordance with site policy.

## Default Value

Not applicable. Configuration varies by site policy.

## References

- CIS Controls: 4.5 Use Automated Patch Management And Software Update Tools

## Profile

- Level 1 - Server
- Level 1 - Workstation
