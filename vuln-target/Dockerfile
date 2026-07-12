# Python mimic of vsftpd 2.3.4 (CVE-2011-2523) for authorized local testing.
# Recon/forge speed path only — Metasploit will NOT get a real shell here.
# Prefer Dockerfile.real (compose default under --profile vuln) for MSF E2E.
#   VULN_DOCKERFILE=Dockerfile VULN_TARGET_IMAGE=talon-vuln-mimic \
#     docker compose --profile vuln up -d --build vuln-target
FROM python:3.12-slim
WORKDIR /srv
COPY server.py .
EXPOSE 21 6200
CMD ["python3", "-u", "server.py"]
