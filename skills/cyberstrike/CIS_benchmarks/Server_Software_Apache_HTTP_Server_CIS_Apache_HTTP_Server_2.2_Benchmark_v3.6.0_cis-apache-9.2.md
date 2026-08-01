# stage: report
# category: CIS_benchmarks


# Ensure KeepAlive Is Enabled

## Description

The `KeepAlive` directive controls whether Apache will reuse the same TCP connection per client to process subsequent HTTP requests from that client. It is recommended that the `KeepAlive` directive be set to `On`.

## Rationale

Allowing per-client reuse of TCP sockets reduces the amount of system and network resources required to serve requests. This efficiency gain may improve a server's resiliency to DoS attacks.

## Impact

None documented

## Audit Procedure

Perform the following steps to determine if the recommended state is implemented:

Verify that the `KeepAlive` directive in the Apache configuration either has a value of `On` or is not present. If the directive is not present, the default value is `On`.

## Remediation

Perform the following to implement the recommended state:

Add or modify the `KeepAlive` directive in the Apache configuration to have a value of `On`.

```
KeepAlive On
```

## Default Value

```
KeepAlive On
```

## References

1. https://httpd.apache.org/docs/2.2/mod/core.html#keepalive

## CIS Controls

**Version 6**

9 Limitation and Control of Network Ports, Protocols, and Services
Limitation and Control of Network Ports, Protocols, and Services

**Version 7**

5.1 Establish Secure Configurations
Maintain documented, standard security configuration standards for all authorized
operating systems and software.

## Profile

Level 1 | Scored
