# stage: post_exploit
# category: CIS_benchmarks


# 18.3 Ensure Storage Explorer is using the latest version (Manual)

## Description

Ensure all users accessing Azure Storage resources with Storage Explorer are using the latest version of the software, applying updates promptly to safeguard against new vulnerabilities and benefit from the latest security enhancements.

## Rationale

Using the latest version of Storage Explorer is essential for safeguarding access to Azure Storage resources.

## Impact

Using the latest version of Storage Explorer is free and requires minimal administrative effort.

## Audit Procedure

### Audit from Storage Explorer for MacOS

1. Go to `Storage Explorer`.
2. From the menu bar, click `Microsoft Azure Storage Explorer`.
3. Click `About`.
4. Check the `Version` listed.
5. Compare the installed version with the latest released version from https://github.com/microsoft/AzureStorageExplorer/releases.
6. Ensure that the installed Storage Explorer version matches the latest released version.

## Expected Result

The installed Storage Explorer version should match the latest released version available at https://github.com/microsoft/AzureStorageExplorer/releases.

## Remediation

### Remediate from Storage Explorer for MacOS

1. Go to `Storage Explorer`.
2. From the menu bar, click `Microsoft Azure Storage Explorer`.
3. Click `Check for Updates`.
4. Follow the instructions to install the latest version of Storage Explorer.

## Default Value

By default, new installations of Storage Explorer will utilize the latest released version.

## References

1. https://github.com/microsoft/AzureStorageExplorer/releases

## Profile

Level 1 | Manual
