# stage: report
# category: CIS_benchmarks


# 7.2.7 Enable RFC-recommended Source Route Validation (Scored)

## Profile Applicability

- Level 1

## Description

Setting `net.ipv4.conf.all.rp_filter` and `net.ipv4.conf.default.rp_filter` to 1 forces the Linux kernel to utilize reverse path filtering on a received packet to determine if the packet was valid. Essentially, with reverse path filtering, if the return packet does not go out the same interface that the corresponding source packet came from, the packet is dropped (and logged if `log_martians` is set).

## Rationale

Setting these flags is a good way to deter attackers from sending your server bogus packets that cannot be responded to. One instance where this feature breaks down is if asymmetrical routing is employed. This would occur when using dynamic routing protocols (bgp, ospf, etc) on your system. If you are using asymmetrical routing on your server, you will not be able to enable this feature without breaking the routing.

## Audit Procedure

### Using Command Line

Perform the following to determine if RFC-recommended source route validation is enabled.

```bash
/sbin/sysctl net.ipv4.conf.all.rp_filter
/sbin/sysctl net.ipv4.conf.default.rp_filter
```

## Expected Result

```
net.ipv4.conf.all.rp_filter = 1
net.ipv4.conf.default.rp_filter = 1
```

## Remediation

### Using Command Line

Set the `net.ipv4.conf.all.rp_filter` and `net.ipv4.conf.default.rp_filter` parameters to 1 in `/etc/sysctl.conf`:

```bash
net.ipv4.conf.all.rp_filter=1
net.ipv4.conf.default.rp_filter=1
```

Modify active kernel parameters to match:

```bash
/sbin/sysctl -w net.ipv4.conf.all.rp_filter=1
/sbin/sysctl -w net.ipv4.conf.default.rp_filter=1
/sbin/sysctl -w net.ipv4.route.flush=1
```

## Default Value

Reverse path filtering is disabled by default.

## References

- CIS Ubuntu 12.04 LTS Server Benchmark v1.1.0

## Profile

Level 1 - Scored
