# Talon Production Deployment Guide

Step-by-step guide to deploying Talon in a production environment. Covers
secrets, TLS, networking, backups, monitoring, and high availability.

**Read [ARCHITECTURE.md](ARCHITECTURE.md) first** for the system design and
service inventory. This document covers the operational "how".

---

## 0. Production checklist

Before going live, verify every item:

- [ ] **Secrets**: all required env vars set via a secret manager (not `.env`)
- [ ] **TLS**: reverse proxy terminates TLS for dashboard and API
- [ ] **DB encryption**: `sslmode=require` on `the Talon DB connection env var`
- [ ] **Auth enabled**: `TALON_AUTH_DISABLED` unset or `false`
- [ ] **Network isolation**: core not directly exposed to the internet
- [ ] **Resource limits**: confirmed per service (already in docker-compose.yml)
- [ ] **Log rotation**: confirmed (json-file, 10m x 3, already configured)
- [ ] **Healthchecks**: all services have healthcheck (already configured)
- [ ] **Backup strategy**: Postgres + RabbitMQ data backed up off-host
- [ ] **Monitoring**: health endpoint scraped; container metrics collected
- [ ] **Authorized targets only**: scope-of-work / ROE enforced procedurally
- [ ] **docker.sock**: access restricted; consider rootless Docker or gVisor
- [ ] **Rate limiting**: at the reverse proxy layer

---

## 1. Prerequisites

### 1.1 Host requirements

| Deployment size | RAM | CPU | Disk | Use case |
|----------------|-----|-----|------|----------|
| **Single-host (lab/staging)** | 16 GB | 4 cores | 50 GB SSD | Core stack, no local LLM |
| **Single-host + local LLM** | 32 GB | 8 cores | 100 GB SSD | Add ollama (8 GB) or onnx-slm (4 GB) |
| **Multi-host (prod)** | 16 GB+ per host | 4+ cores per host | 100 GB SSD per host | Separate arsenal/host, data layer, dashboard |

### 1.2 Software

- Docker Engine 27+ with Compose v2
- Go 1.25+ (for host-side CLI build)
- An LLM backend (Bedrock, OpenAI-compatible, Ollama, or ONNX)

### 1.3 Network

All services use `network_mode: host` by design (simplifies reverse shells,
MSF RPC, and ttyd loopback proxy). Each service binds a distinct localhost port:

| Port | Service | Exposure |
|------|---------|----------|
| 8000 | talon-core | **Loopback only** — proxy from dashboard |
| 3000 | dashboard | **Public** — via reverse proxy with TLS |
| 8888 | arsenal-engine API | Loopback only |
| 7681 | arsenal-engine ttyd | Loopback only |
| 5554 | metasploit RPC | Loopback only |
| 5432 | postgres | Loopback only |
| 6380 | redis | Loopback only |
| 5672 | rabbitmq AMQP | Loopback (or private network for multi-host) |
| 15672 | rabbitmq management | Loopback only |
| 8090 | onnx-slm (profile) | Loopback only |
| 11434 | ollama (profile) | Loopback only |

**Firewall rule**: only expose port 443 (HTTPS) via the reverse proxy.
Everything else stays on localhost.

---

## 2. Secrets management

### 2.1 Required secrets (compose will refuse to start without these)

```bash
MSF_PASSWORD=<random 32+ chars>          # Metasploit RPC
RABBITMQ_PASSWORD=<random 32+ chars>     # RabbitMQ
TALON_PG_PASSWORD=<random 32+ chars>     # Postgres
TALON_ADMIN_PASSWORD=<random 16+ chars>  # Dashboard + CLI login
```

Generate strong secrets:

```bash
openssl rand -base64 32
```

### 2.2 LLM credentials (provider-dependent)

```bash
# Bedrock
AWS_ACCESS_KEY_ID=<key>
AWS_SECRET_ACCESS_KEY=<secret>
AWS_REGION=us-east-1

# OpenAI-compatible
OPENAI_API_KEY=<key>
OPENAI_BASE_URL=https://api.example.com/v1
```

### 2.3 Inject secrets via Docker secrets (recommended for prod)

Create a `secrets.env` file (not committed, on disk with `chmod 600`):

