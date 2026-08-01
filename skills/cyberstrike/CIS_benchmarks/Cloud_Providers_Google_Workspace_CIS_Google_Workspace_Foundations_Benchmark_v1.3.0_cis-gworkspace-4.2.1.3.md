# stage: report
# category: CIS_benchmarks


# 4.2.1.3 Ensure internal apps can access Google Workspace APIs

## Profile Applicability

- Enterprise Level 1

## Description

Enable access to Google Workspace APIs for customer-owned / developed applications.

## Rationale

All organization-built internal apps (owned by your organization), can be trusted to access restricted Google Workspace APIs. That way, the organization does not have to trust them all individually.

## Audit

To verify this setting via the Google Workspace Admin Console:

1. Log in to `https://admin.google.com` as an administrator
2. Select `Security`
3. Select `Access and Data Control`
4. Select `API Controls`, then select `App access control`
5. Under `Settings`, verify `Trust internal, domain-owned apps` is `selected`

## Remediation

To configure this setting via the Google Workspace Admin Console:

1. Log in to `https://admin.google.com` as an administrator
2. Select `Security`
3. Select `Access and Data Control`
4. Select `API Controls`, then select `App access control`
5. Under `Settings`, select `Trust internal, domain-owned apps`
6. Select `Save`

## Default Value

`Trust internal, domain-owned apps` is `selected`

## CIS Controls

| Controls Version | Control                                               | IG 1 | IG 2 | IG 3 |
| ---------------- | ----------------------------------------------------- | ---- | ---- | ---- |
| v8               | 3.3 Configure Data Access Control Lists               | x    | x    | x    |
| v7               | 14.6 Protect Information through Access Control Lists | x    | x    | x    |
