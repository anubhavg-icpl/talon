# stage: post_exploit
# category: CIS_benchmarks


# 1.1.21 Ensure that the Kubernetes PKI key file permissions are set to 600 (Manual)

## Profile Applicability

- Level 1 - Master Node

## Description

Ensure that Kubernetes PKI key files have permissions of `600`.

## Rationale

Kubernetes makes use of a number of key files as part of the operation of its components. The permissions on these files should be set to `600` to protect their integrity and confidentiality.

## Impact

None

## Audit

Run the below command (based on the file location on your system) on the Control Plane node. For example,

```bash
stat -c '%a' /etc/kubernetes/pki/*.key
```

Verify that the permissions are `600` or more restrictive.

or

```bash
ls -l /etc/kubernetes/pki/*.key
```

Verify that the permissions are `-rw------`

## Remediation

Run the below command (based on the file location on your system) on the Control Plane node. For example,

```bash
chmod -R 600 /etc/kubernetes/pki/*.key
```

## Default Value

By default, the keys used by Kubernetes are set to have permissions of `600`

## References

1. [https://kubernetes.io/docs/admin/kube-apiserver/](https://kubernetes.io/docs/admin/kube-apiserver/)

## CIS Controls

| Controls Version | Control                                                                   | IG 1 | IG 2 | IG 3 |
| ---------------- | ------------------------------------------------------------------------- | ---- | ---- | ---- |
| v8               | 5.4 Restrict Administrator Privileges to Dedicated Administrator Accounts | \*   | \*   | \*   |
| v7               | 4 Controlled Use of Administrative Privileges                             |      |      |      |
