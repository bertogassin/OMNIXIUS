# OMNIXIUS — Development Guide

Гайд разработчика: порты, команды, окружение, что делать дальше, типичные проблемы.

## Быстрый старт

### Запуск по стеку (Rust → Go)

```bash
# Windows — полный стек (сначала Rust, потом Go)
start.bat

# Только Go API (как раньше)
start-backend.bat

# Только Rust (порт 8081)
start-rust.bat
```

- **Порядок стека:** 1. Rust (http://localhost:8081) — поиск/ранжирование, тяжёлые задачи. 2. Go (http://localhost:3000) — API, сайт, статика.
- Поиск специалистов при `sort=rating` или `sort=distance` при запущенном Rust отдаёт ранжирование в Rust.
- Сайт: корень (index.html, css/, js/) + приложение по адресу **/app** (SPA из `web/dist`). Перед запуском бэкенда соберите фронт: `cd web && npm run build`.

### Открыть в браузере

- Главная: http://localhost:3000/
- Вход: http://localhost:3000/app/login
- Регистрация: http://localhost:3000/app/register
- Дашборд: http://localhost:3000/app
- Маркетплейс: http://localhost:3000/app/marketplace
- Заказы: http://localhost:3000/app/orders
- Найти специалиста: http://localhost:3000/app/find-professional
- Почта: http://localhost:3000/app/mail
- Кошелёк: http://localhost:3000/app/wallet
- Trade: http://localhost:3000/app/trade
- AI: http://localhost:3000/app/ai

### AI (опционально)

```bash
start-ai.bat
# или: cd ai && python -m uvicorn main:app --reload --port 8000
```

Чат на сайте может вызывать ИИ по http://localhost:8000.

## Порты

| Сервис | Порт | Описание |
|--------|------|----------|
| backend-go | 3000 | API + раздача статики (сайт и app/) |
| AI (uvicorn) | 8000 | Опционально, для чата с ИИ |

## Окружение

- **backend-go:** переменные по желанию (DB_PATH, upload dir и т.д.). По умолчанию БД: `db/omnixius.db`.
- **app/config.js:** при необходимости задаётся API URL (по умолчанию тот же хост).
- **web/ (React):** VITE_API_URL в .env при отдельном запуске фронта.

## Команды

```bash
# Бэкенд
cd backend-go && go run .
cd backend-go && go build -o omnixius-api .

# Миграции выполняются автоматически при старте (RunMigrations).

# Сборка React (если используем web/)
cd web && npm install && npm run build
```

## Что делать дальше (по WHAT-WE-TAKE)

1. **Уже сделано:** документация (PROJECT_HISTORY, STRUCTURE, DEVELOPMENT), audit log, машина состояний заказа, срочные заказы (urgent), WebSocket (GET /api/ws с token в query: order:status, notification, relay:message), радиус и рейтинг в поиске специалистов (lat, lng, radius_km, sort=rating|distance), бэкап кошелька (POST /api/wallet/export, /api/wallet/import — AES-256-GCM по паролю), Relay WebSocket (push relay:message при новом сообщении).
2. **Дорабатывать:** рекомендации Spark, нативные приложения (Swift/Kotlin), E2E чат, полная интеграция Redis/PostgreSQL при масштабе.

## Troubleshooting

### Порт 3000 занят

- Закрой другой процесс на 3000 или поменяй порт в backend-go (например, через переменную окружения или флаг).

### Логин не работает

- Убедись, что backend-go запущен на http://localhost:3000 и что запросы идут на тот же origin (CORS настроен под локальную разработку).

### Ошибки БД

- Проверь путь к БД (DB_PATH). При первой установке миграции создадут таблицы. Если схема ломалась — при необходимости переименуй/удали omnixius.db и перезапусти (создастся заново).

### Статика не грузится

- Запускай сайт через backend-go (он отдаёт файлы из корня). Не открывай index.html как file://.

---

*Обновлять при изменении портов, команд или окружения.*
