# System Design Primer — Study Notes

Source: https://github.com/donnemartin/system-design-primer (cloned to `/tmp/system-design-primer`)
Purpose: distilled reference for making architectural decisions about the Talon platform.

## What the repo is
An organized, community-maintained study guide for large-scale system design and
system-design interviews. It is **reference material**, not a runnable system:
there is no Dockerfile, no services, no code runtime. The deliverable is a single
~1,840-line `README.md` plus Python reference solutions (`solutions/`) and Anki
flashcard decks (`resources/flash_cards/`).

## Repo layout
```
README.md              # the entire body of knowledge (EN, JA, zh-Hans, zh-TW)
solutions/
  system_design/       # pastebin, twitter, web_crawler, mint, scaling_aws, ...
  object_oriented_design/  # call_center, parking_lot, lru_cache, online_chat, ...
resources/
  flash_cards/         # *.apkg Anki decks
  study_guide.png
images/                # diagrams referenced by README
```

## Core mental model (the sections that matter for real architecture)

### 1. Performance vs scalability
- **Performance problem**: slow for a single user.
- **Scalability problem**: fast for one user, slow under load.
- Scalable ⇒ adding resources yields proportional performance gains.

### 2. Latency vs throughput
- Aim for **maximal throughput at acceptable latency**. Always profile, never
  tune blind.

### 3. CAP theorem
In the presence of network partitions (which are inevitable) you choose between:
- **CP** — consistency + partition tolerance (atomic reads/writes).
- **AP** — availability + partition tolerance (eventual consistency).
Networks are unreliable, so the real choice is C-vs-A under partition.

### 4. Consistency patterns
- **Weak**: reads may or may not see a write (real-time apps, VoIP, games).
- **Eventual**: reads see the write within ms (DNS, email, highly-available systems).
- **Strong**: reads always see the latest write (file systems, RDBMS, transactions).

### 5. Availability patterns
- **Fail-over**: active-passive (hot/cold standby) or active-active.
- **Replication**: master-slave, master-master (see Database section).
- Availability math: components **in series** multiply (worse); components
  **in parallel** compound (better). Two 99.9% services in series ⇒ 99.8%;
  in parallel ⇒ 99.9999%.
- "Number of nines": 3 nines = 8h 45m downtime/year; 4 nines = 52m/year.

### 6. DNS, CDN, load balancer, reverse proxy
- **DNS**: hierarchical; cache by TTL; managed providers (Cloudflare, Route 53)
  offer weighted, latency-based, geo routing. Vulnerable to DDoS.
- **CDN**: push (you upload on change) vs pull (CDN fetches on first request,
  TTL-driven). Offloads traffic from origin; costs scale with usage.
- **Load balancer**: distributes requests, prevents routing to unhealthy nodes,
  eliminates SPOF. Layer 4 (transport, cheaper) vs Layer 7 (application,
  smarter). SSL termination, session persistence. **LBs enable horizontal
  scaling** — but services must be stateless (sessions live in a shared store).
- **Reverse proxy**: centralized ingress, hides backends, SSL termination,
  compression, caching, static serving. A load balancer is one use of a reverse
  proxy; LB implies multiple backends, reverse proxy does not.

### 7. Application layer / microservices
- Separate the web layer from the application/platform layer so each scales
  independently. Single responsibility, small autonomous services.
- **Service discovery**: Consul, etcd, Zookeeper — name→address→port registry
  plus health checks.
- Cost: deployment and operational complexity go up.

### 8. Databases
- **RDBMS** — ACID. Scaling techniques: master-slave, master-master, federation
  (split by function), sharding (split by key), denormalization, SQL tuning.
- **NoSQL** — BASE (basically available, soft state, eventual consistency).
  Categories: key-value, document, wide-column, graph.
- **Replication trade-offs**: replication lag, slave replay bottleneck, more
  hardware + complexity, potential data loss on master failure.
- **Sharding trade-offs**: app-level routing, hot shards, rebalancing pain,
  cross-shard joins are hard — consistent hashing mitigates.
- **Denormalization**: trade write cost for read speed; great when reads ≫ writes.
- **SQL tuning**: benchmark (ab) + profile (slow query log); tighten schema,
  index strategically, avoid expensive joins, partition hot tables.

### 9. Cache
- Write-through, write-around, write-back; cache-aside (lazy load).
- Eviction: LRU, LFU. Redis/Memcached are the canonical stores.
- Cache the hot 20% that drives 80% of reads; invalidate on write.

### 10. Asynchronism
- **Message queues** (RabbitMQ, SQS): decouple producers from consumers,
  smooth traffic spikes, enable backpressure.
- **Event-driven / pub-sub** (Kafka, SNS): fan-out, replay, decoupled consumers.
- Workers behind a queue let the app layer scale independently of the web layer.

### 11. Communication protocols
- TCP/UDP, HTTP/1.1 vs HTTP/2 vs HTTP/3, WebSocket (full-duplex, used by Talon
  for streaming), gRPC (binary over HTTP/2, strong contracts).
- Long polling / server-sent events (SSE) as lighter-weight streaming options.

### 12. Security (the section Talon cares most about)
- Encrypt **in transit and at rest**.
- Sanitize all user input (XSS, SQLi). Use parameterized queries.
- **Principle of least privilege**.
- References: OWASP Top 10, API Security Checklist.

### 13. Appendix: numbers every engineer should know
```
L1 cache           0.5 ns     Main memory        100 ns
Branch mispredict  5 ns       SSD 4KB random     150 us
L2 cache           7 ns       1GB Ethernet 1KB   10 us
Mutex lock/unlock  25 ns      DC round-trip      500 us
Compress 1KB zippy 10 us      HDD seek           10 ms
                              CA→NL→CA packet    150 ms
```
Rule of thumb: 2,000 round trips/sec within a DC; 6–7 worldwide/sec.
Sequential read: ~4 GB/s RAM, ~1 GB/s SSD, ~100 MB/s 1Gbps, ~30 MB/s HDD.

## How this maps onto Talon
Talon is **not** a hyperscale web system; it is a security-operations platform
with: a Go HTTP core, a Next.js dashboard, Postgres (run history, auth), Redis
(cache), RabbitMQ (relay queue), Metasploit RPC, and an arsenal-engine (Kali
tool runner). The primer patterns that apply directly:

| Pattern                  | Talon applicability                                                 |
|--------------------------|---------------------------------------------------------------------|
| Statelessness + LB       | talon-core instances must stay stateless so they can be replicated. |
| Cache-aside (Redis)      | Already used for health probes / AI analysis / sessions.            |
| Message queue (RabbitMQ) | talon-relay consumes the queue; backpressure protects the core.     |
| RDBMS + replication      | Postgres for auth/run history; add a replica for read scaling.      |
| Reverse proxy            | Put Caddy/Traefik/nginx in front for TLS termination + routing.     |
| Health endpoints         | `GET /health` exists on talon-core; wire compose healthchecks.      |
| Least privilege          | Run containers non-root, drop ALL caps, read-only fs where possible.|
| Encrypt in transit       | TLS at the proxy; sslmode=require on the DB conn in prod.           |
| Backups / durability     | Postgres + RabbitMQ volumes need off-host backup.                   |

## Where the primer is *not* a guide
- It has no opinions on containerization, compose, Kubernetes, observability,
  secrets management, or CI/CD. Those decisions are ours to make.
- The security section is intentionally thin — for a security product we hold
  ourselves to a higher bar (see OWASP, container hardening guides).
