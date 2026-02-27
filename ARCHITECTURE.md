# OMNIXIUS — Architecture & Vision

Global technology ecosystem: social network, professional platform, exchange, decentralized cloud, IXI blockchain, AI analytics, and investment tools.

**Goal:** Self-sufficient ecosystem — communication, earning, investment, learning, and asset management in one application, without dependence on centralized corporations.

**Итог важнее порядка:** Что раскатывать первым — не главное. Главное — чтобы в итоге **сайт, приложение и проект были лучшими во всём мире**: качество, безопасность, скорость, дизайн, масштаб. Каждое решение ведёт к этому итогу.

**Статус:** Проект в начале пути; видение и стек зафиксированы, заготовки по всем направлениям есть в репо. Наращивание — по ROADMAP и ECOSYSTEM.

---

## 1. Four platforms

### OMNIXIUS CONNECT
- Social network + expert platform.
- Profiles, geo-based orders, AI matching, ratings and reviews.
- Chat, audio/video calls, consultation recording, smart contracts.
- Marketplace for goods and services (crypto, FIAT, cards, PayPal), subscriptions.
- Security: smart contracts, AI moderation.

### OMNIXIUS TRADE & FINANCE
- Exchange: crypto, gold, oil, stocks, indices, startups.
- Buyback fund for IXI blockchain economy.
- Copy trading, long-term portfolios, auto-investing with AI.
- Wallets, Apple Pay, Google Pay, PayPal, cards.

### OMNIXIUS REPOSITORIUM
- Decentralized cloud: home servers, rewards for storage and compute.
- IPFS + AI encryption, blockchain access control.
- Compute rental, referral programs.

### OMNIXIUS BLOCKCHAIN (IXI)
**IXI is the blockchain (base platform), not a token.**
- Zero-Knowledge Proofs (ZKP).
- Hybrid Proof-of-Stake + Proof-of-Storage.
- AI transaction monitoring.
- 30-year emission, buyback fund.

---

## 2. Technology stack

**Полный сверхстабильный стек на десятилетия (Rust, Go, Spark, Swift, Kotlin, TypeScript, Python, Solidity, C++, SQL, Bash/PowerShell) и обоснование — см. [STACK.md](STACK.md).**

### Критерии выбора (в приоритетном порядке)
1. **Безопасность** — память, криптография, изоляция, PQC.
2. **Скорость** — отклик, пропускная способность, латентность.
3. **Графика** — рендер, видео, тяжёлые вычисления.
4. **Дизáйн** — UX, интерфейсы, консистентность.
5. **Масштабируемость** — рост нагрузки, данные, узлы.

Языки подобраны под итог **«лучшие во всём мире»**: каждый — там, где даёт максимум по безопасности, скорости, графике, дизайну и масштабу. **Порядок — лучшие в начале** (Rust → Go → Spark → …), чтобы не переделывать. Полный порядок и обоснование — в [STACK.md](STACK.md).

### Языки (порядок = приоритет, первый подходящий — использовать)
| № | Язык / стек | Где | Зачем для итога |
|---|-------------|-----|-----------------|
| 1 | **Rust** | Видео, поиск, тяжёлые движки, блокчейн-ядро | Безопасность памяти без GC, скорость, надёжность — уровень системного софта мирового класса. |
| 2 | **Go** | API, чаты, микросервисы, авторизация | Простота, масштабируемость, один бинарник — быстрый и предсказуемый деплой. |
| 3 | **Spark** (Scala/Java) | Аналитика, рекомендации, Big Data | Масштаб данных и скорость пакетной обработки на уровне лидеров индустрии. |
| 4 | **Swift** | Нативное iOS | Производительность и дизайн (UIKit/SwiftUI) на уровне лучших iOS-приложений. |
| 5 | **Kotlin** | Нативное Android | Современный стек, интерфейсы и экосистема на уровне лучших Android-приложений. |
| 6 | **TypeScript** | Веб (сайт, приложение) | Только TypeScript (React/Vue); JS и HTML убрать по максимуму, минимально только где неизбежно. |
| 7 | **Python** | ИИ, нейросети, скрипты, автоматизация | Экосистема ML и скорость итераций. |
| 8 | **C++ / C#** | Игры (Unreal / Unity), если будут | Графика и производительность. |
| 9 | **SQL** | БД, транзакции, аналитика | Хранение и запросы данных. |
| 10 | **Bash / PowerShell / YAML** | Инфраструктура, DevOps | Автоматизация и масштабирование. |

Блокчейн: **только Rust** (Solidity не используем — Rust важнее).

