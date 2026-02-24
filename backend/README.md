# OMNIXIUS API (Node.js) — **DEPRECATED**

**Используй `backend-go/`** как единственный бэкенд. Контракт API описан в корне репо: **API.md**.

Этот каталог — legacy Node.js-реализация; оставлен для справки. Не дублируй логику здесь и в backend-go.

---

Бэкенд: регистрация, маркетплейс, заказы, внутренняя почта.

## Запуск

```bash
cd backend
npm install
npm run init-db
npm start
```

API: `http://localhost:3000`

## Переменные окружения

- `PORT` — порт (по умолчанию 3000)
- `JWT_SECRET` — секрет для JWT (обязательно смени в продакшене)
- `JWT_EXPIRES` — срок жизни токена (например `7d`)

## Эндпоинты

- `POST /api/auth/register` — регистрация (email, password, name)
- `POST /api/auth/login` — вход (email, password)
- `GET /api/auth/confirm-email?token=...` — подтверждение email
- `POST /api/auth/forgot-password` — запрос сброса пароля
- `POST /api/auth/reset-password` — сброс пароля (token, password)

- `GET /api/users/me` — профиль (auth)
- `PATCH /api/users/me` — обновить имя (auth)
- `GET /api/users/me/orders` — мои заказы (auth)
- `POST /api/users/me/avatar` — загрузить аватар (auth, multipart)

- `GET /api/products` — список (q, category, location, minPrice, maxPrice, sort)
- `GET /api/products/categories` — категории
- `GET /api/products/:id` — товар
- `POST /api/products` — создать (auth, multipart: title, description, price, category, location, image)
- `PATCH /api/products/:id` — обновить (auth)
- `DELETE /api/products/:id` — удалить (auth)

- `GET /api/orders/my` — мои заказы (auth)
- `POST /api/orders` — создать заказ (auth, product_id)
- `PATCH /api/orders/:id` — статус (auth, status)

- `GET /api/conversations` — мои диалоги (auth)
- `POST /api/conversations` — создать/получить диалог (auth, user_id, product_id?)
- `GET /api/messages/conversation/:id` — сообщения (auth)
- `POST /api/messages/conversation/:id` — отправить (auth, body)
- `POST /api/messages/:id/read` — отметить прочитанным (auth)

## Безопасность

- Пароли: bcrypt
- JWT в заголовке `Authorization: Bearer <token>`
- Ограничение попыток входа (5 за 15 мин)
- Rate limit: 200 запросов / 15 мин
- Валидация и экранирование ввода
- Загрузка файлов: только JPEG/PNG/WebP, макс. 5 МБ

## БД

SQLite: `backend/db/omnixius.db`. Схема в `db/schema.sql`. Для продакшена можно заменить на PostgreSQL (замена драйвера и строки подключения).
