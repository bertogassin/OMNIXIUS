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

**Register / login via Go (no JS):** Open `http://localhost:3000/register` or `http://localhost:3000/login`. Submit the form; if `APP_URL` is set (e.g. your GitHub Pages URL), you are redirected to the app with token and API URL saved.

## Run with Docker

```bash
docker build -t omnixius-api .
docker run -p 3000:3000 -e DB_PATH=/data/omnixius.db -v omnixius-data:/data omnixius-api
```

## Env

Copy `.env.example` to `.env` in production. Variables:

- `PORT` — default 3000
- `DB_PATH` — default `db/omnixius.db`
- `APP_URL` — frontend base URL for redirect after register/login (e.g. `https://bertogassin.github.io/OMNIXIUS`)
- `ALLOWED_ORIGINS` — comma-separated origins for CORS; empty = `*` (dev)
- `DILITHIUM_PUBLIC_KEY` / `DILITHIUM_PRIVATE_KEY` — base64 PQC keys (optional; ephemeral if unset)

## Endpoints

Same as before: `/api/auth/*`, `/api/users/*`, `/api/products/*`, `/api/orders/*`, `/api/conversations`, `/api/messages/*`. See PLATFORM.md or the Node backend README for the list.

## DB

SQLite. Schema in `db/schema.sql` (embedded). Versioned migrations in `db/migrations/` (e.g. `003_add_feature.sql`); applied on startup. DB file created on first run.

## CI

GitHub Actions (`.github/workflows/ci.yml`): on push/PR to `main`/`master`, runs `go vet ./...`, `go test ./...`, and `go build` in `backend-go`.
