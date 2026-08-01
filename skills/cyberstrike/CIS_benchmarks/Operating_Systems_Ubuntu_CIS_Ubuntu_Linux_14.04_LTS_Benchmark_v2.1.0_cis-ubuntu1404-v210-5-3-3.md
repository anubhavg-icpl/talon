# stage: report
# category: CIS_benchmarks


# 5.3.3 Ensure password reuse is limited (Scored)

## Profile Applicability

- Level 1 - Server
- Level 1 - Workstation

## Description

The `/etc/security/opasswd` file stores the users' old passwords and can be checked to ensure that users are not recycling recent passwords.

## Rationale

Forcing users not to reuse their past 5 passwords make it less likely that an attacker will be able to guess the password.

Note that these change only apply to accounts configured on the local system.

## Audit Procedure

Run the following commands and ensure the `remember` option is '5' or more and included in all results:

```bash
egrep '^password\s+required\s+pam_pwhistory.so' /etc/pam.d/common-password
```

## Expected Result

```
password required pam_pwhistory.so remember=5
```

## Remediation

Edit the `/etc/pam.d/common-password` file to include the `remember` option and conform to site policy as shown:

```
password required pam_pwhistory.so remember=5
```

## Default Value

No password reuse limit configured by default.

## References

- CIS Controls: 16 - Account Monitoring and Control

## Profile

- Level 1 - Server
- Level 1 - Workstation
