# stage: post_exploit
# category: CIS_benchmarks


# CIS Kubernetes Benchmark v1.12.0 - Control 2.3

## Profile Applicability

- **Level:** 1 - Master Node

## Description

Do not use self-signed certificates for TLS.

## Rationale

etcd is a highly-available key value store used by Kubernetes deployments for persistent storage of all of its REST API objects. These objects are sensitive in nature and should not be available to unauthenticated clients. You should enable the client authentication via valid certificates to secure the access to the etcd service.

## Impact

Clients will not be able to use self-signed certificates for TLS.

## Audit Procedure

Run the following command on the etcd server node:

```bash
ps -ef | grep etcd
```

Verify that if the `--auto-tls` argument exists, it is not set to `true`.

## Remediation

Edit the etcd pod specification file `/etc/kubernetes/manifests/etcd.yaml` on the master node and either remove the `--auto-tls` parameter or set it to `false`.

```
--auto-tls=false
```

## Default Value

By default, `--auto-tls` is set to `false`.

## References

1. https://coreos.com/etcd/docs/latest/op-guide/security.html
2. https://kubernetes.io/docs/admin/etcd/
3. https://coreos.com/etcd/docs/latest/op-guide/configuration.html#auto-tls

## CIS Controls

| Controls Version | Control                                           | IG 1 | IG 2 | IG 3 |
| ---------------- | ------------------------------------------------- | ---- | ---- | ---- |
| v8               | 3.11 Encrypt Sensitive Data at Rest               |      |      |      |
| v7               | 14.4 Encrypt All Sensitive Information in Transit |      |      |      |

## Profile

**Level 1 - Master Node** (Automated)
