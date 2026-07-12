#!/bin/bash
# Restart vsftpd after each backdoor trigger (the malware never returns).
set -e
mkdir -p /usr/share/empty /var/ftp
echo "vuln-target(real): infected vsftpd 2.3.4 on :21; backdoor shell on :6200 after USER *:) + PASS"
while true; do
  /usr/local/sbin/vsftpd /etc/vsftpd.conf || true
  sleep 0.5
done
