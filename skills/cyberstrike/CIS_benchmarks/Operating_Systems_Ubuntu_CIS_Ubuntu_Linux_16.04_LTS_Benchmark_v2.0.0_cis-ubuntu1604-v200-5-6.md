# stage: report
# category: CIS_benchmarks


# Ensure root login is restricted to system console

**Profile Applicability:**

- Level 1 - Server
- Level 1 - Workstation

**Assessment Status:** Manual

## Description

The file `/etc/securetty` contains a list of valid terminals that may be logged in directly as root.

## Rationale

Since the system console has special properties to handle emergency situations, it is important to ensure that the console is in a physically secure location and that unauthorized consoles have not been defined.

## Audit Procedure

### Command Line

```bash
cat /etc/securetty
```

## Remediation

### Command Line

Remove entries for any consoles that are not in a physically secure location.

## References

None

## CIS Controls

| Controls Version | Control                                                                                                                                                                                                                                                                                                  |
| ---------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| v7               | 4.3 Ensure the Use of Dedicated Administrative Accounts - Ensure that all users with administrative account access use a dedicated or secondary account for elevated activities. This account should only be used for administrative activities and not internet browsing, email, or similar activities. |

## Profile

- **Level 1 - Server**
- **Level 1 - Workstation**
