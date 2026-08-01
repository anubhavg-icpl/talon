# stage: report
# category: CIS_benchmarks


# Ensure the KeepAliveTimeout Is Set Properly

## Description

The `KeepAliveTimeout` directive specifies the number of seconds Apache will wait for a subsequent request before closing a connection that is being kept alive.

## Rationale

Reducing the number of seconds that Apache HTTP server will keep unused resources allocated will increase the availability of resources to serve other requests. This efficiency gain may improve a server's resiliency to DoS attacks.

## Impact

None documented

## Audit Procedure

Perform the following steps to determine if the recommended state is implemented:

Verify that the `KeepAliveTimeout` directive in the Apache configuration either has a value of `15` or less or is not present. If the directive is not present, the default value is `15` seconds.

## Remediation

Perform the following to implement the recommended state:

Add or modify the `KeepAliveTimeout` directive in the Apache configuration to have a value of `15` or less.

```
KeepAliveTimeout 15
```

## Default Value

```
KeepAliveTimeout 15
```

## References

1. https://httpd.apache.org/docs/2.2/mod/core.html#keepalivetimeout

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
