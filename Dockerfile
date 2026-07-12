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

FROM alpine:3.20
RUN apk add --no-cache ca-certificates docker-cli
COPY --from=build /out/ /app/
WORKDIR /app
ENV HEXSTRIKE_MCP_PATH=/app/talon-arsenal \
    METASPLOIT_MCP_PATH=/app/talon-strike
# Default command is core; compose overrides for relay.
CMD ["/app/talon-core"]
