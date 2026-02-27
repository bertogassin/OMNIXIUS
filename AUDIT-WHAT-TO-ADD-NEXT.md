# OMNIXIUS — По всем 10 обновлениям: что добавить дальше

Проход по плану из `AUDIT-AND-PLAN-10-UPDATES.md`. Ниже — что уже сделано и **что можно добавить** по каждому пункту.

---

## 1. Язык и ссылки

**Сделано:** 404 (EN, OMNIXIUS, Go home / Mail), links (относительные пути, EN «All links (dev)»), contact → `#map`, статус-блоки на направлений-заглушках заменены на «Coming next».

**Добавить:**
- **Статус-блоки в app/:** Заменить оставшиеся красные `.status-done` («Live. …») на нейтральный стиль (например `.direction-coming` или убрать с прод-страниц): dashboard, learning, marketplace, media, product, product-create, profile-edit. Либо оставить «Live» но без красного акцента.
- **Якорь #screens:** На главной есть только `#map`. В `architecture.html` и `ecosystem.html` ссылки «All screens» ведут на `index.html#screens` — якоря нет. Заменить на `#map` для консистентности.
- **Mobile (android assets):** В `mobile/android/.../www/` остались русские «СДЕЛАНО», `#screens`, жёсткие URL (bertogassin.github.io). При синхронизации с app/ — привести к EN и относительным путям.

---

## 2. Wallet и Trade

**Сделано:** wallet.html (балансы, транзакции, перевод), trade.html (Balance summary с wallet API, Open Wallet, Assets placeholder, Buy/Sell coming soon), Wallet в навбаре и на index.

**Добавить:**
- **Wallet:** Поддержка verify перед переводом (api.wallet.transferVerify), если бэкенд требует.
- **Trade:** Реальный список активов/валют из API (если появится эндпоинт), а не заглушка; ссылки «Buy/Sell» на будущие экраны или модалки.
- **Навигация:** Убедиться, что Wallet есть во всех app-страницах в блоке More (часть уже добавлена).

---

## 3. Remittances

**Сделано:** remittances.html (список + форма New remittance, POST /api/remittances), блок на dashboard «Recent remittances» + «View all → Remittances», пункт в навбаре и на index.

**Добавить:**
- Фильтры по статусу/дате на странице списка.
- Валидация to_identifier (формат email/телефон/ID в зависимости от бэкенда).
- Краткая подсказка «Как отправить перевод» или ссылка на справку.

---

## 4. Settings

**Сделано:** settings.html (смена пароля, Sessions, Devices, Recovery phrase, Notifications), ссылки с dashboard и profile-edit, пункт Settings в навбаре.

**Добавить:**
- После смены пароля — опционально разлогин на других устройствах или явное сообщение «Измените пароль в активных сессиях».
- Recovery: предупреждение «Сохраните фразу в безопасном месте» перед показом/генерацией.
- Раздел «Privacy» или «Data» (экспорт данных), если появится в API.

---

## 5. Уведомления

**Сделано:** notifications.html (список из history, Mark read), пункт Notifications в навбаре app.

**Добавить:**
- **Иконка «колокольчик» в шапке app** со счётчиком непрочитанных: нужен эндпоинт `GET /api/notifications/unread-count` (или считать по history с `read_at === null`). Если бэкенд добавит unread-count — вывести в шапке и ссылку на notifications.html.
- Пагинация или «Load more» на странице истории.
- Фильтр по типу (email/push) или по дате.

---

## 6. Маркетплейс и продукт

**Сделано:** Фильтры в URL, карточки с image wrap и категорией; product — закрытый контент (Subscriber content), рассрочка на Orders (Request installments).

**Добавить:**
- **Marketplace:** Реальные фото товаров (если есть поле image/thumbnail в API), улучшение сетки на мобильном.
- **Product:** Явная кнопка «Request installments» на странице заказа (orders.html уже есть) — дублировать ссылку на product или в корзине, если появятся корзина/чекаут.
- **Gig/услуги:** Гео-подсказки и сортировка по расстоянию — при появлении полей location в API и эндпоинтов геопоиска.

