# Lab vulnerable target for authorized local testing only.
# Single file, two build targets (compose default: real):
#
#   docker compose --profile vuln build vuln-target          # target=real
#   docker compose --profile vuln build --build-arg ...      # see compose
#
# Targets:
#   real  — infected vsftpd 2.3.4 (CVE-2011-2523); Metasploit can open a shell
#           with payload cmd/unix/reverse_bash + LHOST/LPORT
#   mimic — Python banner/backdoor mimic; recon/forge only (no real MSF session)
#
# Trigger (real): USER containing ":)" then PASS → shell on :6200

# ---------------------------------------------------------------------------
# mimic — fast Python stand-in
# ---------------------------------------------------------------------------
FROM python:3.12-slim AS mimic
WORKDIR /srv
COPY server.py .
EXPOSE 21 6200
CMD ["python3", "-u", "server.py"]

# ---------------------------------------------------------------------------
# real — compile nikdubois infected vsftpd 2.3.4
# ---------------------------------------------------------------------------
FROM debian:bullseye-slim AS real
ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update && apt-get install -y --no-install-recommends \
      build-essential ca-certificates git libssl-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /src
RUN git clone --depth 1 https://github.com/nikdubois/vsftpd-2.3.4-infected.git . \
    && chmod +x vsf_findlibs.sh \
    && make LDFLAGS="-Wl,-s -lcrypt" \
    && install -m 755 vsftpd /usr/local/sbin/vsftpd \
    && mkdir -p /usr/share/empty /var/ftp \
    && (useradd -d /var/ftp -s /bin/false ftp || true)

# local_enable=YES so non-anonymous USER names (e.g. letmein:)) reach the
# backdoor check (anonymous-only mode rejects them first).
RUN printf '%s\n' \
      'listen=YES' \
      'anonymous_enable=YES' \
      'local_enable=YES' \
      'write_enable=NO' \
      'anon_root=/var/ftp' \
      'secure_chroot_dir=/usr/share/empty' \
      'xferlog_enable=NO' \
      'connect_from_port_20=YES' \
      > /etc/vsftpd.conf

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh
EXPOSE 21 6200
CMD ["/entrypoint.sh"]
