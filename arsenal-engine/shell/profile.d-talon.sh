# /etc/profile.d/talon-arsenal.sh — login shells (bash -l via ttyd)
# shellcheck shell=sh
export TALON_ARSENAL=1
export TALON_DISTRO=blackarch
export API_PORT="${API_PORT:-8888}"
export TTYD_PORT="${TTYD_PORT:-7681}"
export PATH="/root/go/bin:/usr/local/go/bin:/usr/local/bin:${PATH:-/usr/bin}"

# Prefer our bashrc for interactive login
if [ -n "${BASH_VERSION:-}" ] && [ -f /etc/talon/bashrc ]; then
  # shellcheck disable=SC1091
  . /etc/talon/bashrc
fi
