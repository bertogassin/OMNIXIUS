-- Professionals search: profession, online, distance. Used for "find professional" by map.
ALTER TABLE users ADD COLUMN profession_id TEXT;
ALTER TABLE users ADD COLUMN lat REAL;
ALTER TABLE users ADD COLUMN lng REAL;
ALTER TABLE users ADD COLUMN last_seen_at INTEGER;
CREATE INDEX IF NOT EXISTS idx_users_profession ON users(profession_id);
CREATE INDEX IF NOT EXISTS idx_users_last_seen ON users(last_seen_at);
