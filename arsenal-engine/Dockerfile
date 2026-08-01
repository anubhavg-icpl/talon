# =============================================================================
# Talon Arsenal Engine — BlackArch tool image (not Kali)
# =============================================================================
# Base: Arch Linux + BlackArch strap (blackarch.org)
# Same contracts as before so nothing breaks for Talon:
#   • HTTP API  :${API_PORT:-8888}  → arsenal_engine.py
#   • ttyd PTY  :127.0.0.1:${TTYD_PORT:-7681}/shell  (SSO via talon-core)
#   • health    GET /health
#   • compose   service arsenal-engine, network_mode: host, privileged
#
# Build:
#   docker compose build arsenal-engine
#   docker build --build-arg INSTALL_CLOUD=1 --build-arg INSTALL_BINARY_HEAVY=1 \
#     -t talon-arsenal-engine ./arsenal-engine
# =============================================================================

FROM archlinux:base-devel

LABEL org.opencontainers.image.title="talon-arsenal-engine" \
      org.opencontainers.image.description="BlackArch tool-execution backend for Talon (HexStrike-class arsenal + MCP HTTP API)" \
      org.opencontainers.image.source="https://github.com/anubhavg-icpl/talon"

ARG INSTALL_CLOUD=1
ARG INSTALL_BINARY_HEAVY=1
ARG INSTALL_BROWSER=1
ARG INSTALL_WIRELESS=0
ARG RUSTSCAN_VERSION=2.3.0
ARG TRIVY_VERSION=0.58.1
ARG TERRASCAN_VERSION=1.19.9
ARG KUBE_BENCH_VERSION=0.9.2
ARG TTYD_VERSION=1.7.7
ARG GO_VERSION=1.23.6

ENV GOPATH=/root/go \
    GOBIN=/root/go/bin \
    PATH="/root/go/bin:/usr/local/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin" \
    LANG=C.UTF-8 \
    PYTHONDONTWRITEBYTECODE=1 \
    PYTHONUNBUFFERED=1

# ---------------------------------------------------------------------------
# 0) Arch base + BlackArch strap (official)
# ---------------------------------------------------------------------------
RUN set -eux; \
    pacman -Syu --noconfirm --needed \
      base base-devel git curl wget unzip tar xz jq which \
      python python-pip python-setuptools python-wheel \
      gcc clang make cmake pkgconf \
      openssl libpcap libffi libxml2 libxslt zlib libjpeg-turbo libpng \
      ruby \
      jre-openjdk-headless \
      sudo iproute2 iputils net-tools bind-tools gnu-netcat \
      procps-ng file less vim nano bash \
      ca-certificates \
      fastfetch; \
    # BlackArch keyring + mirrorlist
    curl -fsSL -o /tmp/strap.sh https://blackarch.org/strap.sh; \
    chmod +x /tmp/strap.sh; \
    /tmp/strap.sh; \
    rm -f /tmp/strap.sh; \
    pacman -Syu --noconfirm; \
    # Parallel downloads, no confirm prompts
    sed -i 's/^#ParallelDownloads.*/ParallelDownloads = 5/' /etc/pacman.conf || true; \
    pacman -Scc --noconfirm || true

# Official Go (newer than some Arch mirrors for reproducible go install)
RUN set -eux; \
    arch="$(uname -m)"; \
    case "$arch" in \
      x86_64) goarch=amd64 ;; \
      aarch64|arm64) goarch=arm64 ;; \
      *) goarch=amd64 ;; \
    esac; \
    curl -fsSL "https://go.dev/dl/go${GO_VERSION}.linux-${goarch}.tar.gz" \
      | tar -C /usr/local -xz; \
    go version

# ---------------------------------------------------------------------------
# 1) pacman helper — install each package independently (BlackArch renames often)
# ---------------------------------------------------------------------------
RUN printf '%s\n' \
  '#!/bin/sh' \
  'set -eu' \
  'ok=0; fail=0' \
  'for p in "$@"; do' \
  '  if pacman -S --noconfirm --needed "$p"; then ok=$((ok+1));' \
  '  else echo "WARN: pacman package unavailable: $p"; fail=$((fail+1)); fi' \
  'done' \
  'echo "pac-opt: ok=$ok fail=$fail"' \
  'exit 0' \
  > /usr/local/bin/pac-opt && chmod +x /usr/local/bin/pac-opt

