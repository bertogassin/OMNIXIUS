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
- `JWT_SECRET` — change in production
- `DB_PATH` — default `db/omnixius.db`

## Endpoints

Same as before: `/api/auth/*`, `/api/users/*`, `/api/products/*`, `/api/orders/*`, `/api/conversations`, `/api/messages/*`. See PLATFORM.md or the Node backend README for the list.

## DB

SQLite. Schema in `db/schema.sql` (embedded). DB file created on first run.
