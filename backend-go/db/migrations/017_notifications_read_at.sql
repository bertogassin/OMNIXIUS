-- Ensure table exists (if 014 was skipped) then add read_at
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
  created_at INTEGER DEFAULT (unixepoch()),
  read_at INTEGER
);
CREATE INDEX IF NOT EXISTS idx_notif_queue_user ON notifications_queue(user_id);
CREATE INDEX IF NOT EXISTS idx_notif_queue_status ON notifications_queue(status);
ALTER TABLE notifications_queue ADD COLUMN read_at INTEGER;
