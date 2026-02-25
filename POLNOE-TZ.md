# Полное техническое задание (ТЗ) сайта OMNIXIUS

Один файл: назначение, стек, структура, полный API, дизайн, деплой. Углублённые детали — в ARCHITECTURE.md, ECOSYSTEM.md, ROADMAP.md, QUANTUM_READINESS.md.

---

## 1. Цель и назначение

**Сайт OMNIXIUS** — веб-витрина экосистемы и веб-приложение в одном репозитории.

- **Главная, карта, контакты** — статические страницы (index, ecosystem, architecture, contact).
- **Приложение (app/)** — логин, маркетплейс, почта, заказы, профиль, ИИ-чат, страницы направлений (Connect, Trade, Repositorium, IXI, Learning, Media, Rewards, Startups и др.).
- **Цель экосистемы:** коммуникация, заработок, инвестиции, обучение и активы в одном приложении, один аккаунт.

**Статус:** фундамент заложен; видение и стек зафиксированы; наращивание по ROADMAP.

---

## 2. Стек технологий

| Слой | Технологии |
|------|------------|
| **Фронтенд** | HTML5, CSS3, JavaScript; мультиязычность (i18n); без сборки для текущих страниц. Заготовка React+TS в `web/`. |
| **Бэкенд API** | Go (`backend-go/`). Один сервер: REST API + раздача статики сайта. Порт 3000. |
| **ИИ-чат** | Python, FastAPI (`ai/`). Отдельный сервис, порт 8000. |
| **Деплой сайта** | GitHub Pages. Бэкенд — отдельный хост при необходимости. |

**Общий стек по направлениям (видение):** сервер — Go, Rust, Spark; веб — TS/JS, React/Vue; мобилки — Swift, Kotlin; ИИ — Python; инфра — Bash/Python, YAML. Детально: ARCHITECTURE.md.

---

## 3. Четыре платформы (ядро)

| Платформа | Кратко |
|-----------|--------|
| **CONNECT** | Соцсеть, эксперты, маркетплейс, чат, видео, подписки. |
| **TRADE & FINANCE** | Биржа, кошельки, копи-трейдинг, платежи. |
| **REPOSITORIUM** | Децентрализованное хранилище и compute, IPFS, награды. |
| **BLOCKCHAIN (IXI)** | Консенсус, ZKP, токены, стейкинг, buyback. |

Остальные направления (Learning, Media, Rewards, Startups, Health, SSI и т.д.) подключаются к этим платформам. Карта: ECOSYSTEM.md.

---

## 4. Структура файлов

### Корень

| Файл / папка | Назначение |
|--------------|------------|
| `index.html` | Главная: hero, экраны, направления, о проекте, карта, футер. |
| `ecosystem.html` | Карта экосистемы. |
| `architecture.html` | Архитектура, четыре платформы, стек. |
| `contact.html` | Контакты. |
| `horizon.html` | Демо-лендинг Horizon Mail. |
| `404.html` | Страница 404 в стиле Horizon. |
| `css/style.css` | Основные стили. |
| `css/horizon.css` | Тема Horizon. |
| `js/main.js` | Меню, общие скрипты. |
| `CNAME` | Домен для GitHub Pages (одна строка). |
| `start-backend.bat` | Запуск бэкенда. |
| `start-ai.bat` | Запуск ИИ. |
| `push-to-github.bat` | Push в репозиторий. |

### app/ (веб-приложение)

| Файл | Назначение |
|------|------------|
| `login.html`, `register.html`, `forgot-password.html` | Вход, регистрация, восстановление пароля. |
| `dashboard.html`, `profile.html`, `profile-edit.html` | Профиль. |
| `marketplace.html`, `product.html`, `product-create.html` | Маркетплейс. |
| `mail.html`, `conversation.html` | Почта. |
| `orders.html` | Заказы. |
| `ai.html` | ИИ-чат. |
| `connect.html`, `trade.html`, `ixi.html`, `repositorium.html` | Страницы направлений. |
| `learning.html`, `media.html`, `rewards.html`, `startups.html` | Обучение, медиа, награды, стартапы. |
| `closed-content.html` | Закрытый контент для подписок. |
| `links.html`, `directions.html`, `test-register.html` | Ссылки, направления, тест регистрации. |
| `config.js` | Конфиг (API URL). |
| `js/api.js` | Клиент API. |
| `js/i18n.js` | Мультиязычность. |
| `js/utils.js`, `js/nav-unread.js` | Утилиты, счётчик непрочитанных. |
| `css/app.css`, `css/ai.css` | Стили приложения и ИИ. |

### backend-go/

