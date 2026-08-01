# stage: report
# category: CIS_benchmarks


# 6.9 Ensure FTP Server is not enabled (Not Scored)

## Profile Applicability

- Level 1

## Description

The File Transfer Protocol (FTP) provides networked computers with the ability to transfer files.

## Rationale

FTP does not protect the confidentiality of data or authentication credentials. It is recommended sftp be used if file transfer is required. Unless there is a need to run the system as a FTP server (for example, to allow anonymous downloads), it is recommended that the package be deleted to reduce the potential attack surface.

## Audit Procedure

### Using Command Line

Ensure no start conditions listed for `vsftpd`:

```bash
initctl show-config vsftpd vsftpd
```

## Expected Result

No start conditions should be listed for vsftpd.

## Remediation

### Using Command Line

Remove or comment out start lines in `/etc/init/vsftpd.conf`:

```bash
sed -i 's/^start/#start/' /etc/init/vsftpd.conf
```

## Default Value

Not installed by default on Ubuntu 12.04 LTS Server.

## References

- CIS Ubuntu 12.04 LTS Server Benchmark v1.1.0

## Profile

Level 1 - Not Scored
