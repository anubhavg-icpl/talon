# stage: report
# category: CIS_benchmarks


# 6.3.3.2 Ensure actions as another user are always logged (Automated)

## Profile

- Level 2 - Server
- Level 2 - Workstation

## Description

`sudo` provides users with temporary elevated privileges to perform operations, either as the superuser or another user.

## Rationale

Creating an audit log of users with temporary elevated privileges and the operation(s) they performed is essential to reporting. Administrators will want to correlate the events written to the audit trail with the records written to `sudo`'s logfile to verify if unauthorized commands have been executed.

## Audit Procedure

### Command Line

On disk configuration - Run the following command to check the on disk rules:

```bash
# awk '/^ *-a *always,exit/ \
&&/ -F *arch=b(32|64)/ \
&&(/ -F *auid!=unset/||/ -F *auid!=-1/||/ -F *auid!=4294967295/) \
&&(/ -C *euid!=uid/||/ -C *uid!=euid/) \
&&/ -S *execve/ \
&&(/ key= *[!-~]* *$/||/ -k *[!-~]* *$/)' /etc/audit/rules.d/*.rules
```

Verify the output matches:

```
-a always,exit -F arch=b64 -C euid!=uid -F auid!=unset -S execve -k user_emulation
-a always,exit -F arch=b32 -C euid!=uid -F auid!=unset -S execve -k user_emulation
```

Running configuration - Run the following command to check loaded rules:

```bash
# auditctl -l | awk '/^ *-a *always,exit/ \
&&/ -F *arch=b(32|64)/ \
&&(/ -F *auid!=unset/||/ -F *auid!=-1/||/ -F *auid!=4294967295/) \
&&(/ -C *euid!=uid/||/ -C *uid!=euid/) \
&&/ -S *execve/ \
&&(/ key= *[!-~]* *$/||/ -k *[!-~]* *$/)'
```

Verify the output matches:

```
-a always,exit -F arch=b64 -S execve -C uid!=euid -F auid!=-1 -F key=user_emulation
-a always,exit -F arch=b32 -S execve -C uid!=euid -F auid!=-1 -F key=user_emulation
```

## Expected Result

Both on disk and running configuration should show execve rules tracking user emulation.

## Remediation

### Command Line

Create audit rules - Edit or create a file in the `/etc/audit/rules.d/` directory, ending in `.rules` extension, with the relevant rules to monitor elevated privileges.

Example:

```bash
# printf "
-a always,exit -F arch=b64 -C euid!=uid -F auid!=unset -S execve -k user_emulation
-a always,exit -F arch=b32 -C euid!=uid -F auid!=unset -S execve -k user_emulation
" >> /etc/audit/rules.d/50-user_emulation.rules
```

Load audit rules - Merge and load the rules into active configuration:

```bash
# augenrules --load
```

Check if reboot is required:

```bash
# if [[ $(auditctl -s | grep "enabled") =~ "2" ]]; then printf "Reboot required to load rules\n"; fi
```

## References

1. NIST SP 800-53 Rev. 5: AU-3

## CIS Controls

| Controls Version | Control                                                        | IG 1 | IG 2 | IG 3 |
| ---------------- | -------------------------------------------------------------- | ---- | ---- | ---- |
| v8               | 8.5 Collect Detailed Audit Logs                                |      | X    | X    |
| v7               | 4.9 Log and Alert on Unsuccessful Administrative Account Login |      | X    | X    |

MITRE ATT&CK Mappings: T1562, T1562.006 / TA0004 / M1047
