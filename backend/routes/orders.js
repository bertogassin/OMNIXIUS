import { Router } from 'express';
import { authRequired } from '../middleware/auth.js';
import { db } from '../db/index.js';

const router = Router();

router.post('/', authRequired, (req, res) => {
  const { product_id } = req.body || {};
  const productId = Number(product_id);
  if (!productId) return res.status(400).json({ error: 'Укажите product_id' });
  const product = db.prepare('SELECT id, user_id FROM products WHERE id = ?').get(productId);
  if (!product) return res.status(404).json({ error: 'Товар не найден' });
  if (product.user_id === req.user.id) return res.status(400).json({ error: 'Нельзя заказать свой товар' });
  db.prepare(
    'INSERT INTO orders (product_id, buyer_id, seller_id, status) VALUES (?, ?, ?, ?)'
  ).run(productId, req.user.id, product.user_id, 'pending');
  const order = db.prepare('SELECT * FROM orders WHERE id = last_insert_rowid()').get();
  res.status(201).json(order);
});

router.patch('/:id', authRequired, (req, res) => {
  const order = db.prepare('SELECT * FROM orders WHERE id = ?').get(req.params.id);
  if (!order) return res.status(404).json({ error: 'Заказ не найден' });
  const canUpdate = order.buyer_id === req.user.id || order.seller_id === req.user.id || req.user.role === 'admin';
  if (!canUpdate) return res.status(403).json({ error: 'Нет прав' });
  const { status } = req.body || {};
  if (!['pending', 'confirmed', 'completed', 'cancelled'].includes(status)) return res.status(400).json({ error: 'Некорректный статус' });
  db.prepare('UPDATE orders SET status = ?, updated_at = unixepoch() WHERE id = ?').run(status, req.params.id);
  res.json(db.prepare('SELECT * FROM orders WHERE id = ?').get(req.params.id));
});

router.get('/my', authRequired, (req, res) => {
  const list = db.prepare(`
    SELECT o.*, p.title, p.price, p.image_path,
           (SELECT name FROM users WHERE id = o.seller_id) as seller_name,
           (SELECT name FROM users WHERE id = o.buyer_id) as buyer_name
    FROM orders o
    JOIN products p ON p.id = o.product_id
    WHERE o.buyer_id = ? OR o.seller_id = ?
    ORDER BY o.created_at DESC
  `).all(req.user.id, req.user.id);
  res.json(list);
});

export default router;
