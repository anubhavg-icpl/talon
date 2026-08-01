# stage: post_exploit
# category: CIS_benchmarks


# 1.1.17 Ensure that the controller-manager.conf file permissions are set to 600 or more restrictive (Automated)

## Profile Applicability

- Level 1 - Master Node

## Description

Ensure that the `controller-manager.conf` file has permissions of 600 or more restrictive.

## Rationale

The `controller-manager.conf` file is the kubeconfig file for the Controller Manager. You should restrict its file permissions to maintain the integrity of the file. The file should be writable by only the administrators on the system.

## Impact

None

## Audit

Run the following command (based on the file location on your system) on the Control Plane node. For example,

```bash
stat -c %a /etc/kubernetes/controller-manager.conf
```

Verify that the permissions are `600` or more restrictive.

## Remediation

Run the below command (based on the file location on your system) on the Control Plane node. For example,

```bash
chmod 600 /etc/kubernetes/controller-manager.conf
```

## Default Value

By default, `controller-manager.conf` has permissions of `640`.

## References

1. [https://kubernetes.io/docs/admin/kube-controller-manager/](https://kubernetes.io/docs/admin/kube-controller-manager/)

## CIS Controls

| Controls Version | Control                                               | IG 1 | IG 2 | IG 3 |
| ---------------- | ----------------------------------------------------- | ---- | ---- | ---- |
| v8               | 3.3 Configure Data Access Control Lists               | \*   | \*   | \*   |
| v7               | 14.6 Protect Information through Access Control Lists | \*   | \*   | \*   |
