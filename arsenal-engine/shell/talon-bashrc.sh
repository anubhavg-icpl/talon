# shellcheck shell=bash
# =============================================================================
# Talon BlackArch Arsenal — interactive shell (sourced for login + interactive)
# Carbon black + electric red operator theme
# =============================================================================

# Only for interactive shells
[[ $- != *i* ]] && return

export TERM="${TERM:-xterm-256color}"
export COLORTERM="${COLORTERM:-truecolor}"
export EDITOR="${EDITOR:-vim}"
export HISTCONTROL=ignoreboth:erasedups
export HISTSIZE=50000
export HISTFILESIZE=100000
export HISTTIMEFORMAT='%F %T  '
shopt -s histappend checkwinsize cmdhist 2>/dev/null || true

# Paths
export PATH="/root/go/bin:/usr/local/go/bin:/usr/local/bin:${PATH}"
export GOPATH="${GOPATH:-/root/go}"

# Colors
RED=$'\e[38;5;196m'
ROSE=$'\e[38;5;203m'
DIM=$'\e[38;5;240m'
WHT=$'\e[38;5;255m'
RST=$'\e[0m'
BLD=$'\e[1m'

# Prompt: root@talon-blackarch  ~/path
#          ❯
_talon_git_branch() {
  git rev-parse --abbrev-ref HEAD 2>/dev/null | sed 's/.*/  &/'
}

PROMPT_COMMAND='history -a'
PS1="\[${DIM}\]┌─\[${RST}\]\[${RED}${BLD}\]talon\[${RST}\]\[${DIM}\]@\[${RST}\]\[${WHT}\]blackarch\[${RST}\] \[${ROSE}\]\w\[${DIM}\]\$(_talon_git_branch)\[${RST}\]\n\[${DIM}\]└─\[${RST}\]\[${RED}\]❯\[${RST}\] "

# Aliases — operator speed
alias ll='ls -lah --color=auto'
alias la='ls -A --color=auto'
alias ls='ls --color=auto'
alias grep='grep --color=auto'
alias ff='fastfetch --config /etc/talon/fastfetch.jsonc'
alias tools='arsenal-tool-check'
alias health='curl -fsS http://127.0.0.1:${API_PORT:-8888}/health | head -c 400; echo'
alias ports='ss -ltnp 2>/dev/null || netstat -ltnp 2>/dev/null'
alias ..='cd ..'
alias ...='cd ../..'

# Colorized man if possible
export LESS_TERMCAP_mb=$'\e[1;31m'
export LESS_TERMCAP_md=$'\e[1;31m'
export LESS_TERMCAP_me=$'\e[0m'
export LESS_TERMCAP_se=$'\e[0m'
export LESS_TERMCAP_so=$'\e[38;5;196;48;5;236m'
export LESS_TERMCAP_ue=$'\e[0m'
export LESS_TERMCAP_us=$'\e[1;37m'

# Banner + fastfetch once per tty login (not every subshell)
if [[ -z "${TALON_SHELL_BANNER_SHOWN:-}" ]]; then
  export TALON_SHELL_BANNER_SHOWN=1
  if [[ -r /etc/talon/motd ]]; then
    cat /etc/talon/motd
  fi
  if command -v fastfetch >/dev/null 2>&1; then
    echo
    fastfetch --config /etc/talon/fastfetch.jsonc 2>/dev/null || fastfetch 2>/dev/null || true
  fi
  echo
  printf '%s\n' "${DIM}API ${ROSE}http://127.0.0.1:${API_PORT:-8888}${DIM}  ·  tools ${ROSE}arsenal-tool-check${DIM}  ·  system ${ROSE}ff${RST}"
  echo
fi
