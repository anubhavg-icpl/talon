# stage: report
# category: CIS_benchmarks


# CIS Red Hat OpenShift Container Platform Benchmark v1.9.0 - Control 1.2.31

## Profile Applicability

- **Level:** 1

## Description

OpenShift supported an option called `unsupportedConfigOverrides` that allowed users to opt into unsupported behavior. This option is no longer supported by OpenShift and should not be used.

## Rationale

Users should stop using deprecated and unmaintained features in favor of supported features.

## Impact

None. The feature is set to `null` by default and isn't used by default.

## Audit Procedure

Make sure the `unsupportedConfigOverrides` in your deployment returns `null` using the following command:

```bash
oc get kubeapiserver/cluster -o jsonpath='{.spec.unsupportedConfigOverrides}'
```

The output should return `null`. Any other return value is a finding and you should migrate away from that particular configuration.

## Remediation

None.

## Default Value

By default, OpenShift sets this value to `null` and doesn't support overriding configuration with unsupported features.

## References

1. https://access.redhat.com/solutions/5170671

## CIS Controls

| Controls Version | Control                                               | IG 1 | IG 2 | IG 3 |
| ---------------- | ----------------------------------------------------- | ---- | ---- | ---- |
| v8               | 2.2 Ensure Authorized Software is Currently Supported | \*   | \*   | \*   |
| v7               | 2.2 Ensure Software is Supported by Vendor            | \*   | \*   | \*   |

## MITRE ATT&CK Mappings

This control does not have specific MITRE ATT&CK mappings in the benchmark.

## Profile

**Level 1** (Manual)
