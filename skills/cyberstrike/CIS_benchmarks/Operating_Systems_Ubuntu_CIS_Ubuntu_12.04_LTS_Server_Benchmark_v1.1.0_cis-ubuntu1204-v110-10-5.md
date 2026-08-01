# stage: report
# category: CIS_benchmarks


# 10.5 Lock Inactive User Accounts (Scored)

## Profile Applicability

- Level 1

## Description

User accounts that have been inactive for over a given period of time can be automatically disabled. It is recommended that accounts that are inactive for 35 or more days be disabled.

## Rationale

Inactive accounts pose a threat to system security since the users are not logging in to notice failed login attempts or other anomalies.

## Audit Procedure

### Using Command Line

```bash
useradd -D | grep INACTIVE
```

## Expected Result

```
INACTIVE=35
```

## Remediation

### Using Command Line

```bash
useradd -D -f 35
```

## Default Value

INACTIVE=-1 (disabled)

## References

- CIS Ubuntu 12.04 LTS Server Benchmark v1.1.0

## Profile

Level 1 - Scored
