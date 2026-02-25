-- B3 Trade stub: user balance (internal units) for subscriptions, rewards, payouts.
CREATE TABLE IF NOT EXISTS user_balances (
  user_id INTEGER PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
  balance REAL NOT NULL DEFAULT 0,
  updated_at INTEGER DEFAULT (unixepoch())
);
