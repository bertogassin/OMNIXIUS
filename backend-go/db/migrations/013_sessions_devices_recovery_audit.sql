-- §1.1 Auth: sessions, devices; §1.2 recovery; §1.9 Audit (doc v4.0)
-- Sessions: one row per login; token references session_id
CREATE TABLE IF NOT EXISTS sessions (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  device_name TEXT NOT NULL DEFAULT 'unknown',
  created_at INTEGER DEFAULT (unixepoch()),
  expires_at INTEGER NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_sessions_expires_at ON sessions(expires_at);

-- Devices: trusted devices (e.g. for recovery flow / device list)
CREATE TABLE IF NOT EXISTS devices (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  public_key_or_credential_id TEXT,
  last_used INTEGER,
  created_at INTEGER DEFAULT (unixepoch())
);
CREATE INDEX IF NOT EXISTS idx_devices_user_id ON devices(user_id);

-- Recovery: store hash of Master Recovery Key (BIP-39 phrase) for verify
CREATE TABLE IF NOT EXISTS user_recovery (
  user_id INTEGER PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
  recovery_hash TEXT NOT NULL,
  created_at INTEGER DEFAULT (unixepoch())
);

-- Audit log (§1.9)
CREATE TABLE IF NOT EXISTS audit_log (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER REFERENCES users(id),
  action TEXT NOT NULL,
  resource TEXT NOT NULL,
  resource_id TEXT,
  old_value TEXT,
  new_value TEXT,
  ip TEXT,
  user_agent TEXT,
  created_at INTEGER DEFAULT (unixepoch())
);
CREATE INDEX IF NOT EXISTS idx_audit_log_user_id ON audit_log(user_id);
CREATE INDEX IF NOT EXISTS idx_audit_log_action ON audit_log(action);
CREATE INDEX IF NOT EXISTS idx_audit_log_created_at ON audit_log(created_at);
