# stage: report
# category: CIS_benchmarks


# 5.1.17 Ensure SSH LoginGraceTime is set to one minute or less (Automated)

## Description

The `LoginGraceTime` parameter specifies the time allowed for successful authentication to the SSH server. The longer the Grace period is the more open unauthenticated connections can exist. Like other session controls in this session the Grace Period should be limited to appropriate organizational limits to ensure the service is available for needed access.

## Rationale

Setting the `LoginGraceTime` parameter to a low number will minimize the risk of successful brute force attacks to the SSH server. It will also limit the number of concurrent unauthenticated connections. While the recommended setting is 60 seconds (1 Minute), set the number based on site policy.

## Audit Procedure

Run the following command and verify that output `LoginGraceTime` is between 1 and 60:

```bash
# sshd -T | grep logingracetime

LoginGraceTime 60
```

## Expected Result

Output should show `LoginGraceTime` between 1 and 60.

## Remediation

Edit the `/etc/ssh/sshd_config` file to set the parameter as follows:

```
LoginGraceTime 60
```

## Default Value

LoginGraceTime 120

## CIS Controls

| Controls Version | Control                             | IG 1 | IG 2 | IG 3 |
| ---------------- | ----------------------------------- | ---- | ---- | ---- |
| v7               | 5.1 Establish Secure Configurations | x    | x    | x    |

## Profile

- Level 2 - Server
