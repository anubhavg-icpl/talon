# stage: report
# category: CIS_benchmarks


# 6.3 Ensure print server is not enabled (Not Scored)

## Profile Applicability

- Level 1

## Description

The Common Unix Print System (CUPS) provides the ability to print to both local and network printers. A system running CUPS can also accept print jobs from remote systems and print them to local printers. It also provides a web based remote administration capability.

## Rationale

If the system does not need to print jobs or accept print jobs from other systems, it is recommended that CUPS be disabled to reduce the potential attack surface.

## Audit Procedure

### Using Command Line

Ensure no start conditions listed for `cups`:

```bash
initctl show-config cups cups
```

## Expected Result

No start conditions should be listed for cups.

## Remediation

### Using Command Line

Remove or comment out start lines in `/etc/init/cups.conf`:

```bash
sed -i 's/^start/#start/' /etc/init/cups.conf
```

## Default Value

Enabled by default on Ubuntu 12.04 LTS Server.

## References

- CIS Ubuntu 12.04 LTS Server Benchmark v1.1.0
- http://www.cups.org

## Profile

Level 1 - Not Scored
