# OMNIXIUS — Infrastructure

Platform & Infrastructure Division: containerization, orchestration, IaC, CI/CD, monitoring. See **PLATFORM_INFRASTRUCTURE.md** in repo root.

- **docker-compose.yml** — local dev: Go API, Rust service, AI.
- **docker-compose.monitoring.yml** — optional Prometheus + Grafana (`docker compose -f docker-compose.monitoring.yml up -d`).
- **k8s/** — Kubernetes manifests (namespace placeholder; add Deployments/Services/Ingress as needed).
- **terraform/** — Infrastructure as Code (placeholder; add cloud + K8s modules).
- **monitoring/prometheus.yml** — Prometheus config for scraping metrics.
- **kafka/** — Apache Kafka placeholder for event streaming at scale.
- **scripts/deploy.sh** — deploy script (Bash).
- **scripts/health_check.py** — health check for API, Rust, AI (Python).

Run from repo root:
```bash
cd infra && docker compose up -d --build
python scripts/health_check.py
```
