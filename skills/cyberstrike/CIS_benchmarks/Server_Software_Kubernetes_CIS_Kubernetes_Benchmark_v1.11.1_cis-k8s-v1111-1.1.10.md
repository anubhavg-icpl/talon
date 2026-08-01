# stage: post_exploit
# category: CIS_benchmarks


# CIS Kubernetes Benchmark v1.11.1 - Control 1.1.10

## Profile Applicability

- **Level:** 1 - Master Node

## Description

Ensure that the Container Network Interface files have ownership set to `root:root`.

## Rationale

Container Network Interface provides various networking options for overlay networking. You should consult their documentation and restrict their respective file permissions to maintain the integrity of those files. Those files should be owned by `root:root`.

## Impact

None

## Audit

Run the below command (based on the file location on your system) on the Control Plane node. For example,

```bash
stat -c %U:%G <path/to/cni/files>
```

Verify that the ownership is set to `root:root`.

## Remediation

Run the below command (based on the file location on your system) on the Control Plane node. For example,

```bash
chown root:root <path/to/cni/files>
```

## Default Value

NA

## References

1. https://kubernetes.io/docs/concepts/cluster-administration/networking/

## CIS Controls

| Controls Version | Control                                                                                                                                                                                                                                                                                                                        | IG 1 | IG 2 | IG 3 |
| ---------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | ---- | ---- | ---- |
| v8               | 5.4 Restrict Administrator Privileges to Dedicated Administrator Accounts<br/>Restrict administrator privileges to dedicated administrator accounts on enterprise assets. Conduct general computing activities, such as internet browsing, email, and productivity suite use, from the user's primary, non-privileged account. | ●    | ●    | ●    |
| v7               | 4 Controlled Use of Administrative Privileges<br/>Controlled Use of Administrative Privileges                                                                                                                                                                                                                                  |      |      |      |
