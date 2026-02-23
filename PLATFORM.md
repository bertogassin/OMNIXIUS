# OMNIXIUS platform: Mail, Marketplace, Shop, and the full ecosystem

**Stack: Go only** for backend (see ARCHITECTURE.md — SPARK, RUST, C++, SWIFT, GO).

## What’s live now

1. **Mail** — Internal mail: conversations and messages between users.
2. **Marketplace** — Listings: browse, search, filters, categories. Create/edit/delete listings (title, description, price, photo, category, location).
3. **Shop** — Buy and sell: orders, “Contact seller” (opens Mail), order status. Single login for marketplace, mail, profile, orders.
4. **Profile** — Registration, login, password reset, name, avatar, my orders (buyer/seller). PQC (Dilithium3) auth tokens. Account deletion: DELETE `/api/users/me` (see PRIVACY.md).

The full ecosystem (Connect, Trade & Finance, Repositorium, Blockchain, and more) is in vision and roadmap; new directions are added over time. See **ECOSYSTEM.md**.

## How to run

1. **Backend (Go)**  
   ```bash
   cd backend-go
   go build -o omnixius-api . && ./omnixius-api
   ```  
   API: http://localhost:3000

2. **Front** — Open the site (e.g. `index.html` or a local server). In `app/config.js` set `window.API_URL = 'http://localhost:3000'`.

3. **Check** — Home → “Sign in” → Register/Login → Dashboard, Marketplace, Mail.

## Structure

- `backend-go/` — REST API in **Go** (official backend).
- `backend/` — Legacy Node.js API (deprecated; use backend-go).
- `app/` — App pages (login, register, dashboard, marketplace, mail).
- Site is static; deploy backend separately (OVH, Railway, etc.).

## Security & scale

- Passwords: Argon2id (legacy bcrypt at login). PQC tokens, rate limit, validation, parameterized queries. File uploads: type and size checked. Privacy: PRIVACY.md; DELETE `/api/users/me` for account deletion.
- Ready to scale: swap SQLite for PostgreSQL, add payments, geo, microservices.
