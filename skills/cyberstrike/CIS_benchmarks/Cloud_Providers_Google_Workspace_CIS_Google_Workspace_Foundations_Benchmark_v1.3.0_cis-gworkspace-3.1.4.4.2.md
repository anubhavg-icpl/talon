# stage: report
# category: CIS_benchmarks


# Ensure allow users to add and use incoming webhooks is disabled (Manual)

## Description

Allow users to configure incoming webhooks and developers to call incoming webhooks to post content. Incoming webhooks let you send asynchronous messages into Google Chat from applications that aren't Chat apps.

## Rationale

Webhook usage should be carefully controlled (vetted and approved) since a malicious application could send bogus information to exposed webhooks and ultimately these users.

## Impact

By default users will have exposed webhooks.

## Audit Procedure

### Using Google Workspace Admin Console

1. Log in to https://admin.google.com as an administrator
2. Select `Apps`
3. Select `Google Chat and classic Hangouts`
4. Select `Chat apps`
5. Under `Chat apps access settings`, verify `Allow users to add and use incoming webhooks` is **OFF**

### Expected Result

`Allow users to add and use incoming webhooks` should be OFF.

## Remediation

### Using Google Workspace Admin Console

1. Log in to https://admin.google.com as an administrator
2. Select `Apps`
3. Select `Google Chat and classic Hangouts`
4. Select `Chat apps`
5. Under `Chat apps access settings`, set `Allow users to add and use incoming webhooks` to **OFF**

## Default Value

`Allow users to add and use incoming webhooks` is **ON**

## CIS Controls

This control does not have explicit CIS Controls mappings in the PDF.

## Profile

Level 1
