import Database from 'better-sqlite3';
import { dirname, join } from 'path';
import { fileURLToPath } from 'url';

const __dirname = dirname(fileURLToPath(import.meta.url));
const dbPath = join(__dirname, 'omnixius.db');

export const db = new Database(dbPath);

// Enable foreign keys
db.pragma('foreign_keys = ON');
