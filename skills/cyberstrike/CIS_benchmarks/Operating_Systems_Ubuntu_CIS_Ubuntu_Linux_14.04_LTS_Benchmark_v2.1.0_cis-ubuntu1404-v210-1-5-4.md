# stage: report
# category: CIS_benchmarks


# 1.5.4 Ensure prelink is disabled (Scored)

## Profile Applicability

- Level 1 - Server
- Level 1 - Workstation

## Description

`prelink` is a program that modifies ELF shared libraries and ELF dynamically linked binaries in such a way that the time needed for the dynamic linker to perform relocations at startup significantly decreases.

## Rationale

The prelinking feature can interfere with the operation of AIDE, because it changes binaries. Prelinking can also increase the vulnerability of the system if a malicious user is able to compromise a common library such as libc.

## Audit Procedure

Run the following command and verify `prelink` is not installed:

```bash
dpkg -s prelink
```

## Expected Result

The output should indicate that the package is not installed.

## Remediation

Run the following command to restore binaries to normal:

```bash
prelink -ua
```

Run the following command to uninstall `prelink`:

```bash
apt-get remove prelink
```

## Default Value

`prelink` is not installed by default.

## References

- CIS Controls: 3.5 Use File Integrity Tools For Critical System Files

## Profile

- Level 1 - Server
- Level 1 - Workstation
