# stage: report
# category: CIS_benchmarks


# 1.6.2.2 Ensure all AppArmor Profiles are enforcing (Scored)

## Profile Applicability

- Level 2 - Server
- Level 2 - Workstation

## Description

AppArmor profiles define what resources applications are able to access.

## Rationale

Security configuration requirements vary from site to site. Some sites may mandate a policy that is stricter than the default policy, which is perfectly acceptable. This item is intended to ensure that any policies that exist on the system are activated.

## Audit Procedure

Run the following command and verify that profiles are loaded, no profiles are in complain mode, and no processes are unconfined:

```bash
apparmor_status
```

## Expected Result

```
apparmor module is loaded.
X profiles are loaded.
X profiles are in enforce mode.
0 profiles are in complain mode.
X processes have profiles defined.
X processes are in enforce mode.
0 processes are in complain mode.
0 processes are unconfined but have a profile defined.
```

## Remediation

Run the following command to set all profiles to enforce mode:

```bash
aa-enforce /etc/apparmor.d/*
```

Any unconfined processes may need to have a profile created or activated for them and then be restarted.

## Default Value

AppArmor profiles are loaded in enforce mode by default on Ubuntu.

## References

- CIS Controls: 14.4 Protect Information With Access Control Lists

## Profile

- Level 2 - Server
- Level 2 - Workstation
