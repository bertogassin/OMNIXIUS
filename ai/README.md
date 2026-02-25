# OMNIXIUS — AI (Python)

Корень своего ИИ: нейросети, распознавание. FastAPI, эндпоинт `/chat`. **A1/A2:** агент вызывает backend-go по интентам и возвращает результат.

**Уровень:** начало. Страница чата в браузере: `app/ai.html`; бэкенд здесь.

## Запуск

**Проще всего:** из корня репо дважды кликни **`start-ai.bat`** или в PowerShell:
```powershell
.\start-ai.ps1
```
Скрипты запускают сервер через `python -m uvicorn` (не требуют `uvicorn.exe` в PATH).

**Вручную:**
```bash
cd ai
pip install -r requirements.txt
python -m uvicorn main:app --reload --port 8000
```
- **http://localhost:8000** — API
- **http://localhost:8000/docs** — Swagger
- **http://localhost:8000/health** — проверка

**Если pip падает с ошибкой записи в `Scripts` (WinError 2):**
1. Установи в папку пользователя: `pip install --user -r requirements.txt`
2. Запускай только так: `python -m uvicorn main:app --reload --port 8000` (из папки `ai`) или используй `start-ai.bat`.

---

## Intent catalog (A1/A2)

Чат передаёт в тело запроса `api_token` и `api_url` (из приложения). По тексту сообщения определяется интент; при совпадении вызывается backend-go и ответ форматируется в текст.

| Intent | Example phrases (EN / RU / FR) | Backend call | Description |
|--------|-------------------------------|--------------|-------------|
| **My orders** | "my orders", "show orders", "мои заказы", "заказы", "mes commandes" | `GET /api/orders/my` | Список заказов пользователя. |
| **Create order** | "create order 5", "order product 5", "buy 3", "оформить заказ на продукт 5", "заказать 5" | `POST /api/orders` body `{"product_id": N}` | Создание заказа по ID продукта (число в сообщении). |
| **Conversations summary** | "my messages", "conversations", "письма", "переписки", "сводка по письмам", "mes messages" | `GET /api/conversations` | Краткая сводка: диалоги, последнее сообщение, непрочитанные. |

Дальше (Phase A): больше интентов и действий — см. ROADMAP.md.
