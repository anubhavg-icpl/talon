# stage: post_exploit
# category: CIS_benchmarks


# 5.2.9 Minimize the admission of containers with added capabilities (Manual)

## Profile Applicability

- Level 1 - Master Node

## Description

Do not generally permit containers with capabilities assigned beyond the default set.

## Rationale

Containers run with a default set of capabilities as assigned by the Container Runtime. Capabilities outside this set can be added to containers which could expose them to risks of container breakout attacks.

There should be at least one policy defined which prevents containers with capabilities beyond the default set from launching.

If you need to run containers with additional capabilities, this should be defined in a separate policy and you should carefully check to ensure that only limited service accounts and users are given permission to use that policy.

## Impact

Pods with containers which require capabilities outwith the default set will not be permitted.

## Audit

Ensure that allowedCapabilities is not present in policies for the cluster unless it is set to an empty array.

```bash
kubectl get pods -A -o=jsonpath=$'{range .items[*]}{@.metadata.name}: {@..securityContext}\n{end}'
```

## Remediation

Ensure that `allowedCapabilities` is not present in policies for the cluster unless it is set to an empty array.

## Default Value

By default, there are no restrictions on adding capabilities to containers.

## References

1. https://kubernetes.io/docs/concepts/security/pod-security-standards/
2. https://www.nccgroup.trust/uk/our-research/abusing-privileged-and-unprivileged-linux-containers/

## CIS Controls

| Controls Version | Control                                                                         | IG 1 | IG 2 | IG 3 |
| ---------------- | ------------------------------------------------------------------------------- | ---- | ---- | ---- |
| v8               | 4.8 Uninstall or Disable Unnecessary Services on Enterprise Assets and Software |      | X    | X    |
| v7               | 5.2 Maintain Secure Images                                                      |      | X    | X    |
