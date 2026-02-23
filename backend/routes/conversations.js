import { Router } from 'express';
import { authRequired } from '../middleware/auth.js';
import { db } from '../db/index.js';

const router = Router();

// List my conversations (with last message and unread count)
router.get('/', authRequired, (req, res) => {
  const list = db.prepare(`
    SELECT c.id, c.product_id, c.updated_at,
           (SELECT body FROM messages WHERE conversation_id = c.id ORDER BY created_at DESC LIMIT 1) as last_message,
           (SELECT COUNT(*) FROM messages m WHERE m.conversation_id = c.id AND m.sender_id != ? AND m.read_at IS NULL) as unread
    FROM conversations c
    JOIN conversation_participants cp ON cp.conversation_id = c.id
    WHERE cp.user_id = ?
    ORDER BY c.updated_at DESC
  `).all(req.user.id, req.user.id);
  const withOther = list.map(c => {
    const other = db.prepare(`
      SELECT u.id, u.name, u.email FROM users u
      JOIN conversation_participants cp ON cp.user_id = u.id
      WHERE cp.conversation_id = ? AND u.id != ?
    `).get(c.id, req.user.id);
    const product = c.product_id ? db.prepare('SELECT id, title FROM products WHERE id = ?').get(c.product_id) : null;
    return { ...c, other, product };
  });
  res.json(withOther);
});

// Get or create conversation with user (and optional product)
router.post('/', authRequired, (req, res) => {
  const { user_id, product_id } = req.body || {};
  const otherId = Number(user_id);
  const productId = product_id ? Number(product_id) : null;
  if (!otherId || otherId === req.user.id) return res.status(400).json({ error: 'Укажите user_id другого пользователя' });
  const other = db.prepare('SELECT id FROM users WHERE id = ?').get(otherId);
  if (!other) return res.status(404).json({ error: 'Пользователь не найден' });

  let conv = db.prepare(`
    SELECT c.id FROM conversations c
    JOIN conversation_participants cp1 ON cp1.conversation_id = c.id AND cp1.user_id = ?
    JOIN conversation_participants cp2 ON cp2.conversation_id = c.id AND cp2.user_id = ?
    WHERE (c.product_id IS NULL AND ? IS NULL) OR (c.product_id = ?)
  `).get(req.user.id, otherId, productId, productId);

  if (!conv) {
    db.prepare('INSERT INTO conversations (product_id) VALUES (?)').run(productId);
    const id = db.prepare('SELECT last_insert_rowid() as id').get().id;
    db.prepare('INSERT INTO conversation_participants (conversation_id, user_id) VALUES (?, ?), (?, ?)').run(id, req.user.id, id, otherId);
    conv = { id };
  }
  const full = db.prepare('SELECT * FROM conversations WHERE id = ?').get(conv.id);
  res.json(full);
});

export default router;
