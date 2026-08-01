# stage: report
# category: CIS_benchmarks


# 6.6 Control the maximum size of a POST request that will be parsed for parameters (Scored)

## Description

The `maxPostSize` attribute controls the maximum size of a POST request that will be parsed for parameters. Setting this to a reasonable value can help mitigate denial of service attacks that attempt to send large POST bodies.

## Rationale

Limiting the size of POST requests that will be parsed for parameters helps prevent denial of service conditions where an attacker sends extremely large POST bodies to consume server resources.

## Audit Procedure

Review the `$CATALINA_HOME/conf/server.xml` file. Ensure the `maxPostSize` attribute is set on each Connector:

```bash
$ grep -i maxPostSize $CATALINA_HOME/conf/server.xml
```

The `maxPostSize` should be set to a value appropriate for the applications hosted, typically `2097152` (2MB) or less.

## Remediation

In the `$CATALINA_HOME/conf/server.xml` file, add the `maxPostSize` attribute to each Connector element:

```xml
<Connector ... maxPostSize="2097152" />
```

## Default Value

The default value is `2097152` (2MB) in Tomcat 8.

## References

1. https://tomcat.apache.org/tomcat-8.0-doc/config/http.html

## CIS Controls

**v7:**

- 9.4 Apply Host-based Firewalls or Port Filtering
  - Apply host-based firewalls or port filtering tools on end systems, with a default-deny rule that drops all traffic except those services and ports that are explicitly allowed.

## Profile Applicability

- Level 1