| Путь | Назначение |
|------|------------|
| `main.go` | Роутинг, middleware, раздача статики (SiteRoot), `/uploads`. |
| `config.go` | Конфигурация. |
| `auth_service.go` | Регистрация, логин, confirm-email, forgot/reset password. |
| `product_service.go`, `order_service.go` | Товары, заказы. |
| `conversation_service.go` | Диалоги и сообщения. |
| `balance_service.go`, `subscription_service.go`, `slot_service.go`, `remittance_service.go` | Балансы, подписки, слоты, ремиты. |
| `db/` | БД, schema.sql, migrations/. |
| `pqc/` | Постквантовая криптография (токены). |
| `.env.example`, `Dockerfile`, `README.md` | Переменные окружения, Docker, инструкции. |

### web/

Заготовка клиента на React + TypeScript (`web/index.html`, `web/src/`).

---

## 5. REST API (полный список)

**Base URL:** в проде — твой API; локально — `http://localhost:3000`.  
**Авторизация:** заголовок `Authorization: Bearer <token>`. Токен дают `POST /api/auth/register` и `POST /api/auth/login`.  
**Ошибки:** тело `{"error": "message"}` и соответствующий HTTP-код.

### Публичные (без auth)

| Method | Path | Описание |
|--------|------|----------|
| GET | `/health` | Health check. 200 `{"status":"ok"}` или 503. |
| POST | `/api/auth/register` | Регистрация. Body: `email`, `password` (8–128), `name` (опц.). Возврат: `user`, `token`. |
| POST | `/api/auth/login` | Вход. Body: `email`, `password`. Возврат: `user`, `token`. |
| GET | `/api/auth/confirm-email?token=...` | Подтверждение email. |
| POST | `/api/auth/forgot-password` | Body: `email`. Отправка ссылки сброса (или всегда 200 ради приватности). |
| POST | `/api/auth/reset-password` | Body: `token`, `password` (мин. 8). Сброс пароля по токену из письма. |
| GET | `/api/products` | Список товаров. Query: `q`, `category`, `location`, `minPrice`, `maxPrice`, `service`, `subscription`, `user_id`. |
| GET | `/api/products/categories` | Список категорий. |
| GET | `/api/products/:id` | Один товар. 404 если нет. |
| GET | `/api/products/:id/slots` | Слоты товара. С auth владелец видит все, остальные — только свободные. |
| GET | `/api/users/:id` | Публичный профиль: `id`, `name`, `avatar_path`, `verified`. |

### С авторизацией (Bearer token)

**Пользователь**

| Method | Path | Описание |
|--------|------|----------|
| GET | `/api/users/me` | Текущий пользователь. |
| PATCH | `/api/users/me` | Обновить профиль. Body: `name` (опц.). |
| DELETE | `/api/users/me` | Удалить аккаунт и данные. |
| GET | `/api/users/me/orders` | Мои заказы: `asBuyer`, `asSeller`. |
| GET | `/api/users/me/balance` | Баланс. `{ balance }`. Заглушка Trade. |
| POST | `/api/users/me/balance/credit` | Тестовое пополнение. Body: `amount`. |
| POST | `/api/users/me/avatar` | Загрузка аватара. Form: `avatar` (файл). |

**Товары**

| Method | Path | Описание |
|--------|------|----------|
| POST | `/api/products` | Создать. Form: `title`, `description`, `category`, `location`, `price`, `image`, `is_service`, `is_subscription` (0/1), `closed_content_url` (опц.). |
| PATCH | `/api/products/:id` | Обновить (владелец). |
| GET | `/api/products/:id/closed-content` | Закрытый контент подписки. 403 если нужно подписаться. |
| DELETE | `/api/products/:id` | Удалить (владелец). 204. |
| POST | `/api/products/:id/slots` | Добавить слот (владелец). Body: `slot_at` (Unix). |
| POST | `/api/products/:id/slots/:sid/book` | Забронировать слот. Создаёт заказ, сообщение продавцу. |
| POST | `/api/subscriptions` | Подписаться. Body: `product_id`. Заглушка. |
| GET | `/api/subscriptions/my` | Мои подписки. |

**Заказы**

| Method | Path | Описание |
|--------|------|----------|
| GET | `/api/orders/my` | Мои заказы (плоский список). Есть `installment_plan`. |
| POST | `/api/orders` | Создать. Body: `product_id`, опц. `installment_plan`: `"requested"`. |
| PATCH | `/api/orders/:id` | Обновить. Body: `status` (`pending` \| `confirmed` \| `completed` \| `cancelled`), опц. `installment_plan`. |

**Ремиты (заглушка)**

| Method | Path | Описание |
|--------|------|----------|
| GET | `/api/remittances/my` | Мои запросы на перевод. |
| POST | `/api/remittances` | Создать. Body: `to_identifier`, `amount`, `currency` (опц., по умолч. USD). |

**Диалоги и сообщения (почта)**