### Приоритет 1 (ядро системы — лучшие в начале)
| Стек | Назначение |
|------|------------|
| **Rust** | Видео, поиск, высоконагруженные движки, блокчейн-ядро. |
| **Go** | API, чаты, микросервисы, авторизация. |
| **Spark** (Scala/Java) | Аналитика, рекомендации, Big Data. |

### Приоритет 2 (клиенты)
| Стек | Назначение |
|------|------------|
| **Swift** | Нативное iOS-приложение. |
| **Kotlin** | Нативное Android-приложение. |
| **TypeScript/JavaScript** | Веб-версия (сайт). |

### Приоритет 3 (инфраструктура)
| Стек | Назначение |
|------|------------|
| **Python / Bash / YAML** | Автоматизация, скрипты, DevOps. |

### По направлениям

### Сервер (логика, данные, нагрузка)
| Язык / стек | Назначение |
|-------------|------------|
| **Spark** (Scala/Java) | Аналитика, рекомендации. |
| **Rust** | Видео, поиск, тяжёлые вычисления. |
| **Go** | API, микросервисы, чаты. |

### Веб-сайт (то, что в браузере)
| Язык / стек | Назначение |
|-------------|------------|
| **TypeScript / JavaScript** | Основной язык. |
| **React / Vue** | Фреймворки для интерфейса. |

### Мобильные приложения (телефоны)
| Язык / стек | Назначение |
|-------------|------------|
| **Swift** | iPhone. |
| **Kotlin** | Android. |

### Игры (если будут)
| Язык / стек | Назначение |
|-------------|------------|
| **C++** | Мощные движки (Unreal). |
| **C#** | Unity. |

### ИИ и умные штуки
| Язык / стек | Назначение |
|-------------|------------|
| **Python** | Нейросети, распознавание. |

### Инфраструктура (чтобы работало)
| Язык / стек | Назначение |
|-------------|------------|
| **Bash / Python** | Скрипты. |
| **YAML** | Настройки серверов. |

- Все пользовательские тексты и интерфейсы — мультиязычные; **английский — язык по умолчанию**. Русский не приоритетнее остальных.

**Quantum resistance:** Вся криптография — квантово-устойчивая или с путём миграции. Токены: Dilithium3 (PQC). Пароли: Argon2id. В проде — TLS с PQC hybrid. Blockchain IXI: PQC-подписи и KEM в консенсусе. См. QUANTUM_READINESS.md.

---

## 3. Ecosystem: 4 platforms + other directions

OMNIXIUS is one ecosystem. The **four platforms** are the core; **other directions** (education, startups, media, your other projects) plug into them so everything stays one product with one account.

- **CONNECT** — social, experts, learning, media, community.
- **TRADE & FINANCE** — payments, investments, rewards, subscriptions.
- **REPOSITORIUM** — storage, compute, files, backups.
- **BLOCKCHAIN (IXI)** — tokens, governance, launchpad, staking.

See **ECOSYSTEM.md** for the full map and how to add new directions.

---

## 4. Браузер OMNIXIUS и ИИ OMNIXIUS

- **Браузер OMNIXIUS** — то, что открывается в браузере: сайт `index.html`, приложение `app/` (логин, маркетплейс, почта, заказы, профиль), заготовка веб-клиента на React+TS в `web/`.
- **ИИ OMNIXIUS** — корень своего ИИ: страница в браузере `app/ai.html` (чат), бэкенд `ai/` (Python, FastAPI, эндпоинт `/chat`). Дальше — обучение и подключение своих моделей.

## 5. Documents

- **ARCHITECTURE.md** — этот файл (видение и стек).
- **OMNIXIUS_CHECKLIST.md** — проверка: всё ли из стека и инфраструктуры есть в репозитории.
- **ECOSYSTEM.md** — карта 4 платформ и направлений.
- **PLATFORM_INFRASTRUCTURE.md** — Platform & Infrastructure Division: Docker, Kubernetes, Terraform, CI/CD, Prometheus, Grafana, Kafka.
- **ROADMAP.md** — фазы и приоритеты.

---

## 6. Ten new directions — architecture and roadmap

Единая архитектура под **10 направлений**, которые наращиваются по фазам. Для каждого задано: к какому ядру подключается, какие общие сервисы использует, как стыкуется с остальными.

### 6.1. Общие принципы

- **Один аккаунт** — все направления под одним логином (текущий auth, профиль, KYC где нужно).
- **Единая идентичность** — общий профиль пользователя; SSI/DID (направление 3) со временем даёт верифицируемое «лицо» для всех сервисов.
- **Общие сервисы** — auth (Go), платежи/кошельки (Trade), хранилище (Repositorium), токены/смарт-контракты (IXI), AI (ai/), аналитика (Spark).
- **Данные** — личные данные в Repositorium с контролем доступа; кросс-направления только через явные согласия и API.

