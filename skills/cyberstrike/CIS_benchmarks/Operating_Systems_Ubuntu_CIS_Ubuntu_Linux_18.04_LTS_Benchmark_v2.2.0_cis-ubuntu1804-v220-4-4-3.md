# stage: report
# category: CIS_benchmarks


# CIS Ubuntu Linux 18.04 LTS Benchmark v2.2.0 - Control 4.4.3

## Profile Applicability

- Level 1 - Server
- Level 1 - Workstation

## Description

The `/etc/security/opasswd` file stores the users' old passwords and can be checked to ensure that users are not recycling recent passwords.

- `remember = 24` - Number of old passwords to remember. The user is not able to reuse any of the most recent `remember` passwords.

## Rationale

Forcing users not to reuse their past 24 passwords make it less likely that an attacker will be able to guess the password. Note that these change only apply to accounts configured on the local system.

## Audit Procedure

### Command Line

Run the following command to verify password reuse is limited:

```bash
grep -Pi '^\h*password\h+(requisite|required|sufficient)\h+pam_unix\.so\h+([^#\n\r]+\h+)?remember=([5-9]|[1-9][0-9]+)\b' /etc/pam.d/common-password
```

### Expected Result

The `remember` parameter should be set to 24 or more.

## Remediation

### Command Line

Edit `/etc/pam.d/common-password` and add or modify the `pam_unix.so` or `pam_pwhistory.so` line to include `remember=24`:

```
password required pam_pwhistory.so remember=24
```

OR ensure `pam_unix.so` includes `remember=24`:

```
password [success=1 default=ignore] pam_unix.so obscure use_authtok try_first_pass sha512 remember=24
```

## References

1. NIST SP 800-53 Rev. 5: IA-5(1)

## CIS Controls

v8 - 5.2 Use Unique Passwords - Use unique passwords for all enterprise assets.

v7 - 4.4 Use Unique Passwords.

## Profile Applicability

- Level 1 - Server
- Level 1 - Workstation

## Assessment Status

Automated