| Method | Path | Описание |
|--------|------|----------|
| GET | `/api/conversations` | Список диалогов. Поля: `id`, `product_id`, `product_title`, `updated_at`, `last_message`, `other`, `unread`. |
| GET | `/api/conversations/unread-count` | `{"unread": N}`. |
| GET | `/api/conversations/:id` | Мета одного диалога: `other`, `product_id`, `product_title`. |
| POST | `/api/conversations` | Создать/получить. Body: `user_id`, `product_id` (опц.). |
| GET | `/api/messages/conversation/:id` | Сообщения в диалоге. |
| POST | `/api/messages/conversation/:id` | Отправить. Body: `body`. |
| POST | `/api/messages/:id/read` | Отметить прочитанным. |

### Коды ответов

| Код | Значение |
|-----|----------|
| 400 | Bad request (валидация, тело). |
| 401 | Unauthorized (нет/неверный token). |
| 403 | Forbidden (не участник/владелец). |
| 404 | Not found. |
| 409 | Conflict (например email уже занят). |
| 429 | Too many requests (лимит входа). |
| 500 | Ошибка сервера. |

### Загрузки

Бэкенд отдаёт файлы по `/uploads` (аватары, картинки товаров). Каталог задаётся в конфиге (например `backend-go/uploads`). В приложении URL картинок: base API URL + путь из ответа API.

---

## 6. ИИ-чат (ai/)

- **Фронт:** `app/ai.html`. В запрос к ИИ передаются `api_token` и `api_url`.
- **Бэкенд:** Python, FastAPI, порт 8000. Запуск: `start-ai.bat` или в `ai/`: `python -m uvicorn main:app --reload --port 8000`.
- **Интенты:** «Мои заказы» → `GET /api/orders/my`; «Создать заказ» (напр. «order product 5») → `POST /api/orders`; «Сводка по письмам» → `GET /api/conversations`. Полный каталог: **ai/README.md**.

---

## 7. Переменные окружения (бэкенд)

`PORT`, `DB_PATH`, `ALLOWED_ORIGINS`, `DILITHIUM_PUBLIC_KEY`, `DILITHIUM_PRIVATE_KEY`, `ARGON2_MEMORY`. Пример: `backend-go/.env.example`.

---

## 8. Дизайн

### Основная тема (style.css, app.css)

- Цвета: фон `#0a0a0d`, карточки `#14141c`, акцент `#00d9b4`, приглушённый текст.
- Шрифты: Outfit (текст), Syne (заголовки).
- Компоненты: фиксированный header, hero с градиентом, direction-card с hover, кнопки .btn-primary / .btn-outline, секции, футер.

### Тема Horizon (horizon.css)

- Антрацит `#1E1E2F`, золото `#D4AF37`, платина `#A0A0A0`. Шрифт Darker Grotesque.
- Классы: .horizon, .horizon-glass, .horizon-btn, .horizon-input, .horizon-loader, .horizon-progress.
- Полная концепция для дизайнера: **HORIZON-DESIGN.md**.

---

## 9. Мультиязычность (i18n)

- Файл: `app/js/i18n.js`. Языки: EN, RU, FR, DE, ES, IT, AR, ZH, JA.
- Функция `t(key)`, `apply()` для обновления DOM по атрибуту `data-i18n`.
- Переключение: ссылки `data-lang="en"` и т.д. в header; язык сохраняется, вызывается `apply()`.

---

## 10. Авторизация на фронте

- Токен в `localStorage`: ключ `omnixius_token`.
- На защищённых страницах app при отсутствии токена — редирект на `login.html`.

---

## 11. Деплой и домен

- **GitHub Pages:** Settings → Pages → Source: ветка (напр. `main`), папка root. Сайт: `https://<user>.github.io/OMNIXIUS/`.
- **Свой домен:** в корне файл **CNAME** с одной строкой (домен). В настройках Pages — Custom domain; у регистратора — DNS по инструкции GitHub.

---

## 12. Связанные документы

| Документ | Содержание |
|----------|------------|
| **README.md** | Обзор, быстрый старт, как выложить на GitHub. |
| **ARCHITECTURE.md** | Четыре платформы, стек по направлениям, 10 направлений, фазы. |
| **API.md** | Текстовая версия REST API (дублирует раздел 5 этого ТЗ). |
| **ECOSYSTEM.md** | Карта направлений, связи платформ. |
| **ROADMAP.md** | Фазы и приоритеты. |
| **HORIZON-DESIGN.md** | Концепция дизайна Horizon Mail. |
| **QUANTUM_READINESS.md** | Квантово-устойчивая криптография. |
| **PLATFORM.md**, **PLATFORM_INFRASTRUCTURE.md** | Прод, инфраструктура (Docker, K8s, CI/CD). |
| **TECHNICAL-DOCUMENT.md** | Краткий технический документ (структура + ссылки на API/архитектуру). |
| **ARCHITECTURE-V4.md** | Полная архитектура платформы v4.0 (ядро, модули, vault, инфра, безопасность, дизайн). |
| **IMPLEMENTATION-V4.md** | План реализации по v4: фазы и первые шаги. |

---

*Полное ТЗ в одном файле. При изменении API или структуры — обновить этот документ и при необходимости API.md.*
