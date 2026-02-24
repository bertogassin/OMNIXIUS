# OMNIXIUS — AI (Python)

Корень своего ИИ: нейросети, распознавание. Сейчас — FastAPI с эндпоинтом `/chat` (заглушка); дальше — свои модели.

**Уровень:** начало. Страница чата в браузере: `app/ai.html`; бэкенд здесь.

```bash
pip install -r requirements.txt
uvicorn main:app --reload --port 8000
# http://localhost:8000/docs  — Swagger
# http://localhost:8000/health
```
