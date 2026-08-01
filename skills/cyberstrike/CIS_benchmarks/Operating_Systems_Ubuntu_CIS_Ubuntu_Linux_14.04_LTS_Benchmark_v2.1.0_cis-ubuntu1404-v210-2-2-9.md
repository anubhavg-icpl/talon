# stage: report
# category: CIS_benchmarks


# 2.2.9 Ensure FTP Server is not enabled (Scored)

## Profile Applicability

- Level 1 - Server
- Level 1 - Workstation

## Description

The File Transfer Protocol (FTP) provides networked computers with the ability to transfer files.

## Rationale

FTP does not protect the confidentiality of data or authentication credentials. It is recommended sftp be used if file transfer is required. Unless there is a need to run the system as a FTP server (for example, to allow anonymous downloads), it is recommended that the package be deleted to reduce the potential attack surface.

## Audit Procedure

Run the following commands to verify no start conditions listed for `vsftpd`:

```bash
initctl show-config vsftpd
```

Verify the output shows `vsftpd` with no start conditions.

## Expected Result

The `vsftpd` service should have no start conditions listed.

## Remediation

Remove or comment out start lines in `/etc/init/vsftpd.conf`:

```bash
#start on runlevel [2345] or net-device-up IFACE!=lo
```

**Note:** Additional FTP servers also exist and should be audited.

## Default Value

FTP server is not enabled by default.

## References

- CIS Controls: 9.1 Limit Open Ports, Protocols, and Services

## Profile

- Level 1 - Server
- Level 1 - Workstation
