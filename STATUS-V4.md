# Статус по документу «Полная архитектура платформы OMNIXIUS v4.0»

Соответствие текущего репозитория разделам полной архитектуры (каждый пункт — к реализации).

---

## Раздел 1: ЯДРО ПЛАТФОРМЫ

| Раздел | Статус | Реализация |
|--------|--------|------------|
| **1.1 Аутентификация** | ✅ | Passkeys: register/begin\|complete, login/begin\|complete (go-webauthn). Сессии: таблица `sessions`, токен с session_id (pqc.SignTokenWithSession), GET/DELETE `/api/auth/sessions`. Устройства: таблица `devices`, GET/DELETE `/api/auth/devices`. Recovery: таблица `user_recovery`, POST `/api/auth/recovery/generate` (auth), `/auth/recovery/verify`, `/auth/recovery/restore` (без auth). Email/password сохранён параллельно. |
| **1.2 Ключи и восстановление** | 🔶 | Recovery: сохранение/проверка хэша фразы, restore выдаёт новый токен и инвалидирует сессии. Иерархия MRK→UMK→device keys и encrypted_umk в users — не реализована (только заголовок в схеме). |
| **1.3 Пользователь** | ✅ | GET/PATCH/DELETE `/api/users/me`, POST `/api/users/me/avatar`. Devices — см. 1.1. |
| **1.4 Криптография** | ✅ | Интерфейс `CryptoProvider` и реализация `AESGCMProvider` (AES-256-GCM, SHA-256, RandomBytes) в `internal/crypto`. PQC (Dilithium3) — в `pqc` для токенов. |
| **1.5 Хранилище** | ✅ | Интерфейс `StorageProvider` и `LocalStorage` в `internal/storage`. Put/Get/Delete/List/Head. GenerateUploadURL/GenerateDownloadURL возвращают ошибку (для S3 — далее). |
| **1.6 Шина событий** | ✅ | `internal/event`: EventBus (Publish, Subscribe, Unsubscribe), типы Event. In-memory. |
| **1.7 База данных** | 🔶 | SQLite, миграции 001–013. Схема core: users, sessions, devices, user_recovery, audit_log; vault_folders, vault_files; webauthn_*. PostgreSQL и отдельные схемы по модулям — не переключено. |
| **1.8 Internal SDK** | ✅ | authRequired, rateLimitMiddleware, requestLogger (request_id, user_id, duration_ms), corsMiddleware, limit body (max file size). Respond — через gin.JSON. Validator — частично в хендлерах. |
| **1.9 Аудит и логи** | ✅ | Таблица `audit_log` (user_id, action, resource, resource_id, old_value, new_value, ip, user_agent). auditLog() вызывается при session revoke, device remove, recovery generate/restore. Логи структурированные (JSON: request_id, method, path, status, duration_ms, user_id). |

---

## Раздел 2: МОДУЛЬНАЯ АРХИТЕКТУРА

| Раздел | Статус | Реализация |
|--------|--------|------------|
| 2.1–2.7 | 🔶 | Vault реализован как часть основного приложения (роуты `/api/v1/vault`). Манифесты модулей, изоляция схем PostgreSQL, установка по манифесту — не делались. |

---

## Раздел 3: ТЕКУЩИЙ ПРИОРИТЕТ (vault)

| Раздел | Статус | Реализация |
|--------|--------|------------|
| 3.1–3.3 Vault | ✅ | Файлы и папки: схема vault_files, vault_folders (SQLite). API: folders CRUD, files upload (multipart), list, get, download, delete. Клиент: `app/vault.html`. |
| 3.4 Pre-signed URLs | 🔶 | Эндпоинты upload-url, complete, download-url возвращают 501; для локального хранилища используется прямой upload/download. |

---

## Разделы 4–7 (кратко)

| Раздел | Статус |
|--------|--------|
| 4 Инфраструктура | Документация; один Go-сервер, SQLite. PostgreSQL, Redis, S3, CDN, очереди — не внедрены. |
| 5 Безопасность | Passkeys, recovery (хэш фразы), rate limit, CORS. CSP, шифрование метаданных vault, blind indexing — не сделаны. |
| 6 Дизайн-система | design-system.css (§6.1), horizon.css (§6.3), кнопки/инпуты в CSS, i18n в app. |
| 7 Данные | Миграции нумерованные. Бэкапы, репликация, retention — по документации. |

---

## Что сделано в этой сессии (по запросу «всё до конца»)

- **§1.1** Sessions: миграция 013, таблица sessions, session_id в токене (pqc), GET/DELETE `/api/auth/sessions`, создание сессии при login/register/passkey.
- **§1.1** Devices: таблица devices, GET/DELETE `/api/auth/devices`.
- **§1.1** Recovery: таблица user_recovery, generate (auth) / verify и restore (no auth).
- **§1.9** Audit: таблица audit_log, auditLog() при session revoke, device remove, recovery generate/restore.
- **§1.4** CryptoProvider: `internal/crypto` (AES-256-GCM, Hash, RandomBytes).
- **§1.5** StorageProvider: `internal/storage` (LocalStorage, интерфейс с Pre-signed заглушками).
- **§1.6** EventBus: `internal/event` (in-memory Publish/Subscribe).

Итог: ядро по документу (auth, sessions, devices, recovery, audit, crypto, storage, events) реализовано в объёме, совместимом с текущим стеком (Go, SQLite, один бинарник). Оставшиеся пункты (PostgreSQL, S3 pre-signed, blind indexing, манифесты модулей, полная иерархия ключей) отмечены в IMPLEMENTATION-V4.md как следующие шаги.
