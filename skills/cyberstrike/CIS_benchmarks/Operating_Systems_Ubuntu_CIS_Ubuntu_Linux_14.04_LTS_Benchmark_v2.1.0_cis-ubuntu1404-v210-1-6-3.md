# stage: report
# category: CIS_benchmarks


# 1.6.3 Ensure SELinux or AppArmor are installed (Scored)

## Profile Applicability

- Level 2 - Server
- Level 2 - Workstation

## Description

SELinux and AppArmor provide Mandatory Access Controls.

## Rationale

Without a Mandatory Access Control system installed only the default Discretionary Access Control system will be available.

## Audit Procedure

Run the following commands and verify either SELinux or AppArmor is installed:

```bash
dpkg -s selinux
dpkg -s apparmor
```

## Expected Result

At least one of the packages (`selinux` or `apparmor`) should be installed with status `install ok installed`.

## Remediation

Run one of the following commands to install SELinux or apparmor:

```bash
apt-get install selinux
```

or

```bash
apt-get install apparmor
```

## Default Value

AppArmor is installed by default on Ubuntu systems.

## References

- CIS Controls: 14.4 Protect Information With Access Control Lists

## Profile

- Level 2 - Server
- Level 2 - Workstation
