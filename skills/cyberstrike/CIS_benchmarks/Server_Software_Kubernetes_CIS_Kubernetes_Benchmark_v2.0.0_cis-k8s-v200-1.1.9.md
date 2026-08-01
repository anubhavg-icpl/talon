# stage: post_exploit
# category: CIS_benchmarks


# 1.1.9 Ensure that the Container Network Interface file permissions are set to 600 or more restrictive (Manual)

## Profile Applicability

- Level 1 - Master Node

## Description

Ensure that the Container Network Interface files have permissions of `600` or more restrictive.

## Rationale

Container Network Interface provides various networking options for overlay networking. You should consult their documentation and restrict their respective file permissions to maintain the integrity of those files. Those files should be writable by only the administrators on the system.

## Impact

None

## Audit

Run the below command (based on the file location on your system) on the Control Plane node. For example,

```bash
stat -c %a <path/to/cni/files>
```

Verify that the permissions are `600` or more restrictive.

## Remediation

Run the below command (based on the file location on your system) on the Control Plane node. For example,

```bash
chmod 600 <path/to/cni/files>
```

## Default Value

NA

## References

1. https://kubernetes.io/docs/concepts/cluster-administration/networking/

## CIS Controls

| Controls Version | Control                                                                                                                                                                                                                                                                                                                                                                                                | IG 1 | IG 2 | IG 3 |
| ---------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ | ---- | ---- | ---- |
| v8               | 3.3 Configure Data Access Control Lists - Configure data access control lists based on a user's need to know. Apply data access control lists, also known as access permissions, to local and remote file systems, databases, and applications.                                                                                                                                                        | \*   | \*   | \*   |
| v7               | 14.6 Protect Information through Access Control Lists - Protect all information stored on systems with file system, network share, claims, application, or database specific access control lists. These controls will enforce the principle that only authorized individuals should have access to the information based on their need to access the information as a part of their responsibilities. | \*   | \*   | \*   |
