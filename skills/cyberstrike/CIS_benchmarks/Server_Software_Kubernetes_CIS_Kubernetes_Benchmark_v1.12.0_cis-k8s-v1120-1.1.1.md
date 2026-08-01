# stage: post_exploit
# category: CIS_benchmarks


# 1.1.1 Ensure that the API server pod specification file permissions are set to 600 or more restrictive (Automated)

## Profile Applicability

- Level 1 - Master Node

## Description

Ensure that the API server pod specification file has permissions of `600` or more restrictive.

## Rationale

The API server pod specification file controls various parameters that set the behavior of the API server. You should restrict its file permissions to maintain the integrity of the file. The file should be writable by only the administrators on the system.

## Impact

None

## Audit

Run the below command (based on the file location on your system) on the Control Plane node. For example,

```bash
stat -c %a /etc/kubernetes/manifests/kube-apiserver.yaml
```

Verify that the permissions are `600` or more restrictive.

## Remediation

Run the below command (based on the file location on your system) on the Control Plane node. For example,

```bash
chmod 600 /etc/kubernetes/manifests/kube-apiserver.yaml
```

## Default Value

By default, the `kube-apiserver.yaml` file has permissions of `640`.

## References

1. [https://kubernetes.io/docs/admin/kube-apiserver/](https://kubernetes.io/docs/admin/kube-apiserver/)

## CIS Controls

| Controls Version | Control                                               | IG 1 | IG 2 | IG 3 |
| ---------------- | ----------------------------------------------------- | ---- | ---- | ---- |
| v8               | 3.3 Configure Data Access Control Lists               | \*   | \*   | \*   |
| v7               | 14.6 Protect Information through Access Control Lists | \*   | \*   | \*   |
