# OMNIXIUS — Ecosystem map

Одна экосистема: **what’s live now** (Mail, Marketplace, Shop, Profile, AI) и **полное видение** (Connect, Trade & Finance, Repositorium, Blockchain и др.). Один аккаунт, один продукт. **Уровень:** начало; карта и правила — здесь и в ROADMAP.

---

## 1. What’s live now

| Direction | What it is |
|-----------|------------|
| **Mail** | Internal mail: conversations and messages between users. |
| **Marketplace** | Listings: browse, search, filters, categories. |
| **Shop** | Buy and sell: orders, “Contact seller” (opens Mail), order status. |
| **Profile** | Account, avatar, my orders (buyer/seller). |
| **AI** | Корень своего ИИ: чат в браузере (`app/ai.html`), бэкенд `ai/` (Python). |

Один фронт, один API (backend-go). Новые направления подключаются по этой карте.

---

## 2. Full ecosystem (Connect, Trade, Repositorium, Blockchain, and more)

The “4” are part of the bigger picture; new directions keep being added.

| Platform / area | Role |
|-----------------|------|
| **OMNIXIUS CONNECT** | Social, experts, marketplace, chat, video. People and services. |
| **OMNIXIUS TRADE & FINANCE** | Exchange, wallets, copy trading, investments. Money and assets. |
| **OMNIXIUS REPOSITORIUM** | Decentralized storage and compute. Data and infrastructure. |
| **OMNIXIUS BLOCKCHAIN (IXI)** | Base layer: consensus, ZKP, tokenomics, buyback fund. Trust and economy. |
| **… and more** | Education, media, startups, your other projects — one account. |

---

## 3. How new directions connect

Additional projects and verticals plug into the same ecosystem so everything stays **one product**, not separate apps.

- **Where they attach**
  - **Connect** — anything social, expert, learning, media, community, support.
  - **Trade & Finance** — anything with payments, investments, rewards, subscriptions.
  - **Repositorium** — anything with files, storage, compute, backups, CDN.
  - **Blockchain (IXI)** — anything with tokens, governance, startups, loyalty, staking.

- **Cross-platform directions** (span 2+ platforms)
  - **Education / Learning** — courses and experts (Connect) + payments and rewards (Trade) + storage for materials (Repositorium).
  - **Startups / Launchpad** — fundraising and tokens (Blockchain + Trade) + community and experts (Connect).
  - **Media / Content** — creators and subscriptions (Connect) + payments (Trade) + storage and delivery (Repositorium).
  - **Rewards / Loyalty** — activity and referrals (Connect) + payouts and IXI (Trade + Blockchain) + optional storage rewards (Repositorium).

- **Single account**
  - One login for: Connect, Trade, Repositorium, Blockchain and all “other directions” built on top. No separate accounts per project.

- **Shared stack**
  - Same tech policy (SPARK, RUST, C++, SWIFT, GO), same security (PQC, Argon2id), same docs (ARCHITECTURE, ROADMAP). New directions use the same API and auth.

---

## 4. Adding a new direction

1. Decide which core platform(s) it belongs to (Connect / Trade / Repositorium / Blockchain).
2. Reuse: auth, users, wallets, storage, or blockchain from existing services.
3. Document it here under “Other directions” with 1–2 lines and which platforms it uses.
4. Link it from the site (Ecosystem page or Architecture) so the map stays clear.

---

## 5. Other directions (to expand)

*Full architecture for the 10 new directions: ARCHITECTURE.md §6.*

| Direction | Platforms | Note |
|-----------|-----------|------|
| Mail (live) | Connect | Conversations, messages. |
| Marketplace & Shop (live) | Connect | Listings, orders, contact seller. |
| Learning / Courses | Connect, Trade, Repositorium | Experts, payments, storage. |
| Startups / IXI launchpad | Blockchain, Trade, Connect | Tokens, fundraising, community. |
| Media / Content | Connect, Trade, Repositorium | Creators, subscriptions, storage. |
| Rewards / Loyalty | Connect, Trade, Blockchain | Activity, referrals, IXI payouts. |
| **ESG / Green** | Trade, Connect | Carbon footprint, offsets, green assets. |
| **Health / Telemedicine** | Connect, Repositorium, Trade | Appointments, health data, payments. |
| **SSI / Decentralized identity** | IXI, Connect | Verifiable credentials, one passport. |
| **Embedded finance** | Trade, Connect | Microloans, insurance, installments. |
| **Gig / Local services** | Connect, Trade | Bookings, payments, ratings. |
| **Creator economy** | Connect, Trade, Repositorium | Subscriptions, tips, exclusive content. |
| **AI agents** | AI, Connect, Trade | Actions: book, pay, summarize. |
| **Privacy-first analytics** | Repositorium, Spark | Business insights without PII. |
| **Cross-border / Remittances** | Trade, IXI | Low-fee transfers, stablecoins. |
| **Education credentials** | IXI, Connect, Learning | Verifiable diplomas, badges on-chain. |

---

Use this file as the **single map** of OMNIXIUS: what’s live (Mail, Marketplace, Shop) + full ecosystem (Connect, Trade, Repositorium, Blockchain) + everything else you add.
