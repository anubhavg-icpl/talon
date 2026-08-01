# stage: report
# category: CIS_benchmarks


# 1.6 Ensure AppArmor is installed (Automated)

## Description

SELinux and AppArmor provide Mandatory Access Controls.

## Rationale

Without a Mandatory Access Control system installed only the default Discretionary Access Control system will be available.

## Audit Procedure

Verify that AppArmor is installed with the following command:

```bash
# grep '\"name\": \"apparmor\"' /etc/cos-package-info.json
"name": "apparmor",
```

## Expected Result

The output should show `"name": "apparmor",` confirming AppArmor is installed.

## Remediation

Update to an OS image that includes the apparmor package.

## Additional Information

SELinux and AppArmor both have several package names in use on different distributions. Research the appropriate packages for your environment.

## CIS Controls

| Controls Version | Control                                               | IG 1 | IG 2 | IG 3 |
| ---------------- | ----------------------------------------------------- | ---- | ---- | ---- |
| v8               | 3.3 Configure Data Access Control Lists               | x    | x    | x    |
| v7               | 14.6 Protect Information through Access Control Lists | x    | x    | x    |

## Profile

Level 1 - Server | Automated
