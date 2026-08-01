# stage: report
# category: CIS_benchmarks


# CIS Docker Benchmark v1.7.0 - Control 7.4

## Profile Applicability

- **Level:** 1
- **Type:** Manual
- **Platform:** Docker Swarm

## Description

You should use Docker's in-built secret management command for control of secrets.

## Rationale

Docker has various commands for managing secrets in a swarm cluster.

## Impact

None

## Audit Procedure

On a swarm manager node, you should run the command below and ensure `docker secret` management is used in your environment where this is in line with your IT security policy.

```bash
docker secret ls
```

## Remediation

You should follow the `docker secret` documentation and use it to manage secrets effectively.

## Default Value

Not Applicable

## References

1. https://docs.docker.com/engine/reference/commandline/secret/

## CIS Controls

**v8:**

- **4 Secure Configuration of Enterprise Assets and Software**
  - Establish and maintain the secure configuration of enterprise assets (end-user devices, including portable and mobile; network devices; non-computing/IoT devices; and servers) and software (operating systems and applications).

**v7:**

- **5.1 Establish Secure Configurations**
  - Maintain documented, standard security configuration standards for all authorized operating systems and software.
