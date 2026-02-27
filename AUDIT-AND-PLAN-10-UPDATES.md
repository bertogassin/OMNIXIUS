# OMNIXIUS — Аудит сайта и план 10 масштабных обновлений

Обход выполнен по: главная, все страницы `app/`, бэкенд (backend-go), документы (ARCHITECTURE, ROADMAP, ECOSYSTEM, API.md), заготовка React (web/).

---

## 1. Что некорректно или стоит исправить

| Проблема | Где | Рекомендация |
|----------|-----|--------------|
| **Смешение языков** | «СДЕЛАНО» на dashboard, trade, learning, startups, media, rewards, marketplace; 404 и links — русский текст | По ARCHITECTURE: английский по умолчанию. Убрать/перевести статус-блоки; 404 и links — на EN или в dev-only |
| **links.html** | Жёсткие URL `bertogassin.github.io`, весь текст на русском | Относительные пути; текст на EN или пометка «Dev links» |
| **404.html** | `lang="ru"`, заголовок «Horizon Mail», русский текст | `lang="en"`, бренд OMNIXIUS, нейтральный текст «Page not found» |
| **contact.html** | Ссылка `index.html#screens` | На главной нет якоря `#screens`, есть `#map` — заменить на `#map` или «All sections» |
| **Дублирование навигации** | В app-страницах: и `nav-drop-panel`, и `nav-group-app-links` с одними и теми же ссылками | Оставить один вариант (dropdown на мобильном, список на десктопе) для упрощения вёрстки |
| **Нет единой страницы настроек** | В API: смена пароля, сессии, устройства, recovery, уведомления | Один экран Settings (Security, Sessions, Notifications, Recovery) вместо разбросанных действий |
| **Нет центра уведомлений** | API: `/notifications/history`, `/history/:id`, read | Добавить страницу или блок «Уведомления» (иконка в шапке + список) |
| **Нет UI для создания ремита** | API: POST /api/remittances; на dashboard только список | Отдельная страница «Переводы» с формой «Создать перевод» (to_identifier, amount, currency) |
| **Wallet API без UI** | Backend: /wallet/balances, transactions, transfer, hold, deposit | Баланс есть на dashboard; нет экрана «Кошелёк» с мультивалютой и историей (см. обновление 2) |

**Что уже в порядке:** Verified-бейдж есть (dashboard, marketplace, profile, product). API.md актуален. Закрытый контент (closed-content.html) и подписки заложены. Vault (Crate) с папками/файлами работает. Слоты и бронирование в API и на product.

---

## 2. Каких экранов не хватает

| Экран | Статус | Что нужно |
|-------|--------|-----------|
| **Wallet (Forge)** | Нет | Страница «Кошелёк»: балансы по валютам (GET /wallet/balances), история (GET /wallet/transactions), форма перевода (POST /wallet/transfer). Сейчас только общий balance на dashboard |
| **Remittances (отдельная страница)** | Частично | Список на dashboard есть; нет страницы с формой «Новый перевод» и фильтрами |
| **Trade (Forge) — MVP** | Заглушка | Минимум: список активов/валют, заглушка «Купить/Продать», ссылка на Wallet |
| **Notifications** | Нет | Страница или панель: список уведомлений, «прочитано», настройки (уже в API) |
| **Settings (единая)** | Нет | Объединить: смена пароля, сессии/устройства, recovery codes, настройки уведомлений |
| **Admin** | Нет | API есть (reports, users, ban/unban, stats). Нужна админ-панель (доступ по role) |
| **Learning (Ascent)** | Заглушка | Каталог курсов, карточки экспертов, запись на курс (модели и API — позже) |
| **Startups (Flare)** | Заглушка | Список проектов, карточка «Launchpad» (заглушка), позже — сбор средств |
| **Media (Lens)** | Заглушка | Лента/каталог креаторов, подписки (подключить к /subscriptions/my) |
| **Rewards (Bounty)** | Заглушка | Реферальная ссылка, таблица наград, история (нужны API и модели) |
| **Repositorium (Ark)** | Лендинг | Уже описание + блоки; дальше — «Мой узел», награды (нужен бэкенд) |
| **IXI (Blockchain)** | Заглушка | Блок-эксплорер, стейкинг, токены — после Phase 2 инфраструктуры |

---

## 3. План 10 масштабных обновлений (по порядку)

### 1. Единообразие языка и корректность ссылок
- Перевести все статус-блоки «СДЕЛАНО» на EN или убрать с прод-страниц (оставить опционально в dev).
- 404: `lang="en"`, заголовок «OMNIXIUS», текст «Page not found» / «Go home», ссылки на index и app.
- links.html: относительные ссылки, EN или вынести в `app/dev/links.html` с пометкой.
- contact.html: заменить `#screens` на `#map`.

**Результат:** Единый язык по умолчанию (EN), рабочие якоря, нет битых ссылок.

---

### 2. Страница Wallet (Forge) и доработка Trade
- Новая страница `app/wallet.html`: балансы по валютам (GET /api/wallet/balances), список транзакций (GET /api/wallet/transactions), кнопка «Transfer» (форма с verify при необходимости).
- В навигации (app и index) добавить пункт «Wallet» или «Кошелёк» в группу Forge.
- trade.html: не только текст «Place ready», а блоки «Balance summary» (кратко с dashboard или wallet), «Assets» (заглушка списка), «Quick actions: Buy / Sell (coming soon)».

**Результат:** Пользователь видит свой кошелёк и историю; Forge получает первый рабочий экран.

---

