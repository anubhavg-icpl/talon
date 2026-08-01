# stage: report
# category: CIS_benchmarks


# Ensure per-user outbound gateways is disabled (Manual)

## Description

A per-user outbound gateway is a mail server, other than the Google Workspace mail servers, that delivers outgoing mail for a user in your domain.

## Rationale

Mail sent via external SMTP will circumvent your outbound gateway.

## Impact

Care should be taken before implementation to ensure there is no business need for mail sent via external SMTP gateway.

## Audit Procedure

### Using Google Workspace Admin Console

1. Log in to https://admin.google.com as an administrator
2. Select `Apps`
3. Select `Google Workspace`
4. Select `Gmail`
5. Under `End User Access` - `Allow per-user outbound gateways`, ensure `Allow users to send mail through an external SMTP server when configuring a "from" address hosted outside your email domain` is **unchecked**

### Expected Result

`Allow users to send mail through an external SMTP server when configuring a "from" address hosted outside your email domain` should be unchecked.

## Remediation

### Using Google Workspace Admin Console

1. Log in to https://admin.google.com as an administrator
2. Select `Apps`
3. Select `Google Workspace`
4. Select `Gmail`
5. Under `End User Access` - `Allow per-user outbound gateways`, set `Allow users to send mail through an external SMTP server when configuring a "from" address hosted outside your email domain` to **unchecked**
6. Select `Save`

## Default Value

`Allow users to send mail through an external SMTP server when configuring a "from" address hosted outside your email domain` is **unchecked**

## CIS Controls

This control does not have explicit CIS Controls mappings in the PDF.

## Profile

Level 1
