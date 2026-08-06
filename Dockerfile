# Talon platform image (ONE Dockerfile for all Go binaries).
# Used by compose services talon-core and talon-relay (same image, different CMD).
#
# Contents: talon-core, talon-relay, talon-arsenal, talon-strike, talon (CLI)
# Core/relay spawn arsenal+strike over MCP stdio (HEXSTRIKE_MCP_PATH /
# METASPLOIT_MCP_PATH). Forge needs docker-cli + mounted docker.sock.
#
# Other images (separate, on purpose — different bases):
#   arsenal-engine/Dockerfile  — Kali tool runner
#   kali-msf/Dockerfile        — msfrpcd
#   vuln-target/Dockerfile     — lab target (targets: real | mimic)
FROM golang:1.25-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/talon-core ./cmd/talon-core \
 && CGO_ENABLED=0 go build -o /out/talon-relay ./cmd/talon-relay \
 && CGO_ENABLED=0 go build -o /out/talon-arsenal ./cmd/talon-arsenal \
 && CGO_ENABLED=0 go build -o /out/talon-strike ./cmd/talon-strike \
 && CGO_ENABLED=0 go build -ldflags="-s -w" -o /out/talon ./cmd/talon

# Runtime on glibc (debian-slim), NOT alpine/musl: the Lightpanda browser
# binary is glibc-linked and fails to exec on musl. The Go binaries are static
# (CGO_ENABLED=0) so they run on either base.
FROM debian:bookworm-slim
RUN apt-get update \
    && apt-get install -y --no-install-recommends ca-certificates curl \
    && rm -rf /var/lib/apt/lists/*
# Docker CLI only (daemon comes from the mounted /var/run/docker.sock).
# SECURITY: talon-core runs as root because Forge (codegen sandbox) needs the
# Docker socket. This is an intentional tradeoff for a pentest tool — the socket
# grants full daemon access. Hardening done in docker-compose.yml:
# security_opt: ["no-new-privileges:true"]. Roadmap: switch to a non-root USER
# and grant socket access via the 'docker' group / socket permissions.
COPY --from=docker:27-cli /usr/local/bin/docker /usr/local/bin/docker
# Lightpanda browser-automation MCP (optional): glibc binary — runs on this
# base. talon-core self-gates via LIGHTPANDA_MCP_PATH; `|| echo` keeps the
# build green on fetch failure so the image is never blocked on it.
# SECURITY/SUPPLY-CHAIN: lightpanda is fetched from a nightly release tag with
# no checksum verification. For production, pin to a specific release tag and
# verify the SHA256 checksum before installing.
RUN curl -fsSL -o /usr/local/bin/lightpanda \
      https://github.com/lightpanda-io/browser/releases/download/nightly/lightpanda-x86_64-linux \
    && chmod +x /usr/local/bin/lightpanda \
    && /usr/local/bin/lightpanda version >/dev/null 2>&1 \
    && echo "lightpanda: installed" \
    || echo "lightpanda: install skipped (optional MCP)"
COPY --from=build /out/ /app/
# Methodology skill pack (GET /skills + agent prompt injection).
COPY --from=build /src/skills /app/skills
WORKDIR /app
ENV HEXSTRIKE_MCP_PATH=/app/talon-arsenal \
    METASPLOIT_MCP_PATH=/app/talon-strike \
    TALON_SKILLS_DIR=/app/skills
# Default command is core; compose overrides for relay.
CMD ["/app/talon-core"]
