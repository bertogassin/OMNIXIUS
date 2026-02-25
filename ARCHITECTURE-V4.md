# ПОЛНАЯ АРХИТЕКТУРА ПЛАТФОРМЫ OMNIXIUS v4.0

**Тип документа:** Исполнительная архитектура  
**Назначение:** Единый источник истины для разработки  
**Статус:** Актуально, включает все технические уточнения  

---

## Раздел 0: СТРУКТУРА ДОКУМЕНТА

- 0. СТРУКТУРА ДОКУМЕНТА
- 1. ЯДРО ПЛАТФОРМЫ (auth, ключи, user, crypto, storage, event bus, db, SDK, аудит)
- 2. МОДУЛЬНАЯ АРХИТЕКТУРА (структура модуля, манифест, изоляция схем, scopes, список модулей, API-шаблон, установка)
- 3. ТЕКУЩИЙ ПРИОРИТЕТ (модуль vault: файлы, папки, поиск, pre-signed URLs)
- 4. ИНФРАСТРУКТУРА (компоненты, сеть, API, PostgreSQL, S3, CDN, кэш, очереди)
- 5. БЕЗОПАСНОСТЬ (принципы, шифрование, ключи, восстановление, метаданные, Blind Indexing, Passkeys, защита API, клиент)
- 6. ДИЗАЙН-СИСТЕМА (константы, компоненты Button/Input, тема Horizon, адаптивность, a11y, i18n)
- 7. ДАННЫЕ (типы, миграции, бэкапы, репликация, retention)
- 8–14. Коммуникации, разработка, развёртывание, мониторинг, документация, версионирование, приложения (API, модели, коды ошибок)

*Полный текст разделов 1–7 и фрагментов далее сохранён в репозитории. Ниже — выжимка и ссылки на реализацию.*

---

## 1. ЯДРО ПЛАТФОРМЫ — кратко

