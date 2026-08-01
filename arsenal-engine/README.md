# Talon Arsenal Engine (BlackArch)

Self-built tool-execution backend for Talon. **Base distro: BlackArch (Arch + strap), not Kali.**

## Contracts (do not break)

| Surface | Value |
|---------|--------|
| Container name | `arsenal_engine` |
| HTTP API | `API_PORT` default **8888** → `arsenal_engine.py` |
| Web shell | `TTYD_PORT` default **7681**, loopback, path `/shell` |
| SSO | talon-core reverse-proxies `/shell` |
| Compose | `network_mode: host`, `privileged: true` |
| Health | `GET /health` |

## Why BlackArch

| | BlackArch (this image) | Kali |
|--|------------------------|------|
| Tool catalog | Very large pacman set | Large apt set |
| Base | Arch + blackarch.org strap | Debian rolling |
| Request | User-selected for arsenal | Still used by `kali-msf` / msf_rpc only |

Metasploit RPC remains the separate **`msf_rpc`** service (`kali-msf/`). This image does **not** replace it.

## Build

```bash
docker compose build arsenal-engine
docker compose up -d --force-recreate arsenal-engine

curl -s http://127.0.0.1:8888/health | head
docker exec arsenal_engine arsenal-tool-check
docker exec arsenal_engine cat /etc/os-release
# expect Arch / BlackArch lineage
```

### Build args

| Arg | Default | Meaning |
|-----|---------|---------|
| `INSTALL_CLOUD` | `1` | trivy, terrascan, kube-bench, prowler, scoutsuite, pacu |
| `INSTALL_BINARY_HEAVY` | `1` | pwntools, angr, ropgadget, one-gadget |
| `INSTALL_BROWSER` | `1` | Chromium + chromedriver |
| `INSTALL_WIRELESS` | `0` | aircrack-ng suite |

Lean:

```bash
docker build \
  --build-arg INSTALL_CLOUD=0 \
  --build-arg INSTALL_BINARY_HEAVY=0 \
  -t talon-arsenal-engine ./arsenal-engine
```

## Tool sources

1. **pacman / BlackArch** — nmap, masscan, amass, nuclei, gobuster, sqlmap, hydra, john, hashcat, netexec, responder, gdb, radare2, rustscan, …
2. **go install** — httpx, katana, dalfox, gau, waybackurls, hakrawler, jaeles, …
3. **GitHub releases** — ttyd, trivy, terrascan, kube-bench, docker-bench-security
4. **pip** — Flask/fastmcp stack + volatility3, checkov, kube-hunter, prowler, …
5. **gem** — evil-winrm, zsteg

Missing packages are **skipped** (`pac-opt`) so one rename does not fail the whole image. `arsenal-tool-check` reports OK/MISS.

## Intentionally external

| Tool | Where |
|------|--------|
| Metasploit full | `msf_rpc` / `kali-msf` |
| Ghidra / Burp / IDA | Manual / licensed |
| Hashcat GPU drivers | Host NVIDIA stack |

## Operator shell (Talon-themed BlackArch)

Login PTY (`bash -l` via ttyd) loads:

| File | Role |
|------|------|
| `/etc/talon/motd` | ASCII TALON banner (red) |
| `/etc/talon/fastfetch.jsonc` | **fastfetch** system card |
| `/etc/talon/bashrc` | Red PS1, aliases (`ff`, `tools`, `health`) |
| `/etc/profile.d/talon-arsenal.sh` | Login hook |
| `ttyd-index.html` | xterm carbon + electric red |

Aliases inside the shell:

```bash
ff       # fastfetch (Talon config)
tools    # arsenal-tool-check
health   # curl API /health
```

Dashboard route: **Arsenal Shell** (`/terminal`) — not “Kali”.
