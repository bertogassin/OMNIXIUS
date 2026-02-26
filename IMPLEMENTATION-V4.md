# План реализации по архитектуре v4.0

Источник: **ARCHITECTURE-V4.md** (полная архитектура платформы v4.0). Каждый пункт архитектуры — к реализации.

---

## Фазы (приоритет)

| Фаза | Содержание | Связь с текущим репо |
|------|------------|----------------------|
| **Фаза 0** | Дизайн-система v4 | ✅ Константы в `css/design-system.css`, подключены. Тема Horizon в `css/horizon.css`. |
| **Фаза 1** | Ядро: auth + user | ✅ Passkeys: миграция 012 (webauthn_credentials, webauthn_sessions), go-webauthn: register/begin|complete, login/begin|complete. Регистрация по email+name (passkey-only аккаунт), вход по email + passkey. Токен — тот же PQC (Dilithium3). |
| **Фаза 2** | Ядро: crypto, storage, events, SDK | ✅ CryptoProvider (`internal/crypto`: AES-256-GCM). StorageProvider (`internal/storage`: LocalStorage). EventBus (`internal/event`: in-memory). SDK: request_id, user_id, duration_ms в логере. |
| **Фаза 3** | Модуль vault | ✅ Схема SQLite: `vault_folders`, `vault_files`. API: POST/GET/DELETE folders, POST/GET/GET download/DELETE files (multipart upload). Клиент: `app/vault.html`. Pre-signed URL — 501. |
| **Фаза 4** | Безопасность и данные | ✅ Аудит: таблица `audit_log`, вызовы при session/device/recovery. 🔲 Шифрование метаданных vault, blind indexing, полный key management (UMK) — далее. |
| **Фаза 5** | Инфраструктура и модули | PostgreSQL (схема core + миграции), Redis, S3 (конфиг, CORS). Очереди для фоновых задач. Подключение модулей connect, horizon, trade по манифестам. |

---

## Первые шаги — сделано

1. **Дизайн-система** — `css/design-system.css` создан и подключён в index.html.
2. **Документация** — README, POLNOE-TZ ссылаются на ARCHITECTURE-V4.md и IMPLEMENTATION-V4.md.
3. **Vault** — миграция 011_vault.sql, полные обработчики vault (folders + files, загрузка multipart, скачивание), страница `app/vault.html`. Passkeys — роуты-заглушки (501).

---

## Соответствие разделов архитектуры и задач

| Раздел доки | Что реализовать |
|-------------|------------------|
| 1.1 Auth | ✅ WebAuthn, сессии (sessions + session_id в токене), устройства (devices), recovery (generate/verify/restore). |
| 1.2 Ключи | Иерархия MRK → UMK → device/encryption/signing; хранение encrypted_umk; добавление устройства (QR, одноразовый токен). |
| 1.3 User | GET/PATCH/DELETE /user/me, avatar, devices (trust, delete). |
| 1.4 Crypto | ✅ CryptoProvider в internal/crypto (AES-256-GCM, Hash, RandomBytes). |
| 1.5 Storage | ✅ StorageProvider, LocalStorage в internal/storage; Pre-signed — заглушки. |
| 1.6 Event bus | ✅ EventBus в internal/event (Publish, Subscribe, in-memory). |
| 1.7 DB | SQLite, миграции до 013 (sessions, devices, user_recovery, audit_log). PostgreSQL — в планах. |
| 1.8 SDK | ✅ AuthMiddleware, RateLimit, Logger (request_id, duration_ms), CORS, body limit. |
| 1.9 Аудит | ✅ audit_log таблица, auditLog() при session/device/recovery; структурированные логи. |
| 2.x Модули | Манифесты, изоляция схем, scopes, установка модулей. |
| 3.x Vault | Схема vault.files, vault.folders; API файлов и папок; pre-signed URLs; поиск (blind index). |
| 4.x Инфра | DNS, SSL, API-сервер (Go), pool БД, S3, CDN, Redis, очереди. |
| 5.x Безопасность | Rate limit, CORS, CSP, ключи, восстановление, шифрование метаданных. |
| 6.x Дизайн | Константы, Button, Input, Horizon, адаптивность, a11y, i18n. |
| 7.x Данные | Миграции, бэкапы, репликация, retention. |

---

---

## Часть 2 архитектуры (операционные системы)

См. **ARCHITECTURE-V4-PART2.md**: платёжная система и wallet (§15), trade-схема и заказы с холдом (§15.3), уведомления и WebSocket (§16), поиск и blind indexing (§17), админ и модерация (§18), импорт/экспорт (§19), тестирование (§20). Реализация — по приоритету после завершения ядра и vault.

---

*Обновлять по мере реализации. При изменении ARCHITECTURE-V4.md — синхронизировать этот план.*
