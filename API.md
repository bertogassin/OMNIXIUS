# OMNIXIUS API

REST API приложения OMNIXIUS. Реализация: **backend-go/** (Go). Base URL: `https://your-api.example.com` в проде или `http://localhost:3000` в разработке.

**Quick start:** В папке `backend-go` выполни `go run .` → API на порту 3000. На фронте задай `API_URL` (локально подставляется сам для localhost/file). Проверка: `GET /health`, затем регистрация/логин через приложение.

**Auth:** Большинство эндпоинтов требуют заголовок `Authorization: Bearer <token>`. Токен возвращают `POST /api/auth/register` и `POST /api/auth/login`.

**Errors:** Ответы — HTTP-код и тело `{"error": "message"}`.

---

## Public (no auth)

| Method | Path | Description |
|--------|------|-------------|
| GET | `/health` | Health check. 200 `{"status":"ok"}` or 503 if DB unavailable. |
| POST | `/api/auth/register` | Register. Body: `email`, `password` (8–128 chars), `name` (optional). Returns `user`, `token`. |
| POST | `/api/auth/login` | Login. Body: `email`, `password`. Returns `user`, `token`. |
| GET | `/api/auth/confirm-email?token=...` | Confirm email by token. |
| POST | `/api/auth/forgot-password` | Body: `email`. Sends reset link (or 200 always for privacy). |
| POST | `/api/auth/reset-password` | Body: `token`, `password` (min 8). Reset password with token from email. |
| GET | `/api/products` | List products. Query: `q`, `category`, `location`, `minPrice`, `maxPrice`. |
| GET | `/api/products/categories` | List category names. |
| GET | `/api/products/:id` | Get one product. 404 if not found. |

---

## Auth required

All below require header: `Authorization: Bearer <token>`.

### User

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/users/me` | Current user: `id`, `email`, `role`, `name`, `avatar_path`, `email_verified`. |
| PATCH | `/api/users/me` | Update profile. Body: `name` (optional). |
| DELETE | `/api/users/me` | Delete account and related data. |
| GET | `/api/users/me/orders` | Alias for `GET /api/orders/my`. |
| POST | `/api/users/me/avatar` | Upload avatar. Form: `avatar` (file). |

### Products

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/products` | Create product. Form: `title`, `description`, `category`, `location`, `price`, `image` (file). |
| PATCH | `/api/products/:id` | Update product (owner only). Form: same as create. |
| DELETE | `/api/products/:id` | Delete product (owner only). 204. |

### Orders

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/orders/my` | My orders. Returns `asBuyer`, `asSeller` arrays. |
| POST | `/api/orders` | Create order. Body: `product_id`. |
| PATCH | `/api/orders/:id` | Update order status (buyer/seller). Body: `status` (`pending` \| `confirmed` \| `completed` \| `cancelled`). |

### Conversations & messages

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/conversations` | List my conversations. Each: `id`, `product_id`, `updated_at`, `last_message`, `other`, `unread`. |
| GET | `/api/conversations/unread-count` | `{"unread": N}` total unread messages. |
| POST | `/api/conversations` | Create or get conversation. Body: `user_id`, `product_id` (optional). Returns `id`, `product_id`. |
| GET | `/api/messages/conversation/:id` | List messages in conversation. Participant only. |
| POST | `/api/messages/conversation/:id` | Send message. Body: `body`. |
| POST | `/api/messages/:id/read` | Mark message as read. |

---

## Error codes

| Status | Meaning |
|--------|---------|
| 400 | Bad request (validation, missing body). |
| 401 | Unauthorized (no or invalid token). |
| 403 | Forbidden (not participant/owner). |
| 404 | Not found. |
| 409 | Conflict (e.g. email already registered). |
| 429 | Too many requests (login rate limit). |
| 500 | Server error. |

---

## Env (backend)

`PORT`, `DB_PATH`, `ALLOWED_ORIGINS`, `DILITHIUM_PUBLIC_KEY`, `DILITHIUM_PRIVATE_KEY`, `ARGON2_MEMORY`. See `backend-go/.env.example`.
