# stage: post_exploit
# category: CIS_benchmarks


# 5.1.12 Minimize access to webhook configuration objects (Manual)

## Profile Applicability

- Level 1 - Master Node

## Description

Users with rights to create/modify/delete `validatingwebhookconfigurations` or `mutatingwebhookconfigurations` can control webhooks that can read any object admitted to the cluster, and in the case of mutating webhooks, also mutate admitted objects. This could allow for privilege escalation or disruption of the operation of the cluster.

## Rationale

The ability to manage webhook configuration should be limited.

## Audit

Review the users who have access to `validatingwebhookconfigurations` or `mutatingwebhookconfigurations` objects in the Kubernetes API.

## Remediation

Where possible, remove access to the `validatingwebhookconfigurations` or `mutatingwebhookconfigurations` objects.

## References

1. https://kubernetes.io/docs/concepts/security/rbac-good-practices/#control-admission-webhooks
