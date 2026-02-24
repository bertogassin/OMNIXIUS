# OMNIXIUS — Architecture & Vision

Global technology ecosystem: social network, professional platform, exchange, decentralized cloud, IXI blockchain, AI analytics, and investment tools.

**Goal:** Self-sufficient ecosystem — communication, earning, investment, learning, and asset management in one application, without dependence on centralized corporations.

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

### Приоритет 1 (ядро системы)
| Стек | Назначение |
|------|------------|
| **Rust** | Видео, поиск, высоконагруженные движки. |
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
| **Go** | API, микросервисы, чаты. |
| **Rust** | Видео, поиск, тяжёлые вычисления. |
| **Spark** (Scala/Java) | Аналитика, рекомендации. |

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
