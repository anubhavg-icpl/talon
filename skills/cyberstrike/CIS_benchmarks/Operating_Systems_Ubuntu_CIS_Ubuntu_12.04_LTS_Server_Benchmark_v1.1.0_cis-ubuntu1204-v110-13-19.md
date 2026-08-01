# stage: report
# category: CIS_benchmarks


# 13.19 Check for Presence of User .netrc Files (Scored)

## Profile Applicability

- Level 1

## Description

The `.netrc` file contains data for logging into a remote host for file transfers via FTP.

## Rationale

The `.netrc` file presents a significant security risk since it stores passwords in unencrypted form. Even if FTP is disabled, user accounts may have brought over `.netrc` files from other systems which could pose a risk to those systems.

## Audit Procedure

### Using Command Line

```bash
#!/bin/bash
for dir in `/bin/cat /etc/passwd |\
/usr/bin/awk -F: '{ print $6 }'`; do
    if [ ! -h "$dir/.netrc" -a -f "$dir/.netrc" ]; then
echo ".netrc file $dir/.netrc exists"
    fi
done
```

## Expected Result

No output should be returned. Any output indicates the presence of `.netrc` files in user home directories.

## Remediation

### Using Command Line

Making global modifications to users' files without alerting the user community can result in unexpected outages and unhappy users. Therefore, it is recommended that a monitoring policy be established to report user `.netrc` files and determine the action to be taken in accordance with site policy.

## Default Value

No `.netrc` files exist by default on Ubuntu 12.04.

## References

- CIS Ubuntu 12.04 LTS Server Benchmark v1.1.0

## Profile

Level 1 - Scored
