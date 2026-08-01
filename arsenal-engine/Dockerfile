# =============================================================================
# Talon Arsenal Engine — production Kali tool image
# =============================================================================
# Base: kalilinux/kali-rolling (preferred over BlackArch for Docker)
#   • Stable, well-maintained packages matching arsenal_engine.py names
#   • Smaller/more reproducible than blackarch/blackarch
#   • Same role as kali-msf/ for Metasploit — Python engine + CLI arsenal
#
# Build:
#   docker compose build arsenal-engine
#   # or full extras (cloud + heavy binary Python deps):
#   docker build --build-arg INSTALL_CLOUD=1 --build-arg INSTALL_BINARY_HEAVY=1 \
#     -t talon-arsenal-engine ./arsenal-engine
#
# Health: curl -s http://127.0.0.1:8888/health
# Shell:  ttyd on 127.0.0.1:7681 (SSO via talon-core /shell)
# =============================================================================

FROM kalilinux/kali-rolling

LABEL org.opencontainers.image.title="talon-arsenal-engine" \
      org.opencontainers.image.description="Kali tool-execution backend for Talon (HexStrike-class arsenal + MCP HTTP API)" \
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

ENV DEBIAN_FRONTEND=noninteractive \
    GOPATH=/root/go \
    GOBIN=/root/go/bin \
    PATH="/root/go/bin:/usr/local/go/bin:/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin" \
    LANG=C.UTF-8 \
    PYTHONDONTWRITEBYTECODE=1 \
    PYTHONUNBUFFERED=1

