# Talon Arsenal Engine
#
# Self-built tool-execution backend for Talon Arsenal (cmd/talon-arsenal),
# built from vendored third-party source under MIT License -- see
# LICENSE.upstream in this directory for the required copyright/permission
# notice. This is a backing tool-execution engine, not orchestration code --
# same role kali-msf/ plays for Metasploit. It stays Python because that's
# what the vendored engine is; nothing here is part of the Go rewrite.
FROM kalilinux/kali-rolling

ENV DEBIAN_FRONTEND=noninteractive
ENV GOPATH=/root/go
ENV PATH="${PATH}:/root/go/bin:/usr/local/go/bin"

# ---------------------------------------------------------------------------
# System tool arsenal. Kali's own repos cover the large majority of the
# 150+ tools the engine shells out to -- that's the whole reason to build
# on kali-rolling rather than a bare Debian/Ubuntu base.
# ---------------------------------------------------------------------------
RUN apt-get update && apt-get install -y --no-install-recommends \
    python3 python3-pip python3-venv python3-dev build-essential \
    git curl wget unzip golang-go \
    # Network & reconnaissance
    nmap masscan amass subfinder fierce dnsenum theharvester \
    netexec enum4linux-ng smbmap nbtscan arp-scan responder \
    autorecon \
    # Web application security
    gobuster feroxbuster ffuf dirb dirsearch nuclei nikto sqlmap wpscan \
    wafw00f arjun \
    # Password & authentication
    hydra john hashcat medusa patator hash-identifier \
    # Binary analysis & reverse engineering
    gdb radare2 binwalk checksec foremost steghide exiftool \
    # Forensics
    testdisk sleuthkit \
    && rm -rf /var/lib/apt/lists/*

# ---------------------------------------------------------------------------
# Go-based recon/fuzzing tools not reliably packaged in Kali's apt repos.
# ---------------------------------------------------------------------------
RUN go install -v github.com/projectdiscovery/httpx/cmd/httpx@latest && \
    go install -v github.com/projectdiscovery/katana/cmd/katana@latest && \
    go install -v github.com/hahwul/dalfox/v2@latest && \
    go install -v github.com/lc/gau/v2/cmd/gau@latest && \
    go install -v github.com/tomnomnom/waybackurls@latest && \
    go install -v github.com/hakluke/hakrawler@latest && \
    go install -v github.com/tomnomnom/anew@latest && \
    go install -v github.com/tomnomnom/qsreplace@latest

# ---------------------------------------------------------------------------
# Python dependencies (flask/fastmcp API layer + pwntools/angr binary
# exploitation stack) from the vendored, unmodified requirements.txt.
# ---------------------------------------------------------------------------
COPY requirements.txt /opt/arsenal-engine/requirements.txt
RUN pip3 install --break-system-packages --no-cache-dir -r /opt/arsenal-engine/requirements.txt

# ---------------------------------------------------------------------------
# Python packages available via pip that aren't in Kali's apt repos.
# ---------------------------------------------------------------------------
RUN pip3 install --break-system-packages --no-cache-dir \
    volatility3 paramspider uro x8

# ---------------------------------------------------------------------------
# NOT installed -- documented rather than silently missing. Any tool call
# routed to one of these returns a real "command not found" from the
# vendored server, not a fake success:
#   - Ghidra: multi-GB download + bundled JDK, install manually into the
#     image if binary analysis via Ghidra specifically is needed.
#   - Burp Suite: commercial license required.
#   - Cloud auditors (prowler, scout-suite, checkov, terrascan, trivy,
#     kube-hunter, kube-bench, docker-bench-security, falco): need live
#     cloud/cluster credentials to be useful even if installed; add them
#     (mostly `pip3 install --break-system-packages prowler scoutsuite
#     checkov` plus binary releases for terrascan/trivy/kube-bench) if
#     you're actually running cloud engagements.
#   - OSINT tools needing API keys (shodan-cli, censys-cli,
#     have-i-been-pwned, spiderfoot, maltego, recon-ng, sherlock,
#     social-analyzer): install per-tool and supply your own API keys.
#   - evil-winrm, ophcrack: GUI/Ruby-gem tooling outside this image's
#     scope for now.
# ---------------------------------------------------------------------------

WORKDIR /opt/arsenal-engine
COPY arsenal_engine.py .
COPY LICENSE.upstream ./LICENSE.upstream

ENV API_PORT=8888
EXPOSE 8888

CMD python3 arsenal_engine.py --port "${API_PORT}"
