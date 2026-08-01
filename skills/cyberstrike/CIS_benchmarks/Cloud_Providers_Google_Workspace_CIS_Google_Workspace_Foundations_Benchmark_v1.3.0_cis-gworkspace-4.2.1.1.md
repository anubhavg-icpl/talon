# stage: report
# category: CIS_benchmarks


# 4.2.1.1 Ensure application access to Google services is restricted

## Profile Applicability

- Enterprise Level 2

## Description

Prevent unrestricted application access to Google services.

## Rationale

You can restrict (or leave unrestricted) access to most Workspace services, including Google Cloud Platform services such as Machine Learning. For Gmail and Google Drive, you can specifically restrict access to high-risk scopes (for example, sending Gmail or deleting files in Drive). While users are prompted to consent to apps, if an app uses restricted scopes and you haven't specifically trusted it, users can't add it.

## Impact

The potential impact associated with implementation of this setting is that any previously installed apps that you haven't trusted stop working and tokens are revoked. When a user tries to install an app that has a restricted scope, they're notified that it's blocked.

## Audit

To verify this setting via the Google Workspace Admin Console:

1. Log in to `https://admin.google.com` as an administrator
2. Select `Security`
3. Select `Access and Data Control`
4. Select `API Controls`, then select `App access control`
5. Under `Overview`, select `MANAGE GOOGLE SERVICES`
6. Ensure `ALL` applicable Google Services have `Restricted` in the `Access` column

## Remediation

To configure this setting via the Google Workspace Admin Console:

1. Log in to `https://admin.google.com` as an administrator
2. Select `Security`
3. Select `Access and Data Control`
4. Select `API Controls`, then select `App access control`
5. Under `Overview`, select `MANAGE GOOGLE SERVICES`
6. Select `ALL` applicable Google Services
7. Click `Change access`
8. Select `Restricted: Only trusted apps can access a service`

## Default Value

`Access` is `Unrestricted`

## CIS Controls

| Controls Version | Control                                               | IG 1 | IG 2 | IG 3 |
| ---------------- | ----------------------------------------------------- | ---- | ---- | ---- |
| v8               | 3.3 Configure Data Access Control Lists               | x    | x    | x    |
| v7               | 14.6 Protect Information through Access Control Lists | x    | x    | x    |
