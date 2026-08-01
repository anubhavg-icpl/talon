# stage: report
# category: CIS_benchmarks


# 5.2 Client Encryption

## Profile Applicability

- Level 1 - Cassandra on Linux
- Level 2 - Cassandra on Linux

## Description

Cassandra offers the option to encrypt data in transit between the client and nodes on the cluster. By default client encryption is turned off.

## Rationale

Data in transit between the client and node on the cluster should be encrypted to avoid network snooping, whether legitimate or not.

## Audit

The Cassandra configuration files can be found in the conf directory of tarballs. For packages, the configuration files will be located in /etc/cassandra.
Open up the cassandra.yaml file, look for client_encryption_options section.
Look for enabled: and optional:

```
enabled: true

optional: false
```

If neither is true, then all client connections are unencrypted which makes this a finding.
If enabled is true and optional is false, then all client connections must be encrypted which makes this not a finding.
If enabled is false and optional is true, then enabled wins and all client connections are unencrypted which makes this a finding.
If both are set to true, then both unencrypted and encrypted connections are allowed on the same port which makes this not a finding.

## Remediation

The client encryption should be implemented before anyone accesses the Cassandra server.
To enable the client encryption mechanism:

1. Stop the Cassandra database.
2. If not done so already, build out your keystore and truststore.
3. Modify cassandra.yaml file to modify/add entries under client_encryption_options:

## Default Value

Client encryption is disabled by default.

## References

Not specified in the benchmark.

## CIS Controls

- v8: 3.10 Encrypt Sensitive Data in Transit
- v7: 14.4 Encrypt All Sensitive Information in Transit

## Profile

- Level 1 | Automated
