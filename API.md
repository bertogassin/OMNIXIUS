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
| GET | `/api/products` | List products. Query: `q`, `category`, `location`, `minPrice`, `maxPrice`, `service`, `subscription`, `user_id`. |
| GET | `/api/products/categories` | List category names. |
| GET | `/api/products/:id` | Get one product. 404 if not found. |
| GET | `/api/products/:id/slots` | List slots for product. Optional auth: owner sees all, others see only free. Returns `[{id, product_id, slot_at, status, order_id, created_at}]`. |
| GET | `/api/users/:id` | Public user profile. Returns `id`, `name`, `avatar_path`, `verified` (true if email or phone verified). |

---

## Auth required

All below require header: `Authorization: Bearer <token>`.

### User

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/users/me` | Current user: `id`, `email`, `role`, `name`, `avatar_path`, `email_verified`, `phone_verified`, `verified` (true if either verified). |
| PATCH | `/api/users/me` | Update profile. Body: `name` (optional). |
| DELETE | `/api/users/me` | Delete account and related data. |
| GET | `/api/users/me/orders` | Alias for `GET /api/orders/my`. |
| POST | `/api/users/me/avatar` | Upload avatar. Form: `avatar` (file). |

### Products

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/products` | Create product. Form: `title`, `description`, `category`, `location`, `price`, `image` (file), `is_service`, `is_subscription` (0/1), `closed_content_url` (optional). |
| PATCH | `/api/products/:id` | Update product (owner only). Form: same as create; optional `closed_content_url`. |
| GET | `/api/products/:id/closed-content` | **Auth.** Closed content for subscription listing. Returns `{ url }` if user is subscriber or owner; 403 if must subscribe; 400 if not a subscription or no URL. |
| DELETE | `/api/products/:id` | Delete product (owner only). 204. |
| POST | `/api/products/:id/slots` | Add slot (owner only). Body: `slot_at` (Unix timestamp). For service listings. |
| POST | `/api/products/:id/slots/:sid/book` | Book slot (auth). Creates order, marks slot booked, sends message to seller in Mail. Only for service listings. Returns `{order, slot_id, message}`. |
| POST | `/api/subscriptions` | Subscribe to a subscription listing (stub; payments later). Body: `product_id`. Product must have `is_subscription=1`. Returns `{id, product_id, user_id, status}`. |
| GET | `/api/subscriptions/my` | My active subscriptions. Returns array with `product_id`, `title`, `price`, `seller_name`, etc. |

### Orders

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/orders/my` | My orders (flat list). Each order includes `installment_plan` (`""` or `"requested"`). |
| GET | `/api/users/me/orders` | My orders grouped as `asBuyer`, `asSeller`. Each item includes `id`, `status`, `created_at`, `installment_plan`, `title`, `price`, `image_path`, `seller_name` (or buyer name in asSeller). |
| POST | `/api/orders` | Create order. Body: `product_id` (required), optional `installment_plan`: `"requested"` to request installments at creation. |
| PATCH | `/api/orders/:id` | Update order (buyer or seller). Body: optional `status` (`pending` \| `confirmed` \| `completed` \| `cancelled`), optional `installment_plan`: `"requested"` to record installments request. |

**Embedded finance (B2) — Installments stub:** Buyer can request installments for an order: set `installment_plan` to `"requested"` on create (POST) or later (PATCH). The value is stored and returned in order lists; actual installment flow (pay in parts) will be implemented via Trade later. UI: "Request installments" on order card → PATCH with `installment_plan: "requested"` → show "Installments requested (coming via Trade)".

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

## AI agent intents (app/ai.html → ai service)

The AI chat can perform actions on behalf of the signed-in user when the frontend sends `api_token` and `api_url` in the `/chat` request body. Supported intents (catalog in **ai/README.md**):

- **My orders** — `GET /api/orders/my`
- **Create order** — `POST /api/orders` with `product_id` (parsed from message, e.g. "order product 5")
- **Conversations summary** — `GET /api/conversations`

---

## Env (backend)

`PORT`, `DB_PATH`, `ALLOWED_ORIGINS`, `DILITHIUM_PUBLIC_KEY`, `DILITHIUM_PRIVATE_KEY`, `ARGON2_MEMORY`. See `backend-go/.env.example`.
