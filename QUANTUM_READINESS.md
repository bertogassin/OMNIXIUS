# OMNIXIUS — Quantum resistance

**Requirement:** Site, blockchain (IXI), internal mail, and all components must be **quantum-resistant**. No rush; we have what we need — implement to the maximum.

---

## What is in place

### 1. Auth tokens (post-quantum signatures)

- **Dilithium3** (NIST PQC) is used for auth tokens instead of JWT/HMAC or RSA/ECDSA.
- Token = signed payload (userID + expiry) with Dilithium3; verification uses the public key only.
- **Production:** set env vars:
  - `DILITHIUM_PUBLIC_KEY` — base64-encoded public key
  - `DILITHIUM_PRIVATE_KEY` — base64-encoded private key  
  Generate once with a small script that calls `pqc.GenerateKey()` and prints base64; store keys securely. If not set, the server generates ephemeral keys at startup (tokens invalid after restart).

### 2. Passwords (quantum-resistant KDF)

- **Argon2id** is used for new passwords (register, password reset).
- Argon2id is a memory-hard KDF; quantum attackers get at most a square-root speedup on exhaustive search. Parameters: time=3, memory=64 MiB (configurable via `ARGON2_MEMORY`), threads=2.
- **Legacy:** Existing bcrypt hashes still work at login; new and reset passwords are stored as Argon2id.

### 3. TLS (production)

- Use **TLS 1.3** in production. Prefer hybrid post-quantum key agreement (e.g. Kyber) when your stack supports it (e.g. Cloudflare, modern OpenSSL).
- Enforce HTTPS; no plain HTTP for auth or PII.

### 4. Internal mail and data

- Messages are stored in the DB; transport must be over HTTPS (and thus protected by TLS).
- Optional next step: E2E encryption of message bodies with a PQC KEM (e.g. Kyber) for key agreement — document and add when required.

### 5. Blockchain (IXI)

- IXI is defined as the ecosystem blockchain (see ARCHITECTURE.md).
- For **quantum resistance**, IXI must use post-quantum primitives:
  - **Signatures:** ML-DSA (Dilithium) or equivalent PQC signature scheme instead of ECDSA.
  - **Key agreement / KEM:** ML-KEM (Kyber) or equivalent where key exchange is used.
- Implementation belongs in the IXI node/consensus code (Rust/C++); the API and site use Dilithium for auth as above.

---

## Checklist

| Component        | Status | Notes |
|-----------------|--------|--------|
| Auth tokens     | Done   | Dilithium3 (backend-go/pqc) |
| Passwords       | Done   | Argon2id; bcrypt legacy supported |
| TLS / HTTPS     | Todo   | Enforce in production; prefer PQC hybrid where available |
| Internal mail   | Partial| Over HTTPS; optional E2E with PQC later |
| Blockchain IXI  | Todo   | PQC signatures and KEM in node/consensus |
| Site (static)   | N/A    | Served over HTTPS |

---

## References

- NIST PQC: FIPS 204 (ML-DSA), FIPS 203 (ML-KEM).
- Argon2: RFC 9106.
- IETF: draft for ML-DSA in JOSE/COSE (JWT with PQC).