# ---------------------------------------------------------------------------
# 2) Core BlackArch / Arch security tools (grouped)
# ---------------------------------------------------------------------------

# 2a Network & recon
RUN set -eux; \
    pacman -Sy --noconfirm; \
    /usr/local/bin/pac-opt \
      nmap masscan amass subfinder fierce dnsenum theharvester \
      netexec crackmapexec enum4linux enum4linux-ng smbmap nbtscan arp-scan \
      responder autorecon rpcbind samba \
      whatweb sslscan testssl.sh rustscan; \
    pacman -Scc --noconfirm || true

# 2b Web application
RUN set -eux; \
    pacman -Sy --noconfirm; \
    /usr/local/bin/pac-opt \
      gobuster feroxbuster ffuf dirb dirsearch nuclei nikto sqlmap wpscan \
      wafw00f arjun wfuzz commix; \
    pacman -Scc --noconfirm || true

# 2c Password & auth
RUN set -eux; \
    pacman -Sy --noconfirm; \
    /usr/local/bin/pac-opt \
      hydra john hashcat medusa patator hash-identifier \
      seclists wordlists; \
    pacman -Scc --noconfirm || true

# 2d Binary / RE / forensics
RUN set -eux; \
    pacman -Sy --noconfirm; \
    /usr/local/bin/pac-opt \
      gdb radare2 binwalk checksec foremost steghide perl-image-exiftool \
      binutils strace ltrace \
      testdisk sleuthkit scalpel bulk_extractor \
      outguess; \
    pacman -Scc --noconfirm || true

# 2e Wireless (optional)
RUN set -eux; \
    if [ "${INSTALL_WIRELESS}" = "1" ]; then \
      pacman -Sy --noconfirm; \
      /usr/local/bin/pac-opt aircrack-ng reaver bully kismet wireshark-cli tcpdump; \
      pacman -Scc --noconfirm || true; \
    else \
      pacman -Sy --noconfirm; \
      /usr/local/bin/pac-opt tcpdump wireshark-cli; \
      pacman -Scc --noconfirm || true; \
    fi

# 2f Exploit-DB searchsploit
RUN set -eux; \
    pacman -Sy --noconfirm; \
    /usr/local/bin/pac-opt exploitdb; \
    pacman -Scc --noconfirm || true

# ---------------------------------------------------------------------------
# 3) Browser agent (Chromium)
# ---------------------------------------------------------------------------
RUN set -eux; \
    if [ "${INSTALL_BROWSER}" = "1" ]; then \
      pacman -Sy --noconfirm; \
      /usr/local/bin/pac-opt chromium chromedriver ttf-liberation ttf-dejavu; \
      if [ -x /usr/bin/chromium ] && [ ! -e /usr/bin/google-chrome ]; then \
        ln -sf /usr/bin/chromium /usr/bin/google-chrome; \
      fi; \
      if [ -x /usr/bin/chromium ] && [ ! -e /usr/bin/chromium-browser ]; then \
        ln -sf /usr/bin/chromium /usr/bin/chromium-browser; \
      fi; \
      pacman -Scc --noconfirm || true; \
    fi

