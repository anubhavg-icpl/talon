# Talon Arsenal Engine (Kali)

Self-built tool-execution backend for Talon (`cmd/talon-arsenal` / dashboard Kali shell).

## Why Kali (not BlackArch)?

| | **Kali rolling** (this image) | BlackArch |
|--|------------------------------|-----------|
| Docker base | Official, well-maintained | Community, heavier |
| Package names | Match `arsenal_engine.py` | Often different |
| Reproducibility | Better for CI / compose | More churn |
| Size | Large but controllable | Typically larger |
| Tool coverage | Apt + Go + pip layers | “Everything” but fragile |

**Verdict:** Kali + explicit installs of missing tools is the better production default. BlackArch is fine on bare metal for research, not ideal as a compose service base.

## Build

```bash
# Default production profile (cloud scanners + pwntools/angr + chromium)
docker compose build arsenal-engine

# Leaner (skip cloud binaries + heavy binary Python stack)
docker build \
  --build-arg INSTALL_CLOUD=0 \
  --build-arg INSTALL_BINARY_HEAVY=0 \
  -t talon-arsenal-engine \
  ./arsenal-engine

# With wireless suite
docker build --build-arg INSTALL_WIRELESS=1 -t talon-arsenal-engine ./arsenal-engine
```

## Build args

| Arg | Default | Meaning |
|-----|---------|---------|
| `INSTALL_CLOUD` | `1` | trivy, terrascan, kube-bench, prowler, scoutsuite, pacu, docker-bench |
| `INSTALL_BINARY_HEAVY` | `1` | pwntools, angr, ropgadget, one-gadget, ropper |
| `INSTALL_BROWSER` | `1` | Chromium + chromedriver for browser-agent |
| `INSTALL_WIRELESS` | `0` | aircrack-ng, kismet, reaver, … |

## Run (compose)

```bash
docker compose up -d arsenal-engine
curl -s http://127.0.0.1:8888/health | jq .
docker exec arsenal_engine arsenal-tool-check
```

## Tool layout (how tools get in)

1. **Kali apt** — nmap, masscan, nuclei (also Go), gobuster, sqlmap, hydra, john, hashcat, netexec, responder, gdb, radare2, …
2. **Go install** — httpx, katana, dalfox, gau, waybackurls, hakrawler, jaeles, ffuf/gobuster/nuclei/subfinder (latest upstream)
3. **GitHub releases** — rustscan, ttyd, trivy, terrascan, kube-bench
4. **pip** — flask/fastmcp stack, volatility3, checkov, kube-hunter, prowler, scoutsuite, paramspider, uro, …
5. **gem** — evil-winrm, zsteg

## Intentionally not in the image

| Tool | Why |
|------|-----|
| Metasploit full | Separate `msf_rpc` service (`kali-msf/`) |
| Ghidra | Multi-GB + JDK; mount or install manually if needed |
| Burp Suite | Commercial license |
| IDA / Binary Ninja | Commercial |
| Maltego | License + GUI |
| Hashcat GPU drivers | Host NVIDIA stack; install on GPU nodes |
| OSINT API CLIs needing keys | Optional; install + inject keys per engagement |

## Health categories (`GET /health`)

`arsenal_engine.py` probes `which <tool>` for essential / network / web / password / binary / forensics / cloud / osint groups. Missing optional tools return `"success": false` for that name without crashing the API.

## MCP / HexStrike notes

This container runs **`arsenal_engine.py`** (HTTP API on `:8888`), not the separate `hexstrike_mcp.py` stdio bridge. Point MCP clients at:

```json
{
  "mcpServers": {
    "talon-arsenal": {
      "command": "python3",
      "args": ["hexstrike_mcp.py", "--server", "http://127.0.0.1:8888"]
    }
  }
}
```

(only if you vendor the MCP client script separately).

## Shell (ttyd)

- Loopback `127.0.0.1:7681`, base path `/shell`
- SSO via talon-core reverse proxy
- Theme: carbon black + electric red (`ttyd-index.html`)
