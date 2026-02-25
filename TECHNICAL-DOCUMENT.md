# Технический документ сайта OMNIXIUS

Единый технический документ по структуре, фронтенду, бэкенду и деплою всего сайта. Детали по API, архитектуре и экосистеме — в API.md, ARCHITECTURE.md, ECOSYSTEM.md.

---

## 1. Назначение и стек

**Сайт OMNIXIUS** — веб-витрина экосистемы и SPA-приложение в одном репозитории: главная, маркетплейс, почта, заказы, профиль, ИИ-чат, направления (Connect, Trade, Repositorium, Blockchain, Learning, Media и др.).

| Слой | Технологии |
|------|------------|
| **Фронтенд (сайт)** | HTML5, CSS3, JavaScript; мультиязычность (i18n); без сборки для текущих страниц. |
| **Бэкенд API** | Go (папка `backend-go/`). Один сервер: REST API + раздача статики сайта. |
| **ИИ-чат** | Python (FastAPI, папка `ai/`). Отдельный сервис на порту 8000. |
| **Деплой** | GitHub Pages (статический хостинг); при необходимости — бэкенд на отдельном хосте. |

---

## 2. Структура файлов сайта

### Корень (публичные страницы)

| Файл / папка | Назначение |
|--------------|------------|
| `index.html` | Главная: hero, список экранов, направления, о проекте, карта, футер. |
| `ecosystem.html` | Карта экосистемы: что в проде, полное видение, таблица направлений. |
| `architecture.html` | Архитектура: четыре платформы, стек. |
| `contact.html` | Контакты (email, телефон). |
| `horizon.html` | Демо-лендинг Horizon Mail (тема «Физика доверия»). |
| `404.html` | Страница «Не найдено» в стиле Horizon. |
| `css/style.css` | Основные стили сайта (переменные, header, hero, карточки, футер, медиа-запросы). |
| `css/horizon.css` | Тема Horizon (антрацит, золото, платина; кнопки, стекло, лоадер). |
| `js/main.js` | Общие скрипты (меню, i18n для главной при подключении i18n.js). |
| `CNAME` | Домен для GitHub Pages (одна строка — твой домен). |
| `start-backend.bat` | Запуск бэкенда (Go) локально. |
| `start-ai.bat` | Запуск ИИ-сервиса (Python) локально. |
| `push-to-github.bat` | Скрипт для push в репозиторий (при необходимости). |

### Папка app/ (веб-приложение)

| Файл | Назначение |
|------|------------|
| `login.html`, `register.html`, `forgot-password.html` | Вход, регистрация, восстановление пароля. |
| `dashboard.html`, `profile.html`, `profile-edit.html` | Профиль пользователя, редактирование. |
| `marketplace.html`, `product.html`, `product-create.html` | Маркетплейс: список товаров, карточка товара, создание. |
| `mail.html`, `conversation.html` | Почта: список диалогов, чат с продавцом/покупателем. |
| `orders.html` | Мои заказы (покупатель/продавец). |
| `ai.html` | ИИ-чат (интенты: заказы, сводка по почте). |
| `connect.html`, `trade.html`, `ixi.html`, `repositorium.html` | Страницы направлений (Connect, Trade, IXI, Repositorium). |
| `learning.html`, `media.html`, `rewards.html`, `startups.html` | Направления: обучение, медиа, награды, стартапы. |
| `closed-content.html` | Закрытый контент для подписок. |
| `links.html` | Служебная страница со ссылками на разделы приложения. |
| `directions.html` | Сводная страница направлений (список). |
| `test-register.html` | Тестовая страница регистрации (для разработки). |
| `config.js` | Конфиг (API URL и т.п.). |
| `js/api.js` | Клиент API (auth, products, orders, conversations, messages). |
| `js/i18n.js` | Мультиязычность: ключи EN/RU/FR/DE/ES/IT/AR/ZH/JA, `data-i18n`, `apply()`. |
| `js/utils.js`, `js/nav-unread.js` | Утилиты и счётчик непрочитанных в почте. |
| `css/app.css` | Стили приложения: формы, карточки, почта, заказы, empty/loading. |
| `css/ai.css` | Стили страницы ИИ-чата. |

### Бэкенд

| Путь | Назначение |
|------|------------|
| `backend-go/` | API на Go. Запуск: `go run .` или `start-backend.bat`. Порт по умолчанию 3000. |
| `backend-go/main.go` | Роутинг, middleware, эндпоинты, раздача статики из корня сайта (SiteRoot). |
| `backend-go/config.go` | Конфигурация (порт, БД, CORS, PQC). |
| `backend-go/auth_service.go` | Регистрация, логин, подтверждение email, сброс пароля. |
| `backend-go/product_service.go`, `order_service.go` | Товары, заказы. |
| `backend-go/conversation_service.go` | Диалоги и сообщения (почта). |
| `backend-go/balance_service.go`, `subscription_service.go`, `slot_service.go`, `remittance_service.go` | Балансы, подписки, слоты, ремиты. |
| `backend-go/db/` | Подключение к БД, schema.sql, миграции (migrations/). |
| `backend-go/pqc/` | Постквантовая криптография (токены). |
| `backend-go/.env.example`, `Dockerfile`, `README.md` | Пример переменных окружения, образ Docker, инструкции. |