### 3. Страница «Переводы» (Remittances) и форма создания
- Новая страница `app/remittances.html`: список моих ремитов (как на dashboard) + форма «New remittance» (to_identifier, amount, currency), вызов POST /api/remittances.
- В dashboard оставить краткий блок «Recent remittances» со ссылкой «View all → remittances.html».
- В навигации (More или внутри Forge) добавить «Remittances».

**Результат:** Полный цикл ремитов из UI без вызова API вручную.

---

### 4. Единая страница настроек (Settings)
- Новая страница `app/settings.html` (или разделы: Security, Sessions, Notifications, Recovery).
- Блоки: смена пароля (POST /auth/change-password), активные сессии (GET/DELETE /auth/sessions), устройства (GET/DELETE /auth/devices), recovery codes (generate, restore), настройки уведомлений (GET/PATCH /notifications/settings).
- Ссылка «Settings» в профиле (dashboard, profile-edit) и в шапке (иконка или «More»).

**Результат:** Один вход для безопасности и уведомлений, меньше разбросанных действий.

---

### 5. Центр уведомлений
- Страница `app/notifications.html` или выпадающая панель в шапке: список GET /notifications/history, пометка «прочитано» POST .../read, ссылка на настройки.
- В шапке app: иконка «колокольчик» с счётчиком (если API даёт unread count — при необходимости отдельный эндпоинт или из history).

**Результат:** Пользователь видит уведомления и управляет ими.

---

### 6. Доработка маркетплейса и продукта (Trove)
- Marketplace: улучшение карточек (фото, категории, бейдж Verified уже есть), сохранение фильтров в URL (query params) для шаринга.
- Product: рассрочка — кнопка «Request installments» на заказе, PATCH installment_plan (уже в API); закрытый контент — явная кнопка «View subscriber content» с переходом на closed-content или iframe.
- Категории и «услуги рядом» (Gig): фильтр по location и is_service уже есть; при необходимости — гео-подсказки и сортировка по расстоянию (бэкенд).

**Результат:** Trove и заказы ближе к Phase A (Gig, Creator economy).

---

### 7. Learning (Ascent) и Media (Lens) — первый контент
- Learning: каталог «курсов» — пока заглушка или список из API (если добавить таблицу courses/experts в бэкенд); карточки курсов, кнопка «Subscribe» / «Enroll» (подписка через Trade/подписки).
- Media: лента или каталог «креаторов» — подписки из GET /subscriptions/my; страница «Creator» = профиль + список подписок пользователя; кнопка «Subscriber content» на продуктах с closed_content_url.

**Результат:** Ascent и Lens перестают быть пустыми заглушками; задел под Creator economy.

---

### 8. Админ-панель (первая версия)
- Страница `app/admin.html` (доступ по role, например `admin`): дашборд GET /api/admin/stats, список репортов GET /api/admin/reports с assign/resolve, просмотр пользователя GET /api/admin/users/:id, ban/unban. Роутинг только для авторизованных с ролью admin.

**Результат:** Модерация и поддержка без прямых вызовов API.

---

### 9. Реакт-приложение (web/) — первый полезный слой
- По ARCHITECTURE: при рефакторинге — TS + React. В web/: роутинг (React Router), экраны Login, Dashboard (или главный app), вызовы API (fetch к backend-go), хранение token (sessionStorage/localStorage). Не обязательно переносить все страницы сразу — начать с одного потока (например логин → dashboard → marketplace список).
- Общий хедер/навигация и конфиг API_URL (env или config).

**Результат:** Основа для постепенной миграции с app/*.html на SPA.

---

### 10. Направления-заглушки: единый формат и навигация
- Startups, Rewards, Repositorium (сверх текущего лендинга), IXI: единый шаблон «Direction»: описание, «What’s next», ссылки на Ecosystem/Architecture, кнопки «Notify me» или «Join waitlist» (опционально — форма в БД или email).
- Убрать повторяющиеся «СДЕЛАНО / Place ready» в пользу короткого «Coming next: …» на EN.
- Проверить, что все пункты главной и карты (#map) ведут на актуальные страницы; добавить Wallet и Remittances в карту/навигацию где нужно.

**Результат:** Консистентный вид всех направлений и предсказуемая навигация по экосистеме.

---

## 4. Сводка по приоритетам

| № | Обновление | Затрагивает | Приоритет |
|---|------------|-------------|-----------|
| 1 | Язык и ссылки | 404, links, contact, статус-блоки | Высокий (корректность) |
| 2 | Wallet + Trade | app/wallet.html, trade.html, навигация | Высокий (ценность Forge) |
| 3 | Remittances UI | app/remittances.html, dashboard | Высокий |
| 4 | Settings | app/settings.html, профиль | Высокий (безопасность, UX) |
| 5 | Уведомления | app/notifications.html или панель в шапке | Средний |
| 6 | Маркетплейс и продукт | marketplace, product, заказы/рассрочка | Высокий (Trove/Gig) |
| 7 | Learning + Media | learning.html, media.html, подписки | Средний |
| 8 | Admin | app/admin.html, роль admin | Средний |
| 9 | React web | web/ (роутинг, auth, 1–2 экрана) | Средний (долгосрочно) |
| 10 | Единый формат направлений | startups, rewards, ixi, repositorium, копии | Низкий (консистентность) |

После выполнения плана: язык и ссылки приведены в порядок, появляются полноценные экраны Wallet, Remittances, Settings и уведомлений, усиливаются Trove и Forge, закладывается админка и React, все направления выглядят единообразно.
