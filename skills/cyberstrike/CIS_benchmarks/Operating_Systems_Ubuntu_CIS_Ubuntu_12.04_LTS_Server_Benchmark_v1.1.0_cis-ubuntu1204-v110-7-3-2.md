# stage: report
# category: CIS_benchmarks


# 7.3.2 Disable IPv6 Redirect Acceptance (Not Scored)

## Profile Applicability

- Level 1

## Description

This setting prevents the system from accepting ICMP redirects. ICMP redirects tell the system about alternate routes for sending traffic.

## Rationale

It is recommended that systems not accept ICMP redirects as they could be tricked into routing traffic to compromised machines. Setting hard routes within the system (usually a single default route to a trusted router) protects the system from bad routes.

## Audit Procedure

### Using Command Line

Perform the following to determine if IPv6 redirects are disabled.

```bash
/sbin/sysctl net.ipv6.conf.all.accept_redirects
/sbin/sysctl net.ipv6.conf.default.accept_redirects
```

## Expected Result

```
net.ipv6.conf.all.accept_redirect = 0
net.ipv6.conf.default.accept_redirect = 0
```

## Remediation

### Using Command Line

Set the `net.ipv6.conf.all.accept_redirects` and `net.ipv6.conf.default.accept_redirects` parameters to 0 in `/etc/sysctl.conf`:

```bash
net.ipv6.conf.all.accept_redirects=0
net.ipv6.conf.default.accept_redirects=0
```

Modify active kernel parameters to match:

```bash
/sbin/sysctl -w net.ipv6.conf.all.accept_redirects=0
/sbin/sysctl -w net.ipv6.conf.default.accept_redirects=0
/sbin/sysctl -w net.ipv6.route.flush=1
```

## Default Value

IPv6 redirect acceptance is enabled by default.

## References

- CIS Ubuntu 12.04 LTS Server Benchmark v1.1.0

## Profile

Level 1 - Not Scored
