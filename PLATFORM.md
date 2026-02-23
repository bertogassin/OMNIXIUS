# OMNIXIUS platform: registration, marketplace, mail

**Stack: Go only** for backend (see ARCHITECTURE.md — SPARK, RUST, C++, SWIFT, GO).

## What’s included

1. **Registration & accounts** — Email/password, confirm email, login/logout, password reset. Roles: user, seller, admin. Profile: name, avatar, order history. JWT auth.
2. **Marketplace** — Create/edit/delete listings (title, description, price, photo, category, location). Search and filters. Product page, “Contact seller” (opens internal mail). Orders.
3. **Internal mail** — Conversations and messages, optional product link, read status.
4. **Integration** — Single login for marketplace, mail, profile, orders. “Contact seller” creates a conversation with optional product link.

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

- Passwords: bcrypt. JWT, rate limit, validation, parameterized queries. File uploads: type and size checked.
- Ready to scale: swap SQLite for PostgreSQL, add payments, geo, microservices.
