# OMNIXIUS

Экосистема: сайт, приложение (маркетплейс, почта, заказы, ИИ), полное видение — **один репозиторий, деплой через GitHub**.

**Уровень проекта:** мы **близки к началу** и **далеки от конца**. Фундамент заложен: стек зафиксирован, документация и заготовки по всем направлениям на месте. Дальше — поэтапное наращивание по ROADMAP.

---

**Стек (по направлениям):** сервер — Go, Rust, Spark; веб — TS/JS, React/Vue; мобилки — Swift, Kotlin; игры — C++/C#; ИИ — Python; инфра — Bash/Python, YAML. Полная таблица: **ARCHITECTURE.md** (§2).

**Структура репо:**  
`backend-go/` — основной API (Go); `app/` — веб-приложение (HTML/JS); `web/` — заготовка на React+TypeScript; `services/rust/`, `analytics/spark/` — Rust и Spark; `mobile/ios/`, `mobile/android/` — Swift и Kotlin; `ai/` — ИИ (Python); `infra/` — Docker, K8s, Terraform, мониторинг. `backend/` — legacy Node (deprecated).

**Production:**  
https://bertogassin.github.io/OMNIXIUS/  
https://bertogassin.github.io/OMNIXIUS/app/marketplace.html  
https://bertogassin.github.io/OMNIXIUS/app/ai.html (ИИ)

**Документация:** ARCHITECTURE.md (видение и стек) · PLATFORM.md (что в проде) · ECOSYSTEM.md (карта направлений) · ROADMAP.md (фазы) · API.md (эндпоинты) · PLATFORM_INFRASTRUCTURE.md (Docker, K8s, CI/CD) · OMNIXIUS_CHECKLIST.md (проверка всего в репо).

Деплой: push в репозиторий → GitHub Pages обновляет сайт сам.

---

## Как выложить через GitHub (GitHub Pages)

### 1. Установи Git (если ещё нет)

Скачай и установи: https://git-scm.com/download/win

### 2. Создай репозиторий на GitHub

1. Зайди на https://github.com и войди в аккаунт.
2. Нажми **+** → **New repository**.
3. Имя репозитория: например **omnixius** или **site-omnixius**.
4. Выбери **Public**, галочку "Add a README" можно не ставить.
5. Нажми **Create repository**.

### 3. Загрузи сайт с компьютера

Открой **PowerShell** или **Терминал** и выполни по очереди (подставь свой логин GitHub и имя репозитория):

```bash
cd "C:\Users\zabir\Desktop\SITE OMNIXIUS"

git init
git add index.html contact.html 404.html css js README.md .gitignore
git commit -m "Первый коммит — сайт OMNIXIUS"

git branch -M main
git remote add origin https://github.com/ТВОЙ_ЛОГИН/ИМЯ_РЕПОЗИТОРИЯ.git
git push -u origin main
```

Когда спросит логин и пароль: логин — твой GitHub, пароль — **Personal Access Token** (не обычный пароль).  
Создать токен: GitHub → Settings → Developer settings → Personal access tokens → Generate new token (поставь галочку `repo`).

### 4. Включи GitHub Pages (подробно)

Ты на странице репозитория **OMNIXIUS** на GitHub. Что вносить по шагам:

1. **Зайди в настройки**  
   Вверху страницы репозитория вкладки: Code, Issues, Pull requests… — нажми **Settings** (иконка шестерёнки).

2. **Открой раздел Pages**  
   В левом меню настроек найдите блок **"Code and automation"**. В нём пункт **Pages** — нажми на него.

3. **Блок "Build and deployment"**  
   Сверху будет надпись **Source** и выпадающий список. Сейчас там может быть "Deploy from a branch" или "None".

4. **Что выбрать в Source**  
   В выпадающем списке **Source** выбери:  
   **Deploy from a branch**  
   (не "GitHub Actions").

5. **Что выбрать в Branch**  
   Под Source появится блок **Branch**:
   - В первом выпадающем списке выбери ветку: **main**.
   - Во втором выпадающем списке выбери папку: **/ (root)**.  
   Оставь **main** и **/ (root)** — ничего больше не меняй.

6. **Save**  
   Нажми оранжевую кнопку **Save** (справа в блоке Branch).  
   После этого страница обновится, сверху может появиться синее сообщение: сайт собирается (building). Подожди 1–2 минуты.

7. **Где взять ссылку на сайт**  
   Чуть выше блока Source после сохранения появится зелёная плашка примерно такого вида:  
   **"Your site is live at https://bertogassin.github.io/OMNIXIUS/"**  
   Это и есть адрес твоего сайта. Открой его в браузере.

Итого, что вносить: **Source** → **Deploy from a branch**; **Branch** → **main**; папка → **/ (root)**; затем **Save**.

### 5. Свой домен (необязательно)

В той же вкладке **Pages** внизу есть **Custom domain** — впиши свой домен (например `omnixius.com`). GitHub покажет, какую DNS-запись добавить у регистратора домена.

---

После этого при любых изменениях достаточно выполнить в папке сайта:

```bash
git add .
git commit -m "Обновил сайт"
git push
```

Сайт на GitHub Pages обновится сам.
