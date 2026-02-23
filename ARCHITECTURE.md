# OMNIXIUS — Architecture & Vision

Global technology ecosystem: social network, professional platform, exchange, decentralized cloud, IXI blockchain, AI analytics, and investment tools.

**Goal:** Self-sufficient ecosystem — communication, earning, investment, learning, and asset management in one application, without dependence on centralized corporations.

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

## 2. Technology stack (only)

**Allowed languages and runtimes: SPARK, RUST, C++, SWIFT, GO. No other languages in the codebase.**

| Layer | Technologies |
|-------|--------------|
| Backend / API / Services | **Go**, **Rust**, **C++** |
| Data / Analytics | **Apache Spark** (or SPARK as designated) |
| Mobile (Apple) | **Swift** (iOS) |
| Blockchain / low-level | **Rust**, **C++** |

- No JavaScript/TypeScript, Kotlin, Python, etc. in core services.
- Frontend: all user-facing content is multilingual; **English is the fallback language**. UI and APIs support all languages; Russian is not prioritized over others.

**Quantum resistance:** All crypto must be quantum-resistant or on a clear migration path. Auth tokens: Dilithium3 (PQC). Passwords: Argon2id. TLS in production with PQC hybrid where available. Blockchain IXI: PQC signatures and KEM in consensus. See QUANTUM_READINESS.md.

---

## 3. Documents

- **ARCHITECTURE.md** — this file (vision and architecture).
- **ROADMAP.md** — phases and priorities.
