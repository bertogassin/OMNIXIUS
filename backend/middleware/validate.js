import validator from 'validator';

const sanitize = (s) => (typeof s === 'string' ? validator.escape(s.trim()) : s);

export function validateRegister(req, res, next) {
  const email = sanitize(req.body?.email);
  const password = req.body?.password;
  const name = sanitize(req.body?.name);
  if (!email || !validator.isEmail(email)) return res.status(400).json({ error: 'Некорректный email' });
  if (!password || password.length < 8) return res.status(400).json({ error: 'Пароль не менее 8 символов' });
  req.body.email = email;
  req.body.password = password;
  req.body.name = name || null;
  next();
}

export function validateLogin(req, res, next) {
  const email = sanitize(req.body?.email);
  const password = req.body?.password;
  if (!email || !password) return res.status(400).json({ error: 'Email и пароль обязательны' });
  req.body.email = email;
  req.body.password = password;
  next();
}

export function validateProduct(req, res, next) {
  const title = sanitize(req.body?.title);
  const description = sanitize(req.body?.description);
  const price = Number(req.body?.price);
  const category = sanitize(req.body?.category);
  const location = sanitize(req.body?.location);
  if (!title || title.length < 2) return res.status(400).json({ error: 'Название не менее 2 символов' });
  if (isNaN(price) || price < 0) return res.status(400).json({ error: 'Некорректная цена' });
  if (!category) return res.status(400).json({ error: 'Укажите категорию' });
  req.body.title = title;
  req.body.description = description || '';
  req.body.price = price;
  req.body.category = category;
  req.body.location = location || null;
  next();
}

export function validateMessage(req, res, next) {
  const body = sanitize(req.body?.body);
  if (!body || body.length < 1) return res.status(400).json({ error: 'Текст сообщения обязателен' });
  req.body.body = body;
  next();
}