# ---------------------------------------------------------------------------
# 4) Go-based recon / fuzzers (latest upstream; complements BlackArch)
# ---------------------------------------------------------------------------
ENV GOPROXY=https://proxy.golang.org,direct
RUN set -eu; \
    failed=0; \
    for pkg in \
      github.com/projectdiscovery/httpx/cmd/httpx@latest \
      github.com/projectdiscovery/katana/cmd/katana@latest \
      github.com/projectdiscovery/nuclei/v3/cmd/nuclei@latest \
      github.com/projectdiscovery/subfinder/v2/cmd/subfinder@latest \
      github.com/hahwul/dalfox/v2@latest \
      github.com/lc/gau/v2/cmd/gau@latest \
      github.com/tomnomnom/waybackurls@latest \
      github.com/hakluke/hakrawler@latest \
      github.com/tomnomnom/anew@latest \
      github.com/tomnomnom/qsreplace@latest \
      github.com/ffuf/ffuf/v2@latest \
      github.com/OJ/gobuster/v3@latest \
      github.com/jaeles-project/jaeles@latest \
    ; do \
      echo "=== go install $pkg ==="; \
      if ! go install -v "$pkg"; then \
        echo "WARN: go install failed for $pkg"; \
        failed=$((failed + 1)); \
      fi; \
    done; \
    echo "go install complete (failures=$failed)"; \
    if command -v nuclei >/dev/null 2>&1; then \
      nuclei -update-templates -silent || true; \
    fi

# ---------------------------------------------------------------------------
# 5) Release binaries (ttyd, trivy, terrascan, kube-bench; rustscan if missing)
# ---------------------------------------------------------------------------
RUN set -eux; \
    arch="$(uname -m)"; \
    case "$arch" in aarch64|arm64) ra=aarch64; goa=arm64; ta=ARM64 ;; *) ra=x86_64; goa=amd64; ta=64bit ;; esac; \
    \
    # rustscan fallback if not from BlackArch
    if ! command -v rustscan >/dev/null 2>&1 && [ "$goa" = "amd64" ]; then \
      curl -fsSL -o /tmp/rustscan.deb \
        "https://github.com/RustScan/RustScan/releases/download/${RUSTSCAN_VERSION}/rustscan_${RUSTSCAN_VERSION}_amd64.deb" || true; \
      if [ -f /tmp/rustscan.deb ]; then \
        # extract deb without dpkg on Arch
        cd /tmp && ar x rustscan.deb && tar -xf data.tar.* -C / && \
        command -v rustscan || true; \
      fi; \
      rm -rf /tmp/rustscan.deb /tmp/data.tar.* /tmp/control.tar.* /tmp/debian-binary 2>/dev/null || true; \
    fi; \
    \
    # ttyd
    curl -fsSL -o /usr/local/bin/ttyd \
      "https://github.com/tsl0922/ttyd/releases/download/${TTYD_VERSION}/ttyd.${ra}" \
      && chmod +x /usr/local/bin/ttyd \
      && /usr/local/bin/ttyd --version \
      || echo "WARN: ttyd install skipped"

RUN set -eux; \
    if [ "${INSTALL_CLOUD}" != "1" ]; then echo "INSTALL_CLOUD=0 — skip cloud binaries"; exit 0; fi; \
    arch="$(uname -m)"; \
    case "$arch" in aarch64|arm64) goa=arm64; ta=ARM64 ;; *) goa=amd64; ta=64bit ;; esac; \
    \
    curl -fsSL -o /tmp/trivy.tar.gz \
      "https://github.com/aquasecurity/trivy/releases/download/v${TRIVY_VERSION}/trivy_${TRIVY_VERSION}_Linux-${ta}.tar.gz" \
      && tar -xzf /tmp/trivy.tar.gz -C /usr/local/bin trivy \
      && chmod +x /usr/local/bin/trivy \
      || true; \
    rm -f /tmp/trivy.tar.gz; \
    \
    curl -fsSL -o /tmp/terrascan.tar.gz \
      "https://github.com/tenable/terrascan/releases/download/v${TERRASCAN_VERSION}/terrascan_${TERRASCAN_VERSION}_Linux_x86_64.tar.gz" \
      && tar -xzf /tmp/terrascan.tar.gz -C /usr/local/bin terrascan \
      && chmod +x /usr/local/bin/terrascan \
      || true; \
    rm -f /tmp/terrascan.tar.gz; \
    \
    curl -fsSL -o /tmp/kube-bench.tar.gz \
      "https://github.com/aquasecurity/kube-bench/releases/download/v${KUBE_BENCH_VERSION}/kube-bench_${KUBE_BENCH_VERSION}_linux_${goa}.tar.gz" \
      && tar -xzf /tmp/kube-bench.tar.gz -C /usr/local/bin kube-bench \
      && chmod +x /usr/local/bin/kube-bench \
      || true; \
    rm -f /tmp/kube-bench.tar.gz; \
    \
    git clone --depth 1 https://github.com/docker/docker-bench-security.git /opt/docker-bench-security \
      && printf '%s\n' '#!/bin/sh' 'cd /opt/docker-bench-security && exec sh docker-bench-security.sh "$@"' \
         > /usr/local/bin/docker-bench-security \
      && chmod +x /usr/local/bin/docker-bench-security \
      || true

