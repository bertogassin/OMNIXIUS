-- Trade: payment_status and hold_id for orders (§15.3)
ALTER TABLE orders ADD COLUMN payment_status TEXT NOT NULL DEFAULT 'pending';
ALTER TABLE orders ADD COLUMN hold_id INTEGER;