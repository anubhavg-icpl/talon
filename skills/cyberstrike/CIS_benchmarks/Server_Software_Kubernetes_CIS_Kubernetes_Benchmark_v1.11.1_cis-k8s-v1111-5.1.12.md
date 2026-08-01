# stage: post_exploit
# category: CIS_benchmarks


# 5.1.12 Minimize access to webhook configuration objects (Manual)

## Profile Applicability

- Level 1 - Master Node

## Description

Users with rights to create/modify/delete `validatingwebhookconfigurations` or `mutatingwebhookconfigurations` can control webhooks that can read any object admitted to the cluster, and in the case of mutating webhooks, also mutate admitted objects. This could allow for privilege escalation or disruption of the operation of the cluster.

## Rationale

The ability to manage webhook configuration should be limited.

## Impact

None

## Audit

Review the users who have access to `validatingwebhookconfigurations` or `mutatingwebhookconfigurations` objects in the Kubernetes API.

## Remediation

Where possible, remove access to the `validatingwebhookconfigurations` or `mutatingwebhookconfigurations` objects.

## Default Value

None

## References

1. https://kubernetes.io/docs/concepts/security/rbac-good-practices/#control-admission-webhooks

## CIS Controls

| Controls Version | Control                                           | IG 1 | IG 2 | IG 3 |
| ---------------- | ------------------------------------------------- | ---- | ---- | ---- |
| v8               | 6.8 Define and Maintain Role-Based Access Control |      |      | X    |
| v7               | 14 Controlled Access Based on the Need to Know    |      |      |      |
