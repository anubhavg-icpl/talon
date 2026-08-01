# stage: report
# category: CIS_benchmarks


# 4.5.1 Configure Image Provenance using ImagePolicyWebhook admission controller (Manual)

## Profile Applicability

- Level 2

## Description

Configure Image Provenance for the deployment.

## Rationale

Kubernetes supports plugging in provenance rules to accept or reject the images in deployments. Rules can be configured to ensure that only approved images are deployed in the cluster.

## Impact

Regular maintenance for the provenance configuration should be carried out, based on container image updates.

## Audit

Review the pod definitions in the cluster and verify that image provenance is configured as appropriate.

## Remediation

Follow the Kubernetes documentation and setup image provenance.

## Default Value

By default, image provenance is not set.

## References

1. https://kubernetes.io/docs/concepts/containers/images/
2. https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/

## CIS Controls

| Controls Version | Control                                            | IG 1 | IG 2 | IG 3 |
| ---------------- | -------------------------------------------------- | ---- | ---- | ---- |
| v8               | 4.6 Securely Manage Enterprise Assets and Software | \*   | \*   | \*   |
| v7               | 18 Application Software Security                   |      |      |      |
