# OMNIXIUS — Audit by 13 risk areas

Pass over the project using the “pitfalls when developing with AI” checklist. Status: **done** / **partial** / **todo**.

---

## 1. Architecture

- **Layers:** Auth, products, orders (order_service.go), conversations/messages (conversation_service.go); handlers validate and call service. **Done.**
- **Clean arch:** Same pattern everywhere. **Done.**
- **DB first:** Schema exists (db/schema.sql), migrations run on init. **Done.**

---

## 2. Correctness (tests, edge cases)

- **Tests:** pqc (Sign/Verify, expiry, tamper); auth (login 401, register 201 + user/token); product get 404. **Partial.**
- **Edge cases:** Input length caps added (email, name, product title/description, message body). **Done.**
- **Docs:** README and PLATFORM describe API at a high level. **Partial.**

---

## 3. Security

- **Passwords:** Argon2id (quantum-resistant KDF); legacy bcrypt supported at login. **Done.**
- **Auth tokens:** Dilithium3 (post-quantum signatures), not JWT/HMAC. **Done.**
- **Rate limit:** Global 200 req/15 min; **per-IP login limit 5 per 15 min** added. **Done.**
- **Input:** Validation and length limits on register, product create, message send. **Done.**
- **SQL:** Parameterized queries only. **Done.**
- **HTTPS:** Not enforced in code; must be set in production (reverse proxy). **Todo.**
- **CORS:** Origin from env `ALLOWED_ORIGINS` (comma-separated); empty = `*` for dev. **Partial.**

---

## 4. Scaling

- **Pagination:** Products list has `limit` (default 50, max 100) and `offset`. **Done.**
- **Cache:** No Redis yet. **Todo.**
- **Indexes:** products (user_id, category, created_at); messages (conversation_id); **orders (buyer_id, seller_id)** added. **Done.**
- **Async:** No queue; all synchronous. **Todo.**

---

## 5. Performance

- **Pagination:** In place for products. **Done.**
- **Query shape:** Conversations list batched (one query for “other” user per conv, one for last message). **Done.**
- **Logging:** Request logger added (method, path, status, IP, latency). **Done.**

---

## 6. Dependencies / versions

- **Go:** go.mod, go.sum; versions pinned. **Done.**
- **Upgrades:** No automated dependency checks. **Todo.**

---

## 7. Business logic

- **Model:** Clear enough for MVP (users, products, orders, conversations, messages). **Done.**
- **Strategy:** Not in code; product/legal decisions documented elsewhere. **N/A.**

---

## 8. Logging

- **Request log:** Method, status, path, client IP, latency. **Done.**
- **Structured/levels:** One JSON line per request (level, method, path, status, ip, latency_ms). **Done.**
- **Secrets/PII:** Not logged in handlers. **Done.**
- **Alerts/monitoring:** None. **Todo.**

---

## 9. Languages / stack

- **Stack:** Go only for API; stack doc (ARCHITECTURE.md) limits to Spark, Rust, C++, Swift, Go. **Done.**

---

## 10. AI dependency

- **Process:** Human review and ownership of generated code; not automated in repo. **N/A.**

---

## 11. Database

- **Indexes:** products, messages, orders (buyer_id, seller_id). **Done.**
- **Normalization:** Tables and FKs in place. **Done.**
- **Migrations:** schema_version table + embedded migrations (db/migrations/NNN_name.sql); runner in db.RunMigrations(). **Partial.**

---

## 12. Legal (GDPR, PII, payments)

- **PII:** PRIVACY.md documents stored PII, retention (until account deletion), rights (access/correction/deletion). DELETE `/api/users/me` deletes account and related data (orders, then user; products/conversations/messages via CASCADE). **Partial.**
- **Payments:** Not implemented. **N/A.**

---

## 13. Production (DevOps, CI/CD, resilience)

- **Docker/CI/CD/backup/DDoS:** Dockerfile for backend-go; GitHub Actions CI (build + test on push/PR). **Partial.**

---

## Summary

- **Done:** Security (Argon2id, Dilithium3, rate limits, login limit, input caps), pagination, structured logging, DB indexes, CORS from env, batched conversations, Dockerfile, CI, versioned migrations, **full service layer** (auth, products, orders, conversations/messages), tests (pqc, auth, register, product 404). **Bugfix:** PQC key assignment in config (priv/pub were swapped).
- **Todo:** Cache, async, HTTPS and strict CORS in prod, backup/DDoS, more tests.

Use this file as the checklist for the next passes over the project.