```bash
cp .env.example secrets.env
# Fill in real values
chmod 600 secrets.env
```

Deploy with:

```bash
docker compose --env-file secrets.env up -d --build
```

For enterprise: use HashiCorp Vault, AWS Secrets Manager, or Docker Swarm
secrets to inject env vars without files on disk.

### 2.4 Rotate secrets

| Secret | Rotation steps |
|--------|----------------|
| `TALON_ADMIN_PASSWORD` | Change env, restart core. Old sessions expire on next request (Redis cache TTL: 60s). |
| `TALON_PG_PASSWORD` | Change Postgres password, update env, restart core + postgres. |
| `MSF_PASSWORD` | Change env, restart metasploit + core. |
| `RABBITMQ_PASSWORD` | Change RabbitMQ user password, update env, restart rabbitmq + relay. |
| LLM API key | Change env, restart core. Active runs continue with old key until restart. |

---

## 3. TLS and reverse proxy

Talon does not terminate TLS natively. Put a reverse proxy (Caddy, Traefik,
or nginx) in front of the dashboard. The proxy handles:
- TLS certificate management (Let's Encrypt or custom)
- HTTP-to-HTTPS redirect
- Rate limiting
- Request body size limits

### 3.1 Caddy (simplest — automatic HTTPS)

Create `Caddyfile`:

```caddyfile
talon.example.com {
    reverse_proxy localhost:3000

    # Rate limit: 10 req/s per IP
    rate_limit {
        zone dynamic 10r/s 50
    }

    # Body size limit
    request_body {
        max_size 10MB
    }

    # Security headers
    header {
        Strict-Transport-Security "max-age=31536000; includeSubDomains"
        X-Content-Type-Options nosniff
        X-Frame-Options DENY
        Referrer-Policy strict-origin-when-cross-origin
    }
}
```

Run Caddy:

```bash
docker run -d --name caddy \
  --network host \
  -v $PWD/Caddyfile:/etc/caddy/Caddyfile:ro \
  -v caddy_data:/data \
  caddy:2-alpine
```

### 3.2 Traefik (recommended for multi-service prod)

Add to `docker-compose.yml`:

```yaml
  traefik:
    image: traefik:v3
    container_name: talon_traefik
    command:
      - "--providers.docker=true"
      - "--providers.docker.exposedbydefault=false"
      - "--entrypoints.web.address=:80"
      - "--entrypoints.websecure.address=:443"
      - "--certificatesresolvers.le.acme.email=you@example.com"
      - "--certificatesresolvers.le.acme.storage=/acme.json"
      - "--certificatesresolvers.le.acme.tlschallenge=true"
      - "--entrypoints.web.http.redirections.entrypoint.to=websecure"
      - "--entrypoints.web.http.redirections.entrypoint.scheme=https"
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock:ro
      - traefik_acme:/acme.json
    restart: unless-stopped
    network_mode: host

volumes:
  traefik_acme:
```

Label the dashboard service:

```yaml
  dashboard:
    labels:
      - "traefik.enable=true"
      - "traefik.http.routers.dashboard.rule=Host(`talon.example.com`)"
      - "traefik.http.routers.dashboard.entrypoints=websecure"
      - "traefik.http.routers.dashboard.tls.certresolver=le"
```

---

## 4. Database configuration (Postgres)

### 4.1 Enable SSL on the connection

Change the `the Talon DB connection env var` in your environment:

```bash
# Before (insecure — default in compose for local dev)
the database connection string (see .env.example)

# After (prod)
the database connection string (see .env.example)
```

For full verification (`sslmode=verify-full`), provision a server certificate
on the Postgres container and mount it.

### 4.2 Performance tuning

Tunable via environment variables (already wired in docker-compose.yml):

```bash
PG_SHARED_BUFFERS=256MB    # Default: 128MB. Set to 25% of available RAM.
PG_MAX_CONNECTIONS=100     # Default: 100. Increase if dashboard + relay need more.
```

### 4.3 Backup and restore

**Backup (pg_dump):**

```bash
# Daily snapshot (cron)
docker exec talon_postgres pg_dump -U talon talon | gzip > backup_$(date +%F).sql.gz

# Or from the host (if psql client installed)
pg_dump -h localhost -U talon talon | gzip > backup_$(date +%F).sql.gz
```

**Automated backup script** (`scripts/backup-postgres.sh`):

```bash
#!/bin/bash
set -eu
BACKUP_DIR="${BACKUP_DIR:-./backups}"
mkdir -p "$BACKUP_DIR"
docker exec talon_postgres pg_dump -U talon talon | gzip > "$BACKUP_DIR/talon_$(date +%Y%m%d_%H%M%S).sql.gz"
# Retain last 30 days
find "$BACKUP_DIR" -name "talon_*.sql.gz" -mtime +30 -delete
echo "backup complete: $(ls -t "$BACKUP_DIR" | head -1)"
```

Add to crontab:

```cron
0 2 * * * /home/operator/talon/scripts/backup-postgres.sh >> /var/log/talon-backup.log 2>&1
```

**Restore:**

```bash
gunzip < backup_20240115.sql.gz | docker exec -i talon_postgres psql -U talon talon
```

**Off-host backup**: copy `./backups/` to S3, GCS, or offsite storage:

```bash
aws s3 sync ./backups/ s3://your-bucket/talon-backups/ --exclude "*" --include "*.sql.gz"
```

---

## 5. RabbitMQ configuration

### 5.1 Backup

RabbitMQ data lives in `./rabbitmq-data/`. For definitions (queues, users,
permissions), export via the management API:

```bash
curl -u talon:$RABBITMQ_PASSWORD \
  http://localhost:15672/api/definitions > rabbitmq_definitions.json
```

### 5.2 HA (quorum queue, multi-node)

For production durability, run a 3-node RabbitMQ cluster with quorum queues.
This prevents message loss when a node fails. See the
[RabbitMQ cluster guide](https://www.rabbitmq.com/clustering.html).

---

## 6. Monitoring and observability

### 6.1 Health endpoints

```bash
# Core liveness
curl -s http://localhost:8000/health
# {"service":"talon-core","status":"ok"}

# Per-service health (7 probes: postgres, arsenal, msfrpcd, rabbitmq, ollama, onnx-slm, redis)
curl -s http://localhost:8000/health/services | jq .

# Dashboard liveness
curl -s -o /dev/null -w "%{http_code}" http://localhost:3000/
# 200

# Container health (from compose healthchecks)
docker inspect --format='{{.State.Health.Status}}' talon_core
# healthy
```

### 6.2 Log aggregation

All services log to `json-file` with rotation (10 MB x 3 files). For prod,
ship logs to a central collector:

**Option A — Fluent Bit sidecar:**

```yaml
  fluentbit:
    image: fluent/fluent-bit:3
    container_name: talon_fluentbit
    user: root
    volumes:
      - /var/lib/docker/containers:/var/lib/docker/containers:ro
      - ./fluent-bit.conf:/fluent-bit/etc/fluent-bit.conf:ro
    network_mode: host
    restart: unless-stopped
```

**Option B — Docker log driver:**

Set the log driver globally in `/etc/docker/daemon.json`:

```json
{
  "log-driver": "fluentd",
  "log-opts": { "fluentd-address": "localhost:24224" }
}
```

### 6.3 Metrics (Prometheus + Grafana)

Talon does not yet expose a Prometheus metrics endpoint. Until it does,
monitor at the container level:

- **cAdvisor** for container CPU/memory/network
- **Postgres exporter** for DB metrics
- **RabbitMQ exporter** for queue depth, consumer count
- **Node exporter** for host metrics

Alert on:
- Container restart count > 0 (unexpected restarts)
- Postgres connections > 80% of `max_connections`
- RabbitMQ queue depth > threshold (relay backlog)
- Any service healthcheck `unhealthy` for > 2 minutes
- Host disk usage > 80% (logs, backups, `.pg-data`)

---

## 7. High availability topology

The default single-host deployment has several single points of failure. For
production availability beyond "best effort", use this topology:

```
                    ┌─────────────┐
                    │  Load Balancer │  (AWS ALB, HAProxy, Cloudflare)
                    │  TLS + health  │
                    └──────┬──────┘
                           │
              ┌────────────┼────────────┐
              │            │            │
         Dashboard 1   Dashboard 2   Dashboard N
         (Next.js)     (Next.js)     (Next.js)
              │            │            │
              └────────────┼────────────┘
                           │
                    ┌──────┴──────┐
                    │ talon-core  │  (single active orchestrator instance;
                    │  :8000      │   run state is in-memory)
                    └──────┬──────┘
                           │
              ┌────────────┼────────────┐
              │            │            │
         Relay 1       Relay 2      Relay N   ← scale workers for throughput
         (AMQP)        (AMQP)       (AMQP)
              │            │            │
              └────────────┼────────────┘
                           │
                    ┌──────┴──────┐
                    │  RabbitMQ   │  3-node quorum cluster
                    │  Cluster    │
                    └─────────────┘

         ┌─────────────┐         ┌─────────────┐
         │  Postgres   │ ──WAL──►│  Postgres   │
         │  Primary    │  stream │  Replica    │  ← dashboard reads
         └─────────────┘         └─────────────┘

         ┌─────────────┐
         │    Redis    │  Single (cache only; sentinel not needed)
         │   Sentinel  │  (optional for HA cache)
         └─────────────┘
```

### What to replicate

| Component | HA strategy | Complexity |
|-----------|-------------|------------|
| Dashboard | Stateless, replicate freely behind LB | Low |
| talon-relay | Stateless worker, replicate behind queue | Low |
| RabbitMQ | 3-node quorum cluster | Medium |
| Postgres | Primary + streaming replica(s) | Medium |
| Redis | Single (or Sentinel for HA cache) | Low-Medium |
| **talon-core** | **Single active instance** (orchestrator state is in-memory) | **Cannot HA without code change** |

**talon-core is the constraint**: its orchestrator holds run state in memory.
To run multiple core instances safely, the orchestrator needs a distributed
state store (Postgres-backed run state with optimistic locking, or a separate
coordination service like etcd). This is a roadmap item, not a config change.

### Workaround for core HA today

1. Run one **active** core instance.
2. Run a **standby** core instance with `TALON_CORE_STANDBY=true` (future env)
   that serves the API but does not run the orchestrator.
3. Use keepalived or the LB health check to promote standby on failure.

---

## 8. Docker socket security

`talon-core` and `talon-relay` mount `/var/run/docker.sock` because Forge
(the codegen sandbox) needs to create Docker containers. This grants
effectively root access to the host.

### Mitigation options (in order of preference)

| Option | Security | Effort |
|--------|----------|--------|
| **Rootless Docker** | High — socket has no host root | Medium (reconfigure Docker daemon) |
| **gVisor / runsc** | High — sandboxed runtime | Low (add `--runtime=runsc` to Forge containers) |
| **Socket proxy** ( Tecnativa/docker-socket-proxy) | Medium — limit API calls to container create/start | Low |
| **Dedicated host** | High — isolate arsenal + core on a throwaway host | Medium |

### Socket proxy example

```yaml
  docker-socket-proxy:
    image: tecnativa/docker-socket-proxy:0.3
    container_name: talon_docker_proxy
    environment:
      CONTAINERS: 1
      POST: 1
      # Block everything else: IMAGES, NETWORKS, VOLUMES, EXEC, etc.
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock:ro
    network_mode: host
    restart: unless-stopped
```

Then mount the proxy socket (not the real one) into talon-core. This requires
Forge to use the proxy endpoint, which may need a code change.

---

## 9. Scaling operations

### 9.1 Scale relay workers (increase throughput)

```bash
# Run 3 relay workers
docker compose up -d --scale talon-relay=3
```

Each worker consumes from the same `execute_agent_task` queue. RabbitMQ
distributes work round-robin. No code changes needed.

### 9.2 Scale the dashboard (increase availability)

```bash
docker compose up -d --scale dashboard=2
```

Put both instances behind a load balancer. Sessions are in Postgres/Redis,
so any instance can serve any request.

### 9.3 Scale Postgres reads (read-heavy dashboards)

Add a read replica and route dashboard list/read queries to it:

```sql
-- On the primary, create a replication user
CREATE ROLE replicator WITH REPLICATION LOGIN PASSWORD 'replica_password';
```

Configure streaming replication on a second Postgres instance. Point the
dashboard's read-only queries at the replica via a separate env var (requires
code change to support read/write splitting).

---

## 10. Deployment procedures

### 10.1 First deployment

```bash
# 1. Clone
git clone <repo> talon && cd talon

# 2. Configure secrets
cp .env.example secrets.env
chmod 600 secrets.env
# Edit secrets.env — fill in all required values

# 3. Build and start
docker compose --env-file secrets.env up -d --build

# 4. Verify
docker compose --env-file secrets.env exec talon-core wget -qO- http://localhost:8000/health
docker compose --env-file secrets.env exec talon-core wget -qO- http://localhost:8000/health/services | jq .

# 5. Build CLI
go build -o bin/talon ./cmd/talon

# 6. Smoke test
./bin/talon status
```

### 10.2 Update (zero-downtime for dashboard and relay)

```bash
# Pull latest
git pull

# Rebuild images
docker compose --env-file secrets.env build

# Roll dashboard (zero-downtime if behind LB)
docker compose --env-file secrets.env up -d --no-deps --build dashboard

# Roll relay workers (one at a time if scaled)
docker compose --env-file secrets.env up -d --no-deps --build talon-relay

# Roll core (brief downtime — active runs are marked error on restart)
docker compose --env-file secrets.env up -d --no-deps --build talon-core
```

**Note**: rolling talon-core marks in-flight runs as `error` (crash recovery).
Schedule core rollouts during a maintenance window.

### 10.3 Rollback

```bash
# Revert code
git checkout <previous-tag>

# Rebuild and restart
docker compose --env-file secrets.env up -d --build

# Restore DB if needed
gunzip < backups/talon_20240115_020000.sql.gz | docker exec -i talon_postgres psql -U talon talon
```

---

## 11. Incident response runbook

| Symptom | First check | Fix |
|---------|-------------|-----|
| Dashboard returns 502 | `docker logs talon_dashboard --tail 50` | Restart dashboard; check if core is up |
| `talon status` shows unhealthy core | `docker logs talon_core --tail 100` | Check DB connectivity, LLM provider, MCP children |
| All health probes show offline | Redis down? Network? | `docker restart talon_redis`; check `docker network` |
| Runs stuck in `awaiting_approval` | HITL gate not actioned | `talon run approve <id>` or `--auto-approve` (lab only) |
| `run_exploit` returns "No session" | msfrpcd down or password mismatch | Check `MSF_PASSWORD` matches; `docker restart msf_rpc` |
| Forge fails with Docker error | docker.sock not mounted or daemon down | Verify socket mount; `docker ps` from inside core |
| Postgres disk full | `.pg-data` grew; backups not rotating | `docker exec talon_postgres vacuumdb -U talon`; clean old backups |
| LLM calls timeout | Provider rate limit or outage | Check `OPENAI_HTTP_TIMEOUT`; increase or switch provider |
| RabbitMQ queue growing | Relay workers down or stuck | `docker compose up -d --scale talon-relay=N`; check relay logs |
| Arsenal tools not found | arsenal-engine rebuild needed | `docker compose build arsenal-engine` |

---

## 12. Compliance and responsible use

- **Authorization**: only test systems you own or have **written authorization**
  to test. Talon enforces this procedurally, not technically.
- **Audit trail**: all runs, tool calls, findings, and reports persist in
  Postgres with timestamps. Export via `talon run export <id>` for evidence.
- **Data retention**: configure a retention policy for `runs` table and
  `talon-data/` traffic logs based on your compliance requirements.
- **Scope enforcement**: use the Engagements module (`/ops`) to define scope
  and ROE before starting operations.
- **Lab isolation**: the `vuln` profile target is intentional malware. Never
  deploy it outside an isolated lab network.

---

## Related documents

| Document | Scope |
|----------|-------|
| [ARCHITECTURE.md](ARCHITECTURE.md) | System design, scalability, security model |
| [PRODUCT.md](PRODUCT.md) | Shipped feature surface |
| [README.md](../README.md) | Getting started, CLI, API reference |
