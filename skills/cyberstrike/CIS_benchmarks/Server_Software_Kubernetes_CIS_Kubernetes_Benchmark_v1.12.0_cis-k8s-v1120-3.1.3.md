# stage: post_exploit
# category: CIS_benchmarks


# CIS Kubernetes Benchmark v1.12.0 - Control 3.1.3

## Profile Applicability

- **Level:** 1 - Master Node

## Description

Kubernetes provides bootstrap tokens which are intended for use by new nodes joining the cluster.

These tokens are not designed for use by end-users they are specifically designed for the purpose of bootstrapping new nodes and not for general authentication.

## Rationale

Bootstrap tokens are not intended for use as a general authentication mechanism and impose constraints on user and group naming that do not facilitate good RBAC design. They also cannot be used with MFA resulting in a weak authentication mechanism being available.

## Impact

External mechanisms for authentication generally require additional software to be deployed.

## Audit Procedure

Review user access to the cluster and ensure that users are not making use of bootstrap token authentication.

## Remediation

Alternative mechanisms provided by Kubernetes such as the use of OIDC should be implemented in place of bootstrap tokens.

## Default Value

Bootstrap token authentication is not enabled by default and requires an API server parameter to be set.

## References

None

## CIS Controls

| Controls Version | Control                                    | IG 1 | IG 2 | IG 3 |
| ---------------- | ------------------------------------------ | ---- | ---- | ---- |
| v8               | 6.2 Establish an Access Revoking Process   |      |      |      |
| v7               | 16.7 Establish Process for Revoking Access |      |      |      |

## Profile

**Level 1 - Master Node** (Manual)
