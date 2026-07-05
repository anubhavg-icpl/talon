# Multi-stage build producing all four binaries this platform ships.
# The api/worker containers spawn hexstrike-mcp and metasploit-mcp as local
# stdio subprocesses (see internal/mcpclient), so all four live in one image.
FROM golang:1.25-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/api ./cmd/api
RUN CGO_ENABLED=0 go build -o /out/worker ./cmd/worker
RUN CGO_ENABLED=0 go build -o /out/hexstrike-mcp ./cmd/hexstrike-mcp
RUN CGO_ENABLED=0 go build -o /out/metasploit-mcp ./cmd/metasploit-mcp

FROM alpine:3.20
# docker-cli: the codegen tool shells out to `docker` to run generated
# exploit code in a sandboxed sibling container -- mount the host's
# /var/run/docker.sock into this container for that to work (see
# docker-compose.yml).
RUN apk add --no-cache ca-certificates docker-cli
COPY --from=build /out/ /app/
WORKDIR /app
ENV HEXSTRIKE_MCP_PATH=/app/hexstrike-mcp
ENV METASPLOIT_MCP_PATH=/app/metasploit-mcp
