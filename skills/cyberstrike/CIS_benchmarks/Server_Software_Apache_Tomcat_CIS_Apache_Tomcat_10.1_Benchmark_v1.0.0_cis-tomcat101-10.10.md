# stage: report
# category: CIS_benchmarks


# 10.10 Configure maxHttpHeaderSize (Automated)

## Description

The `maxHttpHeaderSize` limits the size of the request and response headers defined in bytes.

## Rationale

Limiting the size of the header request can help protect against Denial of Service (DoS) requests.

## Audit Procedure

Locate each `maxHttpHeaderSize` setting in `$CATALINA_HOME/conf/server.xml` and verify that they are set to `8192`.

```
# grep maxHttpHeaderSize $CATALINA_HOME/conf/server.xml
```

## Remediation

Set `maxHttpHeaderSize` for each connector in `$CATALINA_HOME/conf/server.xml` to the appropriate setting.

```
maxHttpHeaderSize="8192"
```

## Default Value

The default is `8192` bytes.

## References

1. https://tomcat.apache.org/tomcat-9.0-doc/config/http.html

## CIS Controls

| Controls Version | Control                                                                                                                                                                                               | IG 1 | IG 2 | IG 3 |
| ---------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | ---- | ---- | ---- |
| v8               | 12.6 Use of Secure Network Management and Communication Protocols<br>Use secure network management and communication protocols (e.g., 802.1X, Wi-Fi Protected Access 2 (WPA2) Enterprise or greater). |      | ●    | ●    |
| v7               | 5.1 Establish Secure Configurations<br>Maintain documented, standard security configuration standards for all authorized operating systems and software.                                              | ●    | ●    | ●    |

## Profile

Level 2
