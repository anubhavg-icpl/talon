# Local vulnerable target for authorized Talon pipeline testing only.
# Mimics vsftpd 2.3.4 backdoor (CVE-2011-2523). Profile-gated ("vuln") so it
# is never built/shipped in a normal deploy.
FROM python:3.12-slim
WORKDIR /srv
COPY server.py .
EXPOSE 21 6200
CMD ["python3", "-u", "server.py"]
