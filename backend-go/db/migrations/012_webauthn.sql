-- WebAuthn (Passkeys) for ARCHITECTURE-V4 Phase 1
-- Credentials: stored per user; credential_json = full webauthn.Credential as JSON
CREATE TABLE IF NOT EXISTS webauthn_credentials (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  credential_id_base64 TEXT NOT NULL UNIQUE,
  credential_json TEXT NOT NULL,
  created_at INTEGER DEFAULT (unixepoch()),
  updated_at INTEGER DEFAULT (unixepoch())
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_webauthn_cred_id ON webauthn_credentials(credential_id_base64);
CREATE INDEX IF NOT EXISTS idx_webauthn_credentials_user ON webauthn_credentials(user_id);

-- Session store for WebAuthn challenge (register/login begin → complete). session_data = JSON of webauthn.SessionData
CREATE TABLE IF NOT EXISTS webauthn_sessions (
  id TEXT PRIMARY KEY,
  session_data TEXT NOT NULL,
  created_at INTEGER DEFAULT (unixepoch())
);
CREATE INDEX IF NOT EXISTS idx_webauthn_sessions_created ON webauthn_sessions(created_at);
