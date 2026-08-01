# stage: post_exploit
# category: CIS_benchmarks


# 4.1.12 Minimize access to the service account token creation (Manual)

## Profile Applicability

- Level 1

## Description

Users with rights to create new service account tokens at a cluster level, can create long-lived privileged credentials in the cluster. This could allow for privilege escalation and persistent access to the cluster, even if the users account has been revoked.

## Rationale

The ability to create service account tokens should be limited.

## Audit

Review the users who have access to create the `token` sub-resource of `serviceaccount` objects in the Kubernetes API.

## Remediation

Where possible, remove access to the `token` sub-resource of `serviceaccount` objects.

## References

1. [https://kubernetes.io/docs/concepts/security/rbac-good-practices/#token-request](https://kubernetes.io/docs/concepts/security/rbac-good-practices/#token-request)

## CIS Controls

| Controls Version | Control                                           | IG 1 | IG 2 | IG 3 |
| ---------------- | ------------------------------------------------- | ---- | ---- | ---- |
| v8               | 6.8 Define and Maintain Role-Based Access Control |      |      | x    |
