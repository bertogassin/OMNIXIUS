# Cache layer (placeholder)

For large-scale: **Redis** (or compatible) — sessions, rate-limit counters, API response cache, job queues.

- Add `docker-compose.cache.yml` or extend `docker-compose.yml` with Redis when needed.
- backend-go: optional Redis client; env `REDIS_URL`.
- No code here yet; this folder reserves the place in the architecture.
