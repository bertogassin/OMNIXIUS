-- B4 Cross-border / Remittances: stub for transfer requests. Real transfer via Trade/IXI later.
CREATE TABLE IF NOT EXISTS remittances (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  from_user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  to_identifier TEXT NOT NULL,
  amount REAL NOT NULL,
  currency TEXT NOT NULL DEFAULT 'USD',
  status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'completed', 'failed', 'cancelled')),
  created_at INTEGER DEFAULT (unixepoch())
);
CREATE INDEX IF NOT EXISTS idx_remittances_from ON remittances(from_user_id);
CREATE INDEX IF NOT EXISTS idx_remittances_status ON remittances(status);
