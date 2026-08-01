# stage: post_exploit
# category: CIS_benchmarks


# CIS Control 2.9: Ensure Granular Policy Creation (Manual)

## Profile Applicability

- **Level 2**

## Description

Granular policies are meticulously tailored to AWS resources, ensuring precision in access control measures.

## Rationale

Emphasizing granular policies in AWS ensures that access control measures are precisely aligned with the requirements of each resource, bolstering security and minimizing unauthorized access. By tailoring policies to specific resources, organizations can adhere more closely to the principle of least privilege, mitigating risks and maintaining compliance with regulatory standards.

## Audit Procedure

### Via AWS CLI

\`\`\`bash

# Review IAM policies for granularity

aws iam list-policies --scope Local

# Examine policy document

aws iam get-policy-version \
 --policy-arn <POLICY_ARN> \
 --version-id <VERSION_ID>
\`\`\`

## Remediation

Review existing IAM policies to identify those that are overly broad or lack granularity. Refine these policies to restrict permissions to only the resources and actions necessary for each user or group.

## References

1. [IAM Policy Tags](https://docs.aws.amazon.com/tag-editor/latest/userguide/tags-in-iam-policies.html)

## CIS Controls

| Controls Version | Control                                            | IG 1 | IG 2 | IG 3 |
| ---------------- | -------------------------------------------------- | ---- | ---- | ---- |
| v8               | 6.8 Define and Maintain Role-Based Access Control  |      |      | ●    |
| v7               | 16.2 Configure Centralized Point of Authentication |      | ●    | ●    |
