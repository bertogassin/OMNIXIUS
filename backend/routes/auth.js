import { Router } from 'express';
import bcrypt from 'bcryptjs';
import jwt from 'jsonwebtoken';
import crypto from 'crypto';
import { db } from '../db/index.js';
import { config } from '../config.js';
import { validateRegister, validateLogin } from '../middleware/validate.js';

const router = Router();
const loginAttempts = new Map();

function getLoginAttempts(ip) {
  const now = Date.now();
  const entry = loginAttempts.get(ip) || { count: 0, resetAt: now + config.loginWindowMs };
  if (now > entry.resetAt) entry.count = 0;
  return entry;
}

router.post('/register', validateRegister, (req, res) => {
  const { email, password, name } = req.body;
  const existing = db.prepare('SELECT id FROM users WHERE email = ?').get(email);
  if (existing) return res.status(409).json({ error: 'Email уже зарегистрирован' });

  const hash = bcrypt.hashSync(password, config.bcryptRounds);
  const verifyToken = crypto.randomBytes(32).toString('hex');
  const now = Math.floor(Date.now() / 1000);
  db.prepare(
    'INSERT INTO users (email, password_hash, name, role, email_verify_token) VALUES (?, ?, ?, ?, ?)'
  ).run(email, hash, name, 'user', verifyToken);
  const user = db.prepare('SELECT id, email, role, name, created_at FROM users WHERE email = ?').get(email);
  // TODO: send email with verify link using verifyToken
  const token = jwt.sign({ userId: user.id }, config.jwtSecret, { expiresIn: config.jwtExpires });
  res.status(201).json({ user: { id: user.id, email: user.email, role: user.role, name: user.name }, token });
});

router.post('/login', validateLogin, (req, res) => {
  const ip = req.ip || req.connection?.remoteAddress || 'unknown';
  const entry = getLoginAttempts(ip);
  if (entry.count >= config.maxLoginAttempts) {
    return res.status(429).json({ error: 'Слишком много попыток входа. Попробуйте позже.' });
  }

  const { email, password } = req.body;
  const user = db.prepare('SELECT id, email, password_hash, role, name, avatar_path, email_verified FROM users WHERE email = ?').get(email);
  if (!user || !bcrypt.compareSync(password, user.password_hash)) {
    entry.count++;
    loginAttempts.set(ip, entry);
    return res.status(401).json({ error: 'Неверный email или пароль' });
  }
  loginAttempts.delete(ip);
  const token = jwt.sign({ userId: user.id }, config.jwtSecret, { expiresIn: config.jwtExpires });
  delete user.password_hash;
  res.json({ user, token });
});

router.post('/logout', (req, res) => {
  res.json({ ok: true });
});

// Confirm email (simple: by token in query)
router.get('/confirm-email', (req, res) => {
  const token = req.query.token;
  if (!token) return res.status(400).json({ error: 'Нужен токен' });
  const r = db.prepare('UPDATE users SET email_verified = 1, email_verify_token = NULL WHERE email_verify_token = ?').run(token);
  if (r.changes === 0) return res.status(400).json({ error: 'Недействительный токен' });
  res.json({ ok: true });
});

// Forgot password: set reset token (in production send email with link)
router.post('/forgot-password', (req, res) => {
  const email = (req.body?.email || '').trim().toLowerCase();
  if (!email) return res.status(400).json({ error: 'Укажите email' });
  const user = db.prepare('SELECT id FROM users WHERE email = ?').get(email);
  if (!user) return res.json({ ok: true }); // don't leak existence
  const resetToken = crypto.randomBytes(32).toString('hex');
  const expires = Math.floor(Date.now() / 1000) + 3600; // 1h
  db.prepare('UPDATE users SET reset_token = ?, reset_token_expires = ? WHERE id = ?').run(resetToken, expires, user.id);
  // TODO: send email with reset link
  res.json({ ok: true });
});

router.post('/reset-password', (req, res) => {
  const { token, password } = req.body || {};
  if (!token || !password || password.length < 8) return res.status(400).json({ error: 'Токен и пароль (≥8 символов) обязательны' });
  const user = db.prepare('SELECT id FROM users WHERE reset_token = ? AND reset_token_expires > ?').get(token, Math.floor(Date.now() / 1000));
  if (!user) return res.status(400).json({ error: 'Недействительный или просроченный токен' });
  const hash = bcrypt.hashSync(password, config.bcryptRounds);
  db.prepare('UPDATE users SET password_hash = ?, reset_token = NULL, reset_token_expires = NULL WHERE id = ?').run(hash, user.id);
  res.json({ ok: true });
});

export default router;
