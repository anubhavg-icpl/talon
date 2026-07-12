# Multi-stage build producing the Talon platform binaries.
# talon-core/talon-relay spawn talon-arsenal and talon-strike as local
# stdio subprocesses (see internal/mcpclient). talon is the operator CLI.
FROM golang:1.25-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/talon-core ./cmd/talon-core
RUN CGO_ENABLED=0 go build -o /out/talon-relay ./cmd/talon-relay
RUN CGO_ENABLED=0 go build -o /out/talon-arsenal ./cmd/talon-arsenal
RUN CGO_ENABLED=0 go build -o /out/talon-strike ./cmd/talon-strike
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /out/talon ./cmd/talon

FROM alpine:3.20
# docker-cli: the forge (codegen) tool shells out to `docker` to run
# generated exploit code in a sandboxed sibling container -- mount the
# host's /var/run/docker.sock into this container for that to work (see
# docker-compose.yml).
RUN apk add --no-cache ca-certificates docker-cli
COPY --from=build /out/ /app/
WORKDIR /app
ENV HEXSTRIKE_MCP_PATH=/app/talon-arsenal
ENV METASPLOIT_MCP_PATH=/app/talon-strike
