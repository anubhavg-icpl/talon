# stage: post_exploit
# category: CIS_benchmarks


# 1.2.6 Ensure that the --authorization-mode argument is not set to AlwaysAllow (Automated)

## Profile Applicability

- Level 1 - Master Node

## Description

Do not always authorize all requests.

## Rationale

The API Server, can be configured to allow all requests. This mode should not be used on any production cluster.

## Impact

Only authorized requests will be served.

## Audit

Run the following command on the Control Plane node:

```bash
ps -ef | grep kube-apiserver
```

Verify that the `--authorization-mode` argument exists and is not set to `AlwaysAllow`.

## Remediation

Edit the API server pod specification file `/etc/kubernetes/manifests/kube-apiserver.yaml` on the Control Plane node and set the `--authorization-mode` parameter to values other than `AlwaysAllow`. One such example could be as below.

```
--authorization-mode=RBAC
```

## Default Value

By default, `AlwaysAllow` is not enabled.

## References

1. [https://kubernetes.io/docs/reference/command-line-tools-reference/kube-apiserver/](https://kubernetes.io/docs/reference/command-line-tools-reference/kube-apiserver/)
2. [https://kubernetes.io/docs/admin/authorization/](https://kubernetes.io/docs/admin/authorization/)

## CIS Controls

| Controls Version | Control                                                            | IG 1 | IG 2 | IG 3 |
| ---------------- | ------------------------------------------------------------------ | ---- | ---- | ---- |
| v8               | 3.3 Configure Data Access Control Lists                            | \*   | \*   | \*   |
| v7               | 9.2 Ensure Only Approved Ports, Protocols and Services Are Running |      | \*   | \*   |
