-- A3 Gig: products can be marked as service (for "services nearby" filter).
ALTER TABLE products ADD COLUMN is_service INTEGER NOT NULL DEFAULT 0;
