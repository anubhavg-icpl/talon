# stage: report
# category: CIS_benchmarks


# 5.1.21 Ensure SSH MaxStartups is configured (Automated)

## Description

The `MaxStartups` parameter specifies the maximum number of concurrent unauthenticated connections to the SSH daemon.

## Rationale

To protect a system from denial of service due to a large number of pending authentication connection attempts, use the rate limiting function of MaxStartups to protect availability of sshd logins and prevent overwhelming the daemon.

## Audit Procedure

Run the following command and verify that output `MaxStartups` is `10:30:60` or matches site policy:

```bash
# sshd -T | grep -i maxstartups
# maxstartups 10:30:60
```

## Expected Result

Output should show `maxstartups 10:30:60` or match site policy.

## Remediation

Edit the `/etc/ssh/sshd_config` file to set the parameter as follows:

```
maxstartups 10:30:60
```

## CIS Controls

| Controls Version | Control                             | IG 1 | IG 2 | IG 3 |
| ---------------- | ----------------------------------- | ---- | ---- | ---- |
| v7               | 5.1 Establish Secure Configurations | x    | x    | x    |

## Profile

- Level 2 - Server
