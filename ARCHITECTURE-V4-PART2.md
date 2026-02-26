# OMNIXIUS — Полная архитектура v4.0. ЧАСТЬ 2

Операционные и интеграционные системы (платежи, уведомления, поиск, админ, импорт/экспорт).

---

## Раздел 15: ПЛАТЁЖНАЯ СИСТЕМА И ФИНАНСЫ

- **15.1** Архитектура: клиент → Payment Gateway (Stripe/Adyen) → Wallet Service (Go) → Ledger DB.
- **15.2 Модуль wallet**: схемы `wallet.balances`, `wallet.transactions`, `wallet.holds`, `wallet.deposit_addresses`; API балансов, транзакций, пополнения (фиат/крипто), вывода, переводов, холдирования (hold/release/capture).
- **15.2.3** Интерфейс `PaymentGateway` (CreateDeposit, HandleWebhook, CreateWithdrawal, GetTransactionStatus); пример Stripe Checkout.
- **15.3 Модуль trade**: схемы `trade.products`, `trade.product_media`, `trade.orders`, `trade.reviews`, `trade.disputes`, `trade.dispute_messages`; API товаров, заказов, отзывов, споров.
- **15.3.3** Процесс покупки: создание заказа → холд в wallet; подтверждение доставки → capture hold, перевод продавцу, комиссия платформы.

---

## Раздел 16: УВЕДОМЛЕНИЯ И ПУШИ

- **16.1** Модули → Notification Service → WebSocket / Email / Mobile Push.
- **16.2** Схемы: `notifications.templates`, `notifications.queue`, `notifications.user_settings`, `notifications.push_tokens`; API настроек, push-токенов, истории, теста.
- **16.2.3** Сервис: Send(), процесс очереди (pending → sent/failed), доставка по каналам (websocket, email, push).
- **16.2.4** WebSocket Hub: регистрация клиентов по userID, рассылка по UserID, exclude по deviceID.

---

## Раздел 17: ПОИСКОВАЯ СИСТЕМА

- **17.1** Search Service (Go) + опционально Elasticsearch/Meilisearch.
- **17.2 Blind Indexing**: на клиенте — генерация хэшей терминов (HMAC с user key), префиксы и n-граммы; на сервере — таблица `vault.search_index` (user_id, file_id, term_hash, weight), поиск по term_hash.
- **17.3** SearchService.SearchVault: приём захэшированных терминов от клиента, фильтры (folderId, mimeType, fromDate), сортировка, пагинация.

---

## Раздел 18: АДМИНИСТРИРОВАНИЕ И МОДЕРАЦИЯ

- **18.1** Схемы: `admin.users`, `admin.audit_log`, `admin.reports`, `admin.bans`; API админов, дашборда, пользователей (ban/unban/verify/warn), файлов, жалоб, аудита, системы (status, maintenance, logs).
- **18.1.3** ModerationService: CheckContent (AI), CreateReport, автоматические репорты и блокировки при critical.

---

## Раздел 19: ИМПОРТ/ЭКСПОРТ И МИГРАЦИЯ

- **19.1** Импорт: источники (Google, Outlook, IMAP, file); ImportService.StartImport, воркер ImportTask (ReadMessages → import в Horizon/Vault).
- **19.2** Экспорт (GDPR): ExportService.RequestExport, processExport по модулям (vault, connect, horizon, wallet), форматы json/html/csv, выгрузка в S3, ссылка на 7 дней, уведомление.

---

## Раздел 20: ТЕСТИРОВАНИЕ И КАЧЕСТВО

- **20.1** Пирамида тестирования (документ обрезан).

---

*Источник: вторая часть полной архитектуры OMNIXIUS v4.0. Реализация — по приоритету в IMPLEMENTATION-V4.md.*
