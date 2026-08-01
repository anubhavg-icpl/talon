# stage: post_exploit
# category: CIS_benchmarks


# Ensure Amazon ECS clusters are tagged

## Description

Ensure all Amazon ECS clusters have resource tags to facilitate asset management, tracking, and compliance.

## Rationale

Consistent tagging supports compliance and helps identify unauthorized or misconfigured resources.

## Impact

There is minimal administrative overhead associated with implementing and maintaining resource tags.

## Audit Procedure

### Using AWS Console

1. Login to the ECS console using https://console.aws.amazon.com/ecs/.
2. In the left panel, click `Clusters`.
3. Click the name of a cluster.
4. Click `Tags`.
5. Ensure at least one tag is listed that does not begin with `aws:`. Tags prefixed with `aws:` are AWS-managed.
6. Repeat steps 1-5 for each ECS cluster.

### Using AWS CLI

Run the following command to list clusters:

```
aws ecs list-clusters
```

Run the following command to view the tags for a cluster:

```
aws ecs list-tags-for-resource --resource-arn <service-arn>
```

Ensure that tags are returned that do not begin with `aws:`. Tags prefixed with `aws:` are AWS-managed.
Repeat for each cluster.

## Expected Result

All ECS clusters should have at least one user-defined tag (not prefixed with `aws:`).

## Remediation

### Using AWS Console

1. Login to the ECS console using https://console.aws.amazon.com/ecs/.
2. In the left panel, click `Clusters`.
3. Click the name of a cluster.
4. Click `Tags`.
5. Click `Manage tags`.
6. Click `Add tag`.
7. Provide a `Key` and optional `Value` for the tag.
8. Click `Save`.
9. Repeat steps 1-8 for each ECS cluster requiring remediation.

### Using AWS CLI

Use the `aws ecs tag-resource` command to add tags to a cluster:

```
aws ecs tag-resource --resource-arn <cluster-arn> --tags key=<TagKey>,value=<TagValue>
```

## Default Value

By default, Amazon ECS clusters have only AWS-managed tags.

## References

1. https://docs.aws.amazon.com/AmazonECS/latest/developerguide/ecs-using-tags.html
2. https://awscli.amazonaws.com/v2/documentation/api/latest/reference/ecs/list-clusters.html
3. https://awscli.amazonaws.com/v2/documentation/api/latest/reference/ecs/list-tags-for-resource.html

## CIS Controls

| Controls Version | Control                                                        | IG 1 | IG 2 | IG 3 |
| ---------------- | -------------------------------------------------------------- | ---- | ---- | ---- |
| v8               | 1.1 Establish and Maintain Detailed Enterprise Asset Inventory | x    | x    | x    |
| v7               | 1.4 Maintain Detailed Asset Inventory                          | x    | x    | x    |

## Profile

Level 1 | Automated
