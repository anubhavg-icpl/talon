# stage: post_exploit
# category: mitre_attack


# T1070.004 File Deletion

> **Sub-technique of:** T1070

## High-Level Description

Adversaries may delete files left behind by the actions of their intrusion activity. Malware, tools, or other non-native files dropped or created on a system by an adversary (ex: Ingress Tool Transfer) may leave traces to indicate to what was done within a network and how. Removal of these files can occur during an intrusion, or as part of a post-intrusion process to minimize the adversary's footprint.

There are tools available from the host operating system to perform cleanup, but adversaries may use other tools as well. Examples of built-in Command and Scripting Interpreter functions include <code>del</code> on Windows, <code>rm</code> or <code>unlink</code> on Linux and macOS, and `rm` on ESXi.

## Kill Chain Phase

- Defense Evasion (TA0005)

**Platforms:** ESXi, Linux, macOS, Windows

## What to Check

- [ ] Identify if File Deletion technique is applicable to target environment
- [ ] Check ESXi systems for indicators of File Deletion
- [ ] Check Linux systems for indicators of File Deletion
- [ ] Check macOS systems for indicators of File Deletion
- [ ] Assess detection coverage (1 detection strategies)

## How to Test

### Atomic Red Team Tests

The following tests are from [Atomic Red Team](https://github.com/redcanaryco/atomic-red-team) and provide actionable ways to test this technique:

### Atomic Test 1: Delete a single file - FreeBSD/Linux/macOS

Delete a single file from the temporary directory

**Supported Platforms:** linux, macos

```bash
rm -f #{file_to_delete}
```

**Dependencies:**

- The file must exist in order to be deleted

### Atomic Test 2: Delete an entire folder - FreeBSD/Linux/macOS

Recursively delete the temporary directory and all files contained within it

**Supported Platforms:** linux, macos

```bash
rm -rf #{folder_to_delete}
```

**Dependencies:**

- The folder must exist in order to be deleted

### Atomic Test 3: Overwrite and delete a file with shred

Use the `shred` command to overwrite the temporary file and then delete it

**Supported Platforms:** linux

```bash
shred -u #{file_to_shred}
```

**Dependencies:**

- Check if file already exists

### Atomic Test 4: Delete a single file - Windows cmd

Delete a single file from the temporary directory using cmd.exe.
Upon execution, no output will be displayed. Use File Explorer to verify the file was deleted.

**Supported Platforms:** windows

```cmd
del /f #{file_to_delete}
```

**Dependencies:**

- The file to delete must exist on disk at specified location (#{file_to_delete})

### Atomic Test 5: Delete an entire folder - Windows cmd

Recursively delete a folder in the temporary directory using cmd.exe.
Upon execution, no output will be displayed. Use File Explorer to verify the folder was deleted.

**Supported Platforms:** windows

```cmd
rmdir /s /q #{folder_to_delete}
```

**Dependencies:**

- The file to delete must exist on disk at specified location (#{folder_to_delete})

### Manual Testing

If Atomic Red Team tests are not applicable, manually verify the technique by:

1. **Identify Attack Surface**: Determine if the target environment is susceptible to File Deletion by examining the target platforms (ESXi, Linux, macOS).

2. **Assess Existing Defenses**: Review whether mitigations for T1070.004 are in place. If defenses are absent or misconfigured, this technique may be exploitable.

3. **Execute Test**: Use tools and methods described in the MITRE ATT&CK page and external references below.

## Remediation Guide

No specific mitigations documented for this technique.

## Detection

### Behavioral Detection of Malicious File Deletion

## Risk Assessment

| Finding                            | Severity | Impact          |
| ---------------------------------- | -------- | --------------- |
| File Deletion technique applicable | Low      | Defense Evasion |

## CWE Categories

| CWE ID  | Title                        |
| ------- | ---------------------------- |
| CWE-693 | Protection Mechanism Failure |

## References

- [Microsoft SDelete July 2016](https://docs.microsoft.com/en-us/sysinternals/downloads/sdelete)
- [Atomic Red Team - T1070.004](https://github.com/redcanaryco/atomic-red-team/tree/master/atomics/T1070.004)
- [MITRE ATT&CK - T1070.004](https://attack.mitre.org/techniques/T1070/004)
