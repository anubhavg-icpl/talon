# stage: report
# category: CIS_benchmarks


# 3.1.2.1.2.3 Ensure shared drive file access is restricted to members only

## Level

L1

## Profile Applicability

- Enterprise Level 1

## Type

Manual

## Description

Shared drive file access should be restricted to that shared drive's members.

## Rationale

Preventing unauthorized users from access sensitive data is paramount in preventing unauthorized or unintentional information disclosures.

## Impact

Disabling this feature will prevent shared drive non-members from accessing content in shared drives where they are not a member.

## Audit

To verify this setting via the Google Workspace Admin Console:

1. Log in to `https://admin.google.com` as an administrator
2. Select `Apps`
3. Select `Google Workspace`
4. Select `Drive and Docs`
5. Select `Sharing settings`
6. Under `Shared drive creation`, ensure `Allow people who aren't shared drive members to be added to files` is `unchecked`

## Remediation

To configure this setting via the Google Workspace Admin Console:

1. Log in to `https://admin.google.com` as an administrator
2. Select `Apps`
3. Select `Google Workspace`
4. Select `Drive and Docs`
5. Select `Sharing settings`
6. Under `Shared drive creation`, set `Allow people who aren't shared drive members to be added to files` to `unchecked`
7. Select `Save`

## Default Value

`Allow people who aren't shared drive members to be added to files` is `checked`

## CIS Controls

| Controls Version | Control                                               | IG 1 | IG 2 | IG 3 |
| ---------------- | ----------------------------------------------------- | ---- | ---- | ---- |
| v8               | 3.3 Configure Data Access Control Lists               | x    | x    | x    |
| v7               | 14.6 Protect Information through Access Control Lists | x    | x    | x    |
