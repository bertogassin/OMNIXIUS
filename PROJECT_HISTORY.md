# OMNIXIUS — История проекта / Project History

**Назначение:** хронология разработки, авторство, ключевые решения. Хранить в репо и обновлять при значимых изменениях.

**Repository:** https://github.com/bertogassin/OMNIXIUS

---

## Русская версия

### Начало

- Проект **OMNIXIUS** — глобальная технологическая экосистема: соцсеть, платформа специалистов, маркетплейс, обмен, децентрализованное облако, блокчейн IXI, ИИ.
- Цель: сайт, приложение и проект — **лучшие во всём мире** (безопасность, скорость, дизайн, масштаб).
- Стек зафиксирован в ARCHITECTURE.md и .cursor/rules/stack.mdc: Go (API), Rust (видео, поиск), Spark (аналитика), TypeScript/React (веб), Swift/Kotlin (мобильные).

### Основные этапы

1. **Фундамент:** backend-go (Go API), статический сайт и app/ (HTML/JS), заготовки web/ (React), services/rust, analytics/spark, mobile (Swift/Kotlin), ai (Python), infra.
2. **Обновления 8–10:** админ-панель, React-слой, единый шаблон направлений (Startups, Rewards, IXI, Repositorium).
3. **Профессии:** 24 профессии с иконками, выбор на главной, localStorage, фильтр маркетплейса, воркспейс Design (хаб, услуги, портфолио).
4. **Направления:** главная — только профессии и услуги; карта разделов и навигация по центру убраны; полная навигация — в навбаре приложения и в настройках (Main directions).
5. **Заказы и специалисты:** уведомление продавцу при заказе; страница заказа с Accept/Decline; GET /orders/:id; поиск специалистов по профессии и онлайн; профессия/локация/last_seen в users; heartbeat и last_seen в auth.
6. **WHAT-WE-TAKE:** список того, что забираем из BOLH и лучших практик (документация, рейтинг, радиус, WebSocket, audit, бэкап кошелька, чат, рекомендации Spark, мобильные, инфра).
7. **Текущий шаг:** реализация всего списка без откладывания — документация (PROJECT_HISTORY, STRUCTURE, DEVELOPMENT), audit log, машина состояний заказа, срочные заказы, WebSocket (order:status, уведомления), радиус и рейтинг в поиске специалистов, бэкап кошелька, развитие Relay по WebSocket.

### Ключевые решения

- Один активный бэкенд — **Go** (backend-go); без серверной логии на Node.
- Профессии хранятся на клиенте (localStorage) и синхронизируются в users.profession_id при заходе в приложение.
- Заказы: pending → confirmed | cancelled → completed; валидация переходов на бэкенде.
- Масштаб: готовность к PostgreSQL, Redis, горизонтальному масштабированию, Spark для рекомендаций.

---

## English version

### Start

- **OMNIXIUS** — global technology ecosystem: social network, professional platform, marketplace, exchange, decentralized cloud, IXI blockchain, AI.
- Goal: site, app, and project — **best in the world** (security, speed, design, scale).
- Stack: Go (API), Rust, Spark, TypeScript/React, Swift/Kotlin; see ARCHITECTURE.md.

### Major phases

1. Foundation: backend-go, static site and app/, web/ (React), rust/spark/mobile/ai/infra stubs.
2. Updates 8–10: admin panel, React layer, unified direction template.
3. Professions: 24 professions, choice on main page, Design workspace (hub, services, portfolio).
4. Navigation: main page focused on professions and services; full nav in app header and Settings.
5. Orders and professionals: seller notification, order detail with Accept/Decline, GET /orders/:id, professionals search (profession, online), user profession/location/last_seen.
6. WHAT-WE-TAKE: full list of adopted practices (docs, rating, radius, WebSocket, audit, wallet backup, chat, Spark, mobile, infra).
7. Current: implementing the full list — docs, audit log, order state machine, urgent orders, WebSocket, radius and rating, wallet backup, Relay WebSocket.

### Key decisions

- Single backend: Go (backend-go).
- Professions: client localStorage + server profession_id sync.
- Orders: pending → confirmed|cancelled → completed; backend validates transitions.
- Scale: PostgreSQL, Redis, horizontal scaling, Spark for recommendations.

---

*Документ обновлять при значимых изменениях. / Update this document when making significant changes.*
