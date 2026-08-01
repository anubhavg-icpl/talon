# stage: report
# category: CIS_benchmarks


# 4.2.3 Minimize the admission of containers wishing to share the host IPC namespace (Automated)

## Profile Applicability

- Level 1

## Description

Do not generally permit containers to be run with the `hostIPC` flag set to true.

## Rationale

A container running in the host's IPC namespace can use IPC to interact with processes outside the container.

There should be at least one admission control policy defined which does not permit containers to share the host IPC namespace.

If you need to run containers which require hostIPC, this should be defined in a separate policy and you should carefully check to ensure that only limited service accounts and users are given permission to use that policy.

## Impact

Pods defined with `spec.hostIPC: true` will not be permitted unless they are run under a specific policy.

## Audit

List the policies in use for each namespace in the cluster, ensure that each policy disallows the admission of `hostIPC` containers.

Search for the hostIPC Flag: In the YAML output, look for the `hostIPC` setting under the spec section to check if it is set to `true`.

```bash
kubectl get pods -A -o json \
| jq -r '
  .items[]
  | select((.spec.hostIPC // false) == true)
  | "\(.metadata.namespace)/\(.metadata.name)\thostIPC=true"
'
```

OR to include pod UID and node name:

```bash
kubectl get pods -A -o json \
| jq -r '
  .items[]
  | select((.spec.hostIPC // false) == true)
  | "\(.metadata.namespace)/\(.metadata.name)\tnode=\(.spec.nodeName // "-")\thostIPC=true"
'
```

When creating a Pod Security Policy, ["kube-system"] namespaces are excluded by default.

## Remediation

Add policies to each namespace in the cluster which has user workloads to restrict the admission of `hostIPC` containers.

## Default Value

By default, there are no restrictions on the creation of `hostIPC` containers.

## References

1. https://kubernetes.io/docs/concepts/security/pod-security-admission/
2. https://docs.cloud.oracle.com/en-us/iaas/Content/ContEng/Concepts/contengoverview.htm

## CIS Controls

| Controls Version | Control                                                       | IG 1 | IG 2 | IG 3 |
| ---------------- | ------------------------------------------------------------- | ---- | ---- | ---- |
| v8               | 3.12 Segment Data Processing and Storage Based on Sensitivity |      | x    | x    |
| v7               | 14.1 Segment the Network Based on Sensitivity                 |      | x    | x    |

## MITRE ATT&CK Mappings

| Techniques / Sub-techniques | Tactics | Mitigations |
| --------------------------- | ------- | ----------- |
| T1098                       | TA0003  | M1030       |
