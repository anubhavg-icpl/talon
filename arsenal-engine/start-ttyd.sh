#!/bin/sh
# Talon Kali shell (ttyd) — carbon black + electric red xterm theme.
# Loopback-only; SSO enforced by talon-core /shell reverse proxy.

set -eu

PORT="${TTYD_PORT:-7681}"
INDEX="${TTYD_INDEX:-/usr/local/share/ttyd-index.html}"

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

# Fallback without custom index (default gray theme)
exec ttyd \
  -p "${PORT}" \
  -i 127.0.0.1 \
  -b /shell \
  -W \
  -T xterm-256color \
  bash -l
