# stage: report
# category: CIS_benchmarks


# 5.4 Ensure root login is restricted to system console (Manual)

## Description

The file `/etc/securetty` contains a list of valid terminals that may be logged in directly as root.

## Rationale

Since the system console has special properties to handle emergency situations, it is important to ensure that the console is in a physically secure location and that unauthorized consoles have not been defined.

## Audit Procedure

```bash
# cat /etc/securetty
```

## Expected Result

Only consoles that are in a physically secure location should be listed in `/etc/securetty`.

## Remediation

Remove entries for any consoles that are not in a physically secure location.

## CIS Controls

| Controls Version | Control                                                                                                                                                                                                                                                                                                                          | IG 1 | IG 2 | IG 3 |
| ---------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---- | ---- | ---- |
| v8               | **5.4 Restrict Administrator Privileges to Dedicated Administrator Accounts** - Restrict administrator privileges to dedicated administrator accounts on enterprise assets. Conduct general computing activities, such as internet browsing, email, and productivity suite use, from the user's primary, non-privileged account. | x    | x    | x    |
| v7               | **5.1 Establish Secure Configurations** - Maintain documented, standard security configuration standards for all authorized operating systems and software.                                                                                                                                                                      | x    | x    | x    |

## Profile

- Level 1 - Server
