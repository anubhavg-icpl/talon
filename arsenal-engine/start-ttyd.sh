#!/bin/sh
# Talon BlackArch arsenal shell (ttyd) — carbon + electric red + fastfetch MOTD.
# Loopback-only; SSO enforced by talon-core /shell reverse proxy.

set -eu

PORT="${TTYD_PORT:-7681}"
INDEX="${TTYD_INDEX:-/usr/local/share/ttyd-index.html}"

# bash -l loads /etc/profile → profile.d/talon-arsenal.sh → MOTD + fastfetch + PS1
if [ -f "$INDEX" ]; then
  exec ttyd \
    -p "${PORT}" \
    -i 127.0.0.1 \
    -b /shell \
    -W \
    -T xterm-256color \
    -I "${INDEX}" \
    bash -l
fi

exec ttyd \
  -p "${PORT}" \
  -i 127.0.0.1 \
  -b /shell \
  -W \
  -T xterm-256color \
  bash -l
