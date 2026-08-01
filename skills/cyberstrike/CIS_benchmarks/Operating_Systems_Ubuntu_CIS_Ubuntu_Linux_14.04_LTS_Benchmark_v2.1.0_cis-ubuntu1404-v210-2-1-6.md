# stage: report
# category: CIS_benchmarks


# 2.1.6 Ensure rsh server is not enabled (Scored)

## Profile Applicability

- Level 1 - Server
- Level 1 - Workstation

## Description

The Berkeley `rsh-server` (`rsh`, `rlogin`, `rexec`) package contains legacy services that exchange credentials in clear-text.

## Rationale

These legacy services contain numerous security exposures and have been replaced with the more secure SSH package.

## Audit Procedure

Verify the rsh services are not enabled. Run the following commands and verify results are as indicated:

```bash
grep -R "^shell" /etc/inetd.*
grep -R "^login" /etc/inetd.*
grep -R "^exec" /etc/inetd.*
```

No results should be returned.

Check `/etc/xinetd.conf` and `/etc/xinetd.d/*` and verify all `rsh`, `rlogin`, and `rexec` services have `disable = yes` set.

## Expected Result

No output should be returned from any of the grep commands. All rsh, rlogin, and rexec services in xinetd should have `disable = yes`.

## Remediation

Comment out or remove any lines starting with `shell`, `login`, or `exec` from `/etc/inetd.conf` and `/etc/inetd.d/*`.

Set `disable = yes` on all `rsh`, `rlogin`, and `rexec` services in `/etc/xinetd.conf` and `/etc/xinetd.d/*`.

## Default Value

rsh services are not enabled by default.

## References

- CIS Controls: 3.4 Use Only Secure Channels For Remote System Administration
- CIS Controls: 9.1 Limit Open Ports, Protocols, and Services

## Profile

- Level 1 - Server
- Level 1 - Workstation
