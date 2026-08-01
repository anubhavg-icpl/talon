# stage: report
# category: CIS_benchmarks


# 3.1 Set a nondeterministic Shutdown command value (Scored)

## Description

Tomcat listens on TCP port 8005 to accept shutdown requests. By connecting to this port and sending the SHUTDOWN command, all applications within Tomcat are halted. The shutdown port is not exposed to the network as it is bound to the loopback interface. It is recommended that a nondeterministic value be set for the shutdown attribute in `$CATALINA_HOME/conf/server.xml`.

## Rationale

Setting the shutdown attribute to a nondeterministic value will prevent malicious local users from shutting down Tomcat.

## Audit Procedure

Perform the following to determine if the shutdown port is configured to use the default shutdown command:

1. Ensure the shutdown attribute in `$CATALINA_HOME/conf/server.xml` is not set to SHUTDOWN.

```bash
$ cd $CATALINA_HOME/conf
$ grep 'shutdown[[:space:]]*=[[:space:]]*"SHUTDOWN"' server.xml
```

## Remediation

1. Update the shutdown attribute in `$CATALINA_HOME/conf/server.xml` as follows:

```xml
<Server port="8005" shutdown="NONDETERMINISTICVALUE">
```

Note: NONDETERMINISTICVALUE should be replaced with a sequence of random characters.

## Default Value

The default value for the shutdown attribute is SHUTDOWN.

## References

1. http://tomcat.apache.org/tomcat-8.0-doc/config/server.html

## CIS Controls

- Not mapped in this benchmark version

## Profile Applicability

- Level 1
