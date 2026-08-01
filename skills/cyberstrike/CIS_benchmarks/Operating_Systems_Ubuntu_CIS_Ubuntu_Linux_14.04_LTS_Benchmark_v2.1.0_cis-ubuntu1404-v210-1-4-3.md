# stage: report
# category: CIS_benchmarks


# 1.4.3 Ensure authentication required for single user mode (Scored)

## Profile Applicability

- Level 1 - Server
- Level 1 - Workstation

## Description

Single user mode is used for recovery when the system detects an issue during boot or by manual selection from the bootloader.

## Rationale

Requiring authentication in single user mode prevents an unauthorized user from rebooting the system into single user to gain root privileges without credentials.

## Audit Procedure

Perform the following to determine if a password is set for the root user:

```bash
grep ^root:[*\!]: /etc/shadow
```

No results should be returned.

## Expected Result

No output should be returned. If output is returned, the root account does not have a password set.

## Remediation

Run the following command and follow the prompts to set a password for the root user:

```bash
passwd root
```

## Default Value

Not applicable.

## References

- CIS Controls: 5.1 Minimize And Sparingly Use Administrative Privileges

## Profile

- Level 1 - Server
- Level 1 - Workstation
