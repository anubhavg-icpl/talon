# stage: report
# category: CIS_benchmarks


# 2.2.2 Ensure X Window System is not installed (Scored)

## Profile Applicability

- Level 1 - Server

## Description

The X Window System provides a Graphical User Interface (GUI) where users can have multiple windows in which to run programs and various add on. The X Windows system is typically used on workstations where users login, but not on servers where users typically do not login.

## Rationale

Unless your organization specifically requires graphical login access via X Windows, remove it to reduce the potential attack surface.

## Audit Procedure

Run the following command and verify X Windows System is not installed:

```bash
dpkg -l xserver-xorg*
```

## Expected Result

No xserver-xorg packages should be installed.

## Remediation

Run the following command to remove the X Windows System packages:

```bash
apt-get remove xserver-xorg*
```

## Default Value

X Window System may be installed on workstation installations.

## References

- CIS Controls: 2 Inventory of Authorized and Unauthorized Software

## Profile

- Level 1 - Server
