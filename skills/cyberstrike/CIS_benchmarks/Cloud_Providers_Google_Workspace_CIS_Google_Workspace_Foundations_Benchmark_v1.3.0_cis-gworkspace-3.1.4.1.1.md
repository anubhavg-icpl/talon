# stage: report
# category: CIS_benchmarks


# Ensure external filesharing in Google Chat and Hangouts is disabled (Manual)

## Description

Control how files are shared externally in Google Chat and Hangouts.

## Rationale

Files often contain confidential information, and some organizations, particularly in regulated industries, need to control the flow of this information within and outside of their organization.

## Impact

Users will not be able to share files via chat externally.

## Audit Procedure

### Using Google Workspace Admin Console

1. Log in to https://admin.google.com as an administrator
2. Select `Apps`
3. Select `Google Chat and classic Hangouts`
4. Select `Chat File Sharing`
5. Under `Setting`, verify `External filesharing` is set to `No files`

### Expected Result

`External filesharing` should be set to `No files`.

## Remediation

### Using Google Workspace Admin Console

1. Log in to https://admin.google.com as an administrator
2. Select `Apps`
3. Select `Google Chat and classic Hangouts`
4. Select `Chat File Sharing`
5. Under `Setting`, set `External filesharing` to `No files`
6. Select `Save`

## Default Value

`External filesharing` is `Allow all files`

## CIS Controls

| Controls Version | Control                                                                         | IG 1 | IG 2 | IG 3 |
| ---------------- | ------------------------------------------------------------------------------- | ---- | ---- | ---- |
| v8               | 4.8 Uninstall or Disable Unnecessary Services on Enterprise Assets and Software |      | x    | x    |

## Profile

Level 1