# ---------------------------------------------------------------------------
# 6) Ruby gems
# ---------------------------------------------------------------------------
RUN set -eux; \
    gem install --no-document evil-winrm zsteg || \
      gem install --no-document evil-winrm || true

# ---------------------------------------------------------------------------
# 7) Python API + security libraries
# ---------------------------------------------------------------------------
WORKDIR /opt/arsenal-engine
COPY requirements.txt /opt/arsenal-engine/requirements.txt

# Arch: system packages may conflict; --break-system-packages is supported on pip 23+
RUN pip3 install --break-system-packages --no-cache-dir --ignore-installed \
      -r /opt/arsenal-engine/requirements.txt \
    || pip install --break-system-packages --no-cache-dir --ignore-installed \
      -r /opt/arsenal-engine/requirements.txt

RUN set -eux; \
    PIP="pip3"; command -v pip3 >/dev/null || PIP="pip"; \
    $PIP install --break-system-packages --no-cache-dir --ignore-installed \
      httpie uro paramspider checkov kube-hunter shodan censys \
      || true; \
    if [ "${INSTALL_CLOUD}" = "1" ]; then \
      $PIP install --break-system-packages --no-cache-dir --ignore-installed \
        prowler scoutsuite pacu \
        || true; \
    fi; \
    $PIP install --break-system-packages --no-cache-dir --ignore-installed \
      volatility3 \
      || true

RUN set -eux; \
    if [ "${INSTALL_BINARY_HEAVY}" = "1" ]; then \
      PIP="pip3"; command -v pip3 >/dev/null || PIP="pip"; \
      $PIP install --break-system-packages --no-cache-dir --ignore-installed \
        'pwntools>=4.10.0,<5.0.0' \
        'angr>=9.2.0,<10.0.0' \
        ropgadget one-gadget \
        || true; \
      $PIP install --break-system-packages --no-cache-dir --ignore-installed ropper || true; \
    fi

# Aliases arsenal_engine health probes expect
RUN set -eux; \
    if command -v nxc >/dev/null 2>&1 && [ ! -e /usr/local/bin/netexec ]; then \
      ln -sf "$(command -v nxc)" /usr/local/bin/netexec; \
    fi; \
    if command -v netexec >/dev/null 2>&1 && [ ! -e /usr/local/bin/nxc ]; then \
      ln -sf "$(command -v netexec)" /usr/local/bin/nxc; \
    fi; \
    if command -v crackmapexec >/dev/null 2>&1; then \
      [ -e /usr/local/bin/netexec ] || ln -sf "$(command -v crackmapexec)" /usr/local/bin/netexec; \
      [ -e /usr/local/bin/nxc ] || ln -sf "$(command -v crackmapexec)" /usr/local/bin/nxc; \
    fi; \
    if command -v volatility3 >/dev/null 2>&1 && [ ! -e /usr/local/bin/vol ]; then \
      ln -sf "$(command -v volatility3)" /usr/local/bin/vol; \
    fi; \
    # exiftool path
    if command -v exiftool >/dev/null 2>&1; then true; \
    elif command -v perl-image-exiftool >/dev/null 2>&1; then true; \
    elif [ -x /usr/bin/vendor_perl/exiftool ]; then \
      ln -sf /usr/bin/vendor_perl/exiftool /usr/local/bin/exiftool; \
    fi; \
    true

