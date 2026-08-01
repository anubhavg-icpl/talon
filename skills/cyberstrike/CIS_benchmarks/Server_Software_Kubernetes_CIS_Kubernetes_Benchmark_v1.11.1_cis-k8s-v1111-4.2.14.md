# stage: post_exploit
# category: CIS_benchmarks


# 4.2.14 Ensure that the --seccomp-default parameter is set to true (Manual)

## Profile Applicability

- Level 1 - Worker Node

## Description

Ensure that the Kubelet enforces the use of the RuntimeDefault seccomp profile

## Rationale

By default, Kubernetes disables the seccomp profile which ships with most container runtimes. Setting this parameter will ensure workloads running on the node are protected by the runtime's seccomp profile.

## Impact

Setting this will remove some rights from pods running on the node.

## Audit

Review the Kubelet's start-up parameters for the value of `--seccomp-default`, and check the Kubelet configuration file for the `seccompDefault`. If neither of these values is set, then the seccomp profile is not in use.

## Remediation

Set the parameter, either via the `--seccomp-default` command line parameter or the `seccompDefault` configuration file setting.

## Default Value

By default the seccomp profile is not enabled.

## References

1. [https://kubernetes.io/docs/tutorials/security/seccomp/#enable-the-use-of-runtimedefault-as-the-default-seccomp-profile-for-all-workloads](https://kubernetes.io/docs/tutorials/security/seccomp/#enable-the-use-of-runtimedefault-as-the-default-seccomp-profile-for-all-workloads)

## CIS Controls

| Controls Version | Control                                                                                       | IG 1 | IG 2 | IG 3 |
| ---------------- | --------------------------------------------------------------------------------------------- | ---- | ---- | ---- |
| v8               | 13.7 Deploy a Host-Based Intrusion Prevention Solution                                        |      |      | x    |
| v7               | 5.1 Establish Secure Configurations                                                           | x    | x    | x    |
| v7               | 5.2 Maintain Secure Images                                                                    |      | x    | x    |
| v7               | 11.4 Install the Latest Stable Version of Any Security-related Updates on All Network Devices | x    | x    | x    |