Подробно по эндпоинтам — **API.md**.

### Заготовка веб-клиента (React)

| Путь | Назначение |
|------|------------|
| `web/` | Заготовка клиента на React + TypeScript. |
| `web/index.html`, `web/src/` | Точка входа и исходники (при сборке — отдельный деплой или интеграция в app). |

---

## 3. Дизайн и темы

### Основная тема (style.css, app.css)

- **Цвета:** тёмный фон (`#0a0a0d`), карточки `#14141c`, акцент бирюзовый `#00d9b4`, приглушённый текст.
- **Шрифты:** Outfit (текст), Syne (заголовки); подключение через Google Fonts в `index.html` и `@import` в `style.css`.
- **Компоненты:** header (fixed), hero с градиентом, direction-card (hover, левая полоска), кнопки .btn-primary / .btn-outline, секции .section, футер.

### Тема Horizon (horizon.css, horizon.html, 404.html)

- **Концепция:** «Физика доверия» — антрацит `#1E1E2F`, золото `#D4AF37`, платина `#A0A0A0`.
- **Шрифт:** Darker Grotesque.
- **Классы:** `.horizon`, `.horizon-glass`, `.horizon-btn`, `.horizon-input`, `.horizon-loader`, `.horizon-progress`.
- Полное описание для дизайнера — **HORIZON-DESIGN.md**.

---

## 4. Мультиязычность (i18n)

- **Файл:** `app/js/i18n.js`. Объекты по языкам (`en`, `ru`, `fr`, …), функция `t(key)`, `apply()` для обновления DOM по `data-i18n`.
- **Использование:** на страницах app подключается `js/i18n.js`; на главной и других статических страницах при наличии скрипта i18n — те же ключи.
- **Переключение:** ссылки `data-lang="en"` и т.д. в header; при клике сохраняется язык и вызывается `apply()`.

---

## 5. API и авторизация

- **Base URL:** в проде — твой API-сервер; локально — `http://localhost:3000`.
- **Авторизация:** `POST /api/auth/register`, `POST /api/auth/login` возвращают `token`. Далее заголовок `Authorization: Bearer <token>`. Токен хранится в `localStorage` (ключ `omnixius_token`).
- **Защищённые страницы:** при отсутствии токена редирект на `login.html` (проверка в скриптах страниц app).
- **Загрузки:** бэкенд отдаёт файлы по пути `/uploads` (аватары, изображения товаров); папка `backend-go/uploads` или значение из конфига. В приложении изображения подставляются через API URL + путь из ответа API.

Краткая таблица эндпоинтов — в **API.md**; полный список продуктов, заказов, почты, слотов, подписок, ремитов — там же.

---

## 6. Деплой и домен

- **GitHub Pages:** репозиторий → Settings → Pages → Source: ветка (например `main`), папка `/ (root)`. После push сайт доступен по `https://<user>.github.io/OMNIXIUS/`.
- **Свой домен:** в корне файл **CNAME** с одной строкой — домен (например `www.omnixius.com`). В настройках Pages указать Custom domain; у регистратора настроить DNS (A/CNAME по инструкции GitHub).

---

## 7. Связанные документы

| Документ | Содержание |
|----------|------------|
| **README.md** | Обзор проекта, быстрый старт, деплой, ссылки. |
| **ARCHITECTURE.md** | Четыре платформы, стек, экосистема, квантовая устойчивость. |
| **API.md** | Полное описание REST API (auth, products, orders, conversations, messages, remittances и т.д.). |
| **ECOSYSTEM.md** | Карта направлений, связи Connect/Trade/Repositorium/Blockchain. |
| **ROADMAP.md** | Фазы развития. |
| **HORIZON-DESIGN.md** | Концепция дизайна Horizon Mail для дизайнера и нейросетей. |
| **PLATFORM.md**, **PLATFORM_INFRASTRUCTURE.md** | Что в проде, инфраструктура (Docker, K8s, CI/CD). |
| **QUANTUM_READINESS.md** | Квантово-устойчивая криптография. |
| **POLNOE-TZ.md** | **Полное ТЗ в одном файле:** назначение, стек, структура, весь API, дизайн, i18n, деплой. |

---

*Документ актуален на момент последнего обновления репозитория. При изменении структуры или стека обнови этот файл и при необходимости README/API.md.*