# ---------------------------------------------------------------------------
# 8) Talon BlackArch operator shell (fastfetch + MOTD + red prompt)
# ---------------------------------------------------------------------------
COPY shell/fastfetch.jsonc /etc/talon/fastfetch.jsonc
COPY shell/talon-bashrc.sh /etc/talon/bashrc
COPY shell/motd /etc/talon/motd
COPY shell/profile.d-talon.sh /etc/profile.d/talon-arsenal.sh
RUN set -eux; \
    chmod 644 /etc/talon/fastfetch.jsonc /etc/talon/bashrc /etc/talon/motd /etc/profile.d/talon-arsenal.sh; \
    # root interactive defaults
    printf '%s\n' \
      '# Talon BlackArch arsenal' \
      '[[ -f /etc/talon/bashrc ]] && . /etc/talon/bashrc' \
      > /root/.bashrc; \
    # ensure login shells pick up profile.d
    if ! grep -q 'profile.d' /root/.bash_profile 2>/dev/null; then \
      printf '%s\n' \
        '[[ -f /etc/profile ]] && . /etc/profile' \
        '[[ -f ~/.bashrc ]] && . ~/.bashrc' \
        > /root/.bash_profile; \
    fi; \
    # hostname for prompt / fastfetch
    echo 'talon-blackarch' > /etc/hostname || true; \
    # fallback if fastfetch package name differs
    if ! command -v fastfetch >/dev/null 2>&1; then \
      pacman -Sy --noconfirm --needed fastfetch || \
      pacman -Sy --noconfirm --needed neofetch || true; \
    fi

# ---------------------------------------------------------------------------
# 9) Application + ttyd theme (API + shell contracts unchanged)
# ---------------------------------------------------------------------------
COPY arsenal_engine.py /opt/arsenal-engine/arsenal_engine.py
COPY LICENSE.upstream /opt/arsenal-engine/LICENSE.upstream
COPY start-ttyd.sh /usr/local/bin/start-ttyd.sh
COPY ttyd-index.html /usr/local/share/ttyd-index.html
RUN chmod +x /usr/local/bin/start-ttyd.sh

RUN set -eu; \
    printf '%s\n' \
      '#!/bin/sh' \
      'tools="nmap masscan rustscan amass subfinder nuclei gobuster feroxbuster ffuf dirb dirsearch nikto sqlmap wpscan httpx katana dalfox wafw00f hydra john hashcat medusa netexec nxc enum4linux-ng responder smbmap gdb radare2 binwalk trivy prowler checkov kube-hunter curl python3 ttyd chromium chromium-browser google-chrome fastfetch"' \
      'ok=0; miss=0' \
      'for t in $tools; do' \
      '  if command -v "$t" >/dev/null 2>&1; then echo "OK  $t"; ok=$((ok+1));' \
      '  else echo "MISS $t"; miss=$((miss+1)); fi' \
      'done' \
      'echo "summary ok=$ok miss=$miss distro=blackarch talon-shell=1"' \
      > /usr/local/bin/arsenal-tool-check && chmod +x /usr/local/bin/arsenal-tool-check; \
    /usr/local/bin/arsenal-tool-check || true

ENV API_PORT=8888 \
    TTYD_PORT=7681 \
    CHROME_BIN=/usr/bin/chromium \
    CHROMIUM_FLAGS="--no-sandbox --disable-dev-shm-usage --disable-gpu --headless=new" \
    TALON_ARSENAL=1 \
    TALON_DISTRO=blackarch

WORKDIR /opt/arsenal-engine
EXPOSE 8888 7681

# Still not in image: Ghidra, Burp, IDA, full MSF (use msf_rpc), Maltego

HEALTHCHECK --interval=30s --timeout=5s --start-period=60s --retries=3 \
  CMD curl -fsS "http://127.0.0.1:${API_PORT}/health" >/dev/null || exit 1

CMD sh -c '/usr/local/bin/start-ttyd.sh & exec python3 arsenal_engine.py --port "${API_PORT}"'
