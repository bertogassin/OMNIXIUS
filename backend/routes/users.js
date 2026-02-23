import { Router } from 'express';
import { authRequired } from '../middleware/auth.js';
import { db } from '../db/index.js';

const router = Router();

router.get('/me', authRequired, (req, res) => {
  const u = db.prepare(
    'SELECT id, email, role, name, avatar_path, email_verified, created_at FROM users WHERE id = ?'
  ).get(req.user.id);
  if (!u) return res.status(404).json({ error: 'Пользователь не найден' });
  res.json(u);
});

router.patch('/me', authRequired, (req, res) => {
  const { name } = req.body || {};
  const nameSafe = typeof name === 'string' ? name.trim().slice(0, 200) : null;
  if (nameSafe !== null) {
    db.prepare('UPDATE users SET name = ?, updated_at = unixepoch() WHERE id = ?').run(nameSafe, req.user.id);
  }
  const u = db.prepare('SELECT id, email, role, name, avatar_path, email_verified, created_at FROM users WHERE id = ?').get(req.user.id);
  res.json(u);
});

router.get('/me/orders', authRequired, (req, res) => {
  const asBuyer = db.prepare(`
    SELECT o.id, o.status, o.created_at, p.title, p.price, p.image_path, u.name as seller_name
    FROM orders o
    JOIN products p ON p.id = o.product_id
    JOIN users u ON u.id = o.seller_id
    WHERE o.buyer_id = ?
    ORDER BY o.created_at DESC
  `).all(req.user.id);
  const asSeller = db.prepare(`
    SELECT o.id, o.status, o.created_at, p.title, p.price, p.image_path, u.name as buyer_name
    FROM orders o
    JOIN products p ON p.id = o.product_id
    JOIN users u ON u.id = o.buyer_id
    WHERE o.seller_id = ?
    ORDER BY o.created_at DESC
  `).all(req.user.id);
  res.json({ asBuyer, asSeller });
});

export default router;
