# OMNIXIUS — проверка: всё в проекте

Краткая сводка: приоритеты, направления, браузер, ИИ, Platform & Infrastructure. Каждый пункт привязан к папке/файлу в репозитории. **Уровень:** начало; по этому чеклисту видно, что заложено и куда наращивать.

---

## Приоритет 1 (ядро системы)

| Стек | Назначение | В репозитории |
|------|------------|----------------|
| **Rust** | Видео, поиск, высоконагруженные движки | `services/rust/` (Cargo, axum, `/health`), Dockerfile |
| **Go** | API, чаты, микросервисы, авторизация | `backend-go/` (Gin, auth, products, orders, conversations), Dockerfile |
| **Spark** (Scala/Java) | Аналитика, рекомендации, Big Data | `analytics/spark/` (sbt, RecommendationsJob) |

---

## Приоритет 2 (клиенты)

| Стек | Назначение | В репозитории |
|------|------------|----------------|
| **Swift** | Нативное iOS | `mobile/ios/` (Package.swift, App.swift) |
| **Kotlin** | Нативное Android | `mobile/android/` (Gradle, MainActivity) |
| **TypeScript/JavaScript** | Веб-версия (сайт) | `app/` (HTML/JS), `web/` (Vite + React + TS) |

---

## Приоритет 3 (инфраструктура)

| Стек | Назначение | В репозитории |
|------|------------|----------------|
| **Python / Bash / YAML** | Автоматизация, скрипты, DevOps | `infra/scripts/` (deploy.sh, health_check.py), `infra/config/app.yaml`, `ai/` (Python) |

---

## По направлениям (дублируют приоритеты + игры, ИИ)

- **Сервер:** Go, Rust, Spark — см. выше.
- **Веб-сайт:** TS/JS, React/Vue — `app/`, `web/`.
- **Мобилки:** Swift, Kotlin — `mobile/ios/`, `mobile/android/`.
- **Игры (если будут):** C++ (Unreal), C# (Unity) — в репо только документированы в ARCHITECTURE.md §2.
- **ИИ:** Python — `ai/` (FastAPI, `/chat`).
- **Инфра:** Bash/Python, YAML — `infra/`.

---

## Браузер OMNIXIUS

- **Сайт и приложение в браузере:** `index.html`, `app/*.html` (логин, регистрация, маркетплейс, почта, заказы, профиль, настройки, чаты).
- **Заготовка веб-клиента (React + TypeScript):** `web/` (Vite, React, App.tsx).

---

## ИИ OMNIXIUS

- **В браузере:** `app/ai.html` — страница чата с ИИ (корень своего ИИ).
- **Бэкенд:** `ai/` — Python, FastAPI, эндпоинты `/health`, `/chat`. Дальше — свои модели.

---

## Platform & Infrastructure Division

| Элемент | В репозитории |
|---------|----------------|
| **Контейнеризация (Docker)** | `backend-go/Dockerfile`, `services/rust/Dockerfile`, `ai/Dockerfile`, `infra/docker-compose.yml` |
| **Оркестрация (Kubernetes)** | `infra/k8s/` (namespace.yaml, README) |
| **Infrastructure as Code (Terraform)** | `infra/terraform/` (main.tf, README) |
| **CI/CD (GitHub Actions)** | `.github/workflows/ci.yml` (сборка и тесты Go API) |
| **Мониторинг (Prometheus, Grafana)** | `infra/docker-compose.monitoring.yml`, `infra/monitoring/prometheus.yml` |
| **Потоки данных (Apache Kafka)** | `infra/kafka/README.md` (заготовка под масштаб) |
| **Кэш (Redis)** | `infra/cache/README.md` (заготовка под масштаб) |

Подробно: **PLATFORM_INFRASTRUCTURE.md**.

---

## Итого

- **Приоритеты 1–3** и **направления** — описаны в ARCHITECTURE.md §2 и реализованы в перечисленных папках.
- **Браузер OMNIXIUS** — `index.html`, `app/`, `web/`.
- **ИИ OMNIXIUS** — `app/ai.html` + `ai/`.
- **Platform & Infrastructure** — Docker, K8s, Terraform, CI/CD, Prometheus, Grafana, Kafka отражены в PLATFORM_INFRASTRUCTURE.md и в папке `infra/`.

---

## Дубликаты и архитектура (проверено)

| Вопрос | Ответ |
|--------|--------|
| **Бэкенд** | Один: **backend-go/** (Go). Папка **backend/** — только README (legacy удалён). |
| **Общая логика в app/** | **Один раз:** `app/js/utils.js` — escapeHtml, formatDate. Подключён в marketplace, product, mail, orders, conversation, ai. Дубликатов в страницах нет. |
| **API, конфиг** | Один источник правды: **API.md** в корне. Конфиг фронта: **app/config.js** (API_URL, AI_URL). |
| **Хедер/нав в каждом HTML** | Для статического сайта без сборки — норма: каждая страница содержит свой header/nav. Не перегруз. |
| **Dockerfile** | Три отдельных (backend-go, services/rust, ai) — по одному на сервис, дубликатов нет. |

Итог: дубликатов кода нет, архитектура без сюрпризов, один бэкенд (Go), общие утилиты в одном файле.

---

Если чего-то не хватает в репо — добавь папку/файл по этой таблице.
