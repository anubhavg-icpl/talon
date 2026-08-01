# stage: post_exploit
# category: CIS_benchmarks


# 1.1.2 Ensure that the API server pod specification file ownership is set to root:root (Automated)

## Profile Applicability

- Level 1 - Master Node

## Description

Ensure that the API server pod specification file ownership is set to `root:root`.

## Rationale

The API server pod specification file controls various parameters that set the behavior of the API server. You should set its file ownership to maintain the integrity of the file. The file should be owned by `root:root`.

## Impact

None

## Audit

Run the below command (based on the file location on your system) on the Control Plane node. For example,

```bash
stat -c %U:%G /etc/kubernetes/manifests/kube-apiserver.yaml
```

Verify that the ownership is set to `root:root`.

## Remediation

Run the below command (based on the file location on your system) on the Control Plane node. For example,

```bash
chown root:root /etc/kubernetes/manifests/kube-apiserver.yaml
```

## Default Value

By default, the `kube-apiserver.yaml` file ownership is set to `root:root`.

## References

1. [https://kubernetes.io/docs/admin/kube-apiserver/](https://kubernetes.io/docs/admin/kube-apiserver/)

## CIS Controls

| Controls Version | Control                                                                                                                                                                                                                                                                                                                          | IG 1 | IG 2 | IG 3 |
| :--------------: | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | :--: | :--: | :--: |
|        v8        | **5.4 Restrict Administrator Privileges to Dedicated Administrator Accounts** - Restrict administrator privileges to dedicated administrator accounts on enterprise assets. Conduct general computing activities, such as internet browsing, email, and productivity suite use, from the user's primary, non-privileged account. |  \*  |  \*  |  \*  |
|        v7        | **4 Controlled Use of Administrative Privileges** - Controlled Use of Administrative Privileges                                                                                                                                                                                                                                  |      |      |      |
