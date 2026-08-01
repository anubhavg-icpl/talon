# stage: report
# category: CIS_benchmarks


# 3.3 Ensure there are no unnecessary roles or excessive privileges

## Profile Applicability

- Level 1 - Cassandra
- Level 1 - Cassandra on Linux

## Description

Verify each role is require and has only the privileges needed to do its job.

## Rationale

Roles which are unneeded, have super user or other potentially excessive privileges may be an avenue for a hacker to gain access to or modify data in the database.

## Audit

As a superuser, retrieve all roles:

```sql
list roles;
```

Retrieve all permissions for all roles

```sql
select * from system_auth.role_permissions;
```

If there are any unnecessary roles or roles with excessive privileges this is a finding.

## Remediation

Remove any unnecessary roles and/or permissions in accordance with organizational needs.

## References

1. http://cassandra.apache.org/doc/latest/cql/security.html

## CIS Controls

**Controls Version v8:**

- 6.8 Define and Maintain Role-Based Access Control
  - Define and maintain role-based access control, through determining and documenting the access rights necessary for each role within the enterprise to successfully carry out its assigned duties. Perform access control reviews of enterprise assets to validate that all privileges are authorized, on a recurring schedule at a minimum annually, or more frequently.

## Profile

- Level 1 | Manual
