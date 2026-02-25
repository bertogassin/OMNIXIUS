-- Vault module (ARCHITECTURE-V4): folders and files per user
CREATE TABLE IF NOT EXISTS vault_folders (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  parent_id INTEGER REFERENCES vault_folders(id) ON DELETE SET NULL,
  created_at INTEGER DEFAULT (unixepoch()),
  updated_at INTEGER DEFAULT (unixepoch())
);
CREATE INDEX IF NOT EXISTS idx_vault_folders_user ON vault_folders(user_id);
CREATE INDEX IF NOT EXISTS idx_vault_folders_parent ON vault_folders(parent_id);

CREATE TABLE IF NOT EXISTS vault_files (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  size_bytes INTEGER NOT NULL,
  mime_type TEXT,
  storage_path TEXT NOT NULL,
  folder_id INTEGER REFERENCES vault_folders(id) ON DELETE SET NULL,
  created_at INTEGER DEFAULT (unixepoch()),
  updated_at INTEGER DEFAULT (unixepoch())
);
CREATE INDEX IF NOT EXISTS idx_vault_files_user ON vault_files(user_id);
CREATE INDEX IF NOT EXISTS idx_vault_files_folder ON vault_files(folder_id);
