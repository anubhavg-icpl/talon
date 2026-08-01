# stage: report
# category: CIS_benchmarks


# Ensure spam filters are not bypassed for internal senders (Manual)

## Description

You can configure your advanced Gmail settings to bypass, or not bypass, spam filters for messages received from internal senders.

## Rationale

Turning off this setting reduces the risk of spoofing and phishing/whaling.

## Impact

Your users will be better protected by filtering their email for spam and minimizing the chances for spoofing and phishing/whaling attacks.

## Audit Procedure

### Using Google Workspace Admin Console

1. Log in to https://admin.google.com as an administrator
2. Select `Apps`
3. Select `Google Workspace`
4. Select `Gmail`
5. Select `Spam, phishing, and malware`
6. Under `Spam`, select `Configure`
7. Ensure `Bypass spam filters for messages received from internal senders.` is **unchecked**

### Expected Result

`Bypass spam filters for messages received from internal senders.` should be unchecked.

## Remediation

### Using Google Workspace Admin Console

1. Log in to https://admin.google.com as an administrator
2. Select `Apps`
3. Select `Google Workspace`
4. Select `Gmail`
5. Select `Spam, phishing, and malware`
6. Under `Spam`, select `Configure`
7. Set `Bypass spam filters for messages received from internal senders.` to **unchecked**
8. Select `Save`

## Default Value

`Bypass spam filters for messages received from internal senders.` is **checked**

## CIS Controls

This control does not have explicit CIS Controls mappings in the PDF.

## Profile

Level 1