---

## 7. Learning и Media

**Сделано:** learning.html — курсы из products (category Course), ссылки на Marketplace и product; media.html — подписки (subscriptions.my), discover creators (products subscription:1), Subscriber content.

**Добавить:**
- **Learning:** Отдельная таблица/API courses или experts — тогда заменить заглушку на полный каталог; кнопка «Enroll»/«Subscribe» с привязкой к подпискам или заказу.
- **Media:** Страница «Creator» (профиль креатора + его подписки/продукты) по id пользователя или slug; глубокие ссылки на product с closed_content_url.
- Общие карточки курсов/креаторов в одном стиле (например как в marketplace).

---

## 8. Admin

**Сделано:** admin.html (проверка role admin), Stats, Reports (список, фильтр, Resolve), User lookup (Ban/Unban), api.admin.* в api.js, ссылка Admin в навбаре и на index.

**Добавить:**
- **Assign:** Кнопка «Assign to me» в списке репортов (reportAssign(id, me.id)).
- Детальный просмотр репорта (reportGet) в модалке или отдельной секции (description, resolution).
- Экспорт списка репортов (CSV) или печать — по желанию.
- Ссылка «Admin» только для пользователей с role admin (скрывать в навбаре для не-админов через JS по api.user.role).

---

## 9. React (web/)

**Сделано:** Роутинг (React Router), Login, Dashboard, Marketplace (список products), API client и auth (token, user), Layout с шапкой, ProtectedRoute, API_URL через env.

**Добавить:**
- Новые экраны в SPA: Orders, Wallet, Notifications, Profile — по одному для постепенной миграции.
- Единый конфиг навбара (ссылки в одном массиве/конфиге) и синхронизация с app/*.html (те же пункты).
- Обработка 401: глобальный перехват (например в api.request) — очистка токена и редирект на /login.
- Build: base path для деплоя в подпапку (например `/app/` или `/web/`) через `vite.config.ts` base.

---

## 10. Направления-заглушки

**Сделано:** Единый шаблон Direction на startups, rewards, ixi (Coming next, What's next, Ecosystem/Architecture/Contact, Join waitlist); repositorium — What's next, ссылки, Join waitlist; index #map — «Coming next: …» для Forge, Ark, IXI; Admin в навбаре направлений; trade — direction-coming.

**Добавить:**
- **Waitlist:** Опционально форма «Notify me» / «Join waitlist» с сохранением email в БД (эндпоинт + таблица waitlist) вместо только ссылки на contact.
- **Repositorium:** Блок «My node» / «Rewards» — заглушка или ссылка на будущий раздел, когда появится бэкенд.
- **Консистентность:** Проверить, что на всех app-страницах в More есть Admin (часть уже добавлена; при добавлении новых страниц — включать Admin в шаблон навбара).
- **Mobile:** При обновлении android assets — применить тот же формат Direction и «Coming next» на соответствующих экранах.

---

## Сводная таблица приоритетов «что добавить»

| № | Что добавить в первую очередь |
|---|------------------------------|
| 1 | Заменить #screens на #map в architecture.html, ecosystem.html; убрать/смягчить красный status-done в app/ |
| 2 | Verify для перевода в Wallet; при появлении API активов — вывести их на Trade |
| 3 | Фильтры и валидация на Remittances |
| 4 | Подсказки по безопасности в Settings (recovery, смена пароля) |
| 5 | Колокольчик с unread в шапке app (после эндпоинта unread-count или расчёта по history) |
| 6 | Фото товаров на маркетплейсе; явная ссылка на рассрочку с product |
| 7 | API/таблицы courses и experts для Learning; страница Creator для Media |
| 8 | Assign to me в Admin; скрывать ссылку Admin для не-админов |
| 9 | Новые экраны в web/ (Orders, Wallet и т.д.); 401 → logout и редирект на login |
| 10 | Waitlist в БД (опционально); синхронизация mobile с форматом Direction |

После этих шагов все 10 обновлений будут доведены до следующего уровня без изменения текущего плана.
