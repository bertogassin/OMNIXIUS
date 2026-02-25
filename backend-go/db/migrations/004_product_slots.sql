-- A4 Gig: booking slots for service listings. Client books a slot → order + notification in Mail.
CREATE TABLE IF NOT EXISTS product_slots (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  product_id INTEGER NOT NULL REFERENCES products(id) ON DELETE CASCADE,
  slot_at INTEGER NOT NULL,
  status TEXT NOT NULL DEFAULT 'free' CHECK (status IN ('free', 'booked')),
  order_id INTEGER REFERENCES orders(id) ON DELETE SET NULL,
  created_at INTEGER DEFAULT (unixepoch())
);
CREATE INDEX IF NOT EXISTS idx_product_slots_product ON product_slots(product_id);
CREATE INDEX IF NOT EXISTS idx_product_slots_status ON product_slots(status);

ALTER TABLE orders ADD COLUMN slot_id INTEGER;
