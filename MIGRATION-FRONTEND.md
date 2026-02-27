# Миграция фронта: JS/HTML убраны до конца

По стеку: **веб только на TypeScript**; JS/HTML убраны по максимуму.

## Текущее состояние

- **web/** — **единственное приложение** (TypeScript + React). Полный API в `api.ts`; страницы: Dashboard, Marketplace, Login, Orders, Order, Find professional, Notifications, Profile, Settings. Раздаётся по **/app** (backend-go отдаёт `web/dist` с SPA fallback).
- **app/** — **удалён** (каталог с HTML/JS убран до конца).
- **js/** — оставлен **минимально** только для лендинга: `professions.js`, `i18n.js`, `main.js` (используются в корневом `index.html`).

## Что делать

1. **Разработка** — только в **web/** (TypeScript). `cd web && npm run dev` (dev-сервер на 5173). Перед деплоем: `npm run build`.
2. **Деплой** — бэкенд раздаёт корень (index.html, css/, js/) и **/app** → `web/dist` (SPA). Сборка обязательна: `cd web && npm run build`.

## Итог

- Фронт приложения = **только TypeScript (web/)**.  
- JS/HTML = только лендинг в корне (index.html + js/professions.js, i18n.js, main.js). Остальное убрано.
