# stage: report
# category: CIS_benchmarks


# CIS Ubuntu Linux 16.04 LTS Benchmark v2.0.0 - Control 6.2.15

## Profile

- **Level:** 1 - Server
- **Level:** 1 - Workstation
- **Assessment Status:** Automated

## Description

Although the `useradd` program will not let you create a duplicate user name, it is possible for an administrator to manually edit the `/etc/passwd` file and change the user name.

## Rationale

If a user is assigned a duplicate user name, it will create and have access to files with the first UID for that username in `/etc/passwd`. For example, if "test4" has a UID of 1000 and a subsequent "test4" entry has a UID of 2000, logging in as "test4" will use UID 1000. Effectively, the UID is shared, which is a security problem.

## Audit Procedure

### Command Line

Run the following script and verify no results are returned:

```bash
#!/bin/bash
cut -d: -f1 /etc/passwd | sort | uniq -d | while read -r x; do
  echo "Duplicate login name $x in /etc/passwd"
done
```

## Expected Result

No output should be returned.

## Remediation

### Command Line

Based on the results of the audit script, establish unique user names for the users. File ownerships will automatically reflect the change as long as the users have unique UIDs.

## Default Value

N/A

## References

N/A

## CIS Controls

| Controls Version | Control                                                              | IG 1 | IG 2 | IG 3 |
| ---------------- | -------------------------------------------------------------------- | ---- | ---- | ---- |
| v7               | 16 Account Monitoring and Control<br/>Account Monitoring and Control |      |      |      |
