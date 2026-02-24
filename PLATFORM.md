# OMNIXIUS platform: Mail, Marketplace, Shop, AI, and the full ecosystem

**Уровень:** Фундамент. Один бэкенд (Go), веб-приложение в браузере, корень ИИ, инфраструктура и стек описаны. Дальше — по ROADMAP.

**Backend:** только **backend-go/** (Go). Стек по направлениям: ARCHITECTURE.md (Rust, Spark, Swift, Kotlin, Python, инфра).

## What’s live now

1. **Mail** — Internal mail: conversations and messages between users.
2. **Marketplace** — Listings: browse, search, filters, categories. Create/edit/delete listings (title, description, price, photo, category, location).
3. **Shop** — Buy and sell: orders, “Contact seller” (opens Mail), order status. Single login for marketplace, mail, profile, orders.
4. **Profile** — Registration, login, password reset, name, avatar, my orders (buyer/seller). PQC (Dilithium3) auth tokens. Account deletion: DELETE `/api/users/me` (see PRIVACY.md).
5. **AI** — Корень своего ИИ: страница в браузере (`app/ai.html`), бэкенд `ai/` (Python, `/chat`). Дальше — свои модели.

Полная экосистема (Connect, Trade & Finance, Repositorium, Blockchain и др.) — в видении и ROADMAP; новые направления подключаются по ECOSYSTEM.md.

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

- `backend-go/` — REST API в **Go** (единственный активный бэкенд). Контракт: **API.md**.
- `backend/` — Legacy Node.js (deprecated).
- `app/` — Страницы приложения (логин, регистрация, дашборд, маркетплейс, почта, заказы, ИИ).
- `web/` — Заготовка веб-клиента на React + TypeScript (Vite).
- `ai/` — ИИ: FastAPI, эндпоинт `/chat`; страница чата — `app/ai.html`.
- Сайт статический; бэкенд деплоится отдельно (OVH, Railway и т.д.). Инфраструктура: **PLATFORM_INFRASTRUCTURE.md**.

## Security & scale

- Passwords: Argon2id (legacy bcrypt at login). PQC tokens, rate limit, validation, parameterized queries. File uploads: type and size checked. Privacy: PRIVACY.md; DELETE `/api/users/me` for account deletion.
- Ready to scale: swap SQLite for PostgreSQL, add payments, geo, microservices.