### 6.2. Архитектура по направлениям

| # | Направление | Ядро (платформы) | Общие сервисы | Кратко |
|---|-------------|------------------|----------------|--------|
| 1 | **ESG / Зелёная повестка** | Trade, Connect | Wallets, analytics, reporting | Углеродный след, офсеты, зелёные активы в Trade; отчёты и рейтинги в Connect. |
| 2 | **Health / Телемедицина** | Connect, Repositorium | Auth, storage, payments (Trade) | Приёмы, рецепты, хранение медданных в Repo; эксперты и записи в Connect; оплата через Trade. |
| 3 | **SSI / Децентрализованная идентичность** | IXI, Connect | Auth, KYC, Repositorium | Верифицируемые credentials на IXI; один «паспорт» для всех направлений; хранение ключей/доказательств в Repo. |
| 4 | **Embedded finance** | Trade, Connect | Wallets, KYC, risk | Микрокредиты, страховка, рассрочка внутри Marketplace/Connect; расчёты и лимиты через Trade. |
| 5 | **Gig / Локальные услуги** | Connect | Auth, payments (Trade), geo, ratings | Оформление «специалисты рядом»: брони, оплата, отзывы; Connect — профили и заказы, Trade — платежи. |
| 6 | **Creator economy / Подписки** | Connect, Trade, Repositorium | Payments, storage, Media (Connect) | Подписки, донаты, закрытый контент; хранение в Repo, выплаты и токены через Trade/IXI. |
| 7 | **AI-агенты** | AI (ai/), Connect, Trade | Auth, API, chat, payments | Агенты выполняют действия: брони, оплаты, сводки; опираются на текущий AI + API Connect/Trade. |
| 8 | **Privacy-first аналитика** | Repositorium, Spark, Connect | Storage, anonymization, dashboards | Отчёты для бизнеса без персональных данных; Spark + дифференциальная приватность; Repo — сырые данные под контролем. |
| 9 | **Cross-border / Ремиты** | Trade, IXI | Wallets, stablecoins, compliance | Низкая комиссия, быстрые переводы; стейблкоины/IXI; KYC/AML через общий auth и Trade. |
| 10 | **Education credentials / Скиллы on-chain** | IXI, Connect, Learning | Auth, Learning, Repositorium | Верифицируемые дипломы/бейджи на IXI; курсы и эксперты в Connect/Learning; хранение доказательств в Repo. |

### 6.3. Точки интеграции с ядром

- **Connect** — профили, эксперты, заказы, чат, гео, рейтинги, медиа, сообщества. Все направления с «людьми» и услугами опираются на Connect.
- **Trade** — кошельки, платежи, KYC/AML, кредиты, страховки, ремиты, подписки. Всё, что связано с деньгами, идёт через Trade.
- **Repositorium** — хранилище, compute, шифрование, бэкапы. Медданные, контент, логи, аналитика без PII живут здесь.
- **IXI** — консенсус, ZKP, токены, смарт-контракты, credentials, стейкинг. Идентичность, сертификаты, лояльность, стартапы — на блокчейне.

### 6.4. Дорожная карта по фазам

Наращивание без «всё сразу»: сначала архитектура и 2–3 направления, потом следующие.

| Фаза | Направления | Фокус |
|------|-------------|--------|
| **Фаза 1** | AI-агенты (7), Gig/локальные услуги (5), Creator economy (6) | Быстрая польза: агенты поверх текущего AI; оформление гига и подписок на базе Connect + Trade. |
| **Фаза 2** | SSI/DID (3), Embedded finance (4), Cross-border/ремиты (9) | Идентичность и деньги: один паспорт, встроенные финансы, ремиты через Trade/IXI. |
| **Фаза 3** | ESG (1), Health (2), Privacy-аналитика (8), Education credentials (10) | Регуляции и доверие: зелёная отчётность, телемедицина, аналитика без PII, сертификаты on-chain. |

После каждой фазы — обновлять ECOSYSTEM.md и интерфейс (карта направлений), чтобы архитектура и продукт шли в ногу.

**Конкретные шаги Phase A** (AI-агенты, Gig, Creator economy) — в ROADMAP.md, блок «Phase A — шаги».

**Phase B (Identity & money):** шаги B1–B4 в ROADMAP.md — верифицированный профиль (B1), рассрочка по заказу (B2), баланс пользователя stub (B3), API ремитов stub (B4). Полноценные SSI на IXI и платежи через Trade — после Phase 1–2 инфраструктуры.
