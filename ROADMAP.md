# OMNIXIUS — Roadmap

Фазы и приоритеты. Стек: **Go, Rust, Spark, Swift, Kotlin, TS/JS, Python** по направлениям (ARCHITECTURE.md).

---

## Текущая фаза: Foundation (мы здесь)

- Сайт и документация (архитектура, видение, стек, чеклист).
- Рабочее приложение: регистрация, логин, маркетплейс, почта, заказы, профиль, ИИ (корень).
- Один бэкенд (Go), контракт API, безопасность (PQC, Argon2id).
- Заготовки по всем направлениям: Rust, Spark, iOS, Android, React+TS, AI, инфра (Docker, K8s, Terraform, CI/CD, мониторинг).

**Уровень:** Близки к началу, далеки от конца. Фундамент заложен; дальше — наращивание по фазам ниже.

---

## Phase 1 — Infrastructure & core

- Инфраструктура (multi-cloud, edge, мониторинг в проде).
- IXI blockchain prototype (consensus, ZKP, economy). **Rust / C++.**

## Phase 2 — Platforms

- OMNIXIUS CONNECT: профили, geo, чат, маркетплейс (расширение). **Go** API, **Swift** iOS, **Kotlin** Android.
- OMNIXIUS TRADE & FINANCE: кошельки, биржа (crypto/FIAT), copy trading (MVP). **Go / Rust.**
- OMNIXIUS REPOSITORIUM: узлы хранения, награды, аренда compute (MVP). **Rust / C++**, **Spark** для аналитики.

## Phase 3 — Integration & scale

- Все платформы связаны через IXI и один аккаунт.
- ИИ: модерация, аналитика, свои модели.
- Мобильные приложения, глобальное масштабирование.

---

IXI — блокчейн-платформа экосистемы, не отдельный токен.
