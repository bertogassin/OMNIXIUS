import Database from 'better-sqlite3';
import { readFileSync } from 'fs';
import { dirname, join } from 'path';
import { fileURLToPath } from 'url';

const __dirname = dirname(fileURLToPath(import.meta.url));
const schema = readFileSync(join(__dirname, 'schema.sql'), 'utf8');
const db = new Database(join(__dirname, 'omnixius.db'));

db.exec(schema);
console.log('DB initialized: backend/db/omnixius.db');
db.close();
