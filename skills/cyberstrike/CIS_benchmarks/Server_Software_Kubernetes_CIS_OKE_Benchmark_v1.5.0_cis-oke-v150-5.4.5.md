# stage: report
# category: CIS_benchmarks


# 5.4.5 Encrypt traffic to HTTPS load balancers with TLS certificates (Manual)

## Profile Applicability

- Level 1

## Description

Encrypt traffic to HTTPS load balancers using TLS certificates.

## Rationale

Encrypting traffic between users and your Kubernetes workload is fundamental to protecting data sent over the web.

## Impact

None specified.

## Audit

Your load balancer vendor can provide details on auditing the certificates and policies required to utilize TLS.

## Remediation

Your load balancer vendor can provide details on configuring HTTPS with TLS.

## Default Value

Not configured by default.

## References

1. [https://docs.cloud.oracle.com/en-us/iaas/Content/ContEng/Tasks/contengcreatingloadbalancer.htm](https://docs.cloud.oracle.com/en-us/iaas/Content/ContEng/Tasks/contengcreatingloadbalancer.htm)

## CIS Controls

| Controls Version | Control                                                                                                                                                                          | IG 1 | IG 2 | IG 3 |
| ---------------- | -------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---- | ---- | ---- |
| v8               | 3.10 Encrypt Sensitive Data in Transit - Encrypt sensitive data in transit. Example implementations can include: Transport Layer Security (TLS) and Open Secure Shell (OpenSSH). |      | x    | x    |
| v7               | 14.4 Encrypt All Sensitive Information in Transit - Encrypt all sensitive information in transit.                                                                                |      | x    | x    |

## MITRE ATT&CK Mappings

| Techniques / Sub-techniques | Tactics | Mitigations |
| --------------------------- | ------- | ----------- |
| T1609                       | TA0002  | M1035       |

## Profile

Level 1 - OKE
