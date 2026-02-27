# OMNIXIUS — что ещё сделать, чтобы было идеально

Чеклист по приоритетам: безопасность → скорость → дизайн → масштаб. Всё по стеку (Rust, Go, TS, см. STACK.md).

---

## Уже есть

- Один бэкенд Go, один фронт (web/ TypeScript), лендинг на TS, mobile грузит тот же код по URL.
- Rust-сервис (health, rank), Go его вызывает при поиске специалистов.
- CI: сборка и тесты Go на push/PR.
- .env.example в backend-go и web, PQC (Dilithium), Argon2id, rate limit, WebSocket (order:status, notifications).
- API.md, DEVELOPMENT.md, ROADMAP с фазами A/B/C.

---

## 1. CI и качество кода (быстро)

| Что | Действие | Файл |
|-----|----------|------|
| **Rust в CI** | Сборка и тесты Rust при push/PR | `.github/workflows/ci.yml` — job `rust` |
| **Frontend в CI** | `npm run build` и линт (TS) | тот же workflow — job `web` |
| **Линтер фронта** | ESLint + типы (уже есть tsc -b) | `web/package.json` — script `lint` |

---

## 2. Тесты (по мере роста)

| Где | Что | Приоритет |
|-----|-----|-----------|
| **Go** | Покрыть хендлеры auth, orders, professionals (уже есть main_test, pqc) | высокий |
| **Rust** | unit-тесты для /rank (сортировка rating/distance) | средний |
| **Web** | Vitest или React Testing Library для критичных страниц (Login, Marketplace) | средний |

---

## 3. Документация API и контракт

| Что | Действие |
|-----|----------|
| **OpenAPI/Swagger** | Сгенерировать из Go (swag, go-swagger) или вести вручную — один открытый spec в репо |
| **API.md** | Держать в актуальном состоянии при добавлении эндпоинтов (remittances, balance, installments по ROADMAP) |

---

## 4. Безопасность и продакшен

| Что | Действие |
|-----|----------|
| **Security headers** | В Go: X-Content-Type-Options, X-Frame-Options, CSP (хотя бы базовые) |
| **CORS** | В проде не `*` — только APP_URL и доверенные домены (ALLOWED_ORIGINS) |
| **Секреты** | PQC-ключи и DILITHIUM_* только из env, не в репо (уже в .env.example только комментарии) |
| **Аудит** | WHAT-WE-TAKE: логировать критические действия (заказ, смена роли) в audit_log или отдельный лог |

---

## 5. Фронт и UX

| Что | Действие |
|-----|----------|
| **Обработка ошибок** | Глобальный error boundary (React), единый формат ошибок API (ApiError) |
| **Загрузка** | Скелетоны или спиннеры на страницах с данными (Marketplace, Orders, Mail) |
| **A11y** | aria-labels у кнопок/форм, фокус при модалках (уже частично в лендинге) |
| **Profile edit** | Отдельная страница или секция «Редактировать профиль» (имя, аватар, телефон) |

---

## 6. Инфра и деплой

| Что | Действие |
|-----|----------|
| **Docker** | docker-compose поднимает Go + Rust + (опционально) AI; статика из web/dist |
| **Health** | Один endpoint типа GET /health (или /ready), который проверяет БД и опционально Rust |
| **Мониторинг** | По мере роста — метрики (Prometheus) и дашборд (Grafana), уже заготовки в infra/ |

---

## 7. По ROADMAP (Phase A/B)

| Шаг | Направление | Из ROADMAP |
|-----|-------------|------------|
| A1–A2 | AI-агенты | В чате ИИ — действия: «покажи заказы», «сводка по почте» через backend-go |
| A3–A4 | Gig | Услуги с гео, слоты, бронирование → заказ + уведомление в Mail |
| A5–A6 | Creator | Подписка (цена/месяц), закрытый контент по подписке |
| B1 | Verified | verified_email, verified_phone, бейдж «Verified» в профиле |
| B2 | Рассрочка | installment_plan у заказа, кнопка «Оформить рассрочку» (заглушка) |
| B3 | Баланс | GET /api/users/me/balance, заглушка пополнения |
| B4–B5 | Ремиты | POST/GET /api/remittances, страница «Мои переводы» |

---

## Итого

- **Сейчас сделать:** CI для Rust и web (см. п.1), при желании — security headers и линт.
- **Дальше:** тесты по приоритету, OpenAPI, пункты Phase A/B по ROADMAP.
- **Идеал:** один контракт API, всё по стеку, безопасность по умолчанию, мониторинг в проде, наращивание по фазам без переписывания.
