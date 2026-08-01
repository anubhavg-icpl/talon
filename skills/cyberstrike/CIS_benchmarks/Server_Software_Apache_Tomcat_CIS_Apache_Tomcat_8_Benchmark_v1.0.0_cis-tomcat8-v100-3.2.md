# stage: report
# category: CIS_benchmarks


# 3.2 Disable the Shutdown port (Not Scored)

## Description

Tomcat listens on TCP port 8005 to accept shutdown requests. By connecting to this port and sending the SHUTDOWN command, all applications within Tomcat are halted. The shutdown port is not exposed to the network as it is bound to the loopback interface. If this functionality is not used, it is recommended that the Shutdown port be disabled.

## Rationale

Disabling the Shutdown port will eliminate the risk of malicious local entities using the shutdown command to disable the Tomcat server.

## Audit Procedure

1. Ensure the port attribute in `$CATALINA_HOME/conf/server.xml` is set to -1.

```bash
$ cd $CATALINA_HOME/conf/
$ grep '<Server[[:space:]]\+[^>]*port[[:space:]]*=[[:space:]]*"-1"' server.xml
```

## Remediation

1. Set the port to -1 in the `$CATALINA_HOME/conf/server.xml` file:

```xml
<Server port="-1" shutdown="$SHUTDOWN">
```

## Default Value

The shutdown port is enabled on TCP port 8005, bound to the loopback address.

## References

1. http://tomcat.apache.org/tomcat-8.0-doc/config/server.html

## CIS Controls

- Not mapped in this benchmark version

## Profile Applicability

- Level 2
