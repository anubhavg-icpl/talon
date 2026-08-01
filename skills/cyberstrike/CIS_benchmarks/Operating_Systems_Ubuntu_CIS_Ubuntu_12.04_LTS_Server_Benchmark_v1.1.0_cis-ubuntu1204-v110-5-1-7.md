# stage: report
# category: CIS_benchmarks


# 5.1.7 Ensure tftp-server is not enabled (Scored)

## Profile Applicability

- Level 1

## Description

Trivial File Transfer Protocol (TFTP) is a simple file transfer protocol, typically used to automatically transfer configuration or boot machines from a boot server. The packages `tftp` and `atftp` are both used to define and support a TFTP server.

## Rationale

TFTP does not support authentication nor does it ensure the confidentiality or integrity of data. It is recommended that TFTP be removed, unless there is a specific need for TFTP. In that case, extreme caution must be used when configuring the services.

## Audit Procedure

### Using Command Line

Ensure the `tftp` service is not enabled:

```bash
grep ^tftp /etc/inetd.conf
```

## Expected Result

No results should be returned.

## Remediation

### Using Command Line

Remove or comment out any `tftp` lines in `/etc/inetd.conf`:

```bash
sed -i 's/^tftp/#tftp/' /etc/inetd.conf
```

## Default Value

Not enabled by default on Ubuntu 12.04 LTS Server.

## References

- CIS Ubuntu 12.04 LTS Server Benchmark v1.1.0

## Profile

Level 1 - Scored
