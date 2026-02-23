import { Router } from 'express';
import { authRequired, roleRequired } from '../middleware/auth.js';
import { validateProduct } from '../middleware/validate.js';
import { db } from '../db/index.js';

const router = Router();

// List: search by title, category, location; filter price, date
router.get('/', (req, res) => {
  const q = (req.query.q || '').trim();
  const category = (req.query.category || '').trim();
  const location = (req.query.location || '').trim();
  const minPrice = req.query.minPrice != null ? Number(req.query.minPrice) : null;
  const maxPrice = req.query.maxPrice != null ? Number(req.query.maxPrice) : null;
  const sort = req.query.sort === 'date' ? 'created_at DESC' : 'created_at DESC';

  let sql = `
    SELECT p.id, p.title, p.price, p.category, p.location, p.image_path, p.created_at,
           u.id as seller_id, u.name as seller_name
    FROM products p
    JOIN users u ON u.id = p.user_id
    WHERE 1=1
  `;
  const params = [];
  if (q) { sql += ' AND (p.title LIKE ? OR p.description LIKE ?)'; params.push(`%${q}%`, `%${q}%`); }
  if (category) { sql += ' AND p.category = ?'; params.push(category); }
  if (location) { sql += ' AND p.location LIKE ?'; params.push(`%${location}%`); }
  if (minPrice != null && !isNaN(minPrice)) { sql += ' AND p.price >= ?'; params.push(minPrice); }
  if (maxPrice != null && !isNaN(maxPrice)) { sql += ' AND p.price <= ?'; params.push(maxPrice); }
  sql += ` ORDER BY p.${sort}`;
  const list = db.prepare(sql).all(...params);
  res.json(list);
});

router.get('/categories', (req, res) => {
  const rows = db.prepare('SELECT DISTINCT category FROM products ORDER BY category').all();
  res.json(rows.map(r => r.category));
});

router.get('/:id', (req, res) => {
  const p = db.prepare(`
    SELECT p.*, u.id as seller_id, u.name as seller_name, u.email as seller_email
    FROM products p
    JOIN users u ON u.id = p.user_id
    WHERE p.id = ?
  `).get(req.params.id);
  if (!p) return res.status(404).json({ error: 'Товар не найден' });
  res.json(p);
});

router.post('/', authRequired, roleRequired('user', 'seller', 'admin'), validateProduct, (req, res) => {
  const { title, description, price, category, location } = req.body;
  const image_path = req.file ? `products/${req.file.filename}` : null;
  db.prepare(
    'INSERT INTO products (user_id, title, description, price, category, location, image_path) VALUES (?, ?, ?, ?, ?, ?, ?)'
  ).run(req.user.id, title, description, price, category, location, image_path);
  const row = db.prepare('SELECT * FROM products WHERE id = last_insert_rowid()').get();
  res.status(201).json(row);
});

router.patch('/:id', authRequired, validateProduct, (req, res) => {
  const product = db.prepare('SELECT user_id FROM products WHERE id = ?').get(req.params.id);
  if (!product) return res.status(404).json({ error: 'Товар не найден' });
  if (product.user_id !== req.user.id && req.user.role !== 'admin') return res.status(403).json({ error: 'Нет прав' });
  const { title, description, price, category, location } = req.body;
  const image_path = req.file ? `products/${req.file.filename}` : undefined;
  const updates = [];
  const params = [];
  if (title !== undefined) { updates.push('title = ?'); params.push(title); }
  if (description !== undefined) { updates.push('description = ?'); params.push(description); }
  if (price !== undefined) { updates.push('price = ?'); params.push(price); }
  if (category !== undefined) { updates.push('category = ?'); params.push(category); }
  if (location !== undefined) { updates.push('location = ?'); params.push(location); }
  if (image_path !== undefined) { updates.push('image_path = ?'); params.push(image_path); }
  if (updates.length === 0) return res.json(db.prepare('SELECT * FROM products WHERE id = ?').get(req.params.id));
  updates.push('updated_at = unixepoch()');
  params.push(req.params.id);
  db.prepare(`UPDATE products SET ${updates.join(', ')} WHERE id = ?`).run(...params);
  res.json(db.prepare('SELECT * FROM products WHERE id = ?').get(req.params.id));
});

router.delete('/:id', authRequired, (req, res) => {
  const product = db.prepare('SELECT user_id FROM products WHERE id = ?').get(req.params.id);
  if (!product) return res.status(404).json({ error: 'Товар не найден' });
  if (product.user_id !== req.user.id && req.user.role !== 'admin') return res.status(403).json({ error: 'Нет прав' });
  db.prepare('DELETE FROM products WHERE id = ?').run(req.params.id);
  res.status(204).send();
});

export default router;