- **auth/** — только Passkeys (WebAuthn), несколько устройств, восстановление через Master Recovery Key (BIP-39, 24 слова). API: `/auth/register/begin|complete`, `/auth/login/begin|complete`, `/auth/sessions`, `/auth/devices`, `/auth/recovery/*`.
- **crypto/keys** — иерархия: Master Recovery Key → User Master Key (UMK) → Device Keys, Encryption Keys, Signing Keys. UMK хранится зашифрованным на сервере.
- **user/** — id, email (только для связи), display_name, avatar, settings, encrypted_umk. API: `/user/me`, `/user/avatar`, `/user/devices`.
- **crypto/** — интерфейс CryptoProvider (Encrypt/Decrypt, ключи, подписи, хэш). Алгоритмы v1: AES-256-GCM, Ed25519, SHA-256; подготовка к PQC (Kyber, Dilithium).
- **storage/** — абстракция StorageProvider (Put/Get/Delete/List, GenerateUploadURL/GenerateDownloadURL). Реализации: S3, Local, IPFS.
- **event/** — шина событий (user.*, module.*, data.*, security.*). EventBus: Publish, Subscribe.
- **db/** — PostgreSQL 15+, goose/golang-migrate, sqlc/pgxpool. Схема core: users, sessions, devices, modules, events.
- **Internal SDK** — AuthMiddleware, RateLimitMiddleware, LimitBodyMiddleware, LoggerMiddleware, RespondJSON/RespondError, Validator.
- **Аудит и логи** — AuditLog (UserID, Action, Resource, OldValue/NewValue, IP, UserAgent); структурированные логи (level, time, request_id, method, path, user_id, status, duration_ms).

---

## 2. МОДУЛЬНАЯ АРХИТЕКТУРА — кратко

- Структура модуля: `api/` (handler, routes, middleware), `service/`, `db/migrations/`, `crypto/`, `events/`, `manifest.json`.
- Манифест: name, version, api.basePath, api.scopes, dependencies, permissions, config, database.schema.
- Каждый модуль — отдельная схема PostgreSQL (например vault, legacy, connect). Без прямых JOIN между схемами.
- Scopes: vault:read|write|delete|share, connect:read|write, trade:read|write и т.д.
- Список модулей: core, vault, legacy, connect, horizon, trade, wallet, ixi, dao, repositorium, learning, media, rewards, startups, health, ssi.
- Установка: проверка зависимостей → миграции → запись в core.modules → регистрация API → событие module.installed.

---

## 3. ТЕКУЩИЙ ПРИОРИТЕТ — vault

- **vault** — цифровой сейф: files, folders, search (blind indexing), sharing (опционально), pre-signed URLs для загрузки/скачивания.
- Схема: vault.files (id, user_id, encrypted_name, encrypted_metadata, search_index, storage_path, size_bytes, folder_id, created_at, updated_at, deleted_at), vault.folders (id, user_id, encrypted_name, parent_id).
- Большие файлы: клиент запрашивает `POST /api/v1/vault/files/upload-url` (name, size, mimeType, folderId) → получает uploadUrl, fileId, expiresAt → загружает в S3 → `POST /api/v1/vault/files/:id/complete`. Скачивание: запрос download URL → GetDownloadURL.

---

## 4. ИНФРАСТРУКТУРА — кратко

- Компоненты: CDN → Load Balancer → API Servers (Go) → PostgreSQL (Primary + Replica), Redis, S3.
- DNS: omnixius.com, api.omnixius.com, cdn.omnixius.com. SSL/TLS 1.3, HSTS.
- API: Go 1.21+, один бинарник, graceful shutdown.
- PostgreSQL: connection pool (pgxpool), max_connections, shared_buffers и т.д.
- S3: структура users/{user_id}/files/, users/{user_id}/avatars/, temp/, system/. CORS для прямых загрузок.
- CDN: Cloudflare, Brotli, cache. Redis: сессии, кэш, rate limit. Очереди: фоновые задачи (email, превью, очистка).

---

## 5. БЕЗОПАСНОСТЬ — кратко

- Принципы: Zero Trust, Defense in Depth, Least Privilege, Secure by Default, Privacy by Design.
- Шифрование: Master Recovery Key → UMK; на сервере — зашифрованный UMK, зашифрованные file keys и метаданные.
- Восстановление: 24 слова → проверка хэша на сервере → получение encrypted UMK → расшифровка на клиенте → новые device keys.
- Encrypted Metadata + Blind Indexing: поиск по хэшам терминов без раскрытия содержимого.
- Passkeys: WebAuthn (register/begin|complete, login/begin|complete), residentKey, userVerification.
- Защита API: rate limiting (Redis), CORS. Клиент: CSP, XSS-защита, HttpOnly cookies.

---

## 6. ДИЗАЙН-СИСТЕМА — константы и темы

См. **css/design-system.css** (константы из §6.1) и **css/horizon.css** (тема Horizon §6.3). Компоненты Button/Input — см. раздел 6.2 в полном источнике; адаптивность (§6.4), a11y (§6.5), i18n (§6.6) — интеграция в app.

---

## 7. ДАННЫЕ — кратко

- Пользовательские данные (файлы, ключи, сообщения, метаданные) — зашифрованы; системные (id, email, размеры, даты, статусы) — видны серверу.
- Аудит: логи входа 90 дней, действия с файлами 30 дней, ошибки 7 дней.
- Миграции: нумерованные .up.sql / .down.sql (goose или golang-migrate).

---

## Связь с текущим репозиторием

- **Текущий бэкенд:** `backend-go/` — auth по email/password, продукты, заказы, почта (conversations/messages). Соответствует части ядра и модулей connect/trade; не реализует Passkeys, vault, шифрование метаданных.
- **Текущий фронт:** `app/` — HTML/JS, i18n, стили. Дизайн-система v4 (§6) может применяться поэтапно (константы, кнопки, инпуты, тема Horizon).
- **План реализации:** см. **IMPLEMENTATION-V4.md**.

*Документ при вставке был обрезан (раздел 7, миграции). Полную версию разделов 8–14 и недостающие фрагменты при необходимости дописать по источнику.*
