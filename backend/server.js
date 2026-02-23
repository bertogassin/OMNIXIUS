import express from 'express';
import cors from 'cors';
import rateLimit from 'express-rate-limit';
import path from 'path';
import { fileURLToPath } from 'url';
import multer from 'multer';
import fs from 'fs';
import { config } from './config.js';
import { db } from './db/index.js';
import { authRequired } from './middleware/auth.js';
import authRoutes from './routes/auth.js';
import usersRoutes from './routes/users.js';
import productsRoutes from './routes/products.js';
import ordersRoutes from './routes/orders.js';
import conversationsRoutes from './routes/conversations.js';
import messagesRoutes from './routes/messages.js';

const __dirname = path.dirname(fileURLToPath(import.meta.url));

const dbPath = path.join(__dirname, 'db', 'omnixius.db');
if (!fs.existsSync(dbPath)) {
  console.log('DB не найден. Выполни: npm run init-db');
  process.exit(1);
}

const app = express();
const limiter = rateLimit({
  windowMs: 15 * 60 * 1000,
  max: 200,
  message: { error: 'Слишком много запросов' },
});
app.use(limiter);
app.use(cors({ origin: true, credentials: true }));
app.use(express.json());

const uploadsDir = path.join(__dirname, config.uploadDir);
[path.join(uploadsDir, 'products'), path.join(uploadsDir, 'avatars')].forEach(dir => {
  if (!fs.existsSync(dir)) fs.mkdirSync(dir, { recursive: true });
});

const storage = multer.diskStorage({
  destination: (req, file, cb) => {
    const sub = file.fieldname === 'avatar' ? 'avatars' : 'products';
    cb(null, path.join(uploadsDir, sub));
  },
  filename: (req, file, cb) => {
    const ext = file.mimetype === 'image/jpeg' ? '.jpg' : file.mimetype === 'image/png' ? '.png' : '.webp';
    cb(null, `${Date.now()}-${Math.random().toString(36).slice(2)}${ext}`);
  },
});
const upload = multer({
  storage,
  limits: { fileSize: config.maxFileSize },
  fileFilter: (req, file, cb) => {
    if (config.allowedImageTypes.includes(file.mimetype)) cb(null, true);
    else cb(new Error('Только JPEG, PNG, WebP'));
  },
});

app.use(`/${config.uploadDir}`, express.static(uploadsDir));

app.use('/api/auth', authRoutes);
app.use('/api/users', usersRoutes);
app.use('/api/orders', ordersRoutes);
app.use('/api/conversations', conversationsRoutes);
app.use('/api/messages', messagesRoutes);

app.post('/api/users/me/avatar', authRequired, upload.single('avatar'), (req, res) => {
  if (!req.file) return res.status(400).json({ error: 'Файл не загружен' });
  const rel = path.relative(uploadsDir, req.file.path).replace(/\\/g, '/');
  db.prepare('UPDATE users SET avatar_path = ?, updated_at = unixepoch() WHERE id = ?').run(rel, req.user.id);
  res.json({ avatar_path: rel });
});

app.use('/api/products', (req, res, next) => {
  if (req.method === 'POST' || req.method === 'PATCH') return upload.single('image')(req, res, next);
  next();
}, productsRoutes);

app.use((err, req, res, next) => {
  if (err instanceof multer.MulterError && err.code === 'LIMIT_FILE_SIZE') {
    return res.status(400).json({ error: 'Файл слишком большой' });
  }
  res.status(500).json({ error: err.message || 'Ошибка сервера' });
});

const PORT = config.port;
app.listen(PORT, () => console.log(`OMNIXIUS API: http://localhost:${PORT}`));
