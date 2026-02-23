import { Router } from 'express';
import { authRequired } from '../middleware/auth.js';
import { validateMessage } from '../middleware/validate.js';
import { db } from '../db/index.js';

const router = Router();

router.get('/conversation/:id', authRequired, (req, res) => {
  const participant = db.prepare('SELECT 1 FROM conversation_participants WHERE conversation_id = ? AND user_id = ?').get(req.params.id, req.user.id);
  if (!participant) return res.status(403).json({ error: 'Нет доступа к диалогу' });
  const messages = db.prepare(`
    SELECT m.id, m.sender_id, m.body, m.read_at, m.created_at, u.name as sender_name
    FROM messages m
    JOIN users u ON u.id = m.sender_id
    WHERE m.conversation_id = ?
    ORDER BY m.created_at ASC
  `).all(req.params.id);
  res.json(messages);
});

router.post('/conversation/:id', authRequired, validateMessage, (req, res) => {
  const participant = db.prepare('SELECT 1 FROM conversation_participants WHERE conversation_id = ? AND user_id = ?').get(req.params.id, req.user.id);
  if (!participant) return res.status(403).json({ error: 'Нет доступа к диалогу' });
  db.prepare('INSERT INTO messages (conversation_id, sender_id, body) VALUES (?, ?, ?)').run(req.params.id, req.user.id, req.body.body);
  db.prepare('UPDATE conversations SET updated_at = unixepoch() WHERE id = ?').run(req.params.id);
  const msg = db.prepare('SELECT * FROM messages WHERE id = last_insert_rowid()').get();
  res.status(201).json(msg);
});

router.post('/:id/read', authRequired, (req, res) => {
  const msg = db.prepare('SELECT id, conversation_id, sender_id FROM messages WHERE id = ?').get(req.params.id);
  if (!msg) return res.status(404).json({ error: 'Сообщение не найдено' });
  if (msg.sender_id === req.user.id) return res.json({ ok: true });
  const participant = db.prepare('SELECT 1 FROM conversation_participants WHERE conversation_id = ? AND user_id = ?').get(msg.conversation_id, req.user.id);
  if (!participant) return res.status(403).json({ error: 'Нет доступа' });
  db.prepare('UPDATE messages SET read_at = unixepoch() WHERE id = ? AND read_at IS NULL').run(req.params.id);
  res.json({ ok: true });
});

export default router;
