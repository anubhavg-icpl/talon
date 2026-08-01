# stage: post_exploit
# category: CIS_benchmarks


# CIS Kubernetes Benchmark v1.12.0 - Control 4.1.4

## Description

If `kube-proxy` is running, ensure that the file ownership of its kubeconfig file is set to `root:root`.

## Rationale

The kubeconfig file for `kube-proxy` controls various parameters for the `kube-proxy` service in the worker node. You should set its file ownership to maintain the integrity of the file. The file should be owned by `root:root`.

## Impact

None

## Audit Procedure

### Step 1: Find the kubeconfig file being used by kube-proxy

```bash
ps -ef | grep kube-proxy
```

If `kube-proxy` is running, get the kubeconfig file location from the `--kubeconfig` parameter.

### Step 2: Verify file ownership

Run the below command (based on the file location on your system) on the each worker node. For example,

```bash
stat -c %U:%G <path><filename>
```

Verify that the ownership is set to `root:root`.

## Remediation

Run the below command (based on the file location on your system) on the each worker node. For example,

```bash
chown root:root <proxy kubeconfig file>
```

## Default Value

By default, `proxy` file ownership is set to `root:root`.

## References

1. https://kubernetes.io/docs/admin/kube-proxy/

## CIS Controls

| Controls Version | Control                                                                                                                                                                  | IG 1 | IG 2 | IG 3 |
| ---------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | ---- | ---- | ---- |
| v8               | 5.4 Restrict Administrator Privileges to Dedicated Administrator Accounts<br>Restrict administrator privileges to dedicated administrator accounts on enterprise assets. |      | ●    | ●    |
| v7               | 4 Controlled Use of Administrative Privileges<br>Controlled Use of Administrative Privileges                                                                             |      |      |      |

## Profile Applicability

- Level 1 - Worker Node

## Assessment Status

Manual
