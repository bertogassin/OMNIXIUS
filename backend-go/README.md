# OMNIXIUS API (Go)

REST API: auth, marketplace, orders, internal mail. **Stack: Go only** (per project language policy).

## Run

```bash
cd backend-go
go build -o omnixius-api .
./omnixius-api
# or: go run .
```

API: `http://localhost:3000`

## Env

- `PORT` — default 3000
- `DB_PATH` — default `db/omnixius.db`
- `ALLOWED_ORIGINS` — comma-separated origins for CORS; empty = `*` (dev)
- `DILITHIUM_PUBLIC_KEY` / `DILITHIUM_PRIVATE_KEY` — base64 PQC keys (optional; ephemeral if unset)

## Endpoints

Same as before: `/api/auth/*`, `/api/users/*`, `/api/products/*`, `/api/orders/*`, `/api/conversations`, `/api/messages/*`. See PLATFORM.md or the Node backend README for the list.

## DB

SQLite. Schema in `db/schema.sql` (embedded). Versioned migrations in `db/migrations/` (e.g. `003_add_feature.sql`); applied on startup. DB file created on first run.

## CI

GitHub Actions (`.github/workflows/ci.yml`): on push/PR to `main`/`master`, runs `go test ./...` and `go build` in `backend-go`.
