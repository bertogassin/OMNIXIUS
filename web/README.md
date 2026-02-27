# OMNIXIUS — Web (React + TypeScript)

React SPA: first useful layer for gradual migration from `app/*.html`. Vite + React + TypeScript + React Router.

**Stack:** ARCHITECTURE.md (TS/JS for web).

## Run

```bash
npm install
npm run dev
# App: http://localhost:5173
```

Backend (API) must be running separately (e.g. `go run .` in `backend-go/` on port 3000). Set API URL:

- **.env:** copy `.env.example` to `.env` and set `VITE_API_URL=http://localhost:3000`
- Or at runtime: `window.__OMNIXIUS_API_URL__ = 'http://localhost:3000'`

## Flow

- **Login** (`/login`) → sign in with email/password; token stored in sessionStorage (or localStorage if "Remember me").
- **Dashboard** (`/`) — protected; shows welcome and link to Marketplace.
- **Marketplace** (`/marketplace`) — protected; lists products from `GET /api/products`; "View" links to static app product page.

Protected routes redirect to `/login` when not authenticated. Header: Dashboard, Marketplace, Sign out.

## Build

```bash
npm run build
# Output: dist/
```

Deploy `dist/` to any static host. Ensure `VITE_API_URL` is set at build time for the production API, or set `window.__OMNIXIUS_API_URL__` before the app loads.
