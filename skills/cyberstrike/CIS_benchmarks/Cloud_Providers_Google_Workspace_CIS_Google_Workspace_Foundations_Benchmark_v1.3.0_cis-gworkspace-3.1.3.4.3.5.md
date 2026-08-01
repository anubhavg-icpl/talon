# stage: report
# category: CIS_benchmarks


# 3.1.3.4.3.5 Ensure groups are protected from inbound emails spoofing your domain

## Overview

| Property                  | Value                                        |
| ------------------------- | -------------------------------------------- |
| **CIS ID**                | 3.1.3.4.3.5                                  |
| **Level**                 | L1                                           |
| **Profile Applicability** | Enterprise Level 1                           |
| **Assessment Type**       | Manual                                       |
| **Section**               | Gmail > Safety > Spoofing and authentication |

## Description

If a group receives an email that is spoofing your domain it is sent to the spam folder.

## Rationale

You should protect your groups from any emails that spoofing your domain.

## Impact

Emails that are spoofing your domain and are received by a group are sent to the spam folder.

## Default Value

Not specified in benchmark.

## Audit

To verify this setting via the Google Workspace Admin Console:

1. Log in to [https://admin.google.com](https://admin.google.com) as an administrator
2. Select **Apps**
3. Select **Google Workspace**
4. Select **Gmail**
5. Under `Safety` - `Spoofing and authentication`, ensure `Protect your Groups from inbound emails spoofing your domain` is **checked**
6. Ensure `Action` is set to `Move email to spam`

## Remediation

To configure this setting via the Google Workspace Admin Console:

1. Log in to [https://admin.google.com](https://admin.google.com) as an administrator
2. Select **Apps**
3. Select **Google Workspace**
4. Select **Gmail**
5. Under `Safety` - `Spoofing and authentication`, set `Protect your Groups from inbound emails spoofing your domain` to **checked**
6. Set `Action` to `Move email to spam`
7. Select **Save**

## CIS Controls

| Controls Version | Control                                                   | IG 1 | IG 2 | IG 3 |
| ---------------- | --------------------------------------------------------- | ---- | ---- | ---- |
| v8               | 9.5 Implement DMARC                                       |      | x    | x    |
| v7               | 7.8 Implement DMARC and Enable Receiver-Side Verification |      | x    | x    |
