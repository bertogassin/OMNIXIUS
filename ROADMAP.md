# OMNIXIUS — Roadmap

Фазы и приоритеты. Стек: **Go, Rust, Spark, Swift, Kotlin, TS/JS, Python** по направлениям (ARCHITECTURE.md).

---

## Текущая фаза: Foundation (мы здесь)

- Сайт и документация (архитектура, видение, стек, чеклист).
- Рабочее приложение: регистрация, логин, маркетплейс, почта, заказы, профиль, ИИ (корень).
- Один бэкенд (Go), контракт API, безопасность (PQC, Argon2id).
- Заготовки по всем направлениям: Rust, Spark, iOS, Android, React+TS, AI, инфра (Docker, K8s, Terraform, CI/CD, мониторинг).

**Уровень:** Близки к началу, далеки от конца. Фундамент заложен; дальше — наращивание по фазам ниже.

---

## Phase 1 — Infrastructure & core

- Инфраструктура (multi-cloud, edge, мониторинг в проде).
- IXI blockchain prototype (consensus, ZKP, economy). **Rust / C++.**

## Phase 2 — Platforms

- OMNIXIUS CONNECT: профили, geo, чат, маркетплейс (расширение). **Go** API, **Swift** iOS, **Kotlin** Android.
- OMNIXIUS TRADE & FINANCE: кошельки, биржа (crypto/FIAT), copy trading (MVP). **Go / Rust.**
- OMNIXIUS REPOSITORIUM: узлы хранения, награды, аренда compute (MVP). **Rust / C++**, **Spark** для аналитики.

## Phase 3 — Integration & scale

- Все платформы связаны через IXI и один аккаунт.
- ИИ: модерация, аналитика, свои модели.
- Мобильные приложения, глобальное масштабирование.

---

## Ten new directions (ARCHITECTURE.md §6)

Архитектура под 10 направлений зафиксирована; наращивание — по фазам ниже.

### Phase A — First value

- **AI-агенты** — действия поверх текущего AI (брони, оплаты, сводки). AI + Connect/Trade API.
- **Gig / Локальные услуги** — оформление «специалисты рядом»: брони, оплата, отзывы. Connect + Trade.
- **Creator economy** — подписки, донаты, закрытый контент. Connect + Trade + Repositorium.

#### Phase A — шаги (что делать по порядку)

| Шаг | Направление | Действие | Стек / место |
|-----|-------------|----------|--------------|
| A1 | AI-агенты | В AI-чате добавить «действия»: запрос к backend-go (список заказов, создание заказа, сводка). Пользователь пишет «покажи мои заказы» — агент дергает API и отвечает. | `ai/` (Python) вызывает backend-go; в чате — кнопка или интент «действие». |
| A2 | AI-агенты | Каталог интентов: показать заказы, создать заказ, кратко по письмам. Документировать в API.md или ai/README. | ai/, backend-go (существующие эндпоинты). |
| A3 | Gig | В Connect/Marketplace: флаг «услуга» у листинга, поле «гео» (город/район). Фильтр «услуги рядом». | backend-go (products/orders), app/marketplace.html, app/connect.html. |
| A4 | Gig | Бронирование слота: дата/время у листинга-услуги, запрос от клиента → создание заказа + уведомление в Mail. | backend-go (orders, mail), app (страница товара/услуги). |
| A5 | Creator economy | Тип листинга «подписка»: цена за месяц, повторяющийся платёж. В профиле креатора — «Подписаться». | backend-go (products, subscriptions stub), Trade (платежи позже). |
| A6 | Creator economy | Закрытый контент: ссылка на файл/страницу только для подписчиков. Repositorium или просто URL с проверкой подписки. | backend-go (check subscription), app (media/learning или отдельная страница). |

После A1–A2 агенты дают быструю пользу; A3–A4 открывают Gig; A5–A6 — базу для креаторов. Дальше Phase B.

### Phase B — Identity & money

- **SSI / Децентрализованная идентичность** — верифицируемые credentials на IXI, один «паспорт».
- **Embedded finance** — микрокредиты, страховка, рассрочка внутри маркетплейса. Trade.
- **Cross-border / Ремиты** — низкая комиссия, стейблкоины/IXI. Trade + IXI.

#### Phase B — шаги (что делать по порядку)

| Шаг | Направление | Действие | Стек / место |
|-----|-------------|----------|--------------|
| B1 | SSI / Identity | Верифицированный профиль: флаги verified_email, verified_phone в БД; в профиле и карточке пользователя бейдж «Verified». Подтверждение email уже есть; опционально — подтверждение телефона (заглушка или SMS). | backend-go (users: verified_*), app (profile, dashboard, карточки продавца). |
| B2 | Embedded finance | Рассрочка по заказу: поле installment_plan у заказа (boolean или enum: full / installments). В UI заказа — кнопка «Оформить рассрочку» → заглушка (модалка «Скоро через Trade» или запись намерения в БД). Документировать сценарий в API.md. | backend-go (orders), app (orders.html, checkout flow). |
| B3 | Trade stub | Баланс пользователя (внутренние единицы): таблица user_balances или поле balance; GET /api/users/me/balance; заглушка пополнения (админ или тестовый кредит). Основа для подписок, наград, выплат. | backend-go (balances), API. |
| B4 | Cross-border / Ремиты | Дизайн API ремитов: POST /api/remittances (from_user_id, to_identifier, amount, currency). Валидация и запись заявки в БД; реальный перевод — позже (Trade/IXI). Документировать в API.md. | backend-go (remittances stub), API.md. |
| B5 | Cross-border / Ремиты | Список моих ремитов: GET /api/remittances/my; отображение в профиле или странице «Мои переводы» (статус, получатель, сумма). | backend-go (remittances), app (dashboard или remittances.html). |

После B1–B2 пользователь получает «доверенный» профиль и задел под рассрочку; B3–B4 закладывают баланс и ремиты под Phase 2 (полноценный Trade/IXI); B5 — просмотр заявок на перевод.

### Phase C — Regulation & trust

- **ESG / Зелёная повестка** — углерод, офсеты, зелёные активы. Trade + Connect.
- **Health / Телемедицина** — приёмы, данные в Repo, оплата через Trade. Connect + Repositorium.
- **Privacy-first аналитика** — отчёты без PII. Spark + Repositorium.
- **Education credentials** — дипломы/бейджи on-chain. IXI + Connect + Learning.

---

IXI — блокчейн-платформа экосистемы, не отдельный токен.
