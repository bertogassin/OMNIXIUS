-- §15.2 wallet.deposit_addresses; §16.2 notifications.push_tokens
CREATE TABLE IF NOT EXISTS wallet_deposit_addresses (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  currency TEXT NOT NULL,
  address TEXT NOT NULL,
  network TEXT NOT NULL DEFAULT 'mainnet',
  created_at INTEGER DEFAULT (unixepoch()),
  last_used_at INTEGER,
  UNIQUE(currency, address)
);
CREATE INDEX IF NOT EXISTS idx_wallet_deposit_user ON wallet_deposit_addresses(user_id);

CREATE TABLE IF NOT EXISTS notifications_push_tokens (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  device_id INTEGER REFERENCES devices(id) ON DELETE CASCADE,
  token TEXT NOT NULL,
  platform TEXT NOT NULL DEFAULT 'web',
  created_at INTEGER DEFAULT (unixepoch()),
  last_used INTEGER DEFAULT (unixepoch())
);
CREATE INDEX IF NOT EXISTS idx_notif_push_user ON notifications_push_tokens(user_id);
