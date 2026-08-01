# stage: report
# category: CIS_benchmarks


# 9.4 Restrict root Login to System Console (Not Scored)

## Profile Applicability

- Level 1

## Description

The file `/etc/securetty` contains a list of valid terminals that may be logged in directly as root.

## Rationale

Since the system console has special properties to handle emergency situations, it is important to ensure that the console is in a physically secure location and that unauthorized consoles have not been defined.

## Audit Procedure

### Using Command Line

```bash
cat /etc/securetty
```

## Expected Result

Only consoles that are in a physically secure location should be listed.

## Remediation

### Using Command Line

Remove entries for any consoles that are not in a physically secure location.

## Default Value

By default, /etc/securetty contains a list of virtual consoles.

## References

- CIS Ubuntu 12.04 LTS Server Benchmark v1.1.0

## Profile

Level 1 - Not Scored
