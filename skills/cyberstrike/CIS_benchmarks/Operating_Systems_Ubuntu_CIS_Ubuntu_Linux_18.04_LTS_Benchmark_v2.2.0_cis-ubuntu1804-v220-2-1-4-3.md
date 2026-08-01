# stage: report
# category: CIS_benchmarks


# 2.1.4.3 Ensure ntp is running as user ntp (Automated)

## Profile

- Level 1 - Server
- Level 1 - Workstation

## Description

The ntp package is installed with a dedicated user account `ntp`. This account is granted the access required by the ntpd daemon.

Note:

- If chrony or systemd-timesyncd are used, ntp should be removed and this section skipped
- This recommendation only applies if ntp is in use on the system
- Only one time synchronization method should be in use on the system

## Rationale

The ntpd daemon should run with only the required privilege.

## Audit Procedure

### Command Line

IF ntp is in use on the system, run the following command to verify the ntpd daemon is being run as the user `ntp`:

```bash
# ps -ef | awk '/[n]tpd/ && $1!="ntp") { print $1 }'
```

Nothing should be returned.

Run the following command to verify the `RUNASUSER=` is set to `ntp` in `/usr/lib/ntp/ntp-systemd-wrapper`:

```bash
# grep -P -- '^\h*RUNASUSER=' /usr/lib/ntp/ntp-systemd-wrapper
RUNASUSER=ntp
```

## Expected Result

- First command: No output (empty result)
- Second command: `RUNASUSER=ntp`

## Remediation

### Command Line

Add or edit the following line in `/usr/lib/ntp/ntp-systemd-wrapper`:

```
RUNASUSER=ntp
```

Run the following command to restart ntp.service:

```bash
# systemctl restart ntp.service
```

OR

If another time synchronization service is in use on the system, run the following command to remove ntp:

```bash
# apt purge ntp
```

## Default Value

user ntp

## References

1. http://www.ntp.org/
2. NIST SP 800-53 Rev. 5: AU-8

## CIS Controls

- v8: 8.4 - Standardize Time Synchronization
- v7: 6.1 - Utilize Three Synchronized Time Sources
