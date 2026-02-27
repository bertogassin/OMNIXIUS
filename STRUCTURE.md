# OMNIXIUS — Structure

Map of folders and files.

## Root

- **README.md** — project overview, GitHub Pages deploy
- **ARCHITECTURE.md** — vision, stack, languages
- **API.md** — REST API contract
- **PROJECT_HISTORY.md** — dev history, decisions
- **STRUCTURE.md** — this file
- **DEVELOPMENT.md** — ports, commands, troubleshooting
- **WHAT-WE-TAKE.md** — adopted practices list
- **start.bat** — запуск по стеку: сначала Rust (8081), затем Go (3000)
- **start-rust.bat** — только Rust-сервис (stack 1)
- **start-backend.bat** — только Go API (port 3000)
- **web/** — **единственный фронт** (TypeScript + React). Сборка: `cd web && npm run build`. Бэкенд раздаёт `web/dist` по пути `/app` (SPA).
- **app/** — **удалён** (JS/HTML убраны до конца). Лендинг: корень (index.html); скрипты профессий и i18n в **js/** (professions.js, i18n.js, main.js).
- **backend-go/** — Go API (stack 2), migrations, orders, notifications, professionals, WebSocket; при RUST_SERVICE_URL вызывает Rust для ранжирования
- **services/rust/** — Rust-сервис (stack 1): /health, /rank; порт 8081
- **css/**, **js/** — global styles and scripts
- **analytics/spark/** — Spark (Big Data) stub
- **ai/** — Python AI service (optional)
- **mobile/ios/**, **mobile/android/** — Swift, Kotlin stubs
- **infra/** — Docker, K8s, Terraform, cache, Kafka

## backend-go/

- **main.go** — routes, handlers, WebSocket hub
- **db/** — DB connection, schema.sql, migrations/
- **order_service.go** — order create, list
- **handlers_wallet_notifications_admin.go** — wallet, notifications, admin

Update this file when adding new modules.