# ---------------------------------------------------------------------------
# 0) Base OS + build toolchain + shell utilities (must succeed)
# ---------------------------------------------------------------------------
RUN set -eux; \
    apt-get update; \
    apt-get install -y --no-install-recommends \
      ca-certificates gnupg apt-transport-https \
      python3 python3-pip python3-venv python3-dev python3-setuptools python3-wheel \
      build-essential gcc g++ make cmake pkg-config \
      git curl wget unzip tar xz-utils jq file procps \
      iproute2 iputils-ping net-tools dnsutils netcat-openbsd \
      libpcap-dev libssl-dev libffi-dev libxml2-dev libxslt1-dev \
      zlib1g-dev libjpeg-dev libpng-dev \
      ruby ruby-dev \
      default-jre-headless \
      sudo less vim-tiny \
      pciutils usbutils \
    ; \
    rm -rf /var/lib/apt/lists/*

# Official Go (Kali's golang-go lags; go install needs a recent toolchain)
RUN set -eux; \
    arch="$(dpkg --print-architecture)"; \
    case "$arch" in \
      amd64) goarch=amd64 ;; \
      arm64) goarch=arm64 ;; \
      *) goarch=amd64 ;; \
    esac; \
    curl -fsSL "https://go.dev/dl/go${GO_VERSION}.linux-${goarch}.tar.gz" \
      | tar -C /usr/local -xz; \
    go version

# ---------------------------------------------------------------------------
# 1) Helper: install apt packages without aborting the layer on missing names
#    (Kali rolling renames packages occasionally)
# ---------------------------------------------------------------------------
RUN printf '%s\n' \
  '#!/bin/sh' \
  'set -eu' \
  'ok=0; fail=0' \
  'for p in "$@"; do' \
  '  if apt-get install -y --no-install-recommends "$p"; then ok=$((ok+1));' \
  '  else echo "WARN: apt package unavailable: $p"; fail=$((fail+1)); fi' \
  'done' \
  'echo "apt-opt: ok=$ok fail=$fail"' \
  'exit 0' \
  > /usr/local/bin/apt-opt && chmod +x /usr/local/bin/apt-opt

# ---------------------------------------------------------------------------
# 2) Core security tools (network, web, password, binary, forensics, wireless)
#    Split into groups so a flaky mirror can be retried per layer.
# ---------------------------------------------------------------------------

# 2a Network & recon
RUN set -eux; \
    apt-get update; \
    /usr/local/bin/apt-opt \
      nmap masscan amass subfinder fierce dnsenum theharvester \
      netexec enum4linux-ng enum4linux smbmap nbtscan arp-scan responder \
      autorecon rpcbind samba-common-bin \
      whatweb sslscan testssl.sh; \
    rm -rf /var/lib/apt/lists/*

# 2b Web application
RUN set -eux; \
    apt-get update; \
    /usr/local/bin/apt-opt \
      gobuster feroxbuster ffuf dirb dirsearch nuclei nikto sqlmap wpscan \
      wafw00f arjun wfuzz commix nosqlmap \
      zaproxy; \
    rm -rf /var/lib/apt/lists/*

# 2c Password & auth
RUN set -eux; \
    apt-get update; \
    /usr/local/bin/apt-opt \
      hydra john hashcat medusa patator hash-identifier hashcat-utils \
      ophcrack seclists wordlists; \
    rm -rf /var/lib/apt/lists/*

# 2d Binary / RE / forensics
RUN set -eux; \
    apt-get update; \
    /usr/local/bin/apt-opt \
      gdb gdb-peda radare2 binwalk checksec foremost steghide exiftool \
      binutils strace ltrace \
      testdisk sleuthkit scalpel bulk-extractor autopsy \
      stegsnow outguess; \
    rm -rf /var/lib/apt/lists/*

# 2e Wireless (optional — large; off by default)
RUN set -eux; \
    if [ "${INSTALL_WIRELESS}" = "1" ]; then \
      apt-get update; \
      /usr/local/bin/apt-opt aircrack-ng reaver bully kismet wireshark-common tshark tcpdump; \
      rm -rf /var/lib/apt/lists/*; \
    else \
      apt-get update; \
      /usr/local/bin/apt-opt tshark tcpdump; \
      rm -rf /var/lib/apt/lists/*; \
    fi

# 2f Metasploit client utils (full RPC is separate msf_rpc service)
RUN set -eux; \
    apt-get update; \
    /usr/local/bin/apt-opt exploitdb; \
    rm -rf /var/lib/apt/lists/*; \
    command -v searchsploit >/dev/null 2>&1 || true

# ---------------------------------------------------------------------------
# 3) Browser agent (Chromium + chromedriver for Selenium)
# ---------------------------------------------------------------------------
RUN set -eux; \
    if [ "${INSTALL_BROWSER}" = "1" ]; then \
      apt-get update; \
      /usr/local/bin/apt-opt \
        chromium chromium-driver \
        fonts-liberation fonts-dejavu-core; \
      # Symlinks Selenium / arsenal often expect
      if [ -x /usr/bin/chromium ] && [ ! -e /usr/bin/google-chrome ]; then \
        ln -sf /usr/bin/chromium /usr/bin/google-chrome; \
      fi; \
      if [ -x /usr/bin/chromium ] && [ ! -e /usr/bin/chromium-browser ]; then \
        ln -sf /usr/bin/chromium /usr/bin/chromium-browser; \
      fi; \
      rm -rf /var/lib/apt/lists/*; \
    fi

# ---------------------------------------------------------------------------
# 4) Go-based recon / fuzzers (not always in Kali apt, or newer upstream)
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
      github.com/dwisiswant0/crlfuzz/cmd/crlfuzz@latest \
    ; do \
      echo "=== go install $pkg ==="; \
      if ! go install -v "$pkg"; then \
        echo "WARN: go install failed for $pkg"; \
        failed=$((failed + 1)); \
      fi; \
    done; \
    echo "go install complete (failures=$failed)"; \
    # Ensure nuclei templates available for first scan
    if command -v nuclei >/dev/null 2>&1; then \
      nuclei -update-templates -silent || true; \
    fi

# ---------------------------------------------------------------------------
# 5) Release binaries (rustscan, trivy, terrascan, kube-bench, ttyd)
# ---------------------------------------------------------------------------
RUN set -eux; \
    arch="$(dpkg --print-architecture)"; \
    case "$arch" in arm64) ra=aarch64 ;; *) ra=x86_64 ;; esac; \
    \
    # rustscan (amd64 deb; skip quietly on other arches)
    if [ "$arch" = "amd64" ]; then \
      curl -fsSL -o /tmp/rustscan.deb \
        "https://github.com/RustScan/RustScan/releases/download/${RUSTSCAN_VERSION}/rustscan_${RUSTSCAN_VERSION}_amd64.deb" \
        || curl -fsSL -o /tmp/rustscan.deb \
        "https://github.com/bee-san/RustScan/releases/download/${RUSTSCAN_VERSION}/rustscan_${RUSTSCAN_VERSION}_amd64.deb" \
        || true; \
      if [ -f /tmp/rustscan.deb ]; then dpkg -i /tmp/rustscan.deb || apt-get -f install -y || true; fi; \
      rm -f /tmp/rustscan.deb; \
    fi; \
    \
    # ttyd
    curl -fsSL -o /usr/local/bin/ttyd \
      "https://github.com/tsl0922/ttyd/releases/download/${TTYD_VERSION}/ttyd.${ra}" \
      && chmod +x /usr/local/bin/ttyd \
      && /usr/local/bin/ttyd --version \
      || echo "WARN: ttyd install skipped"

# Cloud / container scanners (optional but default ON)
RUN set -eux; \
    if [ "${INSTALL_CLOUD}" != "1" ]; then echo "INSTALL_CLOUD=0 — skipping cloud binaries"; exit 0; fi; \
    arch="$(dpkg --print-architecture)"; \
    case "$arch" in arm64) goa=arm64; ta=ARM64 ;; *) goa=amd64; ta=64bit ;; esac; \
    \
    # trivy
    curl -fsSL -o /tmp/trivy.deb \
      "https://github.com/aquasecurity/trivy/releases/download/v${TRIVY_VERSION}/trivy_${TRIVY_VERSION}_Linux-${ta}.deb" \
      && dpkg -i /tmp/trivy.deb || true; \
    rm -f /tmp/trivy.deb; \
    \
    # terrascan
    curl -fsSL -o /tmp/terrascan.tar.gz \
      "https://github.com/tenable/terrascan/releases/download/v${TERRASCAN_VERSION}/terrascan_${TERRASCAN_VERSION}_Linux_x86_64.tar.gz" \
      && tar -xzf /tmp/terrascan.tar.gz -C /usr/local/bin terrascan \
      && chmod +x /usr/local/bin/terrascan \
      || true; \
    rm -f /tmp/terrascan.tar.gz; \
    \
    # kube-bench
    curl -fsSL -o /tmp/kube-bench.tar.gz \
      "https://github.com/aquasecurity/kube-bench/releases/download/v${KUBE_BENCH_VERSION}/kube-bench_${KUBE_BENCH_VERSION}_linux_${goa}.tar.gz" \
      && tar -xzf /tmp/kube-bench.tar.gz -C /usr/local/bin kube-bench \
      && chmod +x /usr/local/bin/kube-bench \
      || true; \
    rm -f /tmp/kube-bench.tar.gz; \
    \
    # kube-hunter via pip later; docker-bench is a script clone
    git clone --depth 1 https://github.com/docker/docker-bench-security.git /opt/docker-bench-security \
      && printf '%s\n' '#!/bin/sh' 'cd /opt/docker-bench-security && exec sh docker-bench-security.sh "$@"' \
         > /usr/local/bin/docker-bench-security \
      && chmod +x /usr/local/bin/docker-bench-security \
      || true

# ---------------------------------------------------------------------------
# 6) Ruby gems (evil-winrm, zsteg)
# ---------------------------------------------------------------------------
RUN set -eux; \
    gem install --no-document evil-winrm zsteg || \
      gem install --no-document evil-winrm || true

# ---------------------------------------------------------------------------
# 7) Python API + security libraries
# ---------------------------------------------------------------------------
WORKDIR /opt/arsenal-engine
COPY requirements.txt /opt/arsenal-engine/requirements.txt

# Core API deps always
RUN pip3 install --break-system-packages --no-cache-dir --ignore-installed \
      -r /opt/arsenal-engine/requirements.txt

# Extra Python security tools used by health-check /api routes
RUN set -eux; \
    pip3 install --break-system-packages --no-cache-dir --ignore-installed \
      httpie \
      uro \
      paramspider \
      checkov \
      kube-hunter \
      shodan \
      censys \
      recon-ng \
      || true; \
    if [ "${INSTALL_CLOUD}" = "1" ]; then \
      pip3 install --break-system-packages --no-cache-dir --ignore-installed \
        prowler \
        scoutsuite \
        pacu \
        || true; \
    fi; \
    pip3 install --break-system-packages --no-cache-dir --ignore-installed \
      volatility3 \
      || true

# Heavy binary-exploitation Python stack (optional — large native wheels)
RUN set -eux; \
    if [ "${INSTALL_BINARY_HEAVY}" = "1" ]; then \
      pip3 install --break-system-packages --no-cache-dir --ignore-installed \
        'pwntools>=4.10.0,<5.0.0' \
        'angr>=9.2.0,<10.0.0' \
        ropgadget \
        one-gadget \
        || true; \
      pip3 install --break-system-packages --no-cache-dir --ignore-installed ropper || true; \
    fi

# Convenience wrappers so `which` finds common aliases arsenal health probes for
RUN set -eux; \
    if command -v nxc >/dev/null 2>&1 && [ ! -e /usr/local/bin/netexec ]; then \
      ln -sf "$(command -v nxc)" /usr/local/bin/netexec; \
    fi; \
    if command -v netexec >/dev/null 2>&1 && [ ! -e /usr/local/bin/nxc ]; then \
      ln -sf "$(command -v netexec)" /usr/local/bin/nxc; \
    fi; \
    if command -v volatility3 >/dev/null 2>&1 && [ ! -e /usr/local/bin/vol ]; then \
      ln -sf "$(command -v volatility3)" /usr/local/bin/vol; \
    fi; \
    # pwninit / one-gadget may install as scripts
    true

# ---------------------------------------------------------------------------
# 8) Application + ttyd theme
# ---------------------------------------------------------------------------
COPY arsenal_engine.py /opt/arsenal-engine/arsenal_engine.py
COPY LICENSE.upstream /opt/arsenal-engine/LICENSE.upstream
COPY start-ttyd.sh /usr/local/bin/start-ttyd.sh
COPY ttyd-index.html /usr/local/share/ttyd-index.html
RUN chmod +x /usr/local/bin/start-ttyd.sh

# Tool inventory helper (used in health / docs; does not fail build)
RUN set -eu; \
    printf '%s\n' \
      '#!/bin/sh' \
      'tools="nmap masscan rustscan amass subfinder nuclei gobuster feroxbuster ffuf dirb dirsearch nikto sqlmap wpscan httpx katana dalfox wafw00f hydra john hashcat medusa netexec nxc enum4linux-ng responder smbmap gdb radare2 binwalk trivy prowler checkov kube-hunter curl python3 ttyd chromium chromium-browser google-chrome"' \
      'ok=0; miss=0' \
      'for t in $tools; do' \
      '  if command -v "$t" >/dev/null 2>&1; then echo "OK  $t"; ok=$((ok+1));' \
      '  else echo "MISS $t"; miss=$((miss+1)); fi' \
      'done' \
      'echo "summary ok=$ok miss=$miss"' \
      > /usr/local/bin/arsenal-tool-check && chmod +x /usr/local/bin/arsenal-tool-check; \
    /usr/local/bin/arsenal-tool-check || true

ENV API_PORT=8888 \
    TTYD_PORT=7681 \
    # Chromium headless defaults for browser-agent
    CHROME_BIN=/usr/bin/chromium \
    CHROMIUM_FLAGS="--no-sandbox --disable-dev-shm-usage --disable-gpu --headless=new"

WORKDIR /opt/arsenal-engine
EXPOSE 8888 7681

# Document intentionally not installed (license / multi-GB / needs external keys):
#   ghidra, burpsuite, ida-free, binary-ninja, maltego, autopsy-gui, falco full,
#   metasploit-framework (use talon msf_rpc service), hashcat GPU drivers
#
# BlackArch alternative: larger package set, less stable Docker base, different
# package names — Kali remains the recommended production base for this engine.

HEALTHCHECK --interval=30s --timeout=5s --start-period=40s --retries=3 \
  CMD curl -fsS "http://127.0.0.1:${API_PORT}/health" >/dev/null || exit 1

CMD sh -c '/usr/local/bin/start-ttyd.sh & exec python3 arsenal_engine.py --port "${API_PORT}"'
