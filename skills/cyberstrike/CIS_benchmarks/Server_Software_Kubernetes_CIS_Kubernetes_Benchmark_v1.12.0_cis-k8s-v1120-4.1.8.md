# stage: post_exploit
# category: CIS_benchmarks


# CIS Kubernetes Benchmark v1.12.0 - Control 4.1.8

## Description

Ensure that the certificate authorities file ownership is set to `root:root`.

## Rationale

The certificate authorities file controls the authorities used to validate API requests. You should set its file ownership to maintain the integrity of the file. The file should be owned by `root:root`.

## Impact

None

## Audit Procedure

### Step 1: Find the client-ca-file

Run the following command:

```bash
ps -ef | grep kubelet
```

Find the file specified by the `--client-ca-file` argument.

### Step 2: Verify file ownership

Run the following command:

```bash
stat -c %U:%G <filename>
```

Verify that the ownership is set to `root:root`.

## Remediation

Run the following command to modify the ownership of the `--client-ca-file`:

```bash
chown root:root <filename>
```

## Default Value

By default no `--client-ca-file` is specified.

## References

1. https://kubernetes.io/docs/admin/authentication/#x509-client-certs

## CIS Controls

| Controls Version | Control                                                                                                                                                                  | IG 1 | IG 2 | IG 3 |
| ---------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | ---- | ---- | ---- |
| v8               | 5.4 Restrict Administrator Privileges to Dedicated Administrator Accounts<br>Restrict administrator privileges to dedicated administrator accounts on enterprise assets. |      | ●    | ●    |
| v7               | 4 Controlled Use of Administrative Privileges<br>Controlled Use of Administrative Privileges                                                                             |      |      |      |

## Profile Applicability

- Level 1 - Worker Node

## Assessment Status

Manual
