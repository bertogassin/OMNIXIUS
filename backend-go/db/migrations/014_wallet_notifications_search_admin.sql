-- §15 Wallet: multi-currency balances, transactions, holds
CREATE TABLE IF NOT EXISTS wallet_balances (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  currency TEXT NOT NULL DEFAULT 'USD',
  amount BIGINT NOT NULL DEFAULT 0,
  hold_amount BIGINT NOT NULL DEFAULT 0,
  updated_at INTEGER DEFAULT (unixepoch()),
  UNIQUE(user_id, currency)
);
CREATE INDEX IF NOT EXISTS idx_wallet_balances_user ON wallet_balances(user_id);

CREATE TABLE IF NOT EXISTS wallet_transactions (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL REFERENCES users(id),
  type TEXT NOT NULL,
  currency TEXT NOT NULL,
  amount BIGINT NOT NULL,
  fee BIGINT NOT NULL DEFAULT 0,
  status TEXT NOT NULL DEFAULT 'completed',
  reference_id TEXT,
  metadata TEXT,
  created_at INTEGER DEFAULT (unixepoch()),
  completed_at INTEGER
);
CREATE INDEX IF NOT EXISTS idx_wallet_tx_user ON wallet_transactions(user_id);
CREATE INDEX IF NOT EXISTS idx_wallet_tx_created ON wallet_transactions(created_at);

CREATE TABLE IF NOT EXISTS wallet_holds (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL REFERENCES users(id),
  order_id INTEGER,
  currency TEXT NOT NULL,
  amount BIGINT NOT NULL,
  expires_at INTEGER NOT NULL,
  created_at INTEGER DEFAULT (unixepoch()),
  released_at INTEGER,
  captured_at INTEGER
);
CREATE INDEX IF NOT EXISTS idx_wallet_holds_user ON wallet_holds(user_id);
CREATE INDEX IF NOT EXISTS idx_wallet_holds_expires ON wallet_holds(expires_at);

-- §16 Notifications: queue, user settings
CREATE TABLE IF NOT EXISTS notifications_queue (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL REFERENCES users(id),
  type TEXT NOT NULL,
  channel TEXT NOT NULL DEFAULT 'websocket',
  title TEXT,
  body TEXT NOT NULL,
  data TEXT,
  status TEXT NOT NULL DEFAULT 'pending',
  attempts INTEGER NOT NULL DEFAULT 0,
  max_attempts INTEGER NOT NULL DEFAULT 3,
  scheduled_for INTEGER DEFAULT (unixepoch()),
  sent_at INTEGER,
  created_at INTEGER DEFAULT (unixepoch())
);
CREATE INDEX IF NOT EXISTS idx_notif_queue_user ON notifications_queue(user_id);
CREATE INDEX IF NOT EXISTS idx_notif_queue_status ON notifications_queue(status);

CREATE TABLE IF NOT EXISTS notifications_user_settings (
  user_id INTEGER PRIMARY KEY REFERENCES users(id),
  email_enabled INTEGER NOT NULL DEFAULT 1,
  push_enabled INTEGER NOT NULL DEFAULT 1,
  updated_at INTEGER DEFAULT (unixepoch())
);

-- §17 Vault search: blind index terms
CREATE TABLE IF NOT EXISTS vault_search_index (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL REFERENCES users(id),
  file_id INTEGER NOT NULL REFERENCES vault_files(id) ON DELETE CASCADE,
  term_hash TEXT NOT NULL,
  weight INTEGER NOT NULL DEFAULT 1,
  created_at INTEGER DEFAULT (unixepoch())
);
CREATE INDEX IF NOT EXISTS idx_vault_search_user_term ON vault_search_index(user_id, term_hash);
CREATE INDEX IF NOT EXISTS idx_vault_search_file ON vault_search_index(file_id);

-- §18 Admin: reports, bans
CREATE TABLE IF NOT EXISTS admin_reports (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  reporter_id INTEGER NOT NULL REFERENCES users(id),
  reported_type TEXT NOT NULL,
  reported_id TEXT NOT NULL,
  reason TEXT NOT NULL,
  description TEXT,
  status TEXT NOT NULL DEFAULT 'pending',
  assigned_to INTEGER REFERENCES users(id),
  resolution TEXT,
  created_at INTEGER DEFAULT (unixepoch()),
  resolved_at INTEGER
);
CREATE INDEX IF NOT EXISTS idx_admin_reports_status ON admin_reports(status);

CREATE TABLE IF NOT EXISTS admin_bans (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL REFERENCES users(id),
  banned_by INTEGER NOT NULL REFERENCES users(id),
  reason TEXT NOT NULL,
  expires_at INTEGER,
  created_at INTEGER DEFAULT (unixepoch()),
  lifted_at INTEGER,
  lifted_by INTEGER REFERENCES users(id)
);
CREATE INDEX IF NOT EXISTS idx_admin_bans_user ON admin_bans(user_id);
CREATE INDEX IF NOT EXISTS idx_admin_bans_expires ON admin_bans(expires_at);
