# stage: report
# category: CIS_benchmarks


# Ensure the Lock File Is Secured

## Description

The `LockFile` directive sets the path to the lock file used when Apache uses `fcntl(2)` or `flock(2)` system calls to implement locking. Most Linux systems will default to using semaphores instead, so the directive may not apply. However, in the event a lock file is used, it is important for the lock file to be in a locally mounted directory that is not writable by other users.

## Rationale

If the `LockFile` is placed in a writable directory, other accounts could create a denial of service attack and prevent the server from starting by creating a lock file with the same name.

## Impact

None documented

## Audit Procedure

Perform these steps to verify the lock file is secured properly:

1. Find the directory in which the `LockFile` would be created. The default value is the `ServerRoot/logs` directory.
2. Verify that the lock file directory is not a directory within the Apache `DocumentRoot`.
3. Verify that the lock file directory is on a locally mounted hard drive rather than an NFS mounted file system.
4. Verify that the ownership and group of the directory is `root:root` (or the user under which apache initially starts up if not root).
5. Verify that the permissions on the directory are only writable by root (or the startup user if not root).

## Remediation

Perform these steps to properly secure the lock file:

1. Find the directory in which the `LockFile` would be created. The default value is the `ServerRoot/logs` directory.
2. Modify the directory for the `LockFile` so it is not within the Apache `DocumentRoot` and so it is on a locally mounted hard drive rather than an NFS mounted file system.
3. Change the ownership and group of the directory to be `root:root`.
4. Change the permissions on the directory so it is only writable by root, or the user under which apache initially starts up (default is root).

## Default Value

The default lock file is `logs/accept.lock`.

## References

1. https://httpd.apache.org/docs/2.2/mod/mpm_common.html#lockfile

## CIS Controls

### Version 6

18 Application Software Security

Application Software Security

### Version 7

14.6 Protect Information through Access Control Lists

Protect all information stored on systems with file system, network share, claims, application, or database specific access control lists. These controls will enforce the principle that only authorized individuals should have access to the information based on their need to access the information as a part of their responsibilities.

## Profile

Level 1 | Scored
