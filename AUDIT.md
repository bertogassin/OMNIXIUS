# OMNIXIUS — Audit by 13 risk areas

Pass over the project using the “pitfalls when developing with AI” checklist. Status: **done** / **partial** / **todo**.

---

## 1. Architecture

- **Layers:** Today handlers + DB in one place (backend-go/main.go). No separate service layer. **Partial.**
- **Clean arch:** Not applied yet; suitable as next refactor (handlers → services → repos). **Todo.**
- **DB first:** Schema exists (db/schema.sql), migrations run on init. **Done.**

---

## 2. Correctness (tests, edge cases)

- **Tests:** No unit or integration tests. **Todo.**
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
- **CORS:** `*` in code; restrict origin in production. **Todo.**

---

## 4. Scaling

- **Pagination:** Products list has `limit` (default 50, max 100) and `offset`. **Done.**
- **Cache:** No Redis yet. **Todo.**
- **Indexes:** products (user_id, category, created_at); messages (conversation_id); **orders (buyer_id, seller_id)** added. **Done.**
- **Async:** No queue; all synchronous. **Todo.**

---

## 5. Performance

- **Pagination:** In place for products. **Done.**
- **Query shape:** No N+1 in current handlers; conversations list does one extra query per row (could be batched later). **Partial.**
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
- **Structured/levels:** Plain log line; no JSON/levels. **Todo.**
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
- **Migrations:** Single schema run on init; no versioned migrations yet. **Partial.**

---

## 12. Legal (GDPR, PII, payments)

- **PII:** Email, name, messages stored; no retention or deletion flow in code. **Todo.**
- **Payments:** Not implemented. **N/A.**

---

## 13. Production (DevOps, CI/CD, resilience)

- **Docker/CI/CD/backup/DDoS:** Not in repo. **Todo.**

---

## Summary

- **Done:** Security (bcrypt, JWT, rate limits, login limit, input caps), pagination, request logging, DB indexes, input length limits.
- **Todo:** Tests, service layer, cache, async, HTTPS/CORS in prod, structured logging, migrations, PII retention, DevOps/CI/CD.

Use this file as the checklist for the next passes over the project.
