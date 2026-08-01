# stage: report
# category: CIS_benchmarks


# 6.1 Ensure the X Window system is not installed (Scored)

## Profile Applicability

- Level 1

## Description

The X Window system provides a Graphical User Interface (GUI) where users can have multiple windows in which to run programs and various add on. The X Window system is typically used on desktops where users login, but not on servers where users typically do not login.

## Rationale

Unless your organization specifically requires graphical login access via X Windows, remove it to reduce the potential attack surface.

## Audit Procedure

### Using Command Line

Run the following command:

```bash
dpkg -l xserver-xorg-core*
```

## Expected Result

Ensure no matching packages are listed as installed.

## Remediation

### Using Command Line

Uninstall X Windows:

```bash
apt-get purge xserver-xorg-core*
```

## Default Value

Not installed by default on Ubuntu 12.04 LTS Server.

## References

- CIS Ubuntu 12.04 LTS Server Benchmark v1.1.0

## Profile

Level 1 - Scored
