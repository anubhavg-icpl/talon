# stage: report
# category: CIS_benchmarks


# 5.5 Ensure root login is restricted to system console (Not Scored)

## Profile Applicability

- Level 1 - Server
- Level 1 - Workstation

## Description

The file `/etc/securetty` contains a list of valid terminals that may be logged in directly as root.

## Rationale

Since the system console has special properties to handle emergency situations, it is important to ensure that the console is in a physically secure location and that unauthorized consoles have not been defined.

## Audit Procedure

```bash
cat /etc/securetty
```

## Expected Result

Review the output and verify only physically secure consoles are listed.

## Remediation

Remove entries for any consoles that are not in a physically secure location.

## Default Value

Varies by installation.

## References

- CIS Controls: 5.1 - Minimize And Sparingly Use Administrative Privileges

## Profile

- Level 1 - Server
- Level 1 - Workstation
