# OMNIXIUS

Сайт OMNIXIUS — лендинг на HTML/CSS/JS.

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

### 4. Включи GitHub Pages

1. На странице репозитория на GitHub открой **Settings**.
2. Слева выбери **Pages**.
3. В блоке **Source** выбери **Deploy from a branch**.
4. В **Branch** выбери **main** и папку **/ (root)**.
5. Нажми **Save**.

Через 1–2 минуты сайт будет доступен по адресу:

**https://ТВОЙ_ЛОГИН.github.io/ИМЯ_РЕПОЗИТОРИЯ/**

Например: `https://ivanov.github.io/omnixius/`

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
