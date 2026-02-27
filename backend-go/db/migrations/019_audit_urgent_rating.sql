-- Audit log, urgent orders, user rating (WHAT-WE-TAKE)
CREATE TABLE IF NOT EXISTS audit_log (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER REFERENCES users(id),
  action TEXT NOT NULL,
  entity_type TEXT NOT NULL,
  entity_id TEXT,
  details TEXT,
  ip_address TEXT,
  created_at INTEGER DEFAULT (unixepoch())
);
CREATE INDEX IF NOT EXISTS idx_audit_user ON audit_log(user_id);
CREATE INDEX IF NOT EXISTS idx_audit_created ON audit_log(created_at);

ALTER TABLE orders ADD COLUMN urgent INTEGER DEFAULT 0;
CREATE INDEX IF NOT EXISTS idx_orders_urgent ON orders(urgent);

ALTER TABLE users ADD COLUMN rating_avg REAL;
ALTER TABLE users ADD COLUMN rating_count INTEGER DEFAULT 0;
