# stage: post_exploit
# category: CIS_benchmarks


# 4.1.11 Minimize access to webhook configuration objects (Manual)

## Profile Applicability

- Level 1

## Description

Users with rights to create/modify/delete `validatingwebhookconfigurations` or `mutatingwebhookconfigurations` can control webhooks that can read any object admitted to the cluster, and in the case of mutating webhooks, also mutate admitted objects. This could allow for privilege escalation or disruption of the operation of the cluster.

## Rationale

The ability to manage webhook configuration should be limited.

## Audit

Review the users who have access to `validatingwebhookconfigurations` or `mutatingwebhookconfigurations` objects in the Kubernetes API.

## Remediation

Where possible, remove access to the `validatingwebhookconfigurations` or `mutatingwebhookconfigurations` objects.

## Default Value

N/A

## References

1. [https://kubernetes.io/docs/concepts/security/rbac-good-practices/#control-admission-webhooks](https://kubernetes.io/docs/concepts/security/rbac-good-practices/#control-admission-webhooks)

## CIS Controls

| Controls Version | Control                                           | IG 1 | IG 2 | IG 3 |
| ---------------- | ------------------------------------------------- | ---- | ---- | ---- |
| v8               | 6.8 Define and Maintain Role-Based Access Control |      |      | x    |
